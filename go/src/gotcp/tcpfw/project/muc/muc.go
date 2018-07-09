package main

import (
	"bytes"
	"compress/zlib"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"sync"
	"unicode/utf8"

	simplejson "github.com/bitly/go-simplejson"
	"github.com/gansidui/gotcp"

	"github.com/gansidui/gotcp/libcomm"
	"github.com/gansidui/gotcp/tcpfw/common"
	"github.com/gansidui/gotcp/tcpfw/include/ht_muc"
	"github.com/gansidui/gotcp/tcpfw/project/muc/util"
	mgo "gopkg.in/mgo.v2"

	"log"
	"net"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/protobuf/proto"
	"github.com/jessevdk/go-flags"
	"gopkg.in/ini.v1"
)

type Callback struct{}

var (
	infoLog *log.Logger
	db      *sql.DB
	//ssdb        *gossdb.Connectors
	mongoSess   *mgo.Session
	mcApi       *common.MemcacheApi
	offlineApi  *common.OfflineApiV2
	DbUtil      *util.DbUtil
	roomManager *util.RoomManager
	ticker      *time.Ticker
	reqRecord   map[string]int64
	reqLock     sync.Mutex // sync mutex goroutines use reqRecord
)

const (
	DB_RET_SUCCESS     = 0
	DB_RET_EXEC_FAILED = 1
	DB_RET_NOT_EXIST   = 100
)

const (
	MUC_MESSAGE_LNE_LIMIT = 30
)

// 消息类型
const (
	MT_TEXT       = 0  // 文字
	MT_TRANSLATE  = 1  // 带翻译的文字
	MT_IMAGE      = 2  // 图片
	MT_VOICE      = 3  // 音频
	MT_LOCATE     = 4  // 定位
	MT_PROFILE    = 5  // 用户介绍
	MT_VOICETEXT  = 6  // 语音转文字
	MT_CORRECTION = 7  // 修改句子
	MT_STICKERS   = 8  // 自定义表情
	MT_DOODLE     = 9  // 涂鸦
	MT_VOIP       = 10 // VOIP？
	MT_NOTIFY     = 11 // 公众消息？
	MT_VIDEO      = 12 // 视频 video
	MT_GVOIP      = 13 // group voip
	MT_LINK       = 14
	MT_CARD       = 15
	MT_UNKNOWN    = 100
)

// 推送类型  需要将消息类型转到推送的类型上
const (
	PUSH_TEXT                    = 0
	PUSH_VOICE                   = 1
	PUSH_IMAGE                   = 2
	PUSH_INTRODUCE               = 3
	PUSH_LOCATION                = 4
	PUSH_FRIEND_INVITE           = 5
	PUSH_LANGUAGE_EXCHANGE       = 6
	PUSH_LANGUAGE_EXCHANGE_REPLY = 7
	PUSH_CORRECT_SENTENCE        = 8
	PUSH_STICKERS                = 9
	PUSH_DOODLE                  = 10
	PUSH_GIFT                    = 11
	PUSH_VOIP                    = 12
	PUSH_ACCEPT_INVITE           = 13
	PUSH_VIDEO                   = 14
	PUSH_GVOIP                   = 15 // group voip
	PUSH_LINK                    = 16
	PUSH_CARD                    = 17
	PUSH_FOLLOW                  = 18
	PUSH_REPLY_YOUR_COMMENT      = 19
	PUSH_COMMENTED_YOUR_POST     = 20
	PUSH_CORRECTED_YOUR_POST     = 21
	PUSH_MOMENT_LIKE             = 22
)

const (
	REQTHRESHOLD = 300
)

func GetMessageType(strMsgType string) uint8 {
	if 0 == strings.Compare(strMsgType, "text") {
		return MT_TEXT
	} else if 0 == strings.Compare(strMsgType, "translate") {
		return MT_TRANSLATE
	} else if 0 == strings.Compare(strMsgType, "voice") {
		return MT_VOICE
	} else if 0 == strings.Compare(strMsgType, "image") {
		return MT_IMAGE
	} else if 0 == strings.Compare(strMsgType, "introduction") {
		return MT_PROFILE
	} else if 0 == strings.Compare(strMsgType, "location") {
		return MT_LOCATE
	} else if 0 == strings.Compare(strMsgType, "voice_text") {
		return MT_VOICETEXT
	} else if 0 == strings.Compare(strMsgType, "correction") {
		return MT_CORRECTION
	} else if (0 == strings.Compare(strMsgType, "sticker")) ||
		(0 == strings.Compare(strMsgType, "new_sticker")) {
		return MT_STICKERS
	} else if 0 == strings.Compare(strMsgType, "doodle") {
		return MT_DOODLE
	} else if 0 == strings.Compare(strMsgType, "video") {
		return MT_VIDEO
	} else if 0 == strings.Compare(strMsgType, "link") {
		return MT_LINK
	} else if 0 == strings.Compare(strMsgType, "card") {
		return MT_CARD
	}

	return MT_UNKNOWN
}

func SubString(str string, begin, length int) (substr string) {
	// 将字符串的转换成[]rune
	rs := []rune(str)
	lth := len(rs)

	// 简单的越界判断
	if begin < 0 {
		begin = 0
	}
	if begin >= lth {
		begin = lth
	}
	end := begin + length
	if end > lth {
		end = lth
	}

	// 返回子串
	return string(rs[begin:end])
}

func (this *Callback) OnConnect(c *gotcp.Conn) bool {
	addr := c.GetRawConn().RemoteAddr()
	c.PutExtraData(addr)
	infoLog.Println("OnConnect:", addr)
	return true
}

func (this *Callback) OnMessage(c *gotcp.Conn, p gotcp.Packet) bool {
	packet, ok := p.(*common.HeadV3Packet)
	if !ok { // 不是HeadV3Packet报文
		infoLog.Printf("OnMessage packet can not change to HeadV3packet")
		return false
	}

	// head 为一个new出来的对象指针
	head, err := packet.GetHead()
	if err != nil {
		//SendResp(c, head, uint16(ERR_INVALID_PARAM))
		infoLog.Printf("OnMessage Get head failed", err)
		return false
	}

	infoLog.Printf("OnMessage:[head:%#v] bodyLen=%v \n", head, len(packet.GetBody()))
	//	infoLog.Printf("OnMessage:[%#v] len=%v\n", head, len(packet.GetBody()))
	_, err = packet.CheckPacketValid()
	if err != nil {
		infoLog.Printf("OnMessage Invalid packet", err)
		return false
	}

	// reqBody 也为一个指针
	reqBody := &ht_muc.MucReqBody{}
	err = proto.Unmarshal(packet.GetBody(), reqBody)
	if err != nil {
		infoLog.Printf("OnMessage proto Unmarshal failed")
		return false
	}
	// 统计总的请求量
	attr := "gomuc/total_recv_req_count"
	libcomm.AttrAdd(attr, 1)

	switch ht_muc.MUC_CMD_TYPE(head.Cmd) {
	case ht_muc.MUC_CMD_TYPE_GO_CMD_MUC_CREATE_ROOM:
		go ProcCreateRoom(c, head, reqBody)
	case ht_muc.MUC_CMD_TYPE_GO_CMD_MUC_REMOVE_MEMBER:
		go ProcRemoveMember(c, head, reqBody)
	case ht_muc.MUC_CMD_TYPE_GO_CMD_MUC_USER_QUIT:
		go ProcQuitRoom(c, head, reqBody)
	case ht_muc.MUC_CMD_TYPE_GO_CMD_MUC_MODIFY_ROOMNAME:
		go ProcModifyRoomName(c, head, reqBody)
	case ht_muc.MUC_CMD_TYPE_GO_CMD_MODIFY_MEMBER_NAME:
		go ProcModifyMemberName(c, head, reqBody)
	case ht_muc.MUC_CMD_TYPE_GO_CMD_MODIFY_MUC_PUSH_SETTING:
		go ProcModifyPushSetting(c, head, reqBody)
	case ht_muc.MUC_CMD_TYPE_GO_CMD_GET_MUC_ROOM_INFO:
		go ProcGetRoomInfo(c, head, reqBody)
	case ht_muc.MUC_CMD_TYPE_GO_CMD_MUC_MESSAGE:
		go ProcMucMessage(c, head, reqBody)
	case ht_muc.MUC_CMD_TYPE_GO_CMD_ADD_MUC_TO_CONTACT_LIST:
		go ProcAddRoomToContactList(c, head, reqBody)
	case ht_muc.MUC_CMD_TYPE_GO_CMD_GET_MUC_CONTACT_LIST:
		go ProcGetRoomFromContactList(c, head, reqBody)
	case ht_muc.MUC_CMD_TYPE_GO_CMD_VOIP_BLOCK_SETTING:
		go ProcS2SVoipBlockSetting(c, head, reqBody)
	case ht_muc.MUC_CMD_TYPE_GO_CMD_GVOIP_INVITE_BROADCAST:
		go ProcS2SGvoipInviteBroadCast(c, head, reqBody)
	case ht_muc.MUC_CMD_TYPE_GO_CMD_GVOIP_END_BROADCAST:
		go ProcS2SGvoipEndBroadCast(c, head, reqBody)
	case ht_muc.MUC_CMD_TYPE_GO_CMD_GVOIP_MEMBER_JOIN_BROADCAST, ht_muc.MUC_CMD_TYPE_GO_CMD_GVOIP_MEMBER_LEAVE_BROADCAST:
		go ProcS2SGvoipMemberChangeBroadCast(c, head, reqBody)
	case ht_muc.MUC_CMD_TYPE_GO_CMD_MUC_REQ_JOIN_ROOM:
		go ProcRequestJoinRoom(c, head, reqBody)
	case ht_muc.MUC_CMD_TYPE_GO_CMD_MUC_OPEN_REQ_VERIFY:
		go ProcMucOpenVerify(c, head, reqBody)
	case ht_muc.MUC_CMD_TYPE_GO_CMD_MUC_REQ_JOIN_ROOM_HANDLE:
		go ProcMucJoinRoomHandle(c, head, reqBody)
	case ht_muc.MUC_CMD_TYPE_GO_CMD_MUC_SET_ADMIN_REQ:
		go ProcMucSetAdmin(c, head, reqBody)
	case ht_muc.MUC_CMD_TYPE_GO_CMD_CREATEUSER_AUTHORIZATION_TRANS:
		go ProcCreateUserAuthTrans(c, head, reqBody)
	case ht_muc.MUC_CMD_TYPE_GO_CMD_MUC_GET_ROOM_BASE_INFO:
		go ProcMucGetRoomBaseInfo(c, head, reqBody)
	case ht_muc.MUC_CMD_TYPE_GO_CMD_MUC_SET_ROOM_ANNOUNCEMENT:
		go ProcMucSetRoomAnnouncement(c, head, reqBody)
	case ht_muc.MUC_CMD_TYPE_GO_CMD_QUREY_USR_IS_IN_ROOM:
		go ProcQueryUserIsAlreadyInRoom(c, head, reqBody)
	case ht_muc.MUC_CMD_TYPE_GO_CMD_GET_MUC_GET_QRCODE_INFO:
		go ProcQueryQRcodeInfo(c, head, reqBody)
	default:
		infoLog.Printf("OnMessage UnHandle Cmd =", head.Cmd)
	}
	return true
}

func IsDuplicateSec(roomId, uid uint32, seq uint16) (result bool) {
	seqKey := fmt.Sprintf("%v-%v-%v", roomId, uid, seq)
	reqLock.Lock()
	defer reqLock.Unlock()
	if _, ok := reqRecord[seqKey]; ok {
		result = true
	} else {
		result = false
		// 如果不存在 这添加元素
		reqRecord[seqKey] = time.Now().Unix()
	}
	return result
}

// 1.proc create room
func ProcCreateRoom(c *gotcp.Conn, head *common.HeadV3, reqBody *ht_muc.MucReqBody) bool {
	result := uint16(ht_muc.MUC_RET_CODE_RET_SUCCESS)
	errMsg := []byte("operation success")
	bNeedCall := true
	defer func() {
		if bNeedCall {
			infoLog.Printf("ProcCreateRoom defer SendCreatRoomRetCode")
			SendCreatRoomRetCode(c, head, result, errMsg)
		} else {
			infoLog.Printf("ProcCreateRoom defer not SendCreatRoomRetCode")
		}
	}()
	// 检查输入参数是否为空
	if head == nil || reqBody == nil {
		bNeedCall = false
		return false
	}

	// test is duplicate req
	bRet := IsDuplicateSec(0, head.From, head.Seq)
	if bRet {
		infoLog.Printf("ProcCreateRoom duplicate req uid=%v seq=%v", head.From, head.Seq)
		result = uint16(ht_muc.MUC_RET_CODE_RET_INPUT_PARAM_ERR)
		errMsg = []byte("input param error")
		return false
	}
	// add static
	attr := "gomuc/create_room_req_count"
	libcomm.AttrAdd(attr, 1)

	subReqBody := reqBody.GetCreateRoomReqbody()
	creatUid := subReqBody.GetCreateUid()
	nickName := subReqBody.GetNickName()
	memberList := subReqBody.GetMembers()
	infoLog.Printf("ProcCreateRoom recv from=%v to=%v cmd=%v seq=%v userid=%v count=%v",
		head.From,
		head.To,
		head.Cmd,
		head.Seq,
		creatUid,
		len(memberList))
	// first check input param
	if creatUid == 0 || len(memberList) < 2 {
		infoLog.Printf("ProcCreateRoom invalid param uid=%v count=%v", creatUid, len(memberList))
		result = uint16(ht_muc.MUC_RET_CODE_RET_INPUT_PARAM_ERR)
		errMsg = []byte("input param error")
		return false
	}
	// all client support muc chat no need to check version
	// second check black me list
	blackMeList, err := DbUtil.GetBlackMeList(creatUid)
	if err != nil {
		infoLog.Printf("ProcCreateRoom GetBlackMeList failed uid=%v", creatUid)
	}
	var totalBalck []uint32
	if len(blackMeList) != 0 {
		for i, v := range memberList {
			infoLog.Printf("ProcCreateRoom invite member index=%v uid=%v", i, v)
			if UidIsInSlice(blackMeList, v.GetUid()) {
				infoLog.Printf("ProcCreateRoom uidInfo=%#v is in blackMeList", v)
				totalBalck = append(totalBalck, v.GetUid())
			}
		}
	}
	if len(totalBalck) != 0 {
		subMucRspBody := new(ht_muc.CreateRoomRspBody)
		errMsg = []byte("some one black me")
		subMucRspBody.Status = &ht_muc.MucHeader{Code: proto.Uint32(uint32(ht_muc.MUC_RET_CODE_RET_SOME_ONE_BLACK_ME)), Reason: errMsg}
		subMucRspBody.ListBlackMe = totalBalck
		bNeedCall = false // 无需调用defer 中的函数
		SendCreatRoomResp(c, head, subMucRspBody)
		return false
	}

	// third get vip expire ts from db
	vipExpireTs, err := DbUtil.GetUserVIPExpireTS(creatUid)
	if err != nil {
		infoLog.Printf("ProcCreateRoom GetUserVIPExpireTS failed uid=%v", creatUid)
	}
	tsNow := time.Now().Unix()
	maxMember := util.MUC_MEMBER_LIMIT
	if vipExpireTs > uint64(tsNow) {
		maxMember = util.MUC_MEMBER_LIMIT_VIP
	}
	if len(memberList) > maxMember {
		infoLog.Printf("ProcCreateRoom memberList count=%v excess maxMember=%v", len(memberList), maxMember)
		result = uint16(ht_muc.MUC_RET_CODE_RET_MEMBER_EXEC_LIMIT)
		errMsg = []byte("member over limit")
		return false
	}

	// fourth crate muc room
	var uidList []uint32
	for _, v := range memberList {
		uidList = append(uidList, v.GetUid())
	}
	roomId, roomTS, err := roomManager.CreateMuc(creatUid, uidList, uint32(maxMember))
	if err != nil {
		infoLog.Printf("ProcCreateRoom room manager create room failed create=%v", creatUid)
		result = uint16(ht_muc.MUC_RET_CODE_RET_INTERNAL_ERR)
		errMsg = []byte("internal error")
		return false
	}
	// fiveth broadcast message
	opInfo := &ht_muc.RoomMemberInfo{
		Uid:      proto.Uint32(creatUid),
		NickName: nickName,
	}
	err = roomManager.NotifyInviteMember(opInfo, opInfo, memberList, roomId, uint64(roomTS), uint32(ht_muc.ROOMID_FROM_TYPE_ENUM_FROM_INVITE))
	if err != nil {
		infoLog.Printf("ProcCreateRoom NotifyInviteMember failed uid=%v roomId=%v roomTS=%v", creatUid, roomId, roomTS)
	}

	infoLog.Printf("DEBUG ProcCreateRoom roomId=%v roomTS=%v", roomId, roomTS)
	// sixth send ack
	subMucRspBody := new(ht_muc.CreateRoomRspBody)
	subMucRspBody.Status = &ht_muc.MucHeader{Code: proto.Uint32(uint32(ht_muc.MUC_RET_CODE_RET_SUCCESS))}
	subMucRspBody.RoomId = proto.Uint32(roomId)
	subMucRspBody.RoomTimestamp = proto.Uint64(uint64(roomTS))
	SendCreatRoomResp(c, head, subMucRspBody)
	bNeedCall = false
	return true
}
func SendCreatRoomRetCode(c *gotcp.Conn, reqHead *common.HeadV3, ret uint16, errMsg []byte) bool {
	head := new(common.HeadV3)
	if reqHead != nil {
		*head = *reqHead
	}
	rspBody := new(ht_muc.MucRspBody)
	rspBody.CreateRoomRspbody = new(ht_muc.CreateRoomRspBody)
	subMucRspBody := rspBody.GetCreateRoomRspbody()
	subMucRspBody.Status = &ht_muc.MucHeader{Code: proto.Uint32(uint32(ret)), Reason: errMsg}
	s, err := proto.Marshal(rspBody)
	if err != nil {
		infoLog.Printf("SendCreatRoomResp proto.Marshal failed from=%v to=%v cmd=%v seq=%v",
			head.From,
			head.To,
			head.Cmd,
			head.Seq)
		return false
	}
	head.Len = uint32(common.PacketV3HeadLen + len(s) + 1) //整个报文长度
	head.Ret = ret
	buf := make([]byte, head.Len)
	buf[0] = common.HTV3MagicBegin
	err = common.SerialHeadV3ToSlice(head, buf[1:])
	if err != nil {
		infoLog.Printf("SendCreatRoomResp SerialHeadV3ToSlice failed")
		return false
	}
	copy(buf[common.PacketV3HeadLen:], s) // return code
	buf[head.Len-1] = common.HTV3MagicEnd

	rspPacket := common.NewHeadV3Packet(buf)
	c.AsyncWritePacket(rspPacket, time.Second)
	return true
}
func SendCreatRoomResp(c *gotcp.Conn, reqHead *common.HeadV3, subMucRspBody *ht_muc.CreateRoomRspBody) bool {
	head := new(common.HeadV3)
	if reqHead != nil {
		*head = *reqHead
	}
	rspBody := new(ht_muc.MucRspBody)
	rspBody.CreateRoomRspbody = subMucRspBody
	s, err := proto.Marshal(rspBody)
	if err != nil {
		infoLog.Printf("SendCreatRoomResp proto.Marshal failed from=%v to=%v cmd=%v seq=%v",
			head.From,
			head.To,
			head.Cmd,
			head.Seq)
		return false
	}

	head.Len = uint32(common.PacketV3HeadLen + len(s) + 1) //整个报文长度
	head.Ret = uint16(subMucRspBody.GetStatus().GetCode())
	buf := make([]byte, head.Len)
	buf[0] = common.HTV3MagicBegin
	err = common.SerialHeadV3ToSlice(head, buf[1:])
	if err != nil {
		infoLog.Printf("SendCreatRoomResp SerialHeadV3ToSlice failed")
		return false
	}
	copy(buf[common.PacketV3HeadLen:], s) // return code
	buf[head.Len-1] = common.HTV3MagicEnd

	rspPacket := common.NewHeadV3Packet(buf)
	c.AsyncWritePacket(rspPacket, time.Second)
	return true
}
func UidIsInSlice(uidList []uint32, uid uint32) bool {
	if uid == 0 || len(uidList) == 0 {
		return false
	}
	for _, v := range uidList {
		if uid == v {
			return true
		}
	}

	return false
}

// 2.proc creater invite member
// func ProcInviteMember(c *gotcp.Conn, p gotcp.Packet) bool {
// 	result := uint16(ht_muc.MUC_RET_CODE_RET_SUCCESS)
// 	var head *common.HeadV3
// 	packet, ok := p.(*common.HeadV3Packet)
// 	if !ok {
// 		infoLog.Printf("ProcInviteMember Convert to HeadV3Packet failed")
// 		result = uint16(ht_muc.MUC_RET_CODE_RET_INTERNAL_ERR)
// 		errMsg := []byte("Internal error")
// 		SendInviteMemberRetCode(c, head, result, errMsg)
// 		return false
// 	}
// 	head, err := packet.GetHead()
// 	if err != nil {
// 		infoLog.Printf("ProcInviteMember Get head faild")
// 		result = uint16(ht_muc.MUC_RET_CODE_RET_INTERNAL_ERR)
// 		errMsg := []byte("Internal error")
// 		SendInviteMemberRetCode(c, head, result, errMsg)
// 		return false
// 	}

// 	reqBody := &ht_muc.MucReqBody{}
// 	err = proto.Unmarshal(packet.GetBody(), reqBody)
// 	if err != nil {
// 		infoLog.Printf("ProcInviteMember proto Unmarshal failed")
// 		result = uint16(ht_muc.MUC_RET_CODE_RET_PB_ERR)
// 		errMsg := []byte("pb unmarshal error")
// 		SendInviteMemberRetCode(c, head, result, errMsg)
// 		return false
// 	}

// 	// add static
// 	attr := "gomuc/invite_member_req_count"
// 	libcomm.AttrAdd(attr, 1)

// 	subReqBody := reqBody.GetInviteMemberReqbody()
// 	inviteUid := subReqBody.GetInviteUid()
// 	roomId := subReqBody.GetRoomId()
// 	memberList := subReqBody.GetMembers()
// 	infoLog.Printf("ProcInviteMembe recv from=%v to=%v cmd=%v seq=%v inviteUid=%v  roomId=%v count=%v",
// 		head.From,
// 		head.To,
// 		head.Cmd,
// 		head.Seq,
// 		inviteUid,
// 		roomId,
// 		len(memberList))

