package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/hilmanxcode/web-inventory-go/entities"
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
	var user entities.User

	hash, err := bcrypt.GenerateFromPassword([]byte(r.FormValue("password")), 10)

	if err != nil {
		log.Panic(err.Error())
	}

	user.NamaLengkap = r.FormValue("nama_lengkap")
	user.Email = r.FormValue("email")
	user.Password = string(hash)
	user.Role = "Staff"

	var result = fmt.Sprintf("Nama Lengkap: %s\nEmail: %s\nPassword: %s", user.NamaLengkap, user.Email, user.Password)

	w.Write([]byte(result))
}
