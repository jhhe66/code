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
	msgObj.Set("registration_id", "d5DcADzIF70:APA91bE48uCQjy5K-iB9uFJagtx5BhbkqomPw4zFsiHwwwlw_h6NXM43_YhlaIgiaLGy8ziC8_wB4VNpOvkJjsnah4wM1AiuGknG_4g-EVVvPUDLSrPwhHUCCJNJhA2yyt1QDsAKkinD")
	msgObj.Set("push_type", "text")
	msgObj.Set("from_id", 1946612)
	msgObj.Set("sender", "ken")
	msgObj.Set("push_param", "HelloTalk")
	msgObj.Set("msg_id", "d5789ab0b34afe12p4c22k46q1481888504625464")
	msgObj.Set("sound", 1)
	msgObj.Set("to_id", 5037760)
	msgObj.Set("action_id", 0)
	msgObj.Set("by_at", 0)
	message, err := msgObj.MarshalJSON()
	if err != nil {
		fmt.Printf("MarshalJSON failed")
		return
	}
	err = w.Publish("test", message)
	if err != nil {
		log.Panic("Could not connect")
	}
	fmt.Printf("Publish success")
	w.Stop()
}
