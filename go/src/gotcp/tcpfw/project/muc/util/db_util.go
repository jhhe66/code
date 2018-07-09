// Copyright 2016 songliwei
//
// HelloTalk.inc

package util

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/gansidui/gotcp/tcpfw/include/ht_muc"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/protobuf/proto"
	// "github.com/seefan/gossdb"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Error type
var (
	ErrNilDbObject    = errors.New("not set  object current is nil")
	ErrDbParam        = errors.New("err param error")
	ErrNilNoSqlObject = errors.New("not set nosql object current is nil")
	ErrMemberQuit     = errors.New("member quit state is not 0")
)

var (
	RoomInfoDB        = "roominfo"
	RoomMemberTable   = "room_member"
	RoomMaxOrderTable = "room_max_order"
	pageSize          = 20
)

const (
	ENUM_QUIT_STATE_NORMAL       = 0
	ENUM_QUIT_STATE_SELF_QUIT    = 1
	ENUM_QUIT_STATE_ADMIN_REMOVE = 2
)

const (
	ENUM_NO_IN_CONTACT_LIST = 0 // 没有添加到联系列表中
	ENUM_IN_CONTACT_LIST    = 1 // 添加到联系列表中
)

type MemberInfoStore struct {
	Id          string `bson:"_id"` // 在Mongo中唯一标识一条记录
	RoomId      uint32
	Uid         uint32
	InviteUid   uint32
	NickName    string
	OrderId     uint32
	JoinTS      uint32
	PushSetting uint32
	ContactList uint32
	QuitStat    uint32
	UpdateTS    uint32
}

type MaxOrderInfo struct {
	RoomId   uint32 // 群聊Id
	MaxOrder uint32 // 当前最大加入序号
}
type DbUtil struct {
	db      *sql.DB
	infoLog *log.Logger
	noSqlDb *mgo.Session
}

func NewDbUtil(mysqlDb *sql.DB, mongoSess *mgo.Session, logger *log.Logger) *DbUtil {
	return &DbUtil{
		db:      mysqlDb,
		noSqlDb: mongoSess,
		infoLog: logger,
	}
}

func (this *DbUtil) RefreshSession(err error) {
	if err == io.EOF {
		this.noSqlDb.Refresh()
		this.infoLog.Printf("RefreshSession input err=%s", err)
	}
}

func (this *DbUtil) GetBlackMeList(uid uint32) (blackMeList []uint32, err error) {
	if this.db == nil {
		err = ErrNilDbObject
		return nil, err
	}

	if uid == 0 {
		err = ErrDbParam
		return
	}

	rows, err := this.db.Query("select USERID from HT_BLACK_LIST where BLACKID = ? and FLAG = 1", uid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var userId uint32
		if err := rows.Scan(&userId); err != nil {
			this.infoLog.Printf("GetBlackMeList rows.Scan failed")
			continue
		}
		blackMeList = append(blackMeList, userId)
	}
	return
}

func (this *DbUtil) GetUserVIPExpireTS(uid uint32) (expireTS uint64, err error) {
	if this.db == nil {
		err = ErrNilDbObject
		return
	}
	if uid == 0 {
		err = ErrDbParam
		return
	}

	err = this.db.QueryRow("select EXPIRETIME from HT_PURCHASE_TRANSLATE where USERID = ?", uid).Scan(&expireTS)
	switch {
	case err == sql.ErrNoRows:
		this.infoLog.Printf("GetUserVIPExpireTS not found uid=%v", uid)
		break
	case err != nil:
		this.infoLog.Println("GetUserVIPExpireTS exec failed [uid, err] =", uid, err)
		break
	default:
		this.infoLog.Printf("GetUserVIPExpireTS uid=%v vipts=%v", uid, expireTS)
	}
	return expireTS, err
}

