package main

/*
#cgo CFLAGS: -I../../../libcomm/cinclude
#cgo LDFLAGS: -L../../../libcomm/ -lneocomm
#include "Attr_API.h"
*/
//import "C"
import (
	"crypto/md5"
	"fmt"

	"github.com/BLHT/HT_GOGO/gotcp"
	"github.com/BLHT/HT_GOGO/gotcp/libcomm"
	"github.com/BLHT/HT_GOGO/gotcp/tcpfw/common"
	"github.com/BLHT/HT_GOGO/gotcp/tcpfw/include/ht_favorite"
	"github.com/go-ini/ini"
	flags "github.com/jessevdk/go-flags"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"log"
	"net"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"
	"time"

	"github.com/golang/protobuf/proto"
)

type Callback struct{}

const (
	favorDB          = "favorite"
	favorRecordTable = "favor_record"
	favorTsTable     = "favor_ts"
	pageSize         = 20
)

var (
	infoLog   *log.Logger
	mongoSess *mgo.Session
)

type VoiceInfo struct {
	Url      string
	Duration uint32
	Size     uint32
}

type ImageInfo struct {
	ThumbUrl string
	BigUrl   string
	Width    uint32
	Height   uint32
}

type CorrectInfo struct {
	Content    string
	Correction string
}

type FavorRecord struct {
	Id         bson.ObjectId `bson:"_id"`
	Uid        uint32
	Md5        string
	SrcUid     uint32
	Type       uint32
	Text       string
	Voice      VoiceInfo
	Correction []CorrectInfo
	Image      ImageInfo
	Mid        string
	Tags       []string
}

type FavorTs struct {
	Uid    uint32
	LastTs uint32
}

func (this *Callback) OnConnect(c *gotcp.Conn) bool {
	addr := c.GetRawConn().RemoteAddr()
	c.PutExtraData(addr)
	infoLog.Println("OnConnect:", addr)
	return true
}

func (this *Callback) OnMessage(c *gotcp.Conn, p gotcp.Packet) bool {
	attrName := "favorite_recv_req"
	libcomm.AttrAdd(attrName, 1)
	packet := p.(*common.HeadV3Packet)
	_, err := packet.CheckPacketValid()
	if err != nil {
		infoLog.Println("Invalid packet", err)
		return false
	}

	go ProcData(c, p)
	return true
}

func SendRsp(c *gotcp.Conn, head *common.HeadV3, body *ht_favorite.RspBody, ret uint16) bool {
	rspHead := new(common.HeadV3)
	if head != nil {
		*rspHead = *head
	}
	var outBuf []byte
	var err error
	if body != nil {
		outBuf, err = proto.Marshal(body)
		if err != nil {
			infoLog.Println("Failed to encode RspBody:", err)
			return false
		}
	}

	rspHead.Ret = ret
	//rspHead.Len = len(rspHead) + 2 + body.GetLenth()
	rspHead.Len = common.HeadV3Len + 2 + uint32(len(outBuf))
	sendBuf := make([]byte, rspHead.Len)
	sendBuf[0] = common.HTV3MagicBegin
	err = common.SerialHeadV3ToSlice(rspHead, sendBuf[1:])
	if err != nil {
		infoLog.Println("SerialHeadV3ToSlice failed")
		return false
	}
	copy(sendBuf[1+common.HeadV3Len:], outBuf)
	sendBuf[rspHead.Len-1] = common.HTV3MagicEnd
	infoLog.Printf("ret=%v, rsppkg_len=%v\n", ret, rspHead.Len)
	resp := common.NewHeadV3Packet(sendBuf)
	c.AsyncWritePacket(resp, time.Second)

	return true
}

func ProcData(c *gotcp.Conn, p gotcp.Packet) bool {
	var head *common.HeadV3

	packet, ok := p.(*common.HeadV3Packet)
	if !ok {
		infoLog.Println("Convert to HeadV3Packet failed")
		return false
	}
	head, err := packet.GetHead()
	if err != nil {
		infoLog.Println("ProcData Get head faild")
		return false
	}

	switch ht_favorite.CMD_TYPE(head.Cmd) {
	case ht_favorite.CMD_TYPE_CMD_ADD:
		go procAddReq(c, p)
	case ht_favorite.CMD_TYPE_CMD_DEL:
		go procDelReq(c, p)
	case ht_favorite.CMD_TYPE_CMD_GET:
		go procGetReq(c, p)
	case ht_favorite.CMD_TYPE_CMD_UPDATE_TAG:
		go procUpdateTagReq(c, p)
	case ht_favorite.CMD_TYPE_CMD_SEARCH_TAG:
		go procSearchTagReq(c, p)
	case ht_favorite.CMD_TYPE_CMD_SEARCH_TEXT:
		go procSearchTextReq(c, p)
	default:
		infoLog.Println("Recv unkown cmd:", head.Cmd)
	}

	return true
}

