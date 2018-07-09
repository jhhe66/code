package main
import (
	"time"
)

func main() {
	go work()
	select {}
}

func work() {
	var a string
	for {
		v := "aaaaaa"
		a = a + v
		time.Sleep(time.Millisecond * 1)
	}
}
