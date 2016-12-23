package main
 
import (
    "encoding/json"
    "fmt"
    "log"
    "time"
 
    "gopkg.in/redis.v2"
)
 
func main() {
    done := make(chan struct{})
    go redisv4Pub("one", done)
    go redisv4Sub(done)
    time.Sleep(time.Second * 1)
    close(done)
    println("exiting")
    time.Sleep(time.Second * 1)
}
 
type Notification struct {
    Id    int         `json:"Id"`
    Time  time.Time   `json:"Time"`
    Extra interface{} `json:"Extra"`
}
 
func redisv4Pub(index string, done chan struct{}) {
    redisClient := redis.NewTCPClient(&redis.Options{
        Addr:     fmt.Sprintf("%s:%d", "192.168.8.136", 6379),
        Password: "",
        DB:       0,
    })
    defer redisClient.Close()
    errcmd := redisClient.Ping()
    if errcmd.Err() != nil {
        panic("network error")
    }
 
    channel := "pubc1"
    count := 0
    for {
        count++
        data := &Notification{
            Id:   count,
            Time: time.Now(),
        }
 
        bs, _ := json.Marshal(data)
        cmd := redisClient.Publish(channel, string(bs))
        if cmd.Err() != nil {
            panic("pub error")
        }
        if count == 1 {
            fmt.Printf("string(bs): %+v\n", string(bs))
        }
        select {
        case <-done:
            println("exit pub")
            return
        default:
            fmt.Printf("select default count: %+v\n", count)
 
        }
    }
}
 
func redisv4Sub(done chan struct{}) {
    redisClient := redis.NewTCPClient(&redis.Options{
        Addr:     fmt.Sprintf("%s:%d", "192.168.8.136", 6379),
        Password: "",
        DB:       0,
    })
    errcmd := redisClient.Ping()
    if errcmd.Err() != nil {
        panic("network error")
    }
 
    channel := "pubc1"
    pubsub := redisClient.PubSub()
    defer pubsub.Close()
    pubsub.Subscribe(channel)
 
    var msg interface{}
    var err error
    data := &Notification{}
    for {
        msg, err = pubsub.ReceiveTimeout(time.Second * 2)
        // msg, err = pubsub.Receive()
        if err != nil {
            log.Printf("read timeout %d of redis sub", 2)
        }
 
        switch v := msg.(type) {
        case *redis.Message:
            _ = json.Unmarshal([]byte(v.Payload), data)
            log.Printf("data: %+v\n", data)
            fmt.Printf("time elapse: %+v\n", time.Since(data.Time))
 
        case *redis.Subscription:
            fmt.Printf("%s: %s %d\n", v.Channel, v.Kind, v.Count)
        case error:
            println("error")
            return
        }
        // time.Sleep(time.Second * 2)
 
        select {
        case <-done:
            println("exit sub")
            return
        default:
 
        }
    }
}
