package main
import (
        "runtime"
        "github.com/fluent/fluent-logger-golang/fluent"
)
func main() {
        logger, _ := fluent.New(fluent.Config{})
        tag2 := "myapp.errors"
        defer func() {
                if e := recover(); e != nil {
                        stack := make([]byte, 1<<16)
                        sz := runtime.Stack(stack, true)
                        stackstr := string(stack[:sz])
                        var data3 = map[string]string {
                                "app_id":       "command-underline",
                                "version":      "v2.0.0",
                                "message":      "A slightly modified error occurred",
                                "error_message":"This is my error other messages",
                                "function":     "main",
                                "exception":    stackstr,
                                "environment":  "prod",
                        }
                        _ = logger.Post(tag2, data3)
                }
        }()
        panic("something went horribly wrong")
}
