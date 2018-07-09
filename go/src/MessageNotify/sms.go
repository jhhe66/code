package main

import (
	"io/ioutil"
	log "log4go"
	"net/http"
)

func (s *SmsConf) Send() error {
	//value := url.Values{}

	//value.Add("from", s.From)
	//value.Add("msg", s.Msg)
	//value.Add("pri", "4")
	//value.Add("to", s.To)
	//param := value.Encode()

	//log.Debug("url: %s", param)

	//md5_str := fmt.Sprintf("%x", md5.Sum([]byte(param+s.Key)))
	//log.Debug("md5: %s", md5_str)

	//url := s.Api + param + "&sig=" + md5_str

	log.Debug("url: %s", s.Url)
	rsp, err := http.Get(s.Url)

	if err != nil {
		log.Debug("send error: %v", err)
		return err
	}

	content, e := ioutil.ReadAll(rsp.Body)

	if e == nil {
		log.Debug("HTTP: %s", string(content))
		return nil
	} else {
		log.Debug("E: %v", e)
		return e
	}
}
