package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Default global logger instance
var defaultLogger *jsonLogger

func init() {
	// Initialize the default logger with INFO level
	defaultLogger = newJSONLogger(info)
}

// newJSONLogger creates a new logger with the specified level
func newJSONLogger(level logLevel) *jsonLogger {
	return &jsonLogger{
		out:    os.Stdout,
		level:  level,
		fields: make(map[string]any),
	}
}

// shouldLog determines if a message at the given level should be logged
func (l *jsonLogger) shouldLog(level logLevel) bool {
	// Mapping log levels to numeric values for comparison
	levelValues := map[logLevel]int{
		trace: 0,
		debug: 1,
		info:  2,
		warn:  3,
		error: 4,
		fatal: 5,
	}

	// Get numeric values of the configured and message levels
	configuredValue, configExists := levelValues[l.level]
	messageValue, messageExists := levelValues[level]

	// Default to showing the message if levels are unknown
	if !configExists || !messageExists {
		return true
	}

	// Only show messages at or above the configured level
	return messageValue >= configuredValue
}

// log is the internal logging function that formats and outputs log entries
func (l *jsonLogger) log(level logLevel, msg string, fields map[string]any) {
	// Check if this level should be logged based on the configured level
	if !l.shouldLog(level) {
		return // Skip logging this message
	}

	// Create merged fields with app fields
	var mergedFields map[string]any

	// Only allocate map if we have fields to merge
	if len(l.fields) > 0 || (fields != nil && len(fields) > 0) {
		mapSize := len(l.fields)
		if fields != nil {
			mapSize += len(fields)
		}
		mergedFields = make(map[string]any, mapSize)

		// Add app fields first
		for k, v := range l.fields {
			mergedFields[k] = v
		}

		// Add log-specific fields, potentially overriding app fields
		if fields != nil {
			for k, v := range fields {
				mergedFields[k] = v
			}
		}
	}

	entry := logEntry{
		Timestamp: time.Now().Format(time.RFC3339),
		Level:     string(level),
		Message:   msg,
	}

	if mergedFields != nil && len(mergedFields) > 0 {
		entry.Fields = &mergedFields
	}

	jsonData, err := json.Marshal(entry)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling log entry: %v\n", err)
		return
	}

	fmt.Fprintln(l.out, string(jsonData))

	// If fatal, exit the program
	if level == fatal {
		os.Exit(1)
	}
}
