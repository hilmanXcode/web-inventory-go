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
	CurrentPage    string
	ErrorMessages  []string
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
		Name:     "session_token",
		Value:    sessionToken,
		Expires:  expiresAt,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteDefaultMode,
	})

	return sessionToken

}

func GetSession(key string) (Session, error) {
	mu.RLock()
	defer mu.RUnlock()

	val, exists := SessionData[key]

	if !exists {
		return Session{}, errors.New("Session tidak ditemukan")
	}

	return val, nil
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

func GetAndClearFlash(r *http.Request) (string, []string) {
	cookie, err := r.Cookie("session_token")

	if err != nil {
		return "", nil
	}

	mu.Lock()
	defer mu.Unlock()

	val, err := GetSession(cookie.Value)

	if err != nil {
		return "", nil
	}

	success := val.SuccessMessage
	errors := val.ErrorMessages

	val.SuccessMessage = ""
	val.ErrorMessages = nil
	SessionData[cookie.Value] = val

	return success, errors

}

func CleanupRoutine() {
	for {
		time.Sleep(1 * time.Minute)
		mu.Lock()
		for key, session := range SessionData {
			if session.IsExpired() {
				delete(SessionData, key)
			}
		}
		mu.Unlock()
	}
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
