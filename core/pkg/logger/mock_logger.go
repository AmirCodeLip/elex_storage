package logger

import (
	"fmt"
	"strings"
	"testing"
)

type MockLogger struct {
	t *testing.T
}

func (ml *MockLogger) Info(msg string, fields ...any) {
	// Capture trace fields
	traceFields := captureTrace()

	// Merge trace fields with extra fields
	allFields := append(traceFields, fields...)
	for _, log := range allFields {
		msg += fmt.Sprintf("\n %s", log)
	}
	ml.t.Log(msg, allFields)
}

func (ml *MockLogger) Error(msg string, fields ...any) {
	msg = strings.ReplaceAll(msg, "\"", "'")

	// Capture trace fields
	traceFields := captureTrace()

	// Merge trace fields with extra fields
	allFields := append(traceFields, fields...)
	for _, log := range allFields {
		msg += fmt.Sprintf("\n %s", log)
	}
	ml.t.Error(msg, allFields)
}

func newMockLogger(t *testing.T) *MockLogger {

	return &MockLogger{t: t}
}
