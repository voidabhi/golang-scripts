
/*

base.html
------------

<!doctype html>
<html>
  <head>
  {{template "extra_head" .}}
  </head>
<body>
{{template "nav" .}}

{{template "content" .}}

{{template "extra_footer" .}}
</body>
</html>

*/

/*

about.html


{{define "extra_head"}}{{end}}
{{define "nav"}}
          <li><a href="/">Home</a></li>
          <li class="active"><a href="#">About</a></li>
          <li><a href="/contact">Contact</a></li>
          <li><a href="/privacy">Privacy</a></li>
{{end}}
{{define "content"}}

      <div class="row-fluid marketing">
          <h4>About ...</h4>
          <p>Blah blah</p>
      </div>

{{end}}
{{define "extra_footer"}}{{end}}
{{template "base.html" .}}

*/


// compiles `templates/<templateName>` along with templates/common/* and the functions
// defined in main.funcMap
func compileTemplate(templateName string) *template.Template {
	t := template.New("")

	// compile common templates and include funcMap (see funcs.go)
	t = template.Must(t.Funcs(funcMap).
		ParseGlob("templates/common/*html"))

	// compile particular template along with common tempaltes and funcMap
	return template.Must(t.ParseFiles("templates/" + templateName))
}

// given a template name, returns an http.HandlerFunc that renders it
func renderTemplate(templateName string) http.HandlerFunc {
	templates := compileTemplate(templateName)
	context := getContext(templateName)
	return func(w http.ResponseWriter, r *http.Request) {
		if *recompile {
			templates = compileTemplate(templateName)
			context = getContext(templateName)
		}
		if err := templates.ExecuteTemplate(w, templateName, context); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