// 	// first check input param
// 	if inviteUid == 0 || roomId == 0 || len(memberList) == 0 {
// 		infoLog.Printf("ProcInviteMember invalid param")
// 		result = uint16(ht_muc.MUC_RET_CODE_RET_INPUT_PARAM_ERR)
// 		errMsg := []byte("input param error")
// 		SendInviteMemberRetCode(c, head, result, errMsg)
// 		return false
// 	}
// 	// second check black me list
// 	blackMeList, err := DbUtil.GetBlackMeList(inviteUid)
// 	if err != nil {
// 		infoLog.Printf("ProcCreateRoom GetBlackMeList failed uid=%v", inviteUid)
// 	}
// 	var totalBalck []uint32
// 	if len(blackMeList) != 0 {
// 		for i, v := range memberList {
// 			infoLog.Printf("ProcCreateRoom black me index=%v uid=%v", i, v)
// 			if UidIsInSlice(blackMeList, v.GetUid()) {
// 				infoLog.Printf("ProcCreateRoom uid=%v is in blackMeList", v)
// 				totalBalck = append(totalBalck, v.GetUid())
// 			}
// 		}
// 	}
// 	if len(totalBalck) != 0 {
// 		subMucRspBody := new(ht_muc.InviteMemberRspBody)
// 		errMsg := []byte("some one black me")
// 		subMucRspBody.Status = &ht_muc.MucHeader{Code: proto.Uint32(uint32(ht_muc.MUC_RET_CODE_RET_SOME_ONE_BLACK_ME)), Reason: errMsg}
// 		subMucRspBody.ListBlackMe = totalBalck
// 		SendInviteMemberResp(c, head, subMucRspBody)
// 		return false
// 	}
// 	// third invite member into room
// 	roomTS, _, err := roomManager.InviteMember(roomId, inviteUid, memberList)
// 	if err != nil {
// 		infoLog.Printf("ProcInviteMember roomId=%v inviteId=%v InviteMember failed err=%v", roomId, inviteUid, err)
// 		result = uint16(ht_muc.MUC_RET_CODE_RET_INTERNAL_ERR)
// 		errMsg := []byte("internal error")
// 		SendInviteMemberRetCode(c, head, result, errMsg)
// 		return false
// 	}

// 	// fourth broadcast message
// 	opInfo := &ht_muc.RoomMemberInfo{
// 		Uid:      proto.Uint32(inviteUid),
// 		NickName: subReqBody.GetNickName(),
// 	}
// 	err = roomManager.NotifyInviteMember(opInfo, opInfo, memberList, roomId, uint64(roomTS), uint32(ht_muc.ROOMID_FROM_TYPE_ENUM_FROM_INVITE))
// 	if err != nil {
// 		infoLog.Printf("ProcInviteMember roomId=%v inviteId=%v NotifyInviteMember failed", roomId, inviteUid)
// 	}

// 	infoLog.Printf("DEBUG ProcInviteMember roomId=%v roomTS=%v", roomId, roomTS)
// 	// fiveth send ack
// 	subMucRspBody := new(ht_muc.InviteMemberRspBody)
// 	subMucRspBody.Status = &ht_muc.MucHeader{Code: proto.Uint32(uint32(ht_muc.MUC_RET_CODE_RET_SUCCESS))}
// 	subMucRspBody.RoomTimestamp = proto.Uint64(uint64(roomTS))
// 	SendInviteMemberResp(c, head, subMucRspBody)
// 	return true
// }
// func SendInviteMemberRetCode(c *gotcp.Conn, reqHead *common.HeadV3, ret uint16, errMsg []byte) bool {
// 	head := new(common.HeadV3)
// 	if reqHead != nil {
// 		*head = *reqHead
// 	}
// 	rspBody := new(ht_muc.MucRspBody)
// 	rspBody.InviteMemberRspbody = new(ht_muc.InviteMemberRspBody)
// 	subMucRspBody := rspBody.GetInviteMemberRspbody()
// 	subMucRspBody.Status = &ht_muc.MucHeader{Code: proto.Uint32(uint32(ret)), Reason: errMsg}
// 	s, err := proto.Marshal(rspBody)
// 	if err != nil {
// 		infoLog.Printf("SendInviteMemberRetCode proto.Marshal failed from=%v to=%v cmd=%v seq=%v",
// 			head.From,
// 			head.To,
// 			head.Cmd,
// 			head.Seq)
// 		return false
// 	}
// 	head.Len = uint32(common.PacketV3HeadLen + len(s) + 1) //整个报文长度
// 	head.Ret = ret
// 	buf := make([]byte, head.Len)
// 	buf[0] = common.HTV3MagicBegin
// 	err = common.SerialHeadV3ToSlice(head, buf[1:])
// 	if err != nil {
// 		infoLog.Printf("SendInviteMemberRetCode SerialHeadV3ToSlice failed")
// 		return false
// 	}
// 	copy(buf[common.PacketV3HeadLen:], s) // return code
// 	buf[head.Len-1] = common.HTV3MagicEnd

// 	rspPacket := common.NewHeadV3Packet(buf)
// 	c.AsyncWritePacket(rspPacket, time.Second)
// 	return true
// }
// func SendInviteMemberResp(c *gotcp.Conn, reqHead *common.HeadV3, subMucRspBody *ht_muc.InviteMemberRspBody) bool {
// 	head := new(common.HeadV3)
// 	if reqHead != nil {
// 		*head = *reqHead
// 	}
// 	rspBody := new(ht_muc.MucRspBody)
// 	rspBody.InviteMemberRspbody = subMucRspBody
// 	s, err := proto.Marshal(rspBody)
// 	if err != nil {
// 		infoLog.Printf("SendInviteMemberResp proto.Marshal failed from=%v to=%v cmd=%v seq=%v",
// 			head.From,
// 			head.To,
// 			head.Cmd,
// 			head.Seq)
// 		return false
// 	}

// 	head.Len = uint32(common.PacketV3HeadLen + len(s) + 1) //整个报文长度
// 	head.Ret = uint16(subMucRspBody.GetStatus().GetCode())
// 	buf := make([]byte, head.Len)
// 	buf[0] = common.HTV3MagicBegin
// 	err = common.SerialHeadV3ToSlice(head, buf[1:])
// 	if err != nil {
// 		infoLog.Printf("SendInviteMemberResp SerialHeadV3ToSlice failed")
// 		return false
// 	}
// 	copy(buf[common.PacketV3HeadLen:], s) // return code
// 	buf[head.Len-1] = common.HTV3MagicEnd

// 	rspPacket := common.NewHeadV3Packet(buf)
// 	c.AsyncWritePacket(rspPacket, time.Second)
// 	return true
// }

// 3.Proc remove member
func ProcRemoveMember(c *gotcp.Conn, head *common.HeadV3, reqBody *ht_muc.MucReqBody) bool {
	result := uint16(ht_muc.MUC_RET_CODE_RET_SUCCESS)
	errMsg := []byte("operation success")
	bNeedCall := true
	defer func() {
		if bNeedCall {
			SendRemoveMemberRetCode(c, head, result, errMsg)
		}
	}()
	// 检查输入参数是否为空
	if head == nil || reqBody == nil {
		bNeedCall = false
		return false
	}

	// add static
	attr := "gomuc/remove_member_req_count"
	libcomm.AttrAdd(attr, 1)

	subReqBody := reqBody.GetRemoveMemberReqbody()
	adminUid := subReqBody.GetAdminUid()
	adminName := string(subReqBody.GetAdminName())
	roomId := subReqBody.GetRoomId()
	removeUid := subReqBody.GetRemoveUid()
	removeName := string(subReqBody.GetRemoveName())
	infoLog.Printf("ProcRemoveMembe recv from=%v to=%v cmd=%v seq=%v roomId=%v adminUid=%v removeUid=%v",
		head.From,
		head.To,
		head.Cmd,
		head.Seq,
		roomId,
		adminUid,
		removeUid)
	// first check input param
	if roomId == 0 || removeUid == 0 || adminUid == 0 || removeUid == adminUid {
		infoLog.Printf("ProcRemoveMember invalid param")
		result = uint16(ht_muc.MUC_RET_CODE_RET_INPUT_PARAM_ERR)
		errMsg = []byte("input param error")
		return false

	}

	// second check whether user is admin
	bCreater, err := roomManager.IsCreateUid(roomId, adminUid)
	if err != nil {
		infoLog.Printf("ProcRemoveMember check uid=%v roomManager.IsCreateUid failed err=%v", adminUid, err)
		result = uint16(ht_muc.MUC_RET_CODE_RET_INTERNAL_ERR)
		errMsg = []byte("Internal error")
		return false
	}

	if !bCreater { //如果不是创建者 继续检查是否是管理员
		bAdmin, err := roomManager.IsUserAdmin(roomId, adminUid)
		if err != nil {
			infoLog.Printf("ProcRemoveMember check uid=%v roomManager.IsUserAdmin failed err=%v", adminUid, err)
			result = uint16(ht_muc.MUC_RET_CODE_RET_INTERNAL_ERR)
			errMsg = []byte("Internal error")
			return false
		}
		// 判断用户是否是管理员 如果不是管理员则出错
		if !bAdmin {
			infoLog.Printf("ProcRemoveMember uid=%v is not admin permission denied", adminUid)
			result = uint16(ht_muc.MUC_RET_CODE_RET_PERMISSION_DENIED)
			errMsg = []byte("permission denied")
			return false
		} else {
			// 继续检查被删除的用户是否是管理员
			bRemoveUserAdmind, err := roomManager.IsUserAdmin(roomId, removeUid)
			if err != nil {
				infoLog.Printf("ProcRemoveMember check uid=%v roomManager.IsUserAdmin failed err=%v", removeUid, err)
				result = uint16(ht_muc.MUC_RET_CODE_RET_INTERNAL_ERR)
				errMsg = []byte("Internal error")
				return false
			}
			if bRemoveUserAdmind {
				infoLog.Printf("ProcRemoveMember removeUid=%v is admin permission denied", removeUid)
				result = uint16(ht_muc.MUC_RET_CODE_RET_PERMISSION_DENIED)
				errMsg = []byte("permission denied")
				return false
			}
		}

	} else {
		infoLog.Printf("ProcRemoveMember adminUid=%v is creater", adminUid)
	}

	// third remove member in db
	roomTS, err := roomManager.RemoveMember(roomId, adminUid, removeUid)
	if err != nil {
		infoLog.Printf("ProcRemoveMember roomManager RemoveMember failed roomId=%v adminUid=%v removeUid=%v",
			roomId,
			adminUid,
			removeUid)
		result = uint16(ht_muc.MUC_RET_CODE_RET_INTERNAL_ERR)
		errMsg = []byte("Internal error")
		return false
	}
	// fourth notify remove member

	err = roomManager.NotifyRemoveMember(roomId, adminUid, adminName, removeUid, removeName, uint64(roomTS))
	if err != nil {
		infoLog.Printf("ProcRemoveMember NotifyRemoveMember faield roomId=%v adminUid=%v removeUid%v err=%v",
			roomId,
			adminUid,
			removeUid,
			err)
	}
	infoLog.Printf("DEBUG ProcRemoveMember roomId=%v roomTS=%v", roomId, roomTS)
	// fifth send ack
	subMucRspBody := new(ht_muc.RemoveMemberRspBody)
	subMucRspBody.Status = &ht_muc.MucHeader{Code: proto.Uint32(uint32(ht_muc.MUC_RET_CODE_RET_SUCCESS))}
	subMucRspBody.RoomTimestamp = proto.Uint64(uint64(roomTS))
	bNeedCall = false
	SendRemoveMemberResp(c, head, subMucRspBody)
	return true
}
func SendRemoveMemberRetCode(c *gotcp.Conn, reqHead *common.HeadV3, ret uint16, errMsg []byte) bool {
	head := new(common.HeadV3)
	if reqHead != nil {
		*head = *reqHead
	}
	rspBody := new(ht_muc.MucRspBody)
	rspBody.RemoveMemberRspbody = new(ht_muc.RemoveMemberRspBody)
	subMucRspBody := rspBody.GetRemoveMemberRspbody()
	subMucRspBody.Status = &ht_muc.MucHeader{Code: proto.Uint32(uint32(ret)), Reason: errMsg}
	s, err := proto.Marshal(rspBody)
	if err != nil {
		infoLog.Printf("SendRemoveMemberRetCode proto.Marshal failed from=%v to=%v cmd=%v seq=%v",
			head.From,
			head.To,
			head.Cmd,
			head.Seq)
		return false
	}
	head.Len = uint32(common.PacketV3HeadLen + len(s) + 1) //整个报文长度
	head.Ret = ret
	buf := make([]byte, head.Len)
	buf[0] = common.HTV3MagicBegin
	err = common.SerialHeadV3ToSlice(head, buf[1:])
	if err != nil {
		infoLog.Printf("SendRemoveMemberRetCode SerialHeadV3ToSlice failed")
		return false
	}
	copy(buf[common.PacketV3HeadLen:], s) // return code
	buf[head.Len-1] = common.HTV3MagicEnd

	rspPacket := common.NewHeadV3Packet(buf)
	c.AsyncWritePacket(rspPacket, time.Second)
	return true
}
func SendRemoveMemberResp(c *gotcp.Conn, reqHead *common.HeadV3, subMucRspBody *ht_muc.RemoveMemberRspBody) bool {
	head := new(common.HeadV3)
	if reqHead != nil {
		*head = *reqHead
	}
	rspBody := new(ht_muc.MucRspBody)
	rspBody.RemoveMemberRspbody = subMucRspBody
	s, err := proto.Marshal(rspBody)
	if err != nil {
		infoLog.Printf("SendRemoveMemberResp proto.Marshal failed from=%v to=%v cmd=%v seq=%v",
			head.From,
			head.To,
			head.Cmd,
			head.Seq)
		return false
	}

	head.Len = uint32(common.PacketV3HeadLen + len(s) + 1) //整个报文长度
	head.Ret = uint16(subMucRspBody.GetStatus().GetCode())
	buf := make([]byte, head.Len)
	buf[0] = common.HTV3MagicBegin
	err = common.SerialHeadV3ToSlice(head, buf[1:])
	if err != nil {
		infoLog.Printf("SendRemoveMemberResp SerialHeadV3ToSlice failed")
		return false
	}
	copy(buf[common.PacketV3HeadLen:], s) // return code
	buf[head.Len-1] = common.HTV3MagicEnd

	rspPacket := common.NewHeadV3Packet(buf)
	c.AsyncWritePacket(rspPacket, time.Second)
	return true
}

// 4.proc member quit
func ProcQuitRoom(c *gotcp.Conn, head *common.HeadV3, reqBody *ht_muc.MucReqBody) bool {
	result := uint16(ht_muc.MUC_RET_CODE_RET_SUCCESS)
	errMsg := []byte("operation success")
	bNeedCall := true
	defer func() {
		if bNeedCall {
			SendQuitRoomRetCode(c, head, result, errMsg)
		}
	}()

	// 检查输入参数是否为空
	if head == nil || reqBody == nil {
		bNeedCall = false
		return false
	}

	// add static
	attr := "gomuc/quit_room_req_count"
	libcomm.AttrAdd(attr, 1)

	subReqBody := reqBody.GetQuitRoomReqbody()
	quitUid := subReqBody.GetQuitUid()
	quitName := subReqBody.GetQuitName()
	roomId := subReqBody.GetRoomId()
	infoLog.Printf("ProcQuitRoom recv from=%v to=%v cmd=%v seq=%v roomId=%v quitUid=%v",
		head.From,
		head.To,
		head.Cmd,
		head.Seq,
		roomId,
		quitUid)
	// first check input param
	if roomId == 0 || quitUid == 0 {
		infoLog.Printf("ProcQuitRoom invalid param")
		result = uint16(ht_muc.MUC_RET_CODE_RET_INPUT_PARAM_ERR)
		errMsg = []byte("input param error")
		return false
	}

	// second process member quit
	newCreater, roomTS, err := roomManager.QuitMucRoom(roomId, quitUid)
	if err != nil {
		infoLog.Printf("ProcQuitRoom internal err")
		result = uint16(ht_muc.MUC_RET_CODE_RET_INTERNAL_ERR)
		errMsg = []byte("internal error")
		return false
	}

	// third broadcast message
	err = roomManager.NotifyMemberQuit(roomId, quitUid, string(quitName), newCreater, uint64(roomTS))
	if err != nil {
		infoLog.Printf("ProcQuitRoom NotifyMemberQuit failed roomId=%v quitUid=%v newCreater=%v",
			roomId,
			quitUid,
			newCreater)
	}

	infoLog.Printf("DEBUG ProcQuitRoom roomId=%v roomTS=%v", roomId, roomTS)
	// sixth send ack
	subMucRspBody := new(ht_muc.QuitRoomRspBody)
	subMucRspBody.Status = &ht_muc.MucHeader{Code: proto.Uint32(uint32(ht_muc.MUC_RET_CODE_RET_SUCCESS))}
	subMucRspBody.RoomTimestamp = proto.Uint64(uint64(roomTS))
	bNeedCall = false
	SendQuitRoomRsp(c, head, subMucRspBody)
	return true
}
func SendQuitRoomRetCode(c *gotcp.Conn, reqHead *common.HeadV3, ret uint16, errMsg []byte) bool {
	head := new(common.HeadV3)
	if reqHead != nil {
		*head = *reqHead
	}
	rspBody := new(ht_muc.MucRspBody)
	rspBody.QuitRoomRspbody = new(ht_muc.QuitRoomRspBody)
	subMucRspBody := rspBody.GetQuitRoomRspbody()
	subMucRspBody.Status = &ht_muc.MucHeader{Code: proto.Uint32(uint32(ret)), Reason: errMsg}
	s, err := proto.Marshal(rspBody)
	if err != nil {
		infoLog.Printf("SendQuitRoomRetCode proto.Marshal failed from=%v to=%v cmd=%v seq=%v",
			head.From,
			head.To,
			head.Cmd,
			head.Seq)
		return false
	}
	head.Len = uint32(common.PacketV3HeadLen + len(s) + 1) //整个报文长度
	head.Ret = ret
	buf := make([]byte, head.Len)
	buf[0] = common.HTV3MagicBegin
	err = common.SerialHeadV3ToSlice(head, buf[1:])
	if err != nil {
		infoLog.Printf("SendCreatRoomResp SerialHeadV3ToSlice failed")
		return false
	}
	copy(buf[common.PacketV3HeadLen:], s) // return code
	buf[head.Len-1] = common.HTV3MagicEnd

	rspPacket := common.NewHeadV3Packet(buf)
	c.AsyncWritePacket(rspPacket, time.Second)
	return true
}
func SendQuitRoomRsp(c *gotcp.Conn, reqHead *common.HeadV3, subMucRspBody *ht_muc.QuitRoomRspBody) bool {
	head := new(common.HeadV3)
	if reqHead != nil {
		*head = *reqHead
	}
	rspBody := new(ht_muc.MucRspBody)
	rspBody.QuitRoomRspbody = subMucRspBody
	s, err := proto.Marshal(rspBody)
	if err != nil {
		infoLog.Printf("SendQuitRoomRsp proto.Marshal failed from=%v to=%v cmd=%v seq=%v",
			head.From,
			head.To,
			head.Cmd,
			head.Seq)
		return false
	}

	head.Len = uint32(common.PacketV3HeadLen + len(s) + 1) //整个报文长度
	head.Ret = uint16(subMucRspBody.GetStatus().GetCode())
	buf := make([]byte, head.Len)
	buf[0] = common.HTV3MagicBegin
	err = common.SerialHeadV3ToSlice(head, buf[1:])
	if err != nil {
		infoLog.Printf("SendQuitRoomRsp SerialHeadV3ToSlice failed")
		return false
	}
	copy(buf[common.PacketV3HeadLen:], s) // return code
	buf[head.Len-1] = common.HTV3MagicEnd

	rspPacket := common.NewHeadV3Packet(buf)
	c.AsyncWritePacket(rspPacket, time.Second)
	return true
}

