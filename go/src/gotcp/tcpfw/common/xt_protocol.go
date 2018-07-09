package common

import (
	"encoding/binary"
	"errors"
	"io"
	"net"

	"github.com/gansidui/gotcp"
)

const (
	XTHeadLen       = 20
	PacketXTHeadLen = 20
	EmptyPacktXTLen = 20
	PacketXTLimit   = 128 * 1024
)

// Error type
var (
	XTErrShortLen = errors.New("use byte[] lne is not enough")
	XTErrLenErr   = errors.New("Len error")
)

type XTHead struct {
	Flag     uint8  // 0xF0客户端请求，0xF1 服务器应答, 0xF2  服务器主动发包, 0xF3  客户端应答, 0xF4 服务器之间的包
	Version  uint8  // 版本号  VER_MMEDIA = 4
	CryKey   uint8  // 加密类型  E_NONE_KEY = 0, E_SESSION_KEY = 1, E_RAND_KEY = 2, E_SERV_KEY= 3
	TermType uint8  // 终端类型
	Cmd      uint16 // 命令字
	Seq      uint16 // 序列号
	From     uint32 // uint32_t uiFrom
	To       uint32 // 目的UID TO_SERVER = 0
	Len      uint32 // PayLoad的总长度不包含包头
}

func NewXTHead(buf []byte) (head *XTHead, err error) {
	if len(buf) < XTHeadLen {
		return nil, XTErrLenErr
	}
	head = new(XTHead)
	head.Flag = buf[0]
	head.Version = buf[1]
	head.CryKey = buf[2]
	head.TermType = buf[3]
	head.Cmd = binary.LittleEndian.Uint16(buf[4:6])
	head.Seq = binary.LittleEndian.Uint16(buf[6:8])
	head.From = binary.LittleEndian.Uint32(buf[8:12])
	head.To = binary.LittleEndian.Uint32(buf[12:16])
	head.Len = binary.LittleEndian.Uint32(buf[16:20])
	return head, nil
}

func SerialXTHeadToSlice(head *XTHead, buf []byte) (err error) {
	if len(buf) < XTHeadLen {
		return XTErrShortLen
	}
	buf[0] = head.Flag
	buf[1] = head.Version
	buf[2] = head.CryKey
	buf[3] = head.TermType
	binary.LittleEndian.PutUint16(buf[4:6], head.Cmd)
	binary.LittleEndian.PutUint16(buf[6:8], head.Seq)
	binary.LittleEndian.PutUint32(buf[8:12], head.From)
	binary.LittleEndian.PutUint32(buf[12:16], head.To)
	binary.LittleEndian.PutUint32(buf[16:20], head.Len)
	return nil
}

// XTHead 格式如下
// XTHead + payload

type XTHeadPacket struct {
	buff []byte
}

func (this *XTHeadPacket) Serialize() []byte {
	return this.buff
}

// index:0 Falg field
func (this *XTHeadPacket) GetFlag() uint8 {
	return this.buff[0]
}

// index:1 Version field
func (this *XTHeadPacket) GetVersion() uint8 {
	return this.buff[1]
}

// index:2 Crypto field
func (this *XTHeadPacket) GetKey() uint8 {
	return this.buff[2]
}

// index:3 Termianl type field
func (this *XTHeadPacket) GetTerminalType() uint8 {
	return this.buff[3]
}

// index:4 Command field
func (this *XTHeadPacket) GetCommand() uint16 {
	return binary.LittleEndian.Uint16(this.buff[4:6])
}

// index:5 Sequence field
func (this *XTHeadPacket) GetSeq() uint16 {
	return binary.LittleEndian.Uint16(this.buff[6:8])
}

// index:6 From Uid field
func (this *XTHeadPacket) GetFromUid() uint32 {
	return binary.LittleEndian.Uint32(this.buff[8:12])
}

// index:7 To Uid field
func (this *XTHeadPacket) GetToUid() uint32 {
	return binary.LittleEndian.Uint32(this.buff[12:16])
}

// index:8 Packet Length field length(SOH+HeadV3+PayLoad+EOT)
func (this *XTHeadPacket) GetLength() uint32 {
	return binary.LittleEndian.Uint32(this.buff[16:20])
}

// index:9 Body filed
func (this *XTHeadPacket) GetBody() []byte {
	return this.buff[XTHeadLen:]
}

func (this *XTHeadPacket) GetHead() (head *XTHead, err error) {
	head, err = NewXTHead(this.buff[:XTHeadLen])
	return
}

func (this *XTHeadPacket) CheckXTPacketValid() (bool, error) {
	if this.GetLength() != uint32(len(this.Serialize())-XTHeadLen) {
		return false, XTErrLenErr
	}

	return true, nil
}

func NewXTHeadPacket(buff []byte) *XTHeadPacket {
	p := &XTHeadPacket{}
	p.buff = buff
	return p
}

type XTHeadProtocol struct {
}

func (this *XTHeadProtocol) ReadPacket(conn *net.TCPConn) (gotcp.Packet, error) {

	var (
		headBytes []byte = make([]byte, XTHeadLen)
		length    uint32
	)

	// read length
	if _, err := io.ReadFull(conn, headBytes); err != nil {
		return nil, err
	}

	head, err := NewXTHead(headBytes[:])
	if err != nil {
		return nil, err
	}

	if length = head.Len; length > PacketXTLimit {
		return nil, errors.New("the size of packet is larger than the limit")
	}

	// length 不包含报文头部的长度
	buff := make([]byte, length+XTHeadLen)
	copy(buff[0:XTHeadLen], headBytes)

	// read body ( buff = lengthBytes + body )
	if _, err := io.ReadFull(conn, buff[XTHeadLen:]); err != nil {
		return nil, err
	}

	return NewXTHeadPacket(buff), nil
}
