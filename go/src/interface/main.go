package main

import (
	"fmt"
)

type CounterI interface {
	Inc() uint64
	Dec() uint64
}

type GenCounterS struct {
	cnt uint64
}

// 实现接口的方法尽量使用指针实例的方法，这样是地址传递
// 也可以使用对象实例的方法，这样是值传递
func (c *GenCounterS) Inc() uint64 {
	c.cnt++
	return c.cnt
}

func (c *GenCounterS) Dec() uint64 {
	c.cnt--
	return c.cnt
}

func (c *GenCounterS) Get() uint64 {
	return c.cnt
}

/*
func (c GenCounterS) Inc() uint64 {
	c.cnt++
	return c.cnt
}

func (c GenCounterS) Dec() uint64 {
	c.cnt--
	return c.cnt
}
*/

type ReaderI interface {
	Get() uint64
}

type SamplerI interface {
	CounterI
	ReaderI
}

func main() {
	var ct GenCounterS = GenCounterS{0}
	var sample SamplerI = &ct

	fmt.Printf("Counter: %v\n", ct.Inc())	
	fmt.Printf("Inc: %v\n", Inc(&ct))
	fmt.Printf("Dec: %v\n", Dec(&ct))
	fmt.Printf("sample: %v\n", ct.Get())
	fmt.Printf("sample: %v\n", get(&ct))
	fmt.Printf("sample: %v\n", get(sample))
	fmt.Printf("Assert: %t\n", assert(ct))

}

func assert(i interface{}) bool {
	_, ok := i.(GenCounterS)

	return ok
}

func Inc(c CounterI) uint64 {
	return c.Inc()
}

func Dec(c CounterI) uint64 {
	return c.Dec()
}

func get(s SamplerI) uint64 {
	return s.Get()
}