// 5.proc modify room name
func ProcModifyRoomName(c *gotcp.Conn, head *common.HeadV3, reqBody *ht_muc.MucReqBody) bool {
	result := uint16(ht_muc.MUC_RET_CODE_RET_SUCCESS)
	errMsg := []byte("operation success")
	bNeedCall := true
	defer func() {
		if bNeedCall {
			SendModifyRoomNameRetCode(c, head, result, errMsg)
		}
	}()

	// 检查输入参数是否为空
	if head == nil || reqBody == nil {
		bNeedCall = false
		return false
	}

	// add static
	attr := "gomuc/modify_room_name_req_count"
	libcomm.AttrAdd(attr, 1)

	subReqBody := reqBody.GetModifyRoomNameReqbody()
	opUid := subReqBody.GetOpUid()
	opName := subReqBody.GetOpName()
	roomId := subReqBody.GetRoomId()
	roomName := string(subReqBody.GetRoomName())
	infoLog.Printf("ProcModifyRoomName recv from=%v to=%v cmd=%v seq=%v roomId=%v roomName=%s opUid=%v",
		head.From,
		head.To,
		head.Cmd,
		head.Seq,
		roomId,
		roomName,
		opUid)
	// first check input param
	if roomId == 0 || opUid == 0 || len(roomName) == 0 {
		infoLog.Printf("ProcModifyRoomName invalid param")
		result = uint16(ht_muc.MUC_RET_CODE_RET_INPUT_PARAM_ERR)
		errMsg = []byte("input param error")
		return false
	}

	// second modify room in db and memory
	roomTS, err := roomManager.ModifyRoomName(roomId, opUid, roomName)
	if err != nil {
		infoLog.Printf("ProcModifyRoomName roomManager.ModifyRoomName failed roomId=%v opUid=%v roomName=%s",
			roomId,
			opUid,
			roomName)
		result = uint16(ht_muc.MUC_RET_CODE_RET_INTERNAL_ERR)
		errMsg = []byte("internal error")
		return false
	}
	// third broadcast notifycation
	err = roomManager.NotifyModifyRoomName(roomId, opUid, string(opName), string(roomName), uint64(roomTS))
	if err != nil {
		infoLog.Printf("ProcModifyRoomName roomManager.NotifyModifyRoomName failed roomId=%v opUid=%v roomName=%s",
			roomId,
			opUid,
			roomName)
	}

	infoLog.Printf("DEBUG ProcModifyRoomName roomId=%v roomTS=%v", roomId, roomTS)
	// sixth send ack
	subMucRspBody := new(ht_muc.ModifyRoomNameRspBody)
	subMucRspBody.Status = &ht_muc.MucHeader{Code: proto.Uint32(uint32(ht_muc.MUC_RET_CODE_RET_SUCCESS))}
	subMucRspBody.RoomTimestamp = proto.Uint64(uint64(roomTS))
	bNeedCall = false
	SendModifyRoomNameResp(c, head, subMucRspBody)
	return true
}
func SendModifyRoomNameRetCode(c *gotcp.Conn, reqHead *common.HeadV3, ret uint16, errMsg []byte) bool {
	head := new(common.HeadV3)
	if reqHead != nil {
		*head = *reqHead
	}
	rspBody := new(ht_muc.MucRspBody)
	rspBody.ModifyRoomNameRspbody = new(ht_muc.ModifyRoomNameRspBody)
	subMucRspBody := rspBody.GetModifyRoomNameRspbody()
	subMucRspBody.Status = &ht_muc.MucHeader{Code: proto.Uint32(uint32(ret)), Reason: errMsg}
	s, err := proto.Marshal(rspBody)
	if err != nil {
		infoLog.Printf("SendModifyRoomNameRetCode proto.Marshal failed from=%v to=%v cmd=%v seq=%v",
			head.From,
			head.To,
			head.Cmd,
			head.Seq)
		return false
	}
	head.Len = uint32(common.PacketV3HeadLen + len(s) + 1) //整个报文长度
	head.Ret = ret
	buf := make([]byte, head.Len)
	buf[0] = common.HTV3MagicBegin
	err = common.SerialHeadV3ToSlice(head, buf[1:])
	if err != nil {
		infoLog.Printf("SendModifyRoomNameRetCode SerialHeadV3ToSlice failed")
		return false
	}
	copy(buf[common.PacketV3HeadLen:], s) // return code
	buf[head.Len-1] = common.HTV3MagicEnd

	rspPacket := common.NewHeadV3Packet(buf)
	c.AsyncWritePacket(rspPacket, time.Second)
	return true
}
func SendModifyRoomNameResp(c *gotcp.Conn, reqHead *common.HeadV3, subMucRspBody *ht_muc.ModifyRoomNameRspBody) bool {
	head := new(common.HeadV3)
	if reqHead != nil {
		*head = *reqHead
	}
	rspBody := new(ht_muc.MucRspBody)
	rspBody.ModifyRoomNameRspbody = subMucRspBody
	s, err := proto.Marshal(rspBody)
	if err != nil {
		infoLog.Printf("SendModifyRoomNameResp proto.Marshal failed from=%v to=%v cmd=%v seq=%v",
			head.From,
			head.To,
			head.Cmd,
			head.Seq)
		return false
	}

	head.Len = uint32(common.PacketV3HeadLen + len(s) + 1) //整个报文长度
	head.Ret = uint16(subMucRspBody.GetStatus().GetCode())
	buf := make([]byte, head.Len)
	buf[0] = common.HTV3MagicBegin
	err = common.SerialHeadV3ToSlice(head, buf[1:])
	if err != nil {
		infoLog.Printf("SendModifyRoomNameResp SerialHeadV3ToSlice failed")
		return false
	}
	copy(buf[common.PacketV3HeadLen:], s) // return code
	buf[head.Len-1] = common.HTV3MagicEnd

	rspPacket := common.NewHeadV3Packet(buf)
	c.AsyncWritePacket(rspPacket, time.Second)
	return true
}

// 6.proc modify member name
func ProcModifyMemberName(c *gotcp.Conn, head *common.HeadV3, reqBody *ht_muc.MucReqBody) bool {
	result := uint16(ht_muc.MUC_RET_CODE_RET_SUCCESS)
	errMsg := []byte("operation success")
	bNeedCall := true
	defer func() {
		if bNeedCall {
			SendModifyMemberNameRetCode(c, head, result, errMsg)
		}
	}()

	// 检查输入参数是否为空
	if head == nil || reqBody == nil {
		bNeedCall = false
		return false
	}

	// add static
	attr := "gomuc/modify_member_name_req_count"
	libcomm.AttrAdd(attr, 1)

	subReqBody := reqBody.GetModifyMemberNameReqbody()
	roomId := subReqBody.GetRoomId()
	opUid := subReqBody.GetOpUid()
	opName := string(subReqBody.GetOpName())
	infoLog.Printf("ProcModifyMemberName recv from=%v to=%v cmd=%v seq=%v roomId=%v opUid=%v opName=%v",
		head.From,
		head.To,
		head.Cmd,
		head.Seq,
		roomId,
		opUid,
		opName)
	// first check input param
	if roomId == 0 || opUid == 0 || len(opName) == 0 {
		infoLog.Printf("ProcModifyMemberName invalid param")
		result = uint16(ht_muc.MUC_RET_CODE_RET_INPUT_PARAM_ERR)
		errMsg = []byte("input param error")
		return false
	}

	// second modify member name
	roomTS, err := roomManager.ModifyMemberName(roomId, opUid, opName)
	if err != nil {
		infoLog.Printf("ProcModifyMemberName failed roomId=%v opUid=%v opName=%s", roomId, opUid, opName)
		result = uint16(ht_muc.MUC_RET_CODE_RET_INTERNAL_ERR)
		errMsg = []byte("internal error")
		return false
	}

	// third broadcast member name chage
	err = roomManager.NotifyModifyMemberName(roomId, opUid, opName, uint64(roomTS))
	if err != nil {
		infoLog.Printf("ProcModifyMemberName NotifyModifyMemberName faild roomId=%v opUid=%v opName=%v",
			roomId,
			opUid,
			opName)
	}

	infoLog.Printf("DEBUG ProcModifyMemberName roomId=%v roomTS=%v", roomId, roomTS)
	// sixth send ack
	subMucRspBody := new(ht_muc.ModifyMemberNameRspBody)
	subMucRspBody.Status = &ht_muc.MucHeader{Code: proto.Uint32(uint32(ht_muc.MUC_RET_CODE_RET_SUCCESS))}
	subMucRspBody.RoomTimestamp = proto.Uint64(uint64(roomTS))
	bNeedCall = false
	SendModifyMemberNameResp(c, head, subMucRspBody)
	return true
}
func SendModifyMemberNameRetCode(c *gotcp.Conn, reqHead *common.HeadV3, ret uint16, errMsg []byte) bool {
	head := new(common.HeadV3)
	if reqHead != nil {
		*head = *reqHead
	}
	rspBody := new(ht_muc.MucRspBody)
	rspBody.ModifyMemberNameRspbody = new(ht_muc.ModifyMemberNameRspBody)
	subMucRspBody := rspBody.GetModifyMemberNameRspbody()
	subMucRspBody.Status = &ht_muc.MucHeader{Code: proto.Uint32(uint32(ret)), Reason: errMsg}
	s, err := proto.Marshal(rspBody)
	if err != nil {
		infoLog.Printf("SendModifyMemberNameRetCode proto.Marshal failed from=%v to=%v cmd=%v seq=%v",
			head.From,
			head.To,
			head.Cmd,
			head.Seq)
		return false
	}
	head.Len = uint32(common.PacketV3HeadLen + len(s) + 1) //整个报文长度
	head.Ret = ret
	buf := make([]byte, head.Len)
	buf[0] = common.HTV3MagicBegin
	err = common.SerialHeadV3ToSlice(head, buf[1:])
	if err != nil {
		infoLog.Printf("SendModifyMemberNameRetCode SerialHeadV3ToSlice failed")
		return false
	}
	copy(buf[common.PacketV3HeadLen:], s) // return code
	buf[head.Len-1] = common.HTV3MagicEnd

	rspPacket := common.NewHeadV3Packet(buf)
	c.AsyncWritePacket(rspPacket, time.Second)
	return true
}
func SendModifyMemberNameResp(c *gotcp.Conn, reqHead *common.HeadV3, subMucRspBody *ht_muc.ModifyMemberNameRspBody) bool {
	head := new(common.HeadV3)
	if reqHead != nil {
		*head = *reqHead
	}
	rspBody := new(ht_muc.MucRspBody)
	rspBody.ModifyMemberNameRspbody = subMucRspBody
	s, err := proto.Marshal(rspBody)
	if err != nil {
		infoLog.Printf("SendModifyMemberNameResp proto.Marshal failed from=%v to=%v cmd=%v seq=%v",
			head.From,
			head.To,
			head.Cmd,
			head.Seq)
		return false
	}

	head.Len = uint32(common.PacketV3HeadLen + len(s) + 1) //整个报文长度
	head.Ret = uint16(subMucRspBody.GetStatus().GetCode())
	buf := make([]byte, head.Len)
	buf[0] = common.HTV3MagicBegin
	err = common.SerialHeadV3ToSlice(head, buf[1:])
	if err != nil {
		infoLog.Printf("SendModifyMemberNameResp SerialHeadV3ToSlice failed")
		return false
	}
	copy(buf[common.PacketV3HeadLen:], s) // return code
	buf[head.Len-1] = common.HTV3MagicEnd

	rspPacket := common.NewHeadV3Packet(buf)
	c.AsyncWritePacket(rspPacket, time.Second)
	return true
}

// 7.Proc modify pushsetting
func ProcModifyPushSetting(c *gotcp.Conn, head *common.HeadV3, reqBody *ht_muc.MucReqBody) bool {
	result := uint16(ht_muc.MUC_RET_CODE_RET_SUCCESS)
	errMsg := []byte("operation success")
	bNeedCall := true
	defer func() {
		if bNeedCall {
			SendModifyPushSettingRetCode(c, head, result, errMsg)
		}
	}()

	// 检查输入参数是否为空
	if head == nil || reqBody == nil {
		bNeedCall = false
		return false
	}

	// add static
	attr := "gomuc/modify_push_setting_req_count"
	libcomm.AttrAdd(attr, 1)

	subReqBody := reqBody.GetModifyPushSettingReqbody()
	roomId := subReqBody.GetRoomId()
	opUid := subReqBody.GetOpUid()
	pushSetting := subReqBody.GetPushSetting()
	infoLog.Printf("ProcModifyPushSetting recv from=%v to=%v cmd=%v seq=%v roomId=%v opUid=%v pushSetting=%v",
		head.From,
		head.To,
		head.Cmd,
		head.Seq,
		roomId,
		opUid,
		pushSetting)
	// first check input param
	if roomId == 0 || opUid == 0 || pushSetting > 1 {
		infoLog.Printf("ProcModifyPushSetting invalid param")
		result = uint16(ht_muc.MUC_RET_CODE_RET_INPUT_PARAM_ERR)
		errMsg = []byte("input param error")
		return false
	}

	// second modify push setting in db
	err := roomManager.ModifyPushSetting(roomId, opUid, pushSetting)
	if err != nil {
		infoLog.Printf("ProcModifyPushSetting modify faield roomId=%v opUid=%v pushSettig=%v", roomId, opUid, pushSetting)
		result = uint16(ht_muc.MUC_RET_CODE_RET_INTERNAL_ERR)
		errMsg = []byte("internal error")
		return false
	}

	// third send ack
	subMucRspBody := new(ht_muc.ModifyPushSettingRspBody)
	subMucRspBody.Status = &ht_muc.MucHeader{Code: proto.Uint32(uint32(ht_muc.MUC_RET_CODE_RET_SUCCESS))}
	bNeedCall = false // 无需调用defer函数
	SendModifyPushSettingResp(c, head, subMucRspBody)
	return true
}
func SendModifyPushSettingRetCode(c *gotcp.Conn, reqHead *common.HeadV3, ret uint16, errMsg []byte) bool {
	head := new(common.HeadV3)
	if reqHead != nil {
		*head = *reqHead
	}
	rspBody := new(ht_muc.MucRspBody)
	rspBody.ModifyPushSettingRspbody = new(ht_muc.ModifyPushSettingRspBody)
	subMucRspBody := rspBody.GetModifyPushSettingRspbody()
	subMucRspBody.Status = &ht_muc.MucHeader{Code: proto.Uint32(uint32(ret)), Reason: errMsg}
	s, err := proto.Marshal(rspBody)
	if err != nil {
		infoLog.Printf("SendModifyPushSettingRetCode proto.Marshal failed from=%v to=%v cmd=%v seq=%v",
			head.From,
			head.To,
			head.Cmd,
			head.Seq)
		return false
	}
	head.Len = uint32(common.PacketV3HeadLen + len(s) + 1) //整个报文长度
	head.Ret = ret
	buf := make([]byte, head.Len)
	buf[0] = common.HTV3MagicBegin
	err = common.SerialHeadV3ToSlice(head, buf[1:])
	if err != nil {
		infoLog.Printf("SendModifyPushSettingRetCode SerialHeadV3ToSlice failed")
		return false
	}
	copy(buf[common.PacketV3HeadLen:], s) // return code
	buf[head.Len-1] = common.HTV3MagicEnd

	rspPacket := common.NewHeadV3Packet(buf)
	c.AsyncWritePacket(rspPacket, time.Second)
	return true
}
func SendModifyPushSettingResp(c *gotcp.Conn, reqHead *common.HeadV3, subMucRspBody *ht_muc.ModifyPushSettingRspBody) bool {
	head := new(common.HeadV3)
	if reqHead != nil {
		*head = *reqHead
	}
	rspBody := new(ht_muc.MucRspBody)
	rspBody.ModifyPushSettingRspbody = subMucRspBody
	s, err := proto.Marshal(rspBody)
	if err != nil {
		infoLog.Printf("SendModifyPushSettingResp proto.Marshal failed from=%v to=%v cmd=%v seq=%v",
			head.From,
			head.To,
			head.Cmd,
			head.Seq)
		return false
	}

	head.Len = uint32(common.PacketV3HeadLen + len(s) + 1) //整个报文长度
	head.Ret = uint16(subMucRspBody.GetStatus().GetCode())
	buf := make([]byte, head.Len)
	buf[0] = common.HTV3MagicBegin
	err = common.SerialHeadV3ToSlice(head, buf[1:])
	if err != nil {
		infoLog.Printf("SendModifyPushSettingResp SerialHeadV3ToSlice failed")
		return false
	}
	copy(buf[common.PacketV3HeadLen:], s) // return code
	buf[head.Len-1] = common.HTV3MagicEnd

	rspPacket := common.NewHeadV3Packet(buf)
	c.AsyncWritePacket(rspPacket, time.Second)
	return true
}

// 8.Proc Get Room Info
func ProcGetRoomInfo(c *gotcp.Conn, head *common.HeadV3, reqBody *ht_muc.MucReqBody) bool {
	result := uint16(ht_muc.MUC_RET_CODE_RET_SUCCESS)
	errMsg := []byte("operation success")
	bNeedCall := true
	defer func() {
		if bNeedCall {
			SendGetRoomInfoRetCode(c, head, result, errMsg)
		}
	}()

	// 检查输入参数是否为空
	if head == nil || reqBody == nil {
		bNeedCall = false
		return false
	}

	// add static
	attr := "gomuc/get_room_info_req_count"
	libcomm.AttrAdd(attr, 1)

	subReqBody := reqBody.GetGetRoomInfoReqbody()
	roomId := subReqBody.GetRoomId()
	opUid := subReqBody.GetOpUid()
	roomTS := subReqBody.GetRoomTimestamp()
	infoLog.Printf("ProcGetRoomInfo recv from=%v to=%v cmd=%v seq=%v roomId=%v opUid=%v roomTS=%v",
		head.From,
		head.To,
		head.Cmd,
		head.Seq,
		roomId,
		opUid,
		roomTS)
	// first check input param
	if roomId == 0 || opUid == 0 {
		infoLog.Printf("ProcGetRoomInfo invalid param")
		result = uint16(ht_muc.MUC_RET_CODE_RET_INPUT_PARAM_ERR)
		errMsg = []byte("input param error")
		return false
	}

	// second get room info
	roomInfo, err := roomManager.GetRoomInfo(roomId)
	if err != nil {
		infoLog.Printf("ProcGetRoomInfo faield roomId=%v opUid=%v roomTS=%v", roomId, opUid, roomTS)
		result = uint16(ht_muc.MUC_RET_CODE_RET_INTERNAL_ERR)
		errMsg = []byte("internal error")
		return false
	}

	// third send ack
	subMucRspBody := new(ht_muc.GetRoomInfoRspBody)
	subMucRspBody.Status = &ht_muc.MucHeader{Code: proto.Uint32(uint32(ht_muc.MUC_RET_CODE_RET_SUCCESS))}
	var memberInfo []*ht_muc.RoomMemberInfo
	memberList := roomInfo.MemberList
	var pushSetting uint32 = 0
	for _, v := range memberList {
		if v.Uid == opUid {
			pushSetting = v.PushSetting
		}
		iterm := &ht_muc.RoomMemberInfo{
			Uid:      proto.Uint32(v.Uid),
			NickName: []byte(v.NickName),
		}
		memberInfo = append(memberInfo, iterm)
	}

	var adminList []uint32
	for _, v := range roomInfo.AdminList {
		if v != 0 {
			adminList = append(adminList, v)
		}
	}
	verify := ht_muc.VERIFY_STAT(roomInfo.VerifyStat)
	subMucRspBody.RoomInfo = &ht_muc.RoomInfoBody{
		RoomId:       proto.Uint32(roomId),
		CreateUid:    proto.Uint32(roomInfo.CreateUid),
		ListAdminUid: adminList,
		AdminLimit:   proto.Uint32(roomInfo.AdminLimit),
		RoomLimit:    proto.Uint32(roomInfo.MemberLimit),
		RoomName:     []byte(roomInfo.RoomName),
		RoomDesc:     []byte(roomInfo.RoomDesc),
		VerifyStat:   &verify,
		Announcement: &ht_muc.AnnoType{
			PublishUid:  proto.Uint32(roomInfo.Announcement.PublishUid),
			PublishTs:   proto.Uint32(roomInfo.Announcement.PublishTS),
			AnnoContent: []byte(roomInfo.Announcement.AnnoContect),
		},
		RoomTimestamp: proto.Uint64(uint64(roomInfo.RoomTS)),
		PushSetting:   proto.Uint32(pushSetting),
		Members:       memberInfo,
	}
	infoLog.Printf("DEBUG ProcGetRoomInfo roomId=%v creatUid=%v admminList=%v adminLimit=%v roomLimit=%v roomName=%s roomDesc=%s verifyStat=%v",
		roomId,
		roomInfo.CreateUid,
		adminList,
		roomInfo.AdminLimit,
		roomInfo.MemberLimit,
		roomInfo.RoomName,
		roomInfo.RoomDesc,
		verify)
	infoLog.Printf("DEBUG ProcGetRoomInfo publishUid=%v publishTS=%v AnnoCotent=%s roomTS=%v pushSetting=%v members=%v",
		roomInfo.Announcement.PublishUid,
		roomInfo.Announcement.PublishTS,
		roomInfo.Announcement.AnnoContect,
		roomInfo.RoomTS,
		pushSetting,
		roomInfo.MemberList)
	bNeedCall = false
	SendGetRoomInfoResp(c, head, subMucRspBody)
	return true
}
func SendGetRoomInfoRetCode(c *gotcp.Conn, reqHead *common.HeadV3, ret uint16, errMsg []byte) bool {
	head := new(common.HeadV3)
	if reqHead != nil {
		*head = *reqHead
	}
	rspBody := new(ht_muc.MucRspBody)
	rspBody.GetRoomInfoRspbody = new(ht_muc.GetRoomInfoRspBody)
	subMucRspBody := rspBody.GetGetRoomInfoRspbody()
	subMucRspBody.Status = &ht_muc.MucHeader{Code: proto.Uint32(uint32(ret)), Reason: errMsg}
	s, err := proto.Marshal(rspBody)
	if err != nil {
		infoLog.Printf("SendGetRoomInfoRetCode proto.Marshal failed from=%v to=%v cmd=%v seq=%v",
			head.From,
			head.To,
			head.Cmd,
			head.Seq)
		return false
	}
	head.Len = uint32(common.PacketV3HeadLen + len(s) + 1) //整个报文长度
	head.Ret = ret
	buf := make([]byte, head.Len)
	buf[0] = common.HTV3MagicBegin
	err = common.SerialHeadV3ToSlice(head, buf[1:])
	if err != nil {
		infoLog.Printf("SendGetRoomInfoRetCode SerialHeadV3ToSlice failed")
		return false
	}
	copy(buf[common.PacketV3HeadLen:], s) // return code
	buf[head.Len-1] = common.HTV3MagicEnd

	rspPacket := common.NewHeadV3Packet(buf)
	c.AsyncWritePacket(rspPacket, time.Second)
	return true
}
func SendGetRoomInfoResp(c *gotcp.Conn, reqHead *common.HeadV3, subMucRspBody *ht_muc.GetRoomInfoRspBody) bool {
	head := new(common.HeadV3)
	if reqHead != nil {
		*head = *reqHead
	}
	rspBody := new(ht_muc.MucRspBody)
	rspBody.GetRoomInfoRspbody = subMucRspBody
	s, err := proto.Marshal(rspBody)
	if err != nil {
		infoLog.Printf("SendGetRoomInfoResp proto.Marshal failed from=%v to=%v cmd=%v seq=%v",
			head.From,
			head.To,
			head.Cmd,
			head.Seq)
		return false
	}

	head.Len = uint32(common.PacketV3HeadLen + len(s) + 1) //整个报文长度
	head.Ret = uint16(subMucRspBody.GetStatus().GetCode())
	buf := make([]byte, head.Len)
	buf[0] = common.HTV3MagicBegin
	err = common.SerialHeadV3ToSlice(head, buf[1:])
	if err != nil {
		infoLog.Printf("SendGetRoomInfoResp SerialHeadV3ToSlice failed")
		return false
	}
	copy(buf[common.PacketV3HeadLen:], s) // return code
	buf[head.Len-1] = common.HTV3MagicEnd

	rspPacket := common.NewHeadV3Packet(buf)
	c.AsyncWritePacket(rspPacket, time.Second)
	return true
}

