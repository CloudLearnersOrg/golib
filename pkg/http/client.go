package http

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Client wraps http.Client with tracing and logging capabilities
type Client struct {
	*http.Client
}

// NewClient creates a new HTTP client with tracing and logging middleware
func NewClient(baseClient *http.Client) *Client {
	if baseClient == nil {
		baseClient = http.DefaultClient
	}

	baseClient.Transport = newLoggingRoundTripper(baseClient.Transport)
	return &Client{Client: baseClient}
}

// OutgoingRequest performs an outgoing HTTP request with "outgoing" attributes to be logged.
func (c *Client) OutgoingRequest(ctx *gin.Context, method, url string, body io.Reader, headers map[string]string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(
		context.WithValue(ctx.Request.Context(), ginContextKey, ctx),
		method,
		url,
		body,
	)
	if err != nil {
		return nil, err
	}

	// Set headers if provided
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	return c.Do(req)
}

type roundTripperFunc func(*http.Request) (*http.Response, error)

// Used for http.RoundTripper interface
func (f roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

// NewLoggingRoundTripper creates a new round tripper that logs outgoing requests
func newLoggingRoundTripper(next http.RoundTripper) http.RoundTripper {
	if next == nil {
		next = http.DefaultTransport
	}

	return roundTripperFunc(func(req *http.Request) (*http.Response, error) {
		start := time.Now()
		traceID := extractTraceID(req.Context())

		if traceID != "" {
			req.Header.Set("X-Trace-ID", traceID)
		}

		resp, err := next.RoundTrip(req)
		return outgoing(req, err, resp, traceID, time.Since(start))
	})
}

func extractTraceID(ctx context.Context) string {
	if c, exists := ctx.Value(ginContextKey).(*gin.Context); exists {
		if id, exists := c.Get("X-Trace-ID"); exists {
			if traceId, ok := id.(string); ok {
				return traceId
			}
		}
	}

	// If not found, try to get from the parent context
	if parent, ok := ctx.(*gin.Context); ok {
		if id, exists := parent.Get("X-Trace-ID"); exists {
			if traceId, ok := id.(string); ok {
				return traceId
			}
		}
	}

	// Last resort: check if it's in request headers
	if gctx, ok := ctx.(*gin.Context); ok {
		return gctx.GetHeader("X-Trace-ID")
	}

	return ""
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
		if resp != nil && resp.StatusCode >= 400 {
			attrs = append(attrs, "http.response.status_code", resp.StatusCode)
		}

		slog.Error("outgoing request failed", append(attrs, "error", err)...)
		return resp, err
	}

	attrs = append(attrs, "http.response.status_code", resp.StatusCode)
	slog.Info("outgoing request completed", attrs...)
	return resp, nil
}
