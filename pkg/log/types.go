package log

import "io"

type logLevel string

const (
	tracelevel logLevel = "TRACE"
	debuglevel logLevel = "DEBUG"
	infolevel  logLevel = "INFO"
	warnlevel  logLevel = "WARN"
	errorlevel logLevel = "ERROR"
	fatallevel logLevel = "FATAL"
)

// jsonLogger is the internal logger implementation
type jsonLogger struct {
	level       logLevel
	out         io.Writer
	initialized bool
	fields      map[string]any
}

// logEntry represents a single log message structure
type logEntry struct {
	Timestamp string          `json:"timestamp"`
	Level     string          `json:"level"`
	Message   string          `json:"message"`
	Fields    *map[string]any `json:"fields,omitempty"`
}
