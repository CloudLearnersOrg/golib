// Package logger provides a Gin middleware for structured logging of HTTP requests.
// It automatically tracks request duration, generates trace IDs, and logs requests
// with appropriate severity levels based on response status codes.
//
// Features:
//   - Automatic trace ID generation and propagation
//   - Request/Response timing measurement
//   - Status code based log levels
//   - Structured logging with consistent fields
//
// Basic Usage:
//
//	router := gin.New()
//	router.Use(logger.Middleware())
//
// Log Fields:
// Each log entry includes the following structured fields:
//   - trace_id: Unique identifier for request tracing
//   - http.response.status_code: HTTP status code
//   - http.request.method: HTTP method (GET, POST, etc.)
//   - http.route: Request path
//   - server.address: Server hostname
//   - http.response.latency: Request processing duration
//
// Log Levels:
// The middleware uses different log levels based on the response status:
//   - INFO: For successful responses (2xx, 3xx)
//   - WARN: For client errors (4xx)
//   - ERROR: For server errors (5xx)
//
// Trace ID Handling:
// The middleware handles trace IDs in the following way:
//  1. Checks for existing X-Trace-ID in request headers
//  2. Generates new UUID if no trace ID exists
//  3. Sets trace ID in Gin context for downstream use
//  4. Adds trace ID to response headers
//
// Example Log Output:
//
//	{
//	    "level": "INFO",
//	    "msg": "incoming request completed",
//	    "trace_id": "550e8400-e29b-41d4-a716-446655440000",
//	    "http.response.status_code": 200,
//	    "http.request.method": "GET",
//	    "http.route": "/api/users",
//	    "server.address": "localhost:8080",
//	    "http.response.latency": "125ms"
//	}
package logger
