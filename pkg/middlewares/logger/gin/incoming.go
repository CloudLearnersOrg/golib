package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/CloudLearnersOrg/golib/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// IncomingLogger returns a new logger middleware with default configuration
func IncomingLogger() *Middleware {
	return &Middleware{
		logRequestBody:  false,
		logResponseBody: false,
		logHeaders:      false,
	}
}

// Handler implements the Gin middleware functionality
func (m *Middleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Get or generate trace ID
		traceID := c.GetHeader("X-Trace-ID")
		if traceID == "" {
			traceID = uuid.New().String()
		}

		// Add trace ID to the context
		ContextWithTraceID(c, traceID)

		// Add trace ID to response headers
		c.Header("X-Trace-ID", traceID)

		// Capture request body if enabled
		var requestBody string
		if m.logRequestBody && c.Request.Body != nil {
			bodyBytes, err := io.ReadAll(c.Request.Body)
			if err != nil {
				if err := c.Error(err); err != nil {
					// Log the error using our logger package
					logger.Warnf("Failed to set error in gin context", map[string]any{
						"error": err.Error(),
					})

					return
				}

				logger.Warnf("Failed to read request body: %s", map[string]any{
					"error": err.Error(),
				})
				return
			}

			requestBody = string(bodyBytes)
			// Restore the body for further processing
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		// Create response body buffer if enabled
		var responseBodyBuffer *bytes.Buffer
		if m.logResponseBody {
			responseBodyBuffer = &bytes.Buffer{}
			blw := &bodyLogWriter{body: responseBodyBuffer, ResponseWriter: c.Writer}
			c.Writer = blw
		}

		// Process request
		c.Next()

		// Calculate duration
		duration := time.Since(start)

		// Create structured log entry
		logEntry := map[string]interface{}{
			"timestamp":                 time.Now().Format(time.RFC3339),
			"trace.id":                  traceID,
			"http.request.method":       c.Request.Method,
			"http.route":                c.FullPath(),
			"server.address":            c.Request.Host,
			"http.response.status_code": c.Writer.Status(),
			"http.response.latency":     duration.String(),
		}

		// Add request body if captured
		if m.logRequestBody && requestBody != "" {
			var prettyJSON bytes.Buffer
			if err := json.Indent(&prettyJSON, []byte(requestBody), "", "  "); err == nil {
				logEntry["http.request.body"] = prettyJSON.String()
			} else {
				logEntry["http.request.body"] = requestBody
			}
		}

		// Add response body if captured
		if m.logResponseBody && responseBodyBuffer != nil {
			respBody := responseBodyBuffer.String()
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
			for name, values := range c.Request.Header {
				headers[name] = strings.Join(values, ", ")
			}
			logEntry["http.request.headers"] = headers
		}

		// Add custom fields
		for _, field := range m.fields {
			logEntry[field.Key] = field.Value
		}

		// Add error if present
		if len(c.Errors) > 0 {
			logEntry["error"] = c.Errors.String()
		}

		// For now, just print to stdout
		// In a real implementation, you'd use a proper logger
		logJSON, _ := json.MarshalIndent(logEntry, "", "  ")
		fmt.Printf("incoming request: %s\n", logJSON)
	}
}

// bodyLogWriter is a custom gin.ResponseWriter that captures the response body
type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
