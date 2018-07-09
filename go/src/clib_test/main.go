// clib_test project main.go
package main

/*
#cgo CFLAGS: -I../include
#cgo LDFLAGS: -L../libs -lcfun

#include "cfun.h"
*/
import "C"

import (
	"clib_warper"
	"fmt"
)

func main() {
	fmt.Printf("%d\n", C.sum(1, 2))
	fmt.Printf("%d\n", clib_warper.Sum(3, 3))
}
