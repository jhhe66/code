// cast project cast.go
package cast

import (
	"bytes"
	"encoding/binary"
)

func Byte2uint64(b []byte) (uint64, bool) {
	buff := bytes.NewBuffer(b[:8])

	var result uint64

	error := binary.Read(buff, binary.BigEndian, &result)

	if error == nil {
		return result, true
	} else {
		return 0, false
	}

	return 0, false
}

func Byte2int64(b []byte) (int64, bool) {

	buff := bytes.NewBuffer(b[:8])

	var result int64

	error := binary.Read(buff, binary.BigEndian, &result)

	if error == nil {
		return result, true
	} else {
		return 0, false
	}

	return 0, false
}

func Byte2int32(b []byte) (int32, bool) {
	buff := bytes.NewBuffer(b[:4])

	var result int32

	error := binary.Read(buff, binary.BigEndian, &result)

	if error == nil {
		return result, true
	} else {
		return 0, false
	}

	return 0, false
}

func Byte2int16(b []byte) (int16, bool) {
	buff := bytes.NewBuffer(b[:2])

	var result int16

	error := binary.Read(buff, binary.BigEndian, &result)

	if error == nil {
		return result, true
	} else {
		return 0, false
	}

	return 0, false
}

func Byte2int8(b []byte) (int8, bool) {
	buff := bytes.NewBuffer(b)

	var result int8

	error := binary.Read(buff, binary.BigEndian, &result)

	if error == nil {
		return result, true
	} else {
		return 0, false
	}

	return 0, false
}

func Int82byte(i int8, b []byte) bool {
	buff := bytes.NewBuffer(b[:1])

	error := binary.Write(buff, binary.BigEndian, &i)

	if error == nil {
		buff.Next(1) // 移动四个byte

		bbb := buff.Next(1)

		b[0] = bbb[0]

		return true
	} else {
		return false
	}

	return false
}

func Int162byte(i int16, b []byte) bool {
	buff := bytes.NewBuffer(b[:2])

	error := binary.Write(buff, binary.BigEndian, &i)

	if error == nil {
		buff.Next(2) // 移动2个byte

		bbb := buff.Next(2)

		b[0] = bbb[0]
		b[1] = bbb[1]

		return true
	} else {
		return false
	}

	return false
}

func Int322byte(i int32, b []byte) bool {
	buff := bytes.NewBuffer(b[:4])

	error := binary.Write(buff, binary.BigEndian, &i)

	if error == nil {
		bbb := buff.Next(4)

		b[0] = bbb[0]
		b[1] = bbb[1]
		b[2] = bbb[2]
		b[3] = bbb[3]

		return true
	} else {
		return false
	}

	return false
}

func Int642byte(i int64, b []byte) bool {
	buff := bytes.NewBuffer(b[:8])

	error := binary.Write(buff, binary.BigEndian, &i)

	if error == nil {
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
}
