// TcpConnector project main.go
package main

import (
	"fmt"
	"net"
	"time"
)

var (
	rAddr *net.TCPAddr
	err   error
)

func main() {
	fmt.Println("Connector Working...")

	rAddr, err = net.ResolveTCPAddr("tcp4", "192.168.100.141:7576")

	run(5000)

	time.Sleep(200 * time.Second)
}

func connect(n uint) {
	conn, err := net.DialTCP("tcp4", nil, rAddr)

	if err != nil {
		fmt.Printf("Connecting failed: %v\n", err)
	} else {
		fmt.Printf("connecting success: %d\n", n)
	}

	time.Sleep(200 * time.Second)

	fmt.Println("connection closeing...")

	conn.Close()
}

func run(nums uint) {
	var i uint

	for i = 0; i < nums; i++ {
		go connect(i)
	}
}