// 9.ProcMucMessage
func ProcMucMessage(c *gotcp.Conn, head *common.HeadV3, reqBody *ht_muc.MucReqBody) bool {
	result := uint16(ht_muc.MUC_RET_CODE_RET_SUCCESS)
	errMsg := []byte("operation success")
	bNeedCall := true
	defer func() {
		if bNeedCall {
			SendMucMessageRetCode(c, head, result, errMsg)
		}
	}()

	// 检查输入参数是否为空
	if head == nil || reqBody == nil {
		bNeedCall = false
		return false
	}

	// add static
	attr := "gomuc/muc_msg_req_count"
	libcomm.AttrAdd(attr, 1)

	// 获取muc msg构造解压对象
	subReqBody := reqBody.GetMucMessageReqbody()
	msg := subReqBody.GetMsg()
	compressBuff := bytes.NewBuffer(msg)
	r, err := zlib.NewReader(compressBuff)
	defer r.Close()
	if err != nil {
		infoLog.Printf("ProcMucMessage zlib.NewReader failed err=%v", err)
		result = uint16(ht_muc.MUC_RET_CODE_RET_INTERNAL_ERR)
		errMsg = []byte("Internal error")
		return false
	}
	// 解压缩
	unCompressSlice, err := ioutil.ReadAll(r)
	if err != nil {
		infoLog.Printf("ProcMucMessage zlib. unCompress failed")
		result = uint16(ht_muc.MUC_RET_CODE_RET_INTERNAL_ERR)
		errMsg = []byte("Internal error")
		return false
	}
	// 解压之后的内容构造json
	rootObj, err := simplejson.NewJson(unCompressSlice)
	if err != nil {
		infoLog.Printf("ProcMucMessage simplejson new packet error", err)
		result = uint16(ht_muc.MUC_RET_CODE_RET_INTERNAL_ERR)
		errMsg = []byte("NewJson failed")
		return false
	}
	// 群固有属性
	roomId := uint32(rootObj.Get("room_id").MustInt64(0))
	strMsgId := rootObj.Get("msg_id").MustString()
	strSendTime := rootObj.Get("send_time").MustString()
	strMsgType := rootObj.Get("msg_type").MustString()
	msgType := GetMessageType(strMsgType)
	// strMsgModel := rootObj.Get("msg_model").MustString()
	fromId := uint32(rootObj.Get("sender_id").MustInt64(0))
	// userProfile := rootObj.Get("sender_ts").MustUint64(0)
	strNickName := rootObj.Get("sender_name").MustString()
	remindType := rootObj.Get("at").Get("remind_type").MustInt(0)
	var remindList []uint32
	// remindType == 2 查看数组
	if remindType == 2 {
		tempList := rootObj.Get("at").Get("remind_list").MustArray()
		infoLog.Printf("ProcMucMessage remind_list=%v", tempList)
		for _, v := range tempList {
			tempUid, err := v.(json.Number).Int64()
			if err != nil {
				infoLog.Printf("ProcMucMessage remind_list get Uid faild")
				continue
			}
			remindList = append(remindList, uint32(tempUid))
		}
	}

	infoLog.Printf("ProcMucMessage recv from=%v to=%v cmd=%v seq=%v roomId=%v strMsgId=%s strMsgType=%s msgType=%v fromId=%v strNickName=%s sendTime=%s",
		head.From,
		head.To,
		head.Cmd,
		head.Seq,
		roomId,
		strMsgId,
		strMsgType,
		msgType,
		fromId,
		strNickName,
		strSendTime)
	// first check input param
	if head.From == 0 || head.To == 0 || roomId == 0 || fromId == 0 || head.From != fromId || head.To != roomId || msgType == MT_UNKNOWN {
		infoLog.Printf("ProcModifyPushSetting invalid param")
		result = uint16(ht_muc.MUC_RET_CODE_RET_INPUT_PARAM_ERR)
		errMsg = []byte("input param error")
		return false
	}

	pushInfo := util.MucPushInfo{NickName: strNickName, MsgId: strMsgId}
	switch msgType {
	case MT_TEXT:
		strMessageText := rootObj.Get("text").Get("text").MustString()
		if utf8.RuneCountInString(strMessageText) > MUC_MESSAGE_LNE_LIMIT { // 字符串超过长度截取部分子串
			pushInfo.PushParam = SubString(strMessageText, 0, MUC_MESSAGE_LNE_LIMIT)
		} else {
			pushInfo.PushParam = strMessageText
		}
		pushInfo.PushType = PUSH_TEXT
	case MT_TRANSLATE:
		strMessageText := rootObj.Get("translate").Get("src_text").MustString()
		if utf8.RuneCountInString(strMessageText) > MUC_MESSAGE_LNE_LIMIT { // 字符串超过长度截取部分子串
			pushInfo.PushParam = SubString(strMessageText, 0, MUC_MESSAGE_LNE_LIMIT)
		} else {
			pushInfo.PushParam = strMessageText
		}
		pushInfo.PushType = PUSH_TEXT
	case MT_PROFILE:
		strIntroduceName := rootObj.Get("introduction").Get("user_profile").Get("nick_name").MustString()
		if utf8.RuneCountInString(strIntroduceName) > MUC_MESSAGE_LNE_LIMIT { // 字符串超过长度截取部分子串
			pushInfo.PushParam = SubString(strIntroduceName, 0, MUC_MESSAGE_LNE_LIMIT)
		} else {
			pushInfo.PushParam = strIntroduceName
		}
		pushInfo.PushType = PUSH_INTRODUCE
	case MT_VOICE, MT_VOICETEXT:
		pushInfo.PushType = PUSH_VOICE

	case MT_IMAGE:
		pushInfo.PushType = PUSH_IMAGE

	case MT_LOCATE:
		pushInfo.PushType = PUSH_LOCATION

	case MT_CORRECTION:
		pushInfo.PushType = PUSH_CORRECT_SENTENCE

	case MT_DOODLE:
		pushInfo.PushType = PUSH_DOODLE

	case MT_STICKERS:
		pushInfo.PushType = PUSH_STICKERS

	case MT_VIDEO:
		pushInfo.PushType = PUSH_VIDEO

	case MT_LINK:
		pushInfo.PushType = PUSH_LINK

	case MT_CARD:
		pushInfo.PushType = PUSH_CARD
	default:
		infoLog.Printf("ProcMucMessage Unhandle type=%v", msgType)
	}

	infoLog.Printf("DEBUG ProcMucMessage msgId=%s", strMsgId)

	// second send ack
	subMucRspBody := new(ht_muc.MucMessageRspBody)
	subMucRspBody.Status = &ht_muc.MucHeader{Code: proto.Uint32(uint32(ht_muc.MUC_RET_CODE_RET_SUCCESS))}
	subMucRspBody.RspDetial = &ht_muc.MucRspDetial{
		MsgId: []byte(strMsgId),
	}
	bNeedCall = false
	SendMucMessageResp(c, head, subMucRspBody)

	// 后续失败都不用返回给客户端 直接返回false

	// 图片URL腾讯云服务的兼容策略(图片服务器新老版本的URL替换) 2015-06-29
	if msgType == MT_IMAGE {
		FixQcloudImageCompatible(rootObj.Get("image"))
	} else if msgType == MT_DOODLE {
		FixQcloudImageCompatible(rootObj.Get("doodle"))
	}

	// 添加消息经过服务器的附加属性
	roomInfo, err := roomManager.GetRoomInfo(roomId)
	if err != nil {
		infoLog.Printf("ProcMucMessage faield roomId=%v fromId=%v roomTS=%v", roomId, fromId, roomInfo.RoomTS)
		return false
	}
	formatTime, ts := GetCurrentFormatTime()
	rootObj.Set("server_time", formatTime)
	rootObj.Set("server_ts", ts)
	rootObj.Set("room_name", roomInfo.RoomName)
	rootObj.Set("room_ts", roomInfo.RoomTS)
	strMsgBody, err := rootObj.MarshalJSON()
	if err != nil {
		infoLog.Printf("ProcMucMessage simpleJson MarshalJSON failed roomId=%v fromId=%v seq=%v",
			roomId,
			fromId,
			head.Seq)
		return false
	}
	infoLog.Printf("ProcMucMessage outJson=%s", strMsgBody)
	// 压缩JSON消息数据
	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	w.Write(strMsgBody)
	w.Close()
	compressSliceByte := b.Bytes()

	// 拼装成新的pb格式
	bcReqBody := new(ht_muc.MucReqBody)
	bcReqBody.MucMessageReqbody = &ht_muc.MucMessageReqBody{Msg: compressSliceByte}
	bcSlice, err := proto.Marshal(bcReqBody)
	if err != nil {
		infoLog.Printf("ProcMucMessage proto.Marshal failed from=%v to=%v cmd=%v seq=%v",
			head.From,
			head.To,
			head.Cmd,
			head.Seq)
		return false
	}

	err = roomManager.BroadcastMucMessageCompatible(roomId, head, bcSlice, &pushInfo, uint32(remindType), remindList, compressSliceByte)
	if err != nil {
		infoLog.Printf("ProcMucMessage BroadcastMucMessage failed roomId=%v fromId=%v to=%v err=%v",
			roomId,
			head.From,
			head.To,
			err)
	}
	return true
}
func FixQcloudImageCompatible(imageObject *simplejson.Json) {
	if imageObject == nil {
		infoLog.Printf("FixQcloudImageCompatible nil imageObject")
		return
	}

	strUrl := imageObject.Get("url").MustString()
	infoLog.Printf("FixQcloudImageCompatible url=%s", strUrl)

	findHelloTalk := strings.Index(strUrl, "hellotalk.")
	var strThumbUrl string
	if findHelloTalk != -1 { // find hellotalk.
		findJpg := strings.LastIndex(strUrl, ".")
		if findJpg != -1 {
			strSuffix := strUrl[findJpg:]
			if strings.Compare(strSuffix, ".jpg") != 0 && strings.Compare(strSuffix, ".png") != 0 &&
				strings.Compare(strSuffix, ".JPG") != 0 && strings.Compare(strSuffix, ".PNG") != 0 {
				infoLog.Printf("FixQcloudImageCompatible Image URL unexpected suffix! url=%s", strUrl)
			}
			strThumbUrl = strUrl[0:findJpg] + "_thum" + strSuffix
			imageObject.Set("thumb_url", strThumbUrl)

		} else {
			infoLog.Printf("FixQcloudImageCompatible Image URL dont find .jpg. URL! url=%s", strUrl)
		}
	} else { // 不包含hellotalk.的URL要以http开头才处理
		if strings.Index(strUrl, "http") == 0 || strings.Index(strUrl, "HTTP") == 0 {
			strNewUrl := "p" + strUrl + "?m.jpg"
			imageObject.Set("url", strNewUrl) // 更新url

			if strings.Index(strUrl, "myqcloud.") != -1 {
				strThumbUrl = strUrl[0 : strings.LastIndex(strUrl, "/")+1]
				strThumbUrl += "scale"
			} else {
				strThumbUrl = strUrl + "/scale"
			}
			imageObject.Set("thumb_url", strThumbUrl) // 设置ThumbUrl
			infoLog.Printf("FixQcloudImageCompatible thumb_url=%s", strThumbUrl)
		} else {
			infoLog.Printf("FixQcloudImageCompatible Unexpecet image url! url=%s", strUrl)
		}
	}
	imageObject.Set("name", strUrl)
}

// iOS 端使用的时间戳是毫秒 所以需要将timeStamp 转成毫秒
func GetCurrentFormatTime() (formatTime string, timeStamp int64) {
	currTime := time.Now().UTC()
	formatTime = fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d",
		currTime.Year(),
		currTime.Month(),
		currTime.Day(),
		currTime.Hour(),
		currTime.Minute(),
		currTime.Second())
	timeStamp = currTime.UnixNano() / (1000 * 1000)
	return
}
func SendMucMessageRetCode(c *gotcp.Conn, reqHead *common.HeadV3, ret uint16, errMsg []byte) bool {
	head := new(common.HeadV3)
	if reqHead != nil {
		*head = *reqHead
	}
	rspBody := new(ht_muc.MucRspBody)
	rspBody.MucMessageRspbody = new(ht_muc.MucMessageRspBody)
	subMucRspBody := rspBody.GetMucMessageRspbody()
	subMucRspBody.Status = &ht_muc.MucHeader{Code: proto.Uint32(uint32(ret)), Reason: errMsg}
	s, err := proto.Marshal(rspBody)
	if err != nil {
		infoLog.Printf("SendMucMessageRetCode proto.Marshal failed from=%v to=%v cmd=%v seq=%v",
			head.From,
			head.To,
			head.Cmd,
			head.Seq)
		return false
	}
	head.Len = uint32(common.PacketV3HeadLen + len(s) + 1) //整个报文长度
	head.Ret = ret
	buf := make([]byte, head.Len)
	buf[0] = common.HTV3MagicBegin
	err = common.SerialHeadV3ToSlice(head, buf[1:])
	if err != nil {
		infoLog.Printf("SendMucMessageRetCode SerialHeadV3ToSlice failed")
		return false
	}
	copy(buf[common.PacketV3HeadLen:], s) // return code
	buf[head.Len-1] = common.HTV3MagicEnd

	rspPacket := common.NewHeadV3Packet(buf)
	c.AsyncWritePacket(rspPacket, time.Second)
	return true
}
func SendMucMessageResp(c *gotcp.Conn, reqHead *common.HeadV3, subMucRspBody *ht_muc.MucMessageRspBody) bool {
	head := new(common.HeadV3)
	if reqHead != nil {
		*head = *reqHead
	}
	rspBody := new(ht_muc.MucRspBody)
	rspBody.MucMessageRspbody = subMucRspBody
	s, err := proto.Marshal(rspBody)
	if err != nil {
		infoLog.Printf("SendMucMessageResp proto.Marshal failed from=%v to=%v cmd=%v seq=%v",
			head.From,
			head.To,
			head.Cmd,
			head.Seq)
		return false
	}

	head.Len = uint32(common.PacketV3HeadLen + len(s) + 1) //整个报文长度
	head.Ret = uint16(subMucRspBody.GetStatus().GetCode())
	buf := make([]byte, head.Len)
	buf[0] = common.HTV3MagicBegin
	err = common.SerialHeadV3ToSlice(head, buf[1:])
	if err != nil {
		infoLog.Printf("SendMucMessageResp SerialHeadV3ToSlice failed")
		return false
	}
	copy(buf[common.PacketV3HeadLen:], s) // return code
	buf[head.Len-1] = common.HTV3MagicEnd

	rspPacket := common.NewHeadV3Packet(buf)
	c.AsyncWritePacket(rspPacket, time.Second)
	return true
}

// 10.Proc add muc to contact list
func ProcAddRoomToContactList(c *gotcp.Conn, head *common.HeadV3, reqBody *ht_muc.MucReqBody) bool {
	result := uint16(ht_muc.MUC_RET_CODE_RET_SUCCESS)
	errMsg := []byte("operation success")
	bNeedCall := true
	defer func() {
		if bNeedCall {
			SendAddRoomToContactListRetCode(c, head, result, errMsg)
		}
	}()

	// 检查输入参数是否为空
	if head == nil || reqBody == nil {
		bNeedCall = false
		return false
	}

	// add static
	attr := "gomuc/add_room_to_contactlist_req_count"
	libcomm.AttrAdd(attr, 1)

	subReqBody := reqBody.GetAddRoomToContactListReqbody()
	roomId := subReqBody.GetRoomId()
	opUid := subReqBody.GetOpUid()
	opType := subReqBody.GetOpType()
	infoLog.Printf("ProcAddRoomToContactList recv from=%v to=%v cmd=%v seq=%v roomId=%v opUid=%v opType=%v",
		head.From,
		head.To,
		head.Cmd,
		head.Seq,
		roomId,
		opUid,
		opType)
	// first check input param
	if roomId == 0 || opUid == 0 || opType > 1 {
		infoLog.Printf("ProcAddRoomToContactList invalid param")
		result = uint16(ht_muc.MUC_RET_CODE_RET_INPUT_PARAM_ERR)
		errMsg = []byte("input param error")
		return false
	}

	// operator room in contact list
	err := roomManager.AddRoomToContactList(roomId, opUid, opType)
	if err != nil {
		infoLog.Printf("ProcAddRoomToContactList modify faield roomId=%v opUid=%v opType=%v err=%v",
			roomId,
			opUid,
			opType,
			err)
		result = uint16(ht_muc.MUC_RET_CODE_RET_INTERNAL_ERR)
		errMsg = []byte("internal error")
		return false
	}

	// third send ack
	subMucRspBody := new(ht_muc.AddRoomToContactListRspBody)
	subMucRspBody.Status = &ht_muc.MucHeader{Code: proto.Uint32(uint32(ht_muc.MUC_RET_CODE_RET_SUCCESS))}
	SendAddRoomToContactListResp(c, head, subMucRspBody)
	bNeedCall = false
	return true
}
func SendAddRoomToContactListRetCode(c *gotcp.Conn, reqHead *common.HeadV3, ret uint16, errMsg []byte) bool {
	head := new(common.HeadV3)
	if reqHead != nil {
		*head = *reqHead
	}
	rspBody := new(ht_muc.MucRspBody)
	rspBody.AddRoomToContactListRspbody = new(ht_muc.AddRoomToContactListRspBody)
	subMucRspBody := rspBody.GetAddRoomToContactListRspbody()
	subMucRspBody.Status = &ht_muc.MucHeader{Code: proto.Uint32(uint32(ret)), Reason: errMsg}
	s, err := proto.Marshal(rspBody)
	if err != nil {
		infoLog.Printf("SendAddRoomToContactListRetCode proto.Marshal failed from=%v to=%v cmd=%v seq=%v",
			head.From,
			head.To,
			head.Cmd,
			head.Seq)
		return false
	}
	head.Len = uint32(common.PacketV3HeadLen + len(s) + 1) //整个报文长度
	head.Ret = ret
	buf := make([]byte, head.Len)
	buf[0] = common.HTV3MagicBegin
	err = common.SerialHeadV3ToSlice(head, buf[1:])
	if err != nil {
		infoLog.Printf("SendAddRoomToContactListRetCode SerialHeadV3ToSlice failed")
		return false
	}
	copy(buf[common.PacketV3HeadLen:], s) // return code
	buf[head.Len-1] = common.HTV3MagicEnd

	rspPacket := common.NewHeadV3Packet(buf)
	c.AsyncWritePacket(rspPacket, time.Second)
	return true
}
func SendAddRoomToContactListResp(c *gotcp.Conn, reqHead *common.HeadV3, subMucRspBody *ht_muc.AddRoomToContactListRspBody) bool {
	head := new(common.HeadV3)
	if reqHead != nil {
		*head = *reqHead
	}
	rspBody := new(ht_muc.MucRspBody)
	rspBody.AddRoomToContactListRspbody = subMucRspBody
	s, err := proto.Marshal(rspBody)
	if err != nil {
		infoLog.Printf("SendAddRoomToContactListResp proto.Marshal failed from=%v to=%v cmd=%v seq=%v",
			head.From,
			head.To,
			head.Cmd,
			head.Seq)
		return false
	}

	head.Len = uint32(common.PacketV3HeadLen + len(s) + 1) //整个报文长度
	head.Ret = uint16(subMucRspBody.GetStatus().GetCode())
	buf := make([]byte, head.Len)
	buf[0] = common.HTV3MagicBegin
	err = common.SerialHeadV3ToSlice(head, buf[1:])
	if err != nil {
		infoLog.Printf("SendAddRoomToContactListResp SerialHeadV3ToSlice failed")
		return false
	}
	copy(buf[common.PacketV3HeadLen:], s) // return code
	buf[head.Len-1] = common.HTV3MagicEnd

	rspPacket := common.NewHeadV3Packet(buf)
	c.AsyncWritePacket(rspPacket, time.Second)
	return true
}

// 11.Get muc contact list
func ProcGetRoomFromContactList(c *gotcp.Conn, head *common.HeadV3, reqBody *ht_muc.MucReqBody) bool {
	result := uint16(ht_muc.MUC_RET_CODE_RET_SUCCESS)
	errMsg := []byte("operation success")
	bNeedCall := true
	defer func() {
		if bNeedCall {
			SendAddRoomToContactListRetCode(c, head, result, errMsg)
		}
	}()

	// 检查输入参数是否为空
	if head == nil || reqBody == nil {
		bNeedCall = false
		return false
	}

	// add static
	attr := "gomuc/get_room_from_contactlist_req_count"
	libcomm.AttrAdd(attr, 1)

	subReqBody := reqBody.GetGetRoomFromContactListReqbody()
	opUid := subReqBody.GetOpUid()
	infoLog.Printf("ProcGetRoomFromContactList recv from=%v to=%v cmd=%v seq=%v opUid=%v",
		head.From,
		head.To,
		head.Cmd,
		head.Seq,
		opUid)
	// first check input param
	if opUid == 0 {
		infoLog.Printf("ProcGetRoomFromContactList invalid param")
		result = uint16(ht_muc.MUC_RET_CODE_RET_INPUT_PARAM_ERR)
		errMsg = []byte("input param error")
		return false
	}

	// operator room in contact list
	roomList, err := roomManager.GetRoomFromContactList(opUid)
	if err != nil {
		infoLog.Printf("ProcGetRoomFromContactList modify faield opUid=%v err=%v",
			opUid,
			err)
		result = uint16(ht_muc.MUC_RET_CODE_RET_INTERNAL_ERR)
		errMsg = []byte("internal error")
		return false
	}

	infoLog.Printf("DEBUG ProcGetRoomFromContactList roomList=%#v", roomList)
	// third send ack
	subMucRspBody := new(ht_muc.GetRoomFromContactListRspBody)
	subMucRspBody.Status = &ht_muc.MucHeader{Code: proto.Uint32(uint32(ht_muc.MUC_RET_CODE_RET_SUCCESS))}
	subMucRspBody.ListRoomInfo = roomList
	SendGetRoomFromContactListResp(c, head, subMucRspBody)
	bNeedCall = false
	return true
}
func SendGetRoomFromContactListRetCode(c *gotcp.Conn, reqHead *common.HeadV3, ret uint16, errMsg []byte) bool {
	head := new(common.HeadV3)
	if reqHead != nil {
		*head = *reqHead
	}
	rspBody := new(ht_muc.MucRspBody)
	rspBody.GetRoomFromContactListRspbody = new(ht_muc.GetRoomFromContactListRspBody)
	subMucRspBody := rspBody.GetGetRoomFromContactListRspbody()
	subMucRspBody.Status = &ht_muc.MucHeader{Code: proto.Uint32(uint32(ret)), Reason: errMsg}
	s, err := proto.Marshal(rspBody)
	if err != nil {
		infoLog.Printf("SendGetRoomFromContactListRetCode proto.Marshal failed from=%v to=%v cmd=%v seq=%v",
			head.From,
			head.To,
			head.Cmd,
			head.Seq)
		return false
	}
	head.Len = uint32(common.PacketV3HeadLen + len(s) + 1) //整个报文长度
	head.Ret = ret
	buf := make([]byte, head.Len)
	buf[0] = common.HTV3MagicBegin
	err = common.SerialHeadV3ToSlice(head, buf[1:])
	if err != nil {
		infoLog.Printf("SendGetRoomFromContactListRetCode SerialHeadV3ToSlice failed")
		return false
	}
	copy(buf[common.PacketV3HeadLen:], s) // return code
	buf[head.Len-1] = common.HTV3MagicEnd

	rspPacket := common.NewHeadV3Packet(buf)
	c.AsyncWritePacket(rspPacket, time.Second)
	return true
}
func SendGetRoomFromContactListResp(c *gotcp.Conn, reqHead *common.HeadV3, subMucRspBody *ht_muc.GetRoomFromContactListRspBody) bool {
	head := new(common.HeadV3)
	if reqHead != nil {
		*head = *reqHead
	}
	rspBody := new(ht_muc.MucRspBody)
	rspBody.GetRoomFromContactListRspbody = subMucRspBody
	s, err := proto.Marshal(rspBody)
	if err != nil {
		infoLog.Printf("SendGetRoomFromContactListResp proto.Marshal failed from=%v to=%v cmd=%v seq=%v",
			head.From,
			head.To,
			head.Cmd,
			head.Seq)
		return false
	}

	head.Len = uint32(common.PacketV3HeadLen + len(s) + 1) //整个报文长度
	head.Ret = uint16(subMucRspBody.GetStatus().GetCode())
	buf := make([]byte, head.Len)
	buf[0] = common.HTV3MagicBegin
	err = common.SerialHeadV3ToSlice(head, buf[1:])
	if err != nil {
		infoLog.Printf("SendGetRoomFromContactListResp SerialHeadV3ToSlice failed")
		return false
	}
	copy(buf[common.PacketV3HeadLen:], s) // return code
	buf[head.Len-1] = common.HTV3MagicEnd

	rspPacket := common.NewHeadV3Packet(buf)
	c.AsyncWritePacket(rspPacket, time.Second)
	return true
}

