package main

import (
	"strings"
	"fmt"
)

func main() {
	sb := new(strings.Builder)
	sz, err := sb.WriteString("chenbo")
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return 
	}

	fmt.Printf("sz: %d\n", sz)

	return
}
