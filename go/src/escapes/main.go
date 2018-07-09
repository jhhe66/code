package main

import (
    "fmt"
)

func main() {
	s := []byte("")
	println(s) // 添加用于打印信息
	
	s1 := append(s, 'a')
	s2 := append(s, 'b')

	//fmt.Println(s1, "===", s2)
	fmt.Println(string(s1), string(s2))
}
