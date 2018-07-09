// type_cast project main.go
package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

func byte2uint64(b []byte) (uint64, bool) {
	fmt.Printf("b: %v\n", b)

	buff := bytes.NewBuffer(b[:8])

	var result uint64

	error := binary.Read(buff, binary.BigEndian, &result)

	if error == nil {
		return result, true
	} else {
		fmt.Printf("cast error: %v\n", error)
		return 0, false
	}

	return 0, false
}

func byte2int64(b []byte) (int64, bool) {
	fmt.Printf("b: %v\n", b)

	buff := bytes.NewBuffer(b[:8])

	var result int64

	error := binary.Read(buff, binary.BigEndian, &result)

	if error == nil {
		return result, true
	} else {
		fmt.Printf("cast error: %v\n", error)
		return 0, false
	}

	return 0, false
}

func byte2int32(b []byte) (int32, bool) {
	fmt.Printf("b: %v\n", b)

	buff := bytes.NewBuffer(b[:4])

	var result int32

	error := binary.Read(buff, binary.BigEndian, &result)

	if error == nil {
		return result, true
	} else {
		fmt.Printf("cast error: %v\n", error)
		return 0, false
	}

	return 0, false
}

func byte2int16(b []byte) (int16, bool) {
	fmt.Printf("b: %v\n", b)

	buff := bytes.NewBuffer(b[:2])

	var result int16

	error := binary.Read(buff, binary.BigEndian, &result)

	if error == nil {
		return result, true
	} else {
		fmt.Printf("cast error: %v\n", error)
		return 0, false
	}

	return 0, false
}

func byte2int8(b []byte) (int8, bool) {
	fmt.Printf("b: %v\n", b)

	buff := bytes.NewBuffer(b)

	var result int8

	error := binary.Read(buff, binary.BigEndian, &result)

	if error == nil {
		return result, true
	} else {
		fmt.Printf("cast error: %v\n", error)
		return 0, false
	}

	return 0, false
}

func int82byte(i int8, b []byte) bool {
	buff := bytes.NewBuffer(b[:1])

	error := binary.Write(buff, binary.BigEndian, &i)

	if error == nil {
		fmt.Printf("len: %v\n", len(buff.Bytes()))

		fmt.Printf("b: %x\n", buff.Next(1)) // 移动四个byte

		bbb := buff.Next(1)

		b[0] = bbb[0]
		//b[2] = bbb[2]
		//b[3] = bbb[3]

		//fmt.Printf("b: %x\n", buff.Next(4))
		return true
	} else {
		return false
	}

	return false
}

func int162byte(i int16, b []byte) bool {
	buff := bytes.NewBuffer(b[:2])

	error := binary.Write(buff, binary.BigEndian, &i)

	if error == nil {
		fmt.Printf("len: %v\n", len(buff.Bytes()))

		fmt.Printf("b: %x\n", buff.Next(2)) // 移动四个byte

		bbb := buff.Next(2)

		b[0] = bbb[0]
		b[1] = bbb[1]
		//b[2] = bbb[2]
		//b[3] = bbb[3]

		//fmt.Printf("b: %x\n", buff.Next(4))
		return true
	} else {
		return false
	}

	return false
}

func int322byte(i int32, b []byte) bool {
	buff := bytes.NewBuffer(b[:4])

	error := binary.Write(buff, binary.BigEndian, &i)

	if error == nil {
		fmt.Printf("len: %v\n", len(buff.Bytes()))

		fmt.Printf("b: %x\n", buff.Next(4)) // 移动四个byte

		bbb := buff.Next(4)

		b[0] = bbb[0]
		b[1] = bbb[1]
		b[2] = bbb[2]
		b[3] = bbb[3]

		//fmt.Printf("b: %x\n", buff.Next(4))
		return true
	} else {
		return false
	}

	return false
}

func int642byte(i int64, b []byte) bool {
	buff := bytes.NewBuffer(b[:8])

	error := binary.Write(buff, binary.BigEndian, &i)

	if error == nil {
		fmt.Printf("len: %v\n", len(buff.Bytes()))

		fmt.Printf("b: %x\n", buff.Next(8)) // 移动四个byte

		bbb := buff.Next(8)

		for k, _ := range bbb {
			b[k] = bbb[k]
		}

		return true
	} else {
		return false
	}

	return false
}

func string2byte(s *string) {
	buffer := [100]byte{}

	sb := buffer[10:20]

	temp := []byte((*s))

	copy(sb, temp)

	sb[5] = 0

	fmt.Printf("sb: %s\n", sb)
	fmt.Printf("sb: %v\n", buffer)
}

func main() {
	b1 := [1]byte{0x33}
	b11 := [1]byte{}

	b2 := [2]byte{0x01, 0x02}

	b22 := [2]byte{}

	b4 := [...]byte{0x01, 0x02, 0x03, 0x04}

	b44 := [4]byte{}

	b88 := [8]byte{}

	b8 := [8]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}

	i, ok := byte2int8(b1[:])

	if ok {
		fmt.Printf("i: 0x%x\n", i)
	} else {
		fmt.Printf("Error\n")
	}

	i2, ok := byte2int16(b2[:])

	if ok {
		fmt.Printf("i2: 0x%x\n", i2)
	} else {
		fmt.Printf("Error\n")
	}

	i4, ok := byte2int32(b4[:])

	if ok {
		fmt.Printf("i2: 0x%d\n", i4)
	} else {
		fmt.Printf("Error\n")
	}

	if ok := int82byte(99, b11[:]); ok {
		fmt.Printf("b11: %v\n", b11)
		i, _ := byte2int8(b22[:])
		fmt.Printf("b11: %d\n", i)
	} else {
		fmt.Printf("Error\n")
	}

	if ok := int162byte(1000, b22[:]); ok {
		fmt.Printf("b22: %v\n", b22)
		i, _ := byte2int16(b22[:])
		fmt.Printf("b22: %d\n", i)
	} else {
		fmt.Printf("Error\n")
	}

	if ok := int322byte(1000, b44[:]); ok {
		fmt.Printf("b44: %v\n", b44)
	} else {
		fmt.Print("Error\n")
	}

	if ok := int642byte(1000, b88[:]); ok {
		fmt.Printf("b88: %v\n", b88)
	} else {
		fmt.Print("Error\n")
	}

	i8, ok := byte2int64(b8[:])

	if ok {
		fmt.Printf("i8: %d\n", i8)
	} else {
		fmt.Print("Error\n")
	}

	var str string = "hello"

	string2byte(&str)

	//var i6 int32 = 99

	//fmt.Printf("i: %v\n", []byte(i6))

}
