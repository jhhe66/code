package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"fmt"
)

type CCSConf struct {
	Id 			uint64 `yaml: "id"`
	App_Id 		uint64 `yaml: "app_id"`
	App_Conf 	string `yaml: conf_file`
	Command  	string `yaml: command`
	Api_Url		string `yaml: api_url`
	Log_Path	string `yaml: "log_path"`
}

func main() {
	content, err := ioutil.ReadFile("cc2.yml")
	if err != nil {
		fmt.Printf("err: %s\n", err.Error())
		return
	}

	cc := CCSConf{}

	err = yaml.Unmarshal(content, &cc)
	if err != nil {
		fmt.Printf("err: %s\n", err.Error())
		return
	}

	fmt.Printf("cc: %v\n", cc)


	ccc := CCSConf{10, 100, "/opt/push", "notify.sh", "http://ifeng.com", "/opt/push/log"}
	
	str, ccc_err := yaml.Marshal(&ccc)
	if ccc_err != nil {
		fmt.Printf("err: %s\n", err.Error())
		return
	}

	fmt.Printf("str: %s\n", str)
	
}
