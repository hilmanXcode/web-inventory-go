package userServices

import (
	"log"
	"net/http"

	viewsconst "github.com/hilmanxcode/web-inventory-go/const"
	"github.com/hilmanxcode/web-inventory-go/database"
	"github.com/hilmanxcode/web-inventory-go/entities"
	usermodel "github.com/hilmanxcode/web-inventory-go/models/userModel"
	"github.com/hilmanxcode/web-inventory-go/sessions"
	"github.com/hilmanxcode/web-inventory-go/utils/httputil"
	"github.com/hilmanxcode/web-inventory-go/utils/jsonutil"
	"github.com/hilmanxcode/web-inventory-go/utils/sessionutil"
	"github.com/hilmanxcode/web-inventory-go/utils/viewsutil"
	"golang.org/x/crypto/bcrypt"
)

func AuthUser(w http.ResponseWriter, r *http.Request, reqs entities.UserLogin) {

	if !sessions.VerifyCSRF(r) {
		httputil.RedirectWithError(w, r, "Invalid CSRF Token", "/")
		return
	}

	result, err := usermodel.GetUserDataWithEmail(reqs.Email, usermodel.UserColumn.Email, usermodel.UserColumn.Password)

	if err != nil {
		httputil.RedirectWithError(w, r, "Email atau Password anda salah!", "/")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(reqs.Password))

	val, errSession := sessionutil.GetSessionData(w, r, "/")

	if err != nil {

		if errSession != nil {
			httputil.RedirectWithError(w, r, "Invalid session", "/")
			return
		}

		val.OldInput = map[string]string{
			"email": reqs.Email,
		}

		sessions.UpdateSession(val.Key, val)

		httputil.RedirectWithError(w, r, "Email atau Password anda salah!", "/")
		return
	}

	val.Email = result.Email

	sessions.UpdateSession(val.Key, val)

	http.Redirect(w, r, "/dashboard", http.StatusFound)

}

func RegisterUser(w http.ResponseWriter, r *http.Request, reqs entities.UserRegister) {
	var errors map[string]string

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

	csrfToken := sessions.GetCSRFToken(w, r)

	var data = map[string]any{
		"success":    successMsg,
		"errors":     string(errorMsgs),
		"csrf_token": csrfToken,
		"oldInput": map[string]string{
			"nama_lengkap": reqs.NamaLengkap,
		},
	}

	viewsutil.ShowView(viewsconst.VIEWS_LOGIN, data, w)
}
