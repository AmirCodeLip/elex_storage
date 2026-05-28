package logger

import "runtime"

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
