package log

func (l *Logger) Debug(v ...interface{}) {
	if LEVEL_DEBUG&l.level == 0 {
		return
	}

	l.print(logPrefixs[LEVEL_DEBUG]+" ", v...)
}

func (l *Logger) Debugln(v ...interface{}) {
	if LEVEL_DEBUG&l.level == 0 {
		return
	}

	l.println(logPrefixs[LEVEL_DEBUG]+" ", v...)
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	if LEVEL_DEBUG&l.level == 0 {
		return
	}

	l.printf(logPrefixs[LEVEL_DEBUG]+" ", format, v...)
}

func (l *Logger) Info(v ...interface{}) {
	if LEVEL_INFO&l.level == 0 {
		return
	}

	l.print(logPrefixs[LEVEL_INFO]+" ", v...)
}

func (l *Logger) Infoln(v ...interface{}) {
	if LEVEL_INFO&l.level == 0 {
		return
	}

	l.println(logPrefixs[LEVEL_INFO]+" ", v...)
}

func (l *Logger) Infof(format string, v ...interface{}) {
	if LEVEL_INFO&l.level == 0 {
		return
	}

	l.printf(logPrefixs[LEVEL_INFO]+" ", format, v...)
}

func (l *Logger) Warn(v ...interface{}) {
	if LEVEL_WARN&l.level == 0 {
		return
	}

	l.print(logPrefixs[LEVEL_WARN]+" ", v...)
}

func (l *Logger) Warnln(v ...interface{}) {
	if LEVEL_WARN&l.level == 0 {
		return
	}

	l.println(logPrefixs[LEVEL_WARN]+" ", v...)
}

func (l *Logger) Warnf(format string, v ...interface{}) {
	if LEVEL_WARN&l.level == 0 {
		return
	}

	l.printf(logPrefixs[LEVEL_WARN]+" ", format, v...)
}

func (l *Logger) Error(v ...interface{}) {
	if LEVEL_ERROR&l.level == 0 {
		return
	}

	l.print(logPrefixs[LEVEL_ERROR]+" ", v...)
}

func (l *Logger) Errorln(v ...interface{}) {
	if LEVEL_ERROR&l.level == 0 {
		return
	}

	l.println(logPrefixs[LEVEL_ERROR]+" ", v...)
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	if LEVEL_ERROR&l.level == 0 {
		return
	}

	l.printf(logPrefixs[LEVEL_ERROR]+" ", format, v...)
}

func (l *Logger) Critical(v ...interface{}) {
	if LEVEL_CRITICAL&l.level == 0 {
		return
	}

	l.print(logPrefixs[LEVEL_CRITICAL]+" ", v...)
}

func (l *Logger) Criticalln(v ...interface{}) {
	if LEVEL_CRITICAL&l.level == 0 {
		return
	}

	l.println(logPrefixs[LEVEL_CRITICAL]+" ", v...)
}

func (l *Logger) Criticalf(format string, v ...interface{}) {
	if LEVEL_CRITICAL&l.level == 0 {
		return
	}

	l.printf(logPrefixs[LEVEL_CRITICAL]+" ", format, v...)
}
