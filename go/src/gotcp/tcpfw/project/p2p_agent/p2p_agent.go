package main

import (
	"github.com/gansidui/gotcp"
	"github.com/gansidui/gotcp/examples/ht_p2p"
	"github.com/gansidui/gotcp/examples/ht_push"
	"github.com/gansidui/gotcp/examples/tcpfw"

	"log"
	"net"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/jessevdk/go-flags"
	"gopkg.in/ini.v1"
)

type Callback struct{}

var (
	infoLog    *log.Logger
	mcApi      *tcpfw.MemcacheApi
	imServer   map[string]*tcpfw.ImServerApi
	workerApi  *tcpfw.P2PWorkerApi
	offlineApi *tcpfw.OfflineApiV2
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
	packet := p.(*tcpfw.HeadV3Packet)
	infoLog.Printf("OnMessage:[%v] [%v]", packet.GetLength(), string(packet.GetBody()))
	head, err := packet.GetHead()
	_, err = packet.CheckPacketValid()
	if err != nil {
		SendResp(c, head, uint16(ht_p2p.RET_CODE_RET_INPUT_PARAM_ERR))
		infoLog.Println("Invalid packet", err)
		return false
	}

	go ProcData(c, p)
	return true
}

func SendResp(c *gotcp.Conn, reqHead *tcpfw.HeadV3, ret uint16) bool {
	head := new(tcpfw.HeadV3)
	if reqHead != nil {
		*head = *reqHead
	}

	head.Len = uint32(tcpfw.EmptyPacktV3Len)
	head.Ret = ret
	buf := make([]byte, head.Len)
	buf[0] = tcpfw.HTV3MagicBegin
	err := tcpfw.SerialHeadV3ToSlice(head, buf[1:])
	if err != nil {
		infoLog.Println("SerialHeadV3ToSlice failed")
		return false
	}

	buf[head.Len-1] = tcpfw.HTV3MagicEnd
	resp := tcpfw.NewHeadV3Packet(buf)
	c.AsyncWritePacket(resp, time.Second)
	return true
}

func ProcData(c *gotcp.Conn, p gotcp.Packet) bool {
	result := uint16(ht_p2p.RET_CODE_RET_SUCCESS)
	var head *tcpfw.HeadV3
	defer SendResp(c, head, result)

	packet, ok := p.(*tcpfw.HeadV3Packet)
	if !ok {
		infoLog.Println("Convert to HeadV3Packet failed")
		result = uint16(ht_p2p.RET_CODE_RET_INTERNAL_ERR)
		return false
	}
	head, err := packet.GetHead()
	if err != nil {
		infoLog.Println("ProcData Get head faild")
		result = uint16(ht_p2p.RET_CODE_RET_INTERNAL_ERR)
		return false
	}

	reqBody := &ht_p2p.P2PMsgBody{}
	err = proto.Unmarshal(packet.GetBody(), reqBody)
	if err != nil {
		infoLog.Println("proto Unmarshal failed")
		result = uint16(ht_p2p.RET_CODE_RET_PB_ERR)
		return false
	}

	stat, err := mcApi.GetUserOnlineStat(head.To)
	if err != nil {
		infoLog.Println("Get user stat failed uid =", head.To)
		result = uint16(ht_p2p.RET_CODE_RET_INTERNAL_ERR)
		return false
	}

	infoLog.Printf("User=%v uid=%v clientType=%v stat=%v SvrIp=%v updateTs=%v version=%v port=%v wid=%v usertype=%v",
		head.To,
		stat.Uid,
		stat.ClientType,
		stat.OnlineStat,
		stat.SvrIp,
		stat.UpdateTs,
		stat.Version,
		stat.Port,
		stat.Wid,
		stat.UserType)

	if stat.OnlineStat == tcpfw.ST_ONLINE {
		svrIp := inet_ntoa(int64(stat.SvrIp)).String()
		infoLog.Printf("svrIP=%s oriIP=%v", svrIp, stat.SvrIp)
		v, ok := imServer[svrIp]
		if !ok { // 不存在直接打印日志返回错误
			result = uint16(ht_p2p.RET_CODE_RET_INTERNAL_ERR)
			infoLog.Printf("not exist IM ip=%v imServer=%v", svrIp, imServer)
			return false
		}
		ret, err := v.SendPacket(head, reqBody.P2PData)
		if err != nil {
			result = uint16(ht_p2p.RET_CODE_RET_INTERNAL_ERR)
			infoLog.Printf("Send packet to IM failed ret=%v err=%s", ret, err)
			return false
		}
		result = uint16(ht_p2p.RET_CODE_RET_SUCCESS)
		infoLog.Printf("recv im ack ret=%v", ret)
	} else { // 不在线首先判断是否丢弃
		justOnline := reqBody.GetJustOnline()
		if justOnline == 1 { // 丢弃
			result = uint16(ht_p2p.RET_CODE_RET_SUCCESS)
			infoLog.Printf("Drop p2p packet cmd=%v from=%v to=%v", head.Cmd, head.From, head.To)
		} else {
			if reqBody.GetPushInfo() != nil { // push info 不为空
				push, err := buildPushContent(head, reqBody.GetPushInfo())
				if err != nil {
					result = uint16(ht_p2p.RET_CODE_RET_INTERNAL_ERR)
					infoLog.Println("build Push Contect faield err =", err)
					return false
				} else {
					ret, err := offlineApi.SendPacketWithHeadV3(head, reqBody.GetP2PData(), push)
					if err != nil {
						result = uint16(ht_p2p.RET_CODE_RET_INTERNAL_ERR)
						infoLog.Println("save offline faield [ret err] =", ret, err)
					} else {
						result = ret
						infoLog.Println("save offline with push succe ret", ret)
					}
				}
			} else { // push info 为空
				ret, err := offlineApi.SendPacketWithHeadV3(head, reqBody.GetP2PData(), nil)
				if err != nil {
					result = uint16(ht_p2p.RET_CODE_RET_INTERNAL_ERR)
					infoLog.Println("save offline faield [ret err] =", ret, err)
				} else {
					result = ret
					infoLog.Println("save offline success ret =", ret)
				}
			}
		}
	}
	return true
}

