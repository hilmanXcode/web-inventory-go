package auth

import (
	"net/http"
	"time"

	"github.com/google/uuid"
)

var Session = map[string]session{}

type session struct {
	Username string
	Expiry   time.Time
}

func (s session) IsExpired() bool {
	return s.Expiry.Before(time.Now())
}

func SetSession(username string, w http.ResponseWriter) {
	sessionToken := uuid.NewString()
	expiresAt := time.Now().Add(120 * time.Second)

	Session[sessionToken] = session{
		Username: username,
		Expiry:   expiresAt,
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken,
		Expires: expiresAt,
	})

}

func GetUsernameSession(cookie string) string {
	val := Session[cookie]

	return val.Username
}

func ClearSession(cookie string) {
	delete(Session, cookie)
}
