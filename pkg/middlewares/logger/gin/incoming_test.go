package middleware

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestIncomingLogger(t *testing.T) {
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
			// Create a new Gin engine
			r := gin.New()

			// Create logger middleware with the test configuration
			logger := IncomingLogger()

			if tt.logRequestBody {
				logger = logger.WithRequestBody()
			}

			if tt.logResponseBody {
				logger = logger.WithResponseBody()
			}

			if tt.logHeaders {
				logger = logger.WithHeaders()
			}

			if tt.customFields != nil {
				logger = logger.WithFields(tt.customFields)
			}

			// Add the logger middleware
			r.Use(logger.Handler())

			// Add a test handler
			r.Any("/test", func(c *gin.Context) {
				// Verify trace ID has been set in context
				traceID := TraceIDFromContext(c)
				if traceID == "" {
					t.Error("No trace ID set in context")
				}

				c.Status(tt.statusCode)
				c.Writer.Write([]byte(tt.responseBody))
			})

			// Create a request
			var reqBody io.Reader
			if tt.requestBody != "" {
				reqBody = bytes.NewBufferString(tt.requestBody)
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
			w := httptest.NewRecorder()

			// Serve the request
			r.ServeHTTP(w, req)

			// Check status code
			if w.Code != tt.statusCode {
				t.Errorf("Handler returned wrong status code: got %v, want %v", w.Code, tt.statusCode)
			}

			// Check response body
			if w.Body.String() != tt.responseBody {
				t.Errorf("Handler returned unexpected body: got %v, want %v", w.Body.String(), tt.responseBody)
			}

			// Check trace ID header in response
			traceIDHeader := w.Header().Get("X-Trace-ID")
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

func TestBodyLogWriter(t *testing.T) {
	t.Run("captures response body", func(t *testing.T) {
		buffer := &bytes.Buffer{}
		gin.SetMode(gin.TestMode)
		c, _ := gin.CreateTestContext(httptest.NewRecorder())

		writer := &bodyLogWriter{
			ResponseWriter: c.Writer,
			body:           buffer,
		}

		testBody := []byte("test response")
		n, err := writer.Write(testBody)

		if err != nil {
			t.Errorf("Write returned error: %v", err)
		}

		if n != len(testBody) {
			t.Errorf("Write returned wrong length: got %v, want %v", n, len(testBody))
		}

		if buffer.String() != string(testBody) {
			t.Errorf("Response body not captured correctly: got %v, want %v", buffer.String(), string(testBody))
		}
	})
}
