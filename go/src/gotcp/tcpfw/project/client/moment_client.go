package main

import (
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/gansidui/gotcp/tcpfw/common"
	"github.com/gansidui/gotcp/tcpfw/include/ht_moment"
	"github.com/golang/protobuf/proto"
	"github.com/jessevdk/go-flags"
	"gopkg.in/ini.v1"
)

var (
	infoLog *log.Logger
)

type Options struct {
	// Example of verbosity with level
	Verbose []bool `short:"v" long:"verbose" description:"Verbose output"`

	// Example of optional value
	Config  string `short:"c" long:"config" description:"config file" optional:"no"`
	Cmd     uint32 `short:"t" long:"command" description:"test commad" optional:"no"`
	Op      uint32 `short:"o" long:"operator" description:"operator add or del"`
	List    uint32 `short:"l" logn:"list" description:"list type"`
	Type    uint32 `short:"y" long:"Latest" description:"get latest or history moment"`
	MntType uint32 `short:"n" long:"moment" description:"moment type"`
	Seq     uint32 `short:"s" long:"seq" description:"specific seq"`
}

var options Options

var parser = flags.NewParser(&options, flags.Default)

func main() {
	// 处理命令行参数
	if _, err := parser.Parse(); err != nil {
		log.Fatalln("parse cmd line failed!")
		os.Exit(1)
	}
	if options.Config == "" {
		log.Fatalln("Must input config file name")
	}

	// log.Println("config name =", options.ClientConf)
	// 读取配置文件
	cfg, err := ini.Load([]byte(""), options.Config)
	if err != nil {
		log.Printf("load config file=%s failed", options.Config)
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

	uid := cfg.Section("TEST_UID").Key("UID").MustInt(1946612)
	headV2Protocol := &common.HeadV2Protocol{}
	reqBody := new(ht_moment.ReqBody)
	head := &common.HeadV2{Cmd: uint32(ht_moment.CMD_TYPE_CMD_OPERATOR_UID),
		Len: uint32(common.EmptyPacktV2Len),
		Uid: uint32(uid)}
	cmd := options.Cmd
	infoLog.Println("cmd =", cmd)
	switch cmd {
	case uint32(ht_moment.CMD_TYPE_CMD_GET_MID_LIST_DOWN):
		{
			head.Cmd = uint32(ht_moment.CMD_TYPE_CMD_GET_MID_LIST_DOWN)
			subReqBody := new(ht_moment.ViewLatestMomentIDReqBody)
			subReqBody.Userid = new(uint32)
			if options.MntType == 1 { // for following
				*(subReqBody.Userid) = 2325928
			} else {
				*(subReqBody.Userid) = uint32(uid)
			}
			subReqBody.Qtype = new(ht_moment.QUERY_TYPE)
			*(subReqBody.Qtype) = ht_moment.QUERY_TYPE(options.MntType)
			subReqBody.LangType = new(uint32)
			*(subReqBody.LangType) = 2 // 母语是中文
			subReqBody.Nationality = []byte("CN")
			subReqBody.LocalMaxMid = []byte("0")
			subReqBody.ReqUserId = new(uint32)
			*(subReqBody.ReqUserId) = uint32(uid)
			reqBody.ViewLatestMntIdReqbody = subReqBody
		}
	case uint32(ht_moment.CMD_TYPE_CMD_GET_MID_LIST_UP):
		{

		}
	case uint32(ht_moment.CMD_TYPE_CMD_OPERATOR_UID):
		{
			head.Cmd = uint32(ht_moment.CMD_TYPE_CMD_OPERATOR_UID)
			subReqBody := new(ht_moment.OperatorUidReqBody)
			subReqBody.ListType = new(ht_moment.OPERATORLIST)
			*(subReqBody.ListType) = ht_moment.OPERATORLIST(options.List) //被隐藏帖子的用户列表
			subReqBody.OpType = new(ht_moment.OPERATORTYPE)
			*(subReqBody.OpType) = ht_moment.OPERATORTYPE(options.Op) //添加操作
			subReqBody.OpUserId = new(uint32)
			*(subReqBody.OpUserId) = 2325928
			subReqBody.Userid = new(uint32)
			*(subReqBody.Userid) = uint32(uid)
			reqBody.OpUidReqbody = subReqBody
		}
	case uint32(ht_moment.CMD_TYPE_CMD_GET_OPERATOR_UID_LIST):
		{
			head.Cmd = uint32(ht_moment.CMD_TYPE_CMD_GET_OPERATOR_UID_LIST)
			subReqBody := new(ht_moment.GetUidListReqBody)
			subReqBody.ListType = new(ht_moment.OPERATORLIST)
			*(subReqBody.ListType) = ht_moment.OPERATORLIST(options.List)
			subReqBody.HideListVer = new(uint32)
			*(subReqBody.HideListVer) = uint32(time.Now().Unix())
			subReqBody.NotShareVer = new(uint32)
			*(subReqBody.NotShareVer) = uint32(time.Now().Unix())
			subReqBody.Userid = new(uint32)
			*(subReqBody.Userid) = uint32(uid)
			reqBody.GetOpUidListReqbody = subReqBody
		}
	case uint32(ht_moment.CMD_TYPE_CMD_BACKEND_OP_UID):
		{
			head.Cmd = uint32(ht_moment.CMD_TYPE_CMD_BACKEND_OP_UID)
			subReqBody := new(ht_moment.BackEndOpUidReqBody)
			subReqBody.ListType = new(ht_moment.OPERATORLIST)
			*(subReqBody.ListType) = ht_moment.OPERATORLIST(options.List) // Moment 指给粉丝看
			subReqBody.OpType = new(ht_moment.OPERATORTYPE)
			*(subReqBody.OpType) = ht_moment.OPERATORTYPE(options.Op) //添加操作
			subReqBody.OpUserId = new(uint32)
			*(subReqBody.OpUserId) = uint32(uid)
			subReqBody.Userid = new(uint32)
			*(subReqBody.Userid) = uint32(uid)
			subReqBody.OpReason = []byte("back server")
			reqBody.BackendOpUidReqbody = subReqBody
		}
	case uint32(ht_moment.CMD_TYPE_CMD_QUERY_REVEAL_ENTRY):
		{
			head.Cmd = uint32(ht_moment.CMD_TYPE_CMD_QUERY_REVEAL_ENTRY)
			subReqBody := new(ht_moment.QueryRevealEntryReqBody)
			subReqBody.ReqUserId = new(uint32)
			*(subReqBody.ReqUserId) = uint32(uid)
			subReqBody.Qtype = make([]ht_moment.QUERY_TYPE, 1)
			subReqBody.Qtype[0] = ht_moment.QUERY_TYPE(options.MntType)
			reqBody.QueryRevealentryReqbody = subReqBody
		}
	case uint32(ht_moment.CMD_TYPE_CMD_GET_USER_INDEX_MOMENT):
		{
			head.Cmd = uint32(ht_moment.CMD_TYPE_CMD_GET_USER_INDEX_MOMENT)
			subReqBody := new(ht_moment.GetUserIndexMomentReqBody)
			subReqBody.Userid = new(uint32)
			*(subReqBody.Userid) = uint32(uid)
			momentIndex := new(ht_moment.MomentUserIndexInfo)
			momentIndex.Userid = new(uint32)
			*(momentIndex.Userid) = uint32(uid)
			momentIndex.Seq = new(uint32)
			*(momentIndex.Seq) = uint32(8)
			subReqBody.UserIndex = make([]*ht_moment.MomentUserIndexInfo, 1)
			subReqBody.UserIndex[0] = momentIndex
			reqBody.GetUserIndexMomentReqbody = subReqBody
		}

	default:
		{
			infoLog.Println("unknow cmd =", cmd)
			return
		}
	}

	payLoad, err := proto.Marshal(reqBody)
	if err != nil {
		infoLog.Println("marshaling error: ", err)
	}

	head.Len = uint32(common.PacketV2HeadLen) + uint32(len(payLoad)) + 1
	buf := make([]byte, head.Len)
	buf[0] = common.HTV2MagicBegin
	err = common.SerialHeadV2ToSlice(head, buf[1:])
	if err != nil {
		infoLog.Println("SerialHeadV3ToSlice failed")
		os.Exit(1)
	}
	copy(buf[common.PacketV2HeadLen:], payLoad)
	buf[head.Len-1] = common.HTV2MagicEnd

	infoLog.Println("payLoad=", common.NewHeadV2Packet(buf).Serialize())
	// write
	conn.Write(common.NewHeadV2Packet(buf).Serialize())

	// read
	p, err := headV2Protocol.ReadPacket(conn)
	if err == nil {
		packet := p.(*common.HeadV2Packet)
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
					infoLog.Printf("stat=%v list_type=%v", subRspBody.GetStatus().GetCode(), subRspBody.GetListType())
				}
			case uint32(ht_moment.CMD_TYPE_CMD_GET_OPERATOR_UID_LIST):
				{
					subRspBody := rspBody.GetGetOpUidListRspbody()
					status := subRspBody.GetStatus()
					infoLog.Println("stat =", status.GetCode())
					listType := subRspBody.GetListType()
					switch listType {

					case ht_moment.OPERATORLIST_OP_HIDE_LIST:
						{
							uidArray := subRspBody.GetHideUidArray().GetUidArray()
							infoLog.Printf("list version=%v", subRspBody.GetHideUidArray().GetListVer())
							for i, v := range uidArray {
								infoLog.Printf("index=%v uid=%v timeStamp=%v", i, v.GetUid(), v.GetOpTimeStamp())
							}
						}
					case ht_moment.OPERATORLIST_OP_NOT_SHARE_LIST:
						{
							uidArray := subRspBody.GetNotShareArray().GetUidArray()
							infoLog.Printf("list version=%v", subRspBody.GetNotShareArray().GetListVer())
							for i, v := range uidArray {
								infoLog.Printf("index=%v uid=%v timeStamp=%v", i, v.GetUid(), v.GetOpTimeStamp())
							}
						}
					}

				}
			case uint32(ht_moment.CMD_TYPE_CMD_BACKEND_OP_UID):
				{
					subRspBody := rspBody.GetBackendOpUidRspbody()
					status := subRspBody.GetStatus()
					infoLog.Printf("status=%v list_type=%v", status.GetCode(), subRspBody.GetListType())
				}
			case uint32(ht_moment.CMD_TYPE_CMD_GET_MID_LIST_DOWN):
				{
					subRspBody := rspBody.GetViewLatestMntIdRspbody()
					status := subRspBody.GetStatus()
					infoLog.Printf("status=%+v pagesize=%d", status, subRspBody.GetPageSize())
					bucketList := subRspBody.GetBucket().GetBucketList()
					for i, v := range bucketList {
						infoLog.Printf("index=%v bucket name=%v index=%v", i, v.GetBucketName(), v.GetIndex())
					}
					idList := subRspBody.GetIdList()
					for i, v := range idList {
						infoLog.Printf("index=%v mid=%v LikedTs=%v CommentTs=%v Deleted=%v OpReason=%v",
							i,
							v.GetMid(),
							v.GetLikedTs(),
							v.GetCommentTs(),
							v.GetDeleted(),
							v.GetOpReason())

					}

				}
			case uint32(ht_moment.CMD_TYPE_CMD_QUERY_REVEAL_ENTRY):
				{
					subRspBody := rspBody.GetQueryRevealentryRspbody()
					status := subRspBody.GetStatus()
					len := len(subRspBody.ListReveal)
					infoLog.Printf("QueryRevealEntry ret status=%v ListRevral size=%d\n", status, len)
					if len > 0 {
						revealBody := subRspBody.ListReveal[0]
						infoLog.Printf("tyep=%v be_reveal=%v\n", revealBody.GetQtype(), revealBody.GetBeReveal())
					} else {
						infoLog.Println("error len = 0")
					}
				}
			case uint32(ht_moment.CMD_TYPE_CMD_GET_USER_INDEX_MOMENT):
				{
					subRspBody := rspBody.GetGetUserIndexMomentRspbody()
					status := subRspBody.GetStatus()
					infoLog.Printf("QueryMid ret status=%v\n", status)
					if status.GetCode() == 0 {
						moments := subRspBody.GetMoment()
						for i := 0; i < len(moments); i++ {
							iterm := moments[i]
							infoLog.Printf("uid=%v mid=%s posttime=%v lat=%v long=%v lang=%v\n", iterm.GetUserid(), iterm.GetMid(), iterm.GetPostTime(), iterm.GetLatitude(), iterm.GetLongitude(), iterm.GetLangType())
						}
					}

				}
			}
		} else {
			infoLog.Println("get HTV2Head failed err =", err)
		}
	}
	conn.Close()
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
