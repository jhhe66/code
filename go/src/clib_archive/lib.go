package main

import "C"
import "fmt"

//export GoCall
func GoCall(buffer *C.char) {
	fmt.Println(C.GoString(buffer))
}

func main() {
}
