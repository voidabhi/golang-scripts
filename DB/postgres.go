package main

import (
        "log"
        "time"
        "os"
        "fmt"
        "database/sql"
        _ "github.com/lib/pq"
)

func main() {
        fmt.Printf("StartTime: %v\n", time.Now())
        var (
                sStmt string = "insert into test (gopher_id, created) values ($1, $2)"
                gophers int = 10
                entries int = 10000
        )

        finishChan := make(chan int)

        for i := 0; i < gophers; i++ {
                go func(c chan int) {
                        db, err := sql.Open("postgres", "host=localhost dbname=testdb sslmode=disable")
                        if err != nil {
                                log.Fatal(err)
                        }
                        defer db.Close()

                        stmt, err := db.Prepare(sStmt)
                        if err != nil {
                                log.Fatal(err)
                        }
                        defer stmt.Close()

                        for j := 0; j < entries; j++ {
                                res, err := stmt.Exec(j, time.Now())
                                if err != nil || res == nil {
                                        log.Fatal(err)
                                }
                        }

                        c <- 1
                }(finishChan)
        }

        finishedGophers := 0
        finishLoop := false
        for {
                if finishLoop {
                        break
                }
                select {
                case n := <-finishChan:
                        finishedGophers += n
                        if finishedGophers == 10 {
                                finishLoop = true
                        }
                }
        }

        fmt.Printf("StopTime: %v\n", time.Now())
}
