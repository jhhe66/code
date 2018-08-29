package  main;

import (
	"fmt"
	//"syscall"
)

func test() {
	fmt.Printf("hello world\n");
}

func preempt() {
	//syscall.Gettid();
}


func main() {
	go test();

	for {
		preempt();
	}
}
