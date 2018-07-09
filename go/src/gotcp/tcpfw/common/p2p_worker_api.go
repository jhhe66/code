package common

import (
	"github.com/gansidui/gotcp"
	"github.com/gansidui/gotcp/libcomm"
)

// Conn exposes a set of callbacks for the various events that occur on a connection
type P2PWorkerApi struct {
	TCPApi
}

func NewP2PWorkerApi(destIp string, destPort int, protocol gotcp.Protocol) *P2PWorkerApi {
	return &P2PWorkerApi{
		TCPApi{
			DestIp:       destIp,
			DestPort:     destPort,
			ReadTimeOut:  1,
			WriteTimeOut: 1,
			protocol:     protocol,
		},
	}
}

func (c *P2PWorkerApi) SendPacket(head *HeadV3, payLoad []uint8) (ret uint16, err error) {
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
		attr := "P2PWorkerApi/inter_writepacket_fail"
		libcomm.AttrAdd(attr, 1)
		ret = uint16(CRetSendFailed)
		return ret, err
	}
	p, err := c.ReadPacket()
	if err != nil {
		attr := "P2PWorkerApi/inter_readpacket_fail"
		libcomm.AttrAdd(attr, 1)
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

func (c *P2PWorkerApi) SendAndRecvPacket(head *HeadV3, payLoad []uint8) (packet gotcp.Packet, err error) {
	head.Len = uint32(PacketV3HeadLen) + uint32(len(payLoad)) + 1
	buf := make([]byte, head.Len)
	buf[0] = HTV3MagicBegin
	err = SerialHeadV3ToSlice(head, buf[1:])
	if err != nil {
		return nil, err
	}
	copy(buf[PacketV3HeadLen:], payLoad)
	buf[head.Len-1] = HTV3MagicEnd

	headV3Packet := NewHeadV3Packet(buf)
	err = c.WritePacket(headV3Packet)
	if err != nil { // 发送失败
		attr := "P2PWorkerApi/inter_writepacket_fail"
		libcomm.AttrAdd(attr, 1)
		return nil, err
	}
	packet, err = c.ReadPacket()
	if err != nil {
		attr := "P2PWorkerApi/inter_readpacket_fail"
		libcomm.AttrAdd(attr, 1)
		return nil, err
	}
	return packet, nil
}
