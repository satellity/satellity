package durable

import (
	"log"
)

type Logger struct {
	*log.Logger
}

func NewLogger() *Logger {
	return &Logger{}
}

func (logger *Logger) Debug(v ...interface{}) {
	log.Println(v...)
}

func (logger *Logger) Debugf(format string, v ...interface{}) {
	log.Printf(format, v...)
}

func (logger *Logger) Info(v ...interface{}) {
	log.Println(v...)
}

func (logger *Logger) Infof(format string, v ...interface{}) {
	log.Printf(format, v...)
}

func (logger *Logger) Error(v ...interface{}) {
	log.Println(v...)
}

func (logger *Logger) Errorf(format string, v ...interface{}) {
	log.Printf(format, v...)
}
