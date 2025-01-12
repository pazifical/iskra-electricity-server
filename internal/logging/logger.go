package logging

import "log"

const (
	ErrorLevel = iota
	WarningLevel
	InfoLevel
	DebugLevel
)

var LogLevel = WarningLevel

func Info(message string) {
	if LogLevel < InfoLevel {
		return
	}
	log.Printf("INFO:    %s", message)
}
func Debug(message string) {
	if LogLevel < DebugLevel {
		return
	}
	log.Printf("DEBUG:   %s", message)
}
func Error(message string) {
	if LogLevel < ErrorLevel {
		return
	}
	log.Printf("ERROR:   %s", message)
}
func Warning(message string) {
	if LogLevel < WarningLevel {
		return
	}
	log.Printf("WARNING: %s", message)
}
