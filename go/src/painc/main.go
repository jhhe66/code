package main

import (
	"fmt"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("painc: %v\n", err)
		}
	}()

	raise()
}

func raise() {
	panic("make a painc.")
}
