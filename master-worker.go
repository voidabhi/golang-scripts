package main

import (
	"fmt"
	"time"
)

type Work struct {
	Fn        func()
	Completed bool
}

func NewWork(fn func()) *Work {
	return &Work{fn, false}
}

func Worker(in chan *Work, out chan *Work) {
	for {
		t := <-in
		t.Fn()
		t.Completed = true
		out <- t
	}
}

func MasterRun(fn func(), jobs int, concurrency int) {
	pending := make(chan *Work)
	done := make(chan *Work)

	go func() {
		for i := 0; i < jobs; i++ {
			pending <- NewWork(fn)
		}
	}()

	for i := 0; i < concurrency; i++ {
		go Worker(pending, done)
	}

	for i := 0; i < jobs; i++ {
		<-done
	}
}

func main() {

	printTime := func() {
		now := time.LocalTime()
		fmt.Println(now.String())
	}

	MasterRun(printTime, 100, 10)

}
