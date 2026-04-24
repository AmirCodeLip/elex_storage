package logger

import (
	"os"
	"runtime"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapLogger struct {
	logger *zap.Logger
}

func captureTrace() []any {
	fields := []any{}

	// Get caller info
	pc, file, line, ok := runtime.Caller(2)
	if !ok {
		return fields
	}

	fn := runtime.FuncForPC(pc)
	funcName := ""
	if fn != nil {
		funcName = fn.Name()
	}

	fields = append(fields,
		"func", funcName,
		"file", file,
		"line", line,
	)

	return fields
}

func (z *ZapLogger) Info(msg string, fields ...any) {
	// Capture trace fields
	traceFields := captureTrace()

	// Merge trace fields with extra fields
	allFields := append(traceFields, fields...)

	z.logger.Sugar().Infow(msg, allFields...)
}

func (z *ZapLogger) Error(msg string, fields ...any) {
	msg = strings.ReplaceAll(msg, "\"", "'")

	// Capture trace fields
	traceFields := captureTrace()

	// Merge trace fields with extra fields
	allFields := append(traceFields, fields...)

	z.logger.Sugar().Errorw(msg, allFields...)
}

func newZapLogger(file *os.File) *ZapLogger {

	// Configure Zap
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(file),
		zapcore.InfoLevel,
	)

	logger := zap.New(core)
	defer logger.Sync()

	return &ZapLogger{logger: logger}
}
