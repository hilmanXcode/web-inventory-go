package main

import (
	"fmt"

	"github.com/hilmanxcode/web-inventory-go/configs"
	"github.com/hilmanxcode/web-inventory-go/entities"
	"github.com/hilmanxcode/web-inventory-go/utils"
)

func main() {

	configs.Connect()

	reqs := entities.User{
		ID:          1337,
		NamaLengkap: "",
		Email:       "hilmanxcode@gmail.com",
		Password:    "TESTING",
		Role:        "Staff",
	}

	err, message := utils.Validate(reqs)

	if err {
		fmt.Println(message)
	}
	// router := routes.SetupRouter()

	// log.Println("Server berjalan di port 8000")

	// err := http.ListenAndServe("0.0.0.0:8000", router)

	// if err != nil {
	// 	log.Fatal(err.Error())
	// }

}
