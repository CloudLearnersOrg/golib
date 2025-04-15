package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// OutgoingLogger wraps http.Client to log outgoing requests
type OutgoingLogger struct {
	client          *http.Client
	logRequestBody  bool
	logResponseBody bool
	logHeaders      bool
	fields          []logField
}

// NewOutgoingLogger creates a new outgoing request logger
func NewOutgoingLogger(client *http.Client) *OutgoingLogger {
	if client == nil {
		client = http.DefaultClient
	}
	return &OutgoingLogger{
		client: client,
	}
}

// WithField adds a custom field to the logger
func (o *OutgoingLogger) WithField(key string, value interface{}) *OutgoingLogger {
	o.fields = append(o.fields, logField{Key: key, Value: value})
	return o
}

// WithFields adds multiple custom fields to the logger
func (o *OutgoingLogger) WithFields(fields map[string]interface{}) *OutgoingLogger {
	for k, v := range fields {
		o.fields = append(o.fields, logField{Key: k, Value: v})
	}
	return o
}

// Do performs an HTTP request and logs the details
func (o *OutgoingLogger) Do(req *http.Request) (*http.Response, error) {
	start := time.Now()

	// Get trace ID from context or generate a new one
	traceID := TraceIDFromContext(req.Context())
	if traceID == "" {
		traceID = req.Header.Get("X-Trace-ID")
	}

	// Add trace ID to outgoing request headers
	req.Header.Set("X-Trace-ID", traceID)

	// Capture request body if enabled
	var requestBody string
	if o.logRequestBody && req.Body != nil {
		bodyBytes, _ := io.ReadAll(req.Body)
		requestBody = string(bodyBytes)
		// Restore the body for the actual request
		req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	}

	// Execute the request
	resp, err := o.client.Do(req)

	// Calculate duration
	duration := time.Since(start)

	// Create structured log entry
	logEntry := map[string]interface{}{
		"timestamp":             time.Now().Format(time.RFC3339),
		"trace.id":              traceID,
		"http.request.method":   req.Method,
		"http.route":            req.URL.Path,
		"server.address":        req.URL.Host,
		"http.response.latency": duration.String(),
	}

	if err != nil {
		logEntry["error"] = err.Error()
	} else {
		logEntry["http.response.status_code"] = resp.StatusCode

		// Capture response body if enabled
		if o.logResponseBody && resp.Body != nil {
			bodyBytes, _ := ioutil.ReadAll(resp.Body)
			// Restore the body for the caller
			resp.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

			// Try to pretty format JSON
			var prettyJSON bytes.Buffer
			if err := json.Indent(&prettyJSON, bodyBytes, "", "  "); err == nil {
				logEntry["http.response.body"] = prettyJSON.String()
			} else {
				logEntry["http.response.body"] = string(bodyBytes)
			}
		}
	}

	// Add request body if captured
	if o.logRequestBody && requestBody != "" {
		var prettyJSON bytes.Buffer
		if err := json.Indent(&prettyJSON, []byte(requestBody), "", "  "); err == nil {
			logEntry["http.request.body"] = prettyJSON.String()
		} else {
			logEntry["http.request.body"] = requestBody
		}
	}

	// Add headers if enabled
	if o.logHeaders {
		reqHeaders := make(map[string]string)
		for name, values := range req.Header {
			reqHeaders[name] = strings.Join(values, ", ")
		}
		logEntry["http.request.headers"] = reqHeaders

		if err == nil {
			respHeaders := make(map[string]string)
			for name, values := range resp.Header {
				respHeaders[name] = strings.Join(values, ", ")
			}
			logEntry["http.response.headers"] = respHeaders
		}
	}

	// Add custom fields
	for _, field := range o.fields {
		logEntry[field.Key] = field.Value
	}

	// For now, just print to stdout
	// In a real implementation, you'd use a proper logger
	logJSON, _ := json.MarshalIndent(logEntry, "", "  ")
	fmt.Printf("outgoing request: %s\n", logJSON)

	return resp, err
}
