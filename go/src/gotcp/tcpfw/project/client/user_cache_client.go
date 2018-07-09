package main

import (
	// "fmt"
	// "github.com/bitly/go-simplejson"
	"log"
	"net"
	"os"
	"strconv"

	"github.com/gansidui/gotcp/tcpfw/common"
	"github.com/gansidui/gotcp/tcpfw/include/ht_user"
	"github.com/golang/protobuf/proto"
	"github.com/jessevdk/go-flags"
	"gopkg.in/ini.v1"
	// "strings"
)

var (
	infoLog *log.Logger
)

type Options struct {
	// Example of verbosity with level
	Verbose []bool `short:"v" long:"verbose" description:"Verbose output"`

	// Example of optional value
	Config string `short:"c" long:"config" description:"config file" optional:"no"`

	Cmd uint32 `short:"t" long:"command" description:"test commad" optional:"no"`
}

var options Options

var parser = flags.NewParser(&options, flags.Default)

func main() {
	// 处理命令行参数
	if _, err := parser.Parse(); err != nil {
		log.Fatalln("parse cmd line failed!")
		os.Exit(1)
	}
	if options.Config == "" || Cmd == 0 {
		log.Fatalln("Must input config file name and Cmd 1:GetUserInfo 2:ModifyUserInfo")
	}

	// log.Println("config name =", options.ClientConf)
	// 读取配置文件
	cfg, err := ini.Load([]byte(""), options.Config)
	if err != nil {
		log.Printf("load config file=%s failed", options.Config)
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

	// 读取ip+port
	serverIp := cfg.Section("OUTER_SERVER").Key("server_ip").MustString("127.0.0.3")
	serverPort := cfg.Section("OUTER_SERVER").Key("server_port").MustInt(8990)

	infoLog.Printf("server_ip=%v server_port=%v\n", serverIp, serverPort)

	tcpAddr, err := net.ResolveTCPAddr("tcp4", serverIp+":"+strconv.Itoa(serverPort))
	checkError(err)
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err)

	uid := cfg.Section("TEST_UID").Key("UID").MustInt(1603139)
	headV2Protocol := &common.HeadV2Protocol{}
	head := &common.HeadV2{Cmd: uint32(ht_user.ProtoCmd(options.Cmd)), Len: uint32(tcpfw.EmptyPacktV2Len), Uid: uint32(uid)}
	reqBody := new(ht_user.ReqBody)
	reqBody.User = make([]*ht_user.UserInfoBody, 1)
	reqBody.User[0] = new(ht_user.UserInfoBody)
	reqBody.User[0].UserID = new(uint32)
	*(reqBody.User[0].UserID) = uint32(uid)

	payLoad, err := proto.Marshal(reqBody)
	if err != nil {
		infoLog.Println("marshaling error: ", err)
	}

	head.Len = uint32(common.PacketV2HeadLen) + uint32(len(payLoad)) + 1
	buf := make([]byte, head.Len)
	buf[0] = common.HTV2MagicBegin
	err = common.SerialHeadV2ToSlice(head, buf[1:])
	if err != nil {
		infoLog.Println("SerialHeadV3ToSlice failed")
		os.Exit(1)
	}
	copy(buf[common.PacketV2HeadLen:], payLoad)
	buf[head.Len-1] = common.HTV2MagicEnd

	infoLog.Println("payLoad=", common.NewHeadV2Packet(buf).Serialize())
	// write
	conn.Write(common.NewHeadV2Packet(buf).Serialize())

	// read
	p, err := headV2Protocol.ReadPacket(conn)
	if err == nil {
		packet := p.(*common.HeadV2Packet)
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
			infoLog.Println("get HTV2Head failed err =", err)
		}
	}

	conn.Close()
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
