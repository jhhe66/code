package main

import (
	"github.com/gansidui/gotcp"
	"github.com/gansidui/gotcp/examples/tcpfw"
	"github.com/jessevdk/go-flags"
	"gopkg.in/ini.v1"
	"log"
	"net"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"
)

const ()

type Callback struct{}

var infoLog *log.Logger

func (this *Callback) OnConnect(c *gotcp.Conn) bool {
	addr := c.GetRawConn().RemoteAddr()
	c.PutExtraData(addr)
	infoLog.Println("OnConnect:", addr)
	return true
}

func (this *Callback) OnMessage(c *gotcp.Conn, p gotcp.Packet) bool {

	headV3Packet := p.(*tcpfw.HeadV3Packet)
	infoLog.Printf("OnMessage:[%v] [%v]\n", headV3Packet.GetLength(), string(headV3Packet.GetBody()))
	_, err = headV3Packet.CheckPacketValid()
	if err != nil {
		infoLog.Println("Invalid packet", err)
		return false
	}
	go ProcData(headV2Packet.GetBody())
	return true
}

func SendResp(c *gotcp.Conn, reqHead *tcpfw.HeadV2, ret uint16) bool {
	head := new(tcpfw.HeadV2)
	head.Version = reqHead.Version
	head.Cmd = reqHead.Cmd
	head.Seq = reqHead.Seq
	head.Uid = reqHead.Uid
	head.Len = tcpfw.EmptyPacktLen
	head.Ret = ret
	head.SysType = reqHead.SysType
	head.Echo[0] = reqHead.Echo[0]
	head.Echo[1] = reqHead.Echo[1]
	head.Echo[2] = reqHead.Echo[2]
	head.Echo[3] = reqHead.Echo[3]
	head.Echo[4] = reqHead.Echo[4]
	head.Echo[5] = reqHead.Echo[5]
	head.Echo[6] = reqHead.Echo[6]
	head.Echo[7] = reqHead.Echo[7]

	payLoad := []byte("acho hi")
	head.Len = uint32(tcpfw.PacketHeadLen) + uint32(len(payLoad)) + 1
	buf := make([]byte, head.Len)
	buf[0] = tcpfw.HTV2MagicBegin
	err := tcpfw.SerialToSlice(head, buf[1:])
	if err != nil {
		infoLog.Println("SerialHeadV2ToSlice failed")
		return false
	}
	copy(buf[tcpfw.PacketHeadLen:], payLoad)
	buf[head.Len-1] = tcpfw.HTV2MagicEnd
	resp := tcpfw.NewHeadV2Packet(buf)
	c.AsyncWritePacket(resp, time.Second)
	return true
}

func ProcData(bodyData []byte) bool {
	return true
}

func (this *Callback) OnClose(c *gotcp.Conn) {
	infoLog.Println("OnClose:", c.GetExtraData())
}

type Options struct {
	// Example of verbosity with level
	Verbose []bool `short:"v" long:"verbose" description:"Verbose output"`

	// Example of optional value
	ServerConf string `short:"c" long:"conf" description:"Server Config" optional:"no"`
}

var options Options

var parser = flags.NewParser(&options, flags.Default)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

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

	// creates a tcp listener
	serverIp := cfg.Section("LOCAL_SERVER").Key("bind_ip").MustString("127.0.0.1")
	serverPort := cfg.Section("LOCAL_SERVER").Key("bind_port").MustInt(8990)
	infoLog.Printf("serverIp=%v serverPort=%v", serverIp, serverPort)
	tcpAddr, err := net.ResolveTCPAddr("tcp4", serverIp+":"+strconv.Itoa(serverPort))
	checkError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	// creates a server
	sendChanLimit := cfg.Section("CHANLIMIT").Key("max_send_chan_count").MustUint(1000)
	recvChanLimit := cfg.Section("CHANLIMIT").Key("max_recv_chan_count").MustUint(1000)
	config := &gotcp.Config{
		PacketSendChanLimit:    uint32(sendChanLimit),
		PacketReceiveChanLimit: uint32(recvChanLimit),
	}
	srv := gotcp.NewServer(config, &Callback{}, &tcpfw.HeadV2Protocol{})

	// starts service
	go srv.Start(listener, time.Second)
	infoLog.Println("listening:", listener.Addr())

	// catchs system signal
	chSig := make(chan os.Signal)
	signal.Notify(chSig, syscall.SIGINT, syscall.SIGTERM)
	infoLog.Println("Signal: ", <-chSig)

	// stops service
	srv.Stop()
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