// 13.s2svoip block setting
func ProcS2SVoipBlockSetting(c *gotcp.Conn, head *common.HeadV3, reqBody *ht_muc.MucReqBody) bool {
	result := uint16(ht_muc.MUC_RET_CODE_RET_SUCCESS)
	errMsg := []byte("operation success")
	bNeedCall := true
	defer func() {
		if bNeedCall {
			SendS2SVoipBlockSettingRetCode(c, head, result, errMsg)
		}
	}()

	// 检查输入参数是否为空
	if head == nil || reqBody == nil {
		bNeedCall = false
		return false
	}
	// add static
	attr := "gomuc/s2s_voip_block_setting_req_count"
	libcomm.AttrAdd(attr, 1)

	subReqBody := reqBody.GetS2SVoipBlockSettingReqbody()
	blockId := subReqBody.GetBlockId()
	blockType := subReqBody.GetBlockType()
	action := subReqBody.GetAction()
	opUid := head.From
	infoLog.Printf("ProcS2SVoipBlockSetting recv from=%v to=%v cmd=%v seq=%v blockId=%v blockType=%v action=%v",
		head.From,
		head.To,
		head.Cmd,
		head.Seq,
		blockId,
		blockType,
		action)
	// first check input param
	if opUid == 0 || blockId == 0 || action > ht_muc.VOIP_BLOCK_SETTING_ENUM_VOIP_UN_KNOW || blockType > ht_muc.VOIP_BLOCK_TYPE_ENUM_UN_KNOW_TYPE {
		infoLog.Printf("ProcS2SVoipBlockSetting invalid param")
		result = uint16(ht_muc.MUC_RET_CODE_RET_INPUT_PARAM_ERR)
		errMsg = []byte("input param error")
		return false
	}

	// operator room in contact list
	err := roomManager.UpdateVoipBlockSetting(opUid, blockId, uint32(blockType), uint32(action))
	if err != nil {
		infoLog.Printf("ProcS2SVoipBlockSetting modify faield opUid=%v blockId=%v blockType=%v action=%v err=%v",
			opUid,
			blockId,
			blockType,
			action,
			err)
		result = uint16(ht_muc.MUC_RET_CODE_RET_INTERNAL_ERR)
		errMsg = []byte("internal error")
		return false
	}

	infoLog.Printf("ProcS2SVoipBlockSetting succ ret=%v", uint32(ht_muc.MUC_RET_CODE_RET_SUCCESS))
	// third send ack
	subMucRspBody := new(ht_muc.S2SVoipBlockSettingRspBody)
	subMucRspBody.Status = &ht_muc.MucHeader{Code: proto.Uint32(uint32(ht_muc.MUC_RET_CODE_RET_SUCCESS))}
	SendS2SVoipBlockSettingResp(c, head, subMucRspBody)
	bNeedCall = false
	return true
}
func SendS2SVoipBlockSettingRetCode(c *gotcp.Conn, reqHead *common.HeadV3, ret uint16, errMsg []byte) bool {
	head := new(common.HeadV3)
	if reqHead != nil {
		*head = *reqHead
	}
	rspBody := new(ht_muc.MucRspBody)
	rspBody.S2SVoipBlockSettingRspbody = new(ht_muc.S2SVoipBlockSettingRspBody)
	subMucRspBody := rspBody.GetS2SVoipBlockSettingRspbody()
	subMucRspBody.Status = &ht_muc.MucHeader{Code: proto.Uint32(uint32(ret)), Reason: errMsg}
	s, err := proto.Marshal(rspBody)
	if err != nil {
		infoLog.Printf("SendS2SVoipBlockSettingRetCode proto.Marshal failed from=%v to=%v cmd=%v seq=%v",
			head.From,
			head.To,
			head.Cmd,
			head.Seq)
		return false
	}
	head.Len = uint32(common.PacketV3HeadLen + len(s) + 1) //整个报文长度
	head.Ret = ret
	buf := make([]byte, head.Len)
	buf[0] = common.HTV3MagicBegin
	err = common.SerialHeadV3ToSlice(head, buf[1:])
	if err != nil {
		infoLog.Printf("SendS2SVoipBlockSettingRetCode SerialHeadV3ToSlice failed")
		return false
	}
	copy(buf[common.PacketV3HeadLen:], s) // return code
	buf[head.Len-1] = common.HTV3MagicEnd

	rspPacket := common.NewHeadV3Packet(buf)
	c.AsyncWritePacket(rspPacket, time.Second)
	return true
}
func SendS2SVoipBlockSettingResp(c *gotcp.Conn, reqHead *common.HeadV3, subMucRspBody *ht_muc.S2SVoipBlockSettingRspBody) bool {
	head := new(common.HeadV3)
	if reqHead != nil {
		*head = *reqHead
	}
	rspBody := new(ht_muc.MucRspBody)
	rspBody.S2SVoipBlockSettingRspbody = subMucRspBody
	s, err := proto.Marshal(rspBody)
	if err != nil {
		infoLog.Printf("SendS2SVoipBlockSettingResp proto.Marshal failed from=%v to=%v cmd=%v seq=%v",
			head.From,
			head.To,
			head.Cmd,
			head.Seq)
		return false
	}

	head.Len = uint32(common.PacketV3HeadLen + len(s) + 1) //整个报文长度
	head.Ret = uint16(subMucRspBody.GetStatus().GetCode())
	buf := make([]byte, head.Len)
	buf[0] = common.HTV3MagicBegin
	err = common.SerialHeadV3ToSlice(head, buf[1:])
	if err != nil {
		infoLog.Printf("SendS2SVoipBlockSettingResp SerialHeadV3ToSlice failed")
		return false
	}
	copy(buf[common.PacketV3HeadLen:], s) // return code
	buf[head.Len-1] = common.HTV3MagicEnd

	rspPacket := common.NewHeadV3Packet(buf)
	c.AsyncWritePacket(rspPacket, time.Second)
	return true
}

// 14.S2S GVOIP invite broadcast
func ProcS2SGvoipInviteBroadCast(c *gotcp.Conn, head *common.HeadV3, reqBody *ht_muc.MucReqBody) bool {
	// 检查输入参数是否为空
	if head == nil || reqBody == nil {
		return false
	}
	// add static
	attr := "gomuc/s2s_gvoip_begin_bc_req_count"
	libcomm.AttrAdd(attr, 1)

	subReqBody := reqBody.GetGvoipInviteBroadcastReqbody()
	createUid := subReqBody.GetCreateUid()
	createName := subReqBody.GetCreateName()
	roomId := subReqBody.GetRoomId()
	channelId := subReqBody.GetChannelId()
	timeStamp := subReqBody.GetTimestamp()

	infoLog.Printf("ProcS2SGvoipInviteBroadCast recv from=%v to=%v cmd=%v seq=%v createUid=%v createName=%s roomId=%v channelId=%s ts=%v",
		head.From,
		head.To,
		head.Cmd,
		head.Seq,
		createUid,
		createName,
		roomId,
		channelId,
		timeStamp)
	// first check input param
	if createUid == 0 || roomId == 0 {
		infoLog.Printf("ProcS2SGvoipInviteBroadCast invalid param")
		return false
	}
	strMsgId := fmt.Sprintf("MSG_%v", timeStamp)
	pushInfo := util.MucPushInfo{PushType: uint8(PUSH_GVOIP), NickName: string(createName), MsgId: strMsgId}

	// second Broad cast
	bcSlice, err := proto.Marshal(reqBody)
	if err != nil {
		infoLog.Printf("ProcS2SGvoipInviteBroadCast proto.Marshal failed from=%v to=%v cmd=%v seq=%v",
			head.From,
			head.To,
			head.Cmd,
			head.Seq)
		return false
	}
	// build xtpacket for old version
	var packetPayLoad []byte
	common.MarshalUint32(createUid, &packetPayLoad)
	common.MarshalSlice(createName, &packetPayLoad)
	common.MarshalUint32(roomId, &packetPayLoad)
	common.MarshalSlice(channelId, &packetPayLoad)
	common.MarshalUint64(timeStamp, &packetPayLoad)

	err = roomManager.BroadcastMucMessageCompatible(roomId, head, bcSlice, &pushInfo, 0, nil, packetPayLoad)
	if err != nil {
		infoLog.Printf("ProcS2SGvoipInviteBroadCast modify faield creatUid=%v createName=%s roomId=%v channelId=%s timestamp=%v err=%v",
			createUid,
			createName,
			roomId,
			channelId,
			timeStamp,
			err)
		return false
	}
	return true
}

// 15.S2S GVOIP End broadcast
func ProcS2SGvoipEndBroadCast(c *gotcp.Conn, head *common.HeadV3, reqBody *ht_muc.MucReqBody) bool {
	// 检查输入参数是否为空
	if head == nil || reqBody == nil {
		return false
	}
	// add static
	attr := "gomuc/s2s_gvoip_end_bc_req_count"
	libcomm.AttrAdd(attr, 1)

	subReqBody := reqBody.GetGvoipEndBreadcastReqbody()
	roomId := subReqBody.GetRoomId()
	channelId := subReqBody.GetChannelId()
	timeStamp := subReqBody.GetTimestamp()

	infoLog.Printf("ProcS2SGvoipEndBroadCast recv from=%v to=%v cmd=%v seq=%v roomId=%v channelId=%s ts=%v",
		head.From,
		head.To,
		head.Cmd,
		head.Seq,
		roomId,
		channelId,
		timeStamp)
	// first check input param
	if head.To == 0 || roomId == 0 || head.To != roomId {
		infoLog.Printf("ProcS2SGvoipEndBroadCast invalid param")
		return false
	}

	// second Broad cast
	bcSlice, err := proto.Marshal(reqBody)
	if err != nil {
		infoLog.Printf("ProcS2SGvoipEndBroadCast proto.Marshal failed from=%v to=%v cmd=%v seq=%v",
			head.From,
			head.To,
			head.Cmd,
			head.Seq)
		return false
	}

	var packetPayLoad []byte
	common.MarshalUint32(roomId, &packetPayLoad)
	common.MarshalSlice(channelId, &packetPayLoad)
	common.MarshalUint64(timeStamp, &packetPayLoad)

	// second Broad cast
	err = roomManager.BroadcastMucMessageCompatible(roomId, head, bcSlice, nil, 0, nil, packetPayLoad)
	if err != nil {
		infoLog.Printf("ProcS2SGvoipEndBroadCast modify faield roomId=%v channelId=%s timestamp=%v err=%v",
			roomId,
			channelId,
			timeStamp,
			err)
		return false
	}
	return true
}

// 16.S2S GVOIP member change broadcast
func ProcS2SGvoipMemberChangeBroadCast(c *gotcp.Conn, head *common.HeadV3, reqBody *ht_muc.MucReqBody) bool {
	// 检查输入参数是否为空
	if head == nil || reqBody == nil {
		return false
	}
	// add static
	attr := "gomuc/s2s_gvoip_member_change_req_count"
	libcomm.AttrAdd(attr, 1)

	subReqBody := reqBody.GetGvoipMemberChangeBraodcastReqbody()
	roomId := subReqBody.GetRoomId()
	channelId := subReqBody.GetChannelId()
	changeUid := subReqBody.GetChangeUid()
	memberCount := subReqBody.GetMemberCount()
	watcherCount := subReqBody.GetTotalWatcherCount()
	totalList := subReqBody.GetTotalWatcherList()

	infoLog.Printf("ProcS2SGvoipMemberChangeBroadCast recv from=%v to=%v cmd=%v seq=%v roomId=%v channelId=%s changeUid=%v memberCount=%v watcheCount=%v",
		head.From,
		head.To,
		head.Cmd,
		head.Seq,
		roomId,
		channelId,
		changeUid,
		memberCount,
		watcherCount)
	// first check input param
	if head.From == 0 || head.To == 0 || roomId == 0 || head.From != changeUid || head.To != roomId {
		infoLog.Printf("ProcS2SGvoipMemberChangeBroadCast invalid param")
		return false
	}

	bcReqBody := &ht_muc.MucReqBody{}
	subReqBody.TotalWatcherList = nil
	bcReqBody.GvoipMemberChangeBraodcastReqbody = subReqBody
	sliceBody, err := proto.Marshal(bcReqBody)
	if err != nil {
		infoLog.Printf("ProcS2SGvoipMemberChangeBroadCast proto marshal err=%v", err)
		return false
	}

	var packetPayLoad []byte
	common.MarshalUint32(roomId, &packetPayLoad)
	common.MarshalSlice(channelId, &packetPayLoad)
	common.MarshalUint32(changeUid, &packetPayLoad)
	common.MarshalUint32(memberCount, &packetPayLoad)
	common.MarshalUint32(watcherCount, &packetPayLoad)

	// second Broad cast
	err = roomManager.MultiCastNotificationCompatible(totalList, head, sliceBody, true, packetPayLoad)
	if err != nil {
		infoLog.Printf("ProcS2SGvoipMemberChangeBroadCast modify faield roomId=%v channelId=%s err=%v",
			roomId,
			channelId,
			err)
		return false
	}
	return true
}

// 17.proc member join request
func ProcRequestJoinRoom(c *gotcp.Conn, head *common.HeadV3, reqBody *ht_muc.MucReqBody) bool {
	result := uint16(ht_muc.MUC_RET_CODE_RET_SUCCESS)
	errMsg := []byte("operation success")
	bNeedCall := true
	defer func() {
		if bNeedCall {
			SendRequestJoinRoomRetCode(c, head, result, errMsg)
		}
	}()

	// 检查输入参数是否为空
	if head == nil || reqBody == nil {
		bNeedCall = false
		return false
	}
	// add static
	attr := "gomuc/join_room_req_count"
	libcomm.AttrAdd(attr, 1)

	subReqBody := reqBody.GetRequestJoinRoomReqbody()
	roomId := subReqBody.GetRoomId()
	inviter := subReqBody.GetInviterInfo()
	inviteeList := subReqBody.GetInviteeInfo()
	roomIdFrom := subReqBody.GetRoomIdFrom()
	msgId := subReqBody.GetMsgId()
	infoLog.Printf("ProcRequestJoinRoom recv from=%v to=%v cmd=%v seq=%v roomId=%v inviter=%v count=%v roomIdFrom=%v msgId=%s",
		head.From,
		head.To,
		head.Cmd,
		head.Seq,
		roomId,
		inviter.GetUid(),
		len(inviteeList),
		roomIdFrom,
		msgId)

	infoLog.Printf("ProcRequestJoinRoom roomId=%v invter=%v inviteeList=%v",
		roomId,
		inviter.GetUid(),
		inviteeList)

	// first check input param
	if roomId == 0 || inviter.GetUid() == 0 || len(inviteeList) == 0 || len(msgId) == 0 {
		infoLog.Printf("ProcRequestJoinRoom invalid param roomId=%v inviterUid=%v membercount=%v msgId=%s",
			roomId,
			inviter.GetUid(),
			len(inviteeList),
			msgId)
		result = uint16(ht_muc.MUC_RET_CODE_RET_INPUT_PARAM_ERR)
		errMsg = []byte("input param error")
		return false
	}
	// 检查邀请者是否已经被 被邀请者拉黑
	blackMeList, err := DbUtil.GetBlackMeList(inviter.GetUid())
	if err != nil {
		infoLog.Printf("ProcRequestJoinRoom GetBlackMeList failed uid=%v", inviter.GetUid())
	}
	var totalBalck []uint32
	if len(blackMeList) != 0 {
		for i, v := range inviteeList {
			infoLog.Printf("ProcRequestJoinRoom black me index=%v uid=%v", i, v.GetUid())
			if UidIsInSlice(blackMeList, v.GetUid()) {
				infoLog.Printf("ProcRequestJoinRoom uid=%v is in blackMeList", v.GetUid())
				totalBalck = append(totalBalck, v.GetUid())
			}
		}
	}
	if len(totalBalck) != 0 {
		subMucRspBody := new(ht_muc.RequestJoinRoomRspBody)
		errMsg := []byte("some one black me")
		subMucRspBody.Status = &ht_muc.MucHeader{Code: proto.Uint32(uint32(ht_muc.MUC_RET_CODE_RET_SOME_ONE_BLACK_ME)), Reason: errMsg}
		subMucRspBody.ListBlackMe = totalBalck
		SendRequestJoinRoomRsp(c, head, subMucRspBody)
		bNeedCall = false
		return false
	}
	// 首先检查是否超过群成员总数限制
	bResult, err := roomManager.IsExceedRoomMemberLimit(roomId, uint32(len(inviteeList)))
	if err != nil {
		infoLog.Printf("ProcRequestJoinRoom internal error roomId=%v err=%v", roomId, err)
		result = uint16(ht_muc.MUC_RET_CODE_RET_INTERNAL_ERR)
		errMsg = []byte("internal error")
		return false
	}
	// 超过群成员总数限制
	if bResult {
		infoLog.Printf("ProcRequestJoinRoom exceed member limit roomId=%v", roomId)
		result = uint16(ht_muc.MUC_RET_CODE_RET_MEMBER_EXEC_LIMIT)
		errMsg = []byte("member exceed limit")
		return false
	}

	// 检查邀请者是否是管理员如果是管理员直接加群并广播
	// 获取群聊确认
	bIsOpenVerify, err := roomManager.IsOpenVerify(roomId)
	if err != nil {
		infoLog.Printf("ProcRequestJoinRoom roomManager.IsOpenVerify exec failed roomId=%v", roomId)
		result = uint16(ht_muc.MUC_RET_CODE_RET_INTERNAL_ERR)
		errMsg = []byte("get verify stat failed")
		return false
	}
	bAdmin, err := roomManager.IsUserAdmin(roomId, inviter.GetUid())
	if err != nil {
		infoLog.Printf("ProcRequestJoinRoom roomManager.IsUserAdmin exec failed roomId=%v", roomId)
		result = uint16(ht_muc.MUC_RET_CODE_RET_INTERNAL_ERR)
		errMsg = []byte("check is admin failed")
		return false
	}

	if !bIsOpenVerify || bAdmin {
		roomTS, err := roomManager.AddMember(roomId, inviter, inviteeList)
		if err != nil {
			if err == util.ErrAlreadyIn {
				// 返回响应加群成功
				subMucRspBody := new(ht_muc.RequestJoinRoomRspBody)
				subMucRspBody.Status = &ht_muc.MucHeader{Code: proto.Uint32(uint32(ht_muc.MUC_RET_CODE_RET_SUCCESS))}
				subMucRspBody.RoomTimestamp = proto.Uint64(uint64(roomTS))
				SendRequestJoinRoomRsp(c, head, subMucRspBody)
				bNeedCall = false
				return true
			} else {
				infoLog.Printf("ProcRequestJoinRoom internal error roomId=%v err=%v", roomId, err)
				result = uint16(ht_muc.MUC_RET_CODE_RET_INTERNAL_ERR)
				errMsg = []byte("internal error")
				return false
			}
		}
		err = roomManager.NotifyInviteMember(inviter, inviter, inviteeList, roomId, uint64(roomTS), roomIdFrom)
		if err != nil {
			infoLog.Printf("ProcRequestJoinRoom NotifyInviteMember roomId=%v inviterId=%v inviteeListCount=%v failed err=%v",
				roomId,
				inviter.GetUid(),
				len(inviteeList),
				err)
		}

		infoLog.Printf("DEBUG ProcRequestJoinRoom roomId=%v roomTS=%v", roomId, roomTS)
		// 返回响应加群成功
		subMucRspBody := new(ht_muc.RequestJoinRoomRspBody)
		subMucRspBody.Status = &ht_muc.MucHeader{Code: proto.Uint32(uint32(ht_muc.MUC_RET_CODE_RET_SUCCESS))}
		subMucRspBody.RoomTimestamp = proto.Uint64(uint64(roomTS))
		SendRequestJoinRoomRsp(c, head, subMucRspBody)
		bNeedCall = false
	} else { // 检查失败或者不是管理员
		// 首先将请求转发到管理员
		// second Broad cast
		reqSlice, err := proto.Marshal(reqBody)
		if err != nil {
			infoLog.Printf("ProcRequestJoinRoom proto.Marshal failed from=%v to=%v cmd=%v seq=%v",
				head.From,
				head.To,
				head.Cmd,
				head.Seq)
			result = uint16(ht_muc.MUC_RET_CODE_RET_INTERNAL_ERR)
			errMsg = []byte("internal error")
			return false
		}
		roomTS, err := roomManager.NotifyAdminRequestJoin(roomId, head, reqSlice)
		if err != nil {
			infoLog.Printf("ProcRequestJoinRoom internal error roomId=%v err=%v", roomId, err)
			result = uint16(ht_muc.MUC_RET_CODE_RET_INTERNAL_ERR)
			errMsg = []byte("internal error")
			return false
		}

		infoLog.Printf("DEBUG ProcRequestJoinRoom roomId=%v roomTS=%v", roomId, roomTS)
		// 返回响应请求成功
		subMucRspBody := new(ht_muc.RequestJoinRoomRspBody)
		subMucRspBody.Status = &ht_muc.MucHeader{Code: proto.Uint32(uint32(ht_muc.MUC_RET_CODE_RET_WAIT_ADMIN_VERIFY))}
		subMucRspBody.RoomTimestamp = proto.Uint64(uint64(roomTS))
		SendRequestJoinRoomRsp(c, head, subMucRspBody)
		bNeedCall = false
	}

	return true
}
func SendRequestJoinRoomRetCode(c *gotcp.Conn, reqHead *common.HeadV3, ret uint16, errMsg []byte) bool {
	head := new(common.HeadV3)
	if reqHead != nil {
		*head = *reqHead
	}
	rspBody := new(ht_muc.MucRspBody)
	rspBody.RequestJoinRoomRspbody = new(ht_muc.RequestJoinRoomRspBody)
	subMucRspBody := rspBody.GetRequestJoinRoomRspbody()
	subMucRspBody.Status = &ht_muc.MucHeader{Code: proto.Uint32(uint32(ret)), Reason: errMsg}
	s, err := proto.Marshal(rspBody)
	if err != nil {
		infoLog.Printf("SendRequestJoinRoomRetCode proto.Marshal failed from=%v to=%v cmd=%v seq=%v",
			head.From,
			head.To,
			head.Cmd,
			head.Seq)
		return false
	}
	head.Len = uint32(common.PacketV3HeadLen + len(s) + 1) //整个报文长度
	head.Ret = ret
	buf := make([]byte, head.Len)
	buf[0] = common.HTV3MagicBegin
	err = common.SerialHeadV3ToSlice(head, buf[1:])
	if err != nil {
		infoLog.Printf("SendRequestJoinRoomRetCode SerialHeadV3ToSlice failed")
		return false
	}
	copy(buf[common.PacketV3HeadLen:], s) // return code
	buf[head.Len-1] = common.HTV3MagicEnd

	rspPacket := common.NewHeadV3Packet(buf)
	c.AsyncWritePacket(rspPacket, time.Second)
	return true
}
func SendRequestJoinRoomRsp(c *gotcp.Conn, reqHead *common.HeadV3, subMucRspBody *ht_muc.RequestJoinRoomRspBody) bool {
	head := new(common.HeadV3)
	if reqHead != nil {
		*head = *reqHead
	}
	rspBody := new(ht_muc.MucRspBody)
	rspBody.RequestJoinRoomRspbody = subMucRspBody
	s, err := proto.Marshal(rspBody)
	if err != nil {
		infoLog.Printf("SendRequestJoinRoomRsp proto.Marshal failed from=%v to=%v cmd=%v seq=%v",
			head.From,
			head.To,
			head.Cmd,
			head.Seq)
		return false
	}

	head.Len = uint32(common.PacketV3HeadLen + len(s) + 1) //整个报文长度
	head.Ret = uint16(subMucRspBody.GetStatus().GetCode())
	buf := make([]byte, head.Len)
	buf[0] = common.HTV3MagicBegin
	err = common.SerialHeadV3ToSlice(head, buf[1:])
	if err != nil {
		infoLog.Printf("SendRequestJoinRoomRsp SerialHeadV3ToSlice failed")
		return false
	}
	copy(buf[common.PacketV3HeadLen:], s) // return code
	buf[head.Len-1] = common.HTV3MagicEnd

	rspPacket := common.NewHeadV3Packet(buf)
	c.AsyncWritePacket(rspPacket, time.Second)
	return true
}

