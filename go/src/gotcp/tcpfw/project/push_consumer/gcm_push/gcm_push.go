package main

import (
	"log"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"

	"github.com/alexjlockwood/gcm"
	simplejson "github.com/bitly/go-simplejson"
	"github.com/jessevdk/go-flags"
	"github.com/nsqio/go-nsq"
	"gopkg.in/ini.v1"
)

var wg *sync.WaitGroup

type Options struct {
	// Example of verbosity with level
	Verbose []bool `short:"v" long:"verbose" description:"Verbose output"`

	// Example of optional value
	ServerConf string `short:"c" long:"conf" description:"Server Config" optional:"no"`
}

const (
	GCM_PUSH_NOPREVIEW     = "no_preview"
	GCM_PUSH_TEXT          = "text"
	GCM_PUSH_VOICE         = "voice"
	GCM_PUSH_IMAGE         = "image"
	GCM_PUSH_LOCATION      = "location"
	GCM_PUSH_INTRODUCE     = "introduce"
	GCM_PUSH_CORRECT       = "correction"
	GCM_PUSH_STICKER       = "sticker"
	GCM_PUSH_DOODLE        = "doodle"
	GCM_PUSH_INVITE        = "friend_invite"
	GCM_PUSH_LX            = "language_exchange"
	GCM_PUSH_LX_REPLY      = "language_exchange_reply" // Reply Code "Declined" "Accepted" "Terminated"
	GCM_PUSH_GIFT          = "gift"
	GCM_PUSH_CALL_INCOMING = "call_incoming"
	GCM_PUSH_CALL_CANCEL   = "call_cancel"
	GCM_PUSH_CALL_MISS     = "call_miss"
	GCM_PUSH_ACCEPT_INVITE = "accept_invite"
	GCM_PUSH_VIDEO         = "video" // add video push
	GCM_PUSH_GVOIP         = "gvoip"
	GCM_PUSH_LINK          = "message_preview_example_no"
	GCM_PUSH_CARD          = "card"

	// 2016-08-26 add by songliwei
	GCM_PUSH_FOLLOW              = "s_has_followed_you"
	GCM_PUSH_REPLY_YOUR_COMMENT  = "s_replied_your_comment"
	GCM_PUSH_COMMENTED_YOUR_POST = "s_commented_your_post"
	GCM_PUSH_CORRECTED_YOUR_POST = "s_corrected_your_post"
)

const (
	APIKey = "AIzaSyCOu1hU2Moz-GuqpjiLavpoNUNwfxrwxBc"
)

var options Options
var parser = flags.NewParser(&options, flags.Default)

func MessageHandle(message *nsq.Message) error {
	log.Printf("MessageHandle Got a message: %v", message)
	wg.Add(1)
	go func() {
		PostData(message)
		wg.Done()
	}()

	return nil
}

func PostData(message *nsq.Message) error {
	log.Println("PostData enter msgId=%v timestamp=%v body=%s", message.ID, message.Timestamp, message.Body)
	// 然后进行下一步无处理
	rootObj, err := simplejson.NewJson(message.Body)
	if err != nil {
		log.Printf("ProcData simplejson new packet error", err)
		return err
	}
	log.Printf("ProcData rootObj=%#v", rootObj)

	regId := rootObj.Get("registration_id").MustString("")
	registerionIds := []string{regId}

	strPushType := rootObj.Get("push_type").MustString("")
	data := map[string]interface{}{"type": strPushType}
	if strPushType != GCM_PUSH_NOPREVIEW {
		data["from_id"] = rootObj.Get("from_id").MustInt64(0)
		data["sender"] = rootObj.Get("sender").MustString("")
	}
	if strPushType == GCM_PUSH_TEXT {
		data["message"] = rootObj.Get("push_param").MustString("")
	}

	if strPushType == GCM_PUSH_INTRODUCE {
		data["user"] = rootObj.Get("push_param").MustString("")
	}

	if strPushType == GCM_PUSH_LX_REPLY {
		data["reply_code"] = rootObj.Get("push_param").MustString("")
	}
	data["sound"] = rootObj.Get("sound").MustInt64(0)
	data["msg_id"] = rootObj.Get("msg_id").MustString("")
	data["to_id"] = rootObj.Get("to_id").MustInt64(0)
	data["actionId"] = rootObj.Get("action_id").MustInt64(0)
	data["byAt"] = rootObj.Get("by_at").MustInt64(0)
	// Create the message to be sent.
	pushMsg := gcm.NewMessage(data, registerionIds...)
	// Create a Sender to send the message.
	sender := &gcm.Sender{ApiKey: APIKey}
	// Send the message and receive the response after at most two retries.
	response, err := sender.Send(pushMsg, 2)
	if err != nil {
		log.Println("Failed to send message:", err)
		return err
	}
	log.Printf("ProcData response = %#v", response)

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
