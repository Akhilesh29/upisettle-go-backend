package logger

import (
	"log"
)

// Logger is a minimal interface to keep logging abstract.
type Logger interface {
	Printf(format string, v ...any)
	Fatalf(format string, v ...any)
}

type stdLogger struct {
	prefix string
}

func (l *stdLogger) Printf(format string, v ...any) {
	log.Printf(l.prefix+format, v...)
}

func (l *stdLogger) Fatalf(format string, v ...any) {
	log.Fatalf(l.prefix+format, v...)
}

// New returns a simple stdlib-backed logger.
func New(env string) Logger {
	prefix := "[" + env + "] "
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	return &stdLogger{prefix: prefix}
}

