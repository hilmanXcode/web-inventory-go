package database

import (
	"database/sql"
	"log"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func Connect(dsn string) {

	var err error
	DB, err = sql.Open("mysql", dsn)

	if err != nil {
		log.Fatalf("Gagal menginisialisasi database: %v", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatalf("Gagal terhubung ke database: %v", err)
	}

	log.Println("Database berhasil terhubung")

}

func GetSingleData(query string, params ...any) *sql.Row {

	row := DB.QueryRow(query, params...)

	return row

}

func SqlQuery(query string, params ...any) (*sql.Rows, error) {

	rows, err := DB.Query(query, params...)

	if err != nil {
		return nil, err
	}

	return rows, nil

}

func InsertQuery(query string, params ...any) (string, error, bool) {

	_, err := DB.Exec(query, params...)

	if err != nil {
		// panic(err.Error())
		var error = string([]byte(err.Error()))
		var isDuplicate = strings.Contains(error, "Duplicate")

		if isDuplicate {
			return "", nil, true
		}

		return "", err, false

	}

	return "Berhasil menambahkan user", nil, false

}
