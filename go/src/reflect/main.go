package main

import (
	"fmt"
	"reflect"
)

func struct_reflect() {
	type T struct {
	    A int
		B string
	}
	
	t := T{23, "skidoo"}
	s := reflect.ValueOf(&t).Elem()
	typeOfT := s.Type()
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		fmt.Printf("%d: %s %s = %v\n", i,
		typeOfT.Field(i).Name, f.Type(), f.Interface())
	}
}

func main() {
	var x int = 100

	v := reflect.ValueOf(x)
	fmt.Printf("%v\n", v.Type())
	fmt.Printf("%d\n", v.Int())
	fmt.Printf("%v\n", v.Interface())
	fmt.Println("settability of v:" , v.CanSet())
	
	var y float64 = 3.4
	p := reflect.ValueOf(&y) // 注意：获取 Y 的地址。
	fmt.Println("type of p:", p.Type())
	fmt.Println("settability of p:" , p.CanSet())

	vv := p.Elem()
	fmt.Println("settability of vv:" , vv.CanSet())
	vv.SetFloat(7.1)
	fmt.Println(vv.Interface())
	fmt.Println(y)

	struct_reflect()
}
