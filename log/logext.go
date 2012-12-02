package log

func (l *Logger) Debug(v ...interface{}) {
	if LEVEL_DEBUG&l.level == 0 {
		return
	}

	l.SetPrefix(logPrefixs[LEVEL_DEBUG] + " ")
	l.Print(v...)
}

func (l *Logger) Debugln(v ...interface{}) {
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
	l.Print(v...)
}

func (l *Logger) Infoln(v ...interface{}) {
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
	l.Print(v...)
}

func (l *Logger) Warnln(v ...interface{}) {
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
	l.Print(v...)
}

func (l *Logger) Errorln(v ...interface{}) {
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
	l.Print(v...)
}

func (l *Logger) Criticalln(v ...interface{}) {
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
