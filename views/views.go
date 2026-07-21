package views

import (
	"html/template"

	viewsconst "github.com/hilmanxcode/web-inventory-go/const"
)

var Cache = map[string]*template.Template{
	viewsconst.VIEWS_LOGIN: generateTemplate(
		"views/auth/login.html",
	),
	viewsconst.VIEWS_DASHBOARD: generateTemplate(
		"views/dashboard/index.html",
	),
}

func generateTemplate(path ...string) *template.Template {
	return template.Must(template.ParseFiles(
		path...,
	))

}
