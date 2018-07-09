package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/gansidui/gotcp"
	"github.com/gansidui/gotcp/libcomm"
	"github.com/gansidui/gotcp/libcrypto"
	"github.com/gansidui/gotcp/tcpfw/common"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jessevdk/go-flags"
	nsq "github.com/nsqio/go-nsq"
	"gopkg.in/ini.v1"
)

type Callback struct{}

var (
	mcApi          *tcpfw.MemcacheApi
	db             *sql.DB
	globalProducer *nsq.Producer
	nsqAPNSTopic   string
	nsqGCMTopic    string
	nsqXinGeTopic  string
	globalRemark   map[string]string
	remarkLock     sync.Mutex
)

// Error type
var (
	ErrNilMcObject      = errors.New("nil memcache object")
	ErrNilDbObject      = errors.New("not set object current is nil")
	ErrParam            = errors.New("err input param")
	ErrNotOpenPush      = errors.New("err not open push")
	ErrNULLPushToken    = errors.New("err push token is empty")
	ErrInvalidPushToken = errors.New("err invalid push token")
	ErrNotFoundRemark   = errors.New("err not found remark")
)

const (
	OFFICE_BEGIN_UID      = 10000
	APNS_TOKEN_LEN        = 64
	ONCE_GETRELATIONLINES = 1000000
)
const (
	CTP2P = 0
	CTMUC = 1
)

const (
	CMD_S2S_MARKNAMECHANGED_NOTIFY     = 0x8025
	CMD_S2S_MARKNAMECHANGED_NOTIFY_ACK = 0x8026
	CMD_S2S_MESSAGE_PUSH               = 0x8027
	CMD_S2S_MESSAGE_PUSH_ACK           = 0x8028
)

const (
	RET_SUCCESS            = 0
	ERR_SYSERR_START       = 100
	ERR_SERVER_BUSY        = 100
	ERR_INTERNAL_ERROR     = 101
	ERR_UNFORMATTED_PACKET = 102
	ERR_NO_ACCESS          = 103
	ERR_INVALID_CLIENT     = 104
	ERR_INVALID_SESSION    = 105
	ERR_INVALID_PARAM      = 106
)

const (
	SERVER_COMM_KEY = "lp$5F@nfN0Oh8I*5"
)

const (
	APNS_PUSH_SOUND_CALLING  = "voipcall.caf"
	APNS_PUSH_SOUND_DEFAULT  = "default"
	PUSH_PARAM_CALLINCOMING  = "voip_call_incoming"
	PUSH_PARAM_CALLIMISS     = "voip_call_miss"
	PUSH_PARAM_CALLICANCEL   = "voip_call_cancel"
	PUSH_PARAM_LX_ACCEPT     = "language_exchange_accept"
	PUSH_PARAM_LX_DECLINED   = "language_exchange_declined"
	PUSH_PARAM_LX_TERMINATED = "language_exchange_terminated"
	PUSH_PARAM_LX_CANCLED    = "language_exchange_cancled"
)

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
	PT_TEXT              = 0
	PT_VOICE             = 1
	PT_PHOTO             = 2
	PT_INDRODUCE         = 3
	PT_LOCATION          = 4
	PT_FRIEND_INVITE     = 5
	PT_LANGUAGE_EXCHANGE = 6
	PT_CORRECT_SENTENCE  = 7
	PT_STICKERS          = 8
	PT_DOODLE            = 9
	PT_GIFT              = 10
	PT_VOIP              = 11
	PT_INVITE_ACCEPT     = 12
	PT_VIDEO             = 13

	PT_GVOIP               = 15
	PT_LINK                = 16
	PT_CARD                = 17
	PT_FOLLOW              = 18
	PT_REPLY_YOUR_COMMENT  = 19
	PT_COMMENTED_YOUR_POST = 20
	PT_CORRECTED_YOUR_POST = 21
	PT_MOMENT_LIKE         = 22
)

type PushMessage struct {
	terminalType uint8
	chatType     uint8
	fromId       uint32
	toId         uint32
	roomId       uint32
	pushType     uint8
	nickName     []byte
	pushContent  []byte
	sound        uint8
	lights       uint8
	messageId    []byte
	actionId     uint32
	byAt         uint8
}

type SimplePushSetting struct {
	token   []byte
	preview uint8
	dndTime bool
	badge   uint16
	isXinGe bool
}

