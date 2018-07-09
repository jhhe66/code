package main

import (
	"fmt"
)

var (
	idx int = 0
)
func main() {
	ppi := make([][]int, 10)
	fmt.Printf("ppi.len: %v cap: %v\n", len(ppi), cap(ppi))
	for k, _ := range ppi {
		ppi[k] = make([]int, 10)
		for kk, _ := range ppi[k] {
			ppi[k][kk] = idx
			idx++
		}
	} 

	fmt.Printf("pp: %v\n", ppi)

	var ppl [][]uint64 // 必须要make 或者 new 才能使用

	fmt.Printf("ppl.len: %v cap: %v\n", len(ppl), cap(ppl))
	ppl = make([][]uint64, 10)	
	for k, _ := range ppl {
		ppl[k] = make([]uint64, 10)
		for kk, _ := range ppl[k] {
			ppl[k][kk] = uint64(idx)
			idx++
		}
	} 
	fmt.Printf("pp: %v\n", ppl)

	var pp *[][]int

	pp = new([][]int)
	*pp = make([][]int, 10)
	fmt.Printf("pp.len: %v cap: %v\n", len(*pp), cap(*pp))
	for k, _ := range *pp {
		(*pp)[k] = make([]int, 10)
		for kk, _ := range (*pp)[k] {
			(*pp)[k][kk] = idx
			idx++
		}
	} 
	fmt.Printf("pp: %v\n", pp)

	return
}
