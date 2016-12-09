
package main

import (
	"fmt"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/urlshortener/v1"
	"net/http"
	"text/template"
)

func main() {
	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/shorten", handleShorten)
	http.HandleFunc("/lengthen", handleLengthen)

	http.ListenAndServe("localhost:8080", nil)
}

var rootHtmlTmpl = template.Must(template.New("rootHtml").Parse(`
<html><body>
<h1>My URL Shrtnr</h1>
{{if .}}{{.}}<br /><br />{{end}}
<form action="/shorten" type="POST">
Shorten this: <input type="text" name="longUrl" />
<input type="submit" value="Give me short URL" />
</form>
<br />
<form action="/lengthen" type="POST">
Expand this: http://goo.gl/<input type="text" name="shortUrl" />
<input type="submit" value="Give me long URL" />
</form>
</body></html>
`))

func handleRoot(w http.ResponseWriter, r *http.Request) {
	rootHtmlTmpl.Execute(w, nil)
}

func handleShorten(w http.ResponseWriter, r *http.Request) {
	longUrl := r.FormValue("longUrl")
	ctx := context.Background()
	client, err := google.DefaultClient(ctx, urlshortener.UrlshortenerScope)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
		return
	}

	urlshortenerSvc, _ := urlshortener.New(client
	url, _ := urlshortenerSvc.Url.Insert(&urlshortener.Url{LongUrl: longUrl}).Do()
	rootHtmlTmpl.Execute(w, fmt.Sprintf("Shortened version of `%s` is: %s", longUrl, url.Id))
}

func handleLengthen(w http.ResponseWriter, r *http.Request) {
	shortUrl := "http://goo.gl/" + r.FormValue("shortUrl")
	ctx := context.Background()
	client, err := google.DefaultClient(ctx, urlshortener.UrlshortenerScope)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
		return
	}

	urlshortenerSvc, _ := urlshortener.New(client)
	url, _ := urlshortenerSvc.Url.Get(shortUrl).Do()
	rootHtmlTmpl.Execute(w, fmt.Sprintf("Longer version of %s is: %s", shortUrl, url.LongUrl))
}
