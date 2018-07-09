// mutli project main.go
package main

import (
	"fmt"
)

func mutli_return() int {
	return 10
}

func main() {
	fmt.Println("Hello World!")

	is := false

	if p, is := mutli_return(); is {
		fmt.Printf("is\n")
	}
}
