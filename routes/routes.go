package routes

import (
	"net/http"

	"github.com/hilmanxcode/web-inventory-go/handlers"
	"github.com/hilmanxcode/web-inventory-go/middleware"
	"github.com/hilmanxcode/web-inventory-go/sessions"
	"github.com/hilmanxcode/web-inventory-go/utils/httputil"
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
	mux.HandleFunc("GET /{$}", middleware.GuestOnly(handlers.LoginPage))
	mux.HandleFunc("POST /", middleware.GuestOnly(handlers.LoginHandler))
	mux.HandleFunc("POST /logout", func(w http.ResponseWriter, r *http.Request) {

		c, err := r.Cookie("session_token")

		if err != nil {
			httputil.RedirectWithError(w, r, "Invalid token", "/")
		}

		if sessions.VerifyCSRF(r) {
			sessions.ClearSession(c.Value, w)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		httputil.RedirectWithError(w, r, "Invalid csrf token", "/")

	})

	// Register Route
	mux.HandleFunc("POST /register", handlers.RegisterHandler)

	// Dashboard Route
	mux.HandleFunc("GET /dashboard", middleware.RequireAuth(handlers.DashboardPage))
	mux.HandleFunc("GET /dashboard/master_barang", middleware.RequireAuth(handlers.MasterBarang))
	mux.HandleFunc("GET /dashboard/barang_masuk", middleware.RequireAuth(handlers.BarangMasuk))
	mux.HandleFunc("GET /dashboard/barang_keluar", middleware.RequireAuth(handlers.BarangKeluar))
	mux.HandleFunc("GET /dashboard/laporan_stok", middleware.RequireAuth(handlers.LaporanStok))

	// Manajemen Barang Route
	mux.HandleFunc("POST /barang/create", middleware.RequireAuth(handlers.CreateBarang))

	return mux

}
