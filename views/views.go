package views

import (
	"html/template"

	constanta "github.com/hilmanxcode/web-inventory-go/const"
)

var Cache = map[string]*template.Template{
	constanta.VIEWS_LOGIN: generateTemplate("views/auth/login.html"),
}

func generateTemplate(path string) *template.Template {
	return template.Must(template.ParseFiles(
		path,
	))

}
