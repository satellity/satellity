package durable

import (
	"log"
)

// Logger is a custom log
type Logger struct {
	*log.Logger
}

// NewLogger create a new Logger
func NewLogger() *Logger {
	return &Logger{}
}

// Debug output log TODO
func (logger *Logger) Debug(v ...interface{}) {
	log.Println(v...)
}

// Debugf output log TODO
func (logger *Logger) Debugf(format string, v ...interface{}) {
	log.Printf(format, v...)
}

// Info output log TODO
func (logger *Logger) Info(v ...interface{}) {
	log.Println(v...)
}

// Infof output log TODO
func (logger *Logger) Infof(format string, v ...interface{}) {
	log.Printf(format, v...)
}

// Error output log TODO
func (logger *Logger) Error(v ...interface{}) {
	log.Println(v...)
}

// Errorf output log TODO
func (logger *Logger) Errorf(format string, v ...interface{}) {
	log.Printf(format, v...)
}
