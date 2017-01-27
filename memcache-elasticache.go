package main

import (
    "bufio"
    "fmt"
    "io"
    "net"
    "os"
    "strconv"
    "strings"

    "github.com/bradfitz/gomemcache/memcache"
)

// Node is a single ElastiCache node
type Node struct {
    URL  string
    Host string
    IP   string
    Port int
}

func main() {
    urls, err := clusterNodes()
    if err != nil {
        fmt.Println(err.Error())
    }

    mc := memcache.New(urls...)
    mc.Set(&memcache.Item{Key: "foo", Value: []byte("my value")})

    it, err := mc.Get("foo")
    if err != nil {
        fmt.Println(err.Error())
        return
    }

    fmt.Printf("%+v", it)
}

func clusterNodes() ([]string, error) {
    conn, err := net.Dial("tcp", elasticache())
    if err != nil {
        return nil, err
    }
    defer conn.Close()

    command := "config get cluster\r\n"
    fmt.Fprintf(conn, command)

    response, err := parseNodes(conn)
    if err != nil {
        return nil, err
    }

    urls, err := parseURLs(response)
    if err != nil {
        return nil, err
    }

    return urls, nil
}

func elasticache() string {
    var endpoint string

    endpoint = os.Getenv("ELASTICACHE_ENDPOINT")
    if len(endpoint) == 0 {
        endpoint = "127.0.0.1:11212"
    }

    return endpoint
}

func parseNodes(conn io.Reader) (string, error) {
    var response string

    count := 0
    location := 3 // AWS docs suggest that nodes will always be listed on line 3

    scanner := bufio.NewScanner(conn)
    for scanner.Scan() {
        count++
        if count == location {
            response = scanner.Text()
        }
        if scanner.Text() == "END" {
            break
        }
    }

    if err := scanner.Err(); err != nil {
        return "", err
    }

    return response, nil
}

func parseURLs(response string) ([]string, error) {
    var urls []string
    var nodes []Node

    items := strings.Split(response, " ")

    for _, v := range items {
        fields := strings.Split(v, "|") // ["host", "ip", "port"]

        port, err := strconv.Atoi(fields[2])
        if err != nil {
            return nil, err
        }

        node := Node{fmt.Sprintf("%s:%d", fields[1], port), fields[0], fields[1], port}
        nodes = append(nodes, node)
        urls = append(urls, node.URL)

        fmt.Printf("Host: %s\n", node.Host)
        fmt.Printf("IP: %s\n", node.IP)
        fmt.Printf("Port: %d\n", node.Port)
        fmt.Printf("URL: %s\n\n", node.URL)
    }

    return urls, nil
}
