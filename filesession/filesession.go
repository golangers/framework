/*
由于跨程序共享session使用的是Gob，在session中有用户自定义类型时，如果另一程序未预先注册此UserType,会报错
因此，暂不考虑跨程序共享session时的资源互斥问题，顾将//(1)处代码注释掉

为避免任何隐藏麻烦，在Encode时如失败不会去尝试gob.Register注册，会直接报错。
以免上一次Encode失败时尝试注册的的类型再Decode时不被识别到，而报错。
顾在程序运行时，尽可能完整的做好测试，再遇到Encode出错时，预注册好此用户自定义类型
注册方式可以：
import "encoding/gob"
func init() {
	gob.Register([UserType的初始化])
}

如果明确知晓session的传输类型，并且在linux系统环境下，可以将//(1)代码解除注释，将//(2)代码注释即可
*/
package filesession

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/gob"
	//(1)
	//"errors"
	//(1)
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sync"
	//(1)
	//"syscall"
	//(1)
	"time"
)

var (
	sessionLock sync.RWMutex
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
	//(1)
	//由于跨程序共享session使用的是Gob，在session中有用户自定义类型时，如果另一程序未预先注册此UserType,会报错
	//因此，暂不考虑跨程序共享session时的资源互斥问题，顾以下代码先注释
	/*f, err := os.OpenFile(filePath, os.O_RDONLY, 0777)
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
	}*/
	//(1)

	//(2)
	var err error
	sessionLock.RLock()
	content, err = ioutil.ReadFile(filePath)
	sessionLock.RUnlock()
	//(2)

	return content, err
}

func writeFile(filePath string, content []byte) error {
	//(1)
	//由于跨程序共享session使用的是Gob，在session中有用户自定义类型时，如果另一程序未预先注册此UserType,会报错
	//因此，暂不考虑跨程序共享session时的资源互斥问题，顾以下代码先注释
	/*f, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0777)
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
	}*/
	//(1)

	//(2)
	sessionLock.Lock()
	err := ioutil.WriteFile(filePath, content, os.ModePerm)
	sessionLock.Unlock()
	//(2)

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
		sessionDir = "./tmp/" + "golangersession/"
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
			if encodeSession, err := encodeGob(session); err == nil {
				writeFile(s.sessionDir+sessionSign+".golanger", encodeSession)
			} else {
				log.Panicln("filesession(encodeGob) error:", err)
			}
		} else {
			s.Clear(sessionSign)
		}
	} else {
		if lsess > 0 {
			if encodeSession, err := encodeGob(session); err == nil {
				sessionSign := s.new(rw)
				writeFile(s.sessionDir+sessionSign+".golanger", encodeSession)
			} else {
				log.Panicln("filesession(encodeGob) error:", err)
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
