package utils

import (
	"strconv"
	"strings"
	"time"
)

type Time struct {
	Layout string
	Zone   string
}

func NewTime() *Time {
	return &Time{
		Layout: "2006-01-02 15:04:05",
		Zone:   "-07:00",
	}
}

func (t *Time) UnixToStr(n int64, layouts ...string) string {
	layout := t.Layout
	if len(layouts) > 0 && layouts[0] != "" {
		layout = layouts[0]
	}

	var ss, ns string
	s := strconv.FormatInt(n, 10)
	l := len(s)
	if l > 10 {
		ss = s[:10]
		ns = s[10:]
		//000000000
		fillLen := 9 - len(ns)
		ns = ns + strings.Repeat("0", fillLen)
	} else {
		ss = s[:10]
	}

	si, _ := strconv.ParseInt(ss, 10, 64)
	ni, _ := strconv.ParseInt(ns, 10, 64)
	tm := time.Unix(si, ni)

	return tm.Format(layout)
}

func (t *Time) StrToUnix(s string) int64 {
	ti, _ := time.Parse(t.Layout+" "+t.Zone, s)
	return ti.Unix()
}
