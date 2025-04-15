package middleware

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestIncomingLogger(t *testing.T) {
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
			// Create a logger middleware with the test configuration
			logger := IncomingLogger()

			if tt.logRequestBody {
				logger = logger.IncomingWithRequestBody()
			}

			if tt.logResponseBody {
				logger = logger.IncomingWithResponseBody()
			}

			if tt.logHeaders {
				logger = logger.IncomingWithHeaders()
			}

			if tt.customFields != nil {
				logger = logger.WithFields(tt.customFields)
			}

			// Create a test handler that will be wrapped
			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Verify trace ID has been set in context
				traceID := TraceIDFromContext(r.Context())
				if traceID == "" {
					t.Error("No trace ID set in context")
				}

				// If request body logging is enabled, check that the body is still readable
				if tt.logRequestBody && tt.requestBody != "" {
					body, err := io.ReadAll(r.Body)
					if err != nil {
						t.Errorf("Failed to read request body: %v", err)
					}

					if string(body) != tt.requestBody {
						t.Errorf("Request body was modified: got %v, want %v", string(body), tt.requestBody)
					}
				}

				// Set response status and body
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.responseBody))
			})

			// Create a request with the test configuration
			var reqBody io.Reader
			if tt.requestBody != "" {
				reqBody = strings.NewReader(tt.requestBody)
			}

			req := httptest.NewRequest(tt.requestMethod, "/test", reqBody)

			// Add headers if specified
			if tt.requestHeaders != nil {
				for k, v := range tt.requestHeaders {
					req.Header.Set(k, v)
				}
			}

			// Add trace ID if specified
			if tt.traceID != "" {
				req.Header.Set("X-Trace-ID", tt.traceID)
			}

			// Create a response recorder
			rr := httptest.NewRecorder()

			// Execute the middleware chain
			handler := logger.Handler(testHandler)
			handler.ServeHTTP(rr, req)

			// Check status code
			if rr.Code != tt.statusCode {
				t.Errorf("Handler returned wrong status code: got %v, want %v", rr.Code, tt.statusCode)
			}

			// Check response body
			if rr.Body.String() != tt.responseBody {
				t.Errorf("Handler returned unexpected body: got %v, want %v", rr.Body.String(), tt.responseBody)
			}

			// Check trace ID header in response
			traceIDHeader := rr.Header().Get("X-Trace-ID")
			if traceIDHeader == "" {
				t.Error("No X-Trace-ID header in response")
			}

			// If a custom trace ID was provided, check it was preserved
			if tt.traceID != "" && traceIDHeader != tt.traceID {
				t.Errorf("Trace ID not preserved: got %v, want %v", traceIDHeader, tt.traceID)
			}
		})
	}
}

func TestResponseWriter(t *testing.T) {
	t.Run("captures response body when enabled", func(t *testing.T) {
		rw := &responseWriter{
			body:         httptest.NewRecorder(),
			captureBody:  true,
			responseBody: &bytes.Buffer{},
		}

		testBody := []byte("test response")
		n, err := rw.Write(testBody)

		if err != nil {
			t.Errorf("Write returned error: %v", err)
		}

		if n != len(testBody) {
			t.Errorf("Write returned wrong length: got %v, want %v", n, len(testBody))
		}

		if rw.statusCode != http.StatusOK {
			t.Errorf("Status code not set to default: got %v, want %v", rw.statusCode, http.StatusOK)
		}

		if rw.responseBody.String() != string(testBody) {
			t.Errorf("Response body not captured correctly: got %v, want %v", rw.responseBody.String(), string(testBody))
		}
	})

	t.Run("doesn't capture response body when disabled", func(t *testing.T) {
		rw := &responseWriter{
			body:         httptest.NewRecorder(),
			captureBody:  false,
			responseBody: &bytes.Buffer{},
		}

		testBody := []byte("test response")
		rw.Write(testBody)

		if rw.responseBody.Len() > 0 {
			t.Errorf("Response body captured when disabled: %v", rw.responseBody.String())
		}
	})

	t.Run("passes header calls to underlying ResponseWriter", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		rw := &responseWriter{
			body: recorder,
		}

		rw.Header().Set("X-Test", "value")

		if recorder.Header().Get("X-Test") != "value" {
			t.Error("Header not set on underlying ResponseWriter")
		}
	})

	t.Run("WriteHeader sets status code", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		rw := &responseWriter{
			body: recorder,
		}

		rw.WriteHeader(http.StatusCreated)

		if rw.statusCode != http.StatusCreated {
			t.Errorf("Status code not set: got %v, want %v", rw.statusCode, http.StatusCreated)
		}

		if recorder.Code != http.StatusCreated {
			t.Errorf("Status code not set on underlying ResponseWriter: got %v, want %v", recorder.Code, http.StatusCreated)
		}
	})
}
