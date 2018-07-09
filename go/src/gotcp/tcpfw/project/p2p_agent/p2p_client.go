package main

import (
	// "fmt"
	// "github.com/bitly/go-simplejson"
	"github.com/gansidui/gotcp/examples/ht_p2p"
	"github.com/gansidui/gotcp/examples/tcpfw"
	"github.com/golang/protobuf/proto"
	"github.com/jessevdk/go-flags"
	"gopkg.in/ini.v1"
	"log"
	"net"
	"os"
	"strconv"
	// "strings"
	// "time"
)

type Options struct {
	// Example of verbosity with level
	Verbose []bool `short:"v" long:"verbose" description:"Verbose output"`

	// Example of optional value
	Count string `short:"c" long:"count" description:"send count" optional:"no"`
}

var options Options

var parser = flags.NewParser(&options, flags.Default)

func main() {
	// 处理命令行参数
	if _, err := parser.Parse(); err != nil {
		log.Fatalln("parse cmd line failed!")
		os.Exit(1)
	}
	tryCount, _ := strconv.Atoi(options.Count)

	// 定义一个文件
	fileName := "p2p_client.log"
	logFile, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	defer logFile.Close()
	if err != nil {
		log.Fatalln("open file error !")
	}
	// 创建一个日志对象
	infoLog := log.New(logFile, "[Info]", log.LstdFlags)
	// 配置log的Flag参数
	infoLog.SetFlags(infoLog.Flags() | log.LstdFlags)

	// 读取配置文件
	cfg, err := ini.Load([]byte(""), "conf.ini")
	if err != nil {
		infoLog.Println("load config file failed")
		return
	}
	// 配置文件只读 设置次标识提升性能
	cfg.BlockMode = false
	// 读取ip+port
	serverIp := cfg.Section("OUTER_SERVER").Key("server_ip").MustString("127.0.0.3")
	serverPort := cfg.Section("OUTER_SERVER").Key("server_port").MustInt(8990)

	infoLog.Printf("server_ip=%v server_port=%v\n", serverIp, serverPort)
	infoLog.Println("trycount =", tryCount)

	tcpAddr, err := net.ResolveTCPAddr("tcp4", serverIp+":"+strconv.Itoa(serverPort))
	checkError(err)
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err)
	for i := 0; i < tryCount; i++ {
		headV3Protocol := &tcpfw.HeadV3Protocol{}
		var head *tcpfw.HeadV3
		head = &tcpfw.HeadV3{Flag: 0xF0, TermType: 1, Cmd: 0x6101, From: 2325928, To: 1946612}
		reqSubBody := new(ht_p2p.P2PMsgBody)
		content := []byte{'H', 'e', 'l', 'l', 'o', ',', 'W', 'o', 'r', 'l', 'd'}
		reqSubBody.P2PData = make([]byte, len(content))
		copy(reqSubBody.P2PData, content)

		reqSubBody.JustOnline = new(uint32)
		*(reqSubBody.JustOnline) = 0
		reqSubBody.DontReply = new(uint32)
		*(reqSubBody.DontReply) = 0
		reqSubBody.PushInfo = new(ht_p2p.PushInfo)
		reqSubBody.PushInfo.PushType = new(uint32)
		*(reqSubBody.PushInfo.PushType) = 19
		nickName := []byte{'s', 'o', 'n', 'g', 'l', 'i', 'w', 'e', 'i'}
		reqSubBody.PushInfo.NickName = make([]byte, len(nickName))
		copy(reqSubBody.PushInfo.NickName, nickName)
		pushContent := make([]byte, tcpfw.EmptyPacktV3Len)
		err = tcpfw.SerialHeadV3ToSlice(head, pushContent)
		reqSubBody.PushInfo.Content = make([]byte, len(pushContent))
		copy(reqSubBody.PushInfo.Content, pushContent)

		payLoad, err := proto.Marshal(reqSubBody)
		if err != nil {
			infoLog.Println("marshaling error: ", err)
		}

		head.Len = uint32(tcpfw.PacketV3HeadLen) + uint32(len(payLoad)) + 1
		buf := make([]byte, head.Len)
		buf[0] = tcpfw.HTV3MagicBegin
		err = tcpfw.SerialHeadV3ToSlice(head, buf[1:])
		if err != nil {
			infoLog.Println("SerialHeadV3ToSlice failed")
			os.Exit(1)
		}
		copy(buf[tcpfw.PacketV3HeadLen:], payLoad)
		buf[head.Len-1] = tcpfw.HTV3MagicEnd

		infoLog.Println("payLoad=", tcpfw.NewHeadV3Packet(buf).Serialize())
		// write
		conn.Write(tcpfw.NewHeadV3Packet(buf).Serialize())

		// read
		p, err := headV3Protocol.ReadPacket(conn)
		if err == nil {
			packet := p.(*tcpfw.HeadV3Packet)
			// infoLog.Printf("Server reply:[%v] [%v]\n", headV2Packet.GetLength(), string(headV2Packet.GetBody()))
			rspHead, err := packet.GetHead()
			if err == nil {
				infoLog.Printf("head=%+v", rspHead)
			} else {
				infoLog.Println("get HTV3Head failed err =", err)
			}
		}
	}

	conn.Close()
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
