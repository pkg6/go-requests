package requests

import (
	"log"
	"os"
)

type Logger struct {
	l    *log.Logger
	logo string
}

func NewLogger(l *log.Logger, logo string) *Logger {
	return &Logger{
		l:    l,
		logo: logo,
	}
}

func DefaultLogger() *Logger {
	return NewLogger(log.New(os.Stderr, "", log.Ldate|log.Lmicroseconds), "GoRequests")
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	l.output("【ERROR】"+l.logo+" "+format, v...)
}

func (l *Logger) Warnf(format string, v ...interface{}) {
	l.output("【WARN】"+l.logo+" "+format, v...)
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	l.output("【DEBUG】"+l.logo+" "+format, v...)
}

func (l *Logger) output(format string, v ...interface{}) {
	if len(v) == 0 {
		l.l.Print(format)
		return
	}
	l.l.Printf(format, v...)
}
