package utils

import (
	"crypto/md5"
	"fmt"
	"io"
	"strings"
)

type Strings string

func NewString(i interface{}) Strings {
	s := Strings(fmt.Sprintf("%v", i))
	return s
}

func (s Strings) String() string {
	return string(s)
}

func (s Strings) Md5() string {
	m := md5.New()
	io.WriteString(m, s.String())

	return fmt.Sprintf("%x", m.Sum(nil))
}

// convert like this: "HelloWorld" to "hello_world"
func (s Strings) SnakeCasedName() string {
	newstr := make([]rune, 0)
	firstTime := true

	for _, chr := range string(s) {
		if isUpper := 'A' <= chr && chr <= 'Z'; isUpper {
			if firstTime == true {
				firstTime = false
			} else {
				newstr = append(newstr, '_')
			}
			chr -= ('A' - 'a')
		}
		newstr = append(newstr, chr)
	}

	return string(newstr)
}

// convert like this: "hello_world" to "HelloWorld"
func (s Strings) TitleCasedName() string {
	newstr := make([]rune, 0)
	upNextChar := true

	for _, chr := range string(s) {
		switch {
		case upNextChar:
			upNextChar = false
			chr -= ('a' - 'A')
		case chr == '_':
			upNextChar = true
			continue
		}

		newstr = append(newstr, chr)
	}

	return string(newstr)
}

func (s Strings) PluralizeString() string {
	str := string(s)
	if strings.HasSuffix(str, "y") {
		str = str[:len(str)-1] + "ie"
	}
	return str + "s"
}
