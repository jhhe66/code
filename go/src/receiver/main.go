package main

import (
	"fmt"
)

type Man struct {
	Id 		int
	Age 	int
	Name 	string
}

type List []int

func (m *Man) setAge(age int) {
	m.Age =  age
}

func (m *Man) setId(id int) {
	m.Id = id
}

func (m *Man) Print() {
	fmt.Printf("%v\n", *m)
}

func (l *List) Add(i int) {
	*l = append(*l, i)
}

func (l *List) Print() {
	fmt.Printf("%v\n", *l)
}

func main() {
	var m Man
	
	m.setAge(10)
	m.setId(100)
	m.Print()

	var l List

	l.Add(10)
	l.Print()
}
