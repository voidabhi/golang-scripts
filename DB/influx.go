package main

import (
    "log"
    "time"

    "github.com/influxdata/influxdb/client/v2"
)

const (
    MyDB = "square_holes"
    username = "bubba"
    password = "bumblebeetuna"
)

func main() {
    // Make client
    c, err := client.NewHTTPClient(client.HTTPConfig{
        Addr: "http://localhost:8086",
        Username: username,
        Password: password,
    })

    if err != nil {
        log.Fatalln("Error: ", err)
    }

    // Create a new point batch
    bp, err := client.NewBatchPoints(client.BatchPointsConfig{
        Database:  MyDB,
        Precision: "s",
    })

    if err != nil {
        log.Fatalln("Error: ", err)
    }

    // Create a point and add to batch
    tags := map[string]string{"cpu": "cpu-total"}
    fields := map[string]interface{}{
        "idle":   10.1,
        "system": 53.3,
        "user":   46.6,
    }
    pt, err := client.NewPoint("cpu_usage", tags, fields, time.Now())

    if err != nil {
        log.Fatalln("Error: ", err)
    }

    bp.AddPoint(pt)

    // Write the batch
    c.Write(bp)
}
