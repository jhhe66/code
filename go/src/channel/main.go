package main


import (
	"fmt"
	"time"
)

func main() {
	c := make(chan int, 10)

	go c_read(c)
	
	go c_write(c, 100)

	time.Sleep(10000 * time.Millisecond) // 10s

	fmt.Printf("server exit...\n")
}

func c_read(c chan int) {
	for ;; {
		fmt.Printf("read: %v\n", <- c)
	}

}

func c_write(c chan int, n int) {
	for i := 0; i < n; i++ {
		c <- i
	}
}