package requests

import (
	"log"
	"os"
)

type Logger interface {
	Errorf(format string, v ...any)
	Warnf(format string, v ...any)
	Debugf(format string, v ...any)
}
type logger struct {
	l *log.Logger
}

var _ Logger = (*logger)(nil)

func DefaultLogger() *logger {
	l := &logger{l: log.New(os.Stderr, "", log.Ldate|log.Lmicroseconds)}
	return l
}

func (l *logger) Errorf(format string, v ...interface{}) {
	l.output("【ERROR】GOREQUEST "+format, v...)
}

func (l *logger) Warnf(format string, v ...interface{}) {
	l.output("【WARN】GOREQUEST "+format, v...)
}

func (l *logger) Debugf(format string, v ...interface{}) {
	l.output("【DEBUG】GOREQUEST "+format, v...)
}

func (l *logger) output(format string, v ...interface{}) {
	if len(v) == 0 {
		l.l.Print(format)
		return
	}
	l.l.Printf(format, v...)
}
