package configs

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func Connect() {

	var dsn = "manz:supersecretpassword@tcp(127.0.0.1:3306)/db_inventory"

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

func SqlQuery(query string, params ...any) (*sql.Rows, error) {

	// if params != nil
	rows, err := DB.Query(query, params...)

	if err != nil {
		return nil, err
	}

	return rows, nil

}
