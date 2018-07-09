package main

import (
	"config"
	log "log4go"
)

const (
	conf_file_path = "server.ini"
)

type RedisConf struct {
	Host string
	Port int
	Key  string
}

type SmsConf struct {
	Api   string
	Msg   string
	From  string
	To    string
	Pri   int
	Appid int
	Key   string
	Url   string
}

type ServerConf struct {
	MaxWorker    int
	ChanneLength int
	Retry        int
}

func (rc *RedisConf) ReadRedisConf() {
	conf, _ := config.ReadDefault(conf_file_path)

	rc.Host, _ = conf.String("redis", "host")
	rc.Port, _ = conf.Int("redis", "port")
	rc.Key, _ = conf.String("redis", "key")
}

func (rc *RedisConf) PrintSelf() {
	log.Debug("--------------------")
	log.Debug("redis->host: %s", rc.Host)
	log.Debug("redis->port: %d", rc.Port)
	log.Debug("redis->key:  %s", rc.Key)
	log.Debug("--------------------")
}

func (sc *SmsConf) ReadSmsConf() {
	log.Debug("-------- ReadSmsConf begin --------")
	conf, err := config.ReadDefault(conf_file_path)
	if err != nil {
		log.Debug("err: %v", err)
		return
	}

	sc.Api, _ = conf.String("sms", "api")
	sc.Appid, _ = conf.Int("sms", "appid")
	sc.From, _ = conf.String("sms", "from")
	sc.Pri, _ = conf.Int("sms", "pri")
	sc.Key, _ = conf.String("sms", "key")

	log.Debug("-------- ReadSmsConf end --------")
}

func (sc *SmsConf) PrintSelf() {
	log.Debug("-------- PrintSelf begin --------")
	log.Debug("sms->api: 	%s", sc.Api)
	log.Debug("sms->appid: %d", sc.Appid)
	log.Debug("sms->key:  	%s", sc.Key)
	log.Debug("sms->from:  %s", sc.From)
	log.Debug("sms->pri:  	%d", sc.Pri)
	log.Debug("-------- PrintSelf end --------")
}

func (sc *ServerConf) ReadServerConf() {
	conf, _ := config.ReadDefault(conf_file_path)

	sc.MaxWorker, _ = conf.Int("server", "max_worker")
	sc.ChanneLength, _ = conf.Int("server", "channel_len")
	sc.Retry, _ = conf.Int("server", "Retry")
}

func (sc *ServerConf) PrintSelf() {
	log.Debug("-------- PrintSelf begin --------")
	log.Debug("server->MaxWorker:%d", sc.MaxWorker)
	log.Debug("server->ChannelLength:%d", sc.ChanneLength)
	log.Debug("server->Retry:%d", sc.Retry)
	log.Debug("-------- PrintSelf end --------")
}
