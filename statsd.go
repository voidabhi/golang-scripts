package main

import (
    "log"
    "os"
    "time"

    "github.com/quipo/statsd"
)

func main() {
    // init
    prefix := "myproject."
    statsdclient := statsd.NewStatsdClient("localhost:8125", prefix)
    err := statsdclient.CreateSocket()
    if nil != err {
        log.Println(err)
        os.Exit(1)
    }
    interval := time.Second * 2 // aggregate stats and flush every 2 seconds
    stats := statsd.NewStatsdBuffer(interval, statsdclient)
    defer stats.Close()

    // not buffered: send immediately
    statsdclient.Incr("mymetric", 4)

    // buffered: aggregate in memory before flushing
    stats.Incr("mymetric", 1)
    stats.Incr("mymetric", 3)
    stats.Incr("mymetric", 1)
    stats.Incr("mymetric", 1)
}
