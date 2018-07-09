package main

import (
	"fmt"
	"log"

	simplejson "github.com/bitly/go-simplejson"
	"github.com/nsqio/go-nsq"
)

func main() {
	config := nsq.NewConfig()
	w, _ := nsq.NewProducer("127.0.0.1:4150", config)
	msgObj := simplejson.New()
	msgObj.Set("chat_type", "CHAT")
	msgObj.Set("token", "363240dd6ee6836c488c84b3ecaeae936f51512c")
	msgObj.Set("title", "ricky")
	msgObj.Set("content", "I will now")
	msgObj.Set("msg_type", "text")
	msgObj.Set("sound", 1)
	msgObj.Set("lights", 1)
	msgObj.Set("from_id", 2325928)
	msgObj.Set("msg_id", "5185584_1481943602214_60")
	msgObj.Set("5035860")
	msgObj.Set("actionid", 0)
	msgObj.Set("byAt", 0)

	message, err := msgObj.MarshalJSON()
	if err != nil {
		fmt.Printf("MarshalJSON failed")
		return
	}
	err = w.Publish("xinge", message)
	if err != nil {
		log.Panic("Could not connect")
	}
	fmt.Printf("Publish success")
	w.Stop()
}
