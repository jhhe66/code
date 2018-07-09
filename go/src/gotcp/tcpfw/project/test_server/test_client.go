package main

import (
	// "fmt"
	// "github.com/bitly/go-simplejson"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/gansidui/gotcp/libcomm"
	"github.com/gansidui/gotcp/tcpfw/common"
	"github.com/gansidui/gotcp/tcpfw/include/ht_user"
	"github.com/golang/protobuf/proto"
	"github.com/jessevdk/go-flags"
	"gopkg.in/ini.v1"
	// "strings"
	// "time"
)

type Options struct {
	// Example of verbosity with level
	Verbose []bool `short:"v" long:"verbose" description:"Verbose output"`

	// Example of optional value
	ServerConf string `short:"c" long:"conf" description:"Server Config" optional:"no"`
}

var options Options
var infoLog *log.Logger

var parser = flags.NewParser(&options, flags.Default)

func main() {
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
	// 定义一个文件
	fileName := cfg.Section("LOG").Key("path").MustString("/home/ht/server.log")
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

	// 读取ip+port
	serverIp := cfg.Section("OUTER_SERVER").Key("server_ip").MustString("127.0.0.3")
	serverPort := cfg.Section("OUTER_SERVER").Key("server_port").MustInt(8990)

	// 读取发送报文次数
	tryCount := cfg.Section("COUNT").Key("tryCount").MustUint64(10)

	infoLog.Printf("server_ip=%v server_port=%v\n", serverIp, serverPort)
	infoLog.Println("trycount =", tryCount)

	tcpAddr, err := net.ResolveTCPAddr("tcp4", serverIp+":"+strconv.Itoa(serverPort))
	checkError(err)
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err)
	headV3Protocol := &common.HeadV3Protocol{}
	var head *common.HeadV3
	head = &common.HeadV3{Flag: 0xF0, TermType: 1, Cmd: 0x6101, From: 2325928, To: 1946612}
	reqBody := new(ht_user.ReqBody)
	reqBody.User = make([]*ht_user.UserInfoBody, 1)
	reqBody.User[0] = new(ht_user.UserInfoBody)
	reqBody.User[0].UserID = new(uint32)
	*(reqBody.User[0].UserID) = 1946612

	payLoad, err := proto.Marshal(reqBody)
	if err != nil {
		infoLog.Println("marshaling error: ", err)
	}

	head.Len = uint32(common.PacketV3HeadLen) + uint32(len(payLoad)) + 1
	buf := make([]byte, head.Len)
	buf[0] = common.HTV3MagicBegin
	err = common.SerialHeadV3ToSlice(head, buf[1:])
	if err != nil {
		infoLog.Println("SerialHeadV3ToSlice failed")
		os.Exit(1)
	}
	copy(buf[common.PacketV3HeadLen:], payLoad)
	buf[head.Len-1] = common.HTV3MagicEnd

	infoLog.Println("payLoad=", common.NewHeadV3Packet(buf).Serialize())

	var i uint64
	for i = 0; i < tryCount; i++ {
		secondBegin := time.Now().UnixNano()
		// write
		conn.Write(common.NewHeadV3Packet(buf).Serialize())
		attr := "Client/send_packet"
		libcomm.AttrAdd(attr, 1)
		// read
		p, err := headV3Protocol.ReadPacket(conn)
		if err == nil {
			packet := p.(*common.HeadV3Packet)
			// infoLog.Printf("Server reply:[%v] [%v]\n", headV2Packet.GetLength(), string(headV2Packet.GetBody()))
			rspHead, err := packet.GetHead()
			if err == nil {
				infoLog.Printf("head=%+v", rspHead)
				rspBody := &ht_user.RspBody{}
				err = proto.Unmarshal(packet.GetBody(), rspBody)
				if err != nil {
					infoLog.Println("proto Unmarshal failed")
					return
				} else {
					userInfo := rspBody.GetUser()[0]
					infoLog.Printf("uid=%v name=%s sex=%v birthday=%v national=%s native=%v learn1=%v learn2=%v learn3=%v\n", userInfo.GetUserID(), userInfo.GetUserName(), userInfo.GetSex(), userInfo.GetBirthday(), userInfo.GetNational(), userInfo.GetNativeLang(), userInfo.GetLearnLang1(), userInfo.GetLearnLang2(), userInfo.GetLearnLang3())
					infoLog.Printf("T2=%v T3=%v nickName=%s vipTs=%v\n", userInfo.GetTeachLang2(), userInfo.GetTeachLang3(), userInfo.GetNickname(), userInfo.GetVipExpireTs())
				}
			} else {
				infoLog.Println("get HTV3Head failed err =", err)
			}
		}
		secondEnd := time.Now().UnixNano()
		if (secondEnd-secondBegin)/1000000 > 10 {
			attr := "Client/slow_process"
			libcomm.AttrAdd(attr, 1)
			infoLog.Println("Forward proc slow")
		}

	}

	conn.Close()
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
