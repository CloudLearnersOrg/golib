package middleware

import (
	"bytes"
	"net/http"
)

type responseWriter struct {
	body         http.ResponseWriter
	statusCode   int
	responseBody *bytes.Buffer
	captureBody  bool
}

// LogField represents a custom log field with a key and value
type logField struct {
	Key   string
	Value interface{}
}
