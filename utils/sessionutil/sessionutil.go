package sessionutil

import (
	"net/http"

	"github.com/hilmanxcode/web-inventory-go/sessions"
)

func GetSessionData(w http.ResponseWriter, r *http.Request, path string) (sessions.Session, error) {

	c, err := r.Cookie("session_token")

	if err != nil {
		sessions.SetSession(sessions.Session{
			ErrorMessages: []string{
				"Invalid Session",
			},
		}, w)
		http.Redirect(w, r, path, http.StatusFound)
		return sessions.Session{}, err
	}

	val, err := sessions.GetSession(c.Value)

	if err != nil {
		sessions.SetSession(sessions.Session{
			ErrorMessages: []string{
				"Invalid session",
			},
		}, w)
		http.Redirect(w, r, path, http.StatusFound)
		return sessions.Session{}, err
	}

	return val, nil

}
