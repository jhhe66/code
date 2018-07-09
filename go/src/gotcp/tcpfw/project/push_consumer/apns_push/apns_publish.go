package main

import (
	"fmt"
	"log"

	simplejson "github.com/bitly/go-simplejson"
	"github.com/nsqio/go-nsq"
)

const (
	CTP2P = 0
	CTMUC = 1
)

const (
	PT_TEXT              = 0
	PT_VOICE             = 1
	PT_PHOTO             = 2
	PT_INDRODUCE         = 3
	PT_LOCATION          = 4
	PT_FRIEND_INVITE     = 5
	PT_LANGUAGE_EXCHANGE = 6
	PT_CORRECT_SENTENCE  = 7
	PT_STICKERS          = 8
	PT_DOODLE            = 9
	PT_GIFT              = 10
	PT_VOIP              = 11
	PT_INVITE_ACCEPT     = 12
	PT_VIDEO             = 13

	PT_GVOIP               = 15
	PT_LINK                = 16
	PT_CARD                = 17
	PT_FOLLOW              = 18
	PT_REPLY_YOUR_COMMENT  = 19
	PT_COMMENTED_YOUR_POST = 20
	PT_CORRECTED_YOUR_POST = 21
	PT_MOMENT_LIKE         = 22
)

const (
	PUSH_SOUND_CALLING = "voipcall.caf"
	PUSH_SOUND_DEFAULT = "default"
)

func main() {
	config := nsq.NewConfig()
	w, _ := nsq.NewProducer("127.0.0.1:4150", config)
	msgObj := simplejson.New()
	msgObj.Set("chat_type", CTP2P)
	msgObj.Set("from_id", 1946612)
	msgObj.Set("to_id", 2325928)
	msgObj.Set("push_type", PT_TEXT)
	msgObj.Set("preview", 1)
	msgObj.Set("token", "c63716dc02001f2ed37d941978e72e6fe499a0c45b3960cf9c964fadc9d57836")
	msgObj.Set("badge_num", 1)
	msgObj.Set("sound", PUSH_SOUND_DEFAULT)
	msgObj.Set("nick_name", "ken")
	msgObj.Set("push_param", "hello songliwei")
	msgObj.Set("msg_id", 0)
	msgObj.Set("expired", 0)
	msgObj.Set("action_id", 0)
	msgObj.Set("by_at", 0)
	message, err := msgObj.MarshalJSON()
	if err != nil {
		fmt.Printf("MarshalJSON failed")
		return
	}
	err = w.Publish("apns", message)
	if err != nil {
		log.Panic("Could not connect")
	}
	fmt.Printf("Publish success")
	w.Stop()
}
