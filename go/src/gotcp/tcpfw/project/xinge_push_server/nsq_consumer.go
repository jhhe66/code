package main

import (
	"log"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"

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

var options Options
var parser = flags.NewParser(&options, flags.Default)

func MessageHandle(message *nsq.Message) error {
	log.Printf("Got a message: %v", message)
	wg.Add(1)
	go func() {
		PostData(nil)
		wg.Done()
	}()

	return nil
}
func PostData(payLoad []byte) error {
	log.Println("enter")
	return nil
}

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
		log.Fatalln("load config file=%s failed", options.ServerConf)
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
		log.Panic("Could not connect")
	}

	// catchs system signal
	chSig := make(chan os.Signal)
	signal.Notify(chSig, syscall.SIGINT, syscall.SIGTERM)
	log.Println("Signal: ", <-chSig)

}
