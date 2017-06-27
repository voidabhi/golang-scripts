#!/usr/bin/env go

package main

import (
	"fmt"
	"time"
)

var (
	ch          chan bool
	timeoutSecs uint = 3
	workSecs    uint = 2
)

func main() {
	timeoutChan := make(chan bool, 1)
	ch = make(chan bool, 1)
	go SlowProcessingWork(timeoutChan)

	select {
	case <-ch:
		fmt.Println("The function worked")
	case <-timeoutChan:
		fmt.Println("Timeout!")
	}
}

func SlowProcessingWork(timeoutChan chan bool) {
	go timeout(timeoutChan)

	time.Sleep(time.Duration(workSecs) * time.Second)
	ch <- true
}

func timeout(timeoutChan chan bool) {
	time.Sleep(time.Duration(timeoutSecs) * time.Second)
	timeoutChan <- true
}