// 18.proc muc open verify
func ProcMucOpenVerify(c *gotcp.Conn, head *common.HeadV3, reqBody *ht_muc.MucReqBody) bool {
	result := uint16(ht_muc.MUC_RET_CODE_RET_SUCCESS)
	errMsg := []byte("operation success")
	bNeedCall := true
	defer func() {
		if bNeedCall {
			SendMucOpenVerifyRetCode(c, head, result, errMsg)
		}
	}()

	// 检查输入参数是否为空
	if head == nil || reqBody == nil {
		bNeedCall = false
		return false
	}

	// add static
	attr := "gomuc/open_verify_req_count"
	libcomm.AttrAdd(attr, 1)

	subReqBody := reqBody.GetMucOpenVerifyReqbody()
	roomId := subReqBody.GetRoomId()
	reqUid := subReqBody.GetReqUid()
	opStat := subReqBody.GetOpType()
	infoLog.Printf("ProcMucOpenVerify recv from=%v to=%v cmd=%v seq=%v roomId=%v reqUid=%v opStat=%v",
		head.From,
		head.To,
		head.Cmd,
		head.Seq,
		roomId,
		reqUid,
		opStat)

	// first check input param
	if roomId == 0 || reqUid == 0 {
		infoLog.Printf("ProcMucOpenVerify invalid param")
		result = uint16(ht_muc.MUC_RET_CODE_RET_INPUT_PARAM_ERR)
		errMsg = []byte("input param error")
		return false
	}
	// 首先检查群求者是否是群创建者
	bResult, err := roomManager.IsCreateUid(roomId, reqUid)
	if err != nil {
		infoLog.Printf("ProcMucOpenVerify roomManager.IsCreateUid() failed roomId=%v reqUid=%v",
			roomId,
			reqUid)
		result = uint16(ht_muc.MUC_RET_CODE_RET_INTERNAL_ERR)
		errMsg = []byte("Internal error")
		return false
	}

	if !bResult {
		// 操作用户不是创建者 返回错误
		infoLog.Printf("ProcMucOpenVerify roomId=%v reqUid=%v is not creater",
			roomId,
			reqUid)
		result = uint16(ht_muc.MUC_RET_CODE_RET_PERMISSION_DENIED)
		errMsg = []byte("Internal error")
		return false
	}
	// 否则设置群验证状态
	roomTS, err := roomManager.UpdateVerifyStat(roomId, uint32(opStat))
	if err != nil {
		infoLog.Printf("ProcMucOpenVerify UpdateVerifyStat failed roomId=%v reqUid=%v opStat=%v",
			roomId,
			reqUid,
			opStat)
		result = uint16(ht_muc.MUC_RET_CODE_RET_INTERNAL_ERR)
		errMsg = []byte("Internal error")
		return false
	}
	err = roomManager.NotifyOpenVerify(roomId, reqUid, uint32(opStat), roomTS)
	if err != nil {
		infoLog.Printf("ProcMucOpenVerify NotifyOpenVerify failed roomId=%v reqUid=%v opStat=%v",
			roomId,
			reqUid,
			opStat)
	}

	infoLog.Printf("DEBUG ProcMucOpenVerify roomId=%v roomTS=%v", roomId, roomTS)
	// 返回响应请求成功
	subMucRspBody := new(ht_muc.MucOpenVerifyRspBody)
	subMucRspBody.Status = &ht_muc.MucHeader{Code: proto.Uint32(uint32(ht_muc.MUC_RET_CODE_RET_SUCCESS))}
	subMucRspBody.RoomTimestamp = proto.Uint64(uint64(roomTS))
	SendMucOpenVerifyRsp(c, head, subMucRspBody)
	bNeedCall = false
	return true
}

func SendMucOpenVerifyRetCode(c *gotcp.Conn, reqHead *common.HeadV3, ret uint16, errMsg []byte) bool {
	head := new(common.HeadV3)
	if reqHead != nil {
		*head = *reqHead
	}
	rspBody := new(ht_muc.MucRspBody)
	rspBody.MucOpenVerifyRspbody = new(ht_muc.MucOpenVerifyRspBody)
	subMucRspBody := rspBody.GetMucOpenVerifyRspbody()
	subMucRspBody.Status = &ht_muc.MucHeader{Code: proto.Uint32(uint32(ret)), Reason: errMsg}
	s, err := proto.Marshal(rspBody)
	if err != nil {
		infoLog.Printf("SendMucOpenVerifyRetCode proto.Marshal failed from=%v to=%v cmd=%v seq=%v",
			head.From,
			head.To,
			head.Cmd,
			head.Seq)
		return false
	}
	head.Len = uint32(common.PacketV3HeadLen + len(s) + 1) //整个报文长度
	head.Ret = ret
	buf := make([]byte, head.Len)
	buf[0] = common.HTV3MagicBegin
	err = common.SerialHeadV3ToSlice(head, buf[1:])
	if err != nil {
		infoLog.Printf("SendMucOpenVerifyRetCode SerialHeadV3ToSlice failed")
		return false
	}
	copy(buf[common.PacketV3HeadLen:], s) // return code
	buf[head.Len-1] = common.HTV3MagicEnd

	rspPacket := common.NewHeadV3Packet(buf)
	c.AsyncWritePacket(rspPacket, time.Second)
	return true
}
func SendMucOpenVerifyRsp(c *gotcp.Conn, reqHead *common.HeadV3, subMucRspBody *ht_muc.MucOpenVerifyRspBody) bool {
	head := new(common.HeadV3)
	if reqHead != nil {
		*head = *reqHead
	}
	rspBody := new(ht_muc.MucRspBody)
	rspBody.MucOpenVerifyRspbody = subMucRspBody
	s, err := proto.Marshal(rspBody)
	if err != nil {
		infoLog.Printf("SendMucOpenVerifyRsp proto.Marshal failed from=%v to=%v cmd=%v seq=%v",
			head.From,
			head.To,
			head.Cmd,
			head.Seq)
		return false
	}

	head.Len = uint32(common.PacketV3HeadLen + len(s) + 1) //整个报文长度
	head.Ret = uint16(subMucRspBody.GetStatus().GetCode())
	buf := make([]byte, head.Len)
	buf[0] = common.HTV3MagicBegin
	err = common.SerialHeadV3ToSlice(head, buf[1:])
	if err != nil {
		infoLog.Printf("SendMucOpenVerifyRsp SerialHeadV3ToSlice failed")
		return false
	}
	copy(buf[common.PacketV3HeadLen:], s) // return code
	buf[head.Len-1] = common.HTV3MagicEnd

	rspPacket := common.NewHeadV3Packet(buf)
	c.AsyncWritePacket(rspPacket, time.Second)
	return true
}

// 19.proc admin process join room request
func ProcMucJoinRoomHandle(c *gotcp.Conn, head *common.HeadV3, reqBody *ht_muc.MucReqBody) bool {
	result := uint16(ht_muc.MUC_RET_CODE_RET_SUCCESS)
	errMsg := []byte("operation success")
	bNeedCall := true
	defer func() {
		if bNeedCall {
			SendMucJoinRoomHandleRetCode(c, head, result, errMsg)
		}
	}()

	// 检查输入参数是否为空
	if head == nil || reqBody == nil {
		bNeedCall = false
		return false
	}

	// add static
	attr := "gomuc/join_room_handle_req_count"
	libcomm.AttrAdd(attr, 1)

	subReqBody := reqBody.GetMcuJoinRoomHandleReqbody()
	roomId := subReqBody.GetRoomId()
	opUidInfo := subReqBody.GetOpUidInfo()
	inviterInfo := subReqBody.GetInviterInfo()
	inviteeInfo := subReqBody.GetInviteeInfo()
	handleType := subReqBody.GetHandleType()
	roomIdFrom := subReqBody.GetRoomIdFrom()
	msgId := subReqBody.GetMsgId()
	infoLog.Printf("ProcMucJoinRoomHandle recv from=%v to=%v cmd=%v seq=%v roomId=%v opUid=%v inviter=%v invitee=%v type=%v roomIdFrom=%v msgId=%s",
		head.From,
		head.To,
		head.Cmd,
		head.Seq,
		roomId,
		opUidInfo.GetUid(),
		inviterInfo.GetUid(),
		inviteeInfo.GetUid(),
		handleType,
		roomIdFrom,
		msgId)

	// first check input param
	if roomId == 0 || opUidInfo.GetUid() == 0 || inviterInfo.GetUid() == 0 || inviteeInfo.GetUid() == 0 || handleType > ht_muc.HANDLE_OP_TYPE_ENUM_HANDLE_ACCEPT {
		infoLog.Printf("ProcMucJoinRoomHandle invalid param")
		result = uint16(ht_muc.MUC_RET_CODE_RET_INPUT_PARAM_ERR)
		errMsg = []byte("input param error")
		return false
	}

	// 然后检查操作者是否是管理员
	bRet, err := roomManager.IsUserAdmin(roomId, opUidInfo.GetUid())
	if err != nil {
		infoLog.Printf("ProcMucJoinRoomHandle exec IsUserAdmin failed roomId=%v opUid=%v", roomId, opUidInfo.GetUid())
		result = uint16(ht_muc.MUC_RET_CODE_RET_INPUT_PARAM_ERR)
		errMsg = []byte("input param error")
		return false
	}

	if !bRet {
		infoLog.Printf("ProcMucJoinRoomHandle roomId=%v opUid=%v is not admin", roomId, opUidInfo.GetUid())
		result = uint16(ht_muc.MUC_RET_CODE_RET_PERMISSION_DENIED)
		errMsg = []byte("permission denied")
		return false
	}

	// 是管理员 继续处理结果
	var roomTS int64 = 0
	switch handleType {
	case ht_muc.HANDLE_OP_TYPE_ENUM_HANDLE_REJECT:
		{ // 被拒绝 需要通知邀请者和其他管理员
			var targetList []uint32
			// 获取群资料时间戳
			roomInfo, err := roomManager.GetRoomInfo(roomId)
			if err != nil {
				infoLog.Printf("ProcMucJoinRoomHandle GetRoomInfo roomId=%v uid=%v failed",
					roomId,
					inviterInfo.GetUid())
				result = uint16(ht_muc.MUC_RET_CODE_RET_INTERNAL_ERR)
				errMsg = []byte("Internal error")
				return false
			}

			if opUidInfo.GetUid() != roomInfo.CreateUid {
				targetList = append(targetList, roomInfo.CreateUid)
			}
			adminList := roomInfo.AdminList
			for _, v := range adminList {
				if opUidInfo.GetUid() != v {
					targetList = append(targetList, v)
				}
			}

			if roomIdFrom == uint32(ht_muc.ROOMID_FROM_TYPE_ENUM_FROM_INVITE) {
				targetList = append(targetList, inviterInfo.GetUid())
			} else if roomIdFrom == uint32(ht_muc.ROOMID_FROM_TYPE_ENUM_FROM_QRCODE) {
				targetList = append(targetList, inviteeInfo.GetUid())
			}

			reqSlice, err := proto.Marshal(reqBody)
			if err != nil {
				infoLog.Printf("ProcMucJoinRoomHandle proto.Marshal failed from=%v to=%v cmd=%v seq=%v",
					head.From,
					head.To,
					head.Cmd,
					head.Seq)
				result = uint16(ht_muc.MUC_RET_CODE_RET_INTERNAL_ERR)
				errMsg = []byte("internal error")
				return false
			}
			err = roomManager.MultiCastNotification(targetList, head, reqSlice, false)
			if err != nil {
				infoLog.Printf("ProcMucJoinRoomHandle MultiCastNotification roomId=%v count=%v failed",
					roomId,
					len(targetList))
				result = uint16(ht_muc.MUC_RET_CODE_RET_INTERNAL_ERR)
				errMsg = []byte("Internal error")
				return false
			}
			roomTS = roomInfo.RoomTS
		}
	case ht_muc.HANDLE_OP_TYPE_ENUM_HANDLE_ACCEPT:
		{ // 接受用户 首先将用户加入群聊 然后进行广播
			addList := []*ht_muc.RoomMemberInfo{inviteeInfo}
			roomTS, err = roomManager.AddMember(roomId, opUidInfo, addList)
			if err != nil {
				if err == util.ErrExecLimit {
					infoLog.Printf("ProcMucJoinRoomHandle error roomId=%v err=%v", roomId, err)
					result = uint16(ht_muc.MUC_RET_CODE_RET_MEMBER_EXEC_LIMIT)
					errMsg = []byte("exceed member limit")
				} else {
					infoLog.Printf("ProcMucJoinRoomHandle internal error roomId=%v err=%v", roomId, err)
					result = uint16(ht_muc.MUC_RET_CODE_RET_INTERNAL_ERR)
					errMsg = []byte("internal error")
				}
				return false
			}
			err = roomManager.NotifyAdminPromotJoin(opUidInfo, inviterInfo, addList, roomId, uint64(roomTS), roomIdFrom, msgId)
			if err != nil {
				infoLog.Printf("ProcMucJoinRoomHandle NotifyInviteMember roomId=%v inviterId=%v invitee=%v failed err=%v",
					roomId,
					opUidInfo.GetUid(),
					inviteeInfo.GetUid(),
					err)

			}
		}
	default:
		infoLog.Printf("ProcMucJoinRoomHandle roomId=%v opUid=%v reqUid=%v unKnow tyep=%v",
			roomId,
			opUidInfo.GetUid(),
			inviteeInfo.GetUid(),
			handleType)
	}

	infoLog.Printf("DEBUG ProcMucJoinRoomHandle roomId=%v roomTS=%v", roomId, roomTS)
	// sixth send ack
	subMucRspBody := new(ht_muc.MucJoinRoomHandleRspBody)
	subMucRspBody.Status = &ht_muc.MucHeader{Code: proto.Uint32(uint32(ht_muc.MUC_RET_CODE_RET_SUCCESS))}
	subMucRspBody.RoomTimestamp = proto.Uint64(uint64(roomTS))
	SendMucJoinRoomHandleResp(c, head, subMucRspBody)
	bNeedCall = false
	return true
}
func SendMucJoinRoomHandleRetCode(c *gotcp.Conn, reqHead *common.HeadV3, ret uint16, errMsg []byte) bool {
	head := new(common.HeadV3)
	if reqHead != nil {
		*head = *reqHead
	}
	rspBody := new(ht_muc.MucRspBody)
	rspBody.McuJoinRoomHandleRspbody = new(ht_muc.MucJoinRoomHandleRspBody)
	subMucRspBody := rspBody.GetMcuJoinRoomHandleRspbody()
	subMucRspBody.Status = &ht_muc.MucHeader{Code: proto.Uint32(uint32(ret)), Reason: errMsg}
	s, err := proto.Marshal(rspBody)
	if err != nil {
		infoLog.Printf("SendMucJoinRoomHandleRetCode proto.Marshal failed from=%v to=%v cmd=%v seq=%v",
			head.From,
			head.To,
			head.Cmd,
			head.Seq)
		return false
	}
	head.Len = uint32(common.PacketV3HeadLen + len(s) + 1) //整个报文长度
	head.Ret = ret
	buf := make([]byte, head.Len)
	buf[0] = common.HTV3MagicBegin
	err = common.SerialHeadV3ToSlice(head, buf[1:])
	if err != nil {
		infoLog.Printf("SendMucJoinRoomHandleRetCode SerialHeadV3ToSlice failed")
		return false
	}
	copy(buf[common.PacketV3HeadLen:], s) // return code
	buf[head.Len-1] = common.HTV3MagicEnd

	rspPacket := common.NewHeadV3Packet(buf)
	c.AsyncWritePacket(rspPacket, time.Second)
	return true
}
func SendMucJoinRoomHandleResp(c *gotcp.Conn, reqHead *common.HeadV3, subMucRspBody *ht_muc.MucJoinRoomHandleRspBody) bool {
	head := new(common.HeadV3)
	if reqHead != nil {
		*head = *reqHead
	}
	rspBody := new(ht_muc.MucRspBody)
	rspBody.McuJoinRoomHandleRspbody = subMucRspBody
	s, err := proto.Marshal(rspBody)
	if err != nil {
		infoLog.Printf("SendMucJoinRoomHandleResp proto.Marshal failed from=%v to=%v cmd=%v seq=%v",
			head.From,
			head.To,
			head.Cmd,
			head.Seq)
		return false
	}

	head.Len = uint32(common.PacketV3HeadLen + len(s) + 1) //整个报文长度
	head.Ret = uint16(subMucRspBody.GetStatus().GetCode())
	buf := make([]byte, head.Len)
	buf[0] = common.HTV3MagicBegin
	err = common.SerialHeadV3ToSlice(head, buf[1:])
	if err != nil {
		infoLog.Printf("SendMucJoinRoomHandleResp SerialHeadV3ToSlice failed")
		return false
	}
	copy(buf[common.PacketV3HeadLen:], s) // return code
	buf[head.Len-1] = common.HTV3MagicEnd

	rspPacket := common.NewHeadV3Packet(buf)
	c.AsyncWritePacket(rspPacket, time.Second)
	return true
}

