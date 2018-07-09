package main

import "fmt"

type S struct {i int}

func (p *S)Get() int {
	return p.i
}

func (p *S)Put(v int) {
	p.i = v
}


type R struct {i int}

func (p *R)Get() int {
	return p.i
}

func (p *R)Put(v int) {
	p.i = v
}


type I interface {
	Get() int
	Put(int)
}

func f(p I) {
	switch t := p.(type) {
		case *S:
			fmt.Printf("*S\n")
			fmt.Printf("%v\n", t)
		case *R:
			fmt.Printf("*P\n")
			fmt.Printf("%v\n", t)
/*		case S:
			fmt.Printf("S\n")
		case R:
			fmt.Printf("P\n")*/
		default:
			fmt.Printf("O\n")
			fmt.Printf("%v\n", t)
	}
}

func g(some interface{}) int {
	return some.(I).Get()
}

func main() {
	var s S = S{1}
	var r R = R{1}

	s.i = 10
	r.i = 20

	g(s)

	f(&s)
	f(&r)

/*	s := new(S)
	r := new(R)

	s.i = 10
	r.i = 20

	f(s)
	f(r)
*/



}
