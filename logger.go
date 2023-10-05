package requests

import (
	"log"
	"os"
)

type Logger struct {
	l *log.Logger
}

func DefaultLogger() *Logger {
	l := &Logger{l: log.New(os.Stderr, "", log.Ldate|log.Lmicroseconds)}
	return l
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	l.output("【ERROR】GOREQUEST "+format, v...)
}

func (l *Logger) Warnf(format string, v ...interface{}) {
	l.output("【WARN】GOREQUEST "+format, v...)
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	l.output("【DEBUG】GOREQUEST "+format, v...)
}

func (l *Logger) output(format string, v ...interface{}) {
	if len(v) == 0 {
		l.l.Print(format)
		return
	}
	l.l.Printf(format, v...)
}
