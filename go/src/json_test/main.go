// json_test project main.go
package main

import (
	"encoding/json"
	"fmt"
)

type ColorGroup struct {
	ID     int
	Name   string
	Colors []string
	Bags    map[string]string `json: "bags"`
}

func main() {
	json_str := `{"ID":1,"Name":"Reds","Colors":["Crimson","Red","Ruby","Maroon"]}`
	//json_str := `{"ID":1,"Name":"Reds","Colors":["Crimson","Red","Ruby","Maroon"], "bags": { "big": "channel", "small": "herms"} }`

	var g ColorGroup

	error := json.Unmarshal([]byte(json_str), &g)

	if error == nil {
		fmt.Printf("%+v\n", g)
		fmt.Printf("Big: %s\n", g.Bags["big"])

		for k,v := range g.Colors {
			fmt.Printf("K: %v V: %v\n", k, v)
		}

		for bags := range g.Bags {
			fmt.Printf("Bags: %v\n", bags)
		}
		for k,v := range g.Bags {
			fmt.Printf("Bags: %v %v\n", k, v)
		}
	}
}
