package log

import (
	"os"
)

var std = New(os.Stderr, "", LstdFlags)

func SetLevel(level int) {
	std.SetLevel(level)
}

func SetFlags(flag int) {
	std.SetFlags(flag)
}

func Debug(v ...interface{}) {
	std.Debug(v...)
}

func Debugln(v ...interface{}) {
	std.Debugln(v...)
}

func Debugf(format string, v ...interface{}) {
	std.Debugf(format, v...)
}

func Info(v ...interface{}) {
	std.Info(v...)
}

func Infoln(v ...interface{}) {
	std.Infoln(v...)
}

func Infof(format string, v ...interface{}) {
	std.Infof(format, v...)
}

func Warn(v ...interface{}) {
	std.Warn(v...)
}

func Warnln(v ...interface{}) {
	std.Warnln(v...)
}

func Warnf(format string, v ...interface{}) {
	std.Warnf(format, v...)
}

func Error(v ...interface{}) {
	std.Error(v...)
}

func Errorln(v ...interface{}) {
	std.Errorln(v...)
}

func Errorf(format string, v ...interface{}) {
	std.Errorf(format, v...)
}

func Critical(v ...interface{}) {
	std.Critical(v...)
}

func Criticalln(v ...interface{}) {
	std.Criticalln(v...)
}

func Criticalf(format string, v ...interface{}) {
	std.Criticalf(format, v...)
}

func Panic(v ...interface{}) {
	std.Panic(v...)
}

func Panicln(v ...interface{}) {
	std.Panicln(v...)
}

func Panicf(format string, v ...interface{}) {
	std.Panicf(format, v...)
}

func Fatal(v ...interface{}) {
	std.Fatal(v...)
}

func Fatalln(v ...interface{}) {
	std.Fatalln(v...)
}

func Fatalf(format string, v ...interface{}) {
	std.Fatalf(format, v...)
}
