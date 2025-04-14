package logger

import (
	"bytes"
	"encoding/json"
	"strings"
	"sync"
	"testing"
	"time"
)

// testLogEntry represents the JSON structure for capturing log output
type testLogEntry struct {
	Timestamp string          `json:"timestamp"`
	Level     string          `json:"level"`
	Message   string          `json:"message"`
	Fields    *map[string]any `json:"fields,omitempty"`
}

func TestLoggerOutputsCorrectJsonFormat(t *testing.T) {
	testCases := []struct {
		name     string
		level    logLevel
		message  string
		fields   map[string]any
		expected func(entry testLogEntry) bool
	}{
		{
			name:    "simple info message without fields",
			level:   info,
			message: "test info message",
			fields:  nil,
			expected: func(entry testLogEntry) bool {
				return entry.Level == string(info) &&
					entry.Message == "test info message" &&
					entry.Fields == nil
			},
		},
		{
			name:    "error message with fields",
			level:   error,
			message: "test error message",
			fields: map[string]any{
				"errorCode": 500,
				"source":    "database",
			},
			expected: func(entry testLogEntry) bool {
				if entry.Level != string(error) || entry.Message != "test error message" || entry.Fields == nil {
					return false
				}
				fields := *entry.Fields
				return fields["errorCode"] == float64(500) && fields["source"] == "database"
			},
		},
		{
			name:    "debug message with mixed type fields",
			level:   debug,
			message: "test debug message",
			fields: map[string]any{
				"bool":   true,
				"string": "value",
				"number": 42.5,
				"array":  []int{1, 2, 3},
			},
			expected: func(entry testLogEntry) bool {
				if entry.Level != string(debug) || entry.Message != "test debug message" || entry.Fields == nil {
					return false
				}
				fields := *entry.Fields
				return fields["bool"] == true &&
					fields["string"] == "value" &&
					fields["number"] == 42.5 &&
					len(fields) == 4 // array will be there but we don't check its exact value
			},
		},
	}

	var wg sync.WaitGroup
	for _, tc := range testCases {
		wg.Add(1)
		go func(tc struct {
			name     string
			level    logLevel
			message  string
			fields   map[string]any
			expected func(entry testLogEntry) bool
		}) {
			defer wg.Done()
			t.Run(tc.name, func(t *testing.T) {
				// Given
				buf := &bytes.Buffer{}
				logger := newJSONLogger(trace) // Use trace to ensure all logs are recorded
				logger.out = buf

				// When
				logger.log(tc.level, tc.message, tc.fields)

				// Then
				var entry testLogEntry
				if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
					t.Fatalf("Failed to unmarshal log entry: %v", err)
				}

				// Verify timestamp format
				_, err := time.Parse(time.RFC3339, entry.Timestamp)
				if err != nil {
					t.Errorf("Invalid timestamp format: %v", err)
				}

				// Verify the expected content
				if !tc.expected(entry) {
					t.Errorf("Log entry doesn't match expectations: %+v", entry)
				}
			})
		}(tc)
	}
	wg.Wait()
}

func TestFieldMergingPrioritizesLogFieldsOverGlobals(t *testing.T) {
	// Given
	buf := &bytes.Buffer{}
	logger := newJSONLogger(info)
	logger.out = buf

	// Global fields
	logger.fields = map[string]any{
		"app":      "test-app",
		"version":  "1.0.0",
		"conflict": "global-value", // This should be overridden
	}

	// When - Log with fields that include an override
	logFields := map[string]any{
		"conflict":   "local-value", // Should override the global value
		"request_id": "req-123",
	}

	logger.log(info, "test message", logFields)

	// Then - Parse and validate the output
	var entry testLogEntry
	if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
		t.Fatalf("Failed to unmarshal log entry: %v", err)
	}

	if entry.Fields == nil {
		t.Fatalf("Expected fields to be present in log output")
	}

	fields := *entry.Fields

	// Check that all fields are present
	expectedFields := map[string]any{
		"app":        "test-app",
		"version":    "1.0.0",
		"conflict":   "local-value", // Should be the local value
		"request_id": "req-123",
	}

	for k, v := range expectedFields {
		if fields[k] != v {
			t.Errorf("Field %s = %v, want %v", k, fields[k], v)
		}
	}
}

