package main

import (
	"log"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"

	"github.com/aiwuTech/xinge"
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
	TYPE_ACTIVITY = 1
	TYPE_URL      = 2
	TYPE_INTENT   = 3
	TYPE_PACKAGE  = 4
)

const (
	AccessId  = "2100093818"
	AccessKey = "A39V82PU7WMC"
	SecretKey = "42deaa73626e172176fbe6a27e730de5"
)

const (
	P2P_CHAT = "CHAT"
	MUC_CHAT = "MUC"
)

var options Options
var parser = flags.NewParser(&options, flags.Default)

var (
	xingeClient = xinge.NewClient(AccessId, 300, AccessKey, SecretKey)
)

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
	log.Printf("PostData enter msgId=%v timestamp=%v body=%s", message.ID, message.Timestamp, message.Body)
	// 然后进行下一步无处理
	rootObj, err := simplejson.NewJson(message.Body)
	if err != nil {
		log.Printf("PostData simplejson new packet error", err)
		return err
	}

	log.Printf("ProcData rootObj=%#v", rootObj)

	var nullString string
	activityAttr := &xinge.ActivityAttr{IF: 0, PF: 0}
	browser := &xinge.Browser{Confirm: 1}
	packetName := &xinge.Package{Confirm: 1}
	action := &xinge.AndroidAction{ActionType: TYPE_ACTIVITY,
		AtyAttr:     activityAttr,
		Browser:     browser,
		PackageName: packetName}
	style := NewStyle(1, rootObj.Get("sound").MustInt(1), 1, 1, 0, rootObj.Get("lights").MustInt(1), 0, 1)
	// 自定义字段
	customObj := map[string]interface{}{}
	customObj["chat_type"] = rootObj.Get("chat_type").MustString("CHAT")
	customObj["sound"] = rootObj.Get("sound").MustInt(1)
	customObj["lights"] = rootObj.Get("lights").MustInt(1)
	customObj["vibrate"] = rootObj.Get("sound").MustInt(0)
	customObj["type"] = rootObj.Get("msg_type").MustString("text")
	customObj["msg_id"] = rootObj.Get("msg_id").MustString("")
	customObj["from_id"] = uint32(rootObj.Get("from_id").MustInt64(0))
	customObj["to_id"] = uint32(rootObj.Get("to_id").MustInt64(0))
	customObj["actionid"] = uint32(rootObj.Get("actionid").MustInt64(0))
	customObj["byAt"] = uint32(rootObj.Get("byAt").MustInt64(0))
	log.Printf("procData customObj=%#v fromId=%v", customObj, uint32(rootObj.Get("from_id").MustInt64(0)))

	xingeMessage := &xinge.AndroidMessage{
		Title:         rootObj.Get("title").MustString(""),
		Content:       rootObj.Get("content").MustString(""),
		AcceptTime:    []*AcceptTime{},
		NotifyId:      0,
		BuilderId:     1,
		Ring:          rootObj.Get("sound").MustInt(1),
		Vibrate:       1,
		Lights:        rootObj.Get("lights").MustInt(1),
		Clearable:     1,
		IconType:      0,
		StyleId:       1,
		Action:        action,
		CustomContent: customObj,
	}
	reqPush := &xinge.ReqPush{
		PushType:     xinge.PushType_single_device,
		DeviceToken:  rootObj.Get("token").MustString(""),
		MessageType:  xinge.MessageType_notify,
		Message:      xingeMessage,
		ExpireTime:   300,
		SendTime:     time.Now(),
		MultiPkgType: xinge.MultiPkg_aid,
		PushEnv:      xinge.PushEnv_android,
		PlatformType: xinge.Platform_android,
		LoopTimes:    2,
		LoopInterval: 7,
		Cli:          xingeClient,
	}

	err := reqPush.Push()
	if err != nil {
		log.Printf("PostData reqPush.Push failed")
		return err
	} else {
		log.Printf("PostData reqPush.Push success")
	}
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
