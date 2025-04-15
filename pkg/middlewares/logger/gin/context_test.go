package middleware

import (
	"testing"

	"github.com/gin-gonic/gin"
)

func TestContextWithTraceID(t *testing.T) {
	tests := []struct {
		name    string
		traceID string
	}{
		{
			name:    "with valid trace ID",
			traceID: "123456789",
		},
		{
			name:    "with empty trace ID",
			traceID: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := gin.CreateTestContext(nil)
			ContextWithTraceID(c, tt.traceID)

			retrievedID := TraceIDFromContext(c)
			if retrievedID != tt.traceID {
				t.Errorf("TraceIDFromContext() = %v, want %v", retrievedID, tt.traceID)
			}
		})
	}
}

func TestTraceIDFromContext(t *testing.T) {
	tests := []struct {
		name     string
		setup    func() *gin.Context
		expected string
	}{
		{
			name: "with trace ID in context",
			setup: func() *gin.Context {
				c, _ := gin.CreateTestContext(nil)
				c.Set(traceIDKey, "123456789")
				return c
			},
			expected: "123456789",
		},
		{
			name: "with nil context",
			setup: func() *gin.Context {
				return nil
			},
			expected: "",
		},
		{
			name: "with context but no trace ID",
			setup: func() *gin.Context {
				c, _ := gin.CreateTestContext(nil)
				return c
			},
			expected: "",
		},
		{
			name: "with non-string trace ID",
			setup: func() *gin.Context {
				c, _ := gin.CreateTestContext(nil)
				c.Set(traceIDKey, 12345)
				return c
			},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := tt.setup()
			got := TraceIDFromContext(ctx)
			if got != tt.expected {
				t.Errorf("TraceIDFromContext() = %v, want %v", got, tt.expected)
			}
		})
	}
}
