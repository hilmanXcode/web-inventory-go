package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/hilmanxcode/web-inventory-go/routes"
)

func main() {

	router := routes.SetupRouter()

	fmt.Println("Server berjalan di port 8000")

	err := http.ListenAndServe("0.0.0.0:8000", router)

	if err != nil {
		log.Fatal(err.Error())
	}

}
