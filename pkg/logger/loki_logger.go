package logger

import (
	"context"
	"elex_storage/pkg/shared_kernel/models"
	"log/slog"
	"strings"

	"github.com/bearsoft-fi/slogloki"
)

type LokiLogger struct {
	logger *slog.Logger
}

func (l *LokiLogger) Info(msg string, fields ...any) {
	// Capture trace fields
	traceFields := captureTrace()

	// Merge trace fields with extra fields
	allFields := append(traceFields, fields...)

	ctx := context.Background()
	l.logger.Log(ctx, slog.LevelInfo, msg, allFields...)
}

func (l *LokiLogger) Error(msg string, fields ...any) {
	msg = strings.ReplaceAll(msg, "\"", "'")

	// Capture trace fields
	traceFields := captureTrace()

	// Merge trace fields with extra fields
	allFields := append(traceFields, fields...)

	ctx := context.Background()
	l.logger.Log(ctx, slog.LevelError, msg, allFields...)
}

func NewLokiLogger(cfg *models.ConfigEnv) (Logger, error) {
	config := slogloki.NewDefaultConfig(cfg.Loki.APIAddress)
	handler, err := slogloki.NewLokiHandler(config, map[string]string{
		"service": cfg.ServiceName,
		"host":    cfg.ServiceName,
	})
	if err != nil {
		return nil, err
	}
	return &LokiLogger{logger: slog.New(handler)}, nil
}
