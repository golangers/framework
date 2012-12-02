package log

import (
	"fmt"
	"io"
	"log"
	"runtime"
	"strings"
)

const (
	Lmodule   = log.Lshortfile << 1
	LstdFlags = log.LstdFlags | log.Lshortfile | Lmodule
)

//日志级别
const (
	LEVEL_DEBUG = 1 << iota
	LEVEL_INFO
	LEVEL_WARN
	LEVEL_ERROR
	LEVEL_CRITICAL
	LEVEL_PANIC
	LEVEL_FATAL
	LEVEL_ALL = LEVEL_DEBUG | LEVEL_INFO | LEVEL_WARN | LEVEL_ERROR | LEVEL_CRITICAL | LEVEL_PANIC | LEVEL_FATAL

	//默认日志级别为
	LEVEL_DEFAULT = LEVEL_INFO
)

var (
	logPrefixs = map[int]string{
		LEVEL_DEBUG:    "[DEBUG]",
		LEVEL_INFO:     "[INFO]",
		LEVEL_WARN:     "[WARN]",
		LEVEL_ERROR:    "[ERROR]",
		LEVEL_CRITICAL: "[CRITICAL]",
		LEVEL_PANIC:    "[PANIC]",
		LEVEL_FATAL:    "[FATAL]",
	}
)

type Logger struct {
	*log.Logger
	level int
}

func moduleOf(file string) string {
	pos := strings.LastIndex(file, "/")
	if pos != -1 {
		pos1 := strings.LastIndex(file[:pos], "/src/")
		if pos1 != -1 {
			return file[pos1+5 : pos]
		}
	}
	return "UNKNOWN"
}

func New(out io.Writer, prefix string, flag int) *Logger {
	return &Logger{
		log.Logger: log.New(out, prefix, flag),
		level:      LEVEL_DEFAULT,
	}
}

func (l *Logger) SetLevel(level int) {
	l.level = level
}

func (l *Logger) Output(calldepth int, s string) error {
	localCalldepth := calldepth + 2
	if Lmodule&l.level != 0 {
		if _, file, _, ok := runtime.Caller(localCalldepth); ok {
			l.SetPrefix(l.Prefix() + "[" + moduleOf(file) + "] ")
		}
	}

	//再加1层，是因为l.loger调用
	calldepth = localCalldepth + 1
	return l.Logger.Output(calldepth, s)
}

func (l *Logger) Println(v ...interface{}) {
	l.Output(2, fmt.Sprintln(v...))
}

func (l *Logger) Printf(format string, v ...interface{}) {
	l.Output(2, fmt.Sprintf(format, v...))
}

func (l *Logger) Debug(v ...interface{}) {
	if LEVEL_DEBUG&l.level == 0 {
		return
	}

	l.SetPrefix(logPrefixs[LEVEL_DEBUG] + " ")
	l.Println(v...)
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	if LEVEL_DEBUG&l.level == 0 {
		return
	}

	l.SetPrefix(logPrefixs[LEVEL_DEBUG] + " ")
	l.Printf(format, v...)
}

func (l *Logger) Info(v ...interface{}) {
	if LEVEL_INFO&l.level == 0 {
		return
	}

	l.SetPrefix(logPrefixs[LEVEL_INFO] + " ")
	l.Println(v...)
}

func (l *Logger) Infof(format string, v ...interface{}) {
	if LEVEL_INFO&l.level == 0 {
		return
	}

	l.SetPrefix(logPrefixs[LEVEL_INFO] + " ")
	l.Printf(format, v...)
}

func (l *Logger) Warn(v ...interface{}) {
	if LEVEL_WARN&l.level == 0 {
		return
	}

	l.SetPrefix(logPrefixs[LEVEL_WARN] + " ")
	l.Println(v...)
}

func (l *Logger) Warnf(format string, v ...interface{}) {
	if LEVEL_WARN&l.level == 0 {
		return
	}

	l.SetPrefix(logPrefixs[LEVEL_WARN] + " ")
	l.Printf(format, v...)
}

func (l *Logger) Error(v ...interface{}) {
	if LEVEL_ERROR&l.level == 0 {
		return
	}

	l.SetPrefix(logPrefixs[LEVEL_ERROR] + " ")
	l.Println(v...)
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	if LEVEL_ERROR&l.level == 0 {
		return
	}

	l.SetPrefix(logPrefixs[LEVEL_ERROR] + " ")
	l.Printf(format, v...)
}

func (l *Logger) Critical(v ...interface{}) {
	if LEVEL_CRITICAL&l.level == 0 {
		return
	}

	l.SetPrefix(logPrefixs[LEVEL_CRITICAL] + " ")
	l.Println(v...)
}

func (l *Logger) Criticalf(format string, v ...interface{}) {
	if LEVEL_CRITICAL&l.level == 0 {
		return
	}

	l.SetPrefix(logPrefixs[LEVEL_CRITICAL] + " ")
	l.Printf(format, v...)
}

func (l *Logger) Panic(v ...interface{}) {
	if LEVEL_PANIC&l.level == 0 {
		return
	}

	l.SetPrefix(logPrefixs[LEVEL_PANIC] + " ")
	l.Panicln(v...)
}

func (l *Logger) Panicf(format string, v ...interface{}) {
	if LEVEL_PANIC&l.level == 0 {
		return
	}

	l.SetPrefix(logPrefixs[LEVEL_PANIC] + " ")
	l.Panicf(format, v...)
}

func (l *Logger) Fatal(v ...interface{}) {
	if LEVEL_FATAL&l.level == 0 {
		return
	}

	l.SetPrefix(logPrefixs[LEVEL_FATAL] + " ")
	l.Fatalln(v...)
}

func (l *Logger) Fatalf(format string, v ...interface{}) {
	if LEVEL_FATAL&l.level == 0 {
		return
	}

	l.SetPrefix(logPrefixs[LEVEL_FATAL] + " ")
	l.Fatalf(format, v...)
}
