package middleware

import (
	"bytes"
	"net/http/httptest"
	"testing"
)

func TestResponseWriterCreation(t *testing.T) {
	recorder := httptest.NewRecorder()
	buffer := &bytes.Buffer{}

	rw := &responseWriter{
		body:         recorder,
		responseBody: buffer,
		captureBody:  true,
	}

	if rw.body != recorder {
		t.Error("responseWriter.body not set correctly")
	}

	if rw.responseBody != buffer {
		t.Error("responseWriter.responseBody not set correctly")
	}

	if !rw.captureBody {
		t.Error("responseWriter.captureBody not set correctly")
	}

	if rw.statusCode != 0 {
		t.Errorf("responseWriter.statusCode should be 0 by default, got %v", rw.statusCode)
	}
}

func TestLogFieldCreation(t *testing.T) {
	field := logField{
		Key:   "testKey",
		Value: "testValue",
	}

	if field.Key != "testKey" {
		t.Errorf("logField.Key not set correctly: got %v, want %v", field.Key, "testKey")
	}

	if field.Value != "testValue" {
		t.Errorf("logField.Value not set correctly: got %v, want %v", field.Value, "testValue")
	}
}
