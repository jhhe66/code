package main

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"

	simplejson "github.com/bitly/go-simplejson"
	"github.com/gansidui/gotcp"
	"github.com/gansidui/gotcp/libcomm"
	"github.com/gansidui/gotcp/tcpfw/common"
	"github.com/jessevdk/go-flags"
	"gopkg.in/ini.v1"
)

type Callback struct{}

var (
	infoLog *log.Logger
)

const (
	RET_CODE_RET_SUCCESS          = 0
	RET_CODE_RET_PB_ERR           = 500
	RET_CODE_RET_INTERNAL_ERR     = 501
	RET_CODE_RET_SESS_TIMEOUT_ERR = 502
	RET_CODE_RET_INPUT_PARAM_ERR  = 503
)

const (
	RESTAPI_PUSHSINGLEDEVICE  = "http://openapi.xg.qq.com/v2/push/single_device"
	RESTAPI_PUSHSINGLEACCOUNT = "http://openapi.xg.qq.com/v2/push/single_account"
	RESTAPI_PUSHACCOUNTLIST   = "http://openapi.xg.qq.com/v2/push/account_list"
	RESTAPI_PUSHALLDEVICE     = "http://openapi.xg.qq.com/v2/push/all_device"
	RESTAPI_PUSHTAGS          = "http://openapi.xg.qq.com/v2/push/tags_device"
	RESTAPI_QUERYPUSHSTATUS   = "http://openapi.xg.qq.com/v2/push/get_msg_status"
	RESTAPI_QUERYDEVICECOUNT  = "http://openapi.xg.qq.com/v2/application/get_app_device_num"
	RESTAPI_QUERYTAGS         = "http://openapi.xg.qq.com/v2/tags/query_app_tags"
	RESTAPI_CANCELTIMINGPUSH  = "http://openapi.xg.qq.com/v2/push/cancel_timing_task"
	RESTAPI_BATCHSETTAG       = "http://openapi.xg.qq.com/v2/tags/batch_set"
	RESTAPI_BATCHDELTAG       = "http://openapi.xg.qq.com/v2/tags/batch_del"
	RESTAPI_QUERYTOKENTAGS    = "http://openapi.xg.qq.com/v2/tags/query_token_tags"
	RESTAPI_QUERYTAGTOKENNUM  = "http://openapi.xg.qq.com/v2/tags/query_tag_token_num"
)

const (
	METHOD_GET  = "get"  //get请求方式
	METHOD_POST = "post" //post请求方式
)

const (
	TYPE_ACTIVITY = 1
	TYPE_URL      = 2
	TYPE_INTENT   = 3
	TYPE_PACKAGE  = 4
)
const (
	TYPE_NOTIFICATION = 1
	TYPE_MESSAGE      = 2
)

const (
	cAccessId  = 2100093818
	cSecretKey = "42deaa73626e172176fbe6a27e730de5"
)

