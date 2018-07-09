// netlib_test project main.go
package main

import (
	"net"
	//"netlib"
	"netlib/handler"
	"netlib/protocol"
	//"bytes"
	//"cast"
	//"encoding/binary"
	"fmt"
	"os"
)

func CheckListen(err error) {
	if err != nil {
		fmt.Printf("Listen Error %v\n", err)
		os.Exit(1)
	}
}

func CheckAccept(err error) {
	if err != nil {
		fmt.Printf("Accept Error: %v\n", err)
	}
}

func HandleConnection(c *handler.TcpConnection) {
	fmt.Printf("New Connection %d\n", c.Id)
	defer c.Close()

	for {
		rLen, e := c.Read()

		fmt.Printf("rLen: %d\n", rLen)

		if rLen > 0 {
			if l, err := c.Parse(); l && err == nil {
				fmt.Printf("Read1: %v\n", c.R[c.R_H:c.R_T])
				pp := protocol.NewPacketByBuffer(c.R[c.R_H:c.R_T], c.R_T-c.R_H)

				//c.R_H = c.R_T

				fmt.Printf("R_H_1: %d\n", c.R_H)

				fmt.Printf("PacketLen: %d\n", pp.GetLen())
				// 调整 c.R_T c.R_H
				c.R_H = c.R_H + uint32(pp.GetLen())

				//check R clear
				c.R_Clear()

				fmt.Printf("Decode: %d\n", pp.Decode())

				fmt.Printf("Read2: %v\n", c.R[:c.R_T])

				fmt.Printf("INT: %d\n", pp.UnpackInt())
				//fmt.Printf("String: %s\n", pp.UnpackString())
				//fmt.Printf("SHORT: %d\n", pp.UnpackShort())
				//fmt.Printf("Byte: %d\n", uint8(pp.UnpackByte()))
				//fmt.Printf("INT64: %d\n", pp.UnpackInt64())

				fmt.Printf("Error: %v\n", e)

				fmt.Printf("Write Message \n")
				wp := protocol.NewPacketByCmd(0x1001, true)

				wp.PackInt(1000)
				s := "1234"
				wp.PackString(&s, 4)
				wp.PackInt(9999)

				wp.Pack()

				wlen, err := c.Write(wp.Buffer[:wp.Length], uint32(wp.Length))

				fmt.Printf("WRITE: %v\n", wp.Buffer[:wp.Length])

				fmt.Printf("wlen: %d\n", wlen)

				if wlen > 0 {
					fmt.Printf("Write: %d\n", wlen)
				} else if err != nil {
					fmt.Printf("W_ERROR: %v\n", err)
					break
				}

			} else if err != nil {
				fmt.Printf("Error: %v\n", err)
				fmt.Printf("connecton colse by server\n")
				break
			} else if !l {
				fmt.Printf("Read3: %v\n", c.R[c.R_H:c.R_T])
				fmt.Printf("R_H: %d\n", c.R_H)
				fmt.Printf("R_T: %d\n", c.R_T)
				fmt.Printf("Error: %v\n", err)
			}
		} else if rLen == 0 {
			//错误处理
			fmt.Printf("connection rest by peer\n")
			fmt.Printf("Error: %v\n", e)
			fmt.Printf("connecton colse by server\n")
			break
		} else if rLen == -1 {
			fmt.Printf("connection rest by peer 1\n")
			fmt.Printf("Error: %v\n", e)
			break
		}
	}
}

var (
	ConnMap handler.TcpConnectionMap = handler.TcpConnectionMap{}
)

func main() {
	listen, err := net.Listen("tcp", ":9999")
	CheckListen(err)

	for {
		conn, err := listen.Accept()
		CheckAccept(err)

		pCon := handler.NewConntion(conn, nil)

		ConnMap[handler.G_CONN_COUNT] = pCon

		go HandleConnection(ConnMap[handler.G_CONN_COUNT])
	}

}
