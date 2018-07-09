package main

import (
	_ "net/http/pprof"
	"net/http"
	"log"
	"time"
)

func main() {
	go func() {
		log.Println(http.ListenAndServe("0.0.0.0:6060", nil))
	}()

	var cnt uint64 = 0

	for {
		time.Sleep(1 * time.Second)
		cnt++
	}
	//select {}
}
