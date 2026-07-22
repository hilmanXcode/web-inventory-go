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
	mux.HandleFunc("POST /", middleware.GuestOnly(handlers.LoginHandler))

	// Register Route
	mux.HandleFunc("POST /register", handlers.RegisterHandler)

	// Dashboard Route
	mux.HandleFunc("GET /dashboard", middleware.RequireAuth(handlers.DashboardPage))
	mux.HandleFunc("GET /dashboard/master_barang", middleware.RequireAuth(handlers.MasterBarang))
	mux.HandleFunc("GET /dashboard/barang_masuk", middleware.RequireAuth(handlers.BarangMasuk))
	mux.HandleFunc("GET /dashboard/barang_keluar", middleware.RequireAuth(handlers.BarangKeluar))
	mux.HandleFunc("GET /dashboard/laporan_stok", middleware.RequireAuth(handlers.LaporanStok))

	return mux

}
