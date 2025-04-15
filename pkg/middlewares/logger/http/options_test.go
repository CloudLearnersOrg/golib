package middleware

import (
	"testing"
)

func TestMiddlewareWithField(t *testing.T) {
	m := IncomingLogger()

	// Add a single field
	m = m.WithField("key", "value")

	if len(m.fields) != 1 {
		t.Errorf("WithField did not add a field: got %v fields", len(m.fields))
	}

	if m.fields[0].Key != "key" || m.fields[0].Value != "value" {
		t.Errorf("WithField added incorrect field: got %+v", m.fields[0])
	}
}

func TestMiddlewareWithFields(t *testing.T) {
	m := IncomingLogger()

	// Add multiple fields
	fields := map[string]interface{}{
		"key1": "value1",
		"key2": 123,
	}

	m = m.WithFields(fields)

	if len(m.fields) != 2 {
		t.Errorf("WithFields did not add correct number of fields: got %v", len(m.fields))
	}

	// Check that both fields were added (order not guaranteed)
	foundKey1 := false
	foundKey2 := false

	for _, field := range m.fields {
		if field.Key == "key1" && field.Value == "value1" {
			foundKey1 = true
		}
		if field.Key == "key2" && field.Value == 123 {
			foundKey2 = true
		}
	}

	if !foundKey1 || !foundKey2 {
		t.Errorf("WithFields did not add correct fields: %+v", m.fields)
	}
}

func TestMiddlewareWithRequestBody(t *testing.T) {
	m := IncomingLogger()

	// Default value should be false
	if m.logRequestBody {
		t.Error("Default logRequestBody should be false")
	}

	// After WithRequestBody, should be true
	m = m.IncomingWithRequestBody()

	if !m.logRequestBody {
		t.Error("WithRequestBody did not enable logRequestBody")
	}
}

func TestMiddlewareWithResponseBody(t *testing.T) {
	m := IncomingLogger()

	// Default value should be false
	if m.logResponseBody {
		t.Error("Default logResponseBody should be false")
	}

	// After WithResponseBody, should be true
	m = m.IncomingWithResponseBody()

	if !m.logResponseBody {
		t.Error("WithResponseBody did not enable logResponseBody")
	}
}

func TestMiddlewareWithHeaders(t *testing.T) {
	m := IncomingLogger()

	// Default value should be false
	if m.logHeaders {
		t.Error("Default logHeaders should be false")
	}

	// After WithHeaders, should be true
	m = m.IncomingWithHeaders()

	if !m.logHeaders {
		t.Error("WithHeaders did not enable logHeaders")
	}
}
