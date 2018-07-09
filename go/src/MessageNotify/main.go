// MessageNotify project main.go
package main

import (
	"encoding/json"
	log "log4go"
	"redis"
	"runtime"
	"time"
)

const (
	log_conf_file = "log.xml"
)

type NotifyMsg struct {
	From string
	To   string
	Url  string
	Time int64
	Type int
}

func Mail_Send(m *Mail) {
	log.Debug("-------- Mail Begin --------")
	err := m.Send()
	if err != nil {
		log.Debug("err: %v", err)
	} else {
		log.Debug("send success.")
	}
	log.Debug("-------- Mail End --------")

}

// {"from":"100013", "to":"18688770050", "msg":"有异常信息", "time":145555000, "type":1}

func NewRedisClient(r *RedisConf) redis.Client {
	spec := new(redis.ConnectionSpec)

	spec.Host(r.Host)
	spec.Port(r.Port)

	client, err := redis.NewSynchClientWithSpec(spec)

	if err != nil {
		return nil
	}

	return client
}

func JobSms(c chan *SmsConf, id uint) {
	for {
		sms, ok := <-c

		if !ok {
			log.Debug("JobSms Read Error: %v", ok)
		}

		err := (*sms).Send()
		if err != nil {
			for i := 0; i < 3; i++ {
				if err := (*sms).Send(); err == nil {
					break
				}
			}
		}
	}
}

func main() {
	var (
		notify_msg NotifyMsg
		msg_id     uint64 = 0
	)

	LogInit(log_conf_file)

	server_conf := new(ServerConf)
	server_conf.ReadServerConf()
	server_conf.PrintSelf()

	runtime.GOMAXPROCS(server_conf.MaxWorker)

	sms := new(SmsConf)
	sms.ReadSmsConf()
	sms.PrintSelf()

	sms_channels := make([]chan *SmsConf, server_conf.MaxWorker) // slice chan
	for i := 0; i < server_conf.MaxWorker; i++ {
		sms_channels[i] = make(chan *SmsConf, server_conf.ChanneLength)
		go JobSms(sms_channels[i], uint(i))
	}

	redis_conf := new(RedisConf)
	redis_conf.ReadRedisConf()
	redis_conf.PrintSelf()
	redis := NewRedisClient(redis_conf)

	if redis == nil {
		log.Debug("redis new error.")
		return
	}

	for {
		if err := redis.Ping(); err != nil {
			log.Debug("redis has down.")
			continue
		}

		msg, err := redis.Lpop(redis_conf.Key)
		if err != nil {
			log.Debug("Lpop error: %v", err)
		} else if len(msg) == 0 {
			time.Sleep(10 * time.Second)
			continue
		} else {
			log.Debug("JSON: %s", msg)
		}

		jerr := json.Unmarshal([]byte(msg), &notify_msg)
		if jerr != nil {
			log.Debug("json parse failed: %v", jerr)
		}

		/* 如果这里To相同的话，由于平台的限制不能向统一个号码30s内发送短信 */
		//if strings.EqualFold(sms.To, notify_msg.To) {
		//	go func(c chan *SmsConf, s *SmsConf) {
		//		sms.Url = notify_msg.Url
		//		time.Sleep(30 * time.Second)
		//		c <- s
		//	}(sms_channels[msg_id%uint64(server_conf.MaxWorker)], sms)
		//} else {
		//	sms.To = notify_msg.To
		//	sms.Url = notify_msg.Url
		//	sms_channels[msg_id%uint64(server_conf.MaxWorker)] <- sms
		//}
		sms.To = notify_msg.To
		sms.Url = notify_msg.Url

		time.Sleep(30 * time.Second)

		sms_channels[msg_id%uint64(server_conf.MaxWorker)] <- sms
		msg_id++
	}
}

func LogInit(conf string) {
	log.LoadConfiguration(conf)
}
