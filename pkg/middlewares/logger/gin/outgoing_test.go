package middleware

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func TestOutgoingLogger(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name            string
		requestBody     string
		responseBody    string
		requestMethod   string
		requestHeaders  map[string]string
		logRequestBody  bool
		logResponseBody bool
		logHeaders      bool
		customFields    map[string]interface{}
		traceID         string
		statusCode      int
	}{
		{
			name:          "basic request with no optional features",
			requestBody:   "",
			responseBody:  "Hello, World!",
			requestMethod: http.MethodGet,
			statusCode:    http.StatusOK,
		},
		{
			name:           "request with body logging",
			requestBody:    `{"test":"value"}`,
			responseBody:   "Response",
			requestMethod:  http.MethodPost,
			logRequestBody: true,
			statusCode:     http.StatusCreated,
		},
		{
			name:            "request with response body logging",
			requestBody:     "",
			responseBody:    `{"response":"data"}`,
			requestMethod:   http.MethodGet,
			logResponseBody: true,
			statusCode:      http.StatusOK,
		},
		{
			name:           "request with headers logging",
			requestBody:    "",
			responseBody:   "Response",
			requestMethod:  http.MethodGet,
			requestHeaders: map[string]string{"X-Custom-Header": "value"},
			logHeaders:     true,
			statusCode:     http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Check if trace ID header is present
				if r.Header.Get("X-Trace-ID") == "" {
					t.Error("No X-Trace-ID header in outgoing request")
				}

				// Write the response
				w.WriteHeader(tt.statusCode)

				n, err := w.Write([]byte(tt.responseBody))
				if err != nil {
					t.Errorf("Failed to write response body: %v", err)
				}

				if n != len(tt.responseBody) {
					t.Errorf("Failed to write complete response body: wrote %d bytes, expected %d bytes", n, len(tt.responseBody))
				}
			}))
			defer server.Close()

			// Create Gin context with default trace ID if not provided
			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			if tt.traceID != "" {
				ContextWithTraceID(c, tt.traceID)
			} else {
				// Generate a default trace ID when not provided
				ContextWithTraceID(c, uuid.New().String())
			}

			// Create logger
			logger := NewOutgoingLogger(server.Client())

			if tt.logRequestBody {
				logger = logger.OutgoingWithRequestBody()
			}
			if tt.logResponseBody {
				logger = logger.OutgoingWithResponseBody()
			}
			if tt.logHeaders {
				logger = logger.OutgoingWithHeaders()
			}
			if tt.customFields != nil {
				logger = logger.WithFields(tt.customFields)
			}

			// Create request
			var reqBody io.Reader
			if tt.requestBody != "" {
				reqBody = bytes.NewBufferString(tt.requestBody)
			}

			req, err := http.NewRequest(tt.requestMethod, server.URL, reqBody)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			// Add headers if specified
			if tt.requestHeaders != nil {
				for k, v := range tt.requestHeaders {
					req.Header.Set(k, v)
				}
			}

			// Perform request
			resp, err := logger.Do(c, req)
			if err != nil {
				t.Fatalf("Request failed: %v", err)
			}
			defer resp.Body.Close()

			// Check response
			if resp.StatusCode != tt.statusCode {
				t.Errorf("Wrong status code: got %v, want %v", resp.StatusCode, tt.statusCode)
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Failed to read response body: %v", err)
			}

			if string(body) != tt.responseBody {
				t.Errorf("Wrong response body: got %v, want %v", string(body), tt.responseBody)
			}
		})
	}
}

func TestOutgoingLoggerWithFields(t *testing.T) {
	logger := NewOutgoingLogger(nil)

	// Test WithField
	logger = logger.WithField("key1", "value1")
	if len(logger.fields) != 1 {
		t.Errorf("WithField failed: got %d fields, want 1", len(logger.fields))
	}

	// Test WithFields
	fields := map[string]interface{}{
		"key2": "value2",
		"key3": 123,
	}
	logger = logger.WithFields(fields)
	if len(logger.fields) != 3 {
		t.Errorf("WithFields failed: got %d fields, want 3", len(logger.fields))
	}
}
