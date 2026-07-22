package main

import (
	"log"
	"net/http"
	"os"

	"github.com/hilmanxcode/web-inventory-go/database"
	"github.com/hilmanxcode/web-inventory-go/routes"
	"github.com/joho/godotenv"
)

func StartNonTLSServer() {
	mux := new(http.ServeMux)
	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Redirecting to https://10.175.125.36/")
		http.Redirect(w, r, "https://10.175.125.36/", http.StatusTemporaryRedirect)
	}))

	http.ListenAndServe(":80", mux)
}

func main() {

	err := godotenv.Load()

	if err != nil {
		panic(err.Error())
	}

	dsn := os.Getenv("DATABASE_DSN")

	database.Connect(dsn)

	go StartNonTLSServer()

	router := routes.SetupRouter()

	log.Println("Server berjalan di port 443")

	// without http
	// err = http.ListenAndServe("0.0.0.0:8000", router)
	// with https
	err = http.ListenAndServeTLS("0.0.0.0:443", "server.crt", "server.key", router)

	if err != nil {
		log.Fatal(err.Error())
	}

}
