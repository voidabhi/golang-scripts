
package main

import (
	"fmt"
	"runtime"
	"time"
)

var (
	INTERVAL_SEC = 10
)

func PrintRoutine1(intervalInSec int) {
	t := time.NewTicker(time.Duration(intervalInSec) * time.Second)
	for _ = range t.C {
		fmt.Println("PrintRoutine1")
	}
}

func PrintRoutine2(intervalInSec int) {
	t := time.NewTicker(time.Duration(intervalInSec) * time.Second)
	for _ = range t.C {
		fmt.Println("PrintRoutine2")
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	go PrintRoutine1(INTERVAL_SEC)
	go PrintRoutine2(INTERVAL_SEC)

	// block forever so that your program won't end
	select {}
}