func procAddReq(c *gotcp.Conn, p gotcp.Packet) bool {
	// parse packet
	rspBody := new(ht_favorite.RspBody)

	result := uint16(ht_favorite.RET_CODE_RET_SUCCESS)
	var head *common.HeadV3
	defer func() {
		SendRsp(c, head, rspBody, result)
	}()
	packet, ok := p.(*common.HeadV3Packet)
	if !ok {
		infoLog.Println("Convert to HeadV3Packet failed")
		result = uint16(ht_favorite.RET_CODE_RET_INPUT_PARAM_ERR)
		return false
	}
	head, err := packet.GetHead()
	if err != nil {
		infoLog.Println("ProcData Get head faild")
		result = uint16(ht_favorite.RET_CODE_RET_INTERNAL_ERR)
		return false
	}
	reqBody := new(ht_favorite.ReqBody)
	err = proto.Unmarshal(packet.GetBody(), reqBody)
	if err != nil {
		infoLog.Println("proto Unmarshal failed")
		result = uint16(ht_favorite.RET_CODE_RET_PB_ERR)
		return false
	}
	//check input
	addReqBody := reqBody.GetAddReqbody()
	if addReqBody == nil {
		infoLog.Println("GetAddReqbody() failed")
		result = uint16(ht_favorite.RET_CODE_RET_INPUT_PARAM_ERR)
		return false
	}
	content := addReqBody.GetContent()
	if content == nil {
		infoLog.Println("GetContent() failed")
		result = uint16(ht_favorite.RET_CODE_RET_INPUT_PARAM_ERR)
		return false
	}
	infoLog.Printf("uid[%v],type[%v], int_type[%v],text[%s],mid[%s]\n", packet.GetFromUid(), content.GetType(), int(content.GetType()), content.GetText(), content.GetMid())

	// write mongodb
	record := new(FavorRecord)
	id := bson.NewObjectId()
	record.Id = id
	record.Uid = head.From
	md5sum := md5.Sum(packet.GetBody())
	var md5slice []byte = md5sum[:]
	record.Md5 = fmt.Sprintf("%x", string(md5slice))
	record.SrcUid = content.GetSrcUid()
	record.Type = uint32(content.GetType())
	record.Text = string(content.GetText())
	record.Voice.Url = string(content.GetVoice().GetUrl())
	record.Voice.Duration = content.GetVoice().GetDuration()
	record.Voice.Size = content.GetVoice().GetSize()
	record.Image.ThumbUrl = string(content.GetImage().GetThumbUrl())
	record.Image.BigUrl = string(content.GetImage().GetBigUrl())
	record.Image.Width = content.GetImage().GetWidth()
	record.Image.Height = content.GetImage().GetHeight()
	record.Mid = string(content.GetMid())
	record.Tags = content.GetTags()
	corrCnt := len(content.GetCorrection())
	record.Correction = make([]CorrectInfo, corrCnt)
	for i, corr := range content.GetCorrection() {
		if corr == nil {
			infoLog.Println("content.GetCorrection is nil")
			break
		}
		record.Correction[i].Content = string(corr.GetContent())
		record.Correction[i].Correction = string(corr.GetCorrection())
	}
	mongoConn := mongoSess.DB(favorDB).C(favorRecordTable)
	err = mongoConn.Insert(record)
	// 如果MD5重复则返回重复收藏的错误
	if err != nil {
		infoLog.Println("Insert failed, err=", err)
		infoLog.Println(err.Code)
		result = uint16(ht_favorite.RET_CODE_RET_REPEAT_ADD)
		return true
	}

	recordTs := new(FavorTs)
	recordTs.Uid = head.From
	recordTs.LastTs = uint32(time.Now().Unix())
	mongoConn2 := mongoSess.DB(favorDB).C(favorTsTable)
	_, err = mongoConn2.Upsert(bson.M{"uid": recordTs.Uid}, recordTs)
	if err != nil {
		infoLog.Println("Upsert failed")
		result = uint16(ht_favorite.RET_CODE_RET_DB_ERR)
		return true
	}

	//send response packet
	rspBody.AddRspbody = new(ht_favorite.AddRspBody)
	rspBody.AddRspbody.Obid = []byte(id.String())
	rspBody.AddRspbody.LastTs = proto.Uint32(recordTs.LastTs)

	return true
}

