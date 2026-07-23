package usermodel

import (
	"fmt"
	"strings"

	"github.com/hilmanxcode/web-inventory-go/database"
)

type users struct {
	Id          string
	NamaLengkap string
	Email       string
	Password    string
	Role        string
}

var UserColumn = users{
	Id:          "id",
	NamaLengkap: "nama_lengkap",
	Email:       "email",
	Password:    "password",
	Role:        "role",
}

func GetUserDataWithEmail(email string, columns ...string) (users, error) {

	query := fmt.Sprintf(`
		SELECT %v FROM users WHERE email = ?
	`, strings.Join(columns, ","))

	var result = users{}

	err := database.GetSingleData(query, email).Scan(&result.Email, &result.Password)

	if err != nil {
		return users{}, err
	}

	return result, nil

}
