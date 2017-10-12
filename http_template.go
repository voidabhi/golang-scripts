package main

import (
	"html/template"
	"net/http"
)

var (
	templates *template.Template
)

func main() {
	// This is the only way I have found to be able to serve files requested in the templates
	http.Handle("/static/img/", http.StripPrefix("/static/img/",
		http.FileServer(http.Dir(path.Join(rootdir, "/static/img/")))))

	http.Handle("/static/css/", http.StripPrefix("/static/css/",
		http.FileServer(http.Dir(path.Join(rootdir, "/static/css/")))))
	
	http.HandleFunc("/", index)
	templates = template.Must(template.ParseFiles("index.html"))
	http.ListenAndServe(":8080", nil)
}

func index(w http.ResponseWriter, r *http.Request) {

	err := templates.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
