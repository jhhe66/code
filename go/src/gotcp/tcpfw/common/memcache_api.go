package common

import (
	"encoding/binary"
	"errors"
	"strconv"

	"github.com/bradfitz/gomemcache/memcache"
)

var (
	ErrPacketShortLen = errors.New("use byte[] length is not enough")
)

const (
	COnlineStatLen        = 70 //70 = 67+2个字节长度 + 1个字节的'\0'
	CMCOnlinePrefix       = "state#"
	CMCMsgCountPrefix     = "msgcount#"       // "msgcount#10888"
	CMCPushTokenPrefix    = "token#"          // "token#10888"
	CMCPushBadgePrefix    = "badge#"          // "badge#10888"
	CMCDndSetPrefix       = "dndset#"         // "dndset#10888"
	CMCNotifySetPrefix    = "notify#"         // "notify#101305"
	CMCOfflineMsgPrefix   = "offlinemessage#" // "offlinemessage#101305"
	CMCLimitInvitePrefix  = "invite#"         // "invite#101305"
	CMCRequestLimitPrefix = "request#"        // "request#101305"
	CMCReportLimitPrefix  = "report#"         // "report#101305"
	CMCPushSettingPrefix  = "push#"           // "push#101305"
	CMCMucConfigPrefix    = "mucconfig#"      // "mucconfig#101305"
)

const (
	ST_OFFLINE   = 0
	ST_ONLINE    = 1
	ST_BACKGROUD = 2
	ST_LOGOUT    = 3
	ST_INIT      = 100 // 初始化
)

const (
	CClientTyepIOS     = 0
	CClientTypeAndroid = 1
)

const (
	CVersion225 = 0x20205
	CVerSion226 = 0x20206
)

type UserState struct {
	Uid        uint32 //用户id
	ClientType uint8  //客户端类型 TERMIAL_TYPE
	OnlineStat uint8  //在线状态
	SvrIp      uint32 //所在服务器IP
	UpdateTs   uint64 //更新时间戳
	Session    []byte //Session
	Version    uint32 //版本
	Port       uint32 //端口
	Wid        uint64 //wid
	UserType   uint8  //用户类型
}

type UserPushSetting struct {
	Uid       uint32
	PushToken []byte
	DndSet    uint8
	DndStart  uint8
	DndEnd    uint8
	TimeZone  uint8
	Alert     uint8
	Preview   uint8
	PushType  uint8
	Follow    uint8
	Moment    uint8
}

func (this *UserState) unMarshall(buf []byte) (err error) {
	if len(buf) < COnlineStatLen {
		err = ErrPacketShortLen
		return
	}
	this.Uid = UnMarshalUint32(&buf)
	this.ClientType = UnMarshalUint8(&buf)
	this.OnlineStat = UnMarshalUint8(&buf)
	this.SvrIp = UnMarshalUint32(&buf)
	this.UpdateTs = UnMarshalUint64(&buf)
	this.Session = UnMarshalSlice(&buf)
	this.Version = UnMarshalUint32(&buf)
	this.Port = UnMarshalUint32(&buf)
	this.Wid = UnMarshalUint64(&buf)
	this.UserType = UnMarshalUint8(&buf)
	return nil
}

func (this *UserState) marshall(buf []byte) (err error) {
	if len(buf) < COnlineStatLen {
		err = ErrPacketShortLen
		return
	}
	MarshalUint32(this.Uid, &buf)
	MarshalUint8(this.ClientType, &buf)
	MarshalUint8(this.OnlineStat, &buf)
	MarshalUint32(this.SvrIp, &buf)
	MarshalUint64(this.UpdateTs, &buf)
	MarshalSlice(this.Session, &buf)
	MarshalUint32(this.Version, &buf)
	MarshalUint32(this.Port, &buf)
	MarshalUint64(this.Wid, &buf)
	MarshalUint8(this.UserType, &buf)
	return nil

}

// Conn exposes a set of callbacks for the various events that occur on a connection
type MemcacheApi struct {
	client *memcache.Client
}

func (c *MemcacheApi) Init(server ...string) {
	c.client = memcache.New(server...)
}

func (c *MemcacheApi) GetUserOnlineStat(uid uint32) (stat *UserState, err error) {
	key := CMCOnlinePrefix + strconv.Itoa(int(uid))
	it, err := c.client.Get(key)
	if err != nil {
		return nil, err
	}
	stat = new(UserState)
	err = stat.unMarshall(it.Value)
	return stat, err
}

func (c *MemcacheApi) GetPushBadge(uid uint32) (badge uint16, err error) {
	key := CMCPushBadgePrefix + strconv.Itoa(int(uid))
	it, err := c.client.Get(key)
	if err != nil {
		return 0, err
	}
	payLoad := it.Value
	_ := UnMarshalUint32(&payLoad)
	badge = UnMarshalUint16(&payLoad)
	return badge, nil
}

func (c *MemcacheApi) GetPushSetting(uid uint32) (pushSetting PushSetting, err error) {
	key := CMCPushSettingPrefix + strconv.Itoa(int(uid))
	it, err := c.client.Get(key)
	if err != nil {
		return nil, err
	}
	payLoad := it.Value
	pushSetting.Uid = UnMarshalUint32(&payLoad)
	pushSetting.PushToken = UnMarshalSlice(&payLoad)
	pushSetting.DndSet = UnMarshalUint8(&payLoad)
	pushSetting.DndStart = UnMarshalUint8(&payLoad)
	pushSetting.DndEnd = UnMarshalUint8(&payLoad)
	pushSetting.TimeZone = UnMarshalUint8(&payLoad)
	pushSetting.Alert = UnMarshalUint8(&payLoad)
	pushSetting.Preview = UnMarshalUint8(&payLoad)
	pushSetting.PushType = UnMarshalUint8(&payLoad)
	if len(payLoad) >= binary.Size(uint8) {
		pushSetting.Follow = UnMarshalUint8(&payLoad)
	} else {
		pushSetting.Follow = 1 // open follow notification
	}

	if len(payLoad) >= binary.Size(uint8) {
		pushSetting.Moment = UnMarshalUint8(&payLoad)
	} else {
		pushSetting.Moment = 1 // has commnet notification
	}

	return pushSetting, nil
}

func (c *MemcacheApi) SetPushSetting(pushSetting PushSetting) (err error) {
	var packetPayLoad []byte
	MarshalUint32(pushSetting.Uid, &packetPayLoad)
	MarshalSlice(pushSetting.PushToken, &packetPayLoad)
	MarshalUint8(pushSetting.DndSet, &packetPayLoad)
	MarshalUint8(pushSetting.DndStart, &packetPayLoad)
	MarshalUint8(pushSetting.DndEnd, &packetPayLoad)
	MarshalUint8(pushSetting.TimeZone, &packetPayLoad)
	MarshalUint8(pushSetting.Alert, &packetPayLoad)
	MarshalUint8(pushSetting.Preview, &packetPayLoad)
	MarshalUint8(pushSetting.PushType, &packetPayLoad)
	MarshalUint8(pushSetting.Follow, &packetPayLoad)
	MarshalUint8(pushSetting.Moment, &packetPayLoad)

	iterm := &memcache.Item{
		Key:   CMCPushSettingPrefix + strconv.Itoa(int(uid)),
		Value: packetPayLoad,
	}
	err := c.Client.Set(iterm)
	return err
}