func procDelReq(c *gotcp.Conn, p gotcp.Packet) bool {
	rspBody := new(ht_favorite.RspBody)
	// parse packet
	result := uint16(ht_favorite.RET_CODE_RET_SUCCESS)
	var head *common.HeadV3
	defer func() {
		SendRsp(c, head, rspBody, result)
	}()
	packet, ok := p.(*common.HeadV3Packet)
	if !ok {
		infoLog.Println("Convert to HeadV3Packet failed")
		result = uint16(ht_favorite.RET_CODE_RET_INPUT_PARAM_ERR)
		return false
	}
	head, err := packet.GetHead()
	if err != nil {
		infoLog.Println("ProcData Get head faild")
		result = uint16(ht_favorite.RET_CODE_RET_INTERNAL_ERR)
		return false
	}
	reqBody := &ht_favorite.ReqBody{}
	err = proto.Unmarshal(packet.GetBody(), reqBody)
	if err != nil {
		infoLog.Println("proto Unmarshal failed")
		result = uint16(ht_favorite.RET_CODE_RET_PB_ERR)
		return false
	}
	delReqBody := reqBody.GetDelReqbody()
	if delReqBody == nil {
		infoLog.Println("GetDelReqbody() failked")
		result = uint16(ht_favorite.RET_CODE_RET_PB_ERR)
		return false
	}
	infoLog.Printf("uid[%v],obid[%s]", packet.GetFromUid(), delReqBody.GetObid())

	//check param
	/*if string(delReqBody.GetObid()) == "" {
		infoLog.Println("Obid in reqbody is nil")
		result = uint16(ht_favorite.RET_CODE_RET_INPUT_PARAM_ERR)
		return false
	}*/
	if !bson.IsObjectIdHex(string(delReqBody.GetObid())) {
		infoLog.Printf("Obid[%s] is invalid hex string\n", delReqBody.GetObid())
		result = uint16(ht_favorite.RET_CODE_RET_INPUT_PARAM_ERR)
		return false
	}

	// write mongodb
	mongoConn := mongoSess.DB(favorDB).C(favorRecordTable)
	obid := bson.ObjectIdHex(string(delReqBody.GetObid()))
	err = mongoConn.RemoveId(obid)
	if err != nil {
		if err == mgo.ErrNotFound {
			infoLog.Printf("Obid[%s] to remove is not found", delReqBody.GetObid())
			result = uint16(ht_favorite.RET_CODE_RET_NOT_EXIST)
			return true
		} else {
			infoLog.Println("Remove failed, err=", err)
			result = uint16(ht_favorite.RET_CODE_RET_DB_ERR)
			return false
		}
	}

	recordTs := new(FavorTs)
	recordTs.Uid = head.From
	recordTs.LastTs = uint32(time.Now().Unix())
	mongoConn2 := mongoSess.DB(favorDB).C(favorTsTable)
	_, err = mongoConn2.Upsert(bson.M{"uid": recordTs.Uid}, recordTs)
	if err != nil {
		infoLog.Println("Insert failed")
		result = uint16(ht_favorite.RET_CODE_RET_DB_ERR)
		return true
	}

	//send response packet
	rspBody.DelRspbody = new(ht_favorite.DelRspBody)
	rspBody.DelRspbody.LastTs = proto.Uint32(recordTs.LastTs)

	return true
}

