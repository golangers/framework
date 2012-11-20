package filesession

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/gob"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"syscall"
	"time"
)

func init() {
	gob.Register([]interface{}{})
	gob.Register(map[int]interface{}{})
	gob.Register(map[string]interface{}{})
	gob.Register(map[interface{}]interface{}{})
}

func encodeGob(obj map[string]interface{}) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(obj)
	if err != nil {
		return []byte(""), err
	}
	return buf.Bytes(), nil
}

func decodeGob(encoded []byte) (map[string]interface{}, error) {
	buf := bytes.NewBuffer(encoded)
	dec := gob.NewDecoder(buf)
	var out map[string]interface{}
	err := dec.Decode(&out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func readFile(filePath string) ([]byte, error) {
	var content []byte
	f, err := os.OpenFile(filePath, os.O_RDONLY, 0777)
	if err == nil {
		fd := int(f.Fd())
		//防止死等待
		//等待10秒，否则报超时，并退出
		time.AfterFunc(10*time.Second, func() {
			syscall.Flock(fd, syscall.LOCK_UN)
			f.Close()
			err = errors.New("wait 10 second to unlock,but timeout")
		})
		if err = syscall.Flock(fd, syscall.LOCK_SH); err == nil {
			if content, err = ioutil.ReadAll(f); err == nil {
				if err = syscall.Flock(fd, syscall.LOCK_UN); err == nil {
					err = f.Close()
				}
			}
		}
	}

	return content, err
}

func writeFile(filePath string, content []byte) error {
	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0777)
	if err == nil {
		fd := int(f.Fd())
		//防止死等待
		//等待10秒，否则报超时，并退出
		time.AfterFunc(10*time.Second, func() {
			syscall.Flock(fd, syscall.LOCK_UN)
			f.Close()
			err = errors.New("wait 10 second to unlock,but timeout")
		})
		if err = syscall.Flock(fd, syscall.LOCK_EX); err == nil {
			if _, err = f.Write(content); err == nil {
				if err = f.Sync(); err == nil {
					if err = syscall.Flock(fd, syscall.LOCK_UN); err == nil {
						err = f.Close()
					}
				}
			}
		}
	}

	return err
}

func getSessionSign() string {
	var n int = 24
	b := make([]byte, n)
	io.ReadFull(rand.Reader, b)

	//return length:32
	return base64.URLEncoding.EncodeToString(b)
}

type SessionManager struct {
	CookieName    string
	expires       int
	sessionDir    string
	timerDuration time.Duration
}

func New(cookieName string, expires int, sessionDir string, timerDuration string) *SessionManager {
	if cookieName == "" {
		cookieName = "GoLangerSession"
	}

	if expires <= 0 {
		expires = 3600 * 24
	}

	if sessionDir == "" {
		sessionDir = os.TempDir() + "golangersession/"
	}

	os.MkdirAll(sessionDir, 0777)

	var dTimerDuration time.Duration

	if td, terr := time.ParseDuration(timerDuration); terr == nil {
		dTimerDuration = td
	} else {
		dTimerDuration, _ = time.ParseDuration("24h")
	}

	s := &SessionManager{
		CookieName:    cookieName,
		expires:       expires,
		sessionDir:    sessionDir,
		timerDuration: dTimerDuration,
	}

	time.AfterFunc(s.timerDuration, func() { s.GC() })

	return s
}

func (s *SessionManager) new(rw http.ResponseWriter) string {
	sessionSign := getSessionSign()

	bCookie := &http.Cookie{
		Name:     s.CookieName,
		Value:    url.QueryEscape(sessionSign),
		Path:     "/",
		HttpOnly: true,
	}

	http.SetCookie(rw, bCookie)

	return sessionSign
}

func (s *SessionManager) Get(rw http.ResponseWriter, req *http.Request) map[string]interface{} {
	m := map[string]interface{}{}

	if c, err := req.Cookie(s.CookieName); err == nil {
		sessionSign := c.Value
		if content, err := readFile(s.sessionDir + sessionSign + ".golanger"); err == nil {
			if len(content) > 0 {
				if dm, err := decodeGob(content); err == nil {
					m = dm
				} else {
					log.Panicln("filesession(decodeGob) error:", err)
				}
			}
		}
	} else {
		s.new(rw)
	}

	return m
}

func (s *SessionManager) Set(session map[string]interface{}, rw http.ResponseWriter, req *http.Request) {
	c, cerr := req.Cookie(s.CookieName)
	lsess := len(session)
	if cerr == nil {
		sessionSign := c.Value
		if lsess > 0 {
			var tryed bool
		TRY1:
			if encodeSession, err := encodeGob(session); err == nil {
				writeFile(s.sessionDir+sessionSign+".golanger", encodeSession)
			} else {
				if tryed {
					log.Panicln("filesession(encodeGob) error:", err)
				} else {
					for _, v := range session {
						gob.Register(v)
					}
					tryed = true
					goto TRY1
				}
			}
		} else {
			s.Clear(sessionSign)
		}
	} else {
		if lsess > 0 {
			var tryed bool
		TRY2:
			if encodeSession, err := encodeGob(session); err == nil {
				sessionSign := s.new(rw)
				writeFile(s.sessionDir+sessionSign+".golanger", encodeSession)
			} else {
				if tryed {
					log.Panicln("filesession(encodeGob) error:", err)
				} else {
					for _, v := range session {
						gob.Register(v)
					}
					tryed = true
					goto TRY2
				}
			}
		}
	}
}

func (s *SessionManager) Len() int64 {
	var slen int64
	if fs, err := filepath.Glob(s.sessionDir + "*.golanger"); err == nil {
		slen = int64(len(fs))
	}

	return slen
}

func (s *SessionManager) Clear(sessionSign string) {
	os.Remove(s.sessionDir + sessionSign + ".golanger")
}

func (s *SessionManager) GC() {
	if f, err := os.Open(s.sessionDir); err == nil {
		if fis, err := f.Readdir(-1); err == nil {
			for _, fi := range fis {
				if fi.ModTime().Unix()+int64(s.expires) <= time.Now().Unix() {
					os.Remove(s.sessionDir + fi.Name())
				}
			}
		}

		defer f.Close()
	}

	time.AfterFunc(s.timerDuration, func() { s.GC() })
}