type StoreUserPushSetting struct {
	uid       uint32
	pushToken []byte
	dndSet    uint8
	dndStart  uint8
	dndEnd    uint8
	timeZone  uint8
	alert     uint8
	preview   uint8
	pushType  uint8
	follow    uint8
	moment    uint8
}

func (this *Callback) OnConnect(c *gotcp.Conn) bool {
	addr := c.GetRawConn().RemoteAddr()
	c.PutExtraData(addr)
	log.Println("OnConnect:", addr)
	return true
}

func (this *Callback) OnMessage(c *gotcp.Conn, p gotcp.Packet) bool {
	packet, ok := p.(*common.XTHeadPacket)
	if !ok { // 不是XTHeadPacket报文
		log.Println("OnMessage packet can not change to xtpacket")
		return false
	}

	head, err := packet.GetHead()
	if err != nil {
		log.Println("OnMessage Get head failed", err)
		return false
	}
	attr := "gopushpro/recv_req_count"
	libcomm.AttrAdd(attr, 1)

	//log.Printf("OnMessage:[%#v] len=%v payLoad=%v\n", head, len(packet.GetBody()), packet.GetBody())
	log.Printf("OnMessage:[%#v] len=%v\n", head, len(packet.GetBody()))
	_, err = packet.CheckXTPacketValid()
	if err != nil {
		SendResp(c, head, uint8(ERR_INVALID_PARAM))
		log.Println("Invalid packet", err)
		return false
	}

	switch head.Cmd {
	case CMD_S2S_MESSAGE_PUSH:
		go ProcMessagePush(c, head, p)
	default:
		log.Println("OnMessage UnHandle Cmd =", head.Cmd)
	}
	return true
}

func SendResp(c *gotcp.Conn, reqHead *common.XTHead, ret uint8) bool {
	if c == nil || reqHead == nil {
		log.Printf("SendResp nil conn=%v reqhead=%v", c, reqHead)
		return false
	}

	head := new(common.XTHead)
	if reqHead != nil {
		*head = *reqHead
	}

	head.Cmd = reqHead.Cmd + 1 // ack cmd = req cmd + 1
	head.Len = 1               // sizeof(uint8)
	buf := make([]byte, common.XTHeadLen+head.Len)
	err := common.SerialXTHeadToSlice(head, buf[:])
	if err != nil {
		log.Printf("SendResp SerialXTHeadToSlice failed")
		return false
	}
	buf[common.XTHeadLen] = ret // return code
	resp := common.NewXTHeadPacket(buf)
	c.AsyncWritePacket(resp, time.Second)
	return true
}

func SendRespWithPayLoad(c *gotcp.Conn, reqHead *common.XTHead, payLoad []byte) bool {
	if c == nil || reqHead == nil {
		log.Printf("SendRespWithPayLoad nil conn=%v reqhead=%v", c, reqHead)
		return false
	}

	head := new(common.XTHead)
	if reqHead != nil {
		*head = *reqHead
	}

	head.Cmd = reqHead.Cmd + 1      // ack cmd = req cmd + 1
	head.Len = uint32(len(payLoad)) //
	buf := make([]byte, common.XTHeadLen+head.Len)
	err := common.SerialXTHeadToSlice(head, buf[:])
	if err != nil {
		log.Println("SerialXTHeadToSlice failed")
		return false
	}
	copy(buf[common.XTHeadLen:], payLoad) // return code
	resp := common.NewXTHeadPacket(buf)
	c.AsyncWritePacket(resp, time.Second)
	return true
}
func ProcMessagePush(c *gotcp.Conn, reqHead *common.XTHead, packet *common.XTHeadPacket) bool {
	if c == nil || reqHead == nil || p == nil {
		log.Printf("ProcMessagePush invalid param c=%v reqHead=%v packet=%v", c, reqHead, p)
		return false
	}
	var body []byte
	if head.CryKey == common.CServKey {
		// 使用Server key 加密
		plainText := libcrypto.TEADecrypt(string(packet.GetBody()), SERVER_COMM_KEY)
		log.Printf("ProcMessagePush plainText len=%v", len(plainText))
		body = []byte(plainText)
	} else if head.CryKey == common.CNoneKey {
		body = packet.GetBody()
	} else {
		log.Printf("ProcMessagePush UnHandle CryKey=%v", head.CryKey)
		return false
	}

	message := PushMessage{
		termType:    common.UnMarshalUint8(&body),
		chatType:    common.UnMarshalUint8(&body),
		fromId:      common.UnMarshalUint32(&body),
		toId:        common.UnMarshalUint32(&body),
		roomId:      common.UnMarshalUint32(&body),
		pushType:    common.UnMarshalUint8(&body),
		nickName:    common.UnMarshalSlice(&body),
		pushContent: common.UnMarshalSlice(&body),
		sound:       common.UnMarshalUint8(&body),
		lights:      common.UnMarshalUint8(&body),
		messageId:   common.UnMarshalSlice(&body),
		actionId:    common.UnMarshalUint32(&body),
		byAt:        common.UnMarshalUint8(&body),
	}

	log.Printf("ProcMessagePush PushMessage=%#v", message)
	pushSetting, err := GetPushSetting(message.toId, message.terminalType, message.roomId)
	if err != nil {
		log.Printf("ProcMessagePush get push setting failed err=%v message=%#v", err, message)
		return false
	}
	if message.terminalType == common.CClientTyepIOS {

	} else {
		if pushSetting.isXinGe {

		} else {

		}
	}
	return true
}

