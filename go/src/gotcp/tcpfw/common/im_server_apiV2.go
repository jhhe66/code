package common

import (
	"time"

	"github.com/gansidui/gotcp"
)

// Conn exposes a set of callbacks for the various events that occur on a connection
type ImServerApiV2 struct {
	pool *Pool
}

func newImServerPool(ip, port string, readTimeout, writeTimeout time.Duration, maxConn int, proto gotcp.Protocol) *Pool {
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

// NewImServerApiV2 的时候需要指定协议protocol 所以在发送HeadV3报文和XTPacket 报文必须要实例化不同的对象来处理
func NewImServerApiV2(ip, port string, readTimeout, writeTimeout time.Duration, protocol gotcp.Protocol, maxConn int) *ImServerApiV2 {
	pool := newImServerPool(ip, port, readTimeout, writeTimeout, maxConn, protocol)
	return &ImServerApiV2{
		pool: pool,
	}
}

func (c *ImServerApiV2) SendPacketToIM(terminalType uint8, cmd uint16, seq uint16, from uint32, to uint32, payLoad []byte) (ret uint16, err error) {
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

func (c *ImServerApiV2) SendPacket(head *HeadV3, payLoad []byte) (ret uint16, err error) {
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
	if err != nil { // 发送失败
		ret = uint16(CRetSendFailed)
		return ret, err
	}

	rspPacket := rsp.(*HeadV3Packet)
	rspHead, err := rspPacket.GetHead()
	if err != nil {
		ret := uint16(CRetUnMarshallFailed)
		return ret, err
	}
	ret = rspHead.Ret
	return ret, nil
}

func (c *ImServerApiV2) SendXTPacket(head *XTHead, payLoad []byte) (err error) {
	head.Len = uint32(len(payLoad)) //
	buf := make([]byte, XTHeadLen+head.Len)
	err = SerialXTHeadToSlice(head, buf[:])
	if err != nil {
		return err
	}
	copy(buf[XTHeadLen:], payLoad) // return code
	reqPacket := NewXTHeadPacket(buf)

	conn := c.pool.Get()
	defer conn.Close()

	rsp, err := conn.Do(reqPacket)
	if err != nil { // 发送失败
		return err
	}

	rspPacket := rsp.(*XTHeadPacket)
	rspHead, err := rspPacket.GetHead()
	if err != nil {
		return err
	}
	if head.Cmd == rspHead.Cmd {
		return nil
	} else {
		return nil
	}
}
