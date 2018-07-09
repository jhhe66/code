// struct project main.go
package main

import (
	"fmt"
)

type People struct {
	id [10]byte
}

type A struct {
	name [2]byte
}

var p People

func main() {
	changed_array(p.id[:], 10)

	for i := 0; i < 10; i++ {
		fmt.Printf("p[%v]: %v \n", i, p.id[i])
	}
}

func changed_array(a []byte, n int) {
	for i := 0; i < n; i++ {
		a[i] = byte(i)
	}
}
