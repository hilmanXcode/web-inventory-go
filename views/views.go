package views

import "html/template"

var Cache = map[string]*template.Template{
	"login": generateTemplate("views/auth/login.html"),
}

func generateTemplate(path string) *template.Template {
	return template.Must(template.ParseFiles(
		path,
	))
}
