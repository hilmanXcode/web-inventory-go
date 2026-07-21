package viewsutil

import (
	"net/http"

	"github.com/hilmanxcode/web-inventory-go/views"
)

func ShowView(view string, data any, w http.ResponseWriter) {
	tmpl, ok := views.Cache[view]

	if !ok {
		http.Error(w, "Template tidak ditemukan", http.StatusInternalServerError)
		return
	}

	err := tmpl.Execute(w, data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
