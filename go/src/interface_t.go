package main

import "fmt"

type S struct {i int}

func (p *S)Get() int {
	return p.i
}

func (p *S)Put(v int) {
	p.i = v
}


type I interface {
	Get() int
	Put(int)
}


func f(p I) {
	fmt.Printf("%d\n", p.Get())
	p.Put(99)
	fmt.Printf("%d\n", p.Get())
}

func main() {
	s := new(S)

	var ss = S{100}

	fmt.Printf("%d\n", ss.i)

	s.i = 1;

	f(s)
	
	f(&ss)
}
