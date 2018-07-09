package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)


type JMQConf struct {
	ID 			string  `yaml: "id"`
	Exchange 	struct {
					Autodel bool `yaml: "autodel"`
					Confirm bool `yaml: "confirm"`
					Durable bool `yaml: "durable"`
					Name	string `yaml: "name"`
					Passive bool `yaml: "passive"`
					Rkey	string `yaml: "rkey"`
				} 
	Queue 		struct {
					Autodel 		bool `yaml: "autodel"`
					Durable 		bool `yaml: "durable"`
					Exclusive 		bool `yaml: "exclusive"`
					Name 			string `yaml: "name"`
					Needack 		bool `yaml: "needack"`
					Passive 		bool `yaml: "passive"`
					Prefetch_cnt 	int `yaml: "prefetch_cnt"`
					Rkey			string `yaml: "rkey"`
				}
	Srv			struct {
					Server 	string `yaml: "server"`		
				}
	Vhost		struct {
					Name 		string `yaml: "name"`
					Password 	string `yaml: "password"`
					User		string `yaml: "user"`
				}
}


func LoadFileAll(file string) ([]byte, error) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Printf("%v", err);
	}
	
	return content, err
}

func main() {
	content, err := LoadFileAll("jmq2.yml")
	if err != nil {
		fmt.Print("error: %v\n", err)
		return
	}

	fmt.Printf("%s\n", content)	

	conf := JMQConf{}

	err = yaml.Unmarshal(content, &conf)
	if err != nil {
		fmt.Printf("Unmarshal error: %v\n", err)
		return 
	}

	fmt.Printf("Conf: %v\n", conf)

	y, _ := yaml.Marshal(&conf)
	fmt.Printf("Yaml: %s\n", y)
}
