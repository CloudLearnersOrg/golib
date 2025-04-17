package log

import (
	"fmt"
	"io"
	"os"
)

// Init initializes the logger with a message
func Init(msg string) {
	defaultLogger.initialized = true
	defaultLogger.log(infolevel, msg, nil)
}

// Tracef logs a message at trace level with fields
func Tracef(msg string, fields map[string]any) {
	defaultLogger.log(tracelevel, msg, fields)
}

// Debugf logs a message at debug level with fields
func Debugf(msg string, fields map[string]any) {
	defaultLogger.log(debuglevel, msg, fields)
}

// Infof logs a message at info level with fields
func Infof(msg string, fields map[string]any) {
	defaultLogger.log(infolevel, msg, fields)
}

// Warnf logs a message at warn level with fields
func Warnf(msg string, fields map[string]any) {
	defaultLogger.log(warnlevel, msg, fields)
}

// Errorf logs a message at error level with fields
func Errorf(msg string, fields map[string]any) {
	defaultLogger.log(errorlevel, msg, fields)
}

// Fatalf logs a message at fatal level with fields and then exits
func Fatalf(msg string, fields map[string]any) {
	defaultLogger.log(fatallevel, msg, fields)
}

// SetFields sets global fields that will be included in all log entries
func SetFields(fields map[string]any) {
	for k, v := range fields {
		defaultLogger.fields[k] = v
	}
}

// SetOutput sets the output destination for the logger
func SetOutput(out io.Writer) {
	defaultLogger.out = out
}
func SetLevel(level string) {
	var lvl logLevel

	// Case-insensitive matching of log level strings
	switch level {
	case "TRACE", "trace":
		lvl = tracelevel
	case "DEBUG", "debug":
		lvl = debuglevel
	case "INFO", "info":
		lvl = infolevel
	case "WARN", "warn", "WARNING", "warning":
		lvl = warnlevel
	case "ERROR", "error":
		lvl = errorlevel
	case "FATAL", "fatal":
		lvl = fatallevel
	default:
		// Default to INFO if invalid level
		lvl = infolevel
		fmt.Fprintf(os.Stderr, "Warning: Unknown log level '%s', defaulting to INFO\n", level)
	}

	defaultLogger.level = lvl
}
