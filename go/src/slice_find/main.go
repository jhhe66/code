package main

import (
	"fmt"
)

func main() {
	uids := make([]uint64, 1000000)

	for idx, _ := range uids {
		uids[idx] = uint64(idx)
		fmt.Printf("idx: %d v: %d\n", idx, uids[idx])
	}
}
