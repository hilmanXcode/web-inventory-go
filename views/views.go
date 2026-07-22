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
		"views/dashboard/template.html",
		"views/dashboard/index.html",
	),

	viewsconst.VIEWS_MASTER_BARANG: generateTemplate(
		"views/dashboard/template.html",
		"views/dashboard/manajemen_barang/index.html",
	),

	viewsconst.VIEWS_BARANG_MASUK: generateTemplate(
		"views/dashboard/template.html",
		"views/dashboard/barang_masuk/index.html",
	),

	viewsconst.VIEWS_BARANG_KELUAR: generateTemplate(
		"views/dashboard/template.html",
		"views/dashboard/barang_keluar/index.html",
	),

	viewsconst.VIEWS_LAPORAN_STOK: generateTemplate(
		"views/dashboard/template.html",
		"views/dashboard/laporan_stok/index.html",
	),
}

func generateTemplate(path ...string) *template.Template {
	return template.Must(template.ParseFiles(
		path...,
	))

}
