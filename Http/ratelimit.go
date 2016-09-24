package main

import (
	"fmt"
	"time"
)

func main() {
	ticker := rateLimit(4, 10)

	go work("A", ticker, 2e9)
	go work("B", ticker, 3e9)
	work("C", ticker, 4e9)
}

func work(name string, ratelimiter chan int, backoff time.Duration) {
	for {
		select {
		case _ = <- ratelimiter:
			fmt.Printf(name)
		default:
			time.Sleep(backoff)
		}
	}
}

func rateLimit(rps int, burst int) chan int {
	c := make(chan int, burst)
	for i := 0; i < burst; i++ {
		c <- 0
	}
	ticker := time.Tick(time.Second/time.Duration(rps))
	go func() {
		for {
			_ = <- ticker
			c <- 0
		}
	}()
	return c
}
