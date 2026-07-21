package routes

import (
	"net/http"

	"github.com/hilmanxcode/web-inventory-go/handlers"
	"github.com/hilmanxcode/web-inventory-go/sessions"
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
	mux.HandleFunc("GET /", handlers.LoginPage)
	mux.HandleFunc("POST /", handlers.LoginHandler)

	// Register Route
	mux.HandleFunc("POST /register", handlers.RegisterHandler)

	mux.HandleFunc("GET /setCookie", func(w http.ResponseWriter, r *http.Request) {
		sessions.SetSession(sessions.Session{
			Username: "hilmanxcode",
		}, w)

		w.Write([]byte("Berhasil menset cookie"))
	})

	mux.HandleFunc("GET /myCookie", func(w http.ResponseWriter, r *http.Request) {

		c, err := r.Cookie("session_token")

		if err != nil {

			if err == http.ErrNoCookie {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			w.WriteHeader(http.StatusBadRequest)
			return
		}

		username, err := sessions.GetUsernameSession(c.Value, w)

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(err.Error()))
			return
		}

		w.Write([]byte("Halo " + username + ", cookie name: " + c.Name))
	})

	mux.HandleFunc("GET /clearCookie", func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("session_token")

		if err != nil {

			if err == http.ErrNoCookie {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			w.WriteHeader(http.StatusBadRequest)
			return
		}

		sessions.ClearSession(c.Value, w)

		w.Write([]byte("Berhasil mendelete sebuah cookie"))
	})

	return mux

}
