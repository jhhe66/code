package main

import (
	// "fmt"
	// "github.com/bitly/go-simplejson"
	"log"
	"net"
	"os"
	"strconv"

	"github.com/gansidui/gotcp/tcpfw/common"
	"github.com/gansidui/gotcp/tcpfw/include/ht_muc"
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

	v3Protocol := &common.HeadV3Protocol{}
	var head *common.HeadV3
	head = &common.HeadV3{Flag: 0xF0,
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
	reqBody := new(ht_muc.MucReqBody)
	switch options.Cmd {
	// 创建群聊
	case 0x7041:
		head.Cmd = 0x7041

		memberInfo1 := &ht_muc.RoomMemberInfo{
			Uid: proto.Uint32(900866),
		}
		memberInfo2 := &ht_muc.RoomMemberInfo{
			Uid: proto.Uint32(109641),
		}

		memberList := []*ht_muc.RoomMemberInfo{
			memberInfo1,
			memberInfo2,
		}

		subReqBody := &ht_muc.CreateRoomReqBody{
			CreateUid: proto.Uint32(1946612),
			NickName:  []byte("songliwei"),
			Members:   memberList,
		}
		reqBody.CreateRoomReqbody = subReqBody

	// 获取添加到群聊列表中的roomId
	case 0x705F:
		head.Cmd = 0x705F
		subReqBody := &ht_muc.GetRoomFromContactListReqBody{
			OpUid: proto.Uint32(1946612),
		}
		reqBody.GetRoomFromContactListReqbody = subReqBody
	// 获取群的详细信息
	case 0x705B:
		head.Cmd = 0x705B
		subReqBody := &ht_muc.GetRoomInfoReqBody{
			OpUid:              proto.Uint32(4667342),
			RoomId:             proto.Uint32(34557),
			storeRoomTimeStamp: proto.Uint64(0),
		}
		reqBody.GetRoomInfoReqbody = subReqBody
	default:
		infoLog.Println("UnKnow input cmd =", options.Cmd)
	}

	payLoad, err = proto.Marshal(reqBody)
	if err != nil {
		infoLog.Printf("proto.Marshal failed from=%v to=%v cmd=%v seq=%v",
			head.From,
			head.To,
			head.Cmd,
			head.Seq)
		return
	}

	head.Len = uint32(common.PacketV3HeadLen + len(payLoad) + 1) //整个报文长度
	buf := make([]byte, head.Len)
	buf[0] = common.HTV3MagicBegin
	err = common.SerialHeadV3ToSlice(head, buf[1:])
	if err != nil {
		infoLog.Println("SerialHeadV3ToSlice failed")
		return
	}
	copy(buf[common.PacketV3HeadLen:], payLoad) // return code
	buf[head.Len-1] = common.HTV3MagicEnd

	infoLog.Printf("len=%v payLaod=%v\n", len(payLoad), payLoad)
	// write
	conn.Write(buf)
	// read
	p, err := v3Protocol.ReadPacket(conn)
	if err == nil {
		rspPacket := p.(*common.HeadV3Packet)
		rspHead, _ := rspPacket.GetHead()
		rspPayLoad := rspPacket.GetBody()
		infoLog.Printf("resp len=%v cmd=%v\n", rspHead.Len, rspHead.Cmd)
		rspBody := &ht_muc.MucRspBody{}
		err = proto.Unmarshal(rspPayLoad, rspBody)
		if err != nil {
			infoLog.Println("proto Unmarshal failed")
			return
		}
		switch rspHead.Cmd {
		case 0x7042:
			subRspBody := rspBody.GetCreateRoomRspbody()
			status := subRspBody.GetStatus()
			infoLog.Printf("CreateRoom rsp code=%v msg=%s roomId=%v ts=%v",
				status.GetCode(),
				status.GetReason(),
				subRspBody.GetRoomId(),
				subRspBody.GetRoomTimestamp())

		case 0x705F:
			subRspBody := rspBody.GetGetRoomFromContactListRspbody()
			status := subRspBody.GetStatus()
			infoLog.Printf("GetGetRoomFromContactListRspbody code=%v reason=%s", status.GetCode(), status.GetReason())
			listRoomInfo := subRspBody.GetListRoomInfo()
			for _, v := range listRoomInfo {
				infoLog.Printf("GetGetRoomFromContactListRspbody roomId=%v createUid=%v adminList=%v roomLimit=%v roomName=%s roomDesc=%s verifyStat=%v publishUid=%v publishTs=%v publishContext=%s roomTimeStamp=%v pushSetting=%v",
					v.GetRoomId(),
					v.GetCreateUid(),
					v.GetAdminLimit(),
					v.GetRoomLimit(),
					v.GetRoomName(),
					v.GetRoomDesc(),
					v.GetVerifyStat(),
					v.GetAnnouncement().GetPublishUid(),
					v.GetAnnouncement().GetPublishTs(),
					v.GetAnnouncement().GetAnnoContent(),
					v.GetRoomTimestamp(),
					v.GetPushSetting())
				for _, member := range v.GetMembers() {
					infoLog.Printf("GetGetRoomFromContactListRspbody member uid=%v name=%s", member.GetUid(), member.GetNickName())
				}

			}
		case 0x705C:
			subRspBody := rspBody.GetRoomInfoRspbody()
			status := subRspBody.GetStatus()
			roomInfo := subRspBody.GetRoomInf()
			infoLog.Printf("GetRoomInfoRspbody code=%v reason=%s", status.GetCode(), status.GetReason())
			infoLog.Printf("GetRoomInfoRespbody roomId=%v createUid=%v adminList=%v roomLimit=%v roomName=%s roomDesc=%s verifyStat=%v publishUid=%v publishTs=%v publishContext=%s roomTimeStamp=%v pushSetting=%v",
				roomInfo.GetRoomId(),
				roomInfo.GetCreateUid(),
				roomInfo.GetAdminLimit(),
				roomInfo.GetRoomLimit(),
				roomInfo.GetRoomName(),
				roomInfo.GetRoomDesc(),
				roomInfo.GetVerifyStat(),
				roomInfo.GetAnnouncement().GetPublishUid(),
				roomInfo.GetAnnouncement().GetPublishTs(),
				roomInfo.GetAnnouncement().GetAnnoContent(),
				roomInfo.GetRoomTimestamp(),
				roomInfo.GetPushSetting())
			for _, member := range roomInfo.GetMembers() {
				infoLog.Printf("GetRoomInfoRspbody member uid=%v name=%s", member.GetUid(), member.GetNickName())
			}

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
