package handlers

import (
	"fmt"
	"log"
	"net/http"

	constanta "github.com/hilmanxcode/web-inventory-go/const"
	"github.com/hilmanxcode/web-inventory-go/database"
	"github.com/hilmanxcode/web-inventory-go/entities"
	"github.com/hilmanxcode/web-inventory-go/utils"
	"golang.org/x/crypto/bcrypt"
)

func LoginPage(w http.ResponseWriter, r *http.Request) {
	utils.ShowView(constanta.VIEWS_LOGIN, nil, w)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var email = r.FormValue("email")
	var password = r.FormValue("password")

	var result = fmt.Sprintf("Email: %s\nPassword: %s", email, password)

	w.Write([]byte(result))
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {

	var reqs = entities.User{
		NamaLengkap: r.FormValue("nama_lengkap"),
		Email:       r.FormValue("email"),
		Password:    r.FormValue("password"),
		Role:        "Staff",
	}

	var errors map[string]string

	invalid, message := utils.Validate(reqs)

	if invalid {

		jsonMessage := utils.MapStringToJson(message, w)

		var data = map[string]any{
			"errors":       string(jsonMessage),
			"registerPage": true,
			"oldInput": map[string]string{
				"nama_lengkap": reqs.NamaLengkap,
				"email":        reqs.Email,
			},
		}

		utils.ShowView(constanta.VIEWS_LOGIN, data, w)

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

	errorMsgs := utils.MapStringToJson(errors, w)

	var data = map[string]any{
		"success": successMsg,
		"errors":  string(errorMsgs),
	}

	utils.ShowView(constanta.VIEWS_LOGIN, data, w)
}