func procGetReq(c *gotcp.Conn, p gotcp.Packet) bool {
	// parse packet
	rspBody := new(ht_favorite.RspBody)
	result := uint16(ht_favorite.RET_CODE_RET_SUCCESS)
	var head *common.HeadV3
	defer func() {
		SendRsp(c, head, rspBody, result)
	}()
	packet, ok := p.(*common.HeadV3Packet)
	if !ok {
		infoLog.Println("Convert to HeadV3Packet failed")
		result = uint16(ht_favorite.RET_CODE_RET_INPUT_PARAM_ERR)
		return false
	}
	head, err := packet.GetHead()
	if err != nil {
		infoLog.Println("ProcData Get head faild")
		result = uint16(ht_favorite.RET_CODE_RET_INTERNAL_ERR)
		return false
	}

	reqBody := &ht_favorite.ReqBody{}
	err = proto.Unmarshal(packet.GetBody(), reqBody)
	if err != nil {
		infoLog.Println("proto Unmarshal failed")
		result = uint16(ht_favorite.RET_CODE_RET_PB_ERR)
		return false
	}

	getReqBody := reqBody.GetGetReqbody()
	if getReqBody == nil {
		infoLog.Println("GetGetReqbody() failked")
		result = uint16(ht_favorite.RET_CODE_RET_PB_ERR)
		return false
	}
	infoLog.Printf("uid[%v],index[%v], last_ts[%v]", packet.GetFromUid(), getReqBody.GetIndex(), getReqBody.GetLastTs())

	// 第一步，比较时间戳，服务端的时间戳可以从objectid中获取
	ts := time.Now().Unix()
	retTs := new(FavorTs)
	bTsExist := false
	mongoConn2 := mongoSess.DB(favorDB).C(favorTsTable)
	err = mongoConn2.Find(bson.M{"uid": head.From}).One(retTs)
	if err != nil {
		bTsExist = false
		infoLog.Println("Find failed!")
		//如果找不到时间戳，就认为时间戳是0，还没有初始化过
	} else {
		bTsExist = true
		if retTs.LastTs == getReqBody.GetLastTs() {
			result = uint16(ht_favorite.RET_CODE_RET_NOT_CHANGE)
			return true
		}
	}

	// 第二步，查询记录，使用limit控制每页条数
	retRecList := []FavorRecord{}
	mongoConn := mongoSess.DB(favorDB).C(favorRecordTable)
	if getReqBody.GetIndex() == "" {
		err = mongoConn.Find(bson.M{"uid": head.From}).Sort("-_id").Limit(pageSize).All(&retRecList)
	} else {
		obid := bson.ObjectIdHex(getReqBody.GetIndex())
		err = mongoConn.Find(bson.M{"uid": head.From, "_id": bson.M{"$lt": obid}}).Sort("-_id").Limit(pageSize).All(&retRecList)
	}
	if err != nil {
		infoLog.Println("Insert failed")
	}

	//step3，update last_ts
	rspBody.GetRspbody = new(ht_favorite.GetRspBody)
	if !bTsExist {
		rspBody.GetRspbody.LastTs = proto.Uint32(uint32(ts))
	} else {
		rspBody.GetRspbody.LastTs = proto.Uint32(retTs.LastTs)
	}
	if len(retRecList) == pageSize {
		rspBody.GetRspbody.More = proto.Uint32(1)
	} else {
		rspBody.GetRspbody.More = proto.Uint32(0)
	}
	if len(retRecList) == 0 {
		result = uint16(ht_favorite.RET_CODE_RET_NO_MORE)
	} else {
		rspBody.GetRspbody.Index = proto.String(retRecList[len(retRecList)-1].Id.String())
		rspBody.GetRspbody.ContentList = make([]*ht_favorite.FavoriteContent, len(retRecList))
		for i := 0; i < len(retRecList); i++ {
			rspBody.GetRspbody.ContentList[i] = new(ht_favorite.FavoriteContent)
			rspBody.GetRspbody.ContentList[i].Obid = proto.String(retRecList[i].Id.String())
			rspBody.GetRspbody.ContentList[i].SrcUid = proto.Uint32(retRecList[i].SrcUid)
			rspBody.GetRspbody.ContentList[i].Type = new(ht_favorite.TYPE_FAVORATE)
			*rspBody.GetRspbody.ContentList[i].Type = ht_favorite.TYPE_FAVORATE(retRecList[i].Type)
			rspBody.GetRspbody.ContentList[i].Text = []byte(retRecList[i].Text)
			rspBody.GetRspbody.ContentList[i].Voice = new(ht_favorite.VoiceBody)
			rspBody.GetRspbody.ContentList[i].Voice.Url = []byte(retRecList[i].Voice.Url)
			rspBody.GetRspbody.ContentList[i].Voice.Duration = proto.Uint32(retRecList[i].Voice.Duration)
			rspBody.GetRspbody.ContentList[i].Voice.Size = proto.Uint32(retRecList[i].Voice.Size)
			rspBody.GetRspbody.ContentList[i].Image = new(ht_favorite.ImageBody)
			rspBody.GetRspbody.ContentList[i].Image.ThumbUrl = []byte(retRecList[i].Image.ThumbUrl)
			rspBody.GetRspbody.ContentList[i].Image.BigUrl = []byte(retRecList[i].Image.BigUrl)
			rspBody.GetRspbody.ContentList[i].Image.Width = proto.Uint32(retRecList[i].Image.Width)
			rspBody.GetRspbody.ContentList[i].Image.Height = proto.Uint32(retRecList[i].Image.Height)
			rspBody.GetRspbody.ContentList[i].Mid = []byte(retRecList[i].Mid)
			rspBody.GetRspbody.ContentList[i].Tags = retRecList[i].Tags
		}
	}

	return true
}

