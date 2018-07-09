package main

import (
	"unsafe"
	"reflect"
	"fmt"
)

type Foo struct {
	Age 	int
	Name 	string
	Sex		int8
}

func main() {
	foo := &Foo{30, "chenbo", 1}
	sb := foo.ToBytes() 
	foo.Age = 40
	fmt.Printf("sb: %v\n", sb)
	foo2 := &Foo{}
	foo2.FromBytes(sb)
	fmt.Printf("foo2: %v\n", foo2)
}

func (f *Foo)ToBytes() []byte {
	FooSz := int(unsafe.Sizeof(Foo{}))
	var x reflect.SliceHeader = reflect.SliceHeader{
										Len:FooSz,
										Cap:FooSz,
										Data:uintptr(unsafe.Pointer(f))}
	return *(*[]byte)(unsafe.Pointer(&x))
}

func (f *Foo)FromBytes(b []byte) bool {
	*f = *(*Foo)(unsafe.Pointer(&b[0]))
	return true;
}

