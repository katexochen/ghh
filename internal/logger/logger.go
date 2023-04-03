package logger

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

// DefaultLogger is the default logger implementation.
type DefaultLogger struct{}

// Infof logs an info message.
func (l *DefaultLogger) Infof(format string, args ...any) {
	log.Printf(fmt.Sprintf("INFO: %s", format), args...)
}

// Infoln logs an info message.
func (l *DefaultLogger) Infoln(args ...any) {
	l.Infof("%s\n", args...)
}

// Warnf logs a warning message.
func (l *DefaultLogger) Warnf(format string, args ...any) {
	log.Printf(fmt.Sprintf("WARN: %s", format), args...)
}

// Warnln logs a warning message.
func (l *DefaultLogger) Warnln(args ...any) {
	l.Warnf("%s\n", args...)
}

// Errorf logs an error message.
func (l *DefaultLogger) Errorf(format string, args ...any) {
	log.Printf(fmt.Sprintf("ERROR: %s", format), args...)
}

// Errorln logs an error message.
func (l *DefaultLogger) Errorln(args ...any) {
	l.Errorf("%s\n", args...)
}

// Debugf logs a debug message.
func (l *DefaultLogger) Debugf(_ string, _ ...any) {}

// Debugln logs a debug message.
func (l *DefaultLogger) Debugln(_ ...any) {}

// PrintJSON logs a JSON representation of v.
func (l *DefaultLogger) PrintJSON(_ string, _ any) {}

// VerboseLogger is a logger implementation that logs debug messages.
type VerboseLogger struct {
	DefaultLogger
}

// Debugf logs a debug message.
func (l *VerboseLogger) Debugf(format string, args ...any) {
	log.Printf(fmt.Sprintf("DEBUG: %s", format), args...)
}

// Debugln logs a debug message.
func (l *VerboseLogger) Debugln(args ...any) {
	l.Debugf("%s\n", args...)
}

// PrintJSON logs a JSON representation of v.
func (l *VerboseLogger) PrintJSON(msg string, v any) {
	l.Infoln(msg + ":")
	w := json.NewEncoder(os.Stdout)
	w.SetIndent("", "  ")
	err := w.Encode(v)
	if err != nil {
		panic(err)
	}
}
