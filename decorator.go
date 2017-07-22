package main

import (
	"fmt"
)

type displayFunc func(s string)

func decorate(f displayFunc) displayFunc {
	return func(s string) {
		fmt.Println("Before")
		f(s)
		fmt.Println("After")
	}
}

func display(s string) {
	fmt.Println(s)
}

func main() {
	display := decorate(display)
	display("In the middle")
}
