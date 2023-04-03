package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type Logger interface {
	Infof(format string, args ...any)
	Infoln(args ...any)
	Warnf(format string, args ...any)
	Warnln(args ...any)
	Errorf(format string, args ...any)
	Errorln(args ...any)
	Debugf(format string, args ...any)
	Debugln(args ...any)
	PrintJSON(msg string, v any)
}

type DefaultLogger struct{}

func (l *DefaultLogger) Infof(format string, args ...any) {
	log.Printf(fmt.Sprintf("INFO: %s", format), args...)
}

func (l *DefaultLogger) Infoln(args ...any) {
	l.Infof("%s\n", args...)
}

func (l *DefaultLogger) Warnf(format string, args ...any) {
	log.Printf(fmt.Sprintf("WARN: %s", format), args...)
}

func (l *DefaultLogger) Warnln(args ...any) {
	l.Warnf("%s\n", args...)
}

func (l *DefaultLogger) Errorf(format string, args ...any) {
	log.Printf(fmt.Sprintf("ERROR: %s", format), args...)
}

func (l *DefaultLogger) Errorln(args ...any) {
	l.Errorf("%s\n", args...)
}

func (l *DefaultLogger) Debugf(format string, args ...any) {}

func (l *DefaultLogger) Debugln(args ...any) {}

func (l *DefaultLogger) PrintJSON(msg string, v any) {}

type VerboseLogger struct {
	DefaultLogger
}

func (l *VerboseLogger) Debugf(format string, args ...any) {
	log.Printf(fmt.Sprintf("DEBUG: %s", format), args...)
}

func (l *VerboseLogger) Debugln(args ...any) {
	l.Debugf("%s\n", args...)
}

func (l *VerboseLogger) PrintJSON(msg string, v any) {
	l.Infoln(msg + ":")
	w := json.NewEncoder(os.Stdout)
	w.SetIndent("", "  ")
	err := w.Encode(v)
	if err != nil {
		panic(err)
	}
}