func procUpdateTagReq(c *gotcp.Conn, p gotcp.Packet) bool {
	// parse packet
	rspBody := new(ht_favorite.RspBody)

	result := uint16(ht_favorite.RET_CODE_RET_SUCCESS)
	var head *common.HeadV3
	defer func() {
		SendRsp(c, head, rspBody, result)
	}()
	packet, ok := p.(*common.HeadV3Packet)
	if !ok {
		infoLog.Println("Convert to HeadV3Packet failed")
		result = uint16(ht_favorite.RET_CODE_RET_INPUT_PARAM_ERR)
		return false
	}
	head, err := packet.GetHead()
	if err != nil {
		infoLog.Println("ProcData Get head faild")
		result = uint16(ht_favorite.RET_CODE_RET_INTERNAL_ERR)
		return false
	}
	reqBody := new(ht_favorite.ReqBody)
	err = proto.Unmarshal(packet.GetBody(), reqBody)
	if err != nil {
		infoLog.Println("proto Unmarshal failed")
		result = uint16(ht_favorite.RET_CODE_RET_PB_ERR)
		return false
	}
	updateTagReqBody := reqBody.GetUpdateTagReqbody()
	if updateTagReqBody == nil {
		infoLog.Println("GetAddReqbody() failked")
		result = uint16(ht_favorite.RET_CODE_RET_PB_ERR)
		return false
	}
	infoLog.Printf("uid[%v], obid[%s], new_tag[%s]\n", head.From, string(updateTagReqBody.Obid), updateTagReqBody.Tags)
	if string(updateTagReqBody.GetObid()) == "" ||
		len(updateTagReqBody.GetTags()) == 0 {
		infoLog.Printf("obid[%s] is nil or tags[%v]\n", updateTagReqBody.GetObid(), updateTagReqBody.GetTags())
		result = uint16(ht_favorite.RET_CODE_RET_INTERNAL_ERR)
		return false
	}
	obid := bson.ObjectIdHex(string(updateTagReqBody.GetObid()))
	retRecord := FavorRecord{}
	mongoConn := mongoSess.DB(favorDB).C(favorRecordTable)
	err = mongoConn.FindId(obid).One(&retRecord)
	// 这里要判断是否找到了
	if err != nil {
		infoLog.Println("find failed")
		result = uint16(ht_favorite.RET_CODE_RET_REPEAT_ADD)
		return true
	}

	retRecord.Tags = updateTagReqBody.Tags
	_, err = mongoConn.Upsert(bson.M{"_id": obid}, retRecord)
	if err != nil {
		infoLog.Println("Upsert favor record failed")
		result = uint16(ht_favorite.RET_CODE_RET_DB_ERR)
		return false
	}
	ts := time.Now().Unix()
	recordTs := new(FavorTs)
	recordTs.Uid = head.From
	recordTs.LastTs = uint32(ts)
	mongoConn2 := mongoSess.DB(favorDB).C(favorTsTable)
	_, err = mongoConn2.Upsert(bson.M{"uid": recordTs.Uid}, recordTs)
	if err != nil {
		infoLog.Println("Upsert ts failed")
		result = uint16(ht_favorite.RET_CODE_RET_DB_ERR)
		return true
	}

	//send response packet
	rspBody.UpdateTagRspbody = new(ht_favorite.UpdateTagRspBody)
	rspBody.UpdateTagRspbody.LastTs = proto.Uint32(uint32(ts))

	return true
}

