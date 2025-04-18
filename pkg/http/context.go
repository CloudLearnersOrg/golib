package http

// Gin Context Key is used to store and retrieve the trace ID from the context
/* graph LR
   A[Incoming Request] --> B[Gin Middleware]
   B -- Store TraceID --> C[Gin Context]
   C -- GinContextKey --> D[Outgoing Request Context]
   D -- Extract TraceID --> E[Downstream Service] */
type contextKey string

const (
	ginContextKey contextKey = "gin"
)
