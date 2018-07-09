package main

import (
	"fmt"
	"encoding/json"
)


type ConfData struct {
    Key     string `json:"configKey"`
	Value   string `json:"configValue"`
}

type DataRsp struct {
	Code        int     `json:"code"`
	Msg         string  `json:"msg"`
	ClientID    string  `json:"clientId"`
	Version     string  `json:"version"`
	Datas   []ConfData  `json:"configData"`
}

func main() {
	//json_str := `{"msg": "当前appI对应的服务信息尚未注册，请联系管理员！", "code": 1}`	
	json_str := `{"msg":"保存绑定关系成功!","code":0,"clientId":"clientid123456789987654326","configData":[{"configKey":"logback.property","configValue":"log4j.rootLogger=info,console,logFile\r\n\r\nlog4j.appender.console=org.apache.log4j.ConsoleAppender\r\nlog4j.appender.console.layout=org.apache.log4j.PatternLayout\r\nlog4j.appender.console.layout.ConversionPattern=%d - xxl-job-admin - %p [%c] - <%m>%n\r\n\r\nlog4j.appender.logFile=org.apache.log4j.DailyRollingFileAppender\r\nlog4j.appender.logFile.File=/opt/push/data/xxl-job/admin/admin.log\r\nlog4j.appender.logFile.layout=org.apache.log4j.PatternLayout\r\nlog4j.appender.logFile.layout.ConversionPattern=%d - xxl-job-admin - %p [%c] - <%m>%n"},{"configKey":"redis.property","configValue":"# Redis settings\r\nredis.pool.maxActive=64\r\nredis.pool.maxIdle=8\r\nredis.pool.maxWait=1000\r\nredis.pool.testOnBorrow=true\r\nredis.pool.testOnReturn=true\r\nredis.ip=192.168.52.79\r\nredis.port=6379\r\nredis.timeout=10000\r\nredis.password=admin"}],"version":"1509440370197"}`	

	data := DataRsp{}

	err := json.Unmarshal([]byte(json_str), &data)
	if err != nil {
		fmt.Printf("error: %s\n", err.Error())
		return
	}
	
	fmt.Printf("data: %v\n", data)
	fmt.Printf("datas: %v\n", data.Datas)


	data2 := DataRsp{}
	conf := ConfData{}
	
	conf.Key = "aaa"
	conf.Value = "bbb"

	data2.Datas = append(data2.Datas, conf)
	buff, _ := json.Marshal(data2)	

	fmt.Printf("buff: %v\n", string(buff))
}
