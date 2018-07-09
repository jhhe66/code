// slice_test project main.go
package main

import (
	"fmt"
)

func main() {
	arr := [5]int{1, 2, 3, 4, 5}

	sa := arr[1:4]

	fmt.Printf("cat(sa): %d\n", cap(sa))
}
