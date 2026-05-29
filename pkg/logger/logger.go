package logger

import (
	"elex_storage/pkg/shared_kernel/models"
	"errors"
	"fmt"
	"testing"
)

type Logger interface {
	Info(msg string, fields ...any)
	Error(msg string, fields ...any)
}

func NewLoggerMock(t *testing.T) Logger {
	return newMockLogger(t)
}

// handleInsertError processes different types of errors
func HandleCommonErr(err error, logger Logger) (bool, error) {
	var commonErr *models.CommonError
	if errors.As(err, &commonErr) {
		// Recognized domain error, return as-is
		return true, commonErr
	}

	// Unexpected database error
	logger.Error("Database error during file insertion: " + err.Error())
	return false, err
}

func RedConsoleLog(a ...any) {
	fmt.Println("\033[1;31m", fmt.Sprint(a...), "\033[0m")
}
