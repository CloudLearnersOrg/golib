package middleware

// LogField represents a custom log field with a key and value
type logField struct {
	Key   string
	Value interface{}
}

// Middleware holds logger configuration
type Middleware struct {
	fields          []logField
	logRequestBody  bool
	logResponseBody bool
	logHeaders      bool
}
