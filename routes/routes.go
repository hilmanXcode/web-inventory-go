package routes

import (
	"net/http"

	"github.com/hilmanxcode/web-inventory-go/auth"
	"github.com/hilmanxcode/web-inventory-go/handlers"
)

func SetupRouter() *http.ServeMux {

	mux := http.NewServeMux()

	mux.HandleFunc("GET /", handlers.HelloWorld)

	mux.HandleFunc("GET /setCookie", func(w http.ResponseWriter, r *http.Request) {
		auth.SetSession("hilmanXcode", w)

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

		username, err := auth.GetUsernameSession(c.Value)

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

		auth.ClearSession(c.Value)

		http.SetCookie(w, &http.Cookie{
			Name:   "session_token",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		})

		w.Write([]byte("Berhasil mendelete sebuah cookie"))
	})

	return mux

}
