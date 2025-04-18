package http

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
		require.NoError(t, err, "failed to write response")
	}))
	defer server.Close()

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/test", nil)
	client := NewClient(nil)

	// When
	resp, err := client.OutgoingRequest(ctx, http.MethodGet, server.URL, nil, nil)
	require.NoError(t, err)
	defer func() {
		err := resp.Body.Close()
		require.NoError(t, err, "failed to close response body")
	}()

	// Then
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Check response body
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.Equal(t, "test response", string(body))
}

func TestOutgoingPostRequest(t *testing.T) {
	// Given
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		body, err := io.ReadAll(r.Body)
		require.NoError(t, err, "failed to read request body")
		assert.Equal(t, "test body", string(body))

		_, err = w.Write([]byte("test response"))
		require.NoError(t, err, "failed to write response")
	}))
	defer server.Close()

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/test", nil)
	client := NewClient(nil)

	// When
	resp, err := client.OutgoingRequest(ctx, http.MethodPost, server.URL, strings.NewReader("test body"), nil)
	require.NoError(t, err)
	defer func() {
		err := resp.Body.Close()
		require.NoError(t, err, "failed to close response body")
	}()

	// Then
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Check response body
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.Equal(t, "test response", string(body))
}

func TestOutgoingRequestWithErrorResponse(t *testing.T) {
	// Given
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte("internal server error"))
		require.NoError(t, err, "failed to write error response")
	}))
	defer server.Close()

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/test", nil)
	client := NewClient(nil)

	// When
	resp, err := client.OutgoingRequest(ctx, http.MethodGet, server.URL, nil, nil)
	require.NoError(t, err)
	defer func() {
		err := resp.Body.Close()
		require.NoError(t, err, "failed to close response body")
	}()

	// Then
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.Equal(t, "internal server error", string(body))
}

func TestOutgoingGetRequestWithHeaders(t *testing.T) {
	// Given
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Accept"))
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))

		_, err := w.Write([]byte("test response"))
		require.NoError(t, err, "failed to write response")
	}))
	defer server.Close()

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/test", nil)
	client := NewClient(nil)

	headers := map[string]string{
		"Accept":        "application/json",
		"Authorization": "Bearer test-token",
	}

	// When
	resp, err := client.OutgoingRequest(ctx, http.MethodGet, server.URL, nil, headers)
	require.NoError(t, err)
	defer func() {
		err := resp.Body.Close()
		require.NoError(t, err, "failed to close response body")
	}()

	// Then
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestOutgoingPostRequestWithHeaders(t *testing.T) {
	// Given
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "application/x-www-form-urlencoded", r.Header.Get("Content-Type"))
		assert.Equal(t, "application/json", r.Header.Get("Accept"))

		body, err := io.ReadAll(r.Body)
		require.NoError(t, err, "failed to read request body")
		assert.Equal(t, "test body", string(body))

		_, err = w.Write([]byte("test response"))
		require.NoError(t, err, "failed to write response")
	}))
	defer server.Close()

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/test", nil)
	client := NewClient(nil)

	headers := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
		"Accept":       "application/json",
	}

	// When
	resp, err := client.OutgoingRequest(ctx, http.MethodPost, server.URL, strings.NewReader("test body"), headers)
	require.NoError(t, err)
	defer func() {
		err := resp.Body.Close()
		require.NoError(t, err, "failed to close response body")
	}()

	// Then
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestOutgoingRequestPreservesTraceID(t *testing.T) {
	// Given
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "test-trace-id", r.Header.Get("X-Trace-ID"))
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/test", nil)
	ctx.Set("X-Trace-ID", "test-trace-id")
	client := NewClient(nil)

	// When
	resp, err := client.OutgoingRequest(ctx, http.MethodGet, server.URL, nil, nil)
	require.NoError(t, err)
	defer resp.Body.Close()

	// Then
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
