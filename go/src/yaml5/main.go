package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type KeysConf struct {
	Key string `yaml:"key"`
	File string `yaml:"file"`
}

type ServiceConf struct {
	Id 			uint64 `yaml: "id"`
	App_Conf 	string `yaml: "app_conf"`
	App_Name	string `yaml: "app_name"`
	Command 	string `yaml: "command"`
	Api_Url 	string `yaml: "api_url"`
	Keys    	[]KeysConf `yaml: "keys"` 
}

type CCSConf struct {
	App_Id 		uint64 `yaml: "app_id"`
	Log_Path	string `yaml: "log_path"`
	Log_Name 	string `yaml: "log_name"`
	Log_Level 	int	   `yaml: "log_level"`
	/*Service		[]struct {
					Id 			uint64 `yaml: "id"`
					App_Conf 	string `yaml: "app_conf"`
					App_Name	string `yaml: "app_name"`
					Command 	string `yaml: "command"`
					Api_Url 	string `yaml: "api_url"`
				}
	*/
	Service 	[]ServiceConf `yaml: "service"`
}

func main() {
	data, err := ioutil.ReadFile("cc.yml")
	if err != nil {
		fmt.Printf("read file error:[%s]\n", err.Error())
		return
	}

	cc := &CCSConf{}
	err = yaml.Unmarshal(data, cc)
	if err != nil {
		fmt.Printf("Unmarshal error:[%s]\n", err.Error())
		return 
	}

	fmt.Printf("CC: %v\n", cc)
	//fmt.Printf("Service[0]: %v\n", cc.Service[0])

	for k, v := range cc.Service {
		fmt.Printf("Service[%d]: %v\n", k, v)
	}
}
