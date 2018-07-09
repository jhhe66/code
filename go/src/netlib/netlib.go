// net project net.go
package netlib

import (
	"fmt"
	"net"
)

func Version() {
	fmt.Println("1.0\n")
}

type MessageEventHandler func(uint16) int32     // 消息命令处理函数
type ConnCloseEventHandler func(net.Conn) int32 //连接关闭处理函数

type TcpServer struct {
	Host        string
	HandleRead  MessageEventHandler
	HandleClose ConnCloseEventHandler
}

func NewTcpServer() *TcpServer {

}

func (self TcpServer) Run() {

}

func (self TcpServer) Init() {

}
