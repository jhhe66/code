package main

import (
	"fmt"
	"runtime"
)

func main() {
	name, file, line, ok := callName(0)
	if ok {
		fmt.Printf("%s:%d-%s\n", file, line, name)
	}
}

func callName(skip int) (name, file string, line int, ok bool) {
	var pc uintptr

	if pc, file, line, ok = runtime.Caller(skip + 2); !ok {
		return
	}

	name = runtime.FuncForPC(pc).Name()
	return
}

