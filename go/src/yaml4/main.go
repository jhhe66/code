package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type CCSVersion struct {
	Version uint64 `yaml: "version"`
}

func main() {
	data, err := ioutil.ReadFile("cc.yml")
	if err != nil {
		fmt.Printf("err: %s\n", err.Error())
		return
	}

	cc := &CCSVersion{}
	err = yaml.Unmarshal(data, &cc)
	if err != nil {
		fmt.Printf("err: %s\n", err.Error())
		return
	}

	fmt.Printf("CC: %v\n", cc)
}
