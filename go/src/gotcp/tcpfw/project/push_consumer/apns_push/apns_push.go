package main

import (
	"crypto/tls"
	"log"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"

	simplejson "github.com/bitly/go-simplejson"
	"github.com/jessevdk/go-flags"
	"github.com/nsqio/go-nsq"
	"github.com/sideshow/apns2"
	"github.com/sideshow/apns2/certificate"
	"gopkg.in/ini.v1"
)

var wg *sync.WaitGroup
var clientManager *apns2.ClientManager
var myCert tls.Certificate

type Options struct {
	// Example of verbosity with level
	Verbose []bool `short:"v" long:"verbose" description:"Verbose output"`

	// Example of optional value
	ServerConf string `short:"c" long:"conf" description:"Server Config" optional:"no"`
}

const (
	//CerPassWd = "yike@helloTalk"
	CerPassWd = "123456"
	TopicId   = "com.helloTalk.helloTalk"
)

const (
	PUSH_PARAM_CALLINCOMING  = "voip_call_incoming"
	PUSH_PARAM_CALLIMISS     = "voip_call_miss"
	PUSH_PARAM_CALLICANCEL   = "voip_call_cancel"
	PUSH_PARAM_LX_ACCEPT     = "language_exchange_accept"
	PUSH_PARAM_LX_DECLINED   = "language_exchange_declined"
	PUSH_PARAM_LX_TERMINATED = "language_exchange_terminated"
	PUSH_PARAM_LX_CANCLED    = "language_exchange_cancled"
)

