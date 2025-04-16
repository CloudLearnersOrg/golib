package httpclient

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	// Given
	var baseClient *http.Client

	// When
	client := NewClient(baseClient)

	// Then
	assert.NotNil(t, client)
	assert.NotNil(t, client.Transport)
}

func TestOutgoingGetRequest(t *testing.T) {
	// Given
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		_, err := w.Write([]byte("test response"))
		if err != nil {
			t.Errorf("failed to write response: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}))
	defer server.Close()

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/test", nil)
	client := NewClient(nil)

	// When
	resp, err := client.OutgoingRequest(ctx, http.MethodGet, server.URL, nil)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestOutgoingPostRequest(t *testing.T) {
	// Given
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("failed to read request body: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		assert.Equal(t, "test body", string(body))

		_, err = w.Write([]byte("test response"))
		if err != nil {
			t.Errorf("failed to write response: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}))
	defer server.Close()

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/test", nil)
	client := NewClient(nil)

	// When
	resp, err := client.OutgoingRequest(ctx, http.MethodPost, server.URL, strings.NewReader("test body"))

	// Then
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

// ...rest of the file remains unchanged...
