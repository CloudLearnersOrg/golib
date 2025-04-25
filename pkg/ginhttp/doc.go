// Package ginhttp provides a custom HTTP client with tracing and logging capabilities for Gin web applications.
// It wraps the standard net/http client and integrates with Gin web framework for
// context propagation and trace ID handling.
//
// The package implements the following main features:
//   - Automatic trace ID propagation across service boundaries
//   - Request/Response logging with structured attributes
//   - Integration with slog for structured logging
//   - Compatible with Gin web framework contexts
//
// Basic usage:
//
//	client := ginhttp.NewClient(nil)
//	resp, err := client.OutgoingRequest(
//		ginCtx,
//		"GET",
//		"https://api.example.com",
//		nil,
//		map[string]string{"Authorization": "Bearer token"},
//	)
//
// Trace ID Propagation:
// The package automatically extracts trace IDs from incoming requests and propagates
// them to outgoing requests through the X-Trace-ID header. The trace ID is extracted
// in the following order:
//  1. From Gin context stored value
//  2. From parent context
//  3. From request headers
//
// Logging:
// All outgoing requests are automatically logged with the following attributes:
//   - trace_id: The propagated trace ID
//   - http.request.method: The HTTP method
//   - http.route: The request path
//   - server.address: The target host
//   - http.response.latency: Request duration
//   - http.response.status_code: Response status code
//   - error: Error message (if request failed)
package ginhttp
