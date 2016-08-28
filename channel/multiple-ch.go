
package main

import (
    "fmt"
    "time"
)

func waitForChannelsToClose(chans ...chan struct{}) {
    t := time.Now()
    for _, v := range chans {
        <-v
        fmt.Printf("%v for chan to close\n", time.Since(t))
    }
    fmt.Printf("%v for channels to close\n", time.Since(t))
}

func simulateProcessingThenClose(c chan struct{}) {
    time.Sleep(50 * time.Millisecond)
    close(c)
}

func main() {

    // Wait for a slice of channels
    var channelList []chan struct{}
    for i := 0; i < 10; i++ {
        channel := make(chan struct{})
        channelList = append(channelList, channel)
        go simulateProcessingThenClose(channel)

    }

    waitForChannelsToClose(channelList...)

    // Or use individual channels
    ch1 := make(chan struct{})
    ch2 := make(chan struct{})

    go simulateProcessingThenClose(ch1)
    go simulateProcessingThenClose(ch2)

    waitForChannelsToClose(ch1, ch2)

    fmt.Println("That's all for now!")
}
