// Copyright 2016 songliwei
//
// HelloTalk.inc

package util

import (
	"errors"
	"log"
	"net"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gansidui/gotcp/libcrypto"
	"github.com/gansidui/gotcp/tcpfw/common"
	"github.com/gansidui/gotcp/tcpfw/include/ht_muc"
	"github.com/golang/protobuf/proto"
)

var (
	ErrNilMCObject = errors.New("not set memcache object current is nil")
	ErrInputParam  = errors.New("err input param error")
	ErrOnlineIpErr = errors.New("err online state ip not exist")
	ErrExecLimit   = errors.New("err member exec limit")
	ErrAlreadyIn   = errors.New("err member already in")
	ErrNotInRoom   = errors.New("err member not int room")
)

const (
	CSaveOffline        = 0
	CSaveOffLineAndPush = 1
	CSendToIMServer     = 2
)

const (
	CNotBeenAt = 0
	CBeenAt    = 1
)

const (
	MUC_MEMBER_LIMIT     = 200
	MUC_MEMBER_LIMIT_VIP = 500
)

const (
	AT_NONE_MEMBER = 0 //没有@用户
	AT_ALL_MEMBER  = 1 //@所有用户
	AT_USER_LIST   = 2 //@remindList中的用户
)

const (
	CT_P2P = 0
	CT_MUC = 1
)

const (
	NICKNAME_LEN    = 128
	CONTENT_LEN     = 128
	NORMALDEC_LEN   = 128
	SERVER_COMM_KEY = "lp$5F@nfN0Oh8I*5"
)

const (
	CMD_MUC_MESSAGE     = 0x7009
	CMD_MUC_MESSAGE_ACK = 0x700A

	CMD_S2S_MESSAGE_PUSH     = 0x8027
	CMD_S2S_MESSAGE_PUSH_ACK = 0x8028

	CMD_NOTIFY_REMOVE_MEMBER  = 0x700D
	CMD_RECEIPT_REMOVE_MEMBER = 0x700E

	CMD_GVOIP_INVITE_BROADCAST               = 0x7103
	CMD_GVOIP_INVITE_BROADCAST_RECEIPT       = 0x7104
	CMD_GVOIP_MEMBER_JOIN_BROADCAST          = 0x7109
	CMD_GVOIP_MEMBER_JOIN_BROADCAST_RECEIPT  = 0x710A
	CMD_GVOIP_MEMBER_LEAVE_BROADCAST         = 0x710D
	CMD_GVOIP_MEMBER_LEAVE_BROADCAST_RECEIPT = 0x710E
	CMD_GVOIP_END_BROADCAST                  = 0x7111
	CMD_GVOIP_END_BROADCAST_RECEIPT          = 0x7112
)

// 群聊消息推送相关信息
type MucPushInfo struct {
	PushType  uint8  // 推送设置
	NickName  string // 推送者呢称
	MsgId     string // msgId
	PushParam string // 推送内容
}

// 每个群成员详细信息
type MemberInfoStruct struct {
	Uid         uint32 // 用户uid
	InvitedUid  uint32 // 邀请用户加入的uid
	NickName    string // 用户昵称
	OrderId     uint32 // 加入群的序号
	JoinTs      uint32 // 加入时间戳
	PushSetting uint32 // 推送设置0:有推送 1:关闭推送
	RoomId      uint32 // 所属的群聊ID
}

type AnnouncementStruct struct {
	PublishUid  uint32 // 群公告发布者
	PublishTS   uint32 // 群公告发布时间
	AnnoContect string // 群公告
}

type RoomInfo struct {
	RoomId       uint32              // 群聊ID
	CreateUid    uint32              // 创建者uid
	AdminList    []uint32            // 管理员列表 为slice 避免后面扩展带来问题
	AdminLimit   uint32              // 管理员个数限制
	RoomName     string              // 群名称
	RoomDesc     string              // 群描述
	MemberLimit  uint32              // 群成员数限制
	MaxOrder     int64               // 当前最大的加入序号
	VerifyStat   uint32              // 加群是否需要管理员确认
	Announcement AnnouncementStruct  // 群公告
	MemberList   []*MemberInfoStruct // 群聊成员详细列表
	RoomTS       int64               // 群资料时间戳
}

type RoomManager struct {
	roomIdToRoomInfo map[uint32]*RoomInfo
	roomIdToUser     map[uint32][]uint32 // RoomId ==> block voip user list
	dbManager        *DbUtil
	sendSeq          uint16
	mcApi            *common.MemcacheApi
	offlineApi       *common.OfflineApiV2
	infoLog          *log.Logger
	roomInfoLock     sync.Mutex // sync mutex goroutines use this.roomIdToRoomInfo
	voipLock         sync.Mutex // sync mutex goroutines use roomIdToUser map
	imServer         map[string]*common.ImServerApiV2
	imOldServer      map[string]*common.ImServerApiV2
}

func NewRoomManager(dbUtil *DbUtil,
	mc *common.MemcacheApi,
	offline *common.OfflineApiV2,
	imServerApi map[string]*common.ImServerApiV2,
	imOldServerApi map[string]*common.ImServerApiV2,
	logger *log.Logger) *RoomManager {
	return &RoomManager{
		roomIdToRoomInfo: map[uint32]*RoomInfo{},
		roomIdToUser:     map[uint32][]uint32{},
		dbManager:        dbUtil,
		sendSeq:          0,
		mcApi:            mc,
		offlineApi:       offline,
		imServer:         imServerApi,
		imOldServer:      imOldServerApi,
		infoLog:          logger,
	}
}

// Convert uint to net.IP http://www.outofmemory.cn
func inet_ntoa(ipnr int64) net.IP {
	var bytes [4]byte
	bytes[0] = byte((ipnr >> 24) & 0xFF)
	bytes[1] = byte((ipnr >> 16) & 0xFF)
	bytes[2] = byte((ipnr >> 8) & 0xFF)
	bytes[3] = byte(ipnr & 0xFF)

	return net.IPv4(bytes[3], bytes[2], bytes[1], bytes[0])
}

// Convert net.IP to int64 ,  http://www.outofmemory.cn
func inet_aton(ipnr net.IP) int64 {
	bits := strings.Split(ipnr.String(), ".")

	b0, _ := strconv.Atoi(bits[0])
	b1, _ := strconv.Atoi(bits[1])
	b2, _ := strconv.Atoi(bits[2])
	b3, _ := strconv.Atoi(bits[3])

	var sum int64

	sum += int64(b0) << 24
	sum += int64(b1) << 16
	sum += int64(b2) << 8
	sum += int64(b3)

	return sum
}

func (this *RoomManager) GetPacketSeq() (packetSeq uint16) {
	this.sendSeq++
	return this.sendSeq
}

func (this *RoomManager) CreateMuc(createUid uint32, memberList []uint32, memberLimit uint32) (roomId uint32, roomTS uint64, err error) {
	roomId, roomTS, err = this.dbManager.CreateMucRoom(createUid, memberList, memberLimit)
	if err != nil {
		this.infoLog.Printf("CreateMuc db err =%v", err)
		return roomId, roomTS, err
	}

	// roomMember, maxOrderId, err := this.dbManager.GetRoomMemberList(roomId)
	// if err != nil {
	// 	this.infoLog.Printf("CreateMuc GetRoomMemberList failed err=%s", err)
	// 	return roomId, roomTS, err
	// }

	// adminList := [10]uint32{createUid}
	// // 在写this.roomIdToRoomInfo 的时候一定要锁住
	// this.roomInfoLock.Lock()
	// defer this.roomInfoLock.Unlock()
	// this.roomIdToRoomInfo[roomId] = &RoomInfo{
	// 	RoomId:      roomId,
	// 	CreateUid:   createUid,
	// 	AdminList:   adminList,
	// 	MemberLimit: memberLimit,
	// 	MaxOrder:    maxOrderId,
	// 	MemberList:  roomMember,
	// 	RoomTS:      int64(roomTS),
	// }
	//
	return roomId, roomTS, nil
}

func (this *RoomManager) InviteMember(roomId, inviteId uint32, memberList []*ht_muc.RoomMemberInfo) (roomTS int64, alreadyIn []uint32, err error) {
	// roomInfo 为一个群信息的指针 可以通过此指针修改群信息
	roomInfo, err := this.GetRoom(roomId)
	if err != nil {
		this.infoLog.Printf("InviteMember GetRoom roomId=%v failed err==%v", roomId, err)
		return 0, nil, err
	}

	// 判断成员是否超过限制
	if uint32(len(memberList)+len(roomInfo.MemberList)) > roomInfo.MemberLimit {
		// 如果当前群成员数限制不等于vip 群成员数限制需要检查用户是否是vip
		if roomInfo.MemberLimit < MUC_MEMBER_LIMIT_VIP {
			this.infoLog.Printf("InviteMember roomId=%v memberLimit=%v not equal viplimiit=%v",
				roomId,
				roomInfo.MemberLimit,
				MUC_MEMBER_LIMIT_VIP)

			vipExpireTs, err := this.dbManager.GetUserVIPExpireTS(roomInfo.CreateUid)
			if err != nil {
				this.infoLog.Printf("InviteMember GetUserVIPExpireTS failed room=%v inviteId=%v err=%v",
					roomId,
					inviteId,
					err)
				// 查询vip过期时间失败直接认为过期 用户成员数超过群限制返回错误
				err = ErrExecLimit
				this.infoLog.Printf("InviteMember roomId=%v inviteId=%v exec limit err=%v", roomId, inviteId, err)
				return roomTS, alreadyIn, err

			}
			tsNow := time.Now().Unix()
			if vipExpireTs > uint64(tsNow) {
				// 是vip会员 但是群信息还是非vip时创建的则更新群成员数限制到数据库同时更新内存
				roomInfo.MemberLimit = MUC_MEMBER_LIMIT_VIP
				err = this.dbManager.UpdateRoomMemberLimit(roomId, roomInfo.MemberLimit)
				if err != nil {
					this.infoLog.Printf("InviteMember UpdateRoomMemberLimit failed room=%v memberLimit=%v err=%v",
						roomId,
						roomInfo.MemberLimit,
						err)
				}
				// 更新群成员数限制后仍然超员了直接返回错误
				if uint32(len(memberList)+len(roomInfo.MemberList)) > roomInfo.MemberLimit {
					err = ErrExecLimit
					this.infoLog.Printf("InviteMember roomId=%v inviteId=%v exec limit err=%v", roomId, inviteId, err)
					return roomTS, alreadyIn, err
				}
			} else {
				// 查询成功但是会员已过期直接返回超员了
				err = ErrExecLimit
				this.infoLog.Printf("InviteMember roomId=%v inviteId=%v exec limit err=%v", roomId, inviteId, err)
				return roomTS, alreadyIn, err
			}
		} else {
			//群成员超过限制做VIP确认
			err = ErrExecLimit
			this.infoLog.Printf("InviteMember roomId=%v inviteId=%v exec limit err=%v", roomId, inviteId, err)
			return roomTS, alreadyIn, err
		}
	}

	// 检测已经在群聊中的待添加用户列表
	alreadyIn, err = this.dbManager.GetMemberAlreadyInMuc(roomId, memberList)
	var filterMemberList []*ht_muc.RoomMemberInfo
	if len(alreadyIn) != 0 {
		for _, v := range memberList {
			if this.UidIsInSlice(alreadyIn, v.GetUid()) {
				// 用户已经在群中了直接过滤掉
				this.infoLog.Printf("InviteMember roomId=%v uid=%v is alreadyIn", roomId, v.GetUid())
				continue
			}
			// 否则添加到过滤列表总
			filterMemberList = append(filterMemberList, v)
		}
	}

	// 过滤之后如果仍然存在用户则添加到群中
	if len(filterMemberList) > 0 {
		roomTS, err = this.dbManager.InviteMember(roomId, inviteId, filterMemberList)
		if err != nil {
			this.infoLog.Printf("InviteMember roomId=%v inviteId=%v exec db InviteMember failed", roomId, inviteId)
			return roomTS, alreadyIn, err
		}
		// 更新内存总的群成员列表
		totalMemberList, maxOrder, err := this.dbManager.GetRoomMemberList(roomId)
		if err != nil {
			this.infoLog.Printf("InviteMember exec db.GetRoomMemberList()failed roomId=%v err=%v", roomId, err)
			return roomTS, alreadyIn, err
		}
		// update memberlist
		this.roomInfoLock.Lock()
		defer this.roomInfoLock.Unlock()
		roomInfo.MaxOrder = maxOrder
		roomInfo.MemberList = totalMemberList
		roomInfo.RoomTS = roomTS
	}
	return roomTS, alreadyIn, err
}

