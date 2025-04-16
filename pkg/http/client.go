package httpclient

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const TraceIDKey = "X-Trace-ID"

// Client wraps http.Client with tracing and logging capabilities
type Client struct {
	*http.Client
}

// NewClient creates a new HTTP client with tracing and logging middleware
func NewClient(baseClient *http.Client) *Client {
	if baseClient == nil {
		baseClient = http.DefaultClient
	}

	baseClient.Transport = NewLoggingRoundTripper(baseClient.Transport)
	return &Client{Client: baseClient}
}

// OutgoingRequest performs an outgoing HTTP request with tracing
func (c *Client) OutgoingRequest(ginCtx *gin.Context, method, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(
		context.WithValue(context.Background(), "GinContextKey", ginCtx),
		method,
		url,
		body,
	)
	if err != nil {
		return nil, err
	}

	return c.Client.Do(req)
}

// Middleware returns a Gin middleware for incoming request logging
func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Get or generate trace ID
		traceID := c.GetHeader(TraceIDKey)
		if traceID == "" {
			traceID = uuid.New().String()
		}

		// Set trace ID in context for further use
		c.Set(TraceIDKey, traceID)
		c.Header(TraceIDKey, traceID)

		// Process request
		c.Next()

		incoming(c, traceID, time.Since(start))
	}
}

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

// NewLoggingRoundTripper creates a new round tripper that logs outgoing requests
func NewLoggingRoundTripper(next http.RoundTripper) http.RoundTripper {
	if next == nil {
		next = http.DefaultTransport
	}

	return roundTripperFunc(func(req *http.Request) (*http.Response, error) {
		start := time.Now()
		traceID := extractTraceID(req.Context())

		if traceID != "" {
			req.Header.Set(TraceIDKey, traceID)
		}

		resp, err := next.RoundTrip(req)
		return outgoing(req, err, resp, traceID, time.Since(start))
	})
}

func extractTraceID(ctx context.Context) string {
	if ginCtx, exists := ctx.Value("GinContextKey").(*gin.Context); exists {
		if id, exists := ginCtx.Get(TraceIDKey); exists {
			return id.(string)
		}
	}
	return ""
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

	log(c.Writer.Status(), "incoming request", attrs...)
	if c.Writer.Status() >= 400 {
		slog.Error("incoming request failed", attrs...)
		return
	}

	slog.Info("incoming request completed", attrs...)
}

func outgoing(req *http.Request, err error, resp *http.Response, traceID string, duration time.Duration) (*http.Response, error) {
	attrs := []any{
		"trace_id", traceID,
		"http.request.method", req.Method,
		"http.route", req.URL.Path,
		"server.address", req.URL.Host,
		"http.response.latency", duration.String(),
	}

	if err != nil {
		slog.Error("outgoing request failed", append(attrs, "error", err)...)
		return resp, err
	}

	attrs = append(attrs, "http.response.status_code", resp.StatusCode)
	log(resp.StatusCode, "outgoing request", attrs...)
	if resp.StatusCode >= 400 {
		slog.Error("outgoing request failed", attrs...)
		return resp, err
	}

	slog.Info("outgoing request completed", attrs...)
	return resp, nil
}

func log(status int, prefix string, attrs ...any) {
	switch {
	case status >= 500:
		slog.Error(prefix+" failed", attrs...)
	case status >= 400:
		slog.Warn(prefix+" warning", attrs...)
	default:
		slog.Info(prefix+" completed", attrs...)
	}
}