func (this *DbUtil) CreateMucRoom(createUid uint32, listMember []uint32, memberLimit uint32) (roomId uint32, roomTS uint64, err error) {
	if this.db == nil {
		err = ErrNilDbObject
		return
	}

	if createUid == 0 || len(listMember) == 0 {
		this.infoLog.Printf("createUid=%u listmember empty", createUid)
		err = ErrDbParam
		return
	}

	// 首先在MySql中创建群
	roomTS = uint64(time.Now().Unix())
	r, err := this.db.Exec("insert into HT_MUC_ROOM (CREATEUSER,CREATETIME,MEMBERLIMIT,EMPTYROOM,ROOMTIMESTAMP,UPDATETIME) values (?, now(), ?, 0, ?, now())",
		createUid,
		memberLimit,
		roomTS)
	if err != nil {
		this.infoLog.Printf("CreateMucRoom insert faield createUid=%v err=%v", createUid, err)
		return
	}
	tempRoomId, err := r.LastInsertId()
	roomId = uint32(tempRoomId)
	if err != nil {
		this.infoLog.Printf("CreateMucRoom Get last insert id faied createUid=%v err=%v", createUid, err)
		return
	}
	this.infoLog.Println("CreateMucRoom succ roomId =", roomId)

	// 将群成员写入Mongo中 首先获取集合对象
	var orderId uint32
	mongoConn := this.noSqlDb.DB(RoomInfoDB).C(RoomMemberTable)
	selfInfo := &MemberInfoStore{
		Id:          bson.NewObjectId().String(),
		RoomId:      roomId,
		Uid:         createUid,
		InviteUid:   createUid,
		NickName:    "",
		OrderId:     orderId,
		JoinTS:      uint32(time.Now().Unix()),
		PushSetting: 0,
		ContactList: 0,
		QuitStat:    0,
		UpdateTS:    uint32(time.Now().Unix()),
	}
	err = mongoConn.Insert(selfInfo)
	// 如果RoomId+Uid 作为唯一所以 如果RoomId+Uid已经存在则插入失败
	if err != nil {
		this.RefreshSession(err)
		this.infoLog.Printf("CreateMucRoom Mongo Insert failed, roomId=%v createUid=%v err=%v",
			roomId,
			createUid,
			err)
		return
	}

	// 更新Mongo中的MaxOrder
	maxOrder := &MaxOrderInfo{RoomId: roomId, MaxOrder: orderId}
	maxOrderConn := this.noSqlDb.DB(RoomInfoDB).C(RoomMaxOrderTable)
	_, err = maxOrderConn.Upsert(bson.M{"roomid": maxOrder.RoomId}, maxOrder)
	if err != nil {
		this.RefreshSession(err)
		this.infoLog.Printf("CreateMucRoom Mongo upsert faield roomId=%v maxOrder=%v err=%v",
			roomId,
			orderId,
			err)
		return
	}

	//再将其它成员添加进去
	change := mgo.Change{
		Update:    bson.M{"$inc": bson.M{"maxorder": 1}},
		ReturnNew: true,
	}
	var maxOrderResult MaxOrderInfo
	for i, v := range listMember {
		this.infoLog.Printf("CreateMucRoom index=%v uid=%v", i, v)
		// 每次将max order 递增1
		_, err := maxOrderConn.Find(bson.M{"roomid": roomId}).Apply(change, &maxOrderResult)
		if err != nil {
			this.infoLog.Printf("CreateMucRoom roomId=%v memberId=%v inc maxorder failed err=%v",
				roomId,
				v,
				err)
			// 执行失败 orderID 跟后一次的保持一致
			orderId = orderId + 1
		} else {
			orderId = maxOrderResult.MaxOrder
		}

		this.infoLog.Printf("CreateMucRoom roomId=%v memberId=%v orderId=%v", roomId, v, orderId)
		strMemberInfo := &MemberInfoStore{
			Id:          bson.NewObjectId().String(),
			RoomId:      roomId,
			Uid:         v,
			InviteUid:   createUid,
			NickName:    "",
			OrderId:     orderId,
			JoinTS:      uint32(time.Now().Unix()),
			PushSetting: 0,
			ContactList: 0,
			QuitStat:    0,
			UpdateTS:    uint32(time.Now().Unix()),
		}

		err = mongoConn.Insert(strMemberInfo)
		if err != nil {
			this.infoLog.Printf("CreateMucRoom Mongo  Insert failed roomId=%v uid=%v err=%v",
				roomId,
				v,
				err)
			continue
		}
	}
	// 更新完毕无需再次更新maxOrder了
	return roomId, roomTS, nil
}

func (this *DbUtil) GetRoomMemberList(roomId uint32) (memberList []*MemberInfoStruct, maxOrderId int64, err error) {
	if this.noSqlDb == nil {
		this.infoLog.Printf("GetRoomMemberList noSqlDb object is nil")
		err = ErrNilNoSqlObject
		return
	}
	//从Mongodb 中获取所有的群成员
	mongoConn := this.noSqlDb.DB(RoomInfoDB).C(RoomMemberTable)
	query := mongoConn.Find(bson.M{"roomid": roomId, "quitstat": uint32(ENUM_QUIT_STATE_NORMAL)})
	query.Sort("orderid") // 按照order的升序排
	var result []MemberInfoStore
	query.All(&result)
	for k, v := range result {
		this.infoLog.Printf("GetRoomMemberList index=%v uid=%v", k, v.Uid)
		iterm := &MemberInfoStruct{
			Uid:         v.Uid,
			InvitedUid:  v.InviteUid,
			NickName:    v.NickName,
			OrderId:     v.OrderId,
			JoinTs:      v.JoinTS,
			PushSetting: v.PushSetting,
			RoomId:      v.RoomId,
		}
		memberList = append(memberList, iterm)
	}
	// 获取maxOrder 集合
	maxOrderConn := this.noSqlDb.DB(RoomInfoDB).C(RoomMaxOrderTable)
	var maxOrderInfo MaxOrderInfo
	err = maxOrderConn.Find(bson.M{"roomid": roomId}).One(&maxOrderInfo)
	if err != nil {
		this.RefreshSession(err)
		this.infoLog.Printf("GetRoomMemberList exec mgo One failed roomId=%v", roomId)
		return memberList, maxOrderId, err
	}
	maxOrderId = int64(maxOrderInfo.MaxOrder)

	return memberList, maxOrderId, err
}

