package main

import (
	_ "fmt" // 这里使用下划线只是调用fmt package 的init function，不能通过fmt引用其方法
)

func main() {
	fmt.Printf("hello world\n")
}
