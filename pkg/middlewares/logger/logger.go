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
	return func(ctx *gin.Context) {
		start := time.Now()

		// Get or generate trace ID
		traceID := ctx.GetHeader("X-Trace-ID")
		if traceID == "" {
			traceID = uuid.New().String()
		}

		// Set trace ID in context for further use
		ctx.Set("X-Trace-ID", traceID)
		ctx.Header("X-Trace-ID", traceID)

		// Process request
		ctx.Next()

		incoming(ctx, traceID, time.Since(start))
	}
}

func incoming(ctx *gin.Context, traceID string, duration time.Duration) {
	attrs := []any{
		"trace_id", traceID,
		"http.response.status_code", ctx.Writer.Status(),
		"http.request.method", ctx.Request.Method,
		"http.route", ctx.Request.URL.Path,
		"server.address", ctx.Request.Host,
		"http.response.latency", duration.String(),
	}

	logger.LogFilteredStatusCode(ctx.Writer.Status(), "incoming request", attrs...)
	if ctx.Writer.Status() >= 400 {
		slog.Error("incoming request failed", attrs...)
		return
	}

	slog.Info("incoming request completed", attrs...)
}
