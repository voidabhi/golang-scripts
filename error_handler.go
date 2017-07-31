package main

import (
	"errors"
	"fmt"
)

var (
	ErrorX = errors.New("error X")
	ErrorY = errors.New("error Y")
)

func GetErrorX() error {
	return ErrorX
}

func main() {
	err := GetErrorX()
	// With switch
	switch err {
	case ErrorX:
		fmt.Println("It's error X")
	case ErrorY:
		fmt.Println("It's error Y")
	default:
		fmt.Printf("Error, %s", err)
	}

	// With if
	if err == ErrorX {
		fmt.Println("It's error X")
	} else if err == ErrorY {
		fmt.Println("It's error Y")
	} else {
		fmt.Printf("Error, %s", err)
	}
}
