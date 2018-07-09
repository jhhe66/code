// conn_deadline project main.go
package main

import (
	"fmt"
	"net"
	"time"
)

func CheckAccept(e error) {

}

func CheckListen(e error) {

}

func HandleConn(c net.Conn) {
	buff := [100]byte{}

	c.SetDeadline(time.Now().Add(10 * time.Second))

	r, error := c.Read(buff[:])

	fmt.Printf("read return: %d err: %v\n", r, error)
}

func main() {
	listen, err := net.Listen("tcp", ":9999")
	CheckListen(err)

	for {
		conn, err := listen.Accept()
		CheckAccept(err)

		go HandleConn(conn)
	}
}
