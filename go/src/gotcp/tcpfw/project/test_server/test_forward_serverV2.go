package main

import (
	"github.com/gansidui/gotcp"
	"github.com/gansidui/gotcp/libcomm"
	"github.com/gansidui/gotcp/tcpfw/common"
	"github.com/gansidui/gotcp/tcpfw/include/ht_user"

	"log"
	"net"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/jessevdk/go-flags"
	"gopkg.in/ini.v1"
)

type Callback struct{}

var (
	infoLog   *log.Logger
	workerApi *common.P2PWorkerApiV2
)

// Convert uint to net.IP http://www.outofmemory.cn
func inet_ntoa(ipnr int64) net.IP {
	var bytes [4]byte
	bytes[0] = byte((ipnr >> 24) & 0xFF)
	bytes[1] = byte((ipnr >> 16) & 0xFF)
	bytes[2] = byte((ipnr >> 8) & 0xFF)
	bytes[3] = byte(ipnr & 0xFF)

	return net.IPv4(bytes[3], bytes[2], bytes[1], bytes[0])
}

// Convert net.IP to int64 ,  http://www.outofmemory.cn
func inet_aton(ipnr net.IP) int64 {
	bits := strings.Split(ipnr.String(), ".")

	b0, _ := strconv.Atoi(bits[0])
	b1, _ := strconv.Atoi(bits[1])
	b2, _ := strconv.Atoi(bits[2])
	b3, _ := strconv.Atoi(bits[3])

	var sum int64

	sum += int64(b0) << 24
	sum += int64(b1) << 16
	sum += int64(b2) << 8
	sum += int64(b3)

	return sum
}

func (this *Callback) OnConnect(c *gotcp.Conn) bool {
	addr := c.GetRawConn().RemoteAddr()
	c.PutExtraData(addr)
	infoLog.Println("OnConnect:", addr)
	return true
}

func (this *Callback) OnMessage(c *gotcp.Conn, p gotcp.Packet) bool {
	packet := p.(*common.HeadV3Packet)
	infoLog.Printf("OnMessage:[%v] [%v]", packet.GetLength(), string(packet.GetBody()))
	head, err := packet.GetHead()
	_, err = packet.CheckPacketValid()
	if err != nil {
		SendResp(c, head, uint16(ht_user.ResultCode_RET_INTERNAL_ERR))
		infoLog.Println("Invalid packet", err)
		return false
	}

	go ProcData(c, p, time.Now().UnixNano())
	return true
}

func SendResp(c *gotcp.Conn, reqHead *common.HeadV3, ret uint16) bool {
	head := new(common.HeadV3)
	if reqHead != nil {
		*head = *reqHead
	}

	head.Len = uint32(common.EmptyPacktV3Len)
	head.Ret = ret
	buf := make([]byte, head.Len)
	buf[0] = common.HTV3MagicBegin
	err := common.SerialHeadV3ToSlice(head, buf[1:])
	if err != nil {
		infoLog.Println("SerialHeadV3ToSlice failed")
		return false
	}

	buf[head.Len-1] = common.HTV3MagicEnd
	resp := common.NewHeadV3Packet(buf)
	c.AsyncWritePacket(resp, time.Second)
	return true
}

func SendRespWithPayLoad(c *gotcp.Conn, reqHead *common.HeadV3, payLoad []byte, ret uint16) bool {
	head := new(common.HeadV3)
	if reqHead != nil {
		*head = *reqHead
	}

	head.Len = uint32(common.PacketV3HeadLen) + uint32(len(payLoad)+1)
	head.Ret = ret
	buf := make([]byte, head.Len)
	buf[0] = common.HTV3MagicBegin
	err := common.SerialHeadV3ToSlice(head, buf[1:])
	if err != nil {
		infoLog.Println("SerialHeadV3ToSlice failed")
		return false
	}
	copy(buf[common.PacketV3HeadLen:], payLoad)
	buf[head.Len-1] = common.HTV3MagicEnd
	resp := common.NewHeadV3Packet(buf)
	c.AsyncWritePacket(resp, time.Second)
	return true
}

func ProcData(c *gotcp.Conn, p gotcp.Packet, nano int64) bool {
	attr := "Forward/proc_message"
	libcomm.AttrAdd(attr, 1)
	result := uint16(ht_user.ResultCode_RET_SUCCESS)
	var head *common.HeadV3
	var payLoad []byte
	defer SendRespWithPayLoad(c, head, payLoad, result)

	packet, ok := p.(*common.HeadV3Packet)
	if !ok {
		attr := "Forward/convert_v3packet_fail"
		libcomm.AttrAdd(attr, 1)
		infoLog.Println("Convert to HeadV3Packet failed")
		result = uint16(ht_user.ResultCode_RET_INTERNAL_ERR)
		return false
	}
	head, err := packet.GetHead()
	if err != nil {
		attr := "Forward/get_headv3_fail"
		libcomm.AttrAdd(attr, 1)
		infoLog.Println("ProcData Get head faild")
		result = uint16(ht_user.ResultCode_RET_INTERNAL_ERR)
		return false
	}
	// forward packet
	rsp, err := workerApi.SendAndRecvPacket(head, packet.GetBody())
	rspPacket := rsp.(*common.HeadV3Packet)
	rspHead, err := rspPacket.GetHead()
	if err != nil {
		result = uint16(ht_user.ResultCode_RET_INTERNAL_ERR)
		infoLog.Println("Forward Get packet head failed")
		return false
	}
	ret := rspHead.Ret
	payLoad = rspPacket.GetBody()
	if err != nil {
		result = uint16(ht_user.ResultCode_RET_INTERNAL_ERR)
		infoLog.Println("Forward Get packet body failed")
		return false
	}
	result = uint16(ht_user.ResultCode_RET_SUCCESS)
	curNano := time.Now().UnixNano()
	if (curNano-nano)/1000000 > 10 {
		attr := "Forward/slow_process"
		libcomm.AttrAdd(attr, 1)
		infoLog.Println("Forward proc slow")
	}

	infoLog.Printf("recv im ack ret=%v", ret)
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

	// 读取P2PWorker 配置
	workerIp := cfg.Section("P2PWORKER").Key("worker_ip").MustString("127.0.0.1")
	workerPort := cfg.Section("P2PWORKER").Key("worker_port").MustString("6")
	infoLog.Printf("p2p worker server ip=%v port=%v", workerIp, workerPort)
	workerApi = common.NewP2PWorkerApiV2(workerIp,
		workerPort,
		time.Minute,
		time.Minute,
		&common.HeadV3Protocol{},
		1000)

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
	srv := gotcp.NewServer(config, &Callback{}, &common.HeadV3Protocol{})

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
