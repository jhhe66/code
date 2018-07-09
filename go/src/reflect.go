package main

import (
			"fmt"
			"reflect"
		)


type Person struct {
	name string "namestr"
	age int
}

func ShowTag(i interface{}) {
	switch t := reflect.TypeOf(i); t.Kind() {
	case reflect.Ptr:
		tag := t.Elem().Field(0).Tag
		fmt.Printf("%v\n", tag)
	}
}

func show(i interface{}) {
	switch t := i.(type) {
	case *Person:
		tt := reflect.TypeOf(i)
		v := reflect.ValueOf(i)
		tag := tt.Elem().Field(0).Tag
		name := v.Elem().Field(0).String()

		fmt.Printf("%v\n", tt)
		fmt.Printf("%v\n", v)
		fmt.Printf("%v\n", tag)
		fmt.Printf("%v\n", name)

		fmt.Printf("%v\n", t)
	}


}

func main() {
	p1 := new(Person)

	ShowTag(p1)
	show(p1)
}