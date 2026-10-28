package serveur

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gofrs/uuid"
)

type session struct {
	userID    string
	expiresAt time.Time
}

var (
	sessions = make(map[string]session)
	mu       sync.Mutex
)

const SessionDuration = 24 * time.Hour

func CreateSession(userID string) string {
	mu.Lock()
	defer mu.Unlock()

	u2, err := uuid.NewV4()
	if err != nil {
		return ""
	}

	sessions[u2.String()] = session{
		userID:    userID,
		expiresAt: time.Now().Add(SessionDuration),
	}
	return u2.String()
}

func Validate(sessionID string) (string, bool) {
	mu.Lock()
	defer mu.Unlock()

	s, existe := sessions[sessionID]
	if !existe {
		return "", false
	}

	if time.Now().After(s.expiresAt) {
		delete(sessions, sessionID)
		return "", false
	}

	return s.userID, true
}

func DeleteSession(sessionID string) {
	mu.Lock()
	defer mu.Unlock()
	delete(sessions, sessionID)
}

func ClientID(r *http.Request) (int, bool) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return 0, false
	}

	userID, valide := Validate(cookie.Value)
	if !valide {
		return 0, false
	}

	id, err := strconv.Atoi(userID)
	if err != nil {
		return 0, false
	}

	return id, true
}
