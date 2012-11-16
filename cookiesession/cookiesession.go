package cookiesession

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"encoding/gob"
	"log"
	"net/http"
	"strings"
	"time"
)

func encodeGob(obj interface{}) (string, error) {
	buf := bytes.NewBuffer(nil)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(obj)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
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

func encodeCookie(content interface{}, key, iv []byte) (string, error) {
	sessionGob, err := encodeGob(content)
	if err != nil {
		return "", err
	}
	//实现动态填充,达到aes.BlockSize的倍数,+4是为了后面提供4个字节来保存字符串长度使用的
	padLen := aes.BlockSize - (len(sessionGob)+4)%aes.BlockSize
	buf := bytes.NewBuffer(nil)
	var sessionLen int32 = (int32)(len(sessionGob))
	binary.Write(buf, binary.BigEndian, sessionLen)
	buf.WriteString(sessionGob)
	buf.WriteString(strings.Repeat("\000", padLen))
	sessionBytes := buf.Bytes()
	aesCipher, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	encrypter := cipher.NewCBCEncrypter(aesCipher, iv)
	encrypter.CryptBlocks(sessionBytes, sessionBytes)
	b64 := base64.URLEncoding.EncodeToString(sessionBytes)
	return b64, nil
}

func decodeCookie(encodedCookie string, key, iv []byte) (map[string]interface{}, error) {
	sessionBytes, err := base64.URLEncoding.DecodeString(encodedCookie)
	if err != nil {
		log.Printf("base64.Decodestring: %s\n", err)
		return nil, err
	}
	aesCipher, err := aes.NewCipher(key)
	if err != nil {
		log.Printf("aes.NewCipher: %s\n", err)
		return nil, err
	}

	decrypter := cipher.NewCBCDecrypter(aesCipher, iv)
	decrypter.CryptBlocks(sessionBytes, sessionBytes)

	buf := bytes.NewBuffer(sessionBytes)
	var gobLen int32
	binary.Read(buf, binary.BigEndian, &gobLen)
	gobBytes := sessionBytes[4 : 4+gobLen]
	session, err := decodeGob(gobBytes)
	if err != nil {
		log.Printf("decodeGob: %s\n", err)
		return nil, err
	}
	return session, nil
}

type SessionManager struct {
	CookieName string
	key        []byte
	iv         []byte
}

func New(cookieName, key string) *SessionManager {
	if cookieName == "" {
		cookieName = "GoLangerCookieSession"
	}

	if key == "" {
		key = "GoLanger Support CookieSession"
	}

	keySha1 := sha1.New()
	keySha1.Write([]byte(key))
	sum := keySha1.Sum(nil)
	return &SessionManager{
		CookieName: cookieName,
		key:        sum[:16],
		iv:         sum[4:],
	}
}

func (s *SessionManager) Get(req *http.Request) map[string]interface{} {
	cookie, err := req.Cookie(s.CookieName)
	if err != nil {
		return map[string]interface{}{}
	}
	session, err := decodeCookie(cookie.Value, s.key, s.iv)
	if err != nil {
		return map[string]interface{}{}
	}

	return session
}

func (s *SessionManager) Set(session map[string]interface{}, rw http.ResponseWriter, req *http.Request) {
	origCookie, err := req.Cookie(s.CookieName)
	var origCookieVal string
	if err != nil {
		origCookieVal = ""
	} else {
		origCookieVal = origCookie.Value
	}

	if len(session) == 0 {
		if origCookieVal != "" {
			cookie := &http.Cookie{
				Name:    s.CookieName,
				Value:   "",
				Path:    "/",
				Expires: time.Unix(0, 0).UTC(),
			}

			http.SetCookie(rw, cookie)
		}
	} else {
		if encoded, err := encodeCookie(session, s.key, s.iv); err == nil {
			if encoded != origCookieVal {
				cookie := &http.Cookie{
					Name:  s.CookieName,
					Value: encoded,
					Path:  "/",
				}

				http.SetCookie(rw, cookie)
			}
		}
	}
}
