package libcrypto

/*
#cgo CFLAGS: -I./
#cgo LDFLAGS: -L./ -lteacrypto
#include <stdlib.h>
#include "TEACrypto.h"
*/
import "C"
import "unsafe"

const (
	CCryptoBufSize = 64 * 1024
)

func TEAEncrypt(plainText string, cryptoKey string) (cryptoText string) {
	pt := C.CString(plainText)
	defer C.free(unsafe.Pointer(pt))
	ptLen := len(plainText)
	ck := C.CString(cryptoKey)
	defer C.free(unsafe.Pointer(ck))

	var cryptoLen C.uint
	cryptoBuf := unsafe.Pointer(C.malloc(CCryptoBufSize))
	defer C.free(cryptoBuf)

	C.xTEAEncryptWithKey(pt, C.uint(ptLen), ck, cryptoBuf, &cryptoLen)
	cryptoText = C.GoStringN(cryptoBuf, cryptoLen)
	return cryptoText
}

func TEADecrypt(cryptoText string, cryptoKey string) (plainText string) {
	ct := C.CString(cryptoText)
	defer C.free(unsafe.Pointer(ct))
	ctLen := len(cryptoText)
	ck := C.CString(cryptoKey)
	defer C.free(unsafe.Pointer(ck))

	var plainLen C.uint
	plainBuf := unsafe.Pointer(C.malloc(CCryptoBufSize))
	defer C.free(plainBuf)
	C.xTEADecryptWithKey(ct, C.uint(ctLen), ck, plainBuf, &plainLen)
	plainText = C.GoStringN(plainBuf, plainLen)
	return plainText
}
