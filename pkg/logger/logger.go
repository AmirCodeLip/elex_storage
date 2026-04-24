package logger

import (
	"elex_storage/pkg/shared_kernel/models"
	"os"
	"path/filepath"
	"testing"
)

type Logger interface {
	Info(msg string, fields ...any)
	Error(msg string, fields ...any)
}

func NewLogger(cfg *models.ConfigEnv) Logger {
	// Clean logger to support linux
	loggerPath := filepath.Clean(cfg.LoggerPath)
	dir := filepath.Dir(loggerPath)

	// Create logger directories
	createrDirErr := os.MkdirAll(dir, os.ModePerm)
	if createrDirErr != nil {
		panic(createrDirErr)
	}

	// Create logger file
	file, err := os.Create(loggerPath)
	if err != nil {
		panic(err)
	}

	var logger Logger = newZapLogger(file)
	return logger
}

func NewLoggerMock(t *testing.T) Logger {
	return newMockLogger(t)
}
