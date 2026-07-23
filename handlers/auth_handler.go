package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	viewsconst "github.com/hilmanxcode/web-inventory-go/const"
	"github.com/hilmanxcode/web-inventory-go/database"
	"github.com/hilmanxcode/web-inventory-go/entities"
	userModel "github.com/hilmanxcode/web-inventory-go/models"
	"github.com/hilmanxcode/web-inventory-go/sessions"
	"github.com/hilmanxcode/web-inventory-go/utils/formutil"
	"github.com/hilmanxcode/web-inventory-go/utils/jsonutil"
	"github.com/hilmanxcode/web-inventory-go/utils/viewsutil"
	"golang.org/x/crypto/bcrypt"
)

func LoginPage(w http.ResponseWriter, r *http.Request) {

	c, err := r.Cookie("session_token")

	if err != nil {

		fmt.Println("MASUK SINI DLU")
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

		var data = map[string]any{
			"errorMsgs":  errorMsgs,
			"successMsg": successMsg,
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

	result, err := userModel.GetUserDataWithEmail(reqs.Email, userModel.UserColumn.Email, userModel.UserColumn.Password)

	c, errCookie := r.Cookie("session_token")
	// Kalau user nya gak ada / session invalid
	if err != nil {

		if errCookie != nil {

			sessions.SetSession(sessions.Session{
				ErrorMessages: []string{
					"Invalid Session",
				},
			}, w)

			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		val, err := sessions.GetSession(c.Value)

		if err != nil {
			sessions.SetSession(sessions.Session{
				ErrorMessages: []string{
					"Invalid Session",
				},
			}, w)

			// mengatur kode http itu sangat penting agar tidak terjadi error
			// kalau kita mau redirect kan lagi ke page itu sendiri kita pakai
			// http.statusfound
			// misal, kita lakuin action post ke /, lalu kita ingin ngeredirect kalau gagal
			// itu ke / lagi, tapi method get, kita pakai http.statusfound, jangan http.statusseeother
			// karna itu akan error, saya mengalaminya sendiri xD.
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		val.ErrorMessages = []string{
			"User tidak ditemukan",
		}

		sessions.UpdateSession(c.Value, val)

		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	fmt.Println(result.Email)

	err = bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(reqs.Password))

	if err != nil {
		val, err := sessions.GetSession(c.Value)

		if err != nil {
			sessions.SetSession(sessions.Session{
				ErrorMessages: []string{
					"Invalid Session",
				},
			}, w)

			// mengatur kode http itu sangat penting agar tidak terjadi error
			// kalau kita mau redirect kan lagi ke page itu sendiri kita pakai
			// http.statusfound
			// misal, kita lakuin action post ke /, lalu kita ingin ngeredirect kalau gagal
			// itu ke / lagi, tapi method get, kita pakai http.statusfound, jangan http.statusseeother
			// karna itu akan error, saya mengalaminya sendiri xD.
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		val.ErrorMessages = []string{
			"Password salah",
		}

		sessions.UpdateSession(c.Value, val)

		http.Redirect(w, r, "/", http.StatusFound)
		return

	} else {
		w.Write([]byte("Password benar"))
	}

	// c, err := r.Cookie("session_token")

	// if err != nil {
	// 	sessions.SetSession(sessions.Session{}, w)

	// 	viewsutil.ShowView(viewsconst.VIEWS_LOGIN, nil, w)
	// 	return
	// }

	// // val, err := sessions.GetSession(c.Value)

	// // val.Username = ""

}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {

	var reqs = entities.UserRegister{
		NamaLengkap: r.FormValue("nama_lengkap"),
		Email:       r.FormValue("email"),
		Password:    r.FormValue("password"),
		Role:        "Staff",
	}

	var errors map[string]string

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

	hash, err := bcrypt.GenerateFromPassword([]byte(reqs.Password), 10)

	if err != nil {
		log.Panic(err.Error())
	}

	successMsg, err, duplicate := database.InsertQuery(`
		INSERT INTO users (nama_lengkap, email, password, role)
		VALUES (?, ?, ?, ?)
	`, reqs.NamaLengkap, reqs.Email, hash, reqs.Role)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if duplicate {
		errors = map[string]string{
			"duplicate_email": "Email telah digunakan",
		}
	}

	errorMsgs := jsonutil.MapStringToJson(errors, w)

	var data = map[string]any{
		"success": successMsg,
		"errors":  string(errorMsgs),
	}

	viewsutil.ShowView(viewsconst.VIEWS_LOGIN, data, w)
}
