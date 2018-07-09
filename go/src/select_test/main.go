package main

import (
	"fmt"
	"time"
)


func stdout(o chan<- int64) {
	defer close(o)
	for {
		o <- time.Now().Unix()	
		time.Sleep(time.Second * 1)
	}
}


func stdin(i <-chan int64) {
	for {
		
	}
}

func main() {
	occ := make([]chan int64, 2)
	for idx := 0; idx < 2; idx++ {
		occ[idx] = make(chan int64, 10)
		go stdout(occ[idx])
	}
	
	sign := make(chan int)

	go func() {
		for {
			select {
				case v1, ok := <- occ[0]:
					if ok {
						fmt.Printf("v1: %v\n", v1)
					} else {
						sign <- 1
					}
				case v2, ok := <- occ[1]:
					if ok {
						fmt.Printf("v2: %v\n", v2)
					} else {
						sign <- 1
					}

				default:
				//fmt.Printf("sleep....\n")
			//	time.Sleep(10 * time.Second)
			}
		}
	}()

	<- sign
}
