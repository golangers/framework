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

func New(cookieName string, expires int, timerDuration string) *SessionManager {
	if cookieName == "" {
		cookieName = "GoLangerSession"
	}

	if expires <= 0 {
		expires = 3600 * 24
	}

	var dTimerDuration time.Duration

	if td, terr := time.ParseDuration(timerDuration); terr == nil {
		dTimerDuration = td
	} else {
		dTimerDuration, _ = time.ParseDuration("24h")
	}

	s := &SessionManager{
		CookieName:    cookieName,
		sessions:      map[string][2]map[string]interface{}{},
		expires:       expires,
		timerDuration: dTimerDuration,
	}

	time.AfterFunc(s.timerDuration, func() { s.GC() })

	return s
}

func (s *SessionManager) Get(rw http.ResponseWriter, req *http.Request) map[string]interface{} {
	var sessionSign string

	s.rmutex.RLock()
	defer s.rmutex.RUnlock()
	if c, err := req.Cookie(s.CookieName); err == nil {
		sessionSign = c.Value
		if sessionValue, ok := s.sessions[sessionSign]; ok {
			return sessionValue[1]
		}

	}

	sessionSign = s.new(rw)
	return s.sessions[sessionSign][1]
}

func (s *SessionManager) Set(rw http.ResponseWriter, req *http.Request) {
	s.rmutex.RLock()
	CookieName := s.CookieName
	s.rmutex.RUnlock()

	if c, err := req.Cookie(CookieName); err == nil {
		sessionSign := c.Value
		s.rmutex.RLock()
		lsess := len(s.sessions[sessionSign][1])
		s.rmutex.RUnlock()

		if lsess == 0 {
			s.Clear(sessionSign)
			bCookie := &http.Cookie{
				Name:     CookieName,
				Value:    "",
				Path:     "/",
				Expires:  time.Now().Add(-3600),
				HttpOnly: true,
			}

			http.SetCookie(rw, bCookie)
		}
	}
}

func (s *SessionManager) Len() int64 {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return int64(len(s.sessions))
}

func (s *SessionManager) new(rw http.ResponseWriter) string {
	timeNano := time.Now().UnixNano()
	s.rmutex.RLock()
	CookieName := s.CookieName
	sessionSign := s.sessionSign()
	s.rmutex.RUnlock()

	s.mutex.Lock()
	s.sessions[sessionSign] = [2]map[string]interface{}{
		map[string]interface{}{
			"create": timeNano,
		},
		map[string]interface{}{},
	}
	s.mutex.Unlock()

	bCookie := &http.Cookie{
		Name:     CookieName,
		Value:    url.QueryEscape(sessionSign),
		Path:     "/",
		HttpOnly: true,
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
