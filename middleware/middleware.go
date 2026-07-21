package middleware

import (
	"fmt"
	"net/http"

	"github.com/hilmanxcode/web-inventory-go/sessions"
)

func RequireAuth(next http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		// var key string
		c, err := r.Cookie("session_token")
		var key string
		if err == http.ErrNoCookie {
			fmt.Println("MASUK SET TERUS")
			key = sessions.SetSession(sessions.Session{}, w)
		}

		// if err != nil {

		// 	if err == http.ErrNoCookie {
		// 		fmt.Println("MASUK NO CUKIS")
		// 		log.Fatal(err.Error())
		// 	}

		// 	// fmt.Println("")

		// 	sessions.SetSession(sessions.Session{
		// 		ErrorMessages: []string{
		// 			"Bad Request",
		// 		},
		// 	}, w)

		// 	http.Redirect(w, r, "/", http.StatusBadRequest)
		// 	return

		// }

		fmt.Println(c.Value)

		// fmt.Println(key)

		username, err := sessions.GetUsernameSession(key, w)

		if err != nil {
			sessions.SetSession(sessions.Session{
				ErrorMessages: []string{
					"Unauthorized: Invalid Session",
				},
			}, w)
			http.Redirect(w, r, "/", http.StatusUnauthorized)
			return
		}

		if username == "" {
			sessions.SetSession(sessions.Session{
				ErrorMessages: []string{
					"Kamu Belum Login",
				},
			}, w)
			http.Redirect(w, r, "/", http.StatusUnauthorized)
			return
		}

		// kalau lolos kita kasih jalan
		next(w, r)

	}

}
