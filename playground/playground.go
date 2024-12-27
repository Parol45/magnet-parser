package main

import (
	"fmt"
)

func Foo() error {
	var err error = nil
	// …
	return err
}

func main() {
	var err error
	err = Foo()
	fmt.Println(err)        // <nil>
	fmt.Println(err == nil) // false
}