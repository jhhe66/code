package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
	"net"
)

var timeout = time.Duration(20 * time.Second)

func dialTimeout(network, addr string) (net.Conn, error) {
	return net.DialTimeout(network, addr, timeout)
}

func main() {
	tr := &http.Transport{
		//使用带超时的连接函数
		Dial: dialTimeout,
		//建立连接后读超时
		ResponseHeaderTimeout: time.Second * 2,
	}
	client := &http.Client{
		Transport: tr,
		//总超时，包含连接读写
		Timeout: timeout,
	}
	
	req, _ := http.NewRequest("GET", "http://www.haiyun.me", nil)
	req.Header.Set("Connection", "keep-alive")
	res, err := client.Do(req)
	if err != nil {
		return
	}
	
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(body))
	for k, v := range res.Header {
		fmt.Println(k, strings.Join(v, ""))
	}

}
