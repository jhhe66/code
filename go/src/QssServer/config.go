// config
package main

import (
	"config"
	//"strconv"
	//"strings"
)

type UrlConf struct {
	Url     []string // slice
	MaxUrls uint16
	Retry   uint16
}

type ServerConf struct {
	MaxRoutines     uint16
	ChannelQueueLen uint16
}

type RedisConf struct {
	Host string
	Port uint16
	Key  string
}

const (
	ConfigFile = "config.ini"
)

func GetUrlConf(u *UrlConf) {
	c, _ := config.ReadDefault(ConfigFile)

	u.Retry, _ = c.Uint16("url", "Retry")
	u.MaxUrls, _ = c.Uint16("url", "MaxUrls")
	u.Url = make([]string, u.MaxUrls)

	for i := uint16(0); i < u.MaxUrls; i++ {
		//urls := []string{"url", strconv.Itoa(int(i + 1))}

		u.Url[i], _ = c.String("url", "url")
	}

}

func GetServerConf(s *ServerConf) {
	c, _ := config.ReadDefault(ConfigFile)

	s.MaxRoutines, _ = c.Uint16("server", "MaxRoutines")
	s.ChannelQueueLen, _ = c.Uint16("server", "ChannelQueueLen")
}

func GetRedisConf(r *RedisConf) {
	c, _ := config.ReadDefault(ConfigFile)

	r.Host, _ = c.String("redis", "host")
	r.Port, _ = c.Uint16("redis", "port")
	r.Key, _ = c.String("redis", "key")
}
