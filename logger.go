package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type Logger struct{}

func (l *Logger) Infof(format string, args ...any) {
	log.Printf(fmt.Sprintf("INFO: %s", format), args...)
}

func (l *Logger) Infoln(args ...any) {
	l.Infof("%s\n", args...)
}

func (l *Logger) Warnf(format string, args ...any) {
	log.Printf(fmt.Sprintf("WARN: %s", format), args...)
}

func (l *Logger) Warnln(args ...any) {
	l.Warnf("%s\n", args...)
}

func (l *Logger) Errorf(format string, args ...any) {
	log.Printf(fmt.Sprintf("ERROR: %s", format), args...)
}

func (l *Logger) Errorln(args ...any) {
	l.Errorf("%s\n", args...)
}

func (l *Logger) Debugf(format string, args ...any) {
	log.Printf(fmt.Sprintf("DEBUG: %s", format), args...)
}

func (l *Logger) Debugln(args ...any) {
	l.Debugf("%s\n", args...)
}

func (l *Logger) PrintJSON(msg string, v any) {
	l.Infoln(msg + ":")
	w := json.NewEncoder(os.Stdout)
	w.SetIndent("", "  ")
	err := w.Encode(v)
	if err != nil {
		panic(err)
	}
}