func (this *DbUtil) GetRoomBaseInfo(roomId uint32) (roomInfo *RoomInfo, err error) {
	if this.db == nil {
		return nil, ErrNilDbObject
	}

	var storeRoomId, storeCreater, storeAdmin, storeAdmin2, storeAdmin3, storeAdminCount, storeVerifyStat sql.NullInt64
	var storePublishUid, storePublishTs sql.NullInt64
	var storePublishContent sql.NullString
	var storeAdmin4, storeAdmin5, storeAdmin6, storeAdmin7, storeAdmin8, storeAdmin9, storeAdmin10, storeMemberLimit sql.NullInt64
	var storeRoomName, storeRoomDesc sql.NullString
	var storeRoomTimeStamp uint64
	err = this.db.QueryRow("select ROOMID, CREATEUSER, ADMINUSER, ADMINUSER2, ADMINUSER3, ADMINCOUNT, INVITEVERIFY, PUBLISHUSER, PUBLISHTS, PUBLISHCONTENT, "+
		"ADMINUSER4, ADMINUSER5, ADMINUSER6, ADMINUSER7 ,ADMINUSER8, ADMINUSER9, ADMINUSER10, MEMBERLIMIT,ROOMNAME,DESCRIPTION,ROOMTIMESTAMP from HT_MUC_ROOM where ROOMID=?;",
		roomId).Scan(&storeRoomId, &storeCreater, &storeAdmin, &storeAdmin2, &storeAdmin3, &storeAdminCount, &storeVerifyStat,
		&storePublishUid, &storePublishTs, &storePublishContent,
		&storeAdmin4, &storeAdmin5, &storeAdmin6, &storeAdmin7, &storeAdmin8, &storeAdmin9, &storeAdmin10,
		&storeMemberLimit, &storeRoomName, &storeRoomDesc, &storeRoomTimeStamp)
	switch {
	case err == sql.ErrNoRows:
		this.infoLog.Printf("GetRoomBaseInfo not found roomId=%v", roomId)
		return nil, err
	case err != nil:
		this.infoLog.Println("GetRoomBaseInfo exec failed [roomId, err] =", roomId, err)
		return nil, err
	default:
		this.infoLog.Printf("GetRoomBaseInfo roomId=%v createuid=%v admin=%v storeAdminCount=%v memberlimit=%v roomName=%s desc=%s",
			roomId,
			storeCreater,
			storeAdmin,
			storeAdminCount,
			storeMemberLimit,
			storeRoomName,
			storeRoomDesc)
	}

	var creater, admin, admin2, admin3, adminCount, verifyStat uint32
	var publishUid, publishTs uint32
	var publishContent string
	var admin4, admin5, admin6, admin7, admin8, admin9, admin10, memberLimit uint32
	var roomName, roomDesc string

	if storeRoomId.Valid {
		roomId = uint32(storeRoomId.Int64)
	}

	if storeCreater.Valid {
		creater = uint32(storeCreater.Int64)
	}

	if storeAdmin.Valid {
		admin = uint32(storeAdmin.Int64)
	}

	if storeAdmin2.Valid {
		admin2 = uint32(storeAdmin2.Int64)
	}

	if storeAdmin3.Valid {
		admin3 = uint32(storeAdmin3.Int64)
	}

	if storeAdminCount.Valid {
		adminCount = uint32(storeAdminCount.Int64)
	}

	if storeVerifyStat.Valid {
		verifyStat = uint32(storeVerifyStat.Int64)
	}

	if storePublishUid.Valid {
		publishUid = uint32(storePublishUid.Int64)
	}

	if storePublishTs.Valid {
		publishTs = uint32(storePublishTs.Int64)
	}

	if storePublishContent.Valid {
		publishContent = storePublishContent.String
	}

	if storeAdmin4.Valid {
		admin4 = uint32(storeAdmin4.Int64)
	}

	if storeAdmin5.Valid {
		admin5 = uint32(storeAdmin5.Int64)
	}

	if storeAdmin6.Valid {
		admin6 = uint32(storeAdmin6.Int64)
	}

	if storeAdmin7.Valid {
		admin7 = uint32(storeAdmin7.Int64)
	}

	if storeAdmin8.Valid {
		admin8 = uint32(storeAdmin8.Int64)
	}

	if storeAdmin9.Valid {
		admin9 = uint32(storeAdmin9.Int64)
	}

	if storeAdmin10.Valid {
		admin10 = uint32(storeAdmin10.Int64)
	}

	if storeMemberLimit.Valid {
		memberLimit = uint32(storeMemberLimit.Int64)
	}

	if storeRoomName.Valid {
		roomName = storeRoomName.String
	}

	if storeRoomDesc.Valid {
		roomDesc = storeRoomDesc.String
	}

	roomInfo = &RoomInfo{
		RoomId:      roomId,
		CreateUid:   creater,
		AdminLimit:  adminCount,
		RoomName:    roomName,
		RoomDesc:    roomDesc,
		MemberLimit: memberLimit,
		VerifyStat:  verifyStat,
		Announcement: AnnouncementStruct{
			PublishUid:  publishUid,
			PublishTS:   publishTs,
			AnnoContect: publishContent,
		},
		RoomTS: int64(storeRoomTimeStamp),
	}

	if admin != 0 {
		roomInfo.AdminList = append(roomInfo.AdminList, admin)
	}

	if admin2 != 0 {
		roomInfo.AdminList = append(roomInfo.AdminList, admin2)
	}

	if admin3 != 0 {
		roomInfo.AdminList = append(roomInfo.AdminList, admin3)
	}

	if admin4 != 0 {
		roomInfo.AdminList = append(roomInfo.AdminList, admin4)
	}

	if admin5 != 0 {
		roomInfo.AdminList = append(roomInfo.AdminList, admin5)
	}

	if admin6 != 0 {
		roomInfo.AdminList = append(roomInfo.AdminList, admin6)
	}

	if admin7 != 0 {
		roomInfo.AdminList = append(roomInfo.AdminList, admin7)
	}

	if admin8 != 0 {
		roomInfo.AdminList = append(roomInfo.AdminList, admin8)
	}

	if admin9 != 0 {
		roomInfo.AdminList = append(roomInfo.AdminList, admin9)
	}

	if admin10 != 0 {
		roomInfo.AdminList = append(roomInfo.AdminList, admin10)
	}

	return roomInfo, nil
}

func (this *DbUtil) UpdateRoomMemberLimit(roomId, memberLimit uint32) (err error) {
	if this.db == nil {
		return ErrNilDbObject
	}

	if roomId == 0 || memberLimit == 0 {
		return ErrDbParam
	}

	_, err = this.db.Exec("update HT_MUC_ROOM set MEMBERLIMIT=?, ROOMTIMESTAMP=?, UPDATETIME=now() where ROOMID=?;",
		memberLimit,
		time.Now().Unix(),
		roomId)
	if err != nil {
		this.infoLog.Printf("UpdateRoomMemberLimit insert faield roomId=%v memberLimit=%v  err=%v", roomId, memberLimit, err)
		return err
	} else {
		return nil
	}
}

