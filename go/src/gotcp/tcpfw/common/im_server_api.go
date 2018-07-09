package common

import (
	"github.com/gansidui/gotcp"
)

// Conn exposes a set of callbacks for the various events that occur on a connection
type ImServerApi struct {
	TCPApi
}

func NewImServerApi(destIp string, destPort int, protocol gotcp.Protocol) *ImServerApi {
	return &ImServerApi{
		TCPApi{
			DestIp:       destIp,
			DestPort:     destPort,
			ReadTimeOut:  1,
			WriteTimeOut: 1,
			protocol:     protocol,
		},
	}
}

func (c *ImServerApi) SendPacketToIM(terminalType uint8,
	cmd uint16,
	seq uint16,
	from uint32,
	to uint32,
	payLoad []uint8) (ret uint16, err error) {
	head := &HeadV3{Flag: uint8(CServToServ),
		CryKey:   uint8(CNoneKey),
		TermType: terminalType,
		Cmd:      cmd,
		Seq:      seq,
		From:     from,
		To:       to,
	}
	ret, err = c.SendPacket(head, payLoad)

	// 如果连接已经断开尝试重发
	return
}

func (c *ImServerApi) SendPacket(head *HeadV3, payLoad []uint8) (ret uint16, err error) {
	head.Len = uint32(PacketV3HeadLen) + uint32(len(payLoad)) + 1
	buf := make([]byte, head.Len)
	buf[0] = HTV3MagicBegin
	err = SerialHeadV3ToSlice(head, buf[1:])
	if err != nil {
		ret = uint16(CRetSendFailed)
		return ret, err
	}
	copy(buf[PacketV3HeadLen:], payLoad)
	buf[head.Len-1] = HTV3MagicEnd

	headV3Packet := NewHeadV3Packet(buf)
	err = c.WritePacket(headV3Packet)
	if err != nil { // 发送失败
		ret = uint16(CRetSendFailed)
		return ret, err
	}
	p, err := c.ReadPacket()
	if err != nil {
		ret = uint16(CRetReadFailed)
		return ret, err
	}

	rspPacket := p.(*HeadV3Packet)
	rspHead, err := rspPacket.GetHead()
	if err != nil {
		ret := uint16(CRetUnMarshallFailed)
		return ret, err
	}
	ret = rspHead.Ret
	return ret, nil
}
