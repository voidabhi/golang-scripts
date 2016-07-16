
// crud mysql operations

package main

import (
    "os"
    "fmt"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
)

func main() {
    user := "root"
    password := ""
    database := "test_go"
    var n int
    var s string

    con, err := sql.Open("mysql", user+":"+password+"@/"+database)
    _, err = con.Exec("create table t(n int, s varchar(256))")
    if err != nil {
        fmt.Println("Can't create table 't'")
        os.Exit(1)
    }
    _, err = con.Exec("insert into t (n, s) values (1, 'foo')")
    _, err = con.Exec("insert into t (n, s) values (?, ?)", 2, "bar")

    fmt.Println("Select one row")
    row := con.QueryRow("select n, s from t where n=1")
    err = row.Scan(&n, &s)
    fmt.Println(n, s)  // 1, foo
   
    row = con.QueryRow("select n, s from t where n=?", 2)
    err = row.Scan(&n, &s)
    fmt.Println(n, s)  // 2, bar

    fmt.Println("Select all rows")
    rows, err := con.Query("select n, s from t")
    for rows.Next() {
        rows.Scan(&n, &s)
        fmt.Println(n, s)
    }

    defer con.Close()
}package main

import (
    "os"
    "fmt"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
)

func main() {
    user := "root"
    password := ""
    database := "test_go"
    var n int
    var s string

    con, err := sql.Open("mysql", user+":"+password+"@/"+database)
    _, err = con.Exec("create table t(n int, s varchar(256))")
    if err != nil {
        fmt.Println("Can't create table 't'")
        os.Exit(1)
    }
    _, err = con.Exec("insert into t (n, s) values (1, 'foo')")
    _, err = con.Exec("insert into t (n, s) values (?, ?)", 2, "bar")

    fmt.Println("Select one row")
    row := con.QueryRow("select n, s from t where n=1")
    err = row.Scan(&n, &s)
    fmt.Println(n, s)  // 1, foo
   
    row = con.QueryRow("select n, s from t where n=?", 2)
    err = row.Scan(&n, &s)
    fmt.Println(n, s)  // 2, bar

    fmt.Println("Select all rows")
    rows, err := con.Query("select n, s from t")
    for rows.Next() {
        rows.Scan(&n, &s)
        fmt.Println(n, s)
    }

    defer con.Close()
}
