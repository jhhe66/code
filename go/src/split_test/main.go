// split_test project main.go
package main

import (
	"fmt"
	"strings"
)

func main() {
	f := strings.SplitN("w b c b c c c ", " ", 3)
	fmt.Printf("l: %d\n", len(f))
	fmt.Printf("l: %d\n", len("hello"))
	fmt.Printf("S: %v\n", f[2])
}
