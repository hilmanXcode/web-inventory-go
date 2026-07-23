package userServices

import (
	"log"
	"net/http"

	viewsconst "github.com/hilmanxcode/web-inventory-go/const"
	"github.com/hilmanxcode/web-inventory-go/database"
	"github.com/hilmanxcode/web-inventory-go/entities"
	userModel "github.com/hilmanxcode/web-inventory-go/models"
	"github.com/hilmanxcode/web-inventory-go/sessions"
	"github.com/hilmanxcode/web-inventory-go/utils/httputil"
	"github.com/hilmanxcode/web-inventory-go/utils/jsonutil"
	"github.com/hilmanxcode/web-inventory-go/utils/sessionutil"
	"github.com/hilmanxcode/web-inventory-go/utils/viewsutil"
	"golang.org/x/crypto/bcrypt"
)

func AuthUser(w http.ResponseWriter, r *http.Request, reqs entities.UserLogin) {
	result, err := userModel.GetUserDataWithEmail(reqs.Email, userModel.UserColumn.Email, userModel.UserColumn.Password)

	if err != nil {
		httputil.RedirectWithError(w, r, "User tidak ditemukan", "/")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(reqs.Password))

	if err != nil {
		val, err := sessionutil.GetSessionData(w, r, "/")

		if err != nil {
			httputil.RedirectWithError(w, r, "Invalid session", "/")
			return
		}

		val.OldInput = map[string]string{
			"email": reqs.Email,
		}

		sessions.UpdateSession(val.Key, val)

		httputil.RedirectWithError(w, r, "Password anda salah", "/")
		return
	}

	w.Write([]byte("Password benar"))
}

func RegisterUser(w http.ResponseWriter, reqs entities.UserRegister) {
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

	var data = map[string]any{
		"success": successMsg,
		"errors":  string(errorMsgs),
	}

	viewsutil.ShowView(viewsconst.VIEWS_LOGIN, data, w)
}