func TestAllLogLevelFunctionsProduceCorrectOutput(t *testing.T) {
	// Given
	origLogger := defaultLogger
	defer func() { defaultLogger = origLogger }()

	tests := []struct {
		name    string
		logFunc func(msg string, fields map[string]any)
		level   string
		message string
		fields  map[string]any
	}{
		{"Tracef logs with TRACE level", Tracef, "TRACE", "trace message", map[string]any{"trace": true}},
		{"Debugf logs with DEBUG level", Debugf, "DEBUG", "debug message", map[string]any{"debug": true}},
		{"Infof logs with INFO level", Infof, "INFO", "info message", map[string]any{"info": true}},
		{"Warnf logs with WARN level", Warnf, "WARN", "warn message", map[string]any{"warn": true}},
		{"Errorf logs with ERROR level", Errorf, "ERROR", "error message", map[string]any{"error": true}},
		// Note: We exclude Fatalf since it would exit the program
	}

	// Use a mutex to synchronize access to the buffer
	var mutex sync.Mutex

	// Run each test case separately to avoid buffer conflicts
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Use a fresh buffer and logger for each test
			buf := &bytes.Buffer{}
			testLogger := newJSONLogger(trace)
			testLogger.out = buf

			// Critical section: update the global logger
			mutex.Lock()
			defaultLogger = testLogger
			mutex.Unlock()

			// When
			tc.logFunc(tc.message, tc.fields)

			// Then
			output := buf.String()
			if len(output) == 0 {
				t.Fatalf("No log output produced")
			}

			var entry testLogEntry
			if err := json.Unmarshal([]byte(output), &entry); err != nil {
				t.Fatalf("Failed to unmarshal log entry: %v\nOutput was: %s", err, output)
			}

			if entry.Level != tc.level {
				t.Errorf("Log level = %s, want %s", entry.Level, tc.level)
			}

			if entry.Message != tc.message {
				t.Errorf("Message = %s, want %s", entry.Message, tc.message)
			}

			if entry.Fields == nil {
				t.Fatalf("Expected fields to be present")
			}

			fields := *entry.Fields
			for k, v := range tc.fields {
				if fields[k] != v {
					t.Errorf("Field %s = %v, want %v", k, fields[k], v)
				}
			}
		})
	}
}

func TestSetOutputChangesDestination(t *testing.T) {
	// Given
	origLogger := defaultLogger
	defer func() { defaultLogger = origLogger }()

	buf1 := &bytes.Buffer{}
	buf2 := &bytes.Buffer{}

	testLogger := newJSONLogger(info)
	testLogger.out = buf1
	defaultLogger = testLogger

	// When
	Infof("message to first buffer", nil)

	// Then
	if buf1.Len() == 0 {
		t.Error("Expected output in first buffer")
	}
	if buf2.Len() != 0 {
		t.Error("Expected second buffer to be empty")
	}

	// Given
	SetOutput(buf2)

	// When
	Infof("message to second buffer", nil)

	// Then
	if buf2.Len() == 0 {
		t.Error("Expected output in second buffer")
	}

	// Verify that both messages were logged correctly
	if !strings.Contains(buf1.String(), "message to first buffer") {
		t.Error("First buffer doesn't contain expected message")
	}

	if !strings.Contains(buf2.String(), "message to second buffer") {
		t.Error("Second buffer doesn't contain expected message")
	}
}

func TestSetFieldsAddsGlobalFields(t *testing.T) {
	// Given
	origLogger := defaultLogger
	defer func() { defaultLogger = origLogger }()

	buf := &bytes.Buffer{}
	testLogger := newJSONLogger(info)
	testLogger.out = buf
	defaultLogger = testLogger

	// When
	SetFields(map[string]any{
		"app": "test-app",
		"env": "testing",
	})

	Infof("test message", map[string]any{"local": "value"})

	// Then
	var entry testLogEntry
	if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
		t.Fatalf("Failed to unmarshal log entry: %v", err)
	}

	if entry.Fields == nil {
		t.Fatalf("Expected fields to be present")
	}

	fields := *entry.Fields

	expectedFields := map[string]string{
		"app":   "test-app",
		"env":   "testing",
		"local": "value",
	}

	for k, v := range expectedFields {
		if fields[k] != v {
			t.Errorf("Field %s = %v, want %v", k, fields[k], v)
		}
	}
}

