// func_param project main.go
package main

import (
	"fmt"
)

func change(i int) {
	i = 100
}

func changeP(p *int) {
	*p = 100
}

func main() {
	i := 99

	fmt.Printf("begin: %d\n", i)

	change(i)

	fmt.Printf("end: %d\n", i)
}
