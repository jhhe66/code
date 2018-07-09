package main

import 
(
	"fmt"
	"encoding/json"
)

type School struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func main() {
	sch := School{ID: 10, Name: "chenbo"}
	jstr, _ := json.Marshal(sch)
	fmt.Printf("%s\n", jstr) 
}