func TestSetLevelFiltersMessages(t *testing.T) {
	// Given
	origLogger := defaultLogger
	defer func() { defaultLogger = origLogger }()

	// Test cases for different log levels
	testCases := []struct {
		name        string
		setLevel    string
		logMessages []struct {
			logFunc  func(msg string, fields map[string]any)
			level    string
			message  string
			expected bool // true if message should be logged, false otherwise
		}
	}{
		{
			name:     "TRACE level shows all logs",
			setLevel: "TRACE",
			logMessages: []struct {
				logFunc  func(msg string, fields map[string]any)
				level    string
				message  string
				expected bool
			}{
				{Tracef, "TRACE", "trace message", true},
				{Debugf, "DEBUG", "debug message", true},
				{Infof, "INFO", "info message", true},
				{Warnf, "WARN", "warn message", true},
				{Errorf, "ERROR", "error message", true},
			},
		},
		{
			name:     "DEBUG level filters out TRACE logs",
			setLevel: "DEBUG",
			logMessages: []struct {
				logFunc  func(msg string, fields map[string]any)
				level    string
				message  string
				expected bool
			}{
				{Tracef, "TRACE", "trace message", false},
				{Debugf, "DEBUG", "debug message", true},
				{Infof, "INFO", "info message", true},
				{Warnf, "WARN", "warn message", true},
				{Errorf, "ERROR", "error message", true},
			},
		},
		{
			name:     "INFO level filters out TRACE and DEBUG logs",
			setLevel: "INFO",
			logMessages: []struct {
				logFunc  func(msg string, fields map[string]any)
				level    string
				message  string
				expected bool
			}{
				{Tracef, "TRACE", "trace message", false},
				{Debugf, "DEBUG", "debug message", false},
				{Infof, "INFO", "info message", true},
				{Warnf, "WARN", "warn message", true},
				{Errorf, "ERROR", "error message", true},
			},
		},
		{
			name:     "WARN level filters out TRACE, DEBUG, and INFO logs",
			setLevel: "WARN",
			logMessages: []struct {
				logFunc  func(msg string, fields map[string]any)
				level    string
				message  string
				expected bool
			}{
				{Tracef, "TRACE", "trace message", false},
				{Debugf, "DEBUG", "debug message", false},
				{Infof, "INFO", "info message", false},
				{Warnf, "WARN", "warn message", true},
				{Errorf, "ERROR", "error message", true},
			},
		},
		{
			name:     "ERROR level filters out all except ERROR and FATAL logs",
			setLevel: "ERROR",
			logMessages: []struct {
				logFunc  func(msg string, fields map[string]any)
				level    string
				message  string
				expected bool
			}{
				{Tracef, "TRACE", "trace message", false},
				{Debugf, "DEBUG", "debug message", false},
				{Infof, "INFO", "info message", false},
				{Warnf, "WARN", "warn message", false},
				{Errorf, "ERROR", "error message", true},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Use a buffer to capture log output
			buf := &bytes.Buffer{}

			// Create a new logger and set it as default
			testLogger := newJSONLogger(info) // Default level will be overridden by SetLevel
			testLogger.out = buf
			defaultLogger = testLogger

			// Set the log level using the public function
			SetLevel(tc.setLevel)

			// Apply each log message and check if it appears in the output
			for _, msg := range tc.logMessages {
				// Clear the buffer before each message
				buf.Reset()

				// When - Log the message
				msg.logFunc(msg.message, map[string]any{"test": true})

				// Then - Check if the message was logged as expected
				output := buf.String()
				wasLogged := len(output) > 0

				if wasLogged != msg.expected {
					if msg.expected {
						t.Errorf("%s level should log %s messages, but didn't.\nOutput: %s",
							tc.setLevel, msg.level, output)
					} else {
						t.Errorf("%s level shouldn't log %s messages, but did.\nOutput: %s",
							tc.setLevel, msg.level, output)
					}
				}

				// If it was logged, verify the content
				if wasLogged {
					var entry testLogEntry
					if err := json.Unmarshal([]byte(output), &entry); err != nil {
						t.Fatalf("Failed to unmarshal log entry: %v", err)
					}

					// Verify level and message content
					if entry.Level != msg.level {
						t.Errorf("Log level = %s, want %s", entry.Level, msg.level)
					}

					if entry.Message != msg.message {
						t.Errorf("Message = %s, want %s", entry.Message, msg.message)
					}
				}
			}
		})
	}
}

// Test case insensitivity in SetLevel
func TestSetLevelCaseInsensitivity(t *testing.T) {
	// Given
	origLogger := defaultLogger
	defer func() { defaultLogger = origLogger }()

	testCases := []struct {
		input         string
		expectedLevel logLevel
	}{
		{"info", info},
		{"INFO", info},
		{"debug", debug},
		{"DEBUG", debug},
		{"trace", trace},
		{"TRACE", trace},
		{"warn", warn},
		{"WARN", warn},
		{"warning", warn},
		{"WARNING", warn},
		{"error", error},
		{"ERROR", error},
		{"fatal", fatal},
		{"FATAL", fatal},
		{"invalid", info}, // Default to info for invalid levels
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			// When
			SetLevel(tc.input)

			// Then
			if defaultLogger.level != tc.expectedLevel {
				t.Errorf("SetLevel(%q) set level to %q, want %q",
					tc.input, defaultLogger.level, tc.expectedLevel)
			}
		})
	}
}
