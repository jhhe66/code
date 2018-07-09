package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
)

type PacketHeader struct {
	Pack_len    uint16   /*整包长度，包括包头+包体，但不包含自身的两个字节*/
	Start       [2]uint8 /*消息开始标记*/
	Main_ver    uint8    /*协议主版本号*/
	Sub_ver     uint8    /*协议子版本号*/
	Main_cmd    uint16   /*主命令字*/
	Code        uint8    /*校验code*/
	Sub_cmd     uint16   /*子命令字*/
	Seq         uint16   /*序列号 首次就填0*/
	Source_type uint8    /*消息来源 */
}

func main() {
	fmt.Println("Starting the server ...")
	// create listener:
	listener, err := net.Listen("tcp", "localhost:4444")
	if err != nil {
		fmt.Println("Error listening", err.Error())
		return // terminate program
	}
	// listen and accept connections from clients:
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting", err.Error())
			return // terminate program
		}
		go doServerStuff(conn)
	}
}

func cat() {
	fmt.Print()

}

func doServerStuff(conn net.Conn) {
	for {
		buf := make([]byte, 512)
		//buf := new(bytes.Buffer)

		n, err := conn.Read(buf)

		if err != nil {
			fmt.Println("Error reading", err.Error())
			return // terminate program
		}

		fmt.Printf("recv: %d\n", n)


		if n >= 14 {
			GetHeader(buf, n)
		}

		fmt.Printf("Received data: %v", string(buf))
	}
}

func Log(format string, v ...interface{}) {
	fmt.Printf(format, v...)
}

func GetHeader(buff []byte, rLen int) {
	Log("----------GetHeader begin -----------------\n")

	bb := bytes.NewBuffer(buff)

	var header PacketHeader

	err := binary.Read(bb, binary.BigEndian, &header)

	// for i := 0; i < rLen; i++ {
	// 	cmd[i] = buff[i]
	// }
	if err == nil {
		fmt.Printf("CMD: %v\n", header)
	} else {
		fmt.Printf("err: %v\n", err)
	}

	Log("----------GetHeader end -----------------\n")
}
