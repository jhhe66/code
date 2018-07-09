package main

import (
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/gansidui/gotcp/examples/ht_moment"
	"github.com/gansidui/gotcp/examples/tcpfw"
	"github.com/golang/protobuf/proto"
	"github.com/jessevdk/go-flags"
	"gopkg.in/ini.v1"
)

type Options struct {
	// Example of verbosity with level
	Verbose []bool `short:"v" long:"verbose" description:"Verbose output"`

	ClientConf string `short:"c" long:"config" description:"config file" optional:"no"`
	Op         uint32 `short:"o" long:"operator" description:"operator add or del"`
	List       uint32 `short:"l" logn:"list" description:"list type"`
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

	// 读取ip+port
	serverIp := cfg.Section("OUTER_SERVER").Key("server_ip").MustString("127.0.0.3")
	serverPort := cfg.Section("OUTER_SERVER").Key("server_port").MustInt(8990)

	infoLog.Printf("server_ip=%v server_port=%v\n", serverIp, serverPort)

	tcpAddr, err := net.ResolveTCPAddr("tcp4", serverIp+":"+strconv.Itoa(serverPort))
	checkError(err)
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err)

	uid := cfg.Section("TEST_UID").Key("UID").MustInt(3249312)
	headV2Protocol := &tcpfw.HeadV2Protocol{}
	reqBody := new(ht_moment.ReqBody)
	head := &tcpfw.HeadV2{Cmd: uint32(ht_moment.CMD_TYPE_CMD_OPERATOR_UID),
		Len: uint32(tcpfw.EmptyPacktV2Len),
		Uid: uint32(uid)}
	beginUid := uint32(2325928)
	for i := 0; i < 1; i++ {
		head.Cmd = uint32(ht_moment.CMD_TYPE_CMD_OPERATOR_UID)
		subReqBody := new(ht_moment.OperatorUidReqBody)
		subReqBody.ListType = new(ht_moment.OPERATORLIST)
		*(subReqBody.ListType) = ht_moment.OPERATORLIST(options.List) //被隐藏帖子的用户列表
		subReqBody.OpType = new(ht_moment.OPERATORTYPE)
		*(subReqBody.OpType) = ht_moment.OPERATORTYPE(options.Op) //添加操作
		subReqBody.OpUserId = new(uint32)
		*(subReqBody.OpUserId) = beginUid
		subReqBody.Userid = new(uint32)
		*(subReqBody.Userid) = uint32(uid)
		reqBody.OpUidReqbody = subReqBody
		infoLog.Printf("uid=%d", beginUid)

		payLoad, err := proto.Marshal(reqBody)
		if err != nil {
			infoLog.Println("marshaling error: ", err)
		}

		head.Len = uint32(tcpfw.PacketV2HeadLen) + uint32(len(payLoad)) + 1
		buf := make([]byte, head.Len)
		buf[0] = tcpfw.HTV2MagicBegin
		err = tcpfw.SerialHeadV2ToSlice(head, buf[1:])
		if err != nil {
			infoLog.Println("SerialHeadV3ToSlice failed")
			os.Exit(1)
		}
		copy(buf[tcpfw.PacketV2HeadLen:], payLoad)
		buf[head.Len-1] = tcpfw.HTV2MagicEnd

		infoLog.Println("payLoad=", tcpfw.NewHeadV2Packet(buf).Serialize())
		// write
		conn.Write(tcpfw.NewHeadV2Packet(buf).Serialize())

		// read
		p, err := headV2Protocol.ReadPacket(conn)
		if err == nil {
			packet := p.(*tcpfw.HeadV2Packet)
			// infoLog.Printf("Server reply:[%v] [%v]\n", headV2Packet.GetLength(), string(headV2Packet.GetBody()))
			rspHead, err := packet.GetHead()
			if err == nil {
				infoLog.Printf("head=%+v", rspHead)
				rspBody := &ht_moment.RspBody{}
				err = proto.Unmarshal(packet.GetBody(), rspBody)
				if err != nil {
					infoLog.Println("proto Unmarshal failed")
					return
				}
				switch rspHead.Cmd {
				case uint32(ht_moment.CMD_TYPE_CMD_OPERATOR_UID):
					{
						subRspBody := rspBody.GetOpUidRspbody()
						infoLog.Printf("stat=%v list_type=%v\n", subRspBody.GetStatus().GetCode(), subRspBody.GetListType())
					}
				default:
					{
						infoLog.Println("Unknow cmd=", rspHead.Cmd)
					}
				}
			} else {
				infoLog.Println("get HTV2Head failed err =", err)
			}
		}
		beginUid++
	}
	return
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
