package main

import (
   "bytes"
   "encoding/binary"
   "fmt"
)
type Head struct {
   Type int32
   Len  int32
}
type Packet struct {
        Head
        Content string
}

func main() {
   h := &Packet{Head{200,300},"string"}
   buf := new(bytes.Buffer)
   if err := binary.Write(buf, binary.BigEndian, h.Head); nil != err {
       fmt.Println(err)
       return
   }
   
   buf.WriteString(h.Content)


   fmt.Printf("% x", buf.Bytes())
}