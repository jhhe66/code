// new_make_test project main.go
package main

import (
	"fmt"
)

func main() {
	pa := new([2]int)

	pa[0] = 0
	pa[1] = 1

	fmt.Printf("pa[0]: %d\n", pa[0])
	fmt.Printf("pa[1]: %d\n", pa[1])

}
