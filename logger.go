package requests

import (
	"fmt"
	"log"
	"os"
)

const (
	loglevelERROR = "ERROR"
	loglevelWARN  = "WARN"
	loglevelDEBUG = "DEBUG"
)

// Color defines a single SGR Code
type color string

const (
	green   color = "\033[97;42m"
	white   color = "\033[90;47m"
	yellow  color = "\033[90;43m"
	red     color = "\033[97;41m"
	blue    color = "\033[97;44m"
	magenta color = "\033[97;45m"
	cyan    color = "\033[97;46m"
	reset   color = "\033[0m"
)

func loglevelColor(level string) string {
	switch level {
	case loglevelERROR:
		return fmt.Sprintf("%s %s %s", red, level, reset)
	case loglevelWARN:
		return fmt.Sprintf("%s %s %s", green, level, reset)
	case loglevelDEBUG:
		return fmt.Sprintf("%s %s %s", blue, level, reset)
	default:
		return level
	}
}

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

func (l *Logger) Errorf(format string, v ...any) {
	l.output(loglevelERROR, format, v...)
}

func (l *Logger) Warnf(format string, v ...any) {
	l.output(loglevelWARN, format, v...)
}

func (l *Logger) Debugf(format string, v ...any) {
	l.output(loglevelDEBUG, format, v...)
}

func (l *Logger) output(level, format string, v ...any) {
	l.l.Printf(fmt.Sprintf("[%s] |%s| ", l.logo, loglevelColor(level))+format+"\n", v...)
}
