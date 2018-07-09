package ice_tag

/*
#cgo CFLAGS: -I../ice_tagalias
#cgo LDFLAGS: -L../ice_tagalias -lice_tagalias -lIce -lIceUtil

#include "ice_tagalias.h"
#include <stdlib.h>
#include <stdio.h>
*/
import "C"
import "unsafe"

type GoIce struct {
	IceHandle C.ice_warper_t
}

func New(conf string, proxy string) *GoIce {
	var ret *GoIce

	conf_path := C.CString(conf)
	proxy_name := C.CString(proxy)
	ret = new(GoIce)
	if ret != nil {
		ret.IceHandle = C.ice_tagalias_init(conf_path, proxy_name)
	}
	
	defer C.free(unsafe.Pointer(conf_path))
	defer C.free(unsafe.Pointer(proxy_name))
	
	return ret
}

func (ice *GoIce) Request(req []byte, res []byte) (int, uint32) {
	var request 	[512]byte
	var response 	[512]byte
	var rsz 		*C.uint

	copy(request[:], req)
	rsz = new(C.uint)
	*rsz = C.uint(len(response))

	ret := C.ice_tagalias_request(ice.IceHandle, (*C.char)(unsafe.Pointer(&request[0])), (*C.char)(unsafe.Pointer(&response[0])), rsz)
	copy(res, response[:*rsz])
	

	return int(ret), uint32(*rsz)
}

func (ice *GoIce) ValidateTags(req []byte, res []byte) (int, uint32) {
	var request 	[512]byte
	var response 	[512]byte
	var rsz 		*C.uint

	copy(request[:], req)
	rsz = new(C.uint)
	*rsz = C.uint(len(response))

	ret := C.ice_tagalias_validateTags(ice.IceHandle, (*C.char)(unsafe.Pointer(&request[0])), (*C.char)(unsafe.Pointer(&response[0])), rsz)
	copy(res, response[:*rsz])
	

	return int(ret), uint32(*rsz)	
}

func (ice *GoIce) IceFree() {
	C.ice_tagalias_free(ice.IceHandle)
}