const (
	APNS_PUSH_LANGUAGE_EXCHANGE = "Language Exchange Request"
	APNS_PUSH_LX_ACCEPT         = "Accepted Language Exchange"
	APNS_PUSH_LX_DECLIND        = "Declined Language Exchange"
	APNS_PUSH_LX_TERMINATED     = "Terminated Language Exchange"
	APNS_PUSH_LX_CANCLED        = "Cancled Language Exchange"

	APNS_PUSH_LANGUAGE_EXCHANGE_V2 = "Language Exchange Request push"
	APNS_PUSH_LX_ACCEPT_V2         = "%@: Exchange request agreed push"
	APNS_PUSH_LX_DECLIND_V2        = "Request Declined push"
	APNS_PUSH_LX_TERMINATED_V2     = "Terminate Exchange push"
	APNS_PUSH_LX_CANCLED_V2        = "Exchange Request Canceled push"

	APNS_PUSH_CALLINCOMING = "Incoming Call"
	APNS_PUSH_CALLMISS     = "Call Canceled"
	APNS_PUSH_CALLCANCEL   = "Call Missed"
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

const (
	NOT_BEEN_AT = 0 //default not at
	BEEN_AT     = 1 //
)

const (
	CTP2P = 0
	CTMUC = 1
)

var options Options
var parser = flags.NewParser(&options, flags.Default)

func MessageHandle(message *nsq.Message) error {
	log.Printf("MessageHandle Got a message: %s", message.Body)
	wg.Add(1)
	go func() {
		PostData(message)
		wg.Done()
	}()

	return nil
}

func PostData(message *nsq.Message) error {
	log.Printf("PostData enter msgId=%v timestamp=%v body=%s", message.ID, message.Timestamp, message.Body)
	// 然后进行下一步无处理
	rootObj, err := simplejson.NewJson(message.Body)
	if err != nil {
		log.Printf("PostData simplejson new packet error", err)
		return err
	}
	log.Printf("PostData rootObj=%#v", rootObj)
	token := rootObj.Get("token").MustString("")
	pushType := rootObj.Get("push_type").MustInt(PT_TEXT)
	preview := rootObj.Get("preview").MustInt(0)
	pushParam := rootObj.Get("push_param").MustString("")
	nickName := rootObj.Get("nick_name").MustString("")
	byAt := rootObj.Get("by_at").MustInt(NOT_BEEN_AT)

	alertObj := simplejson.New()
	var alertString string
	if preview == 0 {
		alertObj.Set("loc-key", "Message Preview Example No")
	} else {
		switch pushType {
		case PT_TEXT:
			if byAt == BEEN_AT {
				alertObj.Set("loc-key", "Someone @ me")
				alertObj.Set("loc-args", []interface{}{})
			} else {
				alertString = nickName + ": " + pushParam
			}
		case PT_PHOTO:
			alertObj.Set("loc-key", "%@: Voice Message")
			alertObj.Set("loc-args", []interface{}{nickName})
		case PT_INDRODUCE:
			alertObj.Set("loc-key", "%1$@: %2$@'s cantact card push")
			alertObj.Set("loc-args", []interface{}{nickName, pushParam})
		case PT_LOCATION:
			alertObj.Set("loc-key", "%@: Location")
			alertObj.Set("loc-args", []interface{}{nickName})
		case PT_FRIEND_INVITE:
			alertObj.Set("loc-key", "%@ wants to add you")
			alertObj.Set("loc-args", []interface{}{nickName})
		case PT_INVITE_ACCEPT:
			alertObj.Set("loc-key", "%@ accepts your request")
			alertObj.Set("loc-args", []interface{}{nickName})
		case PT_LANGUAGE_EXCHANGE:
			if pushParam == APNS_PUSH_LANGUAGE_EXCHANGE ||
				pushParam == APNS_PUSH_LX_ACCEPT ||
				pushParam == APNS_PUSH_LX_DECLIND ||
				pushParam == APNS_PUSH_LX_TERMINATED ||
				pushParam == APNS_PUSH_LX_CANCLED {
				alertString = nickName + ": " + pushParam
			} else {
				alertObj.Set("loc-key", pushParam)
				alertObj.Set("loc-args", []interface{}{nickName})
			}
		case PT_CORRECT_SENTENCE:
			alertObj.Set("loc-key", "%@: Correct Sentences")
			alertObj.Set("loc-args", []interface{}{nickName})
		case PT_DOODLE:
			alertObj.Set("loc-key", "%@: Doodle")
			alertObj.Set("loc-args", []interface{}{nickName})
		case PT_STICKERS:
			alertObj.Set("loc-key", "%@: Stickers")
			alertObj.Set("loc-args", []interface{}{nickName})
		case PT_GIFT:
			alertObj.Set("loc-key", "Gift from %@")
			alertObj.Set("loc-args", []interface{}{nickName})
		case PT_VOIP:
			alertObj.Set("loc-key", "%@: "+pushParam)
			alertObj.Set("loc-args", []interface{}{nickName})
		case PT_VIDEO:
			alertObj.Set("loc-key", "%@: Video")
			alertObj.Set("loc-args", []interface{}{nickName})
		case PT_GVOIP:
			alertObj.Set("loc-key", "Start group call push")
			alertObj.Set("loc-args", []interface{}{nickName})
		case PT_LINK:
			alertObj.Set("loc-key", "Message Preview Example No")
			alertObj.Set("loc-args", []interface{}{nickName})
		case PT_CARD:
			alertObj.Set("loc-key", "%@: Card")
			alertObj.Set("loc-args", []interface{}{nickName})
		case PT_FOLLOW:
			alertObj.Set("loc-key", "%@ has followed you")
			alertObj.Set("loc-args", []interface{}{nickName})
		case PT_REPLY_YOUR_COMMENT:
			alertObj.Set("loc-key", "%@ replied to your comment")
			alertObj.Set("loc-args", []interface{}{nickName})
		case PT_COMMENTED_YOUR_POST:
			alertObj.Set("loc-key", "@ commented your post")
			alertObj.Set("loc-args", []interface{}{nickName})
		case PT_CORRECTED_YOUR_POST:
			alertObj.Set("loc-key", "%@ corrected your post")
			alertObj.Set("loc-args", []interface{}{nickName})
		default:
			log.Printf("UnHandle pushType=%v", pushType)
		}
	}
	apsObj := simplejson.New()
	// add alert obj
	if len(alertString) != 0 {
		apsObj.Set("alert", alertString)
	} else {
		apsObj.Set("alert", alertObj)
	}
	// add badge
	badgeNum := rootObj.Get("badge_num").MustInt(0)
	if badgeNum <= 0 {
		badgeNum = 1
	}
	apsObj.Set("badge", badgeNum)
	//add sound
	apsObj.Set("sound", rootObj.Get("sound").MustInt(0))

	pushMsgObj := simplejson.New()
	// add aps param
	pushMsgObj.Set("aps", apsObj)
	// userid: for the Client UI Jump
	if pushType != PT_FRIEND_INVITE && pushType != PT_VOIP {
		chatType := rootObj.Get("chat_type").MustInt(CTP2P)
		fromId := rootObj.Get("from_id").MustInt(0)

		if chatType == CTMUC {
			pushMsgObj.Set("roomid", fromId)
		} else {
			pushMsgObj.Set("userid", fromId)
		}
	}

	// moment add
	if pushType == PT_REPLY_YOUR_COMMENT ||
		pushType == PT_COMMENTED_YOUR_POST ||
		pushType == PT_CORRECTED_YOUR_POST {
		pushMsgObj.Set("notifyType", "MntCmtNotify")
	} else if pushType == PT_FOLLOW {
		pushMsgObj.Set("notifyType", "FollowingtNotify")
	}
	//add action id
	rootObj.Set("actionid", rootObj.Get("action_id").MustInt(0))
	pushSlice, err := pushMsgObj.MarshalJSON()
	if err != nil {
		log.Printf("PostData MarshalJSON failed")
		return nil
	}

	log.Printf("PostData push message=%s", pushSlice)
	notification := &apns2.Notification{}
	notification.DeviceToken = token
	notification.Topic = TopicId
	notification.Payload = pushSlice // See Payload section below

	client := clientManager.Get(myCert).Production()
	defer clientManager.Add(client)
	res, err := client.Push(notification)
	if err != nil {
		log.Printf("postData client.Push failed err=%v", err)
	} else {
		log.Printf("PostData client.Push ret=%#v", res)
	}

	log.Printf("ProcData response = %#v", res)
	return nil
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// 处理命令行参数
	if _, err := parser.Parse(); err != nil {
		log.Fatalln("main parse cmd line failed!")
	}

	if options.ServerConf == "" {
		log.Fatalln("main Must input config file name")
	}

	// log.Println("config name =", options.ServerConf)
	// 读取配置文件
	cfg, err := ini.Load([]byte(""), options.ServerConf)
	if err != nil {
		log.Fatalln("main load config file=%s failed", options.ServerConf)
		return
	}
	// 配置文件只读 设置此标识提升性能
	cfg.BlockMode = false

	lookupdHost := cfg.Section("LOOKUPD").Key("host").MustString("127.0.0.1:4161")
	topic := cfg.Section("MESSAGE").Key("topic").MustString("test")
	channel := cfg.Section("MESSAGE").Key("chan").MustString("ch")
	cerPath := cfg.Section("CER").Key("path").MustString("./cer.p12")
	myCert, err = certificate.FromP12File(cerPath, CerPassWd)
	if err != nil {
		log.Fatal("Get Cert error:", err)
	}
	clientManager = apns2.NewClientManager()

	// init ansp consumter
	wg = &sync.WaitGroup{}
	config := nsq.NewConfig()
	q, _ := nsq.NewConsumer(topic, channel, config)
	q.AddHandler(nsq.HandlerFunc(MessageHandle))
	err = q.ConnectToNSQLookupd(lookupdHost)
	if err != nil {
		log.Printf("main Could not connect")
	}

	// catchs system signal
	chSig := make(chan os.Signal)
	signal.Notify(chSig, syscall.SIGINT, syscall.SIGTERM)
	log.Println("Signal: ", <-chSig)

}
