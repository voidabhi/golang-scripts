
package main

import (
                "fmt"
                "sync"
                "runtime"
)

func worker(id int, c chan int, wg *sync.WaitGroup) {
                wg.Add(1)
                defer wg.Done()
                for v := range c {
                               fmt.Printf("id %v: %v\n", id, v)
                }
}


func main() {
                c := make(chan int)
                var wg sync.WaitGroup
                for wid:=0; wid<runtime.NumCPU()*2; wid++ {
                               go worker(wid, c, &wg)
                }
                
                for i:=0; i<15; i++ {
                               c <- i
                }
                close(c)
                wg.Wait()
}