func (this *DbUtil) GetMemberAlreadyInMuc(roomId uint32, memberList []*ht_muc.RoomMemberInfo) (alreadyIn []uint32, err error) {
	if this.noSqlDb == nil {
		return nil, ErrNilNoSqlObject
	}
	if roomId == 0 || len(memberList) == 0 {
		return nil, ErrInputParam
	}
	uidlist := make([]uint32, len(memberList))
	for i, v := range memberList {
		uidlist[i] = v.GetUid()
	}
	//从Mongodb 中获取所有的群成员
	mongoConn := this.noSqlDb.DB(RoomInfoDB).C(RoomMemberTable)
	query := mongoConn.Find(bson.M{"roomid": roomId,
		"quitstat": uint32(ENUM_QUIT_STATE_NORMAL),
		"uid":      bson.M{"$in": uidlist}})

	query.Sort("orderid") // 按照order的升序排
	var result []MemberInfoStore
	query.All(&result)
	for k, v := range result {
		this.infoLog.Printf("index=%v uid=%s", k, v.Uid)
		memberUid := v.Uid
		alreadyIn = append(alreadyIn, memberUid)
	}
	return alreadyIn, nil
}

func (this *DbUtil) GetMemberAlreadyInMucInDb(roomId uint32, memberList []*ht_muc.RoomMemberInfo) (alreadyIn []uint32, err error) {
	if this.db == nil {
		return nil, ErrNilDbObject
	}
	if roomId == 0 || len(memberList) == 0 {
		return nil, ErrInputParam
	}
	strMemberList := "("
	for i, v := range memberList {
		if i == 0 {
			strMemberList += fmt.Sprintf("%v", v.GetUid())
		} else {
			strMemberList += "," + fmt.Sprintf("%v", v.GetUid())
		}
	}
	strMemberList = strMemberList + ")"
	this.infoLog.Printf("GetMemberAlreadyInMuc strMemberList=%s", strMemberList)
	rows, err := this.db.Query("select MEMBERID from HT_MUC_MEMBER where ROOMID=? AND QUITSTATE=0 AND MEMBERID in "+strMemberList+";", roomId)
	if err != nil {
		this.infoLog.Printf("GetMemberAlreadyInMuc roomId=%v strMemberList=%s failed\n", roomId, strMemberList)
		return alreadyIn, nil
	}
	defer rows.Close()
	for rows.Next() {
		var uid uint32
		if err := rows.Scan(&uid); err != nil {
			this.infoLog.Println("GetMemberAlreadyInMuc rows.Scan failed roomId=%v", roomId)
			continue
		}
		alreadyIn = append(alreadyIn, uid)
	}
	return alreadyIn, nil
}

func (this *DbUtil) InviteMember(roomId, inviteId uint32, memberList []*ht_muc.RoomMemberInfo) (roomTS int64, err error) {
	if this.noSqlDb == nil {
		err = ErrNilNoSqlObject
		return roomTS, err
	}
	if roomId == 0 || inviteId == 0 || len(memberList) == 0 {
		this.infoLog.Printf("InviteMember roomId=%v inviteId=%v memberCount=%v input error", roomId, inviteId, len(memberList))
		err = ErrDbParam
		return roomTS, err
	}
	roomTS, err = this.UpdateRoomTimeStamp(roomId)
	// 更新room的版本号失败直接返回
	if err != nil {
		this.infoLog.Printf("InviteMember UpdateRoomTimeStamp roomId=%v inviteId=%v memberCount=%v UpdateRoomTimeStamp failed",
			roomId,
			inviteId,
			len(memberList))
		return 0, err
	}
	// 获取群成员集合
	mongoConn := this.noSqlDb.DB(RoomInfoDB).C(RoomMemberTable)
	// 获取maxOrder 集合
	maxOrderConn := this.noSqlDb.DB(RoomInfoDB).C(RoomMaxOrderTable)

	//再将被邀请成员添加进去
	change := mgo.Change{
		Update:    bson.M{"$inc": bson.M{"maxorder": 1}},
		ReturnNew: true,
	}
	var maxOrderResult MaxOrderInfo
	var orderId uint32 = 0
	for i, v := range memberList {
		this.infoLog.Printf("InviteMember roomId=%v index=%v uid=%v", roomId, i, v.GetUid())
		// 每次将max order 递增1
		_, err := maxOrderConn.Find(bson.M{"roomid": roomId}).Apply(change, &maxOrderResult)
		if err != nil {
			this.RefreshSession(err)
			this.infoLog.Printf("InviteMember roomId=%v memberId=%v inc maxorder failed err=%v",
				roomId,
				v.GetUid(),
				err)
			// 执行失败 orderID 跟后一次的保持一致
			orderId = orderId + 1
		} else {
			orderId = maxOrderResult.MaxOrder
		}

		this.infoLog.Printf("InviteMember roomId=%v memberId=%v orderId=%v", roomId, v.GetUid(), orderId)
		strMemberInfo := &MemberInfoStore{
			Id:          bson.NewObjectId().String(),
			RoomId:      roomId,
			Uid:         v.GetUid(),
			InviteUid:   inviteId,
			NickName:    "",
			OrderId:     orderId,
			JoinTS:      uint32(time.Now().Unix()),
			PushSetting: 0,
			ContactList: 0,
			QuitStat:    0,
			UpdateTS:    uint32(time.Now().Unix()),
		}

		err = mongoConn.Insert(strMemberInfo)
		if err != nil {
			err = mongoConn.Update(bson.M{"roomid": roomId, "uid": v.GetUid()},
				bson.M{"$set": bson.M{
					"inviteuid":   inviteId,
					"nickname":    "",
					"orderid":     orderId,
					"joints":      uint32(time.Now().Unix()),
					"pushsetting": 0,
					"contactlist": 0,
					"quitstat":    0,
					"updateTS":    uint32(time.Now().Unix()),
				}})
			if err != nil {
				this.RefreshSession(err)
				this.infoLog.Printf("InviteMember Mongo  UpSert failed roomId=%v uid=%v err=%v",
					roomId,
					v.GetUid(),
					err)
				continue
			}
		}
	}
	// 更新完毕无需再次更新maxOrder了
	return roomTS, nil
}

