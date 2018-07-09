// clib_warper project clib_warper.go
package clib_warper

/*
#cgo CFLAGS: -I../include
#cgo LDFLAGS: -L../libs -lcfun

#include <stdlib.h>
#include <stdio.h>
#include <cfun.h>
*/
import "C"

//export Sum
func Sum(a int, b int) int {
	return int(C.sum(C.int(a), C.int(b)))
}

