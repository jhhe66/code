package main

import (
	"fmt"
)

func toUint(i int) interface{} {
	return i
}

func main() {
	fmt.Printf("%v\n", toUint(10).(int))	
}
