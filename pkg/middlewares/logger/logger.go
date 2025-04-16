package logger

import (
	"log/slog"
	"time"

	"github.com/CloudLearnersOrg/golib/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Middleware returns a Gin middleware for incoming request logging
func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Get or generate trace ID
		traceID := c.GetHeader("X-Trace-ID")
		if traceID == "" {
			traceID = uuid.New().String()
		}

		// Set trace ID in context for further use
		c.Set("X-Trace-ID", traceID)
		c.Header("X-Trace-ID", traceID)

		// Process request
		c.Next()

		incoming(c, traceID, time.Since(start))
	}
}

func incoming(c *gin.Context, traceID string, duration time.Duration) {
	attrs := []any{
		"trace_id", traceID,
		"http.response.status_code", c.Writer.Status(),
		"http.request.method", c.Request.Method,
		"http.route", c.Request.URL.Path,
		"server.address", c.Request.Host,
		"http.response.latency", duration.String(),
	}

	logger.LogFilteredStatusCode(c.Writer.Status(), "incoming request", attrs...)
	if c.Writer.Status() >= 400 {
		slog.Error("incoming request failed", attrs...)
		return
	}

	slog.Info("incoming request completed", attrs...)
}
