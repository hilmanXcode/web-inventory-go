package sessions

import (
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
)

var (
	sessionData = map[string]Session{}
	mu          sync.RWMutex
)

type Session struct {
	Key            string
	Username       string
	SuccessMessage string
	CurrentPage    string
	ErrorMessages  []string
	Expiry         time.Time
	OldInput       map[string]string
}

func (s Session) IsExpired() bool {
	return s.Expiry.Before(time.Now())
}

func SetSession(s Session, w http.ResponseWriter) string {
	sessionToken := uuid.NewString()
	expiresAt := time.Now().Add(12 * time.Hour)

	mu.Lock()
	sessionData[sessionToken] = Session{
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

	val, exists := sessionData[key]

	if !exists {
		return Session{}, errors.New("Session tidak ditemukan")
	}

	return val, nil
}

func GetUsernameSession(cookieToken string, w http.ResponseWriter) (string, error) {
	mu.RLock()
	val, exists := sessionData[cookieToken]
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

func UpdateSession(key string, s Session) error {
	mu.Lock()
	defer mu.Unlock()

	existing, exists := sessionData[key]
	if !exists {
		return errors.New("session tidak ditemukan")
	}

	s.Key = key
	if s.Expiry.IsZero() {
		s.Expiry = existing.Expiry
	}
	sessionData[key] = s

	return nil
}

func GetAndClearOldInput(r *http.Request) (map[string]string, error) {
	cookie, err := r.Cookie("session_token")

	if err != nil {
		return map[string]string{}, err
	}

	mu.Lock()
	defer mu.Unlock()

	val, exists := sessionData[cookie.Value]

	if !exists {
		return map[string]string{}, err
	}

	oldInput := val.OldInput

	val.OldInput = nil
	sessionData[cookie.Value] = val

	return oldInput, nil
}

func GetAndClearFlash(r *http.Request) (string, []string, error) {
	cookie, err := r.Cookie("session_token")

	if err != nil {
		return "", nil, err
	}

	mu.Lock()
	defer mu.Unlock()

	val, exists := sessionData[cookie.Value]

	if !exists {
		return "", nil, err
	}

	success := val.SuccessMessage
	errors := val.ErrorMessages

	val.SuccessMessage = ""
	val.ErrorMessages = nil
	sessionData[cookie.Value] = val

	return success, errors, nil

}

func CleanupRoutine() {
	for {
		time.Sleep(1 * time.Minute)
		mu.Lock()
		for key, session := range sessionData {
			if session.IsExpired() {
				delete(sessionData, key)
			}
		}
		mu.Unlock()
	}
}

func ClearSession(cookie string, w http.ResponseWriter) {
	mu.Lock()
	delete(sessionData, cookie)
	mu.Unlock()

	http.SetCookie(w, &http.Cookie{
		Name:   "session_token",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

}