func GetPushSetting(uid uint32, terminal, pushType uint8) (pushSetting *SimplePushSetting, err error) {
	if mcApi == nil {
		log.Printf("GetPushSetting nil mcApi")
		return nil, ErrNilMCObject
	}
	pushSetting = new(SimplePushSetting)
	// push badge
	if terminal == uint8(common.CClientTyepIOS) {
		badge, err := mcApi.GetPushBadge(uid)
		if err != nil {
			log.Printf("GetPushSetting mcApi.GetPushBadge failed uid=%v err=%v", uid, err)
			badge = 0 // get mc failed set badge zero
		}
		pushSetting.badge = badge
	}
	pushConfig, err := mcApi.GetPushSetting(uid)
	if err != nil {
		log.Printf("GetPushSetting McApi.GetPushSetting load from db failed uid=%v err=%v", uid, err)
		storePushConfig, err := ReloadUserPushSettingByUid(uid)
		if err != nil {
			log.Printf("GetPushSetting ReloadUserPushSettingByUid faild uid=%v err=%v", uid, err)
			return nil, err
		}
		pushConfig.Uid = storePushConfig.uid
		pushConfig.PushToken = storePushConfig.pushToken
		pushConfig.DndSet = storePushConfig.dndSet
		pushConfig.DndStart = storePushConfig.dndStart
		pushConfig.DndEnd = storePushConfig.dndEnd
		pushConfig.TimeZone = storePushConfig.timeZone
		pushConfig.Alert = storePushConfig.alert
		pushConfig.Preview = storePushConfig.preview
		pushConfig.PushType = storePushConfig.pushType
		pushConfig.Follow = storePushConfig.follow
		pushConfig.Moment = storePushConfig.moment

		err := mcApi.SetPushSetting(pushConfig)
		if err != nil {
			log.Printf("GetPushSetting mcApi.SetPushSetting()exec failed uid=%v err=%v", pushConfig.Uid, err)
		} else {
			log.Printf("GetPushSetting mcApi.SetPushSetting success pushConfig=%#v", pushConfig)
		}
	} else {
		log.Printf("GetPushSetting mcApi.GetPushSetting success pushConfig=%#v", pushConfig)
	}
	// whether open push
	if pushConfig.Alert == 0 {
		log.Printf("GetPushSetting close push pushConfig=%#v", pushConfig)
		return nil, ErrNotOpenPush
	} else if pushConfig.Follow == 0 && pushType == PUSH_FOLLOW { // 1:on 0:off
		log.Printf("GetPushSetting close follow notification pushConfig=%#v", pushConfig)
		return nil, ErrNotOpenPush
	} else if pushConfig.Moment == 0 && (pushType == PUSH_REPLY_YOUR_COMMENT ||
		pushType == PUSH_COMMENTED_YOUR_POST ||
		pushType == PUSH_CORRECTED_YOUR_POST) { //1:on 0:off
		log.Printf("GetPushSetting close moment notification pushConfig=%#v", pushConfig)
		return ErrNotOpenPush
	}
	pushSetting.preview = pushConfig.Preview
	if len(pushConfig.PushToken) == 0 {
		log.Printf("GetPushSetting push token is empty uid=%v terminal=%v pushConfig=%#v", uid, terminal, pushConfig)
		return nil, ErrNULLPushToken
	}
	// token 合法性检验
	if terminal == common.CClientTyepIOS {
		if !CheckAPNSToken(pushConfig.PushToken) {
			log.Printf("GetPushSetting invalid apns token uid=%v token=%s", uid, pushConfig.PushToken)
			return nil, ErrInvalidPushToken
		}
		pushSetting.pushToken = pushConfig.PushToken
	} else {
		if len(pushConfig.PushToken) <= 0 {
			log.Printf("GetPushSetting invalid android token uid=%v pushConfig=%#v", uid, pushConfig)
			return nil, ErrInvalidPushToken
		}
		pushSetting.pushToken = pushConfig.PushToken
		pushSetting.isXinGe = false
		// pushtype 为1:信鸽 0:GCM
		pushSetting.isXinGe = (pushConfig.PushType == 1)
	}
	// 是否在免打扰时间内
	if pushConfig.DndSet == 1 {
		localHour := GetLocalHour(pushConfig.TimeZone)
		pushSetting.dndTime = BetweenDNDTime(localHour, pushConfig.DndStart, pushConfig.DndEnd)
	} else {
		pushSetting.dndTime = false
	}
	return pushSetting, nil
}

