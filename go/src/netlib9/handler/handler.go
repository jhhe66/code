package handler

import (
	"cast"
	"errors"
	"fmt"
	"net"
	"netlib9/protocol"
	"sync"
)

const (
	R_BUFFER_SIZE = 8 * 1024
	W_BUFFER_SIZE = 8 * 1024
)

const (
	REQ_HEAD = iota //接收中
	REQ_BODY        //接收完成
	REQ_OK          //完成
	REQ_ERR         //错误
	SENDING         //发送中
)

var G_CONN_COUNT uint64 = 0 //连接计数

type IPC chan *protocol.Packet

type TcpConnectionMap map[uint64]*TcpConnection // UserID -> TcpConnection

//表示连接对象
type TcpConnection struct {
	Conn   net.Conn            //网络连接
	R_C    IPC                 //与routine交互的channel
	UserID uint32              //用户ID
	Id     uint64              //连接标识
	R      [R_BUFFER_SIZE]byte //读缓冲
	R_H    uint32              //读头
	R_T    uint32              //读尾
	W      [W_BUFFER_SIZE]byte //写缓冲
	W_H    uint32              //写头
	W_T    uint32              //写尾
	Status uint8               //包状态
	Locker sync.Mutex          // mutex
}

func NewConntion(c net.Conn, ipc IPC) *TcpConnection {
	nc := new(TcpConnection)

	nc.Conn = c
	nc.R_C = ipc
	nc.Status = REQ_HEAD
	nc.Id = G_CONN_COUNT + 1

	return nc
}

func (self *TcpConnection) Read() (int32, error) {
	fmt.Printf("R_T: %d\n", self.R_T)

	r, e := self.Conn.Read(self.R[self.R_T:])

	if e == nil {
		self.R_T = self.R_T + uint32(r)
		fmt.Printf("R_T: %d\n", self.R_T)
	}

	return int32(r), e
}

func (self *TcpConnection) Write(b []byte, size uint32) (int32, error) {
	copy(self.W[self.W_T:], b[:size])

	self.W_T = self.W_T + size

	var (
		r int
		e error
	)

	for {
		r, e = self.Conn.Write(self.W[self.W_H:self.W_T])

		if e == nil {
			self.W_H = self.W_H + uint32(r)
		}

		//当前的缓冲区已经发完了
		if self.W_Clear() {
			break
		}
	}

	return int32(r), e
}

func (self *TcpConnection) ParseHead() (bool, error) {
	fmt.Printf("ParseHead-R_T: %d\n", self.R_T)
	fmt.Printf("ParseHead-R_H: %d\n", self.R_H)

	if (self.R_T - self.R_H) < 9 {
		return false, nil
	} else {
		//判断是否正确包头
		if string(self.R[self.R_H+2:self.R_H+3]) == "B" && string(self.R[self.R_H+3:self.R_H+4]) == "Y" {
			return true, nil
		} else {
			return false, errors.New("REQ_HEAD ERROR")
		}
	}
}

func (self *TcpConnection) ParseBody() (bool, error) {
	if plen, _ := cast.Byte2int16(self.R[self.R_H : self.R_H+2]); plen <= int16(self.R_T-self.R_H) {
		return true, nil
	} else {
		return false, errors.New("REQ_BODY ERROR")
	}
}

func (self *TcpConnection) Parse() (bool, error) {

	fmt.Printf("CONN-STATUS: %d\n", self.Status)

	if self.Status == REQ_HEAD {
		if l, e := self.ParseHead(); l && e == nil {
			self.Status = REQ_BODY
		} else {
			fmt.Printf("CONN-STATUS-4: %d\n", self.Status)
			return false, e
		}
	}

	if self.Status == REQ_BODY {
		if l, e := self.ParseBody(); l && e == nil {
			self.Status = REQ_HEAD //REQ_OK
			fmt.Printf("CONN-STATUS-2: %d\n", self.Status)
			return true, nil
		} else {
			fmt.Printf("CONN-STATUS-3: %d\n", self.Status)
			return false, e
		}
	}

	fmt.Printf("CONN-STATUS-1: %d\n", self.Status)

	return false, nil
}

func (self *TcpConnection) Close() {
	self.Conn.Close()
	if self.R_C != nil {
		close(self.R_C)
	}

}

func (self *TcpConnection) R_Clear() bool {
	if self.R_H < self.R_T {
		return false
	}

	for i := uint32(0); i < self.R_T; i++ {
		self.R[i] = 0
	}

	self.R_H, self.R_T = 0, 0

	return true
}

func (self *TcpConnection) W_Clear() bool {
	if self.W_H < self.W_T {
		return false
	}

	for i := uint32(0); i < self.W_T; i++ {
		self.W[i] = 0
	}

	self.W_H, self.W_T = 0, 0

	return true
}
