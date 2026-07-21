package main

import (
	"log"
	"net/http"
	"os"

	"github.com/hilmanxcode/web-inventory-go/database"
	"github.com/hilmanxcode/web-inventory-go/routes"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()

	if err != nil {
		panic(err.Error())
	}

	dsn := os.Getenv("DATABASE_DSN")

	database.Connect(dsn)

	router := routes.SetupRouter()

	log.Println("Server berjalan di port 8000")

	err = http.ListenAndServe("0.0.0.0:8000", router)

	if err != nil {
		log.Fatal(err.Error())
	}

}
