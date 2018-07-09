package main

import (
	// "fmt"
	// "github.com/bitly/go-simplejson"
	"log"
	"net"
	"os"
	"strconv"

	"github.com/gansidui/gotcp/tcpfw/common"
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

	Cmd int `short:"t" long:"cmd" description:"Command type" optional:"no"`
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

	infoLog.Printf("server_ip=%v server_port=%v\n", serverIp, serverPort)

	tcpAddr, err := net.ResolveTCPAddr("tcp4", serverIp+":"+strconv.Itoa(serverPort))
	checkError(err)
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err)
	defer conn.Close()

	headXTProtocol := &common.XTHeadProtocol{}
	var head *common.XTHead
	head = &common.XTHead{Flag: 0xF0,
		Version:  4,
		CryKey:   0,
		TermType: 0,
		Cmd:      0x00,
		Seq:      0x1000,
		From:     1946612,
		To:       0,
		Len:      0,
	}

	var payLoad []byte
	switch options.Cmd {
	case 0x2007:
		head.Cmd = 0x2007
		common.MarshalUint32(uint32(1946612), &payLoad)
		common.MarshalSlice([]byte("Chinese"), &payLoad)
		common.MarshalUint8(uint8(1), &payLoad)
		if true {
			common.MarshalSlice([]byte("22.539706"), &payLoad)
			common.MarshalSlice([]byte("113.991826"), &payLoad)
			common.MarshalSlice([]byte("CN"), &payLoad)
			common.MarshalSlice([]byte("Guangdong Sheng"), &payLoad)
			common.MarshalSlice(nil, &payLoad)
			common.MarshalSlice(nil, &payLoad)
			common.MarshalSlice([]byte("Shenzhen Shi"), &payLoad)
			common.MarshalSlice([]byte("Nanshan Qu"), &payLoad)
			common.MarshalSlice([]byte("HuaQiaoCheng"), &payLoad)
			common.MarshalSlice([]byte("google"), &payLoad)
		}
	case 0x2009:
		head.Cmd = 0x2009
		common.MarshalSlice([]byte("Chinese"), &payLoad)
		common.MarshalUint16(uint16(3), &payLoad)
		common.MarshalUint32(uint32(1946612), &payLoad)
		common.MarshalUint32(uint32(2325928), &payLoad)
		common.MarshalUint32(uint32(900866), &payLoad)
	case 0x2035:
		head.Cmd = 0x2035
		common.MarshalSlice([]byte("Chinese"), &payLoad)
		common.MarshalUint32(uint32(1946612), &payLoad)
		common.MarshalUint64(uint64(1477382696), &payLoad)
	default:
		infoLog.Println("UnKnow input cmd =", options.Cmd)

	}
	head.Len = uint32(len(payLoad))
	packetSlice := make([]byte, int(common.PacketXTHeadLen+head.Len))
	common.SerialXTHeadToSlice(head, packetSlice[:])
	copy(packetSlice[common.PacketXTHeadLen:], payLoad)
	infoLog.Printf("len=%v payLaod=%v\n", len(payLoad), payLoad)
	// write
	conn.Write(packetSlice)
	// read
	p, err := headXTProtocol.ReadPacket(conn)
	if err == nil {
		rspPacket := p.(*common.XTHeadPacket)
		rspHead, _ := rspPacket.GetHead()
		rspPayLoad := rspPacket.GetBody()
		infoLog.Printf("resp len=%v cmd=%v\n", rspHead.Len, rspHead.Cmd)
		switch rspHead.Cmd {
		case 0x2008:
			ret := common.UnMarshalUint8(&rspPayLoad)
			allow := common.UnMarshalUint8(&rspPayLoad)
			infoLog.Printf("ret=%v allow=%v\n", ret, allow)
			if allow == uint8(1) {
				country := common.UnMarshalSlice(&rspPayLoad)
				city := common.UnMarshalSlice(&rspPayLoad)
				infoLog.Printf("contry=%s city=%s\n", country, city)
			}
		case 0x200A:
			ret := common.UnMarshalUint8(&rspPayLoad)
			count := common.UnMarshalUint16(&rspPayLoad)
			var i uint16 = 0
			for i = 0; i < count; i++ {
				tagStr := common.UnMarshalSlice(&rspPayLoad)
				targetUid := common.UnMarshalUint32(&tagStr)
				targetTs := common.UnMarshalUint64(&tagStr)
				targetAllow := common.UnMarshalUint8(&tagStr)
				if targetAllow == 1 {
					targetLati := common.UnMarshalSlice(&tagStr)
					targetLong := common.UnMarshalSlice(&tagStr)
					targetCountry := common.UnMarshalSlice(&tagStr)
					targetCity := common.UnMarshalSlice(&tagStr)
					infoLog.Printf("lati=%s long=%s country=%s city=%s\n", targetLati, targetLong, targetCountry, targetCity)
				}
				infoLog.Printf("ret=%v uid=%v ts=%v allow=%v\n", ret, targetUid, targetTs, targetAllow)
			}
		case 0x2036:
			ret := common.UnMarshalUint8(&rspPayLoad)
			tagStr := common.UnMarshalSlice(&rspPayLoad)
			targetUid := common.UnMarshalUint32(&tagStr)
			targetTs := common.UnMarshalUint64(&tagStr)
			targetAllow := common.UnMarshalUint8(&tagStr)
			if targetAllow == 1 {
				targetLati := common.UnMarshalSlice(&tagStr)
				targetLong := common.UnMarshalSlice(&tagStr)
				targetCountry := common.UnMarshalSlice(&tagStr)
				targetCity := common.UnMarshalSlice(&tagStr)
				infoLog.Printf("lati=%s long=%s country=%s city=%s\n", targetLati, targetLong, targetCountry, targetCity)
			}
			infoLog.Printf("ret=%v uid=%v ts=%v allow=%v\n", ret, targetUid, targetTs, targetAllow)
		default:
			infoLog.Println("UnKnow resp cmd =", rspHead.Cmd)
		}

	}
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
