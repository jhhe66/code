package libcomm

/*
#cgo CFLAGS: -I./cinclude
#cgo LDFLAGS: -L./ -lneocomm
#include <stdlib.h>
#include "Attr_API.h"
*/
import "C"
import "unsafe"

func AttrAdd(name string, value int32) {
	cs := C.CString(name)
	defer C.free(unsafe.Pointer(cs))
	C.AttrAdd(cs, C.int(value))
}

func AttrAddAvg(name string, value int32) {
	cs := C.CString(name)
	defer C.free(unsafe.Pointer(cs))
	C.AttrAddAvg(cs, C.int(value))
}

func AttrSet(name string, value int32) {
	cs := C.CString(name)
	defer C.free(unsafe.Pointer(cs))
	C.AttrSet(cs, C.int(value))
}
