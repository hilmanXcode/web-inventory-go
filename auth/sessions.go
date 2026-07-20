package auth

import (
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
)

var (
	Session = map[string]session{}
	mu      sync.RWMutex
)

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

	mu.Lock()
	Session[sessionToken] = session{
		Username: username,
		Expiry:   expiresAt,
	}
	mu.Unlock()

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken,
		Expires: expiresAt,
		Path:    "/",
	})

}

func GetUsernameSession(cookieToken string, w http.ResponseWriter) (string, error) {
	mu.RLock()
	val, exists := Session[cookieToken]
	mu.RUnlock()

	if !exists {
		return "", errors.New("session tidak ditemukan")
	}

	if val.IsExpired() {
		ClearSession(cookieToken, w)
		return "", errors.New("Session sudah expire")
	}

	return val.Username, nil

}

func ClearSession(cookie string, w http.ResponseWriter) {
	mu.Lock()
	delete(Session, cookie)
	mu.Unlock()

	http.SetCookie(w, &http.Cookie{
		Name:   "session_token",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

}
