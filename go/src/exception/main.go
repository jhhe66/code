package main

import (
	"fmt"
	"time"
)

func main() {
    defer func() { //必须要先声明defer，否则不能捕获到panic异常
		fmt.Println("2")
		if err := recover(); err != nil {
			fmt.Println(err) //这里的err其实就是panic传入的内容，bug
		}
		fmt.Println("3")
	}()
	
	f()
}

func f() {
    for {
		fmt.Println("1")
		panic("bug")
		fmt.Println("4") //不会运行的.
		time.Sleep(1 * time.Second)
	}
}
