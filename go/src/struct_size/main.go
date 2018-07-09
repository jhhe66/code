package main

import (
	"encoding/binary"
	"fmt"
	"reflect"
	"unsafe"
)

type A struct {
	a, b uint32
}

type B struct {
	a [2]uint32
	b uint16
}

type PacketHeader struct {
	Len         uint16   /*整包长度，包括包头+包体，但不包含自身的两个字节*/
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
	fmt.Printf("A: %v\n", unsafe.Sizeof(A{}))
	fmt.Printf("B: %v\n", unsafe.Sizeof(B{}))
	fmt.Printf("B: %v\n", reflect.TypeOf(B{}).Size())
	fmt.Printf("B: %v\n", binary.Size(B{}))
	fmt.Printf("C: %v\n", unsafe.Sizeof(PacketHeader{}))
	fmt.Printf("C: %v\n", binary.Size(PacketHeader{}))
	return
}
