package main

import "fmt"

func main() {
	slice := []int{1, 2, 3, 4, 5}

	array := [...]int{6, 7, 8, 9, 10}

	fmt.Printf("%v \n", slice)
	fmt.Printf("%v \n", array)
}