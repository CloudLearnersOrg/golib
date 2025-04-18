# http

The `http` package is written to make wider usage for Gin requests and record incoming and outgoing requests in JSON format.

## Features
- Automatic trace ID generation and propagation
- JSON formatted logging for incoming and outgoing requests
- Request timing and latency tracking
- HTTP status code based log levels
- Gin middleware integration

## Installation
```bash
go get github.com/CloudLearnersOrg/golib/pkg/http
```

## Example Usage

### Basic Setup
```go
package main

import (
    "log/slog"
    "os"

    "github.com/gin-gonic/gin"
    http "github.com/CloudLearnersOrg/golib/pkg/http"
    logger "github.com/CloudLearnersOrg/golib/pkg/middlewares/logger"
)

func main() {
    // Configure structured logging
    logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
    slog.SetDefault(logger)

    // Create Gin router with the middleware
    r := gin.New()
    r.Use(logger.Middleware())

    // Create HTTP client
    client := http.NewClient(nil)

    // Define your routes
    r.GET("/example", func(c *gin.Context) {
        // Make outgoing request
        resp, err := client.OutgoingRequest(
            c,
            "GET",
            "https://api.example.com/data",
            nil,
        )
        if err != nil {
            c.JSON(500, gin.H{"error": err.Error()})
            return
        }
        defer resp.Body.Close()

        // Forward response status
        c.Status(resp.StatusCode)
    })

    r.Run(":8080")
}
```

### Making Different Types of Requests
```go
// GET request
resp, err := client.OutgoingRequest(c, "GET", "https://api.example.com", nil)

// POST request with body
body := strings.NewReader(`{"key": "value"}`)
resp, err := client.OutgoingRequest(c, "POST", "https://api.example.com", body)

// PUT request
resp, err := client.OutgoingRequest(c, "PUT", "https://api.example.com", body)
```

### Example Log Output
```json
{
    "level": "INFO",
    "time": "2024-04-16T10:00:00Z",
    "message": "incoming request completed",
    "trace_id": "550e8400-e29b-41d4-a716-446655440000",
    "http.request.method": "GET",
    "http.route": "/example",
    "server.address": "localhost:8080",
    "http.response.status_code": 200,
    "http.response.latency": "150ms"
}
```

## Request Tracing
The package automatically:
1. Generates a trace ID for new requests
2. Propagates existing trace IDs through the `X-Trace-ID` header
3. Includes trace IDs in all log entries
4. Maintains trace context across service boundaries

## Error Handling
Log levels are automatically set based on HTTP status codes:
- 2xx: INFO level
- 4xx: WARN level
- 5xx: ERROR level

## Custom Client Configuration
```go
customClient := &http.Client{
    Timeout: 30 * time.Second,
    // Add other configurations
}
client := http.NewClient(customClient)
```