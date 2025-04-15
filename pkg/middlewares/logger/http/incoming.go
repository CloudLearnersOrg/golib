package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/CloudLearnersOrg/golib/pkg/logger"
	"github.com/google/uuid"
)

// Middleware holds logger configuration
type Middleware struct {
	fields          []logField
	logRequestBody  bool
	logResponseBody bool
	logHeaders      bool
}

// IncomingLogger returns a new logger middleware with default configuration
func IncomingLogger() *Middleware {
	return &Middleware{
		logRequestBody:  false,
		logResponseBody: false,
		logHeaders:      false,
	}
}

// Handler implements the middleware functionality
func (m *Middleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(body http.ResponseWriter, req *http.Request) {
		// Get the request context which will be canceled when the client disconnects
		ctx := req.Context()

		start := time.Now()

		// Generate or get trace ID
		traceID := req.Header.Get("X-Trace-ID")
		if traceID == "" {
			traceID = uuid.New().String()
		}

		// Add trace ID to the request context and update the request
		ctx = ContextWithTraceID(ctx, traceID)
		req = req.WithContext(ctx) // Correctly update the request with the new context

		// Add trace ID to response headers
		body.Header().Set("X-Trace-ID", traceID)

		// Create our custom response writer
		wrappedWriter := &responseWriter{
			body:         body,
			statusCode:   0,
			responseBody: &bytes.Buffer{},
			captureBody:  m.logResponseBody,
		}

		// Capture request body if enabled
		var requestBody string
		if m.logRequestBody && req.Body != nil {
			bodyBytes, err := io.ReadAll(req.Body)
			if err != nil {
				logger.Warnf("Failed to read request body: %s", map[string]any{
					"error": err.Error(),
				})
				return
			}

			requestBody = string(bodyBytes)
			// Restore the body for further processing
			req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		// Process the request with the updated req containing the trace ID context
		next.ServeHTTP(wrappedWriter, req)

		// If no status code was set, assume 200 OK
		statusCode := wrappedWriter.statusCode
		if statusCode == 0 {
			statusCode = http.StatusOK
		}

		// Calculate duration
		duration := time.Since(start)

		// Create structured log entry
		logEntry := map[string]interface{}{
			"timestamp":                 time.Now().Format(time.RFC3339),
			"trace.id":                  traceID,
			"http.request.method":       req.Method,
			"http.route":                req.URL.Path,
			"server.address":            req.Host,
			"http.response.status_code": statusCode,
			"http.response.latency":     duration.String(),
		}

		// Add request body if captured
		if m.logRequestBody && requestBody != "" {
			// Try to pretty format JSON
			var prettyJSON bytes.Buffer
			if err := json.Indent(&prettyJSON, []byte(requestBody), "", "  "); err == nil {
				logEntry["http.request.body"] = prettyJSON.String()
			} else {
				logEntry["http.request.body"] = requestBody
			}
		}

		// Add response body if captured
		if m.logResponseBody && wrappedWriter.responseBody.Len() > 0 {
			respBody := wrappedWriter.responseBody.String()
			var prettyJSON bytes.Buffer
			if err := json.Indent(&prettyJSON, []byte(respBody), "", "  "); err == nil {
				logEntry["http.response.body"] = prettyJSON.String()
			} else {
				logEntry["http.response.body"] = respBody
			}
		}

		// Add headers if enabled
		if m.logHeaders {
			headers := make(map[string]string)
			for name, values := range req.Header {
				headers[name] = strings.Join(values, ", ")
			}
			logEntry["http.request.headers"] = headers
		}

		// Add custom fields
		for _, field := range m.fields {
			logEntry[field.Key] = field.Value
		}

		// For now, just print to stdout
		// In a real implementation, you'd use a proper logger
		logJSON, err := json.MarshalIndent(logEntry, "", "  ")
		if err != nil {
			logger.Warnf("Failed to marshal log entry: %s", map[string]any{
				"error": err.Error(),
			})

			return
		}

		logger.Infof("incoming request: %s\n", map[string]any{
			"log": string(logJSON),
		})
	})
}

// ServeHTTP implements the http.Handler interface for compatibility with different routers
func (m *Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.Handler) {
	m.Handler(next).ServeHTTP(w, r)
}

func (rw *responseWriter) Header() http.Header {
	return rw.body.Header()
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	// If WriteHeader hasn't been called yet, calling Write implicitly sets the status code to 200
	if rw.statusCode == 0 {
		rw.statusCode = http.StatusOK
	}

	// Capture the response body if enabled
	if rw.captureBody {
		rw.responseBody.Write(b)
	}

	return rw.body.Write(b)
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.body.WriteHeader(statusCode)
}
