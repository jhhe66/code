package main

import (
	"github.com/bitly/go-simplejson"
	"github.com/gansidui/gotcp"
	"github.com/gansidui/gotcp/examples/tcpfw"
	"github.com/garyburd/redigo/redis"
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

func newPool(redisServer string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", redisServer)
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
}

var (
	pool *redis.Pool
)

const (
	// 业务流程
	CMD_LOGIN       = "login"
	CMD_FIRST_LOGIN = "first_login"
	CMD_FORCE_LOGIN = "force_login"
	CMD_RECONNECT   = "reconnect"
	CMD_VOIP        = "voip"
	CMD_VIDEO       = "video"
	CMD_WNS_CONNECT = "wns_connect"
	// IM请求
	CMD_0X1001        = "0x1001"
	CMD_0X1003        = "0x1003"
	CMD_0X1005        = "0x1005"
	CMD_0X4019        = "0x4019"
	CMD_0X401D        = "0x401d"
	CMD_0X401F        = "0x401f"
	CMD_0X4021        = "0x4021"
	CMD_0X4023        = "0x4023"
	CMD_0X411D        = "0x411d"
	CMD_WNS_CONN_STAT = "wns_connect_status"
	CMD_CONN_VOIP_SDK = "conn_voip_sdk"
	//客户端事件
	CMD_POPBUY   = "pop_buy"
	CMD_CLICKBUY = "click_buy"

	LOGSTASH_QUEUE = "client_log"
	LOGHUB_VERSION = "2"
)

type Callback struct{}

var infoLog *log.Logger

func (this *Callback) OnConnect(c *gotcp.Conn) bool {
	addr := c.GetRawConn().RemoteAddr()
	c.PutExtraData(addr)
	infoLog.Println("OnConnect:", addr)
	return true
}

func (this *Callback) OnMessage(c *gotcp.Conn, p gotcp.Packet) bool {

	headV2Packet := p.(*tcpfw.HeadV2Packet)
	// infoLog.Printf("OnMessage:[%v] [%v]\n", headV2Packet.GetLength(), string(headV2Packet.GetBody()))
	head, err := headV2Packet.GetHead()
	// infoLog.Println("OnMessage head =", *head)
	SendResp(c, head, 0)
	_, err = headV2Packet.CheckPacketV2Valid()
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
	head.Len = tcpfw.EmptyPacktV2Len
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
	head.Len = uint32(tcpfw.PacketV2HeadLen) + uint32(len(payLoad)) + 1
	buf := make([]byte, head.Len)
	buf[0] = tcpfw.HTV2MagicBegin
	err := tcpfw.SerialHeadV2ToSlice(head, buf[1:])
	if err != nil {
		infoLog.Println("SerialHeadV2ToSlice failed")
		return false
	}
	copy(buf[tcpfw.PacketV2HeadLen:], payLoad)
	buf[head.Len-1] = tcpfw.HTV2MagicEnd
	resp := tcpfw.NewHeadV2Packet(buf)
	c.AsyncWritePacket(resp, time.Second)
	return true
}

func ProcData(bodyData []byte) bool {
	// 获取一条Redis连接
	redisConn := pool.Get()
	defer redisConn.Close()

	rootObj, err := simplejson.NewJson(bodyData)
	if err != nil {
		infoLog.Println("simplejson new packet error", err)
		return false
	}

	var uid string
	uid, err = rootObj.Get("userid").String()
	if err != nil {
		iUid := rootObj.Get("userid").MustInt()
		uid = strconv.Itoa(iUid)
	}
	terminal, err := rootObj.Get("terminaltype").String()
	if err != nil {
		iTerminal := rootObj.Get("terminaltype").MustInt()
		terminal = strconv.Itoa(iTerminal)
	}

	version := rootObj.Get("version").MustString()
	clientIp := rootObj.Get("client_ip").MustString()
	v := rootObj.Get("data").MustArray()
	count := len(v)
	for i := 0; i < count; i++ {
		vv := rootObj.Get("data").GetIndex(i).MustArray()
		vvLen := len(vv)
		for j := 0; j < vvLen; j++ {
			subValue := rootObj.Get("data").GetIndex(i).GetIndex(j)
			var sPrevCmd, sPrevTs, sPrevCost, sPrevRet, sPrevRetry string
			outObj := simplejson.New()
			cmd := subValue.Get("cmd").MustString()
			if (cmd == CMD_LOGIN) ||
				(cmd == CMD_FORCE_LOGIN) ||
				(cmd == CMD_FIRST_LOGIN) ||
				(cmd == CMD_RECONNECT) ||
				(cmd == CMD_VOIP) ||
				(cmd == CMD_VIDEO) ||
				(cmd == CMD_WNS_CONNECT) {
				outObj.Set("ret", subValue.Get("ret").MustString())
				outObj.Set("cost", subValue.Get("cost").MustString())
				outObj.Set("ts", subValue.Get("ts").MustString())
				outObj.Set("cmd", subValue.Get("cmd").MustString())
				outObj.Set("nt", subValue.Get("nt").MustString())
			} else if (cmd == CMD_0X1001) ||
				(cmd == CMD_0X1003) ||
				(cmd == CMD_0X1005) ||
				(cmd == CMD_WNS_CONN_STAT) {
				outObj.Set("ret", subValue.Get("ret").MustString())
				outObj.Set("cost", subValue.Get("cost").MustString())
				outObj.Set("ts", subValue.Get("ts").MustString())
				outObj.Set("cmd", subValue.Get("cmd").MustString())
				outObj.Set("nt", subValue.Get("nt").MustString())
				outObj.Set("retry", subValue.Get("retry").MustString())
				outObj.Set("wns_code", subValue.Get("wns_code").MustString())
				outObj.Set("wns_ocip", subValue.Get("wns_ocip").MustString())
			} else if (cmd == CMD_0X4019) ||
				(cmd == CMD_0X401D) ||
				(cmd == CMD_0X401F) ||
				(cmd == CMD_0X4021) ||
				(cmd == CMD_0X4023) ||
				(cmd == CMD_0X411D) ||
				(cmd == CMD_CONN_VOIP_SDK) {
				outObj.Set("ret", subValue.Get("ret").MustString())
				outObj.Set("cost", subValue.Get("cost").MustString())
				outObj.Set("ts", subValue.Get("ts").MustString())
				outObj.Set("cmd", subValue.Get("cmd").MustString())
				outObj.Set("nt", subValue.Get("nt").MustString())
				outObj.Set("retry", subValue.Get("retry").MustString())
				outObj.Set("wns_code", subValue.Get("wns_code").MustString())
				outObj.Set("wns_ocip", subValue.Get("wns_ocip").MustString())
				outObj.Set("roomid", subValue.Get("roomid").MustString())
			} else if (cmd == CMD_POPBUY) || (cmd == CMD_CLICKBUY) {
				outObj.Set("ts", subValue.Get("ts").MustString())
				outObj.Set("cmd", subValue.Get("cmd").MustString())
				outObj.Set("pos", subValue.Get("pos").MustString())
			} else {
				infoLog.Println("Unknow cmd=", cmd)
				continue
			}

			msec := transferMsec(outObj.Get("ts").MustString())
			outObj.Set("userid", uid)
			outObj.Set("version", version)
			outObj.Set("terminaltype", terminal)
			outObj.Set("client_ip", clientIp)
			outObj.Set("@version", LOGHUB_VERSION)
			outObj.Set("@timestamp", msec)

			sPacket, err := outObj.MarshalJSON()
			if err != nil {
				infoLog.Println("MarshalJSON failed outObj=", outObj)
				continue
			}

			r, err := redisConn.Do("RPUSH", LOGSTASH_QUEUE, sPacket)
			if err != nil {
				infoLog.Println("redis exec RPUSH failed")
				continue
			} else {
				infoLog.Printf("redis exec RPUSH ret=%v", r)
			}

			// infoLog.Printf("Uid=%v Record=[j:%d:%s]", uid, j, sPacket)
			/*生成一条记录, 只有时间，类型，耗时和前一条不同*/
			if j >= 2 {
				newRecord := simplejson.New()
				if cmd != sPrevCmd {
					sNewCmd := sPrevCmd + "-" + cmd
					iPrevTs, err := strconv.ParseUint(sPrevTs, 10, 64)
					if err != nil {
						infoLog.Printf("strconv.ParseUint failed prevTs=%s", sPrevTs)
						continue
					}

					iPrevCost, err := strconv.ParseUint(sPrevCost, 10, 64)
					if err != nil {
						infoLog.Printf("strconv.ParseUint failed sPrevCost=%s", sPrevCost)
						continue
					}

					iCurTs, err := strconv.ParseUint(subValue.Get("ts").MustString(), 10, 64)
					if err != nil {
						infoLog.Printf("strconv.ParseUint failed curTs=%s", subValue.Get("ts").MustString())
						continue
					}

					sNewTs := strconv.FormatUint(iPrevTs+iPrevCost, 10)
					sCurCost := strconv.FormatUint(iCurTs-(iPrevTs+iPrevCost), 10)
					newRecord.Set("cmd", sNewCmd)
					newRecord.Set("ts", sNewTs)
					newRecord.Set("cost", sCurCost)

					sNewPacket, err := newRecord.MarshalJSON()
					if err != nil {
						infoLog.Println("MarshalJSON failed newRecord=", newRecord)
						continue
					}

					rNew, err := redisConn.Do("RPUSH", LOGSTASH_QUEUE, sNewPacket)
					if err != nil {
						infoLog.Println("redis exec RPUSH new packet failed")
						continue
					} else {
						infoLog.Printf("redis exec RPUSH new packet ret=%v", rNew)
					}
					infoLog.Printf("Uid|%v|Record_add[j:%d:%s]", uid, j, sNewPacket)
				} else {
					if strings.EqualFold(subValue.Get("ret").MustString(), "0") && strings.EqualFold(sPrevRet, "0") {
						infoLog.Printf("Uid|%v|Fatal erorr!Please check the ret field!", uid)
					}
					if strings.EqualFold(subValue.Get("retry").MustString(), sPrevRetry) {
						infoLog.Printf("Uid|%v|Fatal erorr!Please check the retry field!", uid)
					}
					curRetry, _ := strconv.Atoi(subValue.Get("retry").MustString())
					prevRetry, _ := strconv.Atoi(sPrevRetry)
					if curRetry+1 < prevRetry {
						infoLog.Printf("Uid|%v|Fatal erorr!Please check the retry field!front[%s], rear[%s]", uid, subValue.Get("retry").MustString(), sPrevRetry)
					}
				}
			}

			sPrevCmd = cmd
			sPrevTs = subValue.Get("ts").MustString()
			sPrevCost = subValue.Get("cost").MustString()
			sPrevRet = subValue.Get("ret").MustString()
			sPrevRetry = subValue.Get("retry").MustString()
		}
	}
	return true
}

func transferMsec(timeStr string) (str string) {
	u, err := strconv.ParseUint(timeStr, 10, 64)
	if err != nil {
		infoLog.Printf("transferMsec strconv failed time=%s", timeStr)
		return str
	}

	msec := u - (u/1000)*1000
	nsec := msec * 1000 * 1000
	sec := (u - msec) / 1000
	return time.Unix(int64(sec), int64(nsec)).Format(time.RFC3339)
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

	// init redis pool
	redisIp := cfg.Section("REDIS").Key("redis_ip").MustString("127.0.0.1")
	redisPort := cfg.Section("REDIS").Key("redis_port").MustInt(6379)
	infoLog.Printf("redis ip=%v port=%v", redisIp, redisPort)
	pool = newPool(redisIp + ":" + strconv.Itoa(redisPort))

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
