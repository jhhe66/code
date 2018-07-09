package common

// "encoding/binary"

func MarshalUint8(value uint8, packet *[]byte) {
	*packet = append(*packet, value)
}

// Linux 小端存储 低位存低字节 高位存高字节
func MarshalUint16(value uint16, packet *[]byte) {
	*packet = append(*packet, uint8(value&0xFF))
	*packet = append(*packet, uint8((value>>8)&0xFF))
}

func MarshalUint32(value uint32, packet *[]byte) {
	*packet = append(*packet, uint8(value&0xFF))
	*packet = append(*packet, uint8((value>>8)&0xFF))
	*packet = append(*packet, uint8((value>>16)&0xFF))
	*packet = append(*packet, uint8((value>>24)&0xFF))
}

func MarshalUint64(value uint64, packet *[]byte) {
	*packet = append(*packet, uint8(value&0xFF))
	*packet = append(*packet, uint8((value>>8)&0xFF))
	*packet = append(*packet, uint8((value>>16)&0xFF))
	*packet = append(*packet, uint8((value>>24)&0xFF))
	*packet = append(*packet, uint8((value>>32)&0xFF))
	*packet = append(*packet, uint8((value>>40)&0xFF))
	*packet = append(*packet, uint8((value>>48)&0xFF))
	*packet = append(*packet, uint8((value>>56)&0xFF))
}

func MarshalSlice(value []byte, packet *[]byte) {
	length := len(value)
	// 字符串末尾添加'\0'
	length += 1
	value = append(value, 0)
	// 2字节长度+实际内容
	*packet = append(*packet, uint8(length&0xFF))
	*packet = append(*packet, uint8((length>>8)&0xFF))
	// copy 内容 跳过两个字节长度位置
	*packet = append(*packet, value...)
}

func UnMarshalUint8(packet *[]byte) (value uint8) {
	if len(*packet) < 1 {
		return 0
	}

	value = uint8((*packet)[0])
	*packet = (*packet)[1:]
	return
}

func UnMarshalUint16(packet *[]byte) (value uint16) {
	if len(*packet) < 2 {
		return 0
	}

	value = uint16((*packet)[0]) | uint16((*packet)[1])<<8
	*packet = (*packet)[2:]
	return
}

func UnMarshalUint32(packet *[]byte) (value uint32) {
	if len(*packet) < 4 {
		return 0
	}

	value = uint32((*packet)[0]) | uint32((*packet)[1])<<8 | uint32((*packet)[2])<<16 | uint32((*packet)[3])<<24
	*packet = (*packet)[4:]
	return
}

func UnMarshalUint64(packet *[]byte) (value uint64) {
	if len(*packet) < 8 {
		return 0
	}
	value = uint64((*packet)[0]) | uint64((*packet)[1])<<8 | uint64((*packet)[2])<<16 | uint64((*packet)[3])<<24 | uint64((*packet)[4])<<32 | uint64((*packet)[5])<<40 | uint64((*packet)[6])<<48 | uint64((*packet)[7])<<56
	*packet = (*packet)[8:]

	return
}

func UnMarshalSlice(packet *[]byte) (value []byte) {
	if len(*packet) < 2 {
		return nil
	}

	length := uint16((*packet)[0]) | uint16((*packet)[1])<<8
	if length == 0 { // empty slice
		*packet = (*packet)[2:]
		return nil
	}

	// 删除value中的'\0'
	payLoadLen := length - 1
	value = make([]byte, payLoadLen)
	copy(value, (*packet)[2:]) // 跳过长度 2字节
	// 使用2+length 自动跳过'\0'
	*packet = (*packet)[length+2:]
	return
}
