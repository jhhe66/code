// http_get project main.go
package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func http_get(url string) int32 {
	fmt.Println(url)

	rsp, error := http.Get(url)

	if error != nil {
		fmt.Printf("Error: %v\n", error)

		return -1
	}

	content, err := ioutil.ReadAll(rsp.Body)

	if err != nil {
		fmt.Printf("Error: %v\n", error)
	}

	fmt.Printf("html: %s\n", content)

	return 0
}

func main() {
	url := "http://www.ifeng.com"

	go http_get(url)

	time.Sleep(20 * time.Second)
}
