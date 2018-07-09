// OrderSyncTest project main.go
package main

import (
	"fmt"
	"redis"
	"time"
)

var (
	redis_list *[10]redis.Client
)

func main() {
	redis_init()

	for i := 0; i < 10; i++ {
		go Sync(&redis_list[i], i+1)
	}

	for {
		time.Sleep(1000 * time.Millisecond)
	}
}

func redis_init() bool {
	redis_list = new([10]redis.Client)
	var err redis.Error

	var ports = [10]int{
		4501, 4502, 4503, 4504, 4505,
		4506, 4507, 4508, 4509, 4510}

	for i := 0; i < 10; i++ {
		spec := new(redis.ConnectionSpec)

		spec.Host("192.168.100.167")
		spec.Port(ports[i])
		fmt.Printf("Ports[%d]: %d\n", i, ports[i])

		redis_list[i], err = redis.NewSynchClientWithSpec(spec)

		if err != nil {
			return false
		}
	}

	return true
}

func Sync(redis *redis.Client, seq int) {
	var (
		json          string
		format_insert string
		format_update string
	)

	format_insert = "{\"id\": %d, \"sql\":\"INSERT INTO `paycenter_order` (`mid`, `sitemid`, `buyer`, `sid`, `appid`, `pmode`, `pamount`, `pcoins`, `pchips`, `pcard`, `pnum`, `payconfid`, `pcoinsnow`, `pdealno`, `pbankno`, `desc`, `pstarttime`, `pendtime`, `pstatus`, `pamount_rate`, `pamount_unit`, `pamount_usd`, `ext_1`, `ext_2`, `ext_3`, `ext_4`, `ext_5`, `ext_6`, `ext_7`, `ext_8`, `ext_9`, `ext_10`) VALUES (%d, '100002112615471', '100002112615471', 4, 5, 27, 3, 0, 27000, 0, 1, 25296, 0, '0', '0', '0', 1408464000, 0, 0, 1, 'USD', 3, 0, 0, 0, '', '', '', '', '', '', '0101');\", \"time\":140000000, \"type\":0 }"
	format_update = "{\"id\": %d, \"sql\":\"update `paycenter_order` set pstatus=2,ext_7='201' where pid=%d\", \"time\":140000050, \"type\":1 }"

	for i := 1; i <= 100000; i++ {
		if i%10 == seq {
			json = fmt.Sprintf(format_insert, i, i)
			(*redis).Rpush("ORDER_Q", []byte(json))
			json = fmt.Sprintf(format_update, i, i)
			(*redis).Rpush("ORDER_Q", []byte(json))
			time.Sleep(1 * time.Millisecond)
		} else if i%10 == 0 && seq == 10 {
			json = fmt.Sprintf(format_insert, i, i)
			(*redis).Rpush("ORDER_Q", []byte(json))
			json = fmt.Sprintf(format_update, i, i)
			(*redis).Rpush("ORDER_Q", []byte(json))
			time.Sleep(1 * time.Millisecond)
		}

	}
}
