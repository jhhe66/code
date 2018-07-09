package main

import (
	// "github.com/bitly/go-simplejson"
	"database/sql"
	"github.com/gansidui/gotcp/examples/ht_user"
	"github.com/garyburd/redigo/redis"
	_ "github.com/go-sql-driver/mysql"
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

	db, err := sql.Open("mysql", mysqlUser+":"+mysqlPasswd+"@"+"tcp("+mysqlHost+":"+mysqlPort+")/"+mysqlDbName+"?charset=utf8&timeout=90s")
	if err != nil {
		infoLog.Println("open mysql failed")
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	defer db.Close()

	var uid = 10000
	for uid < 5000000 {
		// for uid < 50000 {
		// Execute the query
		rows, err := db.Query("SELECT USERID, NICKNAME FROM HT_USER_BASE WHERE USERID >" + strconv.Itoa(uid) + " limit 1000")
		if err != nil {
			infoLog.Println("mysql Query failed uid =", uid)
			panic(err.Error()) // proper error handling instead of panic in your app
		}

		var sUid, sNickName []byte
		curUid := uid
		infoLog.Println("cur uid =", uid)
		for rows.Next() {
			// get RawBytes from data
			err = rows.Scan(&sUid, &sNickName)
			if err != nil {
				infoLog.Println("mysql Scan failed")
				panic(err.Error()) // proper error handling instead of panic in your app
			}

			iUid, err := strconv.Atoi(string(sUid))
			if err != nil {
				infoLog.Printf("strconv atoi failed uid=%v", sUid)
				continue
			}

			if len(sNickName) == 0 { // nickname 为空跳过设置
				infoLog.Println("empty nickname uid =", iUid)
				continue
			}

			infoLog.Printf("uid=%v nickName=%v", iUid, string(sNickName))

			curUid = iUid // 更新已经获取过的uid

			n, err := c.Do("Get", sUid) // get uid
			if err != nil {
				infoLog.Println("Reids failed Get [err, uid] =", err, iUid)
				continue
			}

			value, err := redis.String(n, err)
			if err != nil {
				infoLog.Println("redis String fail [err, uid] =", err, iUid)
				continue
			}
			reqBody := &ht_user.UserInfoBody{}
			err = proto.Unmarshal([]byte(value), reqBody)
			if err != nil {
				infoLog.Println("proto Unmarshal failed uid=", iUid)
				continue
			}

			// add vip expire ts
			reqBody.Nickname = make([]byte, len(sNickName))
			copy(reqBody.Nickname, sNickName)

			infoLog.Printf("UserID=%v UserName=%v Sex=%v Birthday=%v", *reqBody.UserID, string(reqBody.UserName), *reqBody.Sex, *reqBody.Birthday)
			infoLog.Printf("National=%v, NativeLang=%v LearnLang1=%v LearnLang2=%v LearnLang3=%v TeachLang2=%v TeachLang3=%v",
				string(reqBody.National),
				*reqBody.NativeLang,
				*reqBody.LearnLang1,
				*reqBody.LearnLang2,
				*reqBody.LearnLang3,
				*reqBody.TeachLang2,
				*reqBody.TeachLang3)
			infoLog.Printf("Nickname=%v", string(reqBody.Nickname))

			payLoad, err := proto.Marshal(reqBody)
			if err != nil {
				infoLog.Println("marshaling error: ", err)
			}

			ret, err := c.Do("SET", sUid, payLoad) // get uid
			if err != nil {
				infoLog.Println("Reids failed set [err, uid, ret] =", err, sUid, ret)
				continue
			}

		}

		if curUid == uid {
			uid += 1000
		} else {
			uid = curUid
		}
	}

	// n, err := c.Do("Get", uid)
	// if err != nil {
	// 	infoLog.Println("Reids failed Get err =", err)
	// 	return
	// }

	// value, err := redis.String(n, err)
	// if err != nil {
	// 	infoLog.Println("redis String fail")
	// 	return
	// }
	// reqBody := &ht_user.UserInfoBody{}
	// err = proto.Unmarshal([]byte(value), reqBody)
	// if err != nil {
	// 	infoLog.Println("proto Unmarshal failed")
	// 	return
	// }

	// infoLog.Printf("UserID=%v UserName=%v Sex=%v Birthday=%v", *reqBody.UserID, string(reqBody.UserName), *reqBody.Sex, *reqBody.Birthday)
	// infoLog.Printf("National=%v, NativeLang=%v LearnLang1=%v LearnLang2=%v LearnLang3=%v TeachLang2=%v TeachLang3=%v",
	// 	string(reqBody.National),
	// 	*reqBody.NativeLang,
	// 	*reqBody.LearnLang1,
	// 	*reqBody.LearnLang2,
	// 	*reqBody.LearnLang3,
	// 	*reqBody.TeachLang2,
	// 	*reqBody.TeachLang3)
	// infoLog.Printf("BlackidList=%v, FriendidList=%v FansList=%v Nickname=%v",
	// 	reqBody.BlackidList,
	// 	reqBody.FriendidList,
	// 	reqBody.FansList,
	// 	string(reqBody.Nickname))
	return
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
