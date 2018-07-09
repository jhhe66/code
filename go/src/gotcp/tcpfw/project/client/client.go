package main

import (
	"errors"
	"fmt"
	// "github.com/bitly/go-simplejson"
	"github.com/gansidui/gotcp/examples/ht_moment"
	"github.com/gansidui/gotcp/examples/tcpfw"
	"github.com/golang/protobuf/proto"
	"github.com/jessevdk/go-flags"
	"gopkg.in/ini.v1"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	// "time"
)

type TestUnit struct {
	From, To, Cmd uint32
	OpType        int32
}

func (p *TestUnit) UnmarshalFlag(value string) error {
	parts := strings.Split(value, ",")

	if len(parts) != 4 {
		return errors.New("expected two numbers separated by a ,")
	}

	from, err := strconv.ParseUint(parts[0], 10, 32)

	if err != nil {
		return err
	}

	to, err := strconv.ParseUint(parts[1], 10, 32)

	if err != nil {
		return err
	}

	cmd, err := strconv.ParseUint(parts[2], 10, 32)

	if err != nil {
		return err
	}

	op, err := strconv.ParseInt(parts[3], 10, 32)

	if err != nil {
		return err
	}

	p.From = uint32(from)
	p.To = uint32(to)
	p.Cmd = uint32(cmd)
	p.OpType = int32(op)

	return nil
}

func (p TestUnit) MarshalFlag() (string, error) {
	return fmt.Sprintf("%u,%u,%u,%d", p.From, p.To, p.Cmd, p.OpType), nil
}

type Options struct {
	// Example of verbosity with level
	Verbose []bool `short:"v" long:"verbose" description:"Verbose output"`

	// Example of custom type Marshal/Unmarshal
	TestUnit TestUnit `long:"unit" description:"A from,to,cmd,optype testunit" default:"1946612,2325927,0,1"`
}

var options Options

var parser = flags.NewParser(&options, flags.Default)

func main() {
	// 处理命令行参数
	if _, err := parser.Parse(); err != nil {
		log.Fatalln("parse cmd line failed!")
		os.Exit(1)
	}

	// 定义一个文件
	fileName := "clinet.log"
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
	infoLog.Printf("From=%v To=%v Cmd=%v Op=%v", options.TestUnit.From, options.TestUnit.To, options.TestUnit.Cmd, options.TestUnit.OpType)

	tcpAddr, err := net.ResolveTCPAddr("tcp4", serverIp+":"+strconv.Itoa(serverPort))
	checkError(err)
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err)

	headV2Protocol := &tcpfw.HeadV2Protocol{}

	var head *tcpfw.HeadV2
	var reqBody ht_moment.ReqBody
	uid := options.TestUnit.To
	opType := ht_moment.OPERATOR_TYPE(options.TestUnit.OpType)
	switch options.TestUnit.Cmd {
	case 20:
		head = &tcpfw.HeadV2{Version: 0x04, Cmd: options.TestUnit.Cmd, Uid: options.TestUnit.From}
		reqSubBody := &ht_moment.DoNotShareMyMomentReqBody{UserId: &uid, Type: &opType}
		reqBody.DoNotShareMyMntReqbody = reqSubBody
	case 21:
		head = &tcpfw.HeadV2{Version: 0x04, Cmd: options.TestUnit.Cmd, Uid: options.TestUnit.From}

	case 22:
		head = &tcpfw.HeadV2{Version: 0x04, Cmd: options.TestUnit.Cmd, Uid: options.TestUnit.From}
		reqSubBody := &ht_moment.HideOtherMomentReqBody{UserId: &uid, Type: &opType}
		reqBody.HideOtherMntReqbody = reqSubBody
	case 23:
		head = &tcpfw.HeadV2{Version: 0x04, Cmd: options.TestUnit.Cmd, Uid: options.TestUnit.From}

	default:
		infoLog.Println("Unhandle comand=%v", options.TestUnit.Cmd)
		os.Exit(1)
	}

	payLoad, err := proto.Marshal(&reqBody)
	if err != nil {
		infoLog.Println("marshaling error: ", err)
	}

	head.Len = uint32(tcpfw.PacketHeadLen) + uint32(len(payLoad)) + 1
	buf := make([]byte, head.Len)
	buf[0] = tcpfw.HTV2MagicBegin
	err = tcpfw.SerialToSlice(head, buf[1:])
	if err != nil {
		infoLog.Println("SerialHeadV2ToSlice failed")
		os.Exit(1)
	}
	copy(buf[tcpfw.PacketHeadLen:], payLoad)
	buf[head.Len-1] = tcpfw.HTV2MagicEnd

	infoLog.Println("payLoad=", tcpfw.NewHeadV2Packet(buf).Serialize())
	// write
	conn.Write(tcpfw.NewHeadV2Packet(buf).Serialize())

	// read
	p, err := headV2Protocol.ReadPacket(conn)
	if err == nil {
		headV2Packet := p.(*tcpfw.HeadV2Packet)
		// infoLog.Printf("Server reply:[%v] [%v]\n", headV2Packet.GetLength(), string(headV2Packet.GetBody()))
		var rspBody ht_moment.RspBody
		err = proto.Unmarshal(headV2Packet.GetBody(), &rspBody)
		if err != nil {
			infoLog.Println("unmarshaling error: ", err)
		}
		switch headV2Packet.GetCommand() {
		case 20:
			subRspBody := rspBody.GetDoNotShareMyMntRspbody()
			status := subRspBody.GetStatus()
			infoLog.Printf("code=%v reason=%s", status.GetCode(), string(status.GetReason()))
		case 21:
			subRspBody := rspBody.GetGetNotShareToListRspbody()
			status := subRspBody.GetStatus()
			infoLog.Printf("code=%v reason=%s", status.GetCode(), string(status.GetReason()))
			uidArray := subRspBody.GetToUidArray()
			for k, v := range uidArray {
				infoLog.Printf("key=%d value=[%v]", k, v)
			}

		case 22:
			subRspBody := rspBody.GetHideOtherMntRspbody()
			status := subRspBody.GetStatus()
			infoLog.Printf("code=%v reason=%s", status.GetCode(), string(status.GetReason()))
		case 23:
			subRspBody := rspBody.GetGetIHideListRspbody()
			status := subRspBody.GetStatus()
			infoLog.Printf("code=%v reason=%s", status.GetCode(), string(status.GetReason()))
			uidArray := subRspBody.GetToUidArray()
			for k, v := range uidArray {
				infoLog.Printf("key=%d value=[%v]", k, v)
			}
		default:
			infoLog.Println("Unhand resp cmd=", headV2Packet.GetCommand())
		}
	}

	conn.Close()
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
