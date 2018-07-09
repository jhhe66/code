package main

/*
#cgo CFLAGS: -I../include
#cgo LDFLAGS: -L../libs -lcfun

#include "cfun.h"
*/
import "C"

import "fmt"

func Sum(a int, b int) int {
    return int(C.sum(C.int(a), C.int(b)))
}

func main() {
	fmt.Printf("%d\n", Sum(10, 20))
}
