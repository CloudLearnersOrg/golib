// Package log provides a structured JSON logging system with support for levels,
// fields, and configurable outputs. It implements a global logger instance with
// thread-safe operations.
//
// Log Levels:
// The package supports six logging levels in order of increasing severity:
//   - TRACE: Verbose debugging information
//   - DEBUG: Debugging information
//   - INFO: General operational information (default)
//   - WARN: Warning messages for potentially harmful situations
//   - ERROR: Error messages for serious problems
//   - FATAL: Critical errors that result in program termination
//
// Basic Usage:
//
//	log.Init("Application starting")
//
//	log.Infof("Processing request", map[string]any{
//	    "method": "GET",
//	    "path": "/api/users",
//	})
//
//	if err != nil {
//	    log.Errorf("Request failed", map[string]any{
//	        "error": err,
//	        "status": 500,
//	    })
//	}
//
// Configuration:
// The logger can be configured in several ways:
//
//	// Set global fields that appear in every log entry
//	log.SetFields(map[string]any{
//	    "service": "api",
//	    "version": "1.0.0",
//	})
//
//	// Change log level
//	log.SetLevel("DEBUG")
//
//	// Redirect output
//	file, _ := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
//	log.SetOutput(file)
//
// JSON Output Format:
// All log entries are formatted as JSON with the following structure:
//
//	{
//	    "timestamp": "2024-04-22T15:04:05Z07:00",
//	    "level": "INFO",
//	    "message": "Request processed",
//	    "fields": {
//	        "method": "GET",
//	        "path": "/api/users",
//	        "duration": "125ms"
//	    }
//	}
//
// Thread Safety:
// The logger is safe for concurrent use by multiple goroutines. All logging
// operations are atomic and will not produce interleaved output.
//
// Performance:
// The logger implements several optimizations:
//   - JSON marshaling only occurs if the message will be logged
//   - Fields are allocated only when needed
//   - Log level checks are performed early to avoid unnecessary work
package log
