package handlers

import (
	"log"
	"net/http"

	"github.com/google/uuid"
	viewsconst "github.com/hilmanxcode/web-inventory-go/const"
	"github.com/hilmanxcode/web-inventory-go/entities"
	userServices "github.com/hilmanxcode/web-inventory-go/services"
	"github.com/hilmanxcode/web-inventory-go/sessions"
	"github.com/hilmanxcode/web-inventory-go/utils/formutil"
	"github.com/hilmanxcode/web-inventory-go/utils/jsonutil"
	"github.com/hilmanxcode/web-inventory-go/utils/viewsutil"
)

func LoginPage(w http.ResponseWriter, r *http.Request) {

	c, err := r.Cookie("session_token")

	if err != nil {
		sessions.SetSession(sessions.Session{
			Key: uuid.NewString(),
		}, w)
		viewsutil.ShowView(viewsconst.VIEWS_LOGIN, nil, w)
		return
	} else {

		successMsg, errorMsgs, err := sessions.GetAndClearFlash(r)

		_, err = sessions.GetSession(c.Value)

		if err != nil {

			sessions.SetSession(sessions.Session{
				Key: uuid.NewString(),
			}, w)

		}

		if err != nil {
			sessions.ClearSession(c.Value, w)
			sessions.SetSession(sessions.Session{
				Key: uuid.NewString(),
			}, w)
			viewsutil.ShowView(viewsconst.VIEWS_LOGIN, nil, w)
			return
		}

		oldInput, err := sessions.GetAndClearOldInput(r)

		if err != nil {
			log.Fatal("harusnya sudah ada session dari sini")
		}

		var data = map[string]any{
			"errorMsgs":  errorMsgs,
			"successMsg": successMsg,
			"oldInput":   oldInput,
		}

		viewsutil.ShowView(viewsconst.VIEWS_LOGIN, data, w)

	}

}

func LoginHandler(w http.ResponseWriter, r *http.Request) {

	var reqs = entities.UserLogin{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}

	invalid, message := formutil.Validate(reqs, r)

	if invalid {

		jsonMessage := jsonutil.MapStringToJson(message, w)

		var data = map[string]any{
			"errors": string(jsonMessage),
			"oldInput": map[string]string{
				"email": reqs.Email,
			},
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

		jsonMessage := jsonutil.MapStringToJson(message, w)

		var data = map[string]any{
			"errors":       string(jsonMessage),
			"registerPage": true,
			"oldInput": map[string]string{
				"nama_lengkap": reqs.NamaLengkap,
				"email":        reqs.Email,
			},
		}

		viewsutil.ShowView(viewsconst.VIEWS_LOGIN, data, w)

		return
	}

	userServices.RegisterUser(w, reqs)

}
