package main

import (
	"fmt"
	"net/http"
	"net/http/cgi"
)

func main() {
	if err := cgi.Serve(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := w.Header()
		header.Set("Content-Type", "text/plain; charset=utf-8")
		fmt.Fprintln(w, "Method:", r.Method)
		fmt.Fprintln(w, "URL:", r.URL.String())
		query := r.URL.Query()
		for k := range query {
			fmt.Fprintln(w, "Query", k+":", query.Get(k))
		}
		r.ParseForm()
		form := r.Form
		for k := range form {
			fmt.Fprintln(w, "Form", k+":", form.Get(k))
		}
		post := r.PostForm
		for k := range post {
			fmt.Fprintln(w, "PostForm", k+":", post.Get(k))
		}
		fmt.Fprintln(w, "RemoteAddr:", r.RemoteAddr)
		if referer := r.Referer(); len(referer) > 0 {
			fmt.Fprintln(w, "Referer:", referer)
		}
		if ua := r.UserAgent(); len(ua) > 0 {
			fmt.Fprintln(w, "UserAgent:", ua)
		}
		for _, cookie := range r.Cookies() {
			fmt.Fprintln(w, "Cookie", cookie.Name+":", cookie.Value, cookie.Path, cookie.Domain, cookie.RawExpires)
		}
	})); err != nil {
		fmt.Println(err)
	}
}
