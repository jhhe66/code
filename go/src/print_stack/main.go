package main

import (
	"fmt"
	"runtime/debug"
)

func main() {
	fmt.Printf("%s\n", debug.Stack())	
	debug.PrintStack();	
}
