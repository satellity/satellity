package durable

import (
	"go.uber.org/zap"
)

// Logger is a custom log
type Logger struct {
	*zap.SugaredLogger
}

// NewLogger create a new Logger
func NewLogger(logger *zap.Logger) *Logger {
	return &Logger{logger.Sugar()}
}
