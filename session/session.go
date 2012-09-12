package session

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type SessionManager struct {
	CookieName    string
	rmutex        sync.RWMutex
	mutex         sync.Mutex
	sessions      map[string][2]map[string]interface{}
	expires       int
	timerDuration time.Duration
}

func New(cookieName string, expires int, timerDuration time.Duration) *SessionManager {
	if cookieName == "" {
		cookieName = "GoLangerSession"
	}

	if expires <= 0 {
		expires = 3600
	}

	if timerDuration <= 0 {
		timerDuration, _ = time.ParseDuration("24h")
	}

	s := &SessionManager{
		CookieName:    cookieName,
		sessions:      map[string][2]map[string]interface{}{},
		expires:       expires,
		timerDuration: timerDuration,
	}

	time.AfterFunc(s.timerDuration, func() { s.GC() })

	return s
}

func (s *SessionManager) Get(rw http.ResponseWriter, req *http.Request) map[string]interface{} {
	var sessionSign string

	if c, err := req.Cookie(s.CookieName); err == nil {
		sessionSign = c.Value
		s.rmutex.RLock()
		if sessionValue, ok := s.sessions[sessionSign]; ok {
			s.rmutex.RUnlock()
			return sessionValue[1]
		}

		s.rmutex.RUnlock()
	}

	s.mutex.Lock()
	sessionSign = s.new(rw)
	s.mutex.Unlock()

	return s.sessions[sessionSign][1]
}

func (s *SessionManager) Len() int64 {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return int64(len(s.sessions))
}

func (s *SessionManager) new(rw http.ResponseWriter) string {
	timeNano := time.Now().UnixNano()
	sessionSign := s.sessionSign()
	s.sessions[sessionSign] = [2]map[string]interface{}{
		map[string]interface{}{
			"create": timeNano,
		},
		map[string]interface{}{},
	}

	bCookie := &http.Cookie{
		Name:  s.CookieName,
		Value: url.QueryEscape(sessionSign),
		Path:  "/",
	}

	http.SetCookie(rw, bCookie)

	return sessionSign
}

func (s *SessionManager) Clear(sessionSign string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.sessions, sessionSign)
}

func (s *SessionManager) GC() {
	s.rmutex.RLock()
	for sessionSign, sess := range s.sessions {
		if (sess[0]["create"].(int64) + int64(s.expires)) <= time.Now().Unix() {
			s.mutex.Lock()
			delete(s.sessions, sessionSign)
			s.mutex.Unlock()
		}
	}

	s.rmutex.RUnlock()

	time.AfterFunc(s.timerDuration, func() { s.GC() })
}

func (s *SessionManager) sessionSign() string {
	var n int = 24
	b := make([]byte, n)
	io.ReadFull(rand.Reader, b)

	//return length:32
	return base64.URLEncoding.EncodeToString(b)
}
