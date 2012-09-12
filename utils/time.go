package utils

import (
	"fmt"
	"time"
)

type Time struct {
	TimeFormat string
	TimeZone   string
}

func NewTime() *Time {
	return &Time{
		TimeFormat: "2006-01-02 15:04:05 +00:00",
		TimeZone:   "+00:00",
	}
}

func (t *Time) GetTimeToStr(tm int64, fm ...string) string {
	ti := time.Unix(tm, 0)
	var format string

	if len(fm) < 1 {
		format = "%04d-%02d-%02d %02d:%02d:%02d"
	} else {
		format = fm[0]
	}

	return fmt.Sprintf(format, ti.Year(), ti.Month(), ti.Day(), ti.Hour(), ti.Minute(), ti.Second())
}

func (t *Time) GetStrToTime(s string) int64 {
	s += " " + t.TimeZone
	ti, _ := time.Parse(t.TimeFormat, s)
	return ti.Unix()
}
