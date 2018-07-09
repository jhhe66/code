package main

/*
#cgo CFLAGS: -I../../../libcomm/cinclude
#cgo LDFLAGS: -L../../../libcomm/ -lneocomm
*/
import "C"
import (
	// "fmt"
	// "github.com/bitly/go-simplejson"
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/BLHT/HT_GOGO/gotcp/tcpfw/common"
	"github.com/BLHT/HT_GOGO/gotcp/tcpfw/include/ht_favorite"
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
	Cmd  string `short:"c" long:"cmd" description:"command,1:add, 2:del, 3:get, 4:update tag, 5:search tag" optional:"no"`
	Obid string `short:"i" long:"objectid" description:"objectid" optional:"no"`
	Text string `short:"t" long:"text" description:"text/tag" optional:"no"`
}

var options Options

var parser = flags.NewParser(&options, flags.Default)

func main() {
	// 处理命令行参数
	if _, err := parser.Parse(); err != nil {
		fmt.Println("parse cmd line failed!")
		os.Exit(1)
	}
	cmd, _ := strconv.Atoi(options.Cmd)
	obid := options.Obid
	text := options.Text

	// 读取配置文件
	cfg, err := ini.Load([]byte(""), "test_config.ini")
	if err != nil {
		fmt.Println("load config file failed")
		return
	}
	// 配置文件只读 设置次标识提升性能
	cfg.BlockMode = false
	// 读取ip+port
	serverIp := cfg.Section("OUTER_SERVER").Key("server_ip").MustString("127.0.0.1")
	serverPort := cfg.Section("OUTER_SERVER").Key("server_port").MustInt(8990)

	fmt.Printf("server_ip=%v server_port=%v\n", serverIp, serverPort)

	tcpAddr, err := net.ResolveTCPAddr("tcp4", serverIp+":"+strconv.Itoa(serverPort))
	checkError(err)
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Printf("DialTCP addr[%v] failed!\n", tcpAddr)
		return
	}
	defer conn.Close()
	checkError(err)
	headV3Protocol := &common.HeadV3Protocol{}
	var head *common.HeadV3
	head = &common.HeadV3{Flag: 0xF0, TermType: 1, Cmd: uint16(cmd), From: 1603139, To: 1946612}
	reqBody := new(ht_favorite.ReqBody)

	switch ht_favorite.CMD_TYPE(cmd) {
	case ht_favorite.CMD_TYPE_CMD_ADD:
		reqBody.AddReqbody = new(ht_favorite.AddReqBody)
		//reqBody.AddReqbody = &ht_favorite.AddReqBody{}
		content := new(ht_favorite.FavoriteContent)
		reqBody.AddReqbody.Content = content
		content.SrcUid = proto.Uint32(1604048)
		content.Type = ht_favorite.TYPE_FAVORATE_TYPE_MNT.Enum()
		//content.Text = []byte("Cold weather" + strconv.Itoa(time.Now().Second()))
		content.Voice = new(ht_favorite.VoiceBody)
		content.Voice.Url = []byte("www.voice.com")
		content.Voice.Duration = proto.Uint32(60)
		content.Voice.Size = proto.Uint32(100)
		content.Image = new(ht_favorite.ImageBody)
		content.Image.ThumbUrl = []byte("www.image.com")
		content.Image.BigUrl = []byte("www.image.com")
		content.Image.Width = proto.Uint32(20)
		content.Image.Height = proto.Uint32(40)
		content.Mid = []byte("124000000")
		content.Correction = []*ht_favorite.CorrectContent{
			&ht_favorite.CorrectContent{
				Content:    []byte("c1"),
				Correction: []byte("ccc1"),
			},
			&ht_favorite.CorrectContent{
				Content:    []byte("c2"),
				Correction: []byte("ccc2"),
			},
		}
		content.Tags = []string{"abc", "def"}
		/*content.Tags = make([]string, 3)
		for i := 0; i < len(content.Tags); i++ {
			content.Tags[i] = string('a' + i)
		}*/

		payLoad, err := proto.Marshal(reqBody)
		if err != nil {
			fmt.Println("marshaling error: ", err)
			return
		}

		head.Len = uint32(common.PacketV3HeadLen) + uint32(len(payLoad)) + 1
		buf := make([]byte, head.Len)
		buf[0] = common.HTV3MagicBegin
		err = common.SerialHeadV3ToSlice(head, buf[1:])
		if err != nil {
			fmt.Println("SerialHeadV3ToSlice failed")
			return
		}
		copy(buf[common.PacketV3HeadLen:], payLoad)
		buf[head.Len-1] = common.HTV3MagicEnd
		//fmt.Printf("sendpkg's len=%v\n", len(common.NewHeadV3Packet(buf).Serialize()))

		// write
		conn.Write(common.NewHeadV3Packet(buf).Serialize())

		// read
		p, err := headV3Protocol.ReadPacket(conn)
		if err != nil {
			fmt.Println("ReadPacket failed!")
			return
		}
		packet := p.(*common.HeadV3Packet)
		rspHead, err := packet.GetHead()
		if err != nil {
			fmt.Println("get HTV3Head failed err =", err)
			return
		}
		fmt.Printf("ret=%v, len=%v\n", rspHead.Ret, rspHead.Len)
		if rspHead.Ret != 0 {
			return
		}

		rspBody := new(ht_favorite.RspBody)
		err = proto.Unmarshal(packet.GetBody(), rspBody)
		if err != nil {
			fmt.Println("proto Unmarshal failed")
			return
		}
		addRspBody := rspBody.GetAddRspbody()
		if addRspBody == nil {
			fmt.Println("GetAddRspbody() failed")
			return
		}
		fmt.Printf("obid[%s], lastts[%v]\n", addRspBody.GetObid(), addRspBody.GetLastTs())
	case ht_favorite.CMD_TYPE_CMD_DEL:
		reqBody.DelReqbody = new(ht_favorite.DelReqBody)
		reqBody.DelReqbody.Obid = []byte(obid)

		payLoad, err := proto.Marshal(reqBody)
		if err != nil {
			fmt.Println("marshaling error: ", err)
			return
		}

		head.Len = uint32(common.PacketV3HeadLen) + uint32(len(payLoad)) + 1
		buf := make([]byte, head.Len)
		buf[0] = common.HTV3MagicBegin
		err = common.SerialHeadV3ToSlice(head, buf[1:])
		if err != nil {
			fmt.Println("SerialHeadV3ToSlice failed")
			return
		}
		copy(buf[common.PacketV3HeadLen:], payLoad)
		buf[head.Len-1] = common.HTV3MagicEnd
		//fmt.Printf("sendpkg's len=%v\n", len(common.NewHeadV3Packet(buf).Serialize()))

		// write
		conn.Write(common.NewHeadV3Packet(buf).Serialize())

		// read
		p, err := headV3Protocol.ReadPacket(conn)
		if err != nil {
			fmt.Println("ReadPacket failed!")
			return
		}
		packet := p.(*common.HeadV3Packet)
		rspHead, err := packet.GetHead()
		if err != nil {
			fmt.Println("get HTV3Head failed err =", err)
			return
		}
		fmt.Printf("ret=%v, len=%v\n", rspHead.Ret, rspHead.Len)
		if rspHead.Ret != 0 {
			return
		}

		rspBody := new(ht_favorite.RspBody)
		err = proto.Unmarshal(packet.GetBody(), rspBody)
		if err != nil {
			fmt.Println("proto Unmarshal failed")
			return
		}
		delRspBody := rspBody.GetDelRspbody()
		if delRspBody == nil {
			fmt.Println("GetAddRspbody() failed")
			return
		}
		fmt.Printf("lastts[%v]\n", delRspBody.GetLastTs())

	case ht_favorite.CMD_TYPE_CMD_GET:
		reqBody.GetReqbody = new(ht_favorite.GetReqBody)
		reqBody.GetReqbody.Index = proto.String(obid)

		payLoad, err := proto.Marshal(reqBody)
		if err != nil {
			fmt.Println("marshaling error: ", err)
			return
		}

		head.Len = uint32(common.PacketV3HeadLen) + uint32(len(payLoad)) + 1
		buf := make([]byte, head.Len)
		buf[0] = common.HTV3MagicBegin
		err = common.SerialHeadV3ToSlice(head, buf[1:])
		if err != nil {
			fmt.Println("SerialHeadV3ToSlice failed")
			return
		}
		copy(buf[common.PacketV3HeadLen:], payLoad)
		buf[head.Len-1] = common.HTV3MagicEnd
		//fmt.Printf("sendpkg's len=%v\n", len(common.NewHeadV3Packet(buf).Serialize()))

		// write
		conn.Write(common.NewHeadV3Packet(buf).Serialize())

		// read
		p, err := headV3Protocol.ReadPacket(conn)
		if err != nil {
			fmt.Println("ReadPacket failed!")
			return
		}
		packet := p.(*common.HeadV3Packet)
		rspHead, err := packet.GetHead()
		if err != nil {
			fmt.Println("get HTV3Head failed err =", err)
			return
		}
		fmt.Printf("ret=%v, len=%v\n", rspHead.Ret, rspHead.Len)
		if rspHead.Ret != 0 {
			return
		}

		rspBody := new(ht_favorite.RspBody)
		err = proto.Unmarshal(packet.GetBody(), rspBody)
		if err != nil {
			fmt.Println("proto Unmarshal failed")
			return
		}
		getRspBody := rspBody.GetGetRspbody()
		if getRspBody == nil {
			fmt.Println("GetGetRspbody() failed")
			return
		}
		fmt.Printf("index[%s], more[%v], lastts[%v]\n", getRspBody.GetIndex(), getRspBody.GetMore(), getRspBody.GetLastTs())
		clist := getRspBody.GetContentList()
		for i, content := range clist {
			fmt.Printf("i[%v], obid[%s]\n", i, content.GetObid())
		}
	case ht_favorite.CMD_TYPE_CMD_SEARCH_TAG:
		reqBody.SearchTagReqbody = new(ht_favorite.SearchTagReqBody)
		reqBody.SearchTagReqbody.Tag = proto.String(text)
		reqBody.SearchTagReqbody.Index = proto.String(obid)

		payLoad, err := proto.Marshal(reqBody)
		if err != nil {
			fmt.Println("marshaling error: ", err)
			return
		}

		head.Len = uint32(common.PacketV3HeadLen) + uint32(len(payLoad)) + 1
		buf := make([]byte, head.Len)
		buf[0] = common.HTV3MagicBegin
		err = common.SerialHeadV3ToSlice(head, buf[1:])
		if err != nil {
			fmt.Println("SerialHeadV3ToSlice failed")
			return
		}
		copy(buf[common.PacketV3HeadLen:], payLoad)
		buf[head.Len-1] = common.HTV3MagicEnd
		//fmt.Printf("sendpkg's len=%v\n", len(common.NewHeadV3Packet(buf).Serialize()))

		// write
		conn.Write(common.NewHeadV3Packet(buf).Serialize())

		// read
		p, err := headV3Protocol.ReadPacket(conn)
		if err != nil {
			fmt.Println("ReadPacket failed!")
			return
		}
		packet := p.(*common.HeadV3Packet)
		rspHead, err := packet.GetHead()
		if err != nil {
			fmt.Println("get HTV3Head failed err =", err)
			return
		}
		fmt.Printf("ret=%v, len=%v\n", rspHead.Ret, rspHead.Len)
		if rspHead.Ret != 0 {
			return
		}

		rspBody := new(ht_favorite.RspBody)
		err = proto.Unmarshal(packet.GetBody(), rspBody)
		if err != nil {
			fmt.Println("proto Unmarshal failed")
			return
		}
		searchTagRspBody := rspBody.GetSearchTagRspbody()
		if searchTagRspBody == nil {
			fmt.Println("GetSearchTagRspbody() failed")
			return
		}
		fmt.Printf("index[%s], more[%v], lastts[%v]\n", searchTagRspBody.GetIndex(), searchTagRspBody.GetMore(), searchTagRspBody.GetLastTs())
		clist := searchTagRspBody.GetContentList()
		for i, content := range clist {
			fmt.Printf("i[%v], obid[%s]\n", i, content.GetObid())
		}
	case ht_favorite.CMD_TYPE_CMD_SEARCH_TEXT:
		reqBody.SearchTextReqbody = new(ht_favorite.SearchTextReqBody)
		reqBody.SearchTextReqbody.Text = proto.String(text)
		reqBody.SearchTextReqbody.Index = proto.String(obid)

		payLoad, err := proto.Marshal(reqBody)
		if err != nil {
			fmt.Println("marshaling error: ", err)
			return
		}

		head.Len = uint32(common.PacketV3HeadLen) + uint32(len(payLoad)) + 1
		buf := make([]byte, head.Len)
		buf[0] = common.HTV3MagicBegin
		err = common.SerialHeadV3ToSlice(head, buf[1:])
		if err != nil {
			fmt.Println("SerialHeadV3ToSlice failed")
			return
		}
		copy(buf[common.PacketV3HeadLen:], payLoad)
		buf[head.Len-1] = common.HTV3MagicEnd
		//fmt.Printf("sendpkg's len=%v\n", len(common.NewHeadV3Packet(buf).Serialize()))

		// write
		conn.Write(common.NewHeadV3Packet(buf).Serialize())

		// read
		p, err := headV3Protocol.ReadPacket(conn)
		if err != nil {
			fmt.Println("ReadPacket failed!")
			return
		}
		packet := p.(*common.HeadV3Packet)
		rspHead, err := packet.GetHead()
		if err != nil {
			fmt.Println("get HTV3Head failed err =", err)
			return
		}
		fmt.Printf("ret=%v, len=%v\n", rspHead.Ret, rspHead.Len)
		if rspHead.Ret != 0 {
			return
		}

		rspBody := new(ht_favorite.RspBody)
		err = proto.Unmarshal(packet.GetBody(), rspBody)
		if err != nil {
			fmt.Println("proto Unmarshal failed")
			return
		}
		searchTextRspBody := rspBody.GetSearchTextRspbody()
		if searchTextRspBody == nil {
			fmt.Println("GetSearchTextRspbody() failed")
			return
		}
		fmt.Printf("index[%s], more[%v], lastts[%v]\n", searchTextRspBody.GetIndex(), searchTextRspBody.GetMore(), searchTextRspBody.GetLastTs())
		clist := searchTextRspBody.GetContentList()
		for i, content := range clist {
			fmt.Printf("i[%v], obid[%s]\n", i, content.GetObid())
		}
	case ht_favorite.CMD_TYPE_CMD_UPDATE_TAG:
		reqBody.UpdateTagReqbody = new(ht_favorite.UpdateTagReqBody)
		reqBody.UpdateTagReqbody.Obid = []byte(obid)
		reqBody.UpdateTagReqbody.Tags = []string{"wyf"}

		payLoad, err := proto.Marshal(reqBody)
		if err != nil {
			fmt.Println("marshaling error: ", err)
			return
		}

		head.Len = uint32(common.PacketV3HeadLen) + uint32(len(payLoad)) + 1
		buf := make([]byte, head.Len)
		buf[0] = common.HTV3MagicBegin
		err = common.SerialHeadV3ToSlice(head, buf[1:])
		if err != nil {
			fmt.Println("SerialHeadV3ToSlice failed")
			return
		}
		copy(buf[common.PacketV3HeadLen:], payLoad)
		buf[head.Len-1] = common.HTV3MagicEnd
		//fmt.Printf("sendpkg's len=%v\n", len(common.NewHeadV3Packet(buf).Serialize()))

		// write
		conn.Write(common.NewHeadV3Packet(buf).Serialize())

		// read
		p, err := headV3Protocol.ReadPacket(conn)
		if err != nil {
			fmt.Println("ReadPacket failed!")
			return
		}
		packet := p.(*common.HeadV3Packet)
		rspHead, err := packet.GetHead()
		if err != nil {
			fmt.Println("get HTV3Head failed err =", err)
			return
		}
		fmt.Printf("ret=%v, len=%v\n", rspHead.Ret, rspHead.Len)
		if rspHead.Ret != 0 {
			return
		}

		rspBody := new(ht_favorite.RspBody)
		err = proto.Unmarshal(packet.GetBody(), rspBody)
		if err != nil {
			fmt.Println("proto Unmarshal failed")
			return
		}
		updateTagRspBody := rspBody.GetUpdateTagRspbody()
		if updateTagRspBody == nil {
			fmt.Println("GetUpdateTagRspbody() failed")
			return
		}
		fmt.Printf("lastts[%v]\n", updateTagRspBody.GetLastTs())

	default:
		fmt.Println("Unknown cmd")
	}

}

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