func CheckAPNSToken(pushToken []byte) (result bool) {
	tokenLen := len(pushToken)
	if tokenLen != APNS_TOKEN_LEN {
		return false
	}
	for _, v := range pushToken {
		if !((v >= '0' && v <= '9') || (v >= 'a' && v <= 'f')) {
			return false
		}
	}
	return true
}

func GetLocalHour(timeZone uint8) (localHour uint16) {
	// timeZone 是可能小余0 的需要将无符号转成有符号数字
	realyZone := int8(timeZone)
	// zone 里面记录的是小时需要将小时转换成秒
	localZone := time.FixedZone("UTC", realyZone*3600)
	timeNow := time.Now()
	timeLocal := timeNow.In(localZone)
	localHour = uint16(timeLocal.Hour())
	return localHour
}

func BetweenDNDTime(localHour, startHour, endHour uint16) (ret bool) {
	if startHour <= endHour {
		ret = localHour >= startHour && localHour <= endHour
	} else {
		ret = !(localHour >= endHour && localHour <= startHour)
	}
	return ret
}

func ReloadUserPushSettingByUid(uid uint32) (pushSetting StoreUserPushSetting, err error) {
	if db == nil || uid < OFFICE_BEGIN_UID {
		return nil, ErrDbParam
	}
	var storeUid sql.NullInt64
	var storeToken sql.NullString
	var storeTimeZone, storeDndSet, storeDndStart, storeDndEnd, storeAlert, storePreview, storePushType, storeFollow, storeMoment NullInt64
	err = db.QueryRow("SELECT t1.USERID, t1.APNSTOKEN, t2.TIMEZONE, t3.DNDSET, t3.DNDSTART, t3.DNDEND, t3.NOTIFYALERT, t3.NOTIFYPREVIEW, t1.PUSHTYPE, t3.NOTIFY_FOLLOW, t3.NOTIFY_MOMENT FROM HT_USER_TERMINAL AS t1 LEFT JOIN (HT_USER_BASE AS t2, HT_USER_SETTING AS t3) ON (t2.USERID = t1.USERID and t3.USERID = t1.USERID) WHERE t1.USERID=?", uid).Scan(&storeUid, &storeToken, &storeTimeZone, &storeDndSet, &storeDndStart, &storeDndEnd, &storeAlert, &storePreview, &storePushType, &storeFollow, &storeMoment)
	switch {
	case err == sql.ErrNoRows:
		log.Printf("ReloadUserPushSettingByUid not found uid=%v", uid)
		return nil, err
	case err != nil:
		log.Printf("ReloadUserPushSettingByUid exec failed uid=%v err=%v", uid, err)
		return nil, err
	default:
	}

	if storeUid.Valid {
		pushSetting.Uid = uint32(storeUid.Int64)
	}
	if storeToken.Valid {
		pushSetting.pushToken = []byte(storeToken.String)
	}
	if storeTimeZone.Valid {
		pushSetting.timeZone = uint8(storeTimeZone.Int64)
	}
	if storeDndSet.Valid {
		pushSetting.dndSet = uint8(storeDndSet)
	}
	if storeDndStart.Valid {
		pushSetting.dndStart = uint8(storeDndStart)
	}
	if storeDndEnd.Valid {
		pushSetting.dndEnd = uint8(storeDndEnd)
	}
	if storeAlert.Valid {
		pushSetting.alert = uint8(storeAlert)
	}
	if storePreview.Valid {
		pushSetting.preview = uint8(storePreview)
	}
	if storePushType.Valid {
		pushSetting.pushType = uint8(storePushType)
	}
	if storeFollow.Valid {
		pushSetting.follow = uint8(storeFollow)
	}
	if storeMoment.Valid {
		pushSetting.moment = uint8(storeMoment)
	}
	log.Printf("ReloadUserPushSettingByUid pushSetting=%#v", pushSetting)

	return pushSetting, nil
}

