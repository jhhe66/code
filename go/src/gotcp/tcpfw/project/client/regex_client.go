package main

import (
	"github.com/bitly/go-simplejson"
	"github.com/garyburd/redigo/redis"
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
	outObj := simplejson.New()
	sexWords := "((^|[^A-Za-z|\xc3\x80|\xc3\x81|\xc3\x82|\xc3\x83|\xc3\x84|\xc3\x85|\xc3\x86|\xc3\x87|\xc3\x88|\xc3\x89|\xc3\x8a|\xc3\x8b|\xc3\x8c|\xc3\x8d|\xc3\x8e|\xc3\x8f|\xc3\x90|\xc3\x91|\xc3\x92|\xc3\x93|\xc3\x94|\xc3\x95|\xc3\x96|\xc3\x98|\xc3\x99|\xc3\x9a|\xc3\x9b|\xc3\x9c|\xc3\x9d|\xc3\x9e|\xc3\x9f|\xc3\xa0|\xc3\xa1|\xc3\xa2|\xc3\xa3|\xc3\xa4|\xc3\xa5|\xc3\xa6|\xc3\xa7|\xc3\xa8|\xc3\xa9|\xc3\xaa|\xc3\xab|\xc3\xac|\xc3\xad|\xc3\xae|\xc3\xaf|\xc3\xb1|\xc3\xb2|\xc3\xb3|\xc3\xb4|\xc3\xb5|\xc3\xb6|\xc3\xb8|\xc3\xb9|\xc3\xba|\xc3\xbb|\xc3\xbc|\xc3\xbd|\xc3\xbe|\xc3\xbf])(?<name>(bitch|blowjob|dick|dicks|dildo|fuck|fucks|motherfucker|nigga|penis|whore|Virgin|Vagina|boob|masturbate|cock|condom|genital|horny|jerk off|nipple|wanking|ass|sex|tit|cum|tits|masturbation|virginity|wank|naked|nude|pussy|tetas|\x63\x6f\xc3\xb1\x6f|pechos|polla|culo|sikik|meme|sikeyim|amcik|transar|\x63\xc3\xb3\xc3\xb1\x61|sexe|sexy|puttana|\xd1\x81\xd0\xb5\xd0\xba\xd1\x81\xd0\xbe\xd0\xbc))([^A-Za-z|\xc3\x80|\xc3\x81|\xc3\x82|\xc3\x83|\xc3\x84|\xc3\x85|\xc3\x86|\xc3\x87|\xc3\x88|\xc3\x89|\xc3\x8a|\xc3\x8b|\xc3\x8c|\xc3\x8d|\xc3\x8e|\xc3\x8f|\xc3\x90|\xc3\x91|\xc3\x92|\xc3\x93|\xc3\x94|\xc3\x95|\xc3\x96|\xc3\x98|\xc3\x99|\xc3\x9a|\xc3\x9b|\xc3\x9c|\xc3\x9d|\xc3\x9e|\xc3\x9f|\xc3\xa0|\xc3\xa1|\xc3\xa2|\xc3\xa3|\xc3\xa4|\xc3\xa5|\xc3\xa6|\xc3\xa7|\xc3\xa8|\xc3\xa9|\xc3\xaa|\xc3\xab|\xc3\xac|\xc3\xad|\xc3\xae|\xc3\xaf|\xc3\xb1|\xc3\xb2|\xc3\xb3|\xc3\xb4|\xc3\xb5|\xc3\xb6|\xc3\xb8|\xc3\xb9|\xc3\xba|\xc3\xbb|\xc3\xbc|\xc3\xbd|\xc3\xbe|\xc3\xbf]|$)+?|((?<cname>(\xe6\x93\x8d\xe4\xbd\xa0|\xe4\xb9\xb3\xe6\x88\xbf|\xe8\xa3\xb8\xe4\xbd\x93|\xe5\xa4\x84\xe5\xa5\xb3|\xe5\xa9\x8a\xe5\xad\x90|\xe5\x81\x9a\xe7\x88\xb1|\xe9\xb8\xa1\xe5\xb7\xb4|\xe9\x98\xb4\xe8\x8c\x8e|\xe4\xbd\xa0\xe5\xa6\x88\xe9\x80\xbc|\xe5\xa5\xb3\xe4\xbc\x98|\xe7\xa7\x81\xe5\xa4\x84|\xe5\xb0\x84\xe7\xb2\xbe|\xe3\x81\xbc\xe3\x81\x91|\xe3\x81\x93\xe3\x81\xae\xe3\x82\x84\xe3\x82\x8d\xe3\x81\x86|\xe3\x81\xbe\xe3\x82\x93\xe3\x81\x93|\xe3\x81\x8a\xe3\x81\xa3\xe3\x81\xb1\xe3\x81\x84|\xe3\x82\xbb\xe3\x83\x83\xe3\x82\xaf\xe3\x82\xb9|\xe3\x82\xa8\xe3\x83\x83\xe3\x83\x81|\xe3\x82\xaa\xe3\x83\x8a\xe3\x83\x8b\xe3\x83\xbc|chikubi|\xe3\x83\x87\xe3\x82\xab\xe3\x83\x81\xe3\x83\xb3|\xe3\x83\xa0\xe3\x83\xa9\xe3\x83\xa0\xe3\x83\xa9|\xeb\xb3\x91\xec\x8b\xa0|\xec\x84\xb1\xea\xb8\xb0|\xec\x84\xb9\xec\x8a\xa4))))"
	outObj.Set("sexword", sexWords)

	sPacket, err := outObj.MarshalJSON()
	if err != nil {
		infoLog.Println("MarshalJSON failed outObj=", outObj)
		return
	}

	ret, err := c.Do("SET", "regex_model", sPacket) // get uid
	if err != nil {
		infoLog.Println("Reids failed set [err, ret] =", err, ret)
		return
	}

	n, err := c.Do("Get", "regex_model") // get uid
	strValue, err := redis.String(n, err)
	if err != nil {
		infoLog.Println("Reids failed Get regex_model err =", err)
	} else {
		infoLog.Println("regex_model =", strValue)
	}
	return
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
