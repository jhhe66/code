// timer_test project main.go
package main

import (
	"fmt"
	"time"
)

func callback() {
	fmt.Printf("Timer Callback\n")
}

func callback_2(i int) {
	fmt.Printf("i: %d\n", i)
}

func main() {
	t := time.AfterFunc(10*time.Second, callback_2)

	if t != nil {
	}

	time.Sleep(1000 * time.Second)
}
