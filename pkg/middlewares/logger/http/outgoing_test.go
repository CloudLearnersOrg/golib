package middleware

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
)

func TestOutgoingLogger(t *testing.T) {
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
		{
			name:          "request with custom trace ID",
			requestBody:   "",
			responseBody:  "Response",
			requestMethod: http.MethodGet,
			traceID:       "custom-trace-id",
			statusCode:    http.StatusOK,
		},
		{
			name:          "request with custom fields",
			requestBody:   "",
			responseBody:  "Response",
			requestMethod: http.MethodGet,
			customFields:  map[string]interface{}{"custom": "field", "numeric": 123},
			statusCode:    http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test server to handle the request
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Check if trace ID header is set
				if r.Header.Get("X-Trace-ID") == "" {
					t.Error("No X-Trace-ID header in outgoing request")
				}

				// If a specific trace ID was provided, check it was used
				if tt.traceID != "" && r.Header.Get("X-Trace-ID") != tt.traceID {
					t.Errorf("Wrong trace ID: got %v, want %v", r.Header.Get("X-Trace-ID"), tt.traceID)
				}

				// Read request body if it has content
				if tt.requestBody != "" {
					body, err := io.ReadAll(r.Body)
					if err != nil {
						t.Errorf("Failed to read request body: %v", err)
					}

					if string(body) != tt.requestBody {
						t.Errorf("Request body incorrect: got %v, want %v", string(body), tt.requestBody)
					}
				}

				// Set response status and body
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.responseBody))
			}))
			defer server.Close()

			// Create an HTTP client
			client := server.Client()

			// Create the outgoing logger
			logger := NewOutgoingLogger(client)

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

			// Create the request
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

			// Add trace ID to context or use default
			var ctx context.Context
			if tt.traceID != "" {
				ctx = ContextWithTraceID(context.Background(), tt.traceID)
			} else {
				// Set a default trace ID for tests that don't specify one
				defaultTraceID := uuid.New().String()
				ctx = ContextWithTraceID(context.Background(), defaultTraceID)
			}
			req = req.WithContext(ctx)

			// Make the request
			resp, err := logger.Do(req)
			if err != nil {
				t.Fatalf("Request failed: %v", err)
			}
			defer resp.Body.Close()

			// Check response status code
			if resp.StatusCode != tt.statusCode {
				t.Errorf("Wrong status code: got %v, want %v", resp.StatusCode, tt.statusCode)
			}

			// Check response body
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Failed to read response body: %v", err)
			}

			if string(body) != tt.responseBody {
				t.Errorf("Response body incorrect: got %v, want %v", string(body), tt.responseBody)
			}
		})
	}
}

func TestOutgoingLoggerWithField(t *testing.T) {
	logger := NewOutgoingLogger(nil)

	// Add a single field
	logger = logger.WithField("key", "value")

	if len(logger.fields) != 1 {
		t.Errorf("WithField did not add a field: got %v fields", len(logger.fields))
	}

	if logger.fields[0].Key != "key" || logger.fields[0].Value != "value" {
		t.Errorf("WithField added incorrect field: got %+v", logger.fields[0])
	}
}

func TestOutgoingLoggerWithFields(t *testing.T) {
	logger := NewOutgoingLogger(nil)

	// Add multiple fields
	fields := map[string]interface{}{
		"key1": "value1",
		"key2": 123,
	}

	logger = logger.WithFields(fields)

	if len(logger.fields) != 2 {
		t.Errorf("WithFields did not add correct number of fields: got %v", len(logger.fields))
	}

	// Check that both fields were added (order not guaranteed)
	foundKey1 := false
	foundKey2 := false

	for _, field := range logger.fields {
		if field.Key == "key1" && field.Value == "value1" {
			foundKey1 = true
		}
		if field.Key == "key2" && field.Value == 123 {
			foundKey2 = true
		}
	}

	if !foundKey1 || !foundKey2 {
		t.Errorf("WithFields did not add correct fields: %+v", logger.fields)
	}
}
