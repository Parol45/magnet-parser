package main

import (
	"container/list"
	"fmt"
)

func Foo() error {
	var err error = nil
	// …
	return err
}

type CustomStringType string

const (
	q CustomStringType = "1"
	w CustomStringType = "2"
	e CustomStringType = "3"
	r CustomStringType = "4"
)

func main() {
	l := list.List{}
	l.Init()

	ch := make(chan int)
	close(ch)

	var arr = []int{1}
	fmt.Printf("%v\n", arr[1:])

	var asdfg any
	asdfg = "qwerty"
	res, b := asdfg.(int64)
	fmt.Printf("%v %v\n", res, b)

	var cst CustomStringType
	cst = ""
	fmt.Printf("%v \n", cst)

	m := map[string]string{}
	println(m["1234"])

	t := []int{1, 2, 3, 4, 5}
	sl := t[1:3]
	t[2] = 10000
	for _, item := range sl {
		println(item)
	}
}
