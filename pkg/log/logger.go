package log

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
	defaultLogger = newJSONLogger(infolevel)
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
		tracelevel: 0,
		debuglevel: 1,
		infolevel:  2,
		warnlevel:  3,
		errorlevel: 4,
		fatallevel: 5,
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
	// Fixed S1009: Removed unnecessary nil check since len() for nil maps is defined as zero
	if len(l.fields) > 0 || len(fields) > 0 {
		mapSize := len(l.fields) + len(fields)
		mergedFields = make(map[string]any, mapSize)

		// Add app fields first
		for k, v := range l.fields {
			mergedFields[k] = v
		}

		// Fixed S1031: Removed unnecessary nil check around range
		// since ranging over a nil map is valid and produces 0 iterations
		for k, v := range fields {
			mergedFields[k] = v
		}
	}

	entry := logEntry{
		Timestamp: time.Now().Format(time.RFC3339),
		Level:     string(level),
		Message:   msg,
	}

	// Fixed S1009: Removed unnecessary nil check
	if len(mergedFields) > 0 {
		entry.Fields = &mergedFields
	}

	jsonData, err := json.Marshal(entry)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling log entry: %v\n", err)
		return
	}

	// Fixed errcheck: Now checking error from fmt.Fprintln
	if _, err := fmt.Fprintln(l.out, string(jsonData)); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing log entry: %v\n", err)
		return
	}

	// If fatal, exit the program
	if level == fatallevel {
		os.Exit(1)
	}
}