const (
	P2P_CHAT = "CHAT"
	MUC_CHAT = "MUC"
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

type ActionAttribute struct {
	sendIf int
	sendPf int
}

func NewActionAttribute(inputIf, inputPf int) (result *ActionAttribute) {
	result = new(ActionAttribute)
	result.sendIf = inputIf
	result.sendPf = inputPf
	return result
}

func (this *ActionAttribute) SerializeAttrToJson() (outJson *simplejson.Json) {
	outJson = simplejson.New()
	outJson.Set("if", this.sendIf)
	outJson.Set("pf", this.sendPf)
	return outJson
}

type Browser struct {
	sendUrl     string
	sendConfirm int
}

func NewBrowser(inputUrl string, inputConfim int) (result *Browser) {
	result = new(Browser)
	result.sendUrl = inputUrl
	result.sendConfirm = inputConfim
	return result
}

func (this *Browser) SerializeBrowserToJson() (outJson *simplejson.Json) {
	outJson = simplejson.New()
	outJson.Set("url", this.sendUrl)
	outJson.Set("confirm", this.sendConfirm)
	return outJson
}

type PacketName struct {
	packageDownloadUrl string
	sendConfirm        int
	name               string
}

func NewPacketName(inputDownloadUrl string, inputConfim int, inputName string) (result *PacketName) {
	result = new(PacketName)
	result.packageDownloadUrl = inputDownloadUrl
	result.sendConfirm = inputConfim
	result.name = inputName
	return result
}

func (this *PacketName) SerializePacketNameToJson() (outJson *simplejson.Json) {
	outJson = simplejson.New()
	outJson.Set("packageDownloadUrl", this.packageDownloadUrl)
	outJson.Set("confirm", this.sendConfirm)
	outJson.Set("packageName", this.name)
	return outJson
}

type Style struct {
	builderId int
	ring      int
	vibrate   int
	clearable int
	nId       int
	lights    int
	iconType  int
	styleId   int
}

// ring default 0
// vibrate default 0
// clearable default 1
// nId default 0
// lights default 1
// iconType default 0
// styleId default 1
func NewStyle(builderId, ring, vibrate, clearable, nId, lights, iconType, styleId int) (result *Style) {
	result = new(Style)
	result.builderId = builderId
	result.ring = ring
	result.vibrate = vibrate
	result.clearable = clearable
	result.nId = nId
	result.lights = lights
	result.iconType = iconType
	result.styleId = styleId
	return result
}

type Action struct {
	actionType int
	activity   string
	atyAttr    *ActionAttribute
	browser    *Browser
	packetName *PacketName
	intent     string
}

func NewAction(inType int, inActivity string, inAtyAttr *ActionAttribute, inBrowser *Browser, inPacketName *PacketName, intent string) (result *Action) {
	result = new(Action)
	result.actionType = inType
	result.activity = inActivity
	result.atyAttr = inAtyAttr
	result.browser = inBrowser
	result.packetName = inPacketName
	result.intent = intent
	return result
}

func (this *Action) SerialzeActionToJson() (result *simplejson.Json) {
	result = simplejson.New()
	result.Set("action_type", this.actionType)
	result.Set("browser", this.browser.SerializeBrowserToJson())
	result.Set("activity", this.activity)
	result.Set("intent", this.intent)
	result.Set("aty_attr", this.atyAttr.SerializeAttrToJson())
	result.Set("package_name", this.packetName.SerializePacketNameToJson())
	return result
}

type Message struct {
	sendTitle    string
	sendContent  string
	expireTime   int64
	sendTime     string
	acceptTimes  []interface{}
	msgType      int
	sendStyle    *Style
	sendAction   *simplejson.Json
	sendCustom   *simplejson.Json
	loopInterval int
	loopTimes    int
}

func NewMessage(inTitle, inContent string,
	inExpireTime int64,
	inSendTime string,
	inAcceptTime []interface{},
	inMsgType int,
	inStyle *Style,
	inAction *simplejson.Json,
	inCustom *simplejson.Json) (result *Message) {
	result = new(Message)
	result.sendTitle = inTitle
	result.sendContent = inContent
	result.expireTime = inExpireTime
	result.sendTime = inSendTime
	result.acceptTimes = inAcceptTime
	result.msgType = inMsgType
	result.sendStyle = inStyle
	result.sendAction = inAction
	result.sendCustom = inCustom
	return result
}

func (this *Message) SerializeMessageToJson() (result *simplejson.Json) {
	result = simplejson.New()
	result.Set("title", this.sendTitle)
	result.Set("content", this.sendContent)
	result.Set("accept_time", this.acceptTimes)
	result.Set("builder_id", this.sendStyle.builderId)
	result.Set("ring", this.sendStyle.ring)
	result.Set("vibrate", this.sendStyle.vibrate)
	result.Set("clearable", this.sendStyle.clearable)
	result.Set("n_id", this.sendStyle.nId)
	result.Set("lights", this.sendStyle.lights)
	result.Set("icon_type", this.sendStyle.iconType)
	result.Set("style_id", this.sendStyle.styleId)
	result.Set("action", this.sendAction)
	result.Set("custom_content", this.sendCustom)
	return result
}

func (this *Callback) OnConnect(c *gotcp.Conn) bool {
	addr := c.GetRawConn().RemoteAddr()
	c.PutExtraData(addr)
	infoLog.Println("OnConnect:", addr)
	return true
}

func (this *Callback) OnMessage(c *gotcp.Conn, p gotcp.Packet) bool {
	packet, ok := p.(*common.HeadV3Packet)
	if !ok {
		infoLog.Printf("OnMessage Convert to HeadV3Packet failed")
		return false
	}
	head, err := packet.GetHead()
	if err != nil {
		infoLog.Printf("OnMessage packet.GetHead faild")
		return false
	}

	infoLog.Printf("OnMessage:head=%v packetLen=%v", head, packet.GetLength())
	_, err = packet.CheckPacketValid()
	if err != nil {
		SendResp(c, head, uint16(RET_CODE_RET_INPUT_PARAM_ERR))
		infoLog.Println("Invalid packet", err)
		return false
	}

	// 统计总的请求量
	attr := "goxinge/total_recv_req_count"
	libcomm.AttrAdd(attr, 1)
	// 继续进行后续处理
	go ProcData(c, head, packet.GetBody())
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

func ProcData(c *gotcp.Conn, reqHead *common.HeadV3, payLoad []byte) bool {
	// 首先异步返回响应
	SendResp(c, reqHead, uint16(RET_CODE_RET_SUCCESS))

	// 然后进行下一步无处理
	rootObj, err := simplejson.NewJson(payLoad)
	if err != nil {
		infoLog.Printf("ProcData simplejson new packet error", err)
		return false
	}

	infoLog.Printf("ProcData rootObj=%#v", rootObj)

	var nullString string
	actionAttr := NewActionAttribute(0, 0)
	browser := NewBrowser(nullString, 1)
	packetName := NewPacketName(nullString, 1, nullString)
	action := NewAction(TYPE_ACTIVITY, nullString, actionAttr, browser, packetName, nullString)
	style := NewStyle(1, rootObj.Get("sound").MustInt(1), 1, 1, 0, rootObj.Get("lights").MustInt(1), 0, 1)
	// 自定义字段
	customObj := simplejson.New()
	customObj.Set("chat_type", rootObj.Get("chat_type").MustString("CHAT"))
	customObj.Set("sound", rootObj.Get("sound").MustInt(1))
	customObj.Set("lights", rootObj.Get("lights").MustInt(1))
	customObj.Set("vibrate", rootObj.Get("sound").MustInt(0))
	customObj.Set("type", rootObj.Get("msg_type").MustString("text"))
	customObj.Set("msg_id", rootObj.Get("msg_id").MustString(""))
	customObj.Set("from_id", uint32(rootObj.Get("from_id").MustInt64(0)))
	customObj.Set("to_id", uint32(rootObj.Get("to_id").MustInt64(0)))
	customObj.Set("actionid", uint32(rootObj.Get("actionid").MustInt64(0)))
	customObj.Set("byAt", uint32(rootObj.Get("byAt").MustInt64(0)))
	infoLog.Printf("procData customObj=%#v fromId=%v", customObj, uint32(rootObj.Get("from_id").MustInt64(0)))
	message := NewMessage(rootObj.Get("title").MustString(""),
		rootObj.Get("content").MustString(""),
		0,
		"2013-12-19 17:49:00",
		[]interface{}{},
		TYPE_NOTIFICATION,
		style,
		action.SerialzeActionToJson(),
		customObj)

	messageSlice, err := message.SerializeMessageToJson().MarshalJSON()
	if err != nil {
		infoLog.Printf("procData MarshalJSON failed messag=%v", message)
		return false
	}
	timeNow := time.Now().Unix()
	params := map[string]string{
		"access_id":    strconv.Itoa(cAccessId),
		"expire_time":  strconv.Itoa(int(message.expireTime)),
		"send_time":    message.sendTime,
		"device_token": rootObj.Get("token").MustString(""),
		"message_type": strconv.Itoa(message.msgType),
		"message":      string(messageSlice),
		"timestamp":    strconv.Itoa(int(timeNow)),
		"environment":  "0",
	}
	CallRestful(RESTAPI_PUSHSINGLEDEVICE, params)
	return true

}

func CallRestful(reqUrl string, params map[string]string) {
	sign, err := GenerateSign(METHOD_POST, reqUrl, cSecretKey, params)
	if err != nil {
		infoLog.Printf("CallRestful GenerateSign failed err=%v", err)
		return
	}
	params["sign"] = sign
	data := url.Values{}
	for k, v := range params {
		data.Add(k, v)
	}

	infoLog.Printf("CallRestful reqUrl=%s params=%#v", reqUrl, params)

	encodeParams := data.Encode()
	PushTokenHTAndroid(reqUrl, encodeParams)
}

func GenerateSign(method, reqUrl, secretKey string, params map[string]string) (sign string, err error) {
	upMethod := strings.ToUpper(method)
	parsedUrl, err := url.Parse(reqUrl)
	if err != nil {
		return sign, err
	}
	if len(parsedUrl.Host) != 0 && len(parsedUrl.Path) != 0 {
		reqUrl = parsedUrl.Host + parsedUrl.Path
		infoLog.Printf("GenerateSign new url=%s host=%s path=%s", reqUrl, parsedUrl.Host, parsedUrl.Path)
	}

	var paramStr string
	if len(params) != 0 {
		keySlice := make([]string, len(params))
		index := 0
		for k, _ := range params {
			keySlice[index] = k
			index++
		}
		sort.Strings(keySlice)
		var buffer bytes.Buffer
		for _, v := range keySlice {
			buffer.WriteString(v)
			buffer.WriteString("=")
			buffer.WriteString(params[v])
		}
		paramStr = buffer.String()
	}

	md5String := upMethod + reqUrl + paramStr + secretKey
	hashMd5 := md5.New()
	io.WriteString(hashMd5, md5String)
	sign = fmt.Sprintf("%x", hashMd5.Sum(nil))
	// infoLog.Printf("GenerateSign sign=%s", sign)
	return sign, nil
}

func PushTokenHTAndroid(url, encodeParams string) {
	resp, err := http.Post(url, "application/x-www-form-urlencoded", strings.NewReader(encodeParams))
	if err != nil {
		infoLog.Printf("PushTokenHTAndroid http.Post failed err=%v", err)
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
		infoLog.Printf("PushTokenHTAndroid http.Post err=%v", err)
		return
	}

	infoLog.Printf("PushTokenHTAndroid body=%s", string(body))
	rspJson, err := simplejson.NewJson(body)
	if err != nil {
		infoLog.Printf("PushTokenHTAndroid simplejson.NewJson failed err=%v", err)
		return
	}

	retCode := rspJson.Get("ret_code").MustInt(0)
	if retCode == 0 {
		// 成功总的请求量
		attr := "goxinge/success_req_count"
		libcomm.AttrAdd(attr, 1)
		infoLog.Printf("PushTokenHTAndroid success ret=%v", retCode)
	} else {
		// 失败总的请求量
		attr := "goxinge/failed_req_count"
		libcomm.AttrAdd(attr, 1)
		infoLog.Printf("PushTokenHTAndroid failed ret=%v", retCode)
	}
}

func (this *Callback) OnClose(c *gotcp.Conn) {
	infoLog.Println("OnClose:", c.GetExtraData())
}

func md5Test() {
	h := md5.New()
	io.WriteString(h, "POSTopenapi.xg.qq.com/v2/push/single_deviceaccess_id=2100093818device_token=582cce9ff8ccd2c642bfbc08266a4edbbd7221e0environment=0expire_time=0message")
	infoLog.Printf("md5Sum=%x", h.Sum(nil))
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
	srv := gotcp.NewServer(config, &Callback{}, &common.HeadV3Protocol{})

	// starts service
	go srv.Start(listener, time.Second)
	infoLog.Println("listening:", listener.Addr())

	//MD5TEST
	// md5Test()

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
