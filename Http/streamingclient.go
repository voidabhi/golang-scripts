package main
 
import (
    "fmt"
    "io"
    "net"
    "time"
)
 
const (
    addr                     = "localhost:3000"
    maxConcurrentConnections = 100
)
 
func sendRequest() {
    con, err := net.Dial("tcp", addr)
    if err != nil {
        fmt.Printf("Could not connect to %q: %v.", addr, err)
    }
    // this is needed for HTTP Servers to respond at all
    fmt.Fprintf(con, "GET / HTTP/1.0\r\n\r\n")
 
    buffer := make([]byte, 10)
    err = nil
    for {
        read, err := con.Read(buffer)
        if read > 0 {
            fmt.Printf(string(buffer[:read]))
        }
 
        if err != nil {
            if err == io.EOF {
                break
            }
 
            fmt.Printf("Error reading: %v.", err)
        }
    }
}
func sendRequests(abort <-chan time.Time) {
    requestDone := make(chan bool)
    currentRequests := 0
 
    for {
        select {
        case <-abort:
            return
        case <-requestDone:
            currentRequests--
        default:
            if currentRequests < maxConcurrentConnections {
                currentRequests++
                go func(requestDone chan bool) {
                    sendRequest()
                    requestDone <- true
                }(requestDone)
            }
        }
    }
}
 
func main() {
    duration := time.After(10 * time.Second)
 
    go sendRequests(duration)
    <-duration
}
