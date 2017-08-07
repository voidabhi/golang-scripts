package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

type MyFace interface {
	A()
}

type Cat struct{}
type Dog struct{}

func (c *Cat) A() {
	fmt.Println("Meow")
}

func (d *Dog) A() {
	fmt.Println("Woof")
}

func init() {
	// This type must match exactly what youre going to be using,
	// down to whether or not its a pointer
	gob.Register(&Cat{})
	gob.Register(&Dog{})
}

func main() {
	network := new(bytes.Buffer)
	enc := gob.NewEncoder(network)

	var inter MyFace
	inter = new(Cat)

	// Note: pointer to the interface
	err := enc.Encode(&inter)
	if err != nil {
		panic(err)
	}

	inter = new(Dog)
	err = enc.Encode(&inter)
	if err != nil {
		panic(err)
	}

	// Now lets get them back out
	dec := gob.NewDecoder(network)

	var get MyFace
	err = dec.Decode(&get)
	if err != nil {
		panic(err)
	}

	// Should meow
	get.A()

	err = dec.Decode(&get)
	if err != nil {
		panic(err)
	}

	// Should woof
	get.A()

}
