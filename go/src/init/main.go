// init project main.go
package main

import (
	"fmt"
)

/* 默认每个package 首先执行的函数，用于初始化 */
func init() {
	fmt.Printf("this is init function\n")
}

func main() {
	fmt.Println("Hello World!")
}