func (this *DbUtil) UpdateRoomTimeStamp(roomId uint32) (roomTS int64, err error) {
	if this.db == nil {
		err = ErrNilDbObject
		return roomTS, err
	}

	if roomId == 0 {
		err = ErrDbParam
		return
	}
	roomTS = time.Now().Unix()
	_, err = this.db.Exec("update HT_MUC_ROOM set ROOMTIMESTAMP = ?, UPDATETIME = now() where ROOMID =?;", roomTS, roomId)
	if err != nil {
		this.infoLog.Printf("UpdateRoomTimeStamp failed roomId=%v err=%v", roomId, err)
		return roomTS, err
	} else {
		return roomTS, nil
	}
}

func (this *DbUtil) RemoveMember(roomId, removeId uint32) (err error) {
	if this.noSqlDb == nil {
		err = ErrNilNoSqlObject
		return err
	}
	if roomId == 0 || removeId == 0 {
		this.infoLog.Printf("RemoveMember roomId=%v removId=%v input error", roomId, removeId)
		err = ErrDbParam
		return err
	}

	// 获取群成员集合
	mongoConn := this.noSqlDb.DB(RoomInfoDB).C(RoomMemberTable)
	mongoConn.Update(bson.M{"roomid": roomId, "uid": removeId}, bson.M{"$set": bson.M{
		"quitstat": uint32(ENUM_QUIT_STATE_ADMIN_REMOVE),
		"updatets": uint32(time.Now().Unix())}})

	return nil
}

func (this *DbUtil) QuitMucRoom(roomId, quitUid uint32, bIsCreater bool, newCreateUid uint32) (roomTS int64, err error) {
	if this.noSqlDb == nil {
		err = ErrNilNoSqlObject
		return roomTS, err
	}

	if this.db == nil {
		err = ErrNilDbObject
		return
	}

	if roomId == 0 || quitUid == 0 {
		this.infoLog.Printf("QuitMucRoom roomId=%v quitUid=%v input error", roomId, quitUid)
		err = ErrDbParam
		return
	}

	// 用户主动退出
	mongoConn := this.noSqlDb.DB(RoomInfoDB).C(RoomMemberTable)
	mongoConn.Update(bson.M{"roomid": roomId, "uid": quitUid}, bson.M{"$set": bson.M{
		"quitstat": uint32(ENUM_QUIT_STATE_SELF_QUIT),
		"updatets": uint32(time.Now().Unix())}})

	// 更新完自己的状态之后如果是管理员退出群聊那么还需要设置新的创建者
	if bIsCreater {
		roomTS, err = this.UpdateRoomCreater(roomId, newCreateUid)
	} else {
		roomTS, err = this.UpdateRoomTimeStamp(roomId)
	}
	if err != nil {
		this.infoLog.Printf("QuitMucRoom update failed roomId=%v quitUid=%v err=%v", roomId, quitUid, err)
	}
	return roomTS, err
}

func (this *DbUtil) UpdateRoomCreater(roomId, createrUid uint32) (roomTS int64, err error) {
	if this.db == nil {
		err = ErrNilDbObject
		return
	}
	if roomId == 0 || createrUid == 0 {
		err = ErrDbParam
		return
	}
	roomTS = time.Now().Unix()
	_, err = this.db.Exec("update HT_MUC_ROOM set CREATEUSER = ?, ROOMTIMESTAMP = ?, UPDATETIME = now() where ROOMID =?;", createrUid, roomTS, roomId)
	if err != nil {
		this.infoLog.Printf("UpdateRoomCreater failed roomId=%v err=%v", roomId, err)
		return roomTS, err
	} else {
		return roomTS, nil
	}
}

func (this *DbUtil) ModifyRoomName(roomId, opUid uint32, roomName string) (roomTS int64, err error) {
	if this.db == nil {
		err = ErrNilDbObject
		return
	}

	if roomId == 0 || opUid == 0 || len(roomName) == 0 {
		this.infoLog.Printf("ModifyRoomName roomId=%v opUid=%v input error", roomId, opUid)
		err = ErrDbParam
		return
	}
	roomTS, err = this.UpdateRoomName(roomId, opUid, roomName)
	if err != nil {
		this.infoLog.Printf("ModifyRoomName exec UpdateRoomName fialed roomId=%v roomName=%s", roomId, roomName)
	}
	return roomTS, err
}

func (this *DbUtil) UpdateRoomName(roomId, opUid uint32, roomName string) (roomTS int64, err error) {
	if this.db == nil {
		err = ErrNilDbObject
		return
	}
	if roomId == 0 || opUid == 0 {
		err = ErrDbParam
		return
	}
	roomTS = time.Now().Unix()
	_, err = this.db.Exec("update HT_MUC_ROOM set ROOMNAME = ?, ROOMTIMESTAMP = ?, UPDATETIME = now() where ROOMID =?;", roomName, roomTS, roomId)
	if err != nil {
		this.infoLog.Printf("UpdateRoomName failed roomId=%v err=%v", roomId, err)
		return roomTS, err
	} else {
		return roomTS, nil
	}
}

func (this *DbUtil) ModifyMemberName(roomId, opUid uint32, opName string) (roomTS int64, err error) {
	if this.db == nil {
		err = ErrNilDbObject
		return
	}

	if this.noSqlDb == nil {
		err = ErrNilNoSqlObject
		return
	}

	if roomId == 0 || opUid == 0 || len(opName) == 0 {
		this.infoLog.Printf("ModifyMemberName roomId=%v opUid=%v input error", roomId, opUid)
		err = ErrDbParam
		return
	}

	mongoConn := this.noSqlDb.DB(RoomInfoDB).C(RoomMemberTable)
	mongoConn.Update(bson.M{"roomid": roomId, "uid": opUid}, bson.M{"$set": bson.M{
		"nickname": opName,
		"updatets": uint32(time.Now().Unix())}})

	// 更新成功之后更新群资料版本号
	roomTS, err = this.UpdateRoomTimeStamp(roomId)
	if err != nil {
		this.infoLog.Printf("ModifyMemberName exec UpdateRoomTimeStamp fialed roomId=%v opUid=%v opName=%s", roomId, opUid, opName)
	}
	return roomTS, err
}

func (this *DbUtil) ModifyPushSetting(roomId, opUid, pushSetting uint32) (err error) {
	if this.noSqlDb == nil {
		err = ErrNilNoSqlObject
		return err
	}

	if roomId == 0 || opUid == 0 || pushSetting > 1 {
		this.infoLog.Printf("ModifyPushSetting roomId=%v opUid=%v pushSetting=%v input error", roomId, opUid, pushSetting)
		err = ErrDbParam
		return err
	}
	mongoConn := this.noSqlDb.DB(RoomInfoDB).C(RoomMemberTable)
	mongoConn.Update(bson.M{"roomid": roomId, "uid": opUid}, bson.M{"$set": bson.M{
		"pushsetting": pushSetting,
		"updatets":    uint32(time.Now().Unix())}})
	return err
}

func (this *DbUtil) GetBlockRoomVoipUserList(roomId uint32) (outList []uint32, err error) {
	if this.db == nil {
		err = ErrNilDbObject
		return
	}
	if roomId == 0 {
		this.infoLog.Printf("GetBlockRoomVoipUserList roomId=%v input error", roomId)
		err = ErrDbParam
		return
	}
	//blockid 的类型，1=普通用户id，2=群id（HT_MUC_ROOM.ROOMID）
	//开关 1=开。0=默认关闭
	rows, err := this.db.Query("select USERID from HT_VOIPBLOCK_LIST where BLOCKID = ? and TYPE = 2 and FLAG = 1;", roomId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var userId uint32
		if err := rows.Scan(&userId); err != nil {
			this.infoLog.Printf("GetBlockRoomVoipUserList rows.Scan failed")
			continue
		}
		outList = append(outList, userId)
	}
	return outList, err
}

func (this *DbUtil) GetVoipRejectSetting(uid uint32) (bReject bool, err error) {
	if this.db == nil {
		err = ErrNilDbObject
		return
	}
	if uid == 0 {
		this.infoLog.Printf("GetVoipRejectSetting uid=%v input error", uid)
		err = ErrDbParam
		return
	}
	var setting uint8
	err = this.db.QueryRow("select VOIPREJECT from HT_USER_SETTING where USERID=? ;", uid).Scan(&setting)
	switch {
	case err == sql.ErrNoRows:
		this.infoLog.Printf("GetVoipRejectSetting not found uid=%v", uid)
		bReject = false
		return bReject, nil
	case err != nil:
		this.infoLog.Println("GetVoipRejectSetting exec failed [uid, err] =", uid, err)
		return false, err
	default:
		this.infoLog.Printf("GetVoipRejectSetting uid=%v voip setting=%v", uid, setting)
		if setting == 1 {
			bReject = true
		} else {
			bReject = false
		}
	}
	return bReject, err
}

func (this *DbUtil) AddRoomToContactList(roomId, opUid, opType uint32) (err error) {
	if this.noSqlDb == nil {
		err = ErrNilNoSqlObject
		return err
	}

	if roomId == 0 || opUid == 0 || opType > 1 {
		this.infoLog.Printf("AddRoomToContactList invalid param roomId=%v opUid=%v opType=%v",
			roomId,
			opUid,
			opType)
		err = ErrDbParam
		return err
	}
	mongoConn := this.noSqlDb.DB(RoomInfoDB).C(RoomMemberTable)
	mongoConn.Update(bson.M{"roomid": roomId, "uid": opUid}, bson.M{"$set": bson.M{
		"contactlist": opType,
		"updatets":    uint32(time.Now().Unix())}})
	return nil
}

func (this *DbUtil) GetAllContactListRoomId(opUid uint32) (roomIdList []uint32, err error) {
	if this.noSqlDb == nil {
		err = ErrNilNoSqlObject
		return
	}

	if opUid == 0 {
		this.infoLog.Printf("GetAllContactListRoomId invalid param opUid=%v",
			opUid)
		err = ErrDbParam
		return
	}

	var result []struct{ RoomId uint32 }
	mongoConn := this.noSqlDb.DB(RoomInfoDB).C(RoomMemberTable)
	err = mongoConn.Find(bson.M{"uid": opUid,
		"contactlist": uint32(ENUM_IN_CONTACT_LIST),
		"quitstat":    uint32(ENUM_QUIT_STATE_NORMAL)}).Select(bson.M{"roomid": 1}).All(&result)
	if err != nil {
		this.RefreshSession(err)
		this.infoLog.Printf("GetAllContactListRoomId opUid=%v exec mongoConn find err=%v", opUid, err)
		return nil, err
	}
	for _, v := range result {
		roomIdList = append(roomIdList, v.RoomId)
	}
	return roomIdList, nil
}

func (this *DbUtil) UpdateVoipBlockList(opUid, blockId, blockType, action uint32) (err error) {
	if this.db == nil {
		err = ErrNilDbObject
		return err
	}

	if opUid == 0 || blockId == 0 {
		this.infoLog.Printf("UpdateVoipBlockList opUid=%v blockId=%v blockType=%v action=%v",
			opUid,
			blockId,
			blockType,
			action)
		err = ErrDbParam
		return err
	}

	_, err = this.db.Exec("insert into HT_VOIPBLOCK_LIST (USERID, BLOCKID, TYPE, FLAG, UPDATETIME) VALUES (?,?,?,?,now()) on duplicate key update FLAG=VALUES(FLAG),UPDATETIME=VALUES(UPDATETIME);",
		opUid,
		blockId,
		blockType,
		action)
	if err != nil {
		this.infoLog.Printf("UpdateVoipBlockList insert faield opUid=%v blockId=%v blockType=%v falg=%v err=%v",
			opUid,
			blockId,
			blockType,
			action,
			err)
		return err
	}
	err = nil
	this.infoLog.Printf("UpdateVoipBlockList insert opUid=%v blockId=%v blockType=%v falg=%v success",
		opUid,
		blockId,
		blockType,
		action)
	return err
}

func (this *DbUtil) UpdateVerifyStat(roomId, verifyStat uint32) (roomTS int64, err error) {
	if this.db == nil {
		err = ErrNilDbObject
		return
	}
	if roomId == 0 {
		this.infoLog.Printf("UpdateVerifyStat roomId=%v verifyStat=%v",
			roomId,
			verifyStat)
		err = ErrDbParam
		return
	}

	roomTS = time.Now().Unix()
	_, err = this.db.Exec("update HT_MUC_ROOM set INVITEVERIFY = ?, ROOMTIMESTAMP = ?, UPDATETIME = now() where ROOMID =?;", verifyStat, roomTS, roomId)
	if err != nil {
		this.infoLog.Printf("UpdateVerifyStat failed roomId=%v err=%v", roomId, err)
		return roomTS, err
	} else {
		return roomTS, nil
	}
}

func (this *DbUtil) SetAdminListWithUint32(roomId uint32, adminList []uint32) (roomTS int64, err error) {
	var newAdminList []*ht_muc.RoomMemberInfo
	for _, v := range adminList {
		newAdminList = append(newAdminList, &ht_muc.RoomMemberInfo{Uid: proto.Uint32(v)})
	}
	roomTS, err = this.SetAdminList(roomId, newAdminList)
	return roomTS, err
}

func (this *DbUtil) SetAdminList(roomId uint32, adminList []*ht_muc.RoomMemberInfo) (roomTS int64, err error) {
	if this.db == nil {
		err = ErrNilDbObject
		return
	}

	if roomId == 0 {
		err = ErrNilDbObject
		return
	}
	roomTS = time.Now().Unix()
	switch len(adminList) {
	case 0:
		_, err = this.db.Exec("update HT_MUC_ROOM set ADMINUSER=0, ADMINUSER2=0, ADMINUSER3=0, ADMINUSER4=0, ADMINUSER5=0, ADMINUSER6=0, ADMINUSER7=0, ADMINUSER8=0, ADMINUSER9=0, ADMINUSER10=0, ROOMTIMESTAMP = ?, UPDATETIME = now() where ROOMID =?;", roomTS, roomId)
	case 1:
		_, err = this.db.Exec("update HT_MUC_ROOM set ADMINUSER=?, ADMINUSER2=0, ADMINUSER3=0, ADMINUSER4=0, ADMINUSER5=0, ADMINUSER6=0, ADMINUSER7=0, ADMINUSER8=0, ADMINUSER9=0, ADMINUSER10=0, ROOMTIMESTAMP = ?, UPDATETIME = now() where ROOMID =?;",
			adminList[0].GetUid(),
			roomTS,
			roomId)
	case 2:
		_, err = this.db.Exec("update HT_MUC_ROOM set ADMINUSER=?, ADMINUSER2=?, ADMINUSER3=0, ADMINUSER4=0, ADMINUSER5=0, ADMINUSER6=0, ADMINUSER7=0, ADMINUSER8=0, ADMINUSER9=0, ADMINUSER10=0, ROOMTIMESTAMP = ?, UPDATETIME = now() where ROOMID =?;",
			adminList[0].GetUid(),
			adminList[1].GetUid(),
			roomTS,
			roomId)
	case 3:
		_, err = this.db.Exec("update HT_MUC_ROOM set ADMINUSER=?, ADMINUSER2=?, ADMINUSER3=?, ADMINUSER4=0, ADMINUSER5=0, ADMINUSER6=0, ADMINUSER7=0, ADMINUSER8=0, ADMINUSER9=0, ADMINUSER10=0, ROOMTIMESTAMP = ?, UPDATETIME = now() where ROOMID =?;",
			adminList[0].GetUid(),
			adminList[1].GetUid(),
			adminList[2].GetUid(),
			roomTS,
			roomId)
	case 4:
		_, err = this.db.Exec("update HT_MUC_ROOM set ADMINUSER=?, ADMINUSER2=?, ADMINUSER3=?, ADMINUSER4=?, ADMINUSER5=0, ADMINUSER6=0, ADMINUSER7=0, ADMINUSER8=0, ADMINUSER9=0, ADMINUSER10=0, ROOMTIMESTAMP = ?, UPDATETIME = now() where ROOMID =?;",
			adminList[0].GetUid(),
			adminList[1].GetUid(),
			adminList[2].GetUid(),
			adminList[3].GetUid(),
			roomTS,
			roomId)
	case 5:
		_, err = this.db.Exec("update HT_MUC_ROOM set ADMINUSER=?, ADMINUSER2=?, ADMINUSER3=?, ADMINUSER4=?, ADMINUSER5=?, ADMINUSER6=0, ADMINUSER7=0, ADMINUSER8=0, ADMINUSER9=0, ADMINUSER10=0, ROOMTIMESTAMP = ?, UPDATETIME = now() where ROOMID =?;",
			adminList[0].GetUid(),
			adminList[1].GetUid(),
			adminList[2].GetUid(),
			adminList[3].GetUid(),
			adminList[4].GetUid(),
			roomTS,
			roomId)
	case 6:
		_, err = this.db.Exec("update HT_MUC_ROOM set ADMINUSER=?, ADMINUSER2=?, ADMINUSER3=?, ADMINUSER4=?, ADMINUSER5=?, ADMINUSER6=?, ADMINUSER7=0, ADMINUSER8=0, ADMINUSER9=0, ADMINUSER10=0, ROOMTIMESTAMP = ?, UPDATETIME = now() where ROOMID =?;",
			adminList[0].GetUid(),
			adminList[1].GetUid(),
			adminList[2].GetUid(),
			adminList[3].GetUid(),
			adminList[4].GetUid(),
			adminList[5].GetUid(),
			roomTS,
			roomId)
	case 7:
		_, err = this.db.Exec("update HT_MUC_ROOM set ADMINUSER=?, ADMINUSER2=?, ADMINUSER3=?, ADMINUSER4=?, ADMINUSER5=?, ADMINUSER6=?, ADMINUSER7=?, ADMINUSER8=0, ADMINUSER9=0, ADMINUSER10=0, ROOMTIMESTAMP = ?, UPDATETIME = now() where ROOMID =?;",
			adminList[0].GetUid(),
			adminList[1].GetUid(),
			adminList[2].GetUid(),
			adminList[3].GetUid(),
			adminList[4].GetUid(),
			adminList[5].GetUid(),
			adminList[6].GetUid(),
			roomTS,
			roomId)
	case 8:
		_, err = this.db.Exec("update HT_MUC_ROOM set ADMINUSER=?, ADMINUSER2=?, ADMINUSER3=?, ADMINUSER4=?, ADMINUSER5=?, ADMINUSER6=?, ADMINUSER7=?, ADMINUSER8=?, ADMINUSER9=0, ADMINUSER10=0, ROOMTIMESTAMP = ?, UPDATETIME = now() where ROOMID =?;",
			adminList[0].GetUid(),
			adminList[1].GetUid(),
			adminList[2].GetUid(),
			adminList[3].GetUid(),
			adminList[4].GetUid(),
			adminList[5].GetUid(),
			adminList[6].GetUid(),
			adminList[7].GetUid(),
			roomTS,
			roomId)
	case 9:
		_, err = this.db.Exec("update HT_MUC_ROOM set ADMINUSER=?, ADMINUSER2=?, ADMINUSER3=?, ADMINUSER4=?, ADMINUSER5=?, ADMINUSER6=?, ADMINUSER7=?, ADMINUSER8=?, ADMINUSER9=?, ADMINUSER10=0, ROOMTIMESTAMP = ?, UPDATETIME = now() where ROOMID =?;",
			adminList[0].GetUid(),
			adminList[1].GetUid(),
			adminList[2].GetUid(),
			adminList[3].GetUid(),
			adminList[4].GetUid(),
			adminList[5].GetUid(),
			adminList[6].GetUid(),
			adminList[7].GetUid(),
			adminList[8].GetUid(),
			roomTS,
			roomId)
	case 10:
		_, err = this.db.Exec("update HT_MUC_ROOM set ADMINUSER=?, ADMINUSER2=?, ADMINUSER3=?, ADMINUSER4=?, ADMINUSER5=?, ADMINUSER6=?, ADMINUSER7=?, ADMINUSER8=?, ADMINUSER9=?, ADMINUSER10=?, ROOMTIMESTAMP = ?, UPDATETIME = now() where ROOMID =?;",
			adminList[0].GetUid(),
			adminList[1].GetUid(),
			adminList[2].GetUid(),
			adminList[3].GetUid(),
			adminList[4].GetUid(),
			adminList[5].GetUid(),
			adminList[6].GetUid(),
			adminList[7].GetUid(),
			adminList[8].GetUid(),
			adminList[9].GetUid(),
			roomTS,
			roomId)
	}
	if err != nil {
		this.infoLog.Printf("SetAdminList failed roomId=%v err=%v", roomId, err)
		return roomTS, err
	} else {
		return roomTS, nil
	}

}