func PushAPNSMessage(message *PushMessage, pushSetting *SimplePushSetting) (err error) {
	if message == nil || pushSetting == nil {
		log.Printf("PushAPNSMessage input param err message=%v, pushSetting=%v", messag, pushSetting)
		return ErrParam
	}
	var soundStr string
	if !pushSetting.dndTime && messag.sound {
		if message.pushType == uint8(PUSH_VOIP) && string(message.pushContent) == PUSH_PARAM_CALLINCOMING {
			soundStr = PUSH_SOUND_CALLING
		} else {
			soundStr = PUSH_SOUND_DEFAULT
		}
	}
	var senderName string = string(message.nickName)
	remarkName, err := FindRemarkName(message.fromId, message.toId)
	if err == nil && len(remarkName) > 0 {
		senderName = remarkName
	}
	var fromId uint32
	if message.chatType == uint8(CTMUC) {
		fromId = message.roomId
	} else {
		fromid = message.fromId
	}

	pushType := PT_TEXT
	pushParam := message.pushContent
	switch int(message.pushType) {
	case PUSH_TEXT:
		pushType = PT_TEXT
	case PUSH_VOICE:
		pushType = PT_VOICE
	case PUSH_IMAGE:
		pushType = PT_PHOTO
	case PUSH_INTRODUCE:
		pushType = PT_INDRODUCE
	case PUSH_LOCATION:
		pushType = PT_LOCATION
	case PUSH_FRIEND_INVITE:
		pushType = PT_FRIEND_INVITE
	case PUSH_LANGUAGE_EXCHANGE, PUSH_LANGUAGE_EXCHANGE_REPLY:
		pushType = PT_LANGUAGE_EXCHANGE

	}

}

func PushXinGeMessage(message *PushMessage, pushSetting *SimplePushSetting) (err error) {

}

func PushGCMMessage(message *PushMessage, pushSetting *SimplePushSetting) (err error) {

}

func (this *Callback) OnClose(c *gotcp.Conn) {
	log.Println("OnClose:", c.GetExtraData())
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

	// log.Println("config name =", options.ServerConf)
	// 读取配置文件
	cfg, err := ini.Load([]byte(""), options.ServerConf)
	if err != nil {
		log.Printf("load config file=%s failed", options.ServerConf)
		return
	}
	// 配置文件只读 设置此标识提升性能
	cfg.BlockMode = false

	// 读取 Memcache Ip and port
	mcIp := cfg.Section("MEMCACHE").Key("mc_ip").MustString("127.0.0.1")
	mcPort := cfg.Section("MEMCACHE").Key("mc_port").MustInt(11211)
	log.Printf("memcache ip=%v port=%v", mcIp, mcPort)
	mcApi = new(tcpfw.MemcacheApi)
	mcApi.Init(mcIp + ":" + strconv.Itoa(mcPort))

	// nsq producer config
	nsqUrl := cfg.Section("NSQ").Key("url").MustString("127.0.0.1:4150")
	config := nsq.NewConfig()
	globalProducer, err := nsq.NewProducer(nsqUrl, config)
	if err != nil {
		log.Fatalf("main nsq.NewProcucer failed nsqUrl=%s", nsqUrl)
	}
	nsqAPNSTopic = cfg.Section("APNS").Key("topic").Mustring("apns")
	nsqGCMTopic = cfg.Section("GCM").Key("topic").Mustring("gcm")
	nsqXinGeTopic = cfg.Section("XinGe").Key("topic").Mustring("xinge")

	// init mysql
	mysqlHost := cfg.Section("MYSQL").Key("mysql_host").MustString("127.0.0.1")
	mysqlUser := cfg.Section("MYSQL").Key("mysql_user").MustString("IMServer")
	mysqlPasswd := cfg.Section("MYSQL").Key("mysql_passwd").MustString("hello")
	mysqlDbName := cfg.Section("MYSQL").Key("mysql_db").MustString("HT_IMDB")
	mysqlPort := cfg.Section("MYSQL").Key("mysql_port").MustString("3306")

	log.Printf("mysql host=%v user=%v passwd=%v dbname=%v port=%v",
		mysqlHost,
		mysqlUser,
		mysqlPasswd,
		mysqlDbName,
		mysqlPort)

	db, err = sql.Open("mysql", mysqlUser+":"+mysqlPasswd+"@"+"tcp("+mysqlHost+":"+mysqlPort+")/"+mysqlDbName+"?charset=utf8&timeout=90s")
	if err != nil {
		infoLog.Println("open mysql failed")
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}

	// load all remark from db
	mapFriendRemark := map[string]string{}
	err := LoadRemarNameFromDB()
	if err != nil {
		log.Printf("main LoadRemarNameFromDB failed err=%s", err)
	}

	// creates a tcp listener
	serverIp := cfg.Section("LOCAL_SERVER").Key("bind_ip").MustString("127.0.0.1")
	serverPort := cfg.Section("LOCAL_SERVER").Key("bind_port").MustInt(8990)
	log.Printf("serverIp=%v serverPort=%v", serverIp, serverPort)
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
	srv := gotcp.NewServer(config, &Callback{}, &common.XTHeadProtocol{})

	// starts service
	go srv.Start(listener, time.Second)
	log.Println("listening:", listener.Addr())

	// catchs system signal
	chSig := make(chan os.Signal)
	signal.Notify(chSig, syscall.SIGINT, syscall.SIGTERM)
	log.Println("Signal: ", <-chSig)

	// stops service
	srv.Stop()
	db.Close()
}