func procSearchTagReq(c *gotcp.Conn, p gotcp.Packet) bool {
	// parse packet
	rspBody := new(ht_favorite.RspBody)
	result := uint16(ht_favorite.RET_CODE_RET_SUCCESS)
	var head *common.HeadV3
	defer func() {
		SendRsp(c, head, rspBody, result)
	}()
	packet, ok := p.(*common.HeadV3Packet)
	if !ok {
		infoLog.Println("Convert to HeadV3Packet failed")
		result = uint16(ht_favorite.RET_CODE_RET_INPUT_PARAM_ERR)
		return false
	}
	head, err := packet.GetHead()
	if err != nil {
		infoLog.Println("ProcData Get head faild")
		result = uint16(ht_favorite.RET_CODE_RET_INTERNAL_ERR)
		return false
	}

	reqBody := &ht_favorite.ReqBody{}
	err = proto.Unmarshal(packet.GetBody(), reqBody)
	if err != nil {
		infoLog.Println("proto Unmarshal failed")
		result = uint16(ht_favorite.RET_CODE_RET_PB_ERR)
		return false
	}

	searchTagReqBody := reqBody.GetSearchTagReqbody()
	if searchTagReqBody == nil {
		infoLog.Println("GetSearchTagReqbody() failked")
		result = uint16(ht_favorite.RET_CODE_RET_PB_ERR)
		return false
	}
	infoLog.Printf("uid[%v],search_tag[%v], last_ts[%v], index[%s]", packet.GetFromUid(), searchTagReqBody.GetTag(), searchTagReqBody.GetLastTs(), searchTagReqBody.GetIndex())

	// 第一步，比较时间戳，服务端的时间戳可以从objectid中获取
	ts := time.Now().Unix()
	retTs := new(FavorTs)
	bTsExist := false
	mongoConn2 := mongoSess.DB(favorDB).C(favorTsTable)
	err = mongoConn2.Find(bson.M{"uid": head.From}).One(retTs)
	if err != nil {
		bTsExist = false
		infoLog.Println("Find failed!")
		//如果找不到时间戳，就认为时间戳是0，还没有初始化过
	} else {
		bTsExist = true
		if retTs.LastTs == searchTagReqBody.GetLastTs() {
			result = uint16(ht_favorite.RET_CODE_RET_NOT_CHANGE)
			return true
		}
	}

	// 第二步，查询记录，使用limit控制每页条数
	retRecList := []FavorRecord{}
	mongoConn := mongoSess.DB(favorDB).C(favorRecordTable)
	if searchTagReqBody.GetIndex() == "" {
		//err = mongoConn.Find(bson.M{"uid": head.From, "$text": bson.M{"$search": searchTagReqBody.GetTag()}}).Sort("-_id").Limit(pageSize).All(&retRecList)
		err = mongoConn.Find(bson.M{"uid": head.From, "tags": searchTagReqBody.GetTag()}).Sort("-_id").Limit(pageSize).All(&retRecList)
	} else {
		obid := bson.ObjectIdHex(searchTagReqBody.GetIndex())
		//err = mongoConn.Find(bson.M{"uid": head.From, "_id": bson.M{"$lt": obid}, "$text": bson.M{"$search": searchTagReqBody.GetTag()}}).Sort("-_id").Limit(pageSize).All(&retRecList)
		err = mongoConn.Find(bson.M{"uid": head.From, "_id": bson.M{"$lt": obid}, "tags": searchTagReqBody.GetTag()}).Sort("-_id").Limit(pageSize).All(&retRecList)
	}
	if err != nil {
		infoLog.Println("Find failed")
	}

	//step3，update last_ts
	rspBody.SearchTagRspbody = new(ht_favorite.SearchTagRspBody)
	if !bTsExist {
		rspBody.SearchTagRspbody.LastTs = proto.Uint32(uint32(ts))
	} else {
		rspBody.SearchTagRspbody.LastTs = proto.Uint32(uint32(retTs.LastTs))
	}
	if len(retRecList) == pageSize {
		rspBody.SearchTagRspbody.More = proto.Uint32(1)
	} else {
		rspBody.SearchTagRspbody.More = proto.Uint32(0)
	}
	if len(retRecList) == 0 {
		result = uint16(ht_favorite.RET_CODE_RET_NO_MORE)
	} else {
		rspBody.SearchTagRspbody.Index = proto.String(retRecList[len(retRecList)-1].Id.String())
		rspBody.SearchTagRspbody.ContentList = make([]*ht_favorite.FavoriteContent, len(retRecList))
		for i := 0; i < len(retRecList); i++ {
			rspBody.SearchTagRspbody.ContentList[i] = new(ht_favorite.FavoriteContent)
			rspBody.SearchTagRspbody.ContentList[i].Obid = proto.String(retRecList[i].Id.String())
			rspBody.SearchTagRspbody.ContentList[i].SrcUid = proto.Uint32(retRecList[i].SrcUid)
			rspBody.SearchTagRspbody.ContentList[i].Type = new(ht_favorite.TYPE_FAVORATE)
			*rspBody.SearchTagRspbody.ContentList[i].Type = ht_favorite.TYPE_FAVORATE(retRecList[i].Type)
			rspBody.SearchTagRspbody.ContentList[i].Text = []byte(retRecList[i].Text)
			rspBody.SearchTagRspbody.ContentList[i].Voice = new(ht_favorite.VoiceBody)
			rspBody.SearchTagRspbody.ContentList[i].Voice.Url = []byte(retRecList[i].Voice.Url)
			rspBody.SearchTagRspbody.ContentList[i].Voice.Duration = proto.Uint32(retRecList[i].Voice.Duration)
			rspBody.SearchTagRspbody.ContentList[i].Voice.Size = proto.Uint32(retRecList[i].Voice.Size)
			rspBody.SearchTagRspbody.ContentList[i].Image = new(ht_favorite.ImageBody)
			rspBody.SearchTagRspbody.ContentList[i].Image.ThumbUrl = []byte(retRecList[i].Image.ThumbUrl)
			rspBody.SearchTagRspbody.ContentList[i].Image.BigUrl = []byte(retRecList[i].Image.BigUrl)
			rspBody.SearchTagRspbody.ContentList[i].Image.Width = proto.Uint32(retRecList[i].Image.Width)
			rspBody.SearchTagRspbody.ContentList[i].Image.Height = proto.Uint32(retRecList[i].Image.Height)
			rspBody.SearchTagRspbody.ContentList[i].Mid = []byte(retRecList[i].Mid)
			rspBody.SearchTagRspbody.ContentList[i].Tags = retRecList[i].Tags
		}
	}

	return true
}

