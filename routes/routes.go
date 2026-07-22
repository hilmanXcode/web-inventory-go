package routes

import (
	"net/http"

	"github.com/hilmanxcode/web-inventory-go/handlers"
	"github.com/hilmanxcode/web-inventory-go/middleware"
)

func SetupRouter() *http.ServeMux {

	mux := http.NewServeMux()

	// For handling an asset file
	mux.Handle("GET /assets/",
		http.StripPrefix(
			"/assets/",
			http.FileServer(http.Dir("static/assets")),
		),
	)

	/* Auth Route */
	// Login Route
	mux.HandleFunc("GET /", middleware.GuestOnly(handlers.LoginPage))
	mux.HandleFunc("POST /", handlers.LoginHandler)

	// Register Route
	mux.HandleFunc("POST /register", handlers.RegisterHandler)

	// Dashboard Route
	mux.HandleFunc("GET /dashboard", middleware.RequireAuth(handlers.DashboardPage))

	return mux

}