func LoadRemarNameFromDB() (err error) {
	if db == nil {
		err = ErrNilDbObject
		log.Printf("LoadRemarNameFromDB nil db object")
		return err
	}
	remarkLock.Lock()
	defer remarkLock.UnLock()
	startLine := 0
	for {
		getLines, err := GetAllFriendRemark(startLine, ONCE_GETRELATIONLINES, &globalRemark)
		if err != nil {
			log.Printf("LoadRemarNameFromDB exec GetAllFriendRemark failed err=%s", err)
			return err
		}
		startLine += getLines
		log.Printf("LoadRemarNameFromDB get lines=%v", getLines)
		if getLines == 0 {
			log.Printf("LoadRemarNameFromDB get 0 lines return")
			break
		}
	}
	return nil
}

func GetAllFriendRemark(startLine, onceGetLines uint32, mapFriendRemark *map[string]string) (updateLines uint32, err error) {
	if db == nil {
		log.Printf("GetAllFriendRemark nil db object")
		return 0, ErrNilDbObject
	}
	var emptyRemark string
	rows, err := db.Query("select USERID,FRIENDID,FRIENDREMARKNAME from HT_FRIEND_RELATION where FRIENDREMARKNAME is not null and FRIENDREMARKNAME <>? limit ?, ?;", emptyRemark, startLine, onceGetLines)
	if err != nil {
		return nil, err
	}
	UpdateLines = 0
	defer rows.Close()
	for rows.Next() {
		var from, to uint32
		var remarkName string
		if err := rows.Scan(&from, &to, &remarkName); err != nil {
			log.Printf("GetAllFriendRemark rows.Scan failed")
			continue
		}
		if len(remarkName) <= 0 {
			log.Printf("GetAllFriendRemark empty remark name from=%v to=%v", from, to)
			continue
		}

		log.Printf("GetAllFriendRemark DEBUG from=%v to=%v remark=%s", from, to, remarkName)
		key := GetRemarkNameKey(from, to)
		mapFriendRemark[key] = remarkName
		updateLines++
	}
	return updateLines, nil
}

// mapRemarkName 中存储的是 from 给to 备注的备注名
// 在to 给 from 发送推送时 需要将to的nickname 修改成备注名
func GetRemarkNameKey(from, to uint32) (key string) {
	key = fmt.Sprintf("%v-%v", from, to)
	return key
}

func FindRemarkName(from, to uint32) (remarkName string, err error) {
	// 需要颠倒from, to 的传参顺序因为是发送消息到to 所以需要看to 给from的备注名
	key := GetRemarkNameKey(to, from)
	remarkLock.Lock()
	defer remarkLock.UnLock()
	if v, ok := globalRemark[key]; ok {
		remarkName = v
		log.Printf("FindRemarkName to=%v from=%v remarkName=%s", from, to, remarkName)
	} else {
		err = ErrNotFoundRemark
		log.Printf("FindRemarkName to=%v from=%v not found", from, to)
	}
	return remarkName, err
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
