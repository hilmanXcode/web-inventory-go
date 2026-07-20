package main

import (
	"log"
	"net/http"

	"github.com/hilmanxcode/web-inventory-go/database"
	"github.com/hilmanxcode/web-inventory-go/routes"
)

func main() {

	database.Connect()

	// reqs := entities.User{
	// 	ID:          1337,
	// 	NamaLengkap: "",
	// 	Email:       "hilmanxcode@gmail.com",
	// 	Password:    "TESTING",
	// 	Role:        "Staff",
	// }

	// err, message := utils.Validate(reqs)

	// if err {
	// 	fmt.Println(message)
	// }
	router := routes.SetupRouter()

	log.Println("Server berjalan di port 8000")

	err := http.ListenAndServe("0.0.0.0:8000", router)

	if err != nil {
		log.Fatal(err.Error())
	}

}
