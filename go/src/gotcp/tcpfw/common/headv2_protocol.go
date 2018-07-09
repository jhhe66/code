package common

import (
	"encoding/binary"
	"errors"
	"io"
	"net"

	"github.com/gansidui/gotcp"
)

const (
	HeadV2Len       = 30
	PacketV2HeadLen = 31
	EmptyPacktV2Len = 32
	PacketV2Limit   = 128 * 1024
	HTV2MagicBegin  = 0x08
	HTV2MagicEnd    = 0x09
)

// Error type
var (
	V2ErrShortLen   = errors.New("use byte[] lne is not enough")
	V2ErrBeginOrEnd = errors.New("must begin with 0x08 and end with 0x09")
	V2ErrLenErr     = errors.New("Len error")
)

type HeadV2 struct {
	Version  uint8    // 版本号  VER_MMEDIA = 4
	Cmd      uint32   // 命令字
	Seq      uint32   // 序列号
	Ret      uint16   // 返回码
	Reserved uint8    // 保留字段
	Len      uint32   // 报文总长度
	Uid      uint32   // 用户id
	SysType  uint16   // 请求来源系统
	Echo     [8]uint8 // 回带字段
}

func NewHeadV2(buf []byte) (head *HeadV2, err error) {
	if len(buf) < HeadV2Len {
		return nil, V2ErrShortLen
	}
	head = new(HeadV2)
	head.Version = buf[0]
	head.Cmd = binary.BigEndian.Uint32(buf[1:5])
	head.Seq = binary.BigEndian.Uint32(buf[5:9])
	head.Ret = binary.BigEndian.Uint16(buf[9:11])
	head.Reserved = buf[11]
	head.Len = binary.BigEndian.Uint32(buf[12:16])
	head.Uid = binary.BigEndian.Uint32(buf[16:20])
	head.SysType = binary.BigEndian.Uint16(buf[20:22])
	for i := 0; i < 8; i++ {
		head.Echo[i] = buf[22+i]
	}
	return head, nil
}

func SerialHeadV2ToSlice(head *HeadV2, buf []byte) (err error) {
	if len(buf) < HeadV2Len {
		return V2ErrShortLen
	}
	buf[0] = head.Version
	binary.BigEndian.PutUint32(buf[1:5], head.Cmd)
	binary.BigEndian.PutUint32(buf[5:9], head.Seq)
	binary.BigEndian.PutUint16(buf[9:11], head.Ret)
	buf[11] = head.Reserved
	binary.BigEndian.PutUint32(buf[12:16], head.Len)
	binary.BigEndian.PutUint32(buf[16:20], head.Uid)
	binary.BigEndian.PutUint16(buf[20:22], head.Ret)
	for i := 0; i < 8; i++ {
		buf[22+i] = head.Echo[i]
	}
	return nil
}

// HeadV3Packet 格式如下
// 0x0a + HTHeadV3 + payload + 0x0b

type HeadV2Packet struct {
	buff []byte
}

func (this *HeadV2Packet) Serialize() []byte {
	return this.buff
}

// index:0 SOH filed
func (this *HeadV2Packet) GetSoh() uint8 {
	return this.buff[0]
}

// index:1 Version field
func (this *HeadV2Packet) GetVersion() uint8 {
	return this.buff[1]
}

// index:2 Command field
func (this *HeadV2Packet) GetCommand() uint32 {
	return binary.BigEndian.Uint32(this.buff[2:6])
}

// index:3 Sequence field
func (this *HeadV2Packet) GetSeq() uint32 {
	return binary.BigEndian.Uint32(this.buff[6:10])
}

// index:4 Ret field
func (this *HeadV2Packet) GetRet() uint32 {
	return binary.BigEndian.Uint32(this.buff[10:12])
}

// index:5 Reserved field
func (this *HeadV2Packet) GetReserved() uint8 {
	return this.buff[12]
}

// index:6 Packet Length field length(SOH+HeadV2+PayLoad+EOT)
func (this *HeadV2Packet) GetLength() uint32 {
	return binary.BigEndian.Uint32(this.buff[13:17])
}

// index:7 uid field
func (this *HeadV2Packet) GetUid() uint32 {
	return binary.BigEndian.Uint32(this.buff[17:21])
}

// index:8 Packet come from field
func (this *HeadV2Packet) GetSystem() uint16 {
	return binary.BigEndian.Uint16(this.buff[21:23])
}

// index:9 Echo field
func (this *HeadV2Packet) GetEcho() []byte {
	return this.buff[23:31]
}

// index:10 Body filed
func (this *HeadV2Packet) GetBody() []byte {
	return this.buff[PacketV2HeadLen : len(this.buff)-1]
}

// index:11 EOT filed
func (this *HeadV2Packet) GetEot() uint8 {
	return this.buff[len(this.buff)-1]
}

func (this *HeadV2Packet) GetHead() (head *HeadV2, err error) {
	head, err = NewHeadV2(this.buff[1:])
	return
}

func (this *HeadV2Packet) CheckPacketV2Valid() (bool, error) {
	if this.GetLength() < PacketV2HeadLen {
		return false, V2ErrShortLen
	}

	if (this.GetSoh() != HTV2MagicBegin) || (this.GetEot() != HTV2MagicEnd) {
		return false, V2ErrBeginOrEnd
	}

	if this.GetLength() != uint32(len(this.Serialize())) {
		return false, V2ErrLenErr
	}

	return true, nil
}

func NewHeadV2Packet(buff []byte) *HeadV2Packet {
	p := &HeadV2Packet{}
	p.buff = buff
	return p
}

type HeadV2Protocol struct {
}

func (this *HeadV2Protocol) ReadPacket(conn *net.TCPConn) (gotcp.Packet, error) {

	var (
		headBytes []byte = make([]byte, PacketV2HeadLen)
		length    uint32
	)

	// read length
	if _, err := io.ReadFull(conn, headBytes); err != nil {
		return nil, err
	}

	head, err := NewHeadV2(headBytes[1:])
	if err != nil {
		return nil, err
	}

	if length = head.Len; length > PacketV2Limit {
		return nil, errors.New("the size of packet is larger than the limit")
	}

	buff := make([]byte, length)
	copy(buff[0:PacketV2HeadLen], headBytes)

	// read body ( buff = lengthBytes + body )
	if _, err := io.ReadFull(conn, buff[PacketV2HeadLen:]); err != nil {
		return nil, err
	}

	return NewHeadV2Packet(buff), nil
}
