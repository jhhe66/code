package main

import (
	. "fmt"
	"time"
)

func main() {	
	list := make(map[int]string, 10)

	for i := 0;i < 10; i++ {
		list[i] = "hello"	
	}
	Printf("size: %d\n", len(list))

	for k, v := range list {
		x := struct {
				k int
				v string} {k, v}
		go func() {
			Printf("[%d]:%s\n", x.k, x.v)
		}()
	}

	time.Sleep(10 * time.Second)
}


