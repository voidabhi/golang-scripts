package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	wait := make(chan bool)
	ctx, cancel := context.WithTimeout(context.Background(), 3200*time.Millisecond)
	defer cancel()
	go func() {

		select {
		case <-ctx.Done():
			fmt.Println("Timeout:", string(time.Now().Format("05.00")))
			err := ctx.Err()
			if err != nil {
				fmt.Println("in <-ctx.Done(): ", err)
			}
			// main app is blocking, waiting to hear about how this went
			wait<-true
			break		
		}

	}()
	time.Sleep(2 * time.Second)
	fmt.Println("first sleep completed, ",string(time.Now().Format("05.00")))
	//cancel()
	time.Sleep(2 * time.Second)
	fmt.Println("after second sleep done, ", string(time.Now().Format("05.00")))
	<-wait	// for the goroutine we started earlier
}
