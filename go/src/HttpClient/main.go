package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"
	"os"
)



func Get(url string) {
	var ts = &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 10 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,

		TLSHandshakeTimeout: 5 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	var client = &http.Client{
		Timeout: time.Second * 30,
		Transport: ts,
	}
	
	response, _ := client.Get(url)
	defer response.Body.Close()

	if response.StatusCode == 200 {
		body, _ := ioutil.ReadAll(response.Body)
		fmt.Println(string(body))
	}
}

func main() {
	if len(os.Args) == 1 {
		fmt.Printf("usage: HttpClient http://xxx.xxx.xxx\n")
		return 
	}
	Get(os.Args[1])
}
