// producer project main.go
package main

import (
	"fmt"
)

func producer(c chan int) {
	for i := 0; i < 10; i++ {
		fmt.Printf("Alice puts product, ID is : %d \n", i)
		c <- i
	}
	defer close(c)
}
func consumer(c chan int) {
	hasMore := true
	var p int
	for hasMore {
		if p, hasMore = <-c; hasMore {
			fmt.Printf("Bob gets product, ID is : %d \n", p)
		}
	}
}

func consumer_loop(c chan int) {
	hasMore := true
	for hasMore {
		if p, hasMore := <-c; hasMore {
			fmt.Printf("Bob gets product, ID is : %d \n", p)
		}
	}
}

func main() {
	c := make(chan int)
	go producer(c)
	consumer(c)
}
