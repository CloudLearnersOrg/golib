package httpclient

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/CloudLearnersOrg/golib/pkg/logger"
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

// OutgoingRequest performs an outgoing HTTP request with tracing
func (c *Client) OutgoingRequest(ctx *gin.Context, method, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(
		context.WithValue(context.Background(), ginContextKey, ctx),
		method,
		url,
		body,
	)
	if err != nil {
		return nil, err
	}

	return c.Do(req)
}

type roundTripperFunc func(*http.Request) (*http.Response, error)

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
			return id.(string)
		}
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
		slog.Error("outgoing request failed", append(attrs, "error", err)...)
		return resp, err
	}

	attrs = append(attrs, "http.response.status_code", resp.StatusCode)
	logger.LogFilteredStatusCode(resp.StatusCode, "outgoing request", attrs...)
	if resp.StatusCode >= 400 {
		slog.Error("outgoing request failed", attrs...)
		return resp, err
	}

	slog.Info("outgoing request completed", attrs...)
	return resp, nil
}
