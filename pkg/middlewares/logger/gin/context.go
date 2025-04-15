package middleware

import (
	"github.com/CloudLearnersOrg/golib/pkg/logger"
	"github.com/gin-gonic/gin"
)

const traceIDKey = "trace.id"

// ContextWithTraceID adds a trace ID to the Gin context
func ContextWithTraceID(c *gin.Context, traceID string) {
	if c == nil || traceID == "" {
		logger.Warnf("ContextWithTraceID: context or trace ID is nil or empty", nil)
		return
	}

	c.Set(traceIDKey, traceID)
}

// TraceIDFromContext retrieves the trace ID from the Gin context
func TraceIDFromContext(c *gin.Context) string {
	if c == nil {
		return ""
	}
	if traceID, exists := c.Get(traceIDKey); exists {
		if id, ok := traceID.(string); ok {
			return id
		}
	}
	return ""
}