func procSearchTextReq(c *gotcp.Conn, p gotcp.Packet) bool {
	// parse packet
	rspBody := new(ht_favorite.RspBody)
	result := uint16(ht_favorite.RET_CODE_RET_SUCCESS)
	var head *common.HeadV3
	defer func() {
		SendRsp(c, head, rspBody, result)
	}()
	packet, ok := p.(*common.HeadV3Packet)
	if !ok {
		infoLog.Println("Convert to HeadV3Packet failed")
		result = uint16(ht_favorite.RET_CODE_RET_INPUT_PARAM_ERR)
		return false
	}
	head, err := packet.GetHead()
	if err != nil {
		infoLog.Println("ProcData Get head faild")
		result = uint16(ht_favorite.RET_CODE_RET_INTERNAL_ERR)
		return false
	}

	reqBody := &ht_favorite.ReqBody{}
	err = proto.Unmarshal(packet.GetBody(), reqBody)
	if err != nil {
		infoLog.Println("proto Unmarshal failed")
		result = uint16(ht_favorite.RET_CODE_RET_PB_ERR)
		return false
	}

	searchTextReqBody := reqBody.GetSearchTextReqbody()
	if searchTextReqBody == nil {
		infoLog.Println("GetSearchTagReqbody() failked")
		result = uint16(ht_favorite.RET_CODE_RET_PB_ERR)
		return false
	}
	infoLog.Printf("uid[%v],search_text[%v], last_ts[%v]", packet.GetFromUid(), searchTextReqBody.GetText(), searchTextReqBody.GetLastTs())

	// 第一步，比较时间戳，服务端的时间戳可以从objectid中获取
	ts := time.Now().Unix()
	retTs := new(FavorTs)
	bTsExist := false
	mongoConn2 := mongoSess.DB(favorDB).C(favorTsTable)
	err = mongoConn2.Find(bson.M{"uid": head.From}).One(retTs)
	if err != nil {
		bTsExist = false
		infoLog.Println("Find failed!")
		//如果找不到时间戳，就认为时间戳是0，还没有初始化过
	} else {
		bTsExist = true
		if retTs.LastTs == searchTextReqBody.GetLastTs() {
			result = uint16(ht_favorite.RET_CODE_RET_NOT_CHANGE)
			return true
		}
	}

	// 第二步，查询记录，使用limit控制每页条数
	retRecList := []FavorRecord{}
	mongoConn := mongoSess.DB(favorDB).C(favorRecordTable)
	if searchTextReqBody.GetIndex() == "" {
		err = mongoConn.Find(bson.M{"uid": head.From, "$text": bson.M{"$search": searchTextReqBody.GetText()}}).Sort("-_id").Limit(pageSize).All(&retRecList)
	} else {
		obid := bson.ObjectIdHex(searchTextReqBody.GetIndex())
		err = mongoConn.Find(bson.M{"uid": head.From, "_id": bson.M{"$lt": obid}, "$text": bson.M{"$search": searchTextReqBody.GetText()}}).Sort("-_id").Limit(pageSize).All(&retRecList)
	}
	if err != nil {
		infoLog.Println("Insert failed")
	}

	//step3，update last_ts
	rspBody.SearchTextRspbody = new(ht_favorite.SearchTextRspBody)
	if !bTsExist {
		rspBody.SearchTextRspbody.LastTs = proto.Uint32(uint32(ts))
	} else {
		rspBody.SearchTextRspbody.LastTs = proto.Uint32(uint32(retTs.LastTs))
	}
	if len(retRecList) == pageSize {
		rspBody.SearchTextRspbody.More = proto.Uint32(1)
	} else {
		rspBody.SearchTextRspbody.More = proto.Uint32(0)
	}
	if len(retRecList) == 0 {
		result = uint16(ht_favorite.RET_CODE_RET_NO_MORE)
	} else {
		rspBody.SearchTextRspbody.Index = proto.String(retRecList[len(retRecList)-1].Id.String())
		rspBody.SearchTextRspbody.ContentList = make([]*ht_favorite.FavoriteContent, len(retRecList))
		for i := 0; i < len(retRecList); i++ {
			rspBody.SearchTextRspbody.ContentList[i] = new(ht_favorite.FavoriteContent)
			rspBody.SearchTextRspbody.ContentList[i].Obid = proto.String(retRecList[i].Id.String())
			rspBody.SearchTextRspbody.ContentList[i].SrcUid = proto.Uint32(retRecList[i].SrcUid)
			rspBody.SearchTextRspbody.ContentList[i].Type = new(ht_favorite.TYPE_FAVORATE)
			*rspBody.SearchTextRspbody.ContentList[i].Type = ht_favorite.TYPE_FAVORATE(retRecList[i].Type)
			rspBody.SearchTextRspbody.ContentList[i].Text = []byte(retRecList[i].Text)
			rspBody.SearchTextRspbody.ContentList[i].Voice = new(ht_favorite.VoiceBody)
			rspBody.SearchTextRspbody.ContentList[i].Voice.Url = []byte(retRecList[i].Voice.Url)
			rspBody.SearchTextRspbody.ContentList[i].Voice.Duration = proto.Uint32(retRecList[i].Voice.Duration)
			rspBody.SearchTextRspbody.ContentList[i].Voice.Size = proto.Uint32(retRecList[i].Voice.Size)
			rspBody.SearchTextRspbody.ContentList[i].Image = new(ht_favorite.ImageBody)
			rspBody.SearchTextRspbody.ContentList[i].Image.ThumbUrl = []byte(retRecList[i].Image.ThumbUrl)
			rspBody.SearchTextRspbody.ContentList[i].Image.BigUrl = []byte(retRecList[i].Image.BigUrl)
			rspBody.SearchTextRspbody.ContentList[i].Image.Width = proto.Uint32(retRecList[i].Image.Width)
			rspBody.SearchTextRspbody.ContentList[i].Image.Height = proto.Uint32(retRecList[i].Image.Height)
			rspBody.SearchTextRspbody.ContentList[i].Mid = []byte(retRecList[i].Mid)
			rspBody.SearchTextRspbody.ContentList[i].Tags = retRecList[i].Tags
		}
	}

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
	//runtime.GOMAXPROCS(runtime.NumCPU())
	runtime.GOMAXPROCS(1)

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
	infoLog.SetFlags(infoLog.Flags() | log.LstdFlags | log.Lshortfile)

	// 创建mongodb对象
	mongo_url := cfg.Section("MONGO").Key("url").MustString("localhost")
	infoLog.Println(mongo_url)
	mongoSess, err = mgo.Dial(mongo_url)
	if err != nil {
		log.Fatalln("connect mongodb failed")
		return
	}
	defer mongoSess.Close()
	mongoSess.SetMode(mgo.Monotonic, true)

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
