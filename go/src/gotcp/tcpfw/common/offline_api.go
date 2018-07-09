package common

import (
	"time"

	"github.com/gansidui/gotcp"
	"github.com/gansidui/gotcp/tcpfw/include/ht_offline"
	"github.com/golang/protobuf/proto"
)

// Conn exposes a set of callbacks for the various events that occur on a connection
type OfflineApi struct {
	TCPApi
}

func NewOfflineApi(destIp string, destPort int, protocol gotcp.Protocol) *OfflineApi {
	return &OfflineApi{
		TCPApi{
			DestIp:       destIp,
			DestPort:     destPort,
			ReadTimeOut:  1,
			WriteTimeOut: 1,
			protocol:     protocol,
		},
	}
}

func (c *OfflineApi) SendPacketWithHeadV3(headV3 *HeadV3, payLoad []uint8, push []byte) (ret uint16, err error) {
	// 存离线消息的包（HTHeadV3的包）
	newHead := new(HeadV3)
	*newHead = *headV3
	newHead.Len = uint32(PacketV3HeadLen) + uint32(len(payLoad)) + 1
	buf := make([]byte, newHead.Len)
	buf[0] = HTV3MagicBegin
	err = SerialHeadV3ToSlice(newHead, buf[1:])
	if err != nil {
		ret = uint16(CRetMarshallFailed)
		return ret, err
	}
	copy(buf[PacketV3HeadLen:], payLoad)
	buf[newHead.Len-1] = HTV3MagicEnd

	// 构造HTHeadV2 报文
	headV2 := &HeadV2{Cmd: uint32(ht_offline.CMD_TYPE_CMD_SAVE_OFFLINE_MSG),
		Uid:     headV3.To,
		SysType: uint16(ht_offline.SYS_TYPE_SYS_VOIP_SERVER),
	}

	reqBody := new(ht_offline.ReqBody)
	reqBody.SaveOfflineMsgReqbody = new(ht_offline.SaveOfflineMsgReqBody)
	subReqBody := reqBody.GetSaveOfflineMsgReqbody()
	subReqBody.FromUid = new(uint32)
	*(subReqBody.FromUid) = headV3.From

	subReqBody.ToUid = new(uint32)
	*(subReqBody.ToUid) = headV3.To

	subReqBody.Cmd = new(uint32)
	*(subReqBody.Cmd) = uint32(headV3.Cmd)

	subReqBody.Format = new(uint32)
	*(subReqBody.Format) = 10 // HTHeadV3的消息格式用10

	subReqBody.Content = make([]byte, newHead.Len)
	copy(subReqBody.Content, buf)

	subReqBody.MsgTime = new(uint32)
	*(subReqBody.MsgTime) = uint32(time.Now().Unix()) // 添加时间戳

	subReqBody.PushType = new(uint32)
	if len(push) > 0 {
		*(subReqBody.PushType) = uint32(ht_offline.SaveOfflineMsgReqBody_EM_PUSH_CONFIRM)
		subReqBody.PushPkg = make([]byte, len(push))
		copy(subReqBody.PushPkg, push)
	} else {
		*(subReqBody.PushType) = uint32(ht_offline.SaveOfflineMsgReqBody_EM_PUSH_DISCARD)
	}

	s, err := proto.Marshal(reqBody)
	if err != nil {
		ret = uint16(CRetMarshallFailed)
		return ret, err
	}

	headV2.Len = uint32(PacketV2HeadLen) + uint32(len(s)) + 1
	headV2Buf := make([]byte, headV2.Len)
	headV2Buf[0] = HTV2MagicBegin
	err = SerialHeadV2ToSlice(headV2, headV2Buf[1:])
	if err != nil {
		ret = uint16(CRetMarshallFailed)
		return ret, err
	}
	copy(headV2Buf[PacketV2HeadLen:], s)
	headV2Buf[headV2.Len-1] = HTV2MagicEnd

	packet := NewHeadV2Packet(headV2Buf)

	err = c.WritePacket(packet)
	if err != nil { // 发送失败
		ret = uint16(CRetSendFailed)
		return ret, err
	}

	rsp, err := c.ReadPacket()
	if err != nil {
		ret = uint16(CRetReadFailed)
		return ret, err
	}

	rspPacket := rsp.(*HeadV2Packet)
	rspHead, err := rspPacket.GetHead()
	if err != nil {
		ret := uint16(CRetUnMarshallFailed)
		return ret, err
	}
	ret = rspHead.Ret
	return ret, nil
}

func (c *OfflineApi) SendPacket(head *HeadV2, payLoad []uint8) (ret uint16, err error) {
	head.Len = uint32(PacketV2HeadLen) + uint32(len(payLoad)) + 1
	buf := make([]byte, head.Len)
	buf[0] = HTV2MagicBegin
	err = SerialHeadV2ToSlice(head, buf[1:])
	if err != nil {
		ret = uint16(CRetSendFailed)
		return ret, err
	}
	copy(buf[PacketV2HeadLen:], payLoad)
	buf[head.Len-1] = HTV2MagicEnd

	packet := NewHeadV2Packet(buf)

	err = c.WritePacket(packet)
	if err != nil { // 发送失败
		ret = uint16(CRetSendFailed)
		return ret, err
	}

	rsp, err := c.ReadPacket()
	if err != nil {
		ret = uint16(CRetReadFailed)
		return ret, err
	}

	rspPacket := rsp.(*HeadV2Packet)
	rspHead, err := rspPacket.GetHead()
	if err != nil {
		ret := uint16(CRetUnMarshallFailed)
		return ret, err
	}
	ret = rspHead.Ret
	return ret, nil
}
