package globallog

import (
	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
)

var logger logr.Logger

// InitLogger initializes the global logger with the desired logging configuration.
func InitLogger() {
	zapLogger, _ := zap.NewProduction()
	logger = zapr.NewLogger(zapLogger)
}

// GetLogger returns the global logger instance.
func GetLogger() logr.Logger {
	return logger
}
