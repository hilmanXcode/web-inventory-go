package handlers

import (
	"fmt"
	"log"
	"net/http"

	constanta "github.com/hilmanxcode/web-inventory-go/const"
	"github.com/hilmanxcode/web-inventory-go/database"
	"github.com/hilmanxcode/web-inventory-go/entities"
	"github.com/hilmanxcode/web-inventory-go/sessions"
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

	myKey := sessions.SetSession(sessions.Session{
		OldInput: map[string]string{
			"nama_lengkap": reqs.NamaLengkap,
			"email":        reqs.Email,
			"password":     reqs.Password,
		},
	}, w)

	invalid, message := utils.Validate(reqs)

	if invalid {

		var data = map[string]any{
			"error":        message,
			"registerPage": true,
			"sessions":     sessions.SessionData[myKey],
		}

		fmt.Println(sessions.SessionData[myKey].OldInput)

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

	var data = map[string]any{
		"success": "Berhasil register akun",
	}

	utils.ShowView(constanta.VIEWS_LOGIN, data, w)
}
