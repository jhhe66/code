package main

import p "fmt"

func tt(i interface{}) {
	switch t := i.(type) {
	case string:
		p.Printf("String\n")
		p.Printf("%v\n", t)
	case int:
		p.Printf("int\n")
		p.Printf("%v\n", t)

	} 
}

func main() {
	var i int = 99

	tt(i)
}