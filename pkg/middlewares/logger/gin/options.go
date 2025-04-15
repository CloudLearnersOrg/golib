package middleware

// WithField adds a custom field to the logger
func (m *Middleware) WithField(key string, value interface{}) *Middleware {
	m.fields = append(m.fields, logField{Key: key, Value: value})
	return m
}

// WithFields adds multiple custom fields to the logger
func (m *Middleware) WithFields(fields map[string]interface{}) *Middleware {
	for k, v := range fields {
		m.fields = append(m.fields, logField{Key: k, Value: v})
	}
	return m
}

// WithRequestBody enables request body logging
func (m *Middleware) WithRequestBody() *Middleware {
	m.logRequestBody = true
	return m
}

// WithResponseBody enables response body logging
func (m *Middleware) WithResponseBody() *Middleware {
	m.logResponseBody = true
	return m
}

// WithHeaders enables header logging
func (m *Middleware) WithHeaders() *Middleware {
	m.logHeaders = true
	return m
}
