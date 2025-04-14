package logger

import "io"

type logLevel string

const (
	trace logLevel = "TRACE"
	debug logLevel = "DEBUG"
	info  logLevel = "INFO"
	warn  logLevel = "WARN"
	error logLevel = "ERROR"
	fatal logLevel = "FATAL"
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
