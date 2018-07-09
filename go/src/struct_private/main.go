package main

import (
	"struct_private/def"
	"fmt"
	"unsafe"
)

func main() {
	hh := new(def.Human)
	page := (*int)(unsafe.Pointer(hh))
	*page = 38
	psex := (*int)(unsafe.Pointer(uintptr(unsafe.Pointer(hh)) + unsafe.Sizeof(int(0))))
	*psex = 2
	pname := (*string)(unsafe.Pointer(uintptr(unsafe.Pointer(hh)) + 2 * unsafe.Sizeof(int(0))))
	*pname = "chenbo"
	
	fmt.Printf("age: %v\n", hh.Age())
	fmt.Printf("sex: %v\n", hh.Sex())
	fmt.Printf("name: %v\n", hh.Name())
}
