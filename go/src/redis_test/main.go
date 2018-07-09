// redis_test project main.go
package main

import (
	"fmt"
	"redis"
)

func main() {
	spec := new(redis.ConnectionSpec)

	//spec.Db(0)
	spec.Host("192.168.100.154")
	spec.Port(6380)

	client, err := redis.NewSynchClientWithSpec(spec)

	if err != nil {
		fmt.Printf("create client error: %v\n", err)
	}

	for {
		value, e := client.Lpop("chenbo")

		if e != nil {
			fmt.Printf("Get Error: %s\n", e)
		}

		if len(value) > 0 {
			fmt.Printf("value: %s\n", value)
		}
	}
}
