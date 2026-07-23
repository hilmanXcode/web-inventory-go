package httputil

import (
	"net/http"

	"github.com/hilmanxcode/web-inventory-go/sessions"
)

func RedirectWithError(w http.ResponseWriter, r *http.Request, errorMessage string, path string) {
	c, errCookie := r.Cookie("session_token")

	if errCookie != nil {
		sessions.SetSession(sessions.Session{
			ErrorMessages: []string{errorMessage},
		}, w)
		http.Redirect(w, r, path, http.StatusFound)
		return
	}

	val, err := sessions.GetSession(c.Value)
	if err != nil {
		sessions.SetSession(sessions.Session{
			ErrorMessages: []string{"Invalid Session", errorMessage},
		}, w)
		http.Redirect(w, r, path, http.StatusFound)
		return
	}

	val.ErrorMessages = []string{errorMessage}
	sessions.UpdateSession(c.Value, val)
	http.Redirect(w, r, path, http.StatusFound)
}
