package main

import "fmt"

const (
	_ 	= iota
	KB float64 = 1 << (10 * iota)
	MB
	GB
	TB
	PB
	EB
	ZB
	YB
)

func main() {
	fmt.Printf("KB: %f\n", KB)
	fmt.Printf("MB: %f\n", MB)
}