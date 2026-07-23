package handlers

import (
	"net/http"

	viewsconst "github.com/hilmanxcode/web-inventory-go/const"
	"github.com/hilmanxcode/web-inventory-go/entities"
	userServices "github.com/hilmanxcode/web-inventory-go/services"
	"github.com/hilmanxcode/web-inventory-go/sessions"
	"github.com/hilmanxcode/web-inventory-go/utils/formutil"
	"github.com/hilmanxcode/web-inventory-go/utils/jsonutil"
	"github.com/hilmanxcode/web-inventory-go/utils/viewsutil"
)

func LoginPage(w http.ResponseWriter, r *http.Request) {

	csrfToken := sessions.GetCSRFToken(w, r)

	data := map[string]any{
		"csrf_token": csrfToken,
	}

	successMsg, errorMsgs, _ := sessions.GetAndClearFlash(r)
	oldInput, _ := sessions.GetAndClearOldInput(r)

	if successMsg != "" {
		data["successMsg"] = successMsg
	}
	if len(errorMsgs) > 0 {
		data["errorMsgs"] = errorMsgs
	}
	if oldInput != nil {
		data["oldInput"] = oldInput
	}

	viewsutil.ShowView(viewsconst.VIEWS_LOGIN, data, w)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {

	var reqs = entities.UserLogin{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}

	invalid, message := formutil.Validate(reqs, r)

	if invalid {

		csrfToken := sessions.GetCSRFToken(w, r)

		jsonMessage := jsonutil.MapStringToJson(message, w)

		var data = map[string]any{
			"errors": string(jsonMessage),
			"oldInput": map[string]string{
				"email": reqs.Email,
			},
			"csrf_token": csrfToken,
		}

		viewsutil.ShowView(viewsconst.VIEWS_LOGIN, data, w)

		return
	}

	userServices.AuthUser(w, r, reqs)

}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {

	var reqs = entities.UserRegister{
		NamaLengkap: r.FormValue("nama_lengkap"),
		Email:       r.FormValue("email"),
		Password:    r.FormValue("password"),
		Role:        "Staff",
	}

	invalid, message := formutil.Validate(reqs, r)

	if invalid {

		csrfToken := sessions.GetCSRFToken(w, r)

		jsonMessage := jsonutil.MapStringToJson(message, w)

		var data = map[string]any{
			"errors":       string(jsonMessage),
			"registerPage": true,
			"oldInput": map[string]string{
				"nama_lengkap": reqs.NamaLengkap,
				"email":        reqs.Email,
			},
			"csrf_token": csrfToken,
		}

		viewsutil.ShowView(viewsconst.VIEWS_LOGIN, data, w)

		return
	}

	userServices.RegisterUser(w, r, reqs)

}
