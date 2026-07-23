package handlers

import (
	"fmt"
	"net/http"

	viewsconst "github.com/hilmanxcode/web-inventory-go/const"
	"github.com/hilmanxcode/web-inventory-go/entities"
	"github.com/hilmanxcode/web-inventory-go/sessions"
	"github.com/hilmanxcode/web-inventory-go/utils/formutil"
	"github.com/hilmanxcode/web-inventory-go/utils/httputil"
	"github.com/hilmanxcode/web-inventory-go/utils/jsonutil"
	"github.com/hilmanxcode/web-inventory-go/utils/viewsutil"
)

func CreateBarang(w http.ResponseWriter, r *http.Request) {

	if !sessions.VerifyCSRF(r) {
		fmt.Println("Masuk sini")
		httputil.RedirectWithError(w, r, "Invalid CSRF Token", "/")
		return
	}

	// fmt.Println(sessions.VerifyCSRF(r))

	// return

	var reqs = entities.BarangCreate{
		Sku:        r.FormValue("sku"),
		NamaBarang: r.FormValue("nama"),
		Kategori:   r.FormValue("kategori"),
		Stock:      r.FormValue("stock"),
	}

	// image := httputil.UploadBase64Handler(w, r, "image")

	invalid, message := formutil.Validate(reqs, r)

	if invalid {

		fmt.Println("MASUK INVALID")
		csrfToken := sessions.GetCSRFToken(w, r)

		jsonMessage := jsonutil.MapStringToJson(message, w)

		fmt.Println(string(jsonMessage))

		var data = map[string]any{
			"errors":     string(jsonMessage),
			"csrf_token": csrfToken,
		}

		viewsutil.ShowView(viewsconst.VIEWS_MASTER_BARANG, data, w)

		return
	}

	w.Write([]byte("LOLOS PENGECEKAN"))

}