func (this *RoomManager) UidIsInSlice(uidList []uint32, uid uint32) bool {
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

func (this *RoomManager) NotifyInviteMember(opInfo *ht_muc.RoomMemberInfo,
	inviterInfo *ht_muc.RoomMemberInfo,
	memberList []*ht_muc.RoomMemberInfo,
	roomId uint32,
	roomTS uint64,
	roomIdFrom uint32) (err error) {
	head := &common.HeadV3{Flag: uint8(common.CServToServ),
		Version:  common.CVerMmedia,
		CryKey:   uint8(common.CNoneKey),
		TermType: uint8(0),
		Cmd:      uint16(ht_muc.MUC_CMD_TYPE_GO_CMD_NOTIFY_INVITE_MEMBER),
		Seq:      this.GetPacketSeq(),
		From:     roomId,
		To:       0,
		Len:      0,
	}

	reqBody := new(ht_muc.MucReqBody)
	reqBody.NotifyInviteMemberReqbody = &ht_muc.NotifyInviteMemberReqBody{
		RoomId:        proto.Uint32(roomId),
		OpInfo:        opInfo,
		InviterInfo:   inviterInfo,
		Members:       memberList,
		NotifyTime:    proto.Uint32(uint32(time.Now().Unix())),
		RoomTimestamp: proto.Uint64(roomTS),
		RoomIdFrom:    proto.Uint32(roomIdFrom),
	}
	notifyBody, err := proto.Marshal(reqBody)
	if err != nil {
		this.infoLog.Printf("NotifyInviteMember proto marshal roomId=%v inviterId=%v err=%v",
			roomId,
			inviterInfo.GetUid(),
			err)
		return err
	}

	//调用广播接口 如果是扫描二维码加群的广播给所有人
	if roomIdFrom == uint32(ht_muc.ROOMID_FROM_TYPE_ENUM_FROM_QRCODE) {
		return this.BroadCastNotification(roomId, 0, head, notifyBody)
	} else {
		//如果是邀请加群的不广播给邀请者
		return this.BroadCastNotification(roomId, inviterInfo.GetUid(), head, notifyBody)
	}
}

func (this *RoomManager) BroadCastNotification(roomId, exceptUid uint32, head *common.HeadV3, notifyBody []byte) (err error) {
	roomInfo, err := this.GetRoom(roomId)
	if err != nil {
		this.infoLog.Printf("BroadCastNotification GetRoom failed roomId=%v", roomId)
		return err
	}
	if this.mcApi == nil {
		this.infoLog.Printf("BroadCastNotification memcache api object not set roomId=%v", roomId)
		return ErrNilMCObject
	}

	memberList := roomInfo.MemberList
	for _, v := range memberList {
		// 屏蔽掉过滤用户
		if exceptUid == v.Uid {
			this.infoLog.Printf("BroadCastNotification roomId=%v exceptUid=%v", roomId, exceptUid)
			continue
		}
		// 用户消息处理方式 默认为存储离线
		var procType int = CSaveOffline
		// 查询用户的在线状态
		onlineStat, err := this.mcApi.GetUserOnlineStat(v.Uid)
		if err == nil {
			procType = this.GetMucMsgProcType(onlineStat, v, false)
		} else {
			this.infoLog.Printf("BroadCastNotification Get msg proc failed roomId=%v uid=%v err=%v", roomId, v.Uid, err)
		}
		// 调整发送头部的to字段
		head.To = v.Uid //消息的接收者
		this.infoLog.Printf("BroadCastNotification Msg proc type=%v", procType)

		switch procType {
		case CSendToIMServer:
			err := this.SendPacketToIMServerRelabile(onlineStat, head, notifyBody)
			if err != nil {
				this.infoLog.Printf("BroadCastNotification SendPacketToIMServerRelabile failed roomId=%v from=%v to=%v seq=%v",
					roomId,
					head.From,
					head.To,
					head.Seq)
				// 发送到IM失败存储离线
				ret, err := this.offlineApi.SendPacketWithHeadV3(head, notifyBody, nil)
				if err != nil {
					this.infoLog.Printf("BroadCastNotification save offline faield ret=%v err=%v", ret, err)
				} else {
					this.infoLog.Printf("BroadCastNotification save offline success roomId=%v from=%v to=%v cmd=%v seq=%v ret=%v",
						roomId,
						head.From,
						head.To,
						head.Cmd,
						head.Seq,
						ret)
				}
			} else {
				this.infoLog.Printf("BroadCastNotification SendPacketToIMServerRelabile success from=%v to=%v cmd=%v seq=%v",
					head.From,
					head.To,
					head.Cmd,
					head.Seq)
			}
		case CSaveOffLineAndPush, CSaveOffline:
			ret, err := this.offlineApi.SendPacketWithHeadV3(head, notifyBody, nil)
			if err != nil {
				this.infoLog.Printf("BroadCastNotification save offline faield ret=%v err=%v", ret, err)
			} else {
				this.infoLog.Printf("BroadCastNotification save offline success roomId=%v from=%v to=%v cmd=%v seq=%v ret=%v",
					roomId,
					head.From,
					head.To,
					head.Cmd,
					head.Seq,
					ret)
			}

		default:
			this.infoLog.Printf("BroadCastNotification Unhandle stat=%v", procType)
		}
	}
	return err
}

func (this *RoomManager) GetRoom(roomId uint32) (roomInfo *RoomInfo, err error) {
	roomInfo, ok := this.roomIdToRoomInfo[roomId]
	if ok {
		this.infoLog.Printf("GetRoom roomId=%v exist in memory", roomId)
		return roomInfo, nil
	} else {
		// 操作this.roomIdToRoomInfo 时一定要锁住
		this.roomInfoLock.Lock()
		defer this.roomInfoLock.Unlock()
		this.infoLog.Printf("GetRoom roomId=%v not found load from db", roomId)
		roomInfo, err := this.dbManager.GetRoomBaseInfo(roomId)
		if err != nil {
			this.infoLog.Printf("GetRoom GetRoomBaseInfo failed roomId=%v", roomId)
			return nil, err
		}
		roomMember, maxOrderId, err := this.dbManager.GetRoomMemberList(roomId)
		if err != nil {
			this.infoLog.Printf("GetRoom GetRoomMemberList failed roomId=%v err=%s", roomId, err)
			return nil, err
		}
		roomInfo.MaxOrder = maxOrderId
		roomInfo.MemberList = roomMember

		this.roomIdToRoomInfo[roomId] = roomInfo
		return roomInfo, nil
	}
}

func (this *RoomManager) GetMucMsgProcType(stat *common.UserState, memberInfo *MemberInfoStruct, bAtUser bool) (procType int) {
	procType = CSaveOffline
	if stat.ClientType == common.CClientTyepIOS {
		procType = this.GetIOSMucMsgProcType(stat, memberInfo, bAtUser)
	} else {
		procType = this.GetAndroidMucMsgProcType(stat, memberInfo, bAtUser)
	}
	return procType
}

func (this *RoomManager) GetIOSMucMsgProcType(stat *common.UserState, memberInfo *MemberInfoStruct, atUser bool) (procType int) {
	procType = CSaveOffline
	if stat.OnlineStat == common.ST_ONLINE {
		procType = CSendToIMServer
	} else if stat.OnlineStat != common.ST_LOGOUT && (memberInfo.PushSetting == 0 || atUser) {
		procType = CSaveOffLineAndPush
	} else {
		procType = CSaveOffline
	}
	return procType
}

func (this *RoomManager) GetAndroidMucMsgProcType(stat *common.UserState, memberInfo *MemberInfoStruct, atUser bool) (procType int) {
	procType = CSaveOffline
	if stat.OnlineStat == common.ST_ONLINE {
		procType = CSendToIMServer
	} else if stat.OnlineStat != common.ST_LOGOUT && (memberInfo.PushSetting == 0 || atUser) {
		procType = CSaveOffLineAndPush
	} else {
		procType = CSaveOffline
	}
	return procType
}

func (this *RoomManager) SendPacketToIMServerRelabile(stat *common.UserState, head *common.HeadV3, payLoad []byte) (err error) {
	if stat.OnlineStat != common.ST_ONLINE {
		this.infoLog.Printf("SendPacketToIMServerRelabile online stat=%v error", stat.OnlineStat)
		return ErrInputParam
	}

	svrIp := inet_ntoa(int64(stat.SvrIp)).String()
	this.infoLog.Printf("SendPacketToIMServerRelabile svrIP=%s oriIP=%v", svrIp, stat.SvrIp)
	v, ok := this.imServer[svrIp]
	if !ok { // 不存在直接打印日志返回错误
		this.infoLog.Printf("SendPacketToIMServerRelabile not exist IM ip=%v imServer=%v", svrIp, this.imServer)
		return ErrOnlineIpErr
	}

	tryCount := 2
	for tryCount > 0 {
		ret, err := v.SendPacket(head, payLoad)
		if err == nil && ret == 0 {
			this.infoLog.Printf("SendPacketToIMServerRelabile send succ from=%v to=%v seq=%v cmd=%v", head.From, head.To, head.Seq, head.Cmd)
			return nil
		} else {
			this.infoLog.Printf("SendPacketToIMServerRelabile failed from=%v to=%v seq=%v cmd=%v err=%v ret=%v", head.From, head.To, head.Seq, head.Cmd, err, ret)
		}
		tryCount--
	}
	return common.ErrReachMaxTryCount
}

func (this *RoomManager) SendPacketToOldIMServerRelabile(stat *common.UserState, head *common.XTHead, payLoad []byte) (err error) {
	if stat.OnlineStat != common.ST_ONLINE {
		this.infoLog.Printf("SendPacketToOldIMServerRelabile online stat=%v error", stat.OnlineStat)
		return ErrInputParam
	}

	svrIp := inet_ntoa(int64(stat.SvrIp)).String()
	this.infoLog.Printf("SendPacketToOldIMServerRelabile svrIP=%s oriIP=%v", svrIp, stat.SvrIp)
	v, ok := this.imOldServer[svrIp]
	if !ok { // 不存在直接打印日志返回错误
		this.infoLog.Printf("SendPacketToOldIMServerRelabile not exist old IM ip=%v imServer=%v", svrIp, this.imServer)
		return ErrOnlineIpErr
	}

	tryCount := 2
	for tryCount > 0 {
		err = v.SendXTPacket(head, payLoad)
		if err == nil {
			this.infoLog.Printf("SendPacketToOldIMServerRelabile send succ from=%v to=%v seq=%v cmd=%v", head.From, head.To, head.Seq, head.Cmd)
			return nil
		} else {
			this.infoLog.Printf("SendPacketToOldIMServerRelabile failed from=%v to=%v seq=%v cmd=%v err=%v", head.From, head.To, head.Seq, head.Cmd, err)
		}
		tryCount--
	}
	return common.ErrReachMaxTryCount
}

func (this *RoomManager) IsCreateUid(roomId, uid uint32) (ret bool, err error) {
	if roomId == 0 || uid == 0 {
		this.infoLog.Printf("IsCreateUid input err roomId=%v uid=%v", roomId, uid)
		return false, ErrInputParam
	}
	// roomInfo 为一个群信息的指针 可以通过此指针修改群信息
	roomInfo, err := this.GetRoom(roomId)
	if err != nil {
		this.infoLog.Printf("IsCreateUid GetRoom roomId=%v failed err=%v", roomId, err)
		return false, err
	}

	if roomInfo.CreateUid == uid {
		this.infoLog.Printf("IsCreateUid roomId=%v uid=%v admin=%v match", roomId, uid, roomInfo.CreateUid)
		return true, nil
	} else {
		this.infoLog.Printf("IsCreateUid roomId=%v uid=%v admin=%v not match", roomId, uid, roomInfo.CreateUid)
		return false, nil
	}
}

func (this *RoomManager) RemoveMember(roomId, operatorId, removeId uint32) (roomTS int64, err error) {
	if roomId == 0 || operatorId == 0 || removeId == 0 {
		this.infoLog.Printf("RemoveMember roomId=%v operatorId=%v removeId=%v input err", roomId, operatorId, removeId)
		return 0, ErrInputParam
	}

	// 首先更新ssdb
	err = this.dbManager.RemoveMember(roomId, removeId)
	if err != nil {
		this.infoLog.Printf("RemoveMember dbManager.RemoveMember faield roomId=%v removeId=%v", roomId, removeId)
		return 0, err
	}
	roomTS, err = this.dbManager.UpdateRoomTimeStamp(roomId)
	if err != nil {
		this.infoLog.Printf("RemoveMember dbManager.UpdateRoomTimeStamp failed roomId=%v", roomId)
		return 0, err
	}
	// roomInfo 为一个群信息的指针 可以通过此指针修改群信息
	roomInfo, err := this.GetRoom(roomId)
	if err != nil {
		this.infoLog.Printf("RemoveMember GetRoom roomId=%v failed err==%v", roomId, err)
		return 0, err
	}

	// 更新内存总的群成员列表
	totalMemberList, maxOrder, err := this.dbManager.GetRoomMemberList(roomId)
	if err != nil {
		this.infoLog.Printf("RemoveMember exec db.GetRoomMemberList()failed roomId=%v err=%v", roomId, err)
		return 0, err
	}
	// update memberlist
	this.roomInfoLock.Lock()
	defer this.roomInfoLock.Unlock()
	roomInfo.MaxOrder = maxOrder
	roomInfo.MemberList = totalMemberList
	roomInfo.RoomTS = roomTS

	// 检查被删除用户是否是管理员 如果是管理员 需要更新管理员列表
	bIsAdmin := false
	adminList := roomInfo.AdminList
	for _, v := range adminList {
		if removeId == v {
			bIsAdmin = true
			break
		}
	}
	if bIsAdmin {
		var newAdminList []uint32
		for _, v := range adminList {
			if removeId != v {
				newAdminList = append(newAdminList, v)
			}
		}

		roomTS, err = this.dbManager.SetAdminListWithUint32(roomId, newAdminList)
		if err != nil {
			this.infoLog.Printf("RemoveMember exec db.SetAdminListWithUint32 failed roomId=%v removeId=%v bIsAdmin=%v",
				roomId,
				removeId,
				bIsAdmin)
		}

		// 设置数据库失败首先更新内存确保内存中的正确
		roomInfo.AdminList = newAdminList
		roomInfo.RoomTS = roomTS
		this.infoLog.Printf("RemoveMember roomId=%v new Admin=%v", roomId, newAdminList)
	}
	return roomTS, nil
}

func (this *RoomManager) NotifyRemoveMember(roomId, createUid uint32, createName string, removeId uint32, removeName string, roomTS uint64) (err error) {
	head := &common.HeadV3{Flag: uint8(common.CServToServ),
		Version:  common.CVerMmedia,
		CryKey:   uint8(common.CNoneKey),
		TermType: uint8(0),
		Cmd:      uint16(ht_muc.MUC_CMD_TYPE_GO_CMD_NOTIFY_REMOVE_MEMBER),
		Seq:      this.GetPacketSeq(),
		From:     roomId,
		To:       0,
		Len:      0,
	}
	reqBody := new(ht_muc.MucReqBody)
	reqBody.NotifyRemoveMemberReqbody = &ht_muc.NotifyRemoveMemberReqBody{
		RoomId:        proto.Uint32(roomId),
		AdminUid:      proto.Uint32(createUid),
		AdminName:     []byte(createName),
		RemoveUid:     proto.Uint32(removeId),
		RemoveName:    []byte(removeName),
		NotifyTime:    proto.Uint32(uint32(time.Now().Unix())),
		RoomTimestamp: proto.Uint64(roomTS),
	}

	notifyBody, err := proto.Marshal(reqBody)
	if err != nil {
		this.infoLog.Printf("NotifyRemoveMember proto marshal roomId=%v createUid=%v  removeUid=%v err=%v",
			roomId,
			createUid,
			removeId,
			err)
		return err
	}

	// build xtpacket for old version
	var packetPayLoad []byte
	common.MarshalUint32(roomId, &packetPayLoad)
	common.MarshalUint32(createUid, &packetPayLoad)
	common.MarshalSlice([]byte(createName), &packetPayLoad)
	common.MarshalUint32(removeId, &packetPayLoad)
	common.MarshalSlice([]byte(removeName), &packetPayLoad)
	common.MarshalUint32(uint32(time.Now().Unix()), &packetPayLoad)
	common.MarshalUint64(uint64(roomTS), &packetPayLoad)

	//首先发送消息通知被删除者
	err = this.SendNotificationToSingleUserCompatible(roomId, removeId, head, notifyBody, packetPayLoad)
	if err != nil {
		this.infoLog.Printf("SendNotificationToSingleUser failed roomId=%v uid=%v", roomId, removeId)
	}

	//调用广播接口此时被删除用户已不再群成员列表中了
	return this.BroadCastNotification(roomId, createUid, head, notifyBody)

}
func (this *RoomManager) SendNotificationToSingleUser(roomId, uid uint32, head *common.HeadV3, notifyBody []byte) (err error) {
	if roomId == 0 || uid == 0 {
		this.infoLog.Printf("SendNotificationToSingleUser input param err roomId=%v uid=%v", roomId, uid)
		err = ErrDbParam
		return err
	}
	if this.mcApi == nil {
		this.infoLog.Printf("SendNotificationToSingleUser memcache api object not set roomId=%v", roomId)
		return ErrNilMCObject
	}

	// 用户消息处理方式 默认为存储离线
	var procType int = CSaveOffline
	memberInfo := &MemberInfoStruct{
		Uid:         uid,
		RoomId:      roomId,
		PushSetting: 0,
	}
	// 查询用户的在线状态
	onlineStat, err := this.mcApi.GetUserOnlineStat(uid)
	if err == nil {
		procType = this.GetMucMsgProcType(onlineStat, memberInfo, false)
	} else {
		this.infoLog.Printf("SendNotificationToSingleUser Get msg proc failed roomId=%v uid=%v err=%v", roomId, uid, err)
	}
	// 调整发送头部的to字段
	head.To = uid //消息的接收者
	this.infoLog.Printf("SendNotificationToSingleUser Msg proc type=%v", procType)
	switch procType {
	case CSendToIMServer:
		err := this.SendPacketToIMServerRelabile(onlineStat, head, notifyBody)
		if err != nil {
			this.infoLog.Printf("SendNotificationToSingleUser SendPacketToIMServerRelabile failed roomId=%v from=%v to=%v seq=%v",
				roomId,
				head.From,
				head.To,
				head.Seq)
			// 发送到IM失败存储离线
			ret, err := this.offlineApi.SendPacketWithHeadV3(head, notifyBody, nil)
			if err != nil {
				this.infoLog.Println("SendNotificationToSingleUser save offline faield [ret err] =", ret, err)
			} else {
				this.infoLog.Printf("SendNotificationToSingleUser save offline success roomId=%v from=%v to=%v cmd=%v seq=%v ret=%v",
					roomId,
					head.From,
					head.To,
					head.Cmd,
					head.Seq,
					ret)
			}
		} else {
			this.infoLog.Printf("SendNotificationToSingleUser SendPacketToIMServerRelabile success from=%v to=%v cmd=%v seq=%v",
				head.From,
				head.To,
				head.Cmd,
				head.Seq)
		}
	case CSaveOffLineAndPush, CSaveOffline:
		ret, err := this.offlineApi.SendPacketWithHeadV3(head, notifyBody, nil)
		if err != nil {
			this.infoLog.Println("SendNotificationToSingleUser save offline faield [ret err] =", ret, err)
		} else {
			this.infoLog.Printf("SendNotificationToSingleUser save offline success roomId=%v from=%v to=%v cmd=%v seq=%v ret=%v",
				roomId,
				head.From,
				head.To,
				head.Cmd,
				head.Seq,
				ret)
		}

	default:
		this.infoLog.Printf("Unhandle stat=%v", procType)
	}
	return err
}

func (this *RoomManager) SendNotificationToSingleUserCompatible(roomId, uid uint32, head *common.HeadV3, notifyBody []byte, oldNotifyBody []byte) (err error) {
	if roomId == 0 || uid == 0 {
		this.infoLog.Printf("SendNotificationToSingleUserCompatible input param err roomId=%v uid=%v", roomId, uid)
		err = ErrDbParam
		return err
	}
	if this.mcApi == nil {
		this.infoLog.Printf("SendNotificationToSingleUserCompatible memcache api object not set roomId=%v", roomId)
		return ErrNilMCObject
	}

	// 用户消息处理方式 默认为存储离线
	var procType int = CSaveOffline
	memberInfo := &MemberInfoStruct{
		Uid:         uid,
		RoomId:      roomId,
		PushSetting: 0,
	}
	// 查询用户的在线状态
	onlineStat, err := this.mcApi.GetUserOnlineStat(uid)
	if err == nil {
		procType = this.GetMucMsgProcType(onlineStat, memberInfo, false)
		if (onlineStat.ClientType == common.CClientTyepIOS && onlineStat.Version > common.CVerSion226) ||
			(onlineStat.ClientType == common.CClientTypeAndroid && onlineStat.Version > common.CVerSion226) {
			this.infoLog.Printf("SendNotificationToSingleUserCompatible to new version roomId=%v toUid=%v clientType=%v version=%v",
				roomId,
				uid,
				onlineStat.ClientType,
				onlineStat.Version)
			this.SendNotifyToSingleUserToNewVersion(roomId, uid, head, notifyBody, procType, onlineStat)
		} else {
			this.infoLog.Printf("SendNotificationToSingleUserCompatible to old version roomId=%v toUid=%v clientType=%v version=%v",
				roomId,
				uid,
				onlineStat.ClientType,
				onlineStat.Version)
			this.SendNotifyToSingleUserToOldVersion(roomId, uid, head, oldNotifyBody, procType, onlineStat)
		}
	} else {
		this.infoLog.Printf("SendNotificationToSingleUser Get msg proc failed roomId=%v uid=%v err=%v", roomId, uid, err)
		this.SaveOldVersionOffline(roomId, head, oldNotifyBody, uid)
	}

	return err
}

func (this *RoomManager) SendNotifyToSingleUserToNewVersion(roomId, uid uint32,
	head *common.HeadV3,
	notifyBody []byte,
	procType int,
	onlineStat *common.UserState) {
	// 调整发送头部的to字段
	head.To = uid //消息的接收者
	this.infoLog.Printf("SendNotificationToSingleUserNewVersion Msg proc type=%v", procType)
	switch procType {
	case CSendToIMServer:
		err := this.SendPacketToIMServerRelabile(onlineStat, head, notifyBody)
		if err != nil {
			this.infoLog.Printf("SendNotificationToSingleUserNewVersion SendPacketToIMServerRelabile failed roomId=%v from=%v to=%v seq=%v",
				roomId,
				head.From,
				head.To,
				head.Seq)
			// 发送到IM失败存储离线
			ret, err := this.offlineApi.SendPacketWithHeadV3(head, notifyBody, nil)
			if err != nil {
				this.infoLog.Println("SendNotificationToSingleUserNewVersion save offline faield [ret err] =", ret, err)
			} else {
				this.infoLog.Printf("SendNotificationToSingleUserNewVersion save offline success roomId=%v from=%v to=%v cmd=%v seq=%v ret=%v",
					roomId,
					head.From,
					head.To,
					head.Cmd,
					head.Seq,
					ret)
			}
		} else {
			this.infoLog.Printf("SendNotificationToSingleUserNewVersion SendPacketToIMServerRelabile success from=%v to=%v cmd=%v seq=%v",
				head.From,
				head.To,
				head.Cmd,
				head.Seq)
		}
	case CSaveOffLineAndPush, CSaveOffline:
		ret, err := this.offlineApi.SendPacketWithHeadV3(head, notifyBody, nil)
		if err != nil {
			this.infoLog.Println("SendNotificationToSingleUserNewVersion save offline faield [ret err] =", ret, err)
		} else {
			this.infoLog.Printf("SendNotificationToSingleUserNewVersion save offline success roomId=%v from=%v to=%v cmd=%v seq=%v ret=%v",
				roomId,
				head.From,
				head.To,
				head.Cmd,
				head.Seq,
				ret)
		}

	default:
		this.infoLog.Printf("Unhandle stat=%v", procType)
	}
	return
}

func (this *RoomManager) SendNotifyToSingleUserToOldVersion(roomId, uid uint32,
	head *common.HeadV3,
	notifyBody []byte,
	procType int,
	onlineStat *common.UserState) {

	var rebuildHeader common.XTHead
	rebuildHeader.Flag = head.Flag
	rebuildHeader.Version = head.Version
	rebuildHeader.CryKey = uint8(common.CNoneKey)
	rebuildHeader.TermType = head.TermType
	rebuildHeader.Cmd = uint16(this.GetOldVersionCmd(head.Cmd))
	rebuildHeader.Seq = head.Seq
	rebuildHeader.From = head.From
	rebuildHeader.To = uid // 设置ToId
	rebuildHeader.Len = uint32(len(notifyBody))

	switch procType {
	case CSendToIMServer:
		err := this.SendPacketToOldIMServerRelabile(onlineStat, &rebuildHeader, notifyBody)
		if err != nil {
			this.infoLog.Printf("SendNotifyToSingleUserToOldVersion SendPacketToOldIMServerRelabile failed roomId=%v from=%v to=%v seq=%v",
				roomId,
				rebuildHeader.From,
				rebuildHeader.To,
				rebuildHeader.Seq)
			// 发送到IM失败存储离线
			ret, err := this.offlineApi.SendPacketWithXTHead(&rebuildHeader, notifyBody, nil)
			if err != nil {
				this.infoLog.Println("SendNotifyToSingleUserToOldVersion SendPacketWithXTHead faield [ret err] =", ret, err)
			} else {
				this.infoLog.Printf("SendNotifyToSingleUserToOldVersion SendPacketWithXTHead success roomId=%v from=%v to=%v cmd=%v seq=%v ret=%v",
					roomId,
					rebuildHeader.From,
					rebuildHeader.To,
					rebuildHeader.Cmd,
					rebuildHeader.Seq,
					ret)
			}
		} else {
			this.infoLog.Printf("SendNotifyToSingleUserToOldVersion SendPacketToOldIMServerRelabile success from=%v to=%v cmd=%v seq=%v",
				rebuildHeader.From,
				rebuildHeader.To,
				rebuildHeader.Cmd,
				rebuildHeader.Seq)
		}
	case CSaveOffLineAndPush, CSaveOffline:
		ret, err := this.offlineApi.SendPacketWithXTHead(&rebuildHeader, notifyBody, nil)
		if err != nil {
			this.infoLog.Println("SendNotifyToSingleUserToOldVersion SendPacketWithXTHead faield [ret err] =", ret, err)
		} else {
			this.infoLog.Printf("SendNotifyToSingleUserToOldVersion SendPacketWithXTHead success roomId=%v from=%v to=%v cmd=%v seq=%v ret=%v",
				roomId,
				rebuildHeader.From,
				rebuildHeader.To,
				rebuildHeader.Cmd,
				rebuildHeader.Seq,
				ret)
		}
	default:
		this.infoLog.Printf("Unhandle stat=%v", procType)
	}
	return
}

func (this *RoomManager) QuitMucRoom(roomId, quitUid uint32) (newCreater uint32, roomTS int64, err error) {
	if roomId == 0 || quitUid == 0 {
		this.infoLog.Printf("QuitRoom roomId=%v  quitUid=%v input err", roomId, quitUid)
		err = ErrInputParam
		return
	}

	// roomInfo 为一个群信息的指针 可以通过此指针修改群信息
	roomInfo, err := this.GetRoom(roomId)
	if err != nil {
		this.infoLog.Printf("QuitRoom GetRoom roomId=%v failed err==%v", roomId, err)
		return 0, 0, err
	}

	this.roomInfoLock.Lock()
	defer this.roomInfoLock.Unlock()
	// 检查用户当前是否在群内 不在直接返回成功
	memberList := roomInfo.MemberList
	bInRoom := false
	for _, v := range memberList {
		if v.Uid == quitUid {
			bInRoom = true
			break
		}
	}
	if !bInRoom {
		this.infoLog.Printf("QuitRoom roomId=%v quitUid=%v not in", roomId, quitUid)
		err = ErrInputParam
		return
	}

	//在群内继续执行判断逻辑
	adminList := roomInfo.AdminList
	bIsCreater := false
	newCreater = roomInfo.CreateUid
	if quitUid == roomInfo.CreateUid {
		bIsCreater = true
		// 没有管理员 在根据order选择
		if len(adminList) == 0 {
			orderId := uint32(1000000)
			for _, v := range memberList {
				if v.OrderId < orderId && v.Uid != roomInfo.CreateUid {
					newCreater = v.Uid
					orderId = v.OrderId
				}
			}
			this.infoLog.Printf("QuitMucRoom newCreater=%v orderId=%v", newCreater, orderId)
		} else {
			//有管理员 这选择第一个管理员成员群的创建者
			for _, v := range adminList {
				if v != roomInfo.CreateUid {
					newCreater = v
					break
				}
			}
			this.infoLog.Printf("QuitMucRoom newCreater=%v", newCreater)
		}
	}

	// 检查退出用户是否是管理员 如果是管理员退出 需要更新管理员列表
	bIsAdmin := false
	for _, v := range adminList {
		if quitUid == v {
			bIsAdmin = true
			break
		}
	}
	if bIsAdmin {
		var newAdminList []uint32
		for _, v := range adminList {
			if quitUid != v {
				newAdminList = append(newAdminList, v)
			}
		}

		roomTS, err = this.dbManager.SetAdminListWithUint32(roomId, newAdminList)
		if err != nil {
			this.infoLog.Printf("QuitMucRoom exec db.SetAdminListWithUint32 failed roomId=%v quitUid=%v bIsAdmin=%v",
				roomId,
				quitUid,
				bIsAdmin)
		}

		// 设置数据库失败首先更新内存确保内存中的正确
		roomInfo.AdminList = newAdminList
		this.infoLog.Printf("QuitMucRoom roomId=%v new Admin=%v", roomId, newAdminList)
	}

	roomTS, err = this.dbManager.QuitMucRoom(roomId, quitUid, bIsCreater, newCreater)
	if err != nil {
		this.infoLog.Printf("QuitMucRoom exec db.QuitMucRoom failed roomId=%v quitUid=%v bIsCreater=%v newCreater=%v",
			roomId,
			quitUid,
			bIsCreater,
			newCreater)
	}
	// 更新完数据少数之后更新内存
	if newCreater != roomInfo.CreateUid {
		roomInfo.CreateUid = newCreater
		this.infoLog.Printf("QuitMucRoom uid=%v become roomId=%v new creater", newCreater, roomId)
	}
	// 更新内存总的群成员列表
	totalMemberList, maxOrder, err := this.dbManager.GetRoomMemberList(roomId)
	if err != nil {
		this.infoLog.Printf("RemoveMember exec db.GetRoomMemberList()failed roomId=%v err=%v", roomId, err)
		return newCreater, roomTS, err
	}

	roomInfo.MaxOrder = maxOrder
	roomInfo.MemberList = totalMemberList
	roomInfo.RoomTS = roomTS
	return newCreater, roomTS, nil
}

func (this *RoomManager) NotifyMemberQuit(roomId, quitUid uint32, quitName string, newCreate uint32, roomTS uint64) (err error) {
	head := &common.HeadV3{Flag: uint8(common.CServToServ),
		Version:  common.CVerMmedia,
		CryKey:   uint8(common.CNoneKey),
		TermType: uint8(0),
		Cmd:      uint16(ht_muc.MUC_CMD_TYPE_GO_CMD_NOTIFY_MEMBER_QUIT),
		Seq:      this.GetPacketSeq(),
		From:     roomId,
		To:       0,
		Len:      0,
	}
	reqBody := new(ht_muc.MucReqBody)
	reqBody.NotifyMemberQuitReqbody = &ht_muc.NotifyMemberQuitReqBody{
		RoomId:        proto.Uint32(roomId),
		QuitUid:       proto.Uint32(quitUid),
		QuitName:      []byte(quitName),
		AdminUid:      proto.Uint32(newCreate),
		NotifyTime:    proto.Uint32(uint32(time.Now().Unix())),
		RoomTimestamp: proto.Uint64(roomTS),
	}
	notifyBody, err := proto.Marshal(reqBody)

	if err != nil {
		this.infoLog.Printf("NotifyMemberQuit proto marshal roomId=%v adminUid=%v  quitId=%v err=%v",
			roomId,
			newCreate,
			quitUid,
			err)
		return err
	}

	this.infoLog.Printf("DEBUG NotifyMemberQuit roomId=%v adminUid=%v  quitId=%v err=%v",
		roomId,
		newCreate,
		quitUid,
		err)
	//调用广播接口
	return this.BroadCastNotification(roomId, quitUid, head, notifyBody)
}

func (this *RoomManager) ModifyRoomName(roomId, opUid uint32, roomName string) (roomTS int64, err error) {
	if roomId == 0 || opUid == 0 || len(roomName) == 0 {
		this.infoLog.Printf("ModifyRoomName roomId=%v  opUid=%v roomName=%s input err", roomId, opUid, roomName)
		err = ErrInputParam
		return
	}
	// roomInfo 为一个群信息的指针 可以通过此指针修改群信息
	roomInfo, err := this.GetRoom(roomId)
	if err != nil {
		this.infoLog.Printf("ModifyRoomName GetRoom roomId=%v failed err==%v", roomId, err)
		return roomTS, err
	}

	// 检查修改者是否以在群聊中如果不在这返回错误
	var bInRoom bool = false
	memberList := roomInfo.MemberList
	for _, v := range memberList {
		if v.Uid == opUid {
			bInRoom = true
		}
	}
	if !bInRoom {
		this.infoLog.Printf("ModifyRoomName user is not in room roomId=%v uid=%v", roomId, opUid)
		err = ErrNotInRoom
		return 0, err
	}

	this.roomInfoLock.Lock()
	defer this.roomInfoLock.Unlock()
	roomTS, err = this.dbManager.ModifyRoomName(roomId, opUid, roomName)
	if err != nil {
		this.infoLog.Printf("ModifyRoomName exec dbManager.ModifyRoomName roomId=%v opUid=%v roomName=%s", roomId, opUid, roomName)
		return roomTS, err
	}

	// 更新内存中的数据
	roomInfo.RoomName = roomName
	roomInfo.RoomTS = roomTS
	return roomTS, nil
}

func (this *RoomManager) NotifyModifyRoomName(roomId, opUid uint32, opName string, roomName string, roomTS uint64) (err error) {
	head := &common.HeadV3{Flag: uint8(common.CServToServ),
		Version:  common.CVerMmedia,
		CryKey:   uint8(common.CNoneKey),
		TermType: uint8(0),
		Cmd:      uint16(ht_muc.MUC_CMD_TYPE_GO_CMD_NOTIFY_ROOMNAME_CHANGED),
		Seq:      this.GetPacketSeq(),
		From:     roomId,
		To:       0,
		Len:      0,
	}

	reqBody := new(ht_muc.MucReqBody)
	reqBody.NotifyRoomNameChangeReqbody = &ht_muc.NotifyRoomNameChangeReqBody{
		RoomId:        proto.Uint32(roomId),
		OpUid:         proto.Uint32(opUid),
		OpName:        []byte(opName),
		RoomName:      []byte(roomName),
		NotifyTime:    proto.Uint32(uint32(time.Now().Unix())),
		RoomTimestamp: proto.Uint64(roomTS),
	}
	notifyBody, err := proto.Marshal(reqBody)

	if err != nil {
		this.infoLog.Printf("NotifyModifyRoomName proto marshal roomId=%v opUid=%v opName=%s err=%v",
			roomId,
			opUid,
			roomName,
			err)
		return err
	}

	//调用广播接口
	return this.BroadCastNotification(roomId, opUid, head, notifyBody)
}

func (this *RoomManager) ModifyMemberName(roomId, opUid uint32, opName string) (roomTS int64, err error) {
	if roomId == 0 || opUid == 0 || len(opName) == 0 {
		this.infoLog.Printf("ModifyMemberName roomId=%v  opUid=%v opName=%s input err", roomId, opUid, opName)
		err = ErrInputParam
		return roomTS, err
	}
	// roomInfo 为一个群信息的指针 可以通过此指针修改群信息
	roomInfo, err := this.GetRoom(roomId)
	if err != nil {
		this.infoLog.Printf("ModifyMemberName GetRoom roomId=%v failed err==%v", roomId, err)
		return roomTS, err
	}

	// 检查修改者是否以在群聊中如果不在这返回错误
	var bInRoom bool = false
	memberList := roomInfo.MemberList
	for _, v := range memberList {
		if v.Uid == opUid {
			bInRoom = true
		}
	}
	if !bInRoom {
		this.infoLog.Printf("ModifyMemberName user is not in room roomId=%v uid=%v", roomId, opUid)
		err = ErrNotInRoom
		return 0, err
	}

	this.roomInfoLock.Lock()
	defer this.roomInfoLock.Unlock()
	roomTS, err = this.dbManager.ModifyMemberName(roomId, opUid, opName)
	if err != nil {
		this.infoLog.Printf("ModifyMemberName exec dbManager.ModifyMemberName roomId=%v opUid=%v opName=%s", roomId, opUid, opName)
		return roomTS, err
	}

	// 更新内存中的数据
	for _, v := range memberList {
		if v.Uid == opUid {
			this.infoLog.Printf("ModifyMemberName uid=%v name=%s", opUid, opName)
			v.NickName = opName
			break
		}
	}
	roomInfo.RoomTS = roomTS
	return roomTS, nil
}

func (this *RoomManager) NotifyModifyMemberName(roomId, opUid uint32, opName string, roomTS uint64) (err error) {
	head := &common.HeadV3{Flag: uint8(common.CServToServ),
		Version:  common.CVerMmedia,
		CryKey:   uint8(common.CNoneKey),
		TermType: uint8(0),
		Cmd:      uint16(ht_muc.MUC_CMD_TYPE_GO_CMD_NOTIFY_MEMBERNAME_CHANGED),
		Seq:      this.GetPacketSeq(),
		From:     roomId,
		To:       0,
		Len:      0,
	}

	reqBody := new(ht_muc.MucReqBody)
	reqBody.NotifyMemberNameChangeReqbody = &ht_muc.NotifyMemberNameChangeReqBody{
		RoomId:        proto.Uint32(roomId),
		OpUid:         proto.Uint32(opUid),
		OpName:        []byte(opName),
		NotifyTime:    proto.Uint32(uint32(time.Now().Unix())),
		RoomTimestamp: proto.Uint64(roomTS),
	}
	notifyBody, err := proto.Marshal(reqBody)

	if err != nil {
		this.infoLog.Printf("NotifyInviteMember proto marshal roomId=%v opUid=%v opName=%v err=%v",
			roomId,
			opUid,
			opName,
			err)
		return err
	}

	//调用广播接口
	return this.BroadCastNotification(roomId, opUid, head, notifyBody)
}

func (this *RoomManager) ModifyPushSetting(roomId, opUid uint32, pushSetting uint32) (err error) {
	if roomId == 0 || opUid == 0 || pushSetting > 1 {
		this.infoLog.Printf("ModifyPushSetting roomId=%v  opUid=%v pushSetting=%s input err", roomId, opUid, pushSetting)
		err = ErrInputParam
		return err
	}
	// roomInfo 为一个群信息的指针 可以通过此指针修改群信息
	roomInfo, err := this.GetRoom(roomId)
	if err != nil {
		this.infoLog.Printf("ModifyPushSetting GetRoom roomId=%v failed err==%v", roomId, err)
		return err
	}

	// 检查修改者是否以在群聊中如果不在这返回错误
	var bInRoom bool = false
	memberList := roomInfo.MemberList
	for _, v := range memberList {
		if v.Uid == opUid {
			bInRoom = true
		}
	}
	if !bInRoom {
		this.infoLog.Printf("ModifyPushSetting user is not in room roomId=%v uid=%v", roomId, opUid)
		err = ErrNotInRoom
		return err
	}

	this.roomInfoLock.Lock()
	defer this.roomInfoLock.Unlock()
	err = this.dbManager.ModifyPushSetting(roomId, opUid, pushSetting)
	if err != nil {
		this.infoLog.Printf("ModifyPushSetting exec dbManager.ModifyPushSetting roomId=%v opUid=%v pushSetting=%v", roomId, opUid, pushSetting)
		return err
	}

	// 更新内存中的数据
	for _, v := range memberList {
		if v.Uid == opUid {
			this.infoLog.Printf("ModifyPushSetting uid=%v pushSetting=%v", opUid, pushSetting)
			v.PushSetting = pushSetting
			break
		}
	}
	return nil
}

func (this *RoomManager) GetRoomInfo(roomId uint32) (roomInfo *RoomInfo, err error) {
	if roomId == 0 {
		this.infoLog.Printf("GetRoomInfo roomId=%v input err", roomId)
		err = ErrInputParam
		return nil, err
	}
	// roomInfo 为一个群信息的指针 可以通过此指针修改群信息
	roomInfo, err = this.GetRoom(roomId)
	if err != nil {
		this.infoLog.Printf("GetRoomInfo GetRoom roomId=%v failed err=%v", roomId, err)
		return nil, err
	}
	return roomInfo, nil
}

func (this *RoomManager) BroadcastMucMessage(roomId uint32,
	head *common.HeadV3,
	msgBody []byte,
	pushInfo *MucPushInfo,
	remindType uint32,
	remindList []uint32) (err error) {
	roomInfo, err := this.GetRoom(roomId)
	if err != nil {
		this.infoLog.Printf("BroadcastMucMessage GetRoom failed roomId=%v", roomId)
		return err
	}
	if this.mcApi == nil {
		this.infoLog.Printf("BroadcastMucMessage memcache api object not set roomId=%v", roomId)
		return ErrNilMCObject
	}

	newRemidList := make([]int, len(remindList))
	for i, v := range remindList {
		newRemidList[i] = int(v)
	}
	sort.Ints(newRemidList)

	memberList := roomInfo.MemberList
	for _, v := range memberList {
		toId := v.Uid
		if toId == head.From { //群消息不需要广播给自己
			this.infoLog.Printf("BroadcastMucMessage roomId=%v toId=%v fromId=%v continue", roomId, toId, head.From)
			continue
		}
		bAtUser := false
		if remindType == AT_ALL_MEMBER {
			bAtUser = true
		} else if remindType == AT_USER_LIST && sort.SearchInts(newRemidList, int(toId)) != len(newRemidList) {
			//如果是@列表中的成员并且用户在newRemindList 中
			bAtUser = true
		}

		rebuildHeader := *head
		rebuildHeader.To = toId // 设置ToId
		rebuildHeader.CryKey = uint8(common.CNoneKey)
		if rebuildHeader.Cmd == uint16(ht_muc.MUC_CMD_TYPE_GO_CMD_GVOIP_INVITE_BROADCAST) ||
			rebuildHeader.Cmd == uint16(ht_muc.MUC_CMD_TYPE_GO_CMD_GVOIP_END_BROADCAST) {
			// nodify the from=group_id
			rebuildHeader.From = roomId
			// group voip need check receiver setting
			// if set reject voip need return
			// 0: receive  1: reject
			if this.GetVoipRejectSettin(toId) == true {
				this.infoLog.Printf("BroadcastMucMessage roomId=%v uid=%v reject voip", roomId, toId)
				continue
			}
			if this.IsUserBlockRoomVoip(roomId, toId) == true {
				this.infoLog.Printf("BroadcastMucMessage roomId=%v uid=%v person reject voip", roomId, toId)
				continue
			}
		}

		// 用户消息处理方式 默认为存储离线
		var procType int = CSaveOffline
		// 查询用户的在线状态
		onlineStat, err := this.mcApi.GetUserOnlineStat(v.Uid)
		if err == nil {
			procType = this.GetMucMsgProcType(onlineStat, v, bAtUser)
		} else {
			this.infoLog.Printf("BroadcastMucMessage Get msg proc failed roomId=%v uid=%v err=%v", roomId, v.Uid, err)
		}

		switch procType {
		case CSendToIMServer:
			err := this.SendPacketToIMServerRelabile(onlineStat, &rebuildHeader, msgBody)
			if err != nil {
				this.infoLog.Printf("BroadcastMucMessage SendPacketToIMServerRelabile failed roomId=%v from=%v to=%v seq=%v",
					roomId,
					rebuildHeader.From,
					rebuildHeader.To,
					rebuildHeader.Seq)
				// 发送到IM失败存储离线
				ret, err := this.offlineApi.SendPacketWithHeadV3(&rebuildHeader, msgBody, nil)
				if err != nil {
					this.infoLog.Println("BroadcastMucMessage CSendToIMServer faield [ret err] =", ret, err)
				} else {
					this.infoLog.Printf("BroadcastMucMessage CSendToIMServer success roomId=%v from=%v to=%v cmd=%v seq=%v ret=%v",
						roomId,
						rebuildHeader.From,
						rebuildHeader.To,
						rebuildHeader.Cmd,
						rebuildHeader.Seq,
						ret)
				}
			} else {
				this.infoLog.Printf("BroadcastMucMessage SendPacketToIMServerRelabile success from=%v to=%v cmd=%v seq=%v",
					rebuildHeader.From,
					rebuildHeader.To,
					rebuildHeader.Cmd,
					rebuildHeader.Seq)
			}
		case CSaveOffLineAndPush:
			var pushPacket []byte
			if pushInfo != nil {
				pushPacket, err = this.BuildPushPacket(onlineStat.ClientType,
					pushInfo.PushType,
					roomId,
					rebuildHeader.From,
					toId,
					pushInfo.NickName,
					pushInfo.PushParam,
					pushInfo.MsgId,
					false)
				if err != nil {
					this.infoLog.Printf("BroadcastMucMessage BuildPushPacket faild roomId=%v fromId=%v toId=%v err=%v",
						roomId,
						rebuildHeader.From,
						rebuildHeader.To,
						err)
				}
			}

			ret, err := this.offlineApi.SendPacketWithHeadV3(&rebuildHeader, msgBody, pushPacket)
			if err != nil {
				this.infoLog.Println("BroadcastMucMessage sCSaveOffLineAndPush faield [ret err] =", ret, err)
			} else {
				this.infoLog.Printf("BroadcastMucMessage CSaveOffLineAndPush success roomId=%v from=%v to=%v cmd=%v seq=%v ret=%v",
					roomId,
					rebuildHeader.From,
					rebuildHeader.To,
					rebuildHeader.Cmd,
					rebuildHeader.Seq,
					ret)
			}
		case CSaveOffline:
			ret, err := this.offlineApi.SendPacketWithHeadV3(&rebuildHeader, msgBody, nil)
			if err != nil {
				this.infoLog.Println("BroadcastMucMessage CSaveOffline faield [ret err] =", ret, err)
			} else {
				this.infoLog.Printf("BroadcastMucMessage CSaveOffline success roomId=%v from=%v to=%v cmd=%v seq=%v ret=%v",
					roomId,
					rebuildHeader.From,
					rebuildHeader.To,
					rebuildHeader.Cmd,
					rebuildHeader.Seq,
					ret)
			}
		default:
			this.infoLog.Printf("Unhandle stat=%v", procType)
		}
	}
	return nil
}

func (this *RoomManager) BroadcastMucMessageCompatible(roomId uint32,
	head *common.HeadV3,
	msgBody []byte,
	pushInfo *MucPushInfo,
	remindType uint32,
	remindList []uint32,
	oldMsgBody []byte) (err error) {
	roomInfo, err := this.GetRoom(roomId)
	if err != nil {
		this.infoLog.Printf("BroadcastMucMessageCompatible GetRoom failed roomId=%v", roomId)
		return err
	}
	if this.mcApi == nil {
		this.infoLog.Printf("BroadcastMucMessageCompatible memcache api object not set roomId=%v", roomId)
		return ErrNilMCObject
	}

	newRemidList := make([]int, len(remindList))
	for i, v := range remindList {
		newRemidList[i] = int(v)
	}
	sort.Ints(newRemidList)

	memberList := roomInfo.MemberList
	// 如果用户不在群聊中 不允许发送群消息
	var bInRoom bool = false
	for _, v := range memberList {
		if v.Uid == head.From {
			bInRoom = true
		}
	}
	if !bInRoom {
		this.infoLog.Printf("BroadcastMucMessageCompatible user is not in room roomId=%v uid=%v", roomId, head.From)
		err = ErrNotInRoom
		return err
	}

	for _, v := range memberList {
		toId := v.Uid
		if toId == head.From { //群消息不需要广播给自己
			this.infoLog.Printf("BroadcastMucMessageCompatible roomId=%v toId=%v fromId=%v continue", roomId, toId, head.From)
			continue
		}
		bAtUser := false
		if remindType == AT_ALL_MEMBER {
			bAtUser = true
		} else if remindType == AT_USER_LIST && sort.SearchInts(newRemidList, int(toId)) != len(newRemidList) {
			//如果是@列表中的成员并且用户在newRemindList 中
			bAtUser = true
		}

		// 用户消息处理方式 默认为存储离线
		var procType int = CSaveOffline
		// 查询用户的在线状态
		onlineStat, err := this.mcApi.GetUserOnlineStat(v.Uid)
		if err == nil {
			procType = this.GetMucMsgProcType(onlineStat, v, bAtUser)
			if (onlineStat.ClientType == common.CClientTyepIOS && onlineStat.Version > common.CVerSion226) ||
				(onlineStat.ClientType == common.CClientTypeAndroid && onlineStat.Version > common.CVerSion226) {
				this.infoLog.Printf("BroadcastMucMessageCompatible to new version roomId=%v toUid=%v clientType=%v version=%v",
					roomId,
					toId,
					onlineStat.ClientType,
					onlineStat.Version)
				this.BroadcastMucMessageToNewVersion(roomId, head, msgBody, pushInfo, procType, toId, onlineStat, bAtUser)
			} else {
				this.infoLog.Printf("BroadcastMucMessageCompatible to old version roomId=%v toUid=%v clientType=%v version=%v",
					roomId,
					toId,
					onlineStat.ClientType,
					onlineStat.Version)
				this.BroadcastMucMessageToOldVersion(roomId, head, oldMsgBody, pushInfo, procType, toId, onlineStat)
			}
		} else {
			this.infoLog.Printf("BroadcastMucMessageCompatible Get msg proc failed safe offline roomId=%v uid=%v err=%v", roomId, v.Uid, err)
			this.SaveOldVersionOffline(roomId, head, oldMsgBody, toId)
		}

	}
	return nil
}

func (this *RoomManager) BroadcastMucMessageToNewVersion(roomId uint32,
	head *common.HeadV3,
	msgBody []byte,
	pushInfo *MucPushInfo,
	procType int,
	toId uint32,
	onlineStat *common.UserState,
	bAtUser bool) {
	rebuildHeader := *head
	rebuildHeader.To = toId // 设置ToId
	rebuildHeader.CryKey = uint8(common.CNoneKey)
	if rebuildHeader.Cmd == uint16(ht_muc.MUC_CMD_TYPE_GO_CMD_GVOIP_INVITE_BROADCAST) ||
		rebuildHeader.Cmd == uint16(ht_muc.MUC_CMD_TYPE_GO_CMD_GVOIP_END_BROADCAST) {
		// nodify the from=group_id
		rebuildHeader.From = roomId
		// group voip need check receiver setting
		// if set reject voip need return
		// 0: receive  1: reject
		if this.GetVoipRejectSettin(toId) == true {
			this.infoLog.Printf("BroadcastMucMessageToNewVersion roomId=%v uid=%v reject voip", roomId, toId)
			return
		}
		if this.IsUserBlockRoomVoip(roomId, toId) == true {
			this.infoLog.Printf("BroadcastMucMessageToNewVersion roomId=%v uid=%v person reject voip", roomId, toId)
			return
		}
	}

	switch procType {
	case CSendToIMServer:
		err := this.SendPacketToIMServerRelabile(onlineStat, &rebuildHeader, msgBody)
		if err != nil {
			this.infoLog.Printf("BroadcastMucMessageToNewVersion SendPacketToIMServerRelabile failed roomId=%v from=%v to=%v seq=%v",
				roomId,
				rebuildHeader.From,
				rebuildHeader.To,
				rebuildHeader.Seq)
			// 发送到IM失败存储离线
			ret, err := this.offlineApi.SendPacketWithHeadV3(&rebuildHeader, msgBody, nil)
			if err != nil {
				this.infoLog.Println("BroadcastMucMessageToNewVersion CSendToIMServer faield [ret err] =", ret, err)
			} else {
				this.infoLog.Printf("BroadcastMucMessageToNewVersion CSendToIMServer success roomId=%v from=%v to=%v cmd=%v seq=%v ret=%v",
					roomId,
					rebuildHeader.From,
					rebuildHeader.To,
					rebuildHeader.Cmd,
					rebuildHeader.Seq,
					ret)
			}
		} else {
			this.infoLog.Printf("BroadcastMucMessageToNewVersion SendPacketToIMServerRelabile success from=%v to=%v cmd=%v seq=%v",
				rebuildHeader.From,
				rebuildHeader.To,
				rebuildHeader.Cmd,
				rebuildHeader.Seq)
		}
	case CSaveOffLineAndPush:
		var pushPacket []byte
		var err error
		if pushInfo != nil {
			pushPacket, err = this.BuildPushPacket(onlineStat.ClientType,
				pushInfo.PushType,
				roomId,
				rebuildHeader.From,
				toId,
				pushInfo.NickName,
				pushInfo.PushParam,
				pushInfo.MsgId,
				bAtUser)
			if err != nil {
				this.infoLog.Printf("BroadcastMucMessageToNewVersion BuildPushPacket faild roomId=%v fromId=%v toId=%v err=%v",
					roomId,
					rebuildHeader.From,
					rebuildHeader.To,
					err)
			}
		}

		ret, err := this.offlineApi.SendPacketWithHeadV3(&rebuildHeader, msgBody, pushPacket)
		if err != nil {
			this.infoLog.Println("BroadcastMucMessageToNewVersion sCSaveOffLineAndPush faield [ret err] =", ret, err)
		} else {
			this.infoLog.Printf("BroadcastMucMessageToNewVersion CSaveOffLineAndPush success roomId=%v from=%v to=%v cmd=%v seq=%v ret=%v",
				roomId,
				rebuildHeader.From,
				rebuildHeader.To,
				rebuildHeader.Cmd,
				rebuildHeader.Seq,
				ret)
		}
	case CSaveOffline:
		ret, err := this.offlineApi.SendPacketWithHeadV3(&rebuildHeader, msgBody, nil)
		if err != nil {
			this.infoLog.Println("BroadcastMucMessageToNewVersion CSaveOffline faield [ret err] =", ret, err)
		} else {
			this.infoLog.Printf("BroadcastMucMessageToNewVersion CSaveOffline success roomId=%v from=%v to=%v cmd=%v seq=%v ret=%v",
				roomId,
				rebuildHeader.From,
				rebuildHeader.To,
				rebuildHeader.Cmd,
				rebuildHeader.Seq,
				ret)
		}
	default:
		this.infoLog.Printf("BroadcastMucMessageToNewVersion Unhandle stat=%v", procType)
	}
	return
}

func (this *RoomManager) BroadcastMucMessageToOldVersion(roomId uint32,
	head *common.HeadV3,
	msgBody []byte,
	pushInfo *MucPushInfo,
	procType int,
	toId uint32,
	onlineStat *common.UserState) {
	var rebuildHeader common.XTHead
	rebuildHeader.Flag = head.Flag
	rebuildHeader.Version = head.Version
	rebuildHeader.CryKey = uint8(common.CNoneKey)
	rebuildHeader.TermType = head.TermType
	rebuildHeader.Cmd = uint16(this.GetOldVersionCmd(head.Cmd))
	rebuildHeader.Seq = head.Seq
	rebuildHeader.From = head.From
	rebuildHeader.To = toId // 设置ToId
	rebuildHeader.Len = uint32(len(msgBody))

	if rebuildHeader.Cmd == uint16(CMD_GVOIP_INVITE_BROADCAST) ||
		rebuildHeader.Cmd == uint16(CMD_GVOIP_END_BROADCAST) {
		// nodify the from=group_id
		rebuildHeader.From = roomId
		// group voip need check receiver setting
		// if set reject voip need return
		// 0: receive  1: reject
		if this.GetVoipRejectSettin(toId) == true {
			this.infoLog.Printf("BroadcastMucMessageToOldVersion roomId=%v uid=%v reject voip", roomId, toId)
			return
		}
		if this.IsUserBlockRoomVoip(roomId, toId) == true {
			this.infoLog.Printf("BroadcastMucMessageToOldVersion roomId=%v uid=%v person reject voip", roomId, toId)
			return
		}
	}

	switch procType {
	case CSendToIMServer:
		err := this.SendPacketToOldIMServerRelabile(onlineStat, &rebuildHeader, msgBody)
		if err != nil {
			this.infoLog.Printf("BroadcastMucMessageToOldVersion SendPacketToOldIMServerRelabile failed roomId=%v from=%v to=%v seq=%v",
				roomId,
				rebuildHeader.From,
				rebuildHeader.To,
				rebuildHeader.Seq)
			// 发送到IM失败存储离线
			ret, err := this.offlineApi.SendPacketWithXTHead(&rebuildHeader, msgBody, nil)
			if err != nil {
				this.infoLog.Println("BroadcastMucMessageToOldVersion SendPacketWithXTHead faield [ret err] =", ret, err)
			} else {
				this.infoLog.Printf("BroadcastMucMessageToOldVersion SendPacketWithXTHead success roomId=%v from=%v to=%v cmd=%v seq=%v ret=%v",
					roomId,
					rebuildHeader.From,
					rebuildHeader.To,
					rebuildHeader.Cmd,
					rebuildHeader.Seq,
					ret)
			}
		} else {
			this.infoLog.Printf("BroadcastMucMessageToOldVersion SendPacketToOldIMServerRelabile success from=%v to=%v cmd=%v seq=%v",
				rebuildHeader.From,
				rebuildHeader.To,
				rebuildHeader.Cmd,
				rebuildHeader.Seq)
		}
	case CSaveOffLineAndPush:
		var pushPacket []byte
		var err error
		if pushInfo != nil {
			pushPacket, err = this.BuildPushPacket(onlineStat.ClientType,
				pushInfo.PushType,
				roomId,
				rebuildHeader.From,
				toId,
				pushInfo.NickName,
				pushInfo.PushParam,
				pushInfo.MsgId,
				false)
			if err != nil {
				this.infoLog.Printf("BroadcastMucMessageToOldVersion BuildPushPacket faild roomId=%v fromId=%v toId=%v err=%v",
					roomId,
					rebuildHeader.From,
					rebuildHeader.To,
					err)
			}
		}

		ret, err := this.offlineApi.SendPacketWithXTHead(&rebuildHeader, msgBody, pushPacket)
		if err != nil {
			this.infoLog.Println("BroadcastMucMessageToOldVersion SendPacketWithXTHead faield [ret err] =", ret, err)
		} else {
			this.infoLog.Printf("BroadcastMucMessageToOldVersion SendPacketWithXTHead success roomId=%v from=%v to=%v cmd=%v seq=%v ret=%v",
				roomId,
				rebuildHeader.From,
				rebuildHeader.To,
				rebuildHeader.Cmd,
				rebuildHeader.Seq,
				ret)
		}
	case CSaveOffline:
		ret, err := this.offlineApi.SendPacketWithXTHead(&rebuildHeader, msgBody, nil)
		if err != nil {
			this.infoLog.Println("BroadcastMucMessageToOldVersion SendPacketWithXTHead faield [ret err] =", ret, err)
		} else {
			this.infoLog.Printf("BroadcastMucMessageToOldVersion SendPacketWithXTHead success roomId=%v from=%v to=%v cmd=%v seq=%v ret=%v",
				roomId,
				rebuildHeader.From,
				rebuildHeader.To,
				rebuildHeader.Cmd,
				rebuildHeader.Seq,
				ret)
		}
	default:
		this.infoLog.Printf("Unhandle stat=%v", procType)
	}
	return
}

func (this *RoomManager) SaveOldVersionOffline(roomId uint32, head *common.HeadV3, msgBody []byte, toId uint32) {
	var rebuildHeader common.XTHead
	rebuildHeader.Flag = head.Flag
	rebuildHeader.Version = head.Version
	rebuildHeader.CryKey = uint8(common.CNoneKey)
	rebuildHeader.TermType = head.TermType
	rebuildHeader.Cmd = uint16(CMD_MUC_MESSAGE)
	rebuildHeader.Seq = head.Seq
	rebuildHeader.From = head.From
	rebuildHeader.To = toId // 设置ToId
	rebuildHeader.Len = uint32(len(msgBody))

	ret, err := this.offlineApi.SendPacketWithXTHead(&rebuildHeader, msgBody, nil)
	if err != nil {
		this.infoLog.Println("SaveOldVersionOffline SendPacketWithXTHead faield [ret err] =", ret, err)
	} else {
		this.infoLog.Printf("SaveOldVersionOffline SendPacketWithXTHead success roomId=%v from=%v to=%v cmd=%v seq=%v ret=%v",
			roomId,
			rebuildHeader.From,
			rebuildHeader.To,
			rebuildHeader.Cmd,
			rebuildHeader.Seq,
			ret)
	}
	return
}

func (this *RoomManager) IsUserBlockRoomVoip(roomId, uid uint32) (bBlock bool) {
	bBlock = false
	if uidSlice, ok := this.roomIdToUser[roomId]; !ok {
		this.voipLock.Lock()
		defer this.voipLock.Unlock()
		outList, err := this.GetBlockRoomVoipUserList(roomId)
		if err != nil {
			this.infoLog.Printf("IsUserBlockRoomVoip roomId=%v uid=%v exec faild", roomId, uid)
			bBlock = true // 执行失败当做是用户阻塞了群voip通知
			return bBlock
		}
		this.roomIdToUser[roomId] = outList // 添加到列表中
		for _, v := range outList {
			if v == uid {
				bBlock = true
			}
		}
	} else { // 已经存在RoomId 直接查找
		for _, v := range uidSlice {
			if v == uid {
				bBlock = true
			}
		}
	}
	return bBlock
}

// 注意这个函数没有将数据库返回的结果添加到roomIdToUser 中
func (this *RoomManager) GetBlockRoomVoipUserList(roomId uint32) (outList []uint32, err error) {
	if roomId == 0 {
		this.infoLog.Printf("GetBlockRoomVoipUserList roomId=%v input err", roomId)
		err = ErrInputParam
		return outList, err
	}
	outList, err = this.dbManager.GetBlockRoomVoipUserList(roomId)
	if err != nil {
		this.infoLog.Printf("GetBlockRoomVoipUserList exec db.GetBlockRoomVoipUserList failed roomId=%v err=%v", roomId, err)
		return nil, err
	}
	return outList, nil
}

func (this *RoomManager) GetVoipRejectSettin(uid uint32) (bReject bool) {
	bReject = false
	bReject, err := this.dbManager.GetVoipRejectSetting(uid)
	if err != nil {
		this.infoLog.Printf("GetVoipRejectSettin uid=%v err=%v faied", uid, err)
	}
	return bReject
}

func (this *RoomManager) BuildPushPacket(terminalType, pushType uint8, roomId, fromId, toId uint32, nickName, pushContent, msgId string, bAtUser bool) (outPacket []byte, err error) {
	this.infoLog.Printf("BuildPushPacket roomId=%v fromId=%v toId=%v terminalType=%v pushType=%v nickName=%s msgId=%s",
		roomId,
		fromId,
		toId,
		terminalType,
		pushType,
		nickName,
		msgId)
	if roomId == 0 || fromId == 0 || toId == 0 || terminalType > 1 {
		this.infoLog.Printf("BuildPushPacket invalid param")
		err = ErrInputParam
		return nil, err
	}
	// MUC推送的默认参数
	var chatType uint8 = uint8(CT_MUC)
	var sound uint8 = uint8(1)
	var lights uint8 = uint8(1)

	limitNickeName := make([]byte, NICKNAME_LEN)
	copy(limitNickeName, []byte(nickName))

	limitPushContent := make([]byte, CONTENT_LEN)
	copy(limitPushContent, []byte(pushContent))

	limitMsgId := make([]byte, CONTENT_LEN)
	copy(limitMsgId, []byte(msgId))

	var packetPayLoad []byte
	common.MarshalUint8(terminalType, &packetPayLoad)
	common.MarshalUint8(chatType, &packetPayLoad)
	common.MarshalUint32(fromId, &packetPayLoad)
	common.MarshalUint32(toId, &packetPayLoad)
	common.MarshalUint32(roomId, &packetPayLoad)
	common.MarshalUint8(pushType, &packetPayLoad)
	common.MarshalSlice(limitNickeName, &packetPayLoad)
	common.MarshalSlice(limitPushContent, &packetPayLoad)
	common.MarshalUint8(sound, &packetPayLoad)
	common.MarshalUint8(lights, &packetPayLoad)
	common.MarshalSlice(limitMsgId, &packetPayLoad)
	var actionId uint32
	var byAt uint8 = uint8(CNotBeenAt)
	if bAtUser {
		byAt = uint8(CBeenAt)
	}
	common.MarshalUint32(actionId, &packetPayLoad)
	common.MarshalUint8(byAt, &packetPayLoad)

	head := &common.XTHead{
		Flag:     common.CServToServ,
		Version:  common.CVerMmedia,
		CryKey:   common.CServKey,
		TermType: 0,
		Cmd:      CMD_S2S_MESSAGE_PUSH,
		Seq:      this.GetPacketSeq(),
		From:     0,
		To:       0,
		Len:      0,
	}
	// 使用Server key 加密
	cryptoText := libcrypto.TEAEncrypt(string(packetPayLoad), SERVER_COMM_KEY)
	this.infoLog.Printf("BuildPushPacket cryptoText len=%v", len(cryptoText))

	head.Len = uint32(len(cryptoText)) //
	outPacket = make([]byte, common.XTHeadLen+head.Len)
	err = common.SerialXTHeadToSlice(head, outPacket[:])
	if err != nil {
		this.infoLog.Println("BuildPushPacket SerialXTHeadToSlice failed")
		return nil, err
	}
	copy(outPacket[common.XTHeadLen:], []byte(cryptoText)) // return code
	return outPacket, nil
}

func (this *RoomManager) AddRoomToContactList(roomId, opUid, opType uint32) (err error) {
	if roomId == 0 || opUid == 0 || opType > 1 {
		this.infoLog.Printf("AddRoomToContactList invalid param roomId=%v opUid=%v opType=%v",
			roomId,
			opUid,
			opType)
		err = ErrInputParam
		return err
	}
	err = this.dbManager.AddRoomToContactList(roomId, opUid, opType)
	if err != nil {
		this.infoLog.Printf("AddRoomToContactList roomId=%v opUid=%v opType=%v exec dbManager.AddRoomToContactList failce",
			roomId,
			opUid,
			opType)
	}
	return err
}

func (this *RoomManager) GetRoomFromContactList(opUid uint32) (roomList []*ht_muc.RoomInfoBody, err error) {
	if opUid == 0 {
		this.infoLog.Printf("GetRoomFromContactList  invalid param opUid=%v", opUid)
		err = ErrInputParam
		return nil, err
	}

	idList, err := this.dbManager.GetAllContactListRoomId(opUid)
	if err != nil {
		this.infoLog.Printf("GetRoomFromContactList exec dbManager.GetAllContactListRoomId faied opUid=%v err=%v",
			opUid,
			err)
		return nil, err
	}
	for _, v := range idList {
		this.infoLog.Printf("GetRoomFromContactList GetRoom roomId=%v", v)
		roomId := v
		roomInfo, err := this.GetRoom(roomId)
		if err != nil {
			this.infoLog.Printf("GetRoomFromContactList GetRoom failed roomId=%v err=%v", roomId, err)
			continue
		}
		var memberInfo []*ht_muc.RoomMemberInfo
		memberList := roomInfo.MemberList
		var pushSetting uint32 = 0
		for _, v := range memberList {
			if v.Uid == opUid {
				pushSetting = v.PushSetting
			}
		}
		verifystat := ht_muc.VERIFY_STAT(roomInfo.VerifyStat)
		roomOut := &ht_muc.RoomInfoBody{
			RoomId:       proto.Uint32(roomId),
			CreateUid:    proto.Uint32(roomInfo.CreateUid),
			ListAdminUid: roomInfo.AdminList,
			RoomLimit:    proto.Uint32(roomInfo.MemberLimit),
			RoomName:     []byte(roomInfo.RoomName),
			RoomDesc:     []byte(roomInfo.RoomDesc),
			VerifyStat:   &verifystat,
			Announcement: &ht_muc.AnnoType{
				PublishUid:  proto.Uint32(roomInfo.Announcement.PublishUid),
				PublishTs:   proto.Uint32(roomInfo.Announcement.PublishTS),
				AnnoContent: []byte(roomInfo.Announcement.AnnoContect),
			},
			RoomTimestamp: proto.Uint64(uint64(roomInfo.RoomTS)),
			PushSetting:   proto.Uint32(pushSetting),
			Members:       memberInfo,
		}
		// 添加到输出的roomList 中
		roomList = append(roomList, roomOut)
	}

	return roomList, nil
}

func (this *RoomManager) UpdateVoipBlockSetting(opUid, blockId, blockType, action uint32) (err error) {
	if opUid == 0 || blockId == 0 {
		this.infoLog.Printf("UpdateVoipBlockSetting input param err")
		err = ErrInputParam
		return err
	}

	// roomInfo 为一个群信息的指针 可以通过此指针修改群信息
	roomInfo, err := this.GetRoom(blockId)
	if err != nil {
		this.infoLog.Printf("UpdateVoipBlockSetting GetRoom roomId=%v failed err==%v", blockId, err)
		return err
	}

	// 检查修改者是否以在群聊中如果不在这返回错误
	var bInRoom bool = false
	memberList := roomInfo.MemberList
	for _, v := range memberList {
		if v.Uid == opUid {
			bInRoom = true
		}
	}
	if !bInRoom {
		this.infoLog.Printf("UpdateVoipBlockSetting user is not in room roomId=%v uid=%v", blockId, opUid)
		err = ErrNotInRoom
		return err
	}

	// 首先更新db中参数
	err = this.dbManager.UpdateVoipBlockList(opUid, blockId, blockType, action)
	if err != nil {
		this.infoLog.Printf("UpdateVoipBlockSetting exec dbManager.UpdateVoipBlockSetting faied opUid=%v blockId=%v blockType=%v action=%v err=%v",
			opUid,
			blockId,
			blockType,
			action,
			err)
		return err
	}

	// 更新内存中 群聊voip设置
	this.voipLock.Lock()
	defer this.voipLock.Unlock()
	if _, ok := this.roomIdToUser[blockId]; !ok { // 没有加载 则从数据库load 读取到的即为最新用户设置
		outList, err := this.GetBlockRoomVoipUserList(blockId)
		if err != nil {
			this.infoLog.Printf("UpdateVoipBlockSetting roomId=%v uid=%v exec faild err=%v", blockId, opUid, err)
			return err
		}
		this.roomIdToUser[blockId] = outList // 添加到列表中
	} else { // 已经存在则根据设置进行修改 屏蔽群voip这添加用户 否则 删除用户
		if action == 1 { // add to block list
			this.roomIdToUser[blockId] = append(this.roomIdToUser[blockId], opUid)
		} else { // otherwise delete from block list
			var blockList []uint32
			for _, v := range this.roomIdToUser[blockId] {
				if v == opUid {
					continue
				}
				blockList = append(blockList, v)
			}
			this.roomIdToUser[blockId] = blockList
		}
	}
	err = nil
	this.infoLog.Printf("RoomManager UpdateVoipBlockSetting opUid=%v blockId=%v blockType=%v action=%v",
		opUid,
		blockId,
		blockType,
		action)
	return err
}

func (this *RoomManager) MultiCastNotificationCompatible(targetList []uint32,
	head *common.HeadV3,
	notifyBody []byte,
	onlyOnline bool,
	oldNotifyBody []byte) (err error) {
	if this.mcApi == nil {
		this.infoLog.Printf("MultiCastNotificationCompatible memcache api object not set ")
		return ErrNilMCObject
	}
	for _, v := range targetList {
		// 用户消息处理方式 默认为存储离线
		var procType int = CSaveOffline
		// 查询用户的在线状态
		onlineStat, err := this.mcApi.GetUserOnlineStat(v)
		if err == nil {
			memberInfo := &MemberInfoStruct{Uid: v}
			procType = this.GetMucMsgProcType(onlineStat, memberInfo, false)
			if (onlineStat.ClientType == common.CClientTyepIOS && onlineStat.Version > common.CVerSion226) ||
				(onlineStat.ClientType == common.CClientTypeAndroid && onlineStat.Version > common.CVerSion226) {
				this.infoLog.Printf("MultiCastNotificationCompatible to new version toUid=%v clientType=%v version=%v",
					v,
					onlineStat.ClientType,
					onlineStat.Version)
				this.MultiCastNotificationToNewVersion(head, notifyBody, procType, v, onlineStat, onlyOnline)
			} else {
				this.infoLog.Printf("MultiCastNotificationCompatible to old version toUid=%v clientType=%v version=%v",
					v,
					onlineStat.ClientType,
					onlineStat.Version)
				this.MultiCastNotificationToOldVersion(head, oldNotifyBody, procType, v, onlineStat, onlyOnline)
			}
		} else {
			this.infoLog.Printf("MultiCastNotificationCompatible Get msg proc failed safe offline uid=%v err=%v", v, err)
			if !onlyOnline {
				// 获取不到在线，根据onlyOnlie 存储老版本的离线
				this.SaveOldVersionOffline(0, head, notifyBody, v)
			}
		}
	}
	return nil
}

func (this *RoomManager) MultiCastNotificationToNewVersion(head *common.HeadV3,
	notifyBody []byte,
	procType int,
	toId uint32,
	onlineStat *common.UserState,
	onlyOnline bool) (err error) {

	// 调整发送头部的to字段
	head.To = toId //消息的接收者
	this.infoLog.Printf("MultiCastNotificationToNewVersion Msg proc type=%v", procType)

	switch procType {
	case CSendToIMServer:
		err := this.SendPacketToIMServerRelabile(onlineStat, head, notifyBody)
		if err != nil {
			this.infoLog.Printf("MultiCastNotificationToNewVersion SendPacketToIMServerRelabile failed  from=%v to=%v seq=%v",
				head.From,
				head.To,
				head.Seq)
			// 发送到IM失败存储离线
			if !onlyOnline {
				ret, err := this.offlineApi.SendPacketWithHeadV3(head, notifyBody, nil)
				if err != nil {
					this.infoLog.Println("MultiCastNotificationToNewVersion save offline faield [ret err] =", ret, err)
				} else {
					this.infoLog.Printf("MultiCastNotificationToNewVersion save offline success from=%v to=%v cmd=%v seq=%v ret=%v",
						head.From,
						head.To,
						head.Cmd,
						head.Seq,
						ret)
				}
			}
		} else {
			this.infoLog.Printf("MultiCastNotificationToNewVersion SendPacketToIMServerRelabile success from=%v to=%v cmd=%v seq=%v",
				head.From,
				head.To,
				head.Cmd,
				head.Seq)
		}
	case CSaveOffLineAndPush, CSaveOffline:
		if !onlyOnline {
			ret, err := this.offlineApi.SendPacketWithHeadV3(head, notifyBody, nil)
			if err != nil {
				this.infoLog.Println("MultiCastNotificationToNewVersion save offline faield [ret err] =", ret, err)
			} else {
				this.infoLog.Printf("MultiCastNotificationToNewVersion save offline success from=%v to=%v cmd=%v seq=%v ret=%v",
					head.From,
					head.To,
					head.Cmd,
					head.Seq,
					ret)
			}
		}
	default:
		this.infoLog.Printf("MultiCastNotificationToNewVersion Unhandle stat=%v", procType)
	}
	return
}

func (this *RoomManager) MultiCastNotificationToOldVersion(head *common.HeadV3,
	notifyBody []byte,
	procType int,
	toId uint32,
	onlineStat *common.UserState,
	onlyOnline bool) {
	var rebuildHeader common.XTHead
	rebuildHeader.Flag = head.Flag
	rebuildHeader.Version = head.Version
	rebuildHeader.CryKey = uint8(common.CNoneKey)
	rebuildHeader.TermType = head.TermType
	rebuildHeader.Cmd = uint16(this.GetOldVersionCmd(head.Cmd))
	rebuildHeader.Seq = head.Seq
	rebuildHeader.From = head.From
	rebuildHeader.To = toId // 设置ToId
	rebuildHeader.Len = uint32(len(notifyBody))

	switch procType {
	case CSendToIMServer:
		err := this.SendPacketToOldIMServerRelabile(onlineStat, &rebuildHeader, notifyBody)
		if err != nil {
			this.infoLog.Printf("MultiCastNotificationToOldVersion SendPacketToOldIMServerRelabile failed from=%v to=%v seq=%v",
				rebuildHeader.From,
				rebuildHeader.To,
				rebuildHeader.Seq)
			// 发送到IM失败存储离线
			if !onlyOnline {
				ret, err := this.offlineApi.SendPacketWithXTHead(&rebuildHeader, notifyBody, nil)
				if err != nil {
					this.infoLog.Println("MultiCastNotificationToOldVersion SendPacketWithXTHead faield [ret err] =", ret, err)
				} else {
					this.infoLog.Printf("MultiCastNotificationToOldVersion SendPacketWithXTHead success from=%v to=%v cmd=%v seq=%v ret=%v",
						rebuildHeader.From,
						rebuildHeader.To,
						rebuildHeader.Cmd,
						rebuildHeader.Seq,
						ret)
				}
			}

		} else {
			this.infoLog.Printf("MultiCastNotificationToOldVersion SendPacketToOldIMServerRelabile success from=%v to=%v cmd=%v seq=%v",
				rebuildHeader.From,
				rebuildHeader.To,
				rebuildHeader.Cmd,
				rebuildHeader.Seq)
		}
	case CSaveOffLineAndPush, CSaveOffline:
		if !onlyOnline {
			ret, err := this.offlineApi.SendPacketWithXTHead(&rebuildHeader, notifyBody, nil)
			if err != nil {
				this.infoLog.Println("MultiCastNotificationToOldVersion SendPacketWithXTHead faield [ret err] =", ret, err)
			} else {
				this.infoLog.Printf("MultiCastNotificationToOldVersion SendPacketWithXTHead success from=%v to=%v cmd=%v seq=%v ret=%v",
					rebuildHeader.From,
					rebuildHeader.To,
					rebuildHeader.Cmd,
					rebuildHeader.Seq,
					ret)
			}
		}

	default:
		this.infoLog.Printf("Unhandle stat=%v", procType)
	}
	return
}

func (this *RoomManager) GetOldVersionCmd(newCmd uint16) (oldCmd uint16) {
	cmdType := ht_muc.MUC_CMD_TYPE(newCmd)
	switch cmdType {
	case ht_muc.MUC_CMD_TYPE_GO_CMD_MUC_MESSAGE:
		oldCmd = uint16(CMD_MUC_MESSAGE)
	case ht_muc.MUC_CMD_TYPE_GO_CMD_GVOIP_INVITE_BROADCAST:
		oldCmd = uint16(CMD_GVOIP_INVITE_BROADCAST)
	case ht_muc.MUC_CMD_TYPE_GO_CMD_GVOIP_MEMBER_JOIN_BROADCAST:
		oldCmd = uint16(CMD_GVOIP_MEMBER_JOIN_BROADCAST)
	case ht_muc.MUC_CMD_TYPE_GO_CMD_GVOIP_MEMBER_LEAVE_BROADCAST:
		oldCmd = uint16(CMD_GVOIP_MEMBER_LEAVE_BROADCAST)
	case ht_muc.MUC_CMD_TYPE_GO_CMD_GVOIP_END_BROADCAST:
		oldCmd = uint16(CMD_GVOIP_END_BROADCAST)
	case ht_muc.MUC_CMD_TYPE_GO_CMD_NOTIFY_REMOVE_MEMBER:
		oldCmd = uint16(CMD_NOTIFY_REMOVE_MEMBER)
	default:
		oldCmd = newCmd
	}
	this.infoLog.Printf("GetOldVersionCmd newCmd=%v change to oldCmd=%v", newCmd, oldCmd)
	return oldCmd
}
func (this *RoomManager) MultiCastNotification(targetList []uint32, head *common.HeadV3, notifyBody []byte, onlyOnline bool) (err error) {
	if this.mcApi == nil {
		this.infoLog.Printf("MultiCastNotification memcache api object not set ")
		return ErrNilMCObject
	}

	for _, v := range targetList {
		// 用户消息处理方式 默认为存储离线
		var procType int = CSaveOffline
		// 查询用户的在线状态
		onlineStat, err := this.mcApi.GetUserOnlineStat(v)
		if err == nil {
			memberInfo := &MemberInfoStruct{Uid: v}
			procType = this.GetMucMsgProcType(onlineStat, memberInfo, false)
		} else {
			this.infoLog.Printf("MultiCastNotification Get msg proc failed uid=%v err=%v", v, err)
		}
		// 调整发送头部的to字段
		head.To = v //消息的接收者
		this.infoLog.Printf("MultiCastNotification Msg proc type=%v", procType)

		switch procType {
		case CSendToIMServer:
			err := this.SendPacketToIMServerRelabile(onlineStat, head, notifyBody)
			if err != nil {
				this.infoLog.Printf("MultiCastNotification SendPacketToIMServerRelabile failed  from=%v to=%v seq=%v",
					head.From,
					head.To,
					head.Seq)
				// 发送到IM失败存储离线
				if !onlyOnline {
					ret, err := this.offlineApi.SendPacketWithHeadV3(head, notifyBody, nil)
					if err != nil {
						this.infoLog.Println("MultiCastNotification save offline faield [ret err] =", ret, err)
					} else {
						this.infoLog.Printf("MultiCastNotification save offline success from=%v to=%v cmd=%v seq=%v ret=%v",
							head.From,
							head.To,
							head.Cmd,
							head.Seq,
							ret)
					}
				}
			} else {
				this.infoLog.Printf("MultiCastNotification SendPacketToIMServerRelabile success from=%v to=%v cmd=%v seq=%v",
					head.From,
					head.To,
					head.Cmd,
					head.Seq)
			}
		case CSaveOffLineAndPush, CSaveOffline:
			if !onlyOnline {
				ret, err := this.offlineApi.SendPacketWithHeadV3(head, notifyBody, nil)
				if err != nil {
					this.infoLog.Println("MultiCastNotification save offline faield [ret err] =", ret, err)
				} else {
					this.infoLog.Printf("MultiCastNotification save offline success from=%v to=%v cmd=%v seq=%v ret=%v",
						head.From,
						head.To,
						head.Cmd,
						head.Seq,
						ret)
				}
			}
		default:
			this.infoLog.Printf("MultiCastNotification Unhandle stat=%v", procType)
		}
	}
	return err
}

func (this *RoomManager) IsUserAdmin(roomId, uid uint32) (bResult bool, err error) {
	if roomId == 0 || uid == 0 {
		err = ErrInputParam
		return bResult, err
	}

	roomInfo, err := this.GetRoomInfo(roomId)
	if err != nil {
		this.infoLog.Printf("IsUserAdmin roomId=%v uid=%v GetRoom failed", roomId, uid)
		return bResult, err
	}

	if uid == roomInfo.CreateUid {
		bResult = true
		err = nil
		return bResult, err
	}
	adminList := roomInfo.AdminList
	for _, v := range adminList {
		if uid == v {
			bResult = true
			err = nil
			return bResult, err
		}
	}

	bResult = false
	err = nil
	return bResult, err
}

func (this *RoomManager) IsOpenVerify(roomId uint32) (bResult bool, err error) {
	if roomId == 0 {
		err = ErrInputParam
		return bResult, err
	}
	roomInfo, err := this.GetRoomInfo(roomId)
	if err != nil {
		this.infoLog.Printf("IsOpenVerify GetRoomInf failed roomid=%v err=%v", roomId, err)
		return bResult, err
	}
	if roomInfo.VerifyStat == uint32(ht_muc.VERIFY_STAT_ENUM_NEED_VERIFY) {
		bResult = true
	}
	return bResult, err
}

func (this *RoomManager) IsExceedRoomMemberLimit(roomId, addCount uint32) (bResult bool, err error) {
	if roomId == 0 {
		err = ErrInputParam
		return bResult, err
	}
	// roomInfo 为一个群信息的指针 可以通过此指针修改群信息
	roomInfo, err := this.GetRoom(roomId)
	if err != nil {
		this.infoLog.Printf("IsExceedRoomMemberLimit GetRoom roomId=%v failed err==%v", roomId, err)
		return bResult, err
	}
	// 判断成员是否超过限制
	if uint32(len(roomInfo.MemberList))+addCount > roomInfo.MemberLimit {
		// 如果当前群成员数限制不等于vip 群成员数限制需要检查用户是否是vip
		if roomInfo.MemberLimit < MUC_MEMBER_LIMIT_VIP {
			this.infoLog.Printf("IsExceedRoomMemberLimit roomId=%v memberLimit=%v not equal viplimiit=%v",
				roomId,
				roomInfo.MemberLimit,
				MUC_MEMBER_LIMIT_VIP)

			vipExpireTs, err := this.dbManager.GetUserVIPExpireTS(roomInfo.CreateUid)
			if err != nil {
				this.infoLog.Printf("IsExceedRoomMemberLimit GetUserVIPExpireTS failed room=%v addCount=%v err=%v",
					roomId,
					addCount,
					err)
				// 查询vip过期时间失败直接认为过期 用户成员数超过群限制返回错误
				err = nil
				bResult = true
				this.infoLog.Printf("IsExceedRoomMemberLimit roomId=%v addCount=%v exec limit", roomId, addCount)
				return bResult, err

			}
			tsNow := time.Now().Unix()
			if vipExpireTs > uint64(tsNow) {
				// 是vip会员 但是群信息还是非vip时创建的则更新群成员数限制到数据库同时更新内存
				roomInfo.MemberLimit = MUC_MEMBER_LIMIT_VIP
				err = this.dbManager.UpdateRoomMemberLimit(roomId, roomInfo.MemberLimit)
				if err != nil {
					this.infoLog.Printf("IsExceedRoomMemberLimit UpdateRoomMemberLimit failed room=%v memberLimit=%v err=%v",
						roomId,
						roomInfo.MemberLimit,
						err)
				}
				// 更新群成员数限制后仍然超员了直接返回错误
				if uint32(len(roomInfo.MemberList))+addCount > roomInfo.MemberLimit {
					err = nil
					bResult = true
					this.infoLog.Printf("IsExceedRoomMemberLimit roomId=%v addCount=%v exec limit", roomId, addCount)
					return bResult, err
				} else {
					err = nil
					bResult = false
					return bResult, err
				}
			} else {
				// 查询成功但是会员已过期直接返回超员了
				err = nil
				bResult = true
				this.infoLog.Printf("IsExceedRoomMemberLimit roomId=%v addCount=%v exec limit", roomId, addCount)
				return bResult, err
			}
		} else {
			//群成员超过VIP成员数限制 直接反馈true
			err = nil
			bResult = true
			this.infoLog.Printf("AddMember roomId=%v addCount=%v exec limit", roomId, addCount)
			return bResult, err
		}
	} else {
		bResult = false
	}
	return bResult, err
}

func (this *RoomManager) AddMember(roomId uint32, inviter *ht_muc.RoomMemberInfo, inviteeList []*ht_muc.RoomMemberInfo) (roomTS int64, err error) {
	// roomInfo 为一个群信息的指针 可以通过此指针修改群信息
	roomInfo, err := this.GetRoomInfo(roomId)
	if err != nil {
		this.infoLog.Printf("AddMember GetRoom roomId=%v failed err==%v", roomId, err)
		return roomTS, err
	}
	// 判断成员是否超过限制
	bResult, err := this.IsExceedRoomMemberLimit(roomId, uint32(len(inviteeList)))
	if err != nil || bResult == true {
		this.infoLog.Printf("AddMember member is exceed limit roomId=%v inviterId=%v inviteeCount=%v err=%v",
			roomId,
			inviter.Uid,
			len(inviteeList),
			err)
		err = ErrExecLimit
		return roomTS, err
	}
	// 检查用户是否已经加入群聊中
	var realAddList []*ht_muc.RoomMemberInfo
	memberList := roomInfo.MemberList
	for _, value := range inviteeList {
		var bAlreadyIn bool
		for _, v := range memberList {
			if v.Uid == value.GetUid() {
				this.infoLog.Printf("AddMember roomId=%v member=%v already in", roomId, value.GetUid())
				bAlreadyIn = true
				break
			}
		}
		// 如果用户不在群组中这添加到realAddList
		if !bAlreadyIn {
			realAddList = append(realAddList, value)
		}
	}

	if len(realAddList) == 0 {
		err = ErrAlreadyIn
		this.infoLog.Printf("AddMember roomId=%v all member already in", roomId)
		return roomTS, err
	}

	//添加到群成员列表中
	roomTS, err = this.dbManager.InviteMember(roomId, inviter.GetUid(), realAddList)
	if err != nil {
		this.infoLog.Printf("AddMember roomId=%v inviteId=%v exec db InviteMember failed", roomId, inviter.GetUid())
		return roomTS, err
	}
	// 更新内存总的群成员列表
	totalMemberList, maxOrder, err := this.dbManager.GetRoomMemberList(roomId)
	if err != nil {
		this.infoLog.Printf("AddMember exec db.GetRoomMemberList()failed roomId=%v err=%v", roomId, err)
		return roomTS, err
	}
	// update memberlist
	this.roomInfoLock.Lock()
	defer this.roomInfoLock.Unlock()
	roomInfo, ok := this.roomIdToRoomInfo[roomId]
	if ok {
		roomInfo.MaxOrder = maxOrder
		roomInfo.MemberList = totalMemberList
		roomInfo.RoomTS = roomTS
	} else {
		this.infoLog.Printf("AddMember not found roomId=%v RoomInfo", roomId)
	}

	return roomTS, err
}

func (this *RoomManager) NotifyAdminRequestJoin(roomId uint32, head *common.HeadV3, notifyBody []byte) (roomTS int64, err error) {
	if roomId == 0 {
		this.infoLog.Printf("NotifyAdminRequestJoin input param err roomId=%v", roomId)
		err = ErrInputParam
		return roomTS, err
	}
	// roomInfo 为一个群信息的指针 可以通过此指针修改群信息
	roomInfo, err := this.GetRoomInfo(roomId)
	if err != nil {
		this.infoLog.Printf("NotifyAdminRequestJoin GetRoomInfo roomId=%v failed err==%v", roomId, err)
		return roomTS, err
	}
	adminList := []uint32{roomInfo.CreateUid}
	for _, v := range roomInfo.AdminList {
		if v != 0 {
			if !this.UidIsInSlice(adminList, v) {
				adminList = append(adminList, v)
			}
		}
	}

	roomTS = roomInfo.RoomTS
	err = this.MultiCastNotification(adminList, head, notifyBody, false) // 即发送在线也存储离线
	if err != nil {
		this.infoLog.Printf("NotifyAdminRequestJoin roomId=%v err=%v failed", roomId, err)
	}
	return roomTS, err
}

func (this *RoomManager) NotifyAdminPromotJoin(opInfo *ht_muc.RoomMemberInfo,
	inviterInfo *ht_muc.RoomMemberInfo,
	memberList []*ht_muc.RoomMemberInfo,
	roomId uint32,
	roomTS uint64,
	roomIdFrom uint32,
	msgId []byte) (err error) {
	head := &common.HeadV3{Flag: uint8(common.CServToServ),
		Version:  common.CVerMmedia,
		CryKey:   uint8(common.CNoneKey),
		TermType: uint8(0),
		Cmd:      uint16(ht_muc.MUC_CMD_TYPE_GO_CMD_NOTIFY_INVITE_MEMBER),
		Seq:      this.GetPacketSeq(),
		From:     roomId,
		To:       0,
		Len:      0,
	}

	reqBody := new(ht_muc.MucReqBody)
	reqBody.NotifyInviteMemberReqbody = &ht_muc.NotifyInviteMemberReqBody{
		RoomId:        proto.Uint32(roomId),
		OpInfo:        opInfo,
		InviterInfo:   inviterInfo,
		Members:       memberList,
		NotifyTime:    proto.Uint32(uint32(time.Now().Unix())),
		RoomTimestamp: proto.Uint64(roomTS),
		RoomIdFrom:    proto.Uint32(roomIdFrom),
		MsgId:         msgId,
	}
	notifyBody, err := proto.Marshal(reqBody)

	if err != nil {
		this.infoLog.Printf("NotifyInviteMember proto marshal roomId=%v inviterId=%v err=%v",
			roomId,
			inviterInfo.GetUid(),
			err)
		return err
	}

	//调用广播接口
	//因为是管理员批准的，邀请者和管理员肯定不是同一个人所以无需提出邀请者
	return this.BroadCastNotification(roomId, 0, head, notifyBody)
}

func (this *RoomManager) UpdateVerifyStat(roomId, verifyStat uint32) (roomTS int64, err error) {
	if roomId == 0 {
		this.infoLog.Printf("UpdateVerifyStat input param err roomId=%v", roomId)
		err = ErrInputParam
		return roomTS, err
	}
	// 首先更新db
	roomTS, err = this.dbManager.UpdateVerifyStat(roomId, verifyStat)
	if err != nil {
		this.infoLog.Printf("UpdateVerifyStat dbManager.UpdateVerifyStat roomId=%v failed", roomId)
		return roomTS, err
	}

	// 更新内存
	roomInfo, err := this.GetRoomInfo(roomId)
	if err != nil {
		this.infoLog.Printf("UpdateVerifyStat GetRoomInfo roomId=%v failed", roomId)
		return roomTS, err
	}
	// 加锁
	this.roomInfoLock.Lock()
	defer this.roomInfoLock.Unlock()
	roomInfo.RoomTS = roomTS
	roomInfo.VerifyStat = verifyStat
	this.infoLog.Printf("UpdateVerifyStat set roomInfo.VerifyStat=%v", verifyStat)
	return roomTS, nil
}

func (this *RoomManager) NotifyOpenVerify(roomId, reqUid, verifyStat uint32, roomTS int64) (err error) {
	head := &common.HeadV3{Flag: uint8(common.CServToServ),
		Version:  common.CVerMmedia,
		CryKey:   uint8(common.CNoneKey),
		TermType: uint8(0),
		Cmd:      uint16(ht_muc.MUC_CMD_TYPE_GO_CMD_NOTIFY_OPEN_REQ_VERIFY),
		Seq:      this.GetPacketSeq(),
		From:     roomId,
		To:       0,
		Len:      0,
	}

	verify := ht_muc.VERIFY_STAT(verifyStat)
	reqBody := new(ht_muc.MucReqBody)
	reqBody.NotifyOpenVerifyReqbody = &ht_muc.NotifyOpenVerifyReqBody{
		RoomId:        proto.Uint32(roomId),
		OpType:        &verify,
		NotifyTime:    proto.Uint32(uint32(time.Now().Unix())),
		RoomTimestamp: proto.Uint64(uint64(roomTS)),
	}
	notifyBody, err := proto.Marshal(reqBody)

	if err != nil {
		this.infoLog.Printf("NotifyOpenVerify proto marshal roomId=%v opType=%v err=%v",
			roomId,
			verify,
			err)
		return err
	}

	//调用广播接口
	return this.BroadCastNotification(roomId, reqUid, head, notifyBody)
}

func (this *RoomManager) SetAdminList(roomId, opUid uint32, adminList []*ht_muc.RoomMemberInfo) (roomTS int64, err error) {
	if roomId == 0 || opUid == 0 {
		this.infoLog.Printf("SetAdminList input param err roomId=%v opUid=%v adminSize=%v", roomId, opUid, len(adminList))
		err = ErrInputParam
		return roomTS, err
	}

	// 首先更新db
	roomTS, err = this.dbManager.SetAdminList(roomId, adminList)
	if err != nil {
		this.infoLog.Printf("SetAdminList dbManager.SetAdminList roomId=%v failed", roomId)
		return roomTS, err
	}

	// 更新内存
	roomInfo, err := this.GetRoomInfo(roomId)
	if err != nil {
		this.infoLog.Printf("SetAdminList GetRoomInfo roomId=%v failed", roomId)
		return roomTS, err
	}
	// 加锁
	this.roomInfoLock.Lock()
	defer this.roomInfoLock.Unlock()
	roomInfo.RoomTS = roomTS

	var setAdminList []uint32
	for _, v := range adminList {
		setAdminList = append(setAdminList, v.GetUid())
	}
	roomInfo.AdminList = setAdminList // 更新新设置的管理员
	return roomTS, nil
}

func (this *RoomManager) NotifySetAdmin(roomId, opUid uint32, adminList []*ht_muc.RoomMemberInfo, roomTS int64) (err error) {
	head := &common.HeadV3{Flag: uint8(common.CServToServ),
		Version:  common.CVerMmedia,
		CryKey:   uint8(common.CNoneKey),
		TermType: uint8(0),
		Cmd:      uint16(ht_muc.MUC_CMD_TYPE_GO_CMD_NOTIFY_SET_ADMIN),
		Seq:      this.GetPacketSeq(),
		From:     roomId,
		To:       0,
		Len:      0,
	}
	reqBody := new(ht_muc.MucReqBody)
	reqBody.NotifySetAdminReqbody = &ht_muc.NotifySetAdminReqBody{
		RoomId:        proto.Uint32(roomId),
		Members:       adminList,
		NotifyTime:    proto.Uint32(uint32(time.Now().Unix())),
		RoomTimestamp: proto.Uint64(uint64(roomTS)),
	}
	notifyBody, err := proto.Marshal(reqBody)

	if err != nil {
		this.infoLog.Printf("NotifySetAdmin proto marshal roomId=%v opUid=%v adminSize=%v err=%v",
			roomId,
			opUid,
			len(adminList),
			err)
		return err
	}
	//调用广播接口
	return this.BroadCastNotification(roomId, opUid, head, notifyBody)
}

func (this *RoomManager) SetCreateUid(roomId, opUid, targetUid uint32) (roomTS int64, err error) {
	if roomId == 0 || opUid == 0 {
		this.infoLog.Printf("SetCreateUid input param err roomId=%v opUid=%v targetUid=%v", roomId, opUid, targetUid)
		err = ErrInputParam
		return roomTS, err
	}

	// 首先更新db
	roomTS, err = this.dbManager.UpdateCreateUid(roomId, opUid, targetUid)
	if err != nil {
		this.infoLog.Printf("SetCreateUid dbManager.SetCreateUid roomId=%v failed", roomId)
		return roomTS, err
	}

	// 更新内存
	roomInfo, err := this.GetRoomInfo(roomId)
	if err != nil {
		this.infoLog.Printf("SetCreateUid GetRoomInfo roomId=%v failed", roomId)
		return roomTS, err
	}
	// 加锁
	this.roomInfoLock.Lock()
	defer this.roomInfoLock.Unlock()
	roomInfo.RoomTS = roomTS
	roomInfo.CreateUid = targetUid
	return roomTS, nil
}

func (this *RoomManager) NotifyCreateUserAuthTrans(roomId, opUid uint32, memberInfo *ht_muc.RoomMemberInfo, roomTS int64) (err error) {
	head := &common.HeadV3{Flag: uint8(common.CServToServ),
		Version:  common.CVerMmedia,
		CryKey:   uint8(common.CNoneKey),
		TermType: uint8(0),
		Cmd:      uint16(ht_muc.MUC_CMD_TYPE_GO_CMD_NOTIFY_AUTHORIZATION_TRANS),
		Seq:      this.GetPacketSeq(),
		From:     roomId,
		To:       0,
		Len:      0,
	}

	reqBody := new(ht_muc.MucReqBody)
	reqBody.NotifyCreateUserTransReqbody = &ht_muc.NotifyCreateUserTransReqBody{
		RoomId:        proto.Uint32(roomId),
		Member:        memberInfo,
		NotifyTime:    proto.Uint32(uint32(time.Now().Unix())),
		RoomTimestamp: proto.Uint64(uint64(roomTS)),
	}
	notifyBody, err := proto.Marshal(reqBody)

	if err != nil {
		this.infoLog.Printf("NotifyCreateUserAuthTrans proto marshal roomId=%v opUid=%v err=%v",
			roomId,
			opUid,
			err)
		return err
	}

	//调用广播接口
	return this.BroadCastNotification(roomId, opUid, head, notifyBody)
}

func (this *RoomManager) SetAnnouncement(roomId uint32, anno *ht_muc.AnnoType) (roomTS int64, err error) {
	if roomId == 0 || anno.GetPublishUid() == 0 {
		this.infoLog.Printf("SetAnnouncement input param err roomId=%v opUid=%v", roomId, anno.GetPublishUid())
		err = ErrInputParam
		return roomTS, err
	}

	// 首先更新db
	roomTS, err = this.dbManager.UpdateAnnouncement(roomId, anno.GetPublishUid(), anno.GetPublishTs(), string(anno.GetAnnoContent()))
	if err != nil {
		this.infoLog.Printf("SetAnnouncement dbManager.SetCreateUid roomId=%v failed", roomId)
		return roomTS, err
	}

	// 更新内存
	roomInfo, err := this.GetRoomInfo(roomId)
	if err != nil {
		this.infoLog.Printf("SetAnnouncement GetRoomInfo roomId=%v failed", roomId)
		return roomTS, err
	}
	// 加锁
	this.roomInfoLock.Lock()
	defer this.roomInfoLock.Unlock()
	roomInfo.RoomTS = roomTS
	roomInfo.Announcement = AnnouncementStruct{
		PublishUid:  anno.GetPublishUid(),
		PublishTS:   anno.GetPublishTs(),
		AnnoContect: string(anno.GetAnnoContent()),
	}
	return roomTS, nil
}

func (this *RoomManager) NotifySetAnnouncement(roomId uint32, anno *ht_muc.AnnoType, roomTS int64) (err error) {
	head := &common.HeadV3{Flag: uint8(common.CServToServ),
		Version:  common.CVerMmedia,
		CryKey:   uint8(common.CNoneKey),
		TermType: uint8(0),
		Cmd:      uint16(ht_muc.MUC_CMD_TYPE_GO_CMD_NOTIFY_ROOM_ANNOUNCEMENT),
		Seq:      this.GetPacketSeq(),
		From:     roomId,
		To:       0,
		Len:      0,
	}

	reqBody := new(ht_muc.MucReqBody)
	reqBody.NotifySetRoomAnnouncementReqbody = &ht_muc.NotifySetRoomAnnouncementReqBody{
		RoomId:        proto.Uint32(roomId),
		Announcement:  anno,
		NotifyTime:    proto.Uint32(uint32(time.Now().Unix())),
		RoomTimestamp: proto.Uint64(uint64(roomTS)),
	}
	notifyBody, err := proto.Marshal(reqBody)

	if err != nil {
		this.infoLog.Printf("NotifySetAnnouncement proto marshal roomId=%v opUid=%v err=%v",
			roomId,
			anno.GetPublishUid(),
			err)
		return err
	}

	//调用广播接口
	return this.BroadCastNotification(roomId, anno.GetPublishUid(), head, notifyBody)
}
