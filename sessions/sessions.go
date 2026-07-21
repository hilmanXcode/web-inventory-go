package sessions

import (
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
)

var (
	SessionData = map[string]Session{}
	mu          sync.RWMutex
)

type Session struct {
	Key            string
	Username       string
	SuccessMessage string
	ErrorMessages  []string
	OldInput       map[string]string
	Expiry         time.Time
}

func (s Session) IsExpired() bool {
	return s.Expiry.Before(time.Now())
}

func SetSession(s Session, w http.ResponseWriter) string {
	sessionToken := uuid.NewString()
	expiresAt := time.Now().Add(120 * time.Second)

	mu.Lock()
	SessionData[sessionToken] = Session{
		Key:            sessionToken,
		Username:       s.Username,
		SuccessMessage: s.SuccessMessage,
		ErrorMessages:  s.ErrorMessages,
		Expiry:         expiresAt,
	}
	mu.Unlock()

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken,
		Expires: expiresAt,
		Path:    "/",
	})

	return sessionToken

}

func GetUsernameSession(cookieToken string, w http.ResponseWriter) (string, error) {
	mu.RLock()
	val, exists := SessionData[cookieToken]
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
	delete(SessionData, cookie)
	mu.Unlock()

	http.SetCookie(w, &http.Cookie{
		Name:   "session_token",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

}
