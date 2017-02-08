package main

import (
    "log"
    "net/http"
    "net/http/httputil"
)

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        director := func(req *http.Request) {
            req = r
            req.URL.Scheme = "http"
            req.URL.Host = r.Host
        }
        proxy := &httputil.ReverseProxy{Director: director}
        proxy.ServeHTTP(w, r)
    })
    log.Fatal(http.ListenAndServe(":8181", nil))
}