// 20.proc set admind req
func ProcMucSetAdmin(c *gotcp.Conn, head *common.HeadV3, reqBody *ht_muc.MucReqBody) bool {
	result := uint16(ht_muc.MUC_RET_CODE_RET_SUCCESS)
	errMsg := []byte("operation success")
	bNeedCall := true
	defer func() {
		if bNeedCall {
			SendMucSetAdminRetCode(c, head, result, errMsg)
		}
	}()

	// 检查输入参数是否为空
	if head == nil || reqBody == nil {
		bNeedCall = false
		return false
	}

	// add static
	attr := "gomuc/muc_set_admin_req_count"
	libcomm.AttrAdd(attr, 1)

	subReqBody := reqBody.GetMucSetAdminReqbody()
	roomId := subReqBody.GetRoomId()
	opUid := subReqBody.GetOpUid()
	adminList := subReqBody.GetMembers()
	infoLog.Printf("ProcMucSetAdmin recv from=%v to=%v cmd=%v seq=%v roomId=%v opUid=%v adminSize=%v",
		head.From,
		head.To,
		head.Cmd,
		head.Seq,
		roomId,
		opUid,
		len(adminList))
	// first check input param
	if roomId == 0 || opUid == 0 {
		infoLog.Printf("ProcMucSetAdmin invalid param")
		result = uint16(ht_muc.MUC_RET_CODE_RET_INPUT_PARAM_ERR)
		errMsg = []byte("input param error")
		return false
	}

	// second check whether op user is admin
	bRet, err := roomManager.IsCreateUid(roomId, opUid)
	if err != nil {
		infoLog.Printf("ProcMucSetAdmin exec IsCreateUid failed roomId=%v opUid=%v err=%v", roomId, opUid, err)
		result = uint16(ht_muc.MUC_RET_CODE_RET_INTERNAL_ERR)
		errMsg = []byte("Internal error")
		return false
	}
	// third the op uid is not admin send err rsp
	if !bRet {
		infoLog.Printf("ProcMucSetAdmin user is not creater roomId=%v opUid=%v err=%v", roomId, opUid, err)
		result = uint16(ht_muc.MUC_RET_CODE_RET_PERMISSION_DENIED)
		errMsg = []byte("Internal error")
		return false
	}

	// // second get room Info
	// roomInfo, err := roomManager.GetRoomInfo(roomId)
	// if err != nil {
	// 	infoLog.Printf("ProcMucSetAdmin roomId=%v opUid=%v get room info failed", roomId, opUid)
	// 	result = uint16(ht_muc.MUC_RET_CODE_RET_INTERNAL_ERR)
	// 	errMsg = []byte("Internal error")
	// 	return false
	// }
	// // get prev admin list
	// prevAdminList := roomInfo.AdminList

	// forth set new admin list
	roomTS, err := roomManager.SetAdminList(roomId, opUid, adminList)
	if err != nil {
		infoLog.Printf("ProcMucSetAdmin roomManager.SetAdminList failed roomId=%v opUid=%v adminCount=%v",
			roomId,
			opUid,
			len(adminList))
		result = uint16(ht_muc.MUC_RET_CODE_RET_INTERNAL_ERR)
		errMsg = []byte("internal error")
		return false
	}
	// var notifyList []*ht_muc.RoomMemberInfo
	// for _, v := range adminList {
	// 	if UidIsInSlice(prevAdminList, v.GetUid()) {
	// 		infoLog.Printf("ProcMucSetAdmin uid=%v is already admin not notify", v.GetUid())
	// 	} else {
	// 		notifyList = append(notifyList, v)
	// 	}
	// }
	// fifth broadcast notifycation
	err = roomManager.NotifySetAdmin(roomId, opUid, adminList, roomTS)
	if err != nil {
		infoLog.Printf("ProcMucSetAdmin roomManager.NotifySetAdmin failed roomId=%v opUid=%v",
			roomId,
			opUid)
	}
	infoLog.Printf("DEBUG ProcMucSetAdmin roomId=%v roomTS=%v", roomId, roomTS)
	// sixth send ack
	subMucRspBody := new(ht_muc.MucSetAdminRspBody)
	subMucRspBody.Status = &ht_muc.MucHeader{Code: proto.Uint32(uint32(ht_muc.MUC_RET_CODE_RET_SUCCESS))}
	subMucRspBody.RoomTimestamp = proto.Uint64(uint64(roomTS))
	SendMucSetAdminResp(c, head, subMucRspBody)
	bNeedCall = false
	return true
}
func SendMucSetAdminRetCode(c *gotcp.Conn, reqHead *common.HeadV3, ret uint16, errMsg []byte) bool {
	head := new(common.HeadV3)
	if reqHead != nil {
		*head = *reqHead
	}
	rspBody := new(ht_muc.MucRspBody)
	rspBody.MucSetAdminRspbody = new(ht_muc.MucSetAdminRspBody)
	subMucRspBody := rspBody.GetMucSetAdminRspbody()
	subMucRspBody.Status = &ht_muc.MucHeader{Code: proto.Uint32(uint32(ret)), Reason: errMsg}
	s, err := proto.Marshal(rspBody)
	if err != nil {
		infoLog.Printf("SendMucSetAdminRetCode proto.Marshal failed from=%v to=%v cmd=%v seq=%v",
			head.From,
			head.To,
			head.Cmd,
			head.Seq)
		return false
	}
	head.Len = uint32(common.PacketV3HeadLen + len(s) + 1) //整个报文长度
	head.Ret = ret
	buf := make([]byte, head.Len)
	buf[0] = common.HTV3MagicBegin
	err = common.SerialHeadV3ToSlice(head, buf[1:])
	if err != nil {
		infoLog.Printf("SendMucSetAdminRetCode SerialHeadV3ToSlice failed")
		return false
	}
	copy(buf[common.PacketV3HeadLen:], s) // return code
	buf[head.Len-1] = common.HTV3MagicEnd

	rspPacket := common.NewHeadV3Packet(buf)
	c.AsyncWritePacket(rspPacket, time.Second)
	return true
}
func SendMucSetAdminResp(c *gotcp.Conn, reqHead *common.HeadV3, subMucRspBody *ht_muc.MucSetAdminRspBody) bool {
	head := new(common.HeadV3)
	if reqHead != nil {
		*head = *reqHead
	}
	rspBody := new(ht_muc.MucRspBody)
	rspBody.MucSetAdminRspbody = subMucRspBody
	s, err := proto.Marshal(rspBody)
	if err != nil {
		infoLog.Printf("SendMucSetAdminResp proto.Marshal failed from=%v to=%v cmd=%v seq=%v",
			head.From,
			head.To,
			head.Cmd,
			head.Seq)
		return false
	}

	head.Len = uint32(common.PacketV3HeadLen + len(s) + 1) //整个报文长度
	head.Ret = uint16(subMucRspBody.GetStatus().GetCode())
	buf := make([]byte, head.Len)
	buf[0] = common.HTV3MagicBegin
	err = common.SerialHeadV3ToSlice(head, buf[1:])
	if err != nil {
		infoLog.Printf("SendMucSetAdminResp SerialHeadV3ToSlice failed")
		return false
	}
	copy(buf[common.PacketV3HeadLen:], s) // return code
	buf[head.Len-1] = common.HTV3MagicEnd

	rspPacket := common.NewHeadV3Packet(buf)
	c.AsyncWritePacket(rspPacket, time.Second)
	return true
}

// 21.proc create user right transfer
func ProcCreateUserAuthTrans(c *gotcp.Conn, head *common.HeadV3, reqBody *ht_muc.MucReqBody) bool {
	result := uint16(ht_muc.MUC_RET_CODE_RET_SUCCESS)
	errMsg := []byte("operation success")
	bNeedCall := true
	defer func() {
		if bNeedCall {
			SendCreateUserAuthTransRetCode(c, head, result, errMsg)
		}
	}()

	// 检查输入参数是否为空
	if head == nil || reqBody == nil {
		bNeedCall = false
		return false
	}
	// add static
	attr := "gomuc/create_auth_trans_req_count"
	libcomm.AttrAdd(attr, 1)

	subReqBody := reqBody.GetCreateUserTransReqbody()
	roomId := subReqBody.GetRoomId()
	opUid := subReqBody.GetOpUid()
	memberInfo := subReqBody.GetMember()
	infoLog.Printf("ProcCreateUserAuthTrans recv from=%v to=%v cmd=%v seq=%v roomId=%v opUid=%v memeberUid=%v",
		head.From,
		head.To,
		head.Cmd,
		head.Seq,
		roomId,
		opUid,
		memberInfo.GetUid())
	// first check input param
	if roomId == 0 || opUid == 0 || memberInfo.GetUid() == 0 {
		infoLog.Printf("ProcCreateUserAuthTrans invalid param")
		result = uint16(ht_muc.MUC_RET_CODE_RET_INPUT_PARAM_ERR)
		errMsg = []byte("input param error")
		return false
	}

	// second check whether op user is admin
	bRet, err := roomManager.IsCreateUid(roomId, opUid)
	if err != nil {
		infoLog.Printf("ProcCreateUserAuthTrans exec IsCreateUid failed roomId=%v opUid=%v err=%v", roomId, opUid, err)
		result = uint16(ht_muc.MUC_RET_CODE_RET_INTERNAL_ERR)
		errMsg = []byte("Internal error")
		return false
	}
	// third the op uid is not admin send err rsp
	if !bRet {
		infoLog.Printf("ProcCreateUserAuthTrans user is not creater roomId=%v opUid=%v err=%v", roomId, opUid, err)
		result = uint16(ht_muc.MUC_RET_CODE_RET_PERMISSION_DENIED)
		errMsg = []byte("Internal error")
		return false
	}
	// forth set new admin list
	roomTS, err := roomManager.SetCreateUid(roomId, opUid, memberInfo.GetUid())
	if err != nil {
		infoLog.Printf("ProcCreateUserAuthTrans roomManager.SetAdminList failed roomId=%v opUid=%v memberUid=%v",
			roomId,
			opUid,
			memberInfo.GetUid())
		result = uint16(ht_muc.MUC_RET_CODE_RET_INTERNAL_ERR)
		errMsg = []byte("internal error")
		return false
	}
	// fifth broadcast notifycation
	err = roomManager.NotifyCreateUserAuthTrans(roomId, opUid, memberInfo, roomTS)
	if err != nil {
		infoLog.Printf("ProcCreateUserAuthTrans roomManager.NotifyCreateUserAuthTrans failed roomId=%v opUid=%v",
			roomId,
			opUid)
	}

	infoLog.Printf("DEBUG ProcCreateUserAuthTrans roomId=%v roomTS=%v", roomId, roomTS)
	// sixth send ack
	subMucRspBody := new(ht_muc.CreateUserAuthTransRspBody)
	subMucRspBody.Status = &ht_muc.MucHeader{Code: proto.Uint32(uint32(ht_muc.MUC_RET_CODE_RET_SUCCESS))}
	subMucRspBody.RoomTimestamp = proto.Uint64(uint64(roomTS))
	SendCreateUserAuthTransResp(c, head, subMucRspBody)
	bNeedCall = false
	return true
}
func SendCreateUserAuthTransRetCode(c *gotcp.Conn, reqHead *common.HeadV3, ret uint16, errMsg []byte) bool {
	head := new(common.HeadV3)
	if reqHead != nil {
		*head = *reqHead
	}
	rspBody := new(ht_muc.MucRspBody)
	rspBody.CreateUserTransRspbody = new(ht_muc.CreateUserAuthTransRspBody)
	subMucRspBody := rspBody.GetCreateUserTransRspbody()
	subMucRspBody.Status = &ht_muc.MucHeader{Code: proto.Uint32(uint32(ret)), Reason: errMsg}
	s, err := proto.Marshal(rspBody)
	if err != nil {
		infoLog.Printf("SendCreateUserAuthTransRetCode proto.Marshal failed from=%v to=%v cmd=%v seq=%v",
			head.From,
			head.To,
			head.Cmd,
			head.Seq)
		return false
	}
	head.Len = uint32(common.PacketV3HeadLen + len(s) + 1) //整个报文长度
	head.Ret = ret
	buf := make([]byte, head.Len)
	buf[0] = common.HTV3MagicBegin
	err = common.SerialHeadV3ToSlice(head, buf[1:])
	if err != nil {
		infoLog.Printf("SendCreateUserAuthTransRetCode SerialHeadV3ToSlice failed")
		return false
	}
	copy(buf[common.PacketV3HeadLen:], s) // return code
	buf[head.Len-1] = common.HTV3MagicEnd

	rspPacket := common.NewHeadV3Packet(buf)
	c.AsyncWritePacket(rspPacket, time.Second)
	return true
}
func SendCreateUserAuthTransResp(c *gotcp.Conn, reqHead *common.HeadV3, subMucRspBody *ht_muc.CreateUserAuthTransRspBody) bool {
	head := new(common.HeadV3)
	if reqHead != nil {
		*head = *reqHead
	}
	rspBody := new(ht_muc.MucRspBody)
	rspBody.CreateUserTransRspbody = subMucRspBody
	s, err := proto.Marshal(rspBody)
	if err != nil {
		infoLog.Printf("SendCreateUserAuthTransResp proto.Marshal failed from=%v to=%v cmd=%v seq=%v",
			head.From,
			head.To,
			head.Cmd,
			head.Seq)
		return false
	}

	head.Len = uint32(common.PacketV3HeadLen + len(s) + 1) //整个报文长度
	head.Ret = uint16(subMucRspBody.GetStatus().GetCode())
	buf := make([]byte, head.Len)
	buf[0] = common.HTV3MagicBegin
	err = common.SerialHeadV3ToSlice(head, buf[1:])
	if err != nil {
		infoLog.Printf("SendCreateUserAuthTransResp SerialHeadV3ToSlice failed")
		return false
	}
	copy(buf[common.PacketV3HeadLen:], s) // return code
	buf[head.Len-1] = common.HTV3MagicEnd

	rspPacket := common.NewHeadV3Packet(buf)
	c.AsyncWritePacket(rspPacket, time.Second)
	return true
}

// 22.proc Get room base info
func ProcMucGetRoomBaseInfo(c *gotcp.Conn, head *common.HeadV3, reqBody *ht_muc.MucReqBody) bool {
	result := uint16(ht_muc.MUC_RET_CODE_RET_SUCCESS)
	errMsg := []byte("operation success")
	bNeedCall := true
	defer func() {
		if bNeedCall {
			SendMucGetRoomBaseInfoRetCode(c, head, result, errMsg)
		}
	}()

	// 检查输入参数是否为空
	if head == nil || reqBody == nil {
		bNeedCall = false
		return false
	}

	// add static
	attr := "gomuc/get_room_base_info_req_count"
	libcomm.AttrAdd(attr, 1)

	subReqBody := reqBody.GetMucGetRoomBaseInfoReqbody()
	roomId := subReqBody.GetRoomId()
	opUid := subReqBody.GetReqUid()
	infoLog.Printf("ProcMucGetRoomBaseInfo recv from=%v to=%v cmd=%v seq=%v roomId=%v opUid=%v",
		head.From,
		head.To,
		head.Cmd,
		head.Seq,
		roomId,
		opUid)
	// first check input param
	if roomId == 0 || opUid == 0 {
		infoLog.Printf("ProcMucGetRoomBaseInfo invalid param")
		result = uint16(ht_muc.MUC_RET_CODE_RET_INPUT_PARAM_ERR)
		errMsg = []byte("input param error")
		return false
	}

	// second get room Info
	roomInfo, err := roomManager.GetRoomInfo(roomId)
	if err != nil {
		infoLog.Printf("ProcMucGetRoomBaseInfo roomId=%v opUid=%v get room info failed", roomId, opUid)
		result = uint16(ht_muc.MUC_RET_CODE_RET_INTERNAL_ERR)
		errMsg = []byte("Internal error")
		return false
	}

	// sixth send ack
	subMucRspBody := new(ht_muc.MucGetRoomBaseInfoRspBody)
	subMucRspBody.Status = &ht_muc.MucHeader{Code: proto.Uint32(uint32(ht_muc.MUC_RET_CODE_RET_SUCCESS))}
	subMucRspBody.RoomTimestamp = proto.Uint64(uint64(roomInfo.RoomTS))
	verify := ht_muc.VERIFY_STAT(roomInfo.VerifyStat)
	subMucRspBody.BaseInfo = &ht_muc.RoomBaseInfo{
		CreateUid:    proto.Uint32(roomInfo.CreateUid),
		ListAdminUid: roomInfo.AdminList[:],
		RoomLimit:    proto.Uint32(roomInfo.MemberLimit),
		VerifyStat:   &verify,
		Announcement: &ht_muc.AnnoType{
			PublishUid:  proto.Uint32(roomInfo.Announcement.PublishUid),
			PublishTs:   proto.Uint32(roomInfo.Announcement.PublishTS),
			AnnoContent: []byte(roomInfo.Announcement.AnnoContect),
		},
	}
	infoLog.Printf("DEBUG ProcMucGetRoomBaseInfo baseInfo=%#v", subMucRspBody)

	SendMucGetRoomBaseInfoResp(c, head, subMucRspBody)
	bNeedCall = false
	return true
}
func SendMucGetRoomBaseInfoRetCode(c *gotcp.Conn, reqHead *common.HeadV3, ret uint16, errMsg []byte) bool {
	head := new(common.HeadV3)
	if reqHead != nil {
		*head = *reqHead
	}
	rspBody := new(ht_muc.MucRspBody)
	rspBody.MucGetRoomBaseInfoRspbody = new(ht_muc.MucGetRoomBaseInfoRspBody)
	subMucRspBody := rspBody.GetMucGetRoomBaseInfoRspbody()
	subMucRspBody.Status = &ht_muc.MucHeader{Code: proto.Uint32(uint32(ret)), Reason: errMsg}
	s, err := proto.Marshal(rspBody)
	if err != nil {
		infoLog.Printf("SendMucGetRoomBaseInfoRetCode proto.Marshal failed from=%v to=%v cmd=%v seq=%v",
			head.From,
			head.To,
			head.Cmd,
			head.Seq)
		return false
	}
	head.Len = uint32(common.PacketV3HeadLen + len(s) + 1) //整个报文长度
	head.Ret = ret
	buf := make([]byte, head.Len)
	buf[0] = common.HTV3MagicBegin
	err = common.SerialHeadV3ToSlice(head, buf[1:])
	if err != nil {
		infoLog.Printf("SendMucGetRoomBaseInfoRetCode SerialHeadV3ToSlice failed")
		return false
	}
	copy(buf[common.PacketV3HeadLen:], s) // return code
	buf[head.Len-1] = common.HTV3MagicEnd

	rspPacket := common.NewHeadV3Packet(buf)
	c.AsyncWritePacket(rspPacket, time.Second)
	return true
}
func SendMucGetRoomBaseInfoResp(c *gotcp.Conn, reqHead *common.HeadV3, subMucRspBody *ht_muc.MucGetRoomBaseInfoRspBody) bool {
	head := new(common.HeadV3)
	if reqHead != nil {
		*head = *reqHead
	}
	rspBody := new(ht_muc.MucRspBody)
	rspBody.MucGetRoomBaseInfoRspbody = subMucRspBody
	s, err := proto.Marshal(rspBody)
	if err != nil {
		infoLog.Printf("SendMucGetRoomBaseInfoResp proto.Marshal failed from=%v to=%v cmd=%v seq=%v",
			head.From,
			head.To,
			head.Cmd,
			head.Seq)
		return false
	}

	head.Len = uint32(common.PacketV3HeadLen + len(s) + 1) //整个报文长度
	head.Ret = uint16(subMucRspBody.GetStatus().GetCode())
	buf := make([]byte, head.Len)
	buf[0] = common.HTV3MagicBegin
	err = common.SerialHeadV3ToSlice(head, buf[1:])
	if err != nil {
		infoLog.Printf("SendMucGetRoomBaseInfoResp SerialHeadV3ToSlice failed")
		return false
	}
	copy(buf[common.PacketV3HeadLen:], s) // return code
	buf[head.Len-1] = common.HTV3MagicEnd

	rspPacket := common.NewHeadV3Packet(buf)
	c.AsyncWritePacket(rspPacket, time.Second)
	return true
}

// 23.proc create user right transfer
func ProcMucSetRoomAnnouncement(c *gotcp.Conn, head *common.HeadV3, reqBody *ht_muc.MucReqBody) bool {
	result := uint16(ht_muc.MUC_RET_CODE_RET_SUCCESS)
	errMsg := []byte("operation success")
	bNeedCall := true
	defer func() {
		if bNeedCall {
			SendMucSetRoomAnnoRetCode(c, head, result, errMsg)
		}
	}()

	// 检查输入参数是否为空
	if head == nil || reqBody == nil {
		bNeedCall = false
		return false
	}

	// add static
	attr := "gomuc/set_room_anno_req_count"
	libcomm.AttrAdd(attr, 1)

	subReqBody := reqBody.GetMucSetRoomAnnouncementReqbody()
	roomId := subReqBody.GetRoomId()
	anno := subReqBody.GetAnnouncement()
	infoLog.Printf("ProcMucSetRoomAnnouncement recv from=%v to=%v cmd=%v seq=%v roomId=%v opUid=%v ts=%v content=%s",
		head.From,
		head.To,
		head.Cmd,
		head.Seq,
		roomId,
		anno.GetPublishUid(),
		anno.GetPublishTs(),
		anno.GetAnnoContent())
	// first check input param
	if roomId == 0 || anno.GetPublishUid() == 0 {
		infoLog.Printf("ProcMucSetRoomAnnouncement invalid param")
		result = uint16(ht_muc.MUC_RET_CODE_RET_INPUT_PARAM_ERR)
		errMsg = []byte("input param error")
		return false
	}

	// second check whether op user is admin
	bRet, err := roomManager.IsUserAdmin(roomId, anno.GetPublishUid())
	if err != nil {
		infoLog.Printf("ProcMucSetRoomAnnouncement exec IsUserAdmin failed roomId=%v opUid=%v err=%v", roomId, anno.GetPublishUid(), err)
		result = uint16(ht_muc.MUC_RET_CODE_RET_INTERNAL_ERR)
		errMsg = []byte("Internal error")
		return false
	}
	// third the op uid is not admin send err rsp
	if !bRet {
		infoLog.Printf("ProcMucSetRoomAnnouncement user is not creater roomId=%v opUid=%v err=%v", roomId, anno.GetPublishUid(), err)
		result = uint16(ht_muc.MUC_RET_CODE_RET_PERMISSION_DENIED)
		errMsg = []byte("Internal error")
		return false
	}
	// 群公告发布时间一服务端时间戳为准
	anno.PublishTs = proto.Uint32(uint32(time.Now().Unix()))

	// forth set new announcement
	roomTS, err := roomManager.SetAnnouncement(roomId, anno)
	if err != nil {
		infoLog.Printf("ProcMucSetRoomAnnouncement roomManager.SetAdminList failed roomId=%v publishUid=%v publishTs=%v",
			roomId,
			anno.GetPublishUid(),
			anno.GetPublishTs())
		result = uint16(ht_muc.MUC_RET_CODE_RET_INTERNAL_ERR)
		errMsg = []byte("internal error")
		return false
	}
	// fifth broadcast notifycation
	err = roomManager.NotifySetAnnouncement(roomId, anno, roomTS)
	if err != nil {
		infoLog.Printf("ProcMucSetRoomAnnouncement roomManager.NotifySetAnnouncement failed roomId=%v publishUid=%v",
			roomId,
			anno.GetPublishUid())
	}

	// sixth send ack
	subMucRspBody := new(ht_muc.MucSetRoomAnnouncementRspBody)
	subMucRspBody.Status = &ht_muc.MucHeader{Code: proto.Uint32(uint32(ht_muc.MUC_RET_CODE_RET_SUCCESS))}
	subMucRspBody.RoomTimestamp = proto.Uint64(uint64(roomTS))

	infoLog.Printf("DEBUG ProcMucSetRoomAnnouncement roomId=%v roomTS=%v", roomId, roomTS)
	SendMucSetRoomAnnoResp(c, head, subMucRspBody)
	bNeedCall = false
	return true
}
func SendMucSetRoomAnnoRetCode(c *gotcp.Conn, reqHead *common.HeadV3, ret uint16, errMsg []byte) bool {
	head := new(common.HeadV3)
	if reqHead != nil {
		*head = *reqHead
	}
	rspBody := new(ht_muc.MucRspBody)
	rspBody.MucSetRoomAnnouncementRspbody = new(ht_muc.MucSetRoomAnnouncementRspBody)
	subMucRspBody := rspBody.GetMucSetRoomAnnouncementRspbody()
	subMucRspBody.Status = &ht_muc.MucHeader{Code: proto.Uint32(uint32(ret)), Reason: errMsg}
	s, err := proto.Marshal(rspBody)
	if err != nil {
		infoLog.Printf("SendMucSetRoomAnnoRetCode proto.Marshal failed from=%v to=%v cmd=%v seq=%v",
			head.From,
			head.To,
			head.Cmd,
			head.Seq)
		return false
	}
	head.Len = uint32(common.PacketV3HeadLen + len(s) + 1) //整个报文长度
	head.Ret = ret
	buf := make([]byte, head.Len)
	buf[0] = common.HTV3MagicBegin
	err = common.SerialHeadV3ToSlice(head, buf[1:])
	if err != nil {
		infoLog.Printf("SendMucSetRoomAnnoRetCode SerialHeadV3ToSlice failed")
		return false
	}
	copy(buf[common.PacketV3HeadLen:], s) // return code
	buf[head.Len-1] = common.HTV3MagicEnd

	rspPacket := common.NewHeadV3Packet(buf)
	c.AsyncWritePacket(rspPacket, time.Second)
	return true
}

func SendMucSetRoomAnnoResp(c *gotcp.Conn, reqHead *common.HeadV3, subMucRspBody *ht_muc.MucSetRoomAnnouncementRspBody) bool {
	head := new(common.HeadV3)
	if reqHead != nil {
		*head = *reqHead
	}
	rspBody := new(ht_muc.MucRspBody)
	rspBody.MucSetRoomAnnouncementRspbody = subMucRspBody
	s, err := proto.Marshal(rspBody)
	if err != nil {
		infoLog.Printf("SendMucSetRoomAnnoResp proto.Marshal failed from=%v to=%v cmd=%v seq=%v",
			head.From,
			head.To,
			head.Cmd,
			head.Seq)
		return false
	}

	head.Len = uint32(common.PacketV3HeadLen + len(s) + 1) //整个报文长度
	head.Ret = uint16(subMucRspBody.GetStatus().GetCode())
	buf := make([]byte, head.Len)
	buf[0] = common.HTV3MagicBegin
	err = common.SerialHeadV3ToSlice(head, buf[1:])
	if err != nil {
		infoLog.Printf("SendMucSetRoomAnnoResp SerialHeadV3ToSlice failed")
		return false
	}
	copy(buf[common.PacketV3HeadLen:], s) // return code
	buf[head.Len-1] = common.HTV3MagicEnd

	rspPacket := common.NewHeadV3Packet(buf)
	c.AsyncWritePacket(rspPacket, time.Second)
	return true
}

// 24.Proc QueryUserIsAlreadyInRoom
func ProcQueryUserIsAlreadyInRoom(c *gotcp.Conn, head *common.HeadV3, reqBody *ht_muc.MucReqBody) bool {
	result := uint16(ht_muc.MUC_RET_CODE_RET_SUCCESS)
	errMsg := []byte("operation success")
	bNeedCall := true
	defer func() {
		if bNeedCall {
			SendQueryUserIsAlreadyInRoomRetCode(c, head, result, errMsg)
		}
	}()

	// 检查输入参数是否为空
	if head == nil || reqBody == nil {
		bNeedCall = false
		return false
	}

	// add static
	attr := "gomuc/query_user_is_in_req_count"
	libcomm.AttrAdd(attr, 1)

	subReqBody := reqBody.GetMucGetRoomBaseInfoReqbody()
	roomId := subReqBody.GetRoomId()
	opUid := subReqBody.GetReqUid()
	infoLog.Printf("ProcQueryUserIsAlreadyInRoom recv from=%v to=%v cmd=%v seq=%v roomId=%v opUid=%v",
		head.From,
		head.To,
		head.Cmd,
		head.Seq,
		roomId,
		opUid)
	// first check input param
	if roomId == 0 || opUid == 0 {
		infoLog.Printf("ProcQueryUserIsAlreadyInRoom invalid param")
		result = uint16(ht_muc.MUC_RET_CODE_RET_INPUT_PARAM_ERR)
		errMsg = []byte("input param error")
		return false
	}

	// second get room Info
	roomInfo, err := roomManager.GetRoomInfo(roomId)
	if err != nil {
		infoLog.Printf("ProcQueryUserIsAlreadyInRoom roomId=%v opUid=%v get room info failed", roomId, opUid)
		result = uint16(ht_muc.MUC_RET_CODE_RET_INTERNAL_ERR)
		errMsg = []byte("Internal error")
		return false
	}

	var bInRoom bool = false
	memberList := roomInfo.MemberList
	for _, v := range memberList {
		if v.Uid == opUid {
			bInRoom = true
		}
	}

	subMucRspBody := new(ht_muc.QueryUserIsAlreadyInRoomRspBody)
	subMucRspBody.Status = &ht_muc.MucHeader{Code: proto.Uint32(uint32(ht_muc.MUC_RET_CODE_RET_SUCCESS))}
	subMucRspBody.IsInRoom = proto.Bool(bInRoom)

	infoLog.Printf("DEBUG ProcQueryUserIsAlreadyInRoom subMucRspBody=%#v", subMucRspBody)
	SendQueryUserIsAlreadyInRoomResp(c, head, subMucRspBody)
	bNeedCall = false
	return true
}
func SendQueryUserIsAlreadyInRoomRetCode(c *gotcp.Conn, reqHead *common.HeadV3, ret uint16, errMsg []byte) bool {
	head := new(common.HeadV3)
	if reqHead != nil {
		*head = *reqHead
	}
	rspBody := new(ht_muc.MucRspBody)
	rspBody.QueryUserIsAlreadyInRspbody = new(ht_muc.QueryUserIsAlreadyInRoomRspBody)
	subMucRspBody := rspBody.GetQueryUserIsAlreadyInRspbody()
	subMucRspBody.Status = &ht_muc.MucHeader{Code: proto.Uint32(uint32(ret)), Reason: errMsg}
	s, err := proto.Marshal(rspBody)
	if err != nil {
		infoLog.Printf("SendQueryUserIsAlreadyInRoomRetCode proto.Marshal failed from=%v to=%v cmd=%v seq=%v",
			head.From,
			head.To,
			head.Cmd,
			head.Seq)
		return false
	}
	head.Len = uint32(common.PacketV3HeadLen + len(s) + 1) //整个报文长度
	head.Ret = ret
	buf := make([]byte, head.Len)
	buf[0] = common.HTV3MagicBegin
	err = common.SerialHeadV3ToSlice(head, buf[1:])
	if err != nil {
		infoLog.Printf("SendQueryUserIsAlreadyInRoomRetCode SerialHeadV3ToSlice failed")
		return false
	}
	copy(buf[common.PacketV3HeadLen:], s) // return code
	buf[head.Len-1] = common.HTV3MagicEnd

	rspPacket := common.NewHeadV3Packet(buf)
	c.AsyncWritePacket(rspPacket, time.Second)
	return true
}
func SendQueryUserIsAlreadyInRoomResp(c *gotcp.Conn, reqHead *common.HeadV3, subMucRspBody *ht_muc.QueryUserIsAlreadyInRoomRspBody) bool {
	head := new(common.HeadV3)
	if reqHead != nil {
		*head = *reqHead
	}
	rspBody := new(ht_muc.MucRspBody)
	rspBody.QueryUserIsAlreadyInRspbody = subMucRspBody
	s, err := proto.Marshal(rspBody)
	if err != nil {
		infoLog.Printf("SendQueryUserIsAlreadyInRoomResp proto.Marshal failed from=%v to=%v cmd=%v seq=%v",
			head.From,
			head.To,
			head.Cmd,
			head.Seq)
		return false
	}

	head.Len = uint32(common.PacketV3HeadLen + len(s) + 1) //整个报文长度
	head.Ret = uint16(subMucRspBody.GetStatus().GetCode())
	buf := make([]byte, head.Len)
	buf[0] = common.HTV3MagicBegin
	err = common.SerialHeadV3ToSlice(head, buf[1:])
	if err != nil {
		infoLog.Printf("SendQueryUserIsAlreadyInRoomResp SerialHeadV3ToSlice failed")
		return false
	}
	copy(buf[common.PacketV3HeadLen:], s) // return code
	buf[head.Len-1] = common.HTV3MagicEnd

	rspPacket := common.NewHeadV3Packet(buf)
	c.AsyncWritePacket(rspPacket, time.Second)
	return true
}

// 25.Proc QueryQRcode info
func ProcQueryQRcodeInfo(c *gotcp.Conn, head *common.HeadV3, reqBody *ht_muc.MucReqBody) bool {
	result := uint16(ht_muc.MUC_RET_CODE_RET_SUCCESS)
	errMsg := []byte("operation success")
	bNeedCall := true
	defer func() {
		if bNeedCall {
			SendQueryQRcodeInfoRetCode(c, head, result, errMsg)
		}
	}()

	// 检查输入参数是否为空
	if head == nil || reqBody == nil {
		bNeedCall = false
		return false
	}

	// add static
	attr := "gomuc/query_orcode_req_count"
	libcomm.AttrAdd(attr, 1)

	subReqBody := reqBody.GetQueryQrcodeInfoReqbody()
	roomId := subReqBody.GetRoomId()
	scanUid := subReqBody.GetScanUid()
	shareUid := subReqBody.GetShareId()
	maxShowCnt := subReqBody.GetMaxShowCnt()
	infoLog.Printf("ProcQueryQRcodeInfo recv from=%v to=%v cmd=%v seq=%v roomId=%v scanUid=%v shareUid=%v maxshowCnt=%v",
		head.From,
		head.To,
		head.Cmd,
		head.Seq,
		roomId,
		scanUid,
		shareUid,
		maxShowCnt)
	// first check input param
	if roomId == 0 || scanUid == 0 || shareUid == 0 {
		infoLog.Printf("ProcQueryQRcodeInfo invalid param")
		result = uint16(ht_muc.MUC_RET_CODE_RET_INPUT_PARAM_ERR)
		errMsg = []byte("input param error")
		return false
	}

	// second get room Info
	roomInfo, err := roomManager.GetRoomInfo(roomId)
	if err != nil {
		infoLog.Printf("ProcQueryQRcodeInfo roomId=%v scanUid=%v get room info failed", roomId, scanUid)
		result = uint16(ht_muc.MUC_RET_CODE_RET_INTERNAL_ERR)
		errMsg = []byte("Internal error")
		return false
	}

	memberList := roomInfo.MemberList
	if maxShowCnt > uint32(len(memberList)) {
		maxShowCnt = uint32(len(memberList))
	}
	var bScanInRoom bool
	var bShareInRoom bool
	var subMemberList []uint32

	for i, v := range memberList {
		if v.Uid == scanUid {
			bScanInRoom = true
		}
		if v.Uid == shareUid {
			bShareInRoom = true
		}

		if uint32(i) < maxShowCnt {
			subMemberList = append(subMemberList, v.Uid)
		}
	}

	subMucRspBody := new(ht_muc.QueryQRcodeInfoRspBody)
	subMucRspBody.Status = &ht_muc.MucHeader{Code: proto.Uint32(uint32(ht_muc.MUC_RET_CODE_RET_SUCCESS))}
	subMucRspBody.IsScanerInRoom = proto.Bool(bScanInRoom)
	subMucRspBody.IsSharerInRoom = proto.Bool(bShareInRoom)
	subMucRspBody.RoomName = []byte(roomInfo.RoomName)
	subMucRspBody.MemberCount = proto.Uint32(uint32(len(memberList)))
	subMucRspBody.Userids = subMemberList

	infoLog.Printf("DEBUG ProcQueryQRcodeInfo subMucRspBody=%#v", subMucRspBody)

	SendQueryQRcodeInfoResp(c, head, subMucRspBody)
	bNeedCall = false
	return true
}

func SendQueryQRcodeInfoRetCode(c *gotcp.Conn, reqHead *common.HeadV3, ret uint16, errMsg []byte) bool {
	head := new(common.HeadV3)
	if reqHead != nil {
		*head = *reqHead
	}
	rspBody := new(ht_muc.MucRspBody)
	rspBody.QueryQrcodeInfoRspbody = new(ht_muc.QueryQRcodeInfoRspBody)
	subMucRspBody := rspBody.GetQueryQrcodeInfoRspbody()
	subMucRspBody.Status = &ht_muc.MucHeader{Code: proto.Uint32(uint32(ret)), Reason: errMsg}
	s, err := proto.Marshal(rspBody)
	if err != nil {
		infoLog.Printf("SendQueryQRcodeInfoRetCode proto.Marshal failed from=%v to=%v cmd=%v seq=%v",
			head.From,
			head.To,
			head.Cmd,
			head.Seq)
		return false
	}
	head.Len = uint32(common.PacketV3HeadLen + len(s) + 1) //整个报文长度
	head.Ret = ret
	buf := make([]byte, head.Len)
	buf[0] = common.HTV3MagicBegin
	err = common.SerialHeadV3ToSlice(head, buf[1:])
	if err != nil {
		infoLog.Printf("SendQueryQRcodeInfoRetCode SerialHeadV3ToSlice failed")
		return false
	}
	copy(buf[common.PacketV3HeadLen:], s) // return code
	buf[head.Len-1] = common.HTV3MagicEnd

	rspPacket := common.NewHeadV3Packet(buf)
	c.AsyncWritePacket(rspPacket, time.Second)
	return true
}
func SendQueryQRcodeInfoResp(c *gotcp.Conn, reqHead *common.HeadV3, subMucRspBody *ht_muc.QueryQRcodeInfoRspBody) bool {
	head := new(common.HeadV3)
	if reqHead != nil {
		*head = *reqHead
	}
	rspBody := new(ht_muc.MucRspBody)
	rspBody.QueryQrcodeInfoRspbody = subMucRspBody
	s, err := proto.Marshal(rspBody)
	if err != nil {
		infoLog.Printf("SendQueryQRcodeInfoResp proto.Marshal failed from=%v to=%v cmd=%v seq=%v",
			head.From,
			head.To,
			head.Cmd,
			head.Seq)
		return false
	}

	head.Len = uint32(common.PacketV3HeadLen + len(s) + 1) //整个报文长度
	head.Ret = uint16(subMucRspBody.GetStatus().GetCode())
	buf := make([]byte, head.Len)
	buf[0] = common.HTV3MagicBegin
	err = common.SerialHeadV3ToSlice(head, buf[1:])
	if err != nil {
		infoLog.Printf("SendQueryQRcodeInfoResp SerialHeadV3ToSlice failed")
		return false
	}
	copy(buf[common.PacketV3HeadLen:], s) // return code
	buf[head.Len-1] = common.HTV3MagicEnd

	rspPacket := common.NewHeadV3Packet(buf)
	c.AsyncWritePacket(rspPacket, time.Second)
	return true
}

func (this *Callback) OnClose(c *gotcp.Conn) {
	infoLog.Println("OnClose:", c.GetExtraData())
}

type Options struct {
	// Example of verbosity with level
	Verbose []bool `short:"v" long:"verbose" description:"Verbose output"`

	// Example of optional value
	ServerConf string `short:"c" long:"conf" description:"Server Config" optional:"no"`
}

var options Options

var parser = flags.NewParser(&options, flags.Default)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// 处理命令行参数
	if _, err := parser.Parse(); err != nil {
		log.Fatalln("parse cmd line failed!")
	}

	if options.ServerConf == "" {
		log.Fatalln("Must input config file name")
	}

	// log.Printf("config name =", options.ServerConf)
	// 读取配置文件
	cfg, err := ini.Load([]byte(""), options.ServerConf)
	if err != nil {
		log.Printf("load config file=%s failed", options.ServerConf)
		return
	}
	// 配置文件只读 设置此标识提升性能
	cfg.BlockMode = false
	// 定义一个文件
	fileName := cfg.Section("LOG").Key("path").MustString("/home/ht/server.log")
	logFile, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	defer logFile.Close()
	if err != nil {
		log.Fatalln("open file error !")
		return
	}

	// 创建一个日志对象
	infoLog = log.New(logFile, "[Info]", log.LstdFlags)
	// 配置log的Flag参数
	infoLog.SetFlags(infoLog.Flags() | log.LstdFlags)

	// 读取 Memcache Ip and port
	mcIp := cfg.Section("MEMCACHE").Key("mc_ip").MustString("127.0.0.1")
	mcPort := cfg.Section("MEMCACHE").Key("mc_port").MustInt(11211)
	infoLog.Printf("memcache ip=%v port=%v", mcIp, mcPort)
	mcApi = new(common.MemcacheApi)
	mcApi.Init(mcIp + ":" + strconv.Itoa(mcPort))

	// init mongo
	// 创建mongodb对象
	mongo_url := cfg.Section("MONGO").Key("url").MustString("localhost")
	infoLog.Printf("Mongo url=%s", mongo_url)
	mongoSess, err = mgo.Dial(mongo_url)
	if err != nil {
		log.Fatalln("connect mongodb failed")
		return
	}
	defer mongoSess.Close()
	// Optional. Switch the session to a monotonic behavior.
	mongoSess.SetMode(mgo.Monotonic, true)

	// init mysql
	mysqlHost := cfg.Section("MYSQL").Key("mysql_host").MustString("127.0.0.1")
	mysqlUser := cfg.Section("MYSQL").Key("mysql_user").MustString("IMServer")
	mysqlPasswd := cfg.Section("MYSQL").Key("mysql_passwd").MustString("hello")
	mysqlDbName := cfg.Section("MYSQL").Key("mysql_db").MustString("HT_IMDB")
	mysqlPort := cfg.Section("MYSQL").Key("mysql_port").MustString("3306")

	infoLog.Printf("mysql host=%v user=%v passwd=%v dbname=%v port=%v",
		mysqlHost,
		mysqlUser,
		mysqlPasswd,
		mysqlDbName,
		mysqlPort)

	db, err = sql.Open("mysql", mysqlUser+":"+mysqlPasswd+"@"+"tcp("+mysqlHost+":"+mysqlPort+")/"+mysqlDbName+"?charset=utf8&timeout=90s")
	if err != nil {
		infoLog.Printf("open mysql failed")
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}

	// 读取offline 配置
	offlineIp := cfg.Section("OFFLINE").Key("offline_ip").MustString("127.0.0.1")
	offlinePort := cfg.Section("OFFLINE").Key("offline_port").MustString("0")
	offlineConnLimit := cfg.Section("OFFLINE").Key("pool_limit").MustInt(1000)
	infoLog.Printf("offline server ip=%v port=%v connLimit=%v", offlineIp, offlinePort, offlineConnLimit)
	offlineApi = common.NewOfflineApiV2(offlineIp, offlinePort, 3*time.Second, 3*time.Second, &common.HeadV2Protocol{}, offlineConnLimit)

	// 读取IMServer 配置
	imCount := cfg.Section("IMSERVER").Key("imserver_cnt").MustInt(2)
	imConnLimit := cfg.Section("IMSERVER").Key("pool_limit").MustInt(1000)
	imServer := make(map[string]*common.ImServerApiV2, imCount)
	imOldServer := make(map[string]*common.ImServerApiV2, imCount)
	infoLog.Printf("IMServer Count=%v", imCount)
	for i := 0; i < imCount; i++ {
		ipKey := "imserver_ip_" + strconv.Itoa(i)
		ipOnlineKye := "imserver_ip_online_" + strconv.Itoa(i)
		portKey := "imserver_port_" + strconv.Itoa(i)
		oldPortKey := "imserver_old_port_" + strconv.Itoa(i)
		imIp := cfg.Section("IMSERVER").Key(ipKey).MustString("127.0.0.1")
		imIpOnline := cfg.Section("IMSERVER").Key(ipOnlineKye).MustString("127.0.0.1")
		imPort := cfg.Section("IMSERVER").Key(portKey).MustString("18380")
		imOldPort := cfg.Section("IMSERVER").Key(oldPortKey).MustString("18280")

		infoLog.Printf("im server ip=%v ip_online=%v port=%v oldPort=%v", imIp, imIpOnline, imPort, imOldPort)
		imServer[imIpOnline] = common.NewImServerApiV2(imIp, imPort, 3*time.Second, 3*time.Second, &common.HeadV3Protocol{}, imConnLimit)
		imOldServer[imIpOnline] = common.NewImServerApiV2(imIp, imOldPort, 3*time.Second, 3*time.Second, &common.XTHeadProtocol{}, imConnLimit)
	}
	//创建RoomManager 和 DbUtil 对象
	DbUtil = util.NewDbUtil(db, mongoSess, infoLog)
	roomManager = util.NewRoomManager(DbUtil, mcApi, offlineApi, imServer, imOldServer, infoLog)
	// creates a tcp listener
	serverIp := cfg.Section("LOCAL_SERVER").Key("bind_ip").MustString("127.0.0.1")
	serverPort := cfg.Section("LOCAL_SERVER").Key("bind_port").MustInt(8990)
	infoLog.Printf("serverIp=%v serverPort=%v", serverIp, serverPort)
	tcpAddr, err := net.ResolveTCPAddr("tcp4", serverIp+":"+strconv.Itoa(serverPort))
	checkError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	// creates a server
	sendChanLimit := cfg.Section("CHANLIMIT").Key("max_send_chan_count").MustUint(1000)
	recvChanLimit := cfg.Section("CHANLIMIT").Key("max_recv_chan_count").MustUint(1000)
	config := &gotcp.Config{
		PacketSendChanLimit:    uint32(sendChanLimit),
		PacketReceiveChanLimit: uint32(recvChanLimit),
	}

	// 创建请求滤器
	reqRecord = map[string]int64{}

	srv := gotcp.NewServer(config, &Callback{}, &common.HeadV3Protocol{})

	// starts service
	go srv.Start(listener, time.Second)
	infoLog.Println("listening:", listener.Addr())

	// 启动tikcer
	ticker = time.NewTicker(time.Second * 60)
	go func() {
		index := 0
		for range ticker.C {
			reqLock.Lock()
			curSecond := time.Now().Unix()
			for k, v := range reqRecord {
				infoLog.Printf("main key=%v value=%v", k, v)
				if v+REQTHRESHOLD < curSecond {
					infoLog.Printf("main delete key=%v value=%v", k, v)
					delete(reqRecord, k)
				}
			}
			infoLog.Printf("main remain req size=%v", len(reqRecord))
			reqLock.Unlock()

			infoLog.Printf("main ticker tick index=%v", index)
			index++
		}
	}()

	// catchs system signal
	chSig := make(chan os.Signal)
	signal.Notify(chSig, syscall.SIGINT, syscall.SIGTERM)
	infoLog.Println("Signal: ", <-chSig)

	// stops service
	srv.Stop()
	ticker.Stop()
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
