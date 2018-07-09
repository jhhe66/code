package common

import (
	"github.com/gansidui/gotcp"
	"github.com/gansidui/gotcp/tcpfw/include/ht_push"
	"github.com/golang/protobuf/proto"
)

// Conn exposes a set of callbacks for the various events that occur on a connection
type PushServerApi struct {
	TCPApi
}

func NewPushServerApi(destIp string, destPort int, protocol gotcp.Protocol) *PushServerApi {
	return &PushServerApi{
		TCPApi{
			DestIp:       destIp,
			DestPort:     destPort,
			ReadTimeOut:  1,
			WriteTimeOut: 1,
			protocol:     protocol,
		},
	}
}

func (c *PushServerApi) SendPushPacket(from uint32,
	to uint32,
	chatType uint32,
	pushType uint32,
	roomId uint32,
	nickName []byte,
	content []byte,
	sound uint32,
	light uint32,
	messageId []byte,
	actionId uint32,
	needRsp uint32) (ret uint16, err error) {
	reqHead := &HeadV3{Flag: 0xF2, From: from, To: to, Cmd: uint16(ht_push.CMD_TYPE_CMD_S2S_MESSAGE_PUSH), SysType: uint16(ht_push.SYS_TYPE_SYS_IM_SERVER)}
	var reqBody ht_push.ReqBody
	subReqBody := &ht_push.PushMsgReqBody{
		ChatType:  &chatType,
		PushType:  &pushType,
		RoomId:    &roomId,
		Nickname:  nickName,
		Content:   content,
		Sound:     &sound,
		Light:     &light,
		MessageId: messageId,
		ActionId:  &actionId,
		NeedRsp:   &needRsp,
	}
	reqBody.PushMsgReqbody = subReqBody
	payLoad, err := proto.Marshal(&reqBody)
	if err != nil {
		ret = uint16(ht_push.RET_CODE_RET_PB_ERR)
		return ret, err
	}
	reqHead.Len = uint32(PacketV3HeadLen) + uint32(len(payLoad)) + 1
	buf := make([]byte, reqHead.Len)
	buf[0] = HTV3MagicBegin
	err = SerialHeadV3ToSlice(reqHead, buf[1:])
	if err != nil {
		ret = uint16(ht_push.RET_CODE_RET_INTERNAL_ERR)
		return ret, err
	}
	copy(buf[PacketV3HeadLen:], payLoad)
	buf[reqHead.Len-1] = HTV3MagicEnd

	packet := NewHeadV3Packet(buf)
	err = c.WritePacket(packet)
	if err != nil { // 发送失败
		ret = uint16(ht_push.RET_CODE_RET_INTERNAL_ERR)
		return ret, err
	}
	rsp, err := c.ReadPacket()
	if err != nil {
		ret = uint16(ht_push.RET_CODE_RET_INTERNAL_ERR)
		return ret, err
	}

	rspPacket := rsp.(*HeadV3Packet)
	rspHead, err := rspPacket.GetHead()
	if err != nil {
		ret = uint16(ht_push.RET_CODE_RET_INTERNAL_ERR)
		return ret, err
	}
	ret = rspHead.Ret
	return ret, err
}
