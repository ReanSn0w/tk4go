package tools

import "log"

// Logger is a simple interface for logging
type Logger interface {
	Logf(format string, args ...interface{})
}

// BaseLogger is a simple implementation of Logger
func BaseLogger() Logger {
	return &baseLogger{}
}

type baseLogger struct{}

func (l *baseLogger) Logf(format string, args ...interface{}) {
	log.Printf(format, args...)
}