func (this *DbUtil) UpdateCreateUid(roomId, opUid, targetUid uint32) (roomTS int64, err error) {
	if this.db == nil {
		err = ErrNilDbObject
		return
	}
	if roomId == 0 || opUid == 0 || targetUid == 0 {
		this.infoLog.Printf("UpdateCreateUid roomId=%v opUid=%v targetUid=%v",
			roomId,
			opUid,
			targetUid)
		err = ErrDbParam
		return
	}

	roomTS = time.Now().Unix()
	_, err = this.db.Exec("update HT_MUC_ROOM set CREATEUSER = ?, ROOMTIMESTAMP = ?, UPDATETIME = now() where ROOMID =?;", targetUid, roomTS, roomId)
	if err != nil {
		this.infoLog.Printf("UpdateCreateUid failed roomId=%v err=%v", roomId, err)
		return roomTS, err
	} else {
		return roomTS, nil
	}
}

func (this *DbUtil) UpdateAnnouncement(roomId, publishUid, publishTs uint32, content string) (roomTS int64, err error) {
	if this.db == nil {
		err = ErrNilDbObject
		return
	}
	if roomId == 0 || publishUid == 0 {
		this.infoLog.Printf("UpdateAnnouncement roomId=%v publishUid=%v publishTs=%v content=%s",
			roomId,
			publishUid,
			publishTs,
			content)
		err = ErrDbParam
		return
	}

	roomTS = time.Now().Unix()
	_, err = this.db.Exec("update HT_MUC_ROOM set PUBLISHUSER = ?, PUBLISHTS = ?, PUBLISHCONTENT = ?, ROOMTIMESTAMP = ?, UPDATETIME = now() where ROOMID =?;",
		publishUid,
		publishTs,
		content,
		roomTS,
		roomId)
	if err != nil {
		this.infoLog.Printf("UpdateAnnouncement failed roomId=%v err=%v", roomId, err)
		return roomTS, err
	} else {
		return roomTS, nil
	}
}