func buildPushContent(head *tcpfw.HeadV3, pushInfo *ht_p2p.PushInfo) (push []byte, err error) {
	newHead := &tcpfw.HeadV3{From: head.From, To: head.To, Cmd: uint16(ht_push.CMD_TYPE_CMD_S2S_MESSAGE_PUSH), SysType: uint16(ht_push.SYS_TYPE_SYS_VOIP_SERVER)}
	reqBody := new(ht_push.ReqBody)
	reqBody.PushMsgReqbody = new(ht_push.PushMsgReqBody)
	subReqBody := reqBody.GetPushMsgReqbody()
	subReqBody.ChatType = new(uint32)
	*(subReqBody.ChatType) = 0

	subReqBody.RoomId = new(uint32)
	*(subReqBody.RoomId) = 0

	subReqBody.PushType = new(uint32)
	*(subReqBody.PushType) = pushInfo.GetPushType()

	subReqBody.Nickname = make([]byte, len(pushInfo.GetNickName()))
	copy(subReqBody.Nickname, pushInfo.GetNickName())

	subReqBody.Content = make([]byte, len(pushInfo.GetContent()))
	copy(subReqBody.Content, pushInfo.GetContent())

	subReqBody.Sound = new(uint32)
	*(subReqBody.Sound) = 1

	subReqBody.Light = new(uint32)
	*(subReqBody.Light) = 1

	strTime := time.Now().String()
	subReqBody.MessageId = make([]byte, len(strTime))
	copy(subReqBody.MessageId, []byte(strTime))

	s, err := proto.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	newHead.Len = uint32(tcpfw.PacketV3HeadLen) + uint32(len(s)) + 1
	push = make([]byte, newHead.Len)
	push[0] = tcpfw.HTV3MagicBegin
	err = tcpfw.SerialHeadV3ToSlice(newHead, push[1:])
	if err != nil {
		return nil, err
	}
	copy(push[tcpfw.PacketV3HeadLen:], s)
	push[newHead.Len-1] = tcpfw.HTV3MagicEnd
	return push, nil
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

	// 读取 Memcache Ip and port
	mcIp := cfg.Section("MEMCACHE").Key("mc_ip").MustString("127.0.0.1")
	mcPort := cfg.Section("MEMCACHE").Key("mc_port").MustInt(11211)
	infoLog.Printf("memcache ip=%v port=%v", mcIp, mcPort)
	mcApi = new(tcpfw.MemcacheApi)
	mcApi.Init(mcIp + ":" + strconv.Itoa(mcPort))

	// 读取P2PWorker 配置
	workerIp := cfg.Section("P2PWORKER").Key("worker_ip").MustString("127.0.0.1")
	workerPort := cfg.Section("P2PWORKER").Key("worker_port").MustInt(0)
	infoLog.Printf("p2p worker server ip=%v port=%v", workerIp, workerPort)
	workerApi = tcpfw.NewP2PWorkerApi(workerIp, workerPort, &tcpfw.HeadV3Protocol{})

	// 读取offline 配置
	offlineIp := cfg.Section("OFFLINE").Key("offline_ip").MustString("127.0.0.1")
	offlinePort := cfg.Section("OFFLINE").Key("offline_port").MustString("0")
	infoLog.Printf("offline server ip=%v port=%v", offlineIp, offlinePort)
	offlineApi = tcpfw.NewOfflineApiV2(offlineIp, offlinePort, time.Minute, time.Minute, &tcpfw.HeadV2Protocol{}, 1000)

	// 读取IMServer 配置
	imCount := cfg.Section("IMSERVER").Key("imserver_cnt").MustInt(0)
	imServer = make(map[string]*tcpfw.ImServerApi, imCount)
	infoLog.Printf("IMServer Count=%v", imCount)
	for i := 0; i < imCount; i++ {
		ipKey := "imserver_ip_" + strconv.Itoa(i)
		ipOnlineKye := "imserver_ip_online_" + strconv.Itoa(i)
		portKey := "imserver_port_" + strconv.Itoa(i)
		imIp := cfg.Section("IMSERVER").Key(ipKey).MustString("127.0.0.1")
		imIpOnline := cfg.Section("IMSERVER").Key(ipOnlineKye).MustString("127.0.0.1")
		imPort := cfg.Section("IMSERVER").Key(portKey).MustInt(0)
		infoLog.Printf("im server ip=%v ip_online=%v port=%v", imIp, imIpOnline, imPort)
		imServer[imIpOnline] = tcpfw.NewImServerApi(imIp, imPort, &tcpfw.HeadV3Protocol{})
	}

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
	srv := gotcp.NewServer(config, &Callback{}, &tcpfw.HeadV3Protocol{})

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
