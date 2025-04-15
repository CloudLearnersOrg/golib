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
func (m *Middleware) IncomingWithRequestBody() *Middleware {
	m.logRequestBody = true
	return m
}

// WithResponseBody enables response body logging
func (m *Middleware) IncomingWithResponseBody() *Middleware {
	m.logResponseBody = true
	return m
}

// WithHeaders enables header logging
func (m *Middleware) IncomingWithHeaders() *Middleware {
	m.logHeaders = true
	return m
}

// WithRequestBody enables request body logging
func (o *OutgoingLogger) OutgoingWithRequestBody() *OutgoingLogger {
	o.logRequestBody = true
	return o
}

// WithResponseBody enables response body logging
func (o *OutgoingLogger) OutgoingWithResponseBody() *OutgoingLogger {
	o.logResponseBody = true
	return o
}

// WithHeaders enables header logging
func (o *OutgoingLogger) OutgoingWithHeaders() *OutgoingLogger {
	o.logHeaders = true
	return o
}
