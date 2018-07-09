// server
package main

import (
	"fmt"
	"io/ioutil"
	"log4go"
	"net/http"
	"redis"
	"runtime"
	"strconv"
)

type HttpTask struct {
	Params string
}

var (
	TaskRedisConf  RedisConf
	TaskServerConf ServerConf
	TaskUrlConf    UrlConf
	TaskId         uint64 = 0
)

func InitConf() {
	GetUrlConf(&TaskUrlConf)
	GetServerConf(&TaskServerConf)
	GetRedisConf(&TaskRedisConf)
	runtime.GOMAXPROCS(int(TaskServerConf.MaxRoutines))
}

func Run() {
	var c []chan *HttpTask = make([]chan *HttpTask, TaskServerConf.MaxRoutines)

	client := MakeRedis(TaskRedisConf)

	if client == nil {
		log4go.Debug("redis make error.")
		return
	}

	for i := uint16(0); i < TaskServerConf.MaxRoutines; i++ {
		c[i] = make(chan *HttpTask, TaskServerConf.ChannelQueueLen)

		go DoTask(c[i], i)
	}

	for {
		task, e := GetTask(&client)

		//获取失败，重新创建redis实例
		if e != nil {
			client = MakeRedis(TaskRedisConf)
			continue
		}

		//判断是否取到的内容为空
		if len(task.Params) == 0 {
			continue
		}

		TaskId++

		// 分配任务
		c[TaskId%uint64(TaskServerConf.MaxRoutines)] <- task

		//for i := uint16(0); i < TaskServerConf.MaxRoutines; i++ {
		//	c[i] <- GetTask()
		//}
	}
}

func DoTask(c chan *HttpTask, i uint16) {
	//执行HttpTask

	defer close(c)

	var (
		input   *HttpTask
		ok      bool
		chan_id uint16 = i
	)

	for {
		input, ok = <-c

		if !ok { //channel管理就退出
			log4go.Error("%d channel closed.", chan_id)
			break
		}

		//执行任务
		log4go.Debug("Input[%d]: %s", i, input.Params)

		url := fmt.Sprintf("%s%s", TaskUrlConf.Url[0], input.Params)

		log4go.Debug("URL: %s", url)

		for i := 0; i < TaskUrlConf.Retry; i++ {
			if HttpRequest(&url) {
				log4go.Debug("Request Success.")
				break
			} else {
				log4go.Debug("Request Failed.")
				continue
			}
		}

	}
}

func GetTask(client *redis.Client) (*HttpTask, error) {
	value, e := (*client).Lpop(TaskRedisConf.Key)

	if e != nil {
		log4go.Error("Lpop Error: %v", e)
		return nil, e
	}

	return &HttpTask{Params: string(value)}, nil
}

func MakeRedis(r RedisConf) redis.Client {
	//生成redis的实例
	spec := new(redis.ConnectionSpec)

	spec.Host(r.Host)
	spec.Port(int(r.Port))

	client, err := redis.NewSynchClientWithSpec(spec)

	if err != nil {
		return nil
	}

	return client
}

func HttpRequest(url *string) bool {
	rsp, error := http.Get(*url)

	if error != nil {
		log4go.Error("Response Error: %v", error)

		return false
	}

	content, err := ioutil.ReadAll(rsp.Body)

	if err != nil {
		log4go.Error("Read Error: %v", error)
		return false
	}

	if r, _ := strconv.Atoi(string(content)); r == 0 {
		return true
	} else {
		return false
	}

	return true
}

func ConfTest() {
	log4go.Debug("--------ServerConf--------")
	log4go.Debug("MaxRoutines: %d", TaskServerConf.MaxRoutines)
	log4go.Debug("ChannelQueueLen: %d", TaskServerConf.ChannelQueueLen)
	log4go.Debug("--------ServerConf--------")
	log4go.Debug("--------RedisConf--------")
	log4go.Debug("Host: %s", TaskRedisConf.Host)
	log4go.Debug("Port: %d", TaskRedisConf.Port)
	log4go.Debug("Key: %s", TaskRedisConf.Key)
	log4go.Debug("--------RedisConf--------")
	log4go.Debug("--------UrlConf--------")
	log4go.Debug("MaxUrls: %d", TaskUrlConf.MaxUrls)
	log4go.Debug("Url: %v", TaskUrlConf.Url)
	log4go.Debug("Retry: %d", TaskUrlConf.Retry)
	log4go.Debug("--------UrlConf--------")
}
