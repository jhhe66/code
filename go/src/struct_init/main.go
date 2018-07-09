// struct_init project main.go
package main

import (
//"fmt"
)

type People struct {
	Name [8]byte
	Age  uint16
	Sex  uint8
}

var chenbo People = People{[8]byte{'c'}, 16, 1}

func main() {

}
