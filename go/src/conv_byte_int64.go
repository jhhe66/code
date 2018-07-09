package main

import (
        "fmt"
        "encoding/binary"
        "io"
        "bytes"
)

type Byte64 [8]byte

func (bits *Byte64) Read(p []byte) (r int, e error) {
        if len(p) > 8 {
                e = io.EOF
        }
        for r = 0; r < 8 && r < len(p); r++ {
                p[r] = bits[r]
        }
        return
}

type Header struct {
        Len     uint16
        Mark    [2]uint8
        Cmd     uint16
        Source  uint16
}

// func (header *Header) Read(p []byte) (r int, e error) {
//         if len(p) > 8 {
//                 e = io.EOF
//         }

//         for r = 0; r < 8 && r < len(p); r++ {
//                 p[r] = header[r]
//         }

//         return 
// }

func main() {
        var n int64
        bits := new(Byte64)
        bits[0] = 1
        bits[1] = 2
        binary.Read(bits, binary.LittleEndian, &n)
        fmt.Printf("%v \n", n)
        //binary.Read(bits, binary.BigEndian, &n)

        //buffer := new([8]byte)

        buffer_new := new(bytes.Buffer)//bytes.NewBuffer(buffer)

        header := new(Header)

        header.Len = 1
        header.Mark[0] = 2
        header.Mark[1] = 3
        header.Cmd = 4
        header.Source = 5

        binary.Write(buffer_new, binary.LittleEndian, header)

        // header2 := new(Header)

        // binary.Read(buffer_new, binary.LittleEndian, &header2)

        fmt.Printf("%x \n", buffer_new) //513

        var header2 Header

        binary.Read(buffer_new, binary.LittleEndian, &header2)

        fmt.Printf("%v \n", header2)
}
