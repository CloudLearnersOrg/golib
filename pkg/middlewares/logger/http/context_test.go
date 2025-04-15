package middleware

import (
	"context"
	"testing"
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
			ctx := context.Background()
			newCtx := ContextWithTraceID(ctx, tt.traceID)

			if newCtx == nil {
				t.Fatal("ContextWithTraceID returned nil context")
			}

			retrievedID := TraceIDFromContext(newCtx)
			if retrievedID != tt.traceID {
				t.Errorf("TraceIDFromContext() = %v, want %v", retrievedID, tt.traceID)
			}
		})
	}
}

func TestTraceIDFromContext(t *testing.T) {
	tests := []struct {
		name     string
		ctx      context.Context
		expected string
	}{
		{
			name:     "with trace ID in context",
			ctx:      context.WithValue(context.Background(), traceIDKey, "123456789"),
			expected: "123456789",
		},
		{
			name:     "with nil context",
			ctx:      nil,
			expected: "",
		},
		{
			name:     "with context but no trace ID",
			ctx:      context.Background(),
			expected: "",
		},
		{
			name:     "with non-string trace ID",
			ctx:      context.WithValue(context.Background(), traceIDKey, 12345),
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TraceIDFromContext(tt.ctx)
			if got != tt.expected {
				t.Errorf("TraceIDFromContext() = %v, want %v", got, tt.expected)
			}
		})
	}
}
