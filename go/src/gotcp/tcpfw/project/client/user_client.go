package main

import (
	// "github.com/bitly/go-simplejson"
	"github.com/gansidui/gotcp/examples/ht_user"
	"github.com/garyburd/redigo/redis"
	"github.com/golang/protobuf/proto"
	"github.com/jessevdk/go-flags"
	"gopkg.in/ini.v1"
	"log"
	"os"
	"strconv"
	// "time"
)

type Options struct {
	// Example of verbosity with level
	Verbose []bool `short:"v" long:"verbose" description:"Verbose output"`

	// Example of optional value
	ClientConf string `short:"c" long:"conf" description:"Clinet Config" optional:"no"`
}

var options Options
var parser = flags.NewParser(&options, flags.Default)

var infoLog *log.Logger

func main() {
	// 处理命令行参数
	if _, err := parser.Parse(); err != nil {
		log.Fatalln("parse cmd line failed!")
	}

	if options.ClientConf == "" {
		log.Fatalln("Must input config file name")
	}

	// log.Println("config name =", options.ClientConf)
	// 读取配置文件
	cfg, err := ini.Load([]byte(""), options.ClientConf)
	if err != nil {
		log.Printf("load config file=%s failed", options.ClientConf)
		return
	}
	// 配置文件只读 设置此标识提升性能
	cfg.BlockMode = false
	// 定义一个文件
	fileName := cfg.Section("LOG").Key("path").MustString("/home/ht/clinet.log")
	logFile, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	defer logFile.Close()
	if err != nil {
		log.Fatalln("open file error !")
		return
	}
	uid := cfg.Section("TEST").Key("UID").MustString("200002")

	// 创建一个日志对象
	infoLog = log.New(logFile, "[Info]", log.LstdFlags)
	// 配置log的Flag参数
	infoLog.SetFlags(infoLog.Flags() | log.LstdFlags)

	// init redis pool
	redisIp := cfg.Section("REDIS").Key("redis_ip").MustString("127.0.0.1")
	redisPort := cfg.Section("REDIS").Key("redis_port").MustInt(6379)
	infoLog.Printf("redis ip=%v port=%v", redisIp, redisPort)
	c, err := redis.Dial("tcp", redisIp+":"+strconv.Itoa(redisPort))
	if err != nil {
		// handle error
	}
	defer c.Close()
	n, err := c.Do("Get", uid)
	if err != nil {
		infoLog.Println("Reids failed Get err =", err)
		return
	}

	value, err := redis.String(n, err)
	if err != nil {
		infoLog.Println("redis String fail")
		return
	}

	reqBody := &ht_user.UserInfoBody{}
	err = proto.Unmarshal([]byte(value), reqBody)
	if err != nil {
		infoLog.Println("proto Unmarshal failed")
		return
	}

	infoLog.Printf("UserID=%v UserName=%v Sex=%v Birthday=%v", *reqBody.UserID, string(reqBody.UserName), *reqBody.Sex, *reqBody.Birthday)
	infoLog.Printf("National=%v, NativeLang=%v LearnLang1=%v LearnLang2=%v LearnLang3=%v TeachLang2=%v TeachLang3=%v",
		string(reqBody.National),
		*reqBody.NativeLang,
		*reqBody.LearnLang1,
		*reqBody.LearnLang2,
		*reqBody.LearnLang3,
		*reqBody.TeachLang2,
		*reqBody.TeachLang3)
	if reqBody.VipExpireTs == nil {
		reqBody.VipExpireTs = new(uint32)
		*reqBody.VipExpireTs = 0
	}
	infoLog.Printf("BlackidList=%v, FriendidList=%v FansList=%v Nickname=%v VipExpireTs=%v",
		reqBody.BlackidList,
		reqBody.FriendidList,
		reqBody.FansList,
		string(reqBody.Nickname),
		*reqBody.VipExpireTs)
	return
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
