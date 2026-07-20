package handlers

import (
	"fmt"
	"log"
	"net/http"

	constanta "github.com/hilmanxcode/web-inventory-go/const"
	"github.com/hilmanxcode/web-inventory-go/database"
	"github.com/hilmanxcode/web-inventory-go/entities"
	"github.com/hilmanxcode/web-inventory-go/utils"
	"github.com/hilmanxcode/web-inventory-go/views"
	"golang.org/x/crypto/bcrypt"
)

func LoginPage(w http.ResponseWriter, r *http.Request) {

	tmpl, ok := views.Cache["login"]

	if !ok {
		http.Error(w, "Template tidak ditemukan", http.StatusInternalServerError)
		return
	}

	err := tmpl.Execute(w, nil)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
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

	invalid, message := utils.Validate(reqs)

	if invalid {

		var data = map[string]any{
			"message":      message,
			"registerPage": true,
		}

		utils.ShowView(constanta.VIEWS_LOGIN, data, w)

		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(reqs.Password), 10)

	if err != nil {
		log.Panic(err.Error())
	}

	database.InsertQuery(`
		INSERT INTO users (nama_lengkap, email, password, role)
		VALUES (?, ?, ?, ?)
	`, reqs.NamaLengkap, reqs.Email, hash, reqs.Role)

	var result = fmt.Sprintf("Nama Lengkap: %s\nEmail: %s\nPassword: %s", reqs.NamaLengkap, reqs.Email, hash)

	w.Write([]byte(result))
}
