package middleware

import (
	"net/http"

	"github.com/hilmanxcode/web-inventory-go/sessions"
)

func RequireAuth(next http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		c, err := r.Cookie("session_token")

		if err != nil {
			sessions.SetSession(sessions.Session{
				ErrorMessages: []string{
					"Kamu Belum Login",
				},
			}, w)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		// Cookie ada, cek session dengan value dari cookie
		username, err := sessions.GetUsernameSession(c.Value, w)

		if err != nil {
			sessions.SetSession(sessions.Session{
				ErrorMessages: []string{
					"Unauthorized: Invalid Session",
				},
			}, w)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		if username == "" {
			sessions.SetSession(sessions.Session{
				ErrorMessages: []string{
					"Kamu Belum Login",
				},
			}, w)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		// kalau lolos kita kasih jalan
		next(w, r)

	}

}
