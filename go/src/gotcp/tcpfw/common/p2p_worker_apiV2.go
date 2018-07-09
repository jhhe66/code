package common

import (
	"time"

	"github.com/gansidui/gotcp"
)

// Conn exposes a set of callbacks for the various events that occur on a connection
type P2PWorkerApiV2 struct {
	pool *Pool
}

func newWorkerPool(ip, port string, readTimeout, writeTimeout time.Duration, maxConn int, proto gotcp.Protocol) *Pool {
	return &Pool{
		MaxIdle:     maxConn,
		MaxActive:   maxConn,
		IdleTimeout: 240 * time.Second,
		Dial: func() (Conn, error) {
			c, err := Dial(ip, port, proto, DialReadTimeout(readTimeout), DialWriteTimeout(writeTimeout))
			if err != nil {
				return nil, err
			}
			return c, nil
		},
	}
}

func NewP2PWorkerApiV2(ip, port string, readTimeout, writeTimeout time.Duration, protocol gotcp.Protocol, maxConn int) *P2PWorkerApiV2 {
	pool := newWorkerPool(ip, port, readTimeout, writeTimeout, maxConn, protocol)
	return &P2PWorkerApiV2{
		pool: pool,
	}
}

func (c *P2PWorkerApiV2) SendPacket(head *HeadV3, payLoad []uint8) (ret uint16, err error) {
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
	conn := c.pool.Get()
	defer conn.Close()

	rsp, err := conn.Do(headV3Packet)
	rspPacket, ok := rsp.(*HeadV3Packet)
	if !ok {
		ret := uint16(CRetSendFailed)
		return ret, err
	}
	rspHead, err := rspPacket.GetHead()
	if err != nil {
		ret := uint16(CRetUnMarshallFailed)
		return ret, err
	}
	ret = rspHead.Ret
	return ret, nil
}

func (c *P2PWorkerApiV2) SendAndRecvPacket(head *HeadV3, payLoad []uint8) (packet gotcp.Packet, err error) {
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
	conn := c.pool.Get()
	defer conn.Close()

	packet, err = conn.Do(headV3Packet)
	return packet, nil
}
