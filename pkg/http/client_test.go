package httpclient

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
		w.Write([]byte("test response"))
	}))
	defer server.Close()

	w := httptest.NewRecorder()
	ginCtx, _ := gin.CreateTestContext(w)
	ginCtx.Request = httptest.NewRequest(http.MethodGet, "/test", nil)
	client := NewClient(nil)

	// When
	resp, err := client.OutgoingRequest(ginCtx, http.MethodGet, server.URL, nil)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestOutgoingPostRequest(t *testing.T) {
	// Given
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		body, _ := io.ReadAll(r.Body)
		assert.Equal(t, "test body", string(body))
		w.Write([]byte("test response"))
	}))
	defer server.Close()

	w := httptest.NewRecorder()
	ginCtx, _ := gin.CreateTestContext(w)
	ginCtx.Request = httptest.NewRequest(http.MethodPost, "/test", nil)
	client := NewClient(nil)

	// When
	resp, err := client.OutgoingRequest(ginCtx, http.MethodPost, server.URL, strings.NewReader("test body"))

	// Then
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestMiddlewareWithoutTraceID(t *testing.T) {
	// Given
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(Middleware())
	r.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)

	// When
	r.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEmpty(t, w.Header().Get(TraceIDKey))
}

func TestMiddlewareWithExistingTraceID(t *testing.T) {
	// Given
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(Middleware())
	r.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	existingTraceID := uuid.New().String()
	req.Header.Set(TraceIDKey, existingTraceID)

	// When
	r.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, existingTraceID, w.Header().Get(TraceIDKey))
}

func TestTraceIDPropagationThroughChain(t *testing.T) {
	// Given
	downstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceID := r.Header.Get(TraceIDKey)
		w.Header().Set(TraceIDKey, traceID)
		w.WriteHeader(http.StatusOK)
	}))
	defer downstream.Close()

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(Middleware())
	client := NewClient(nil)

	r.GET("/test", func(c *gin.Context) {
		resp, err := client.OutgoingRequest(c, http.MethodGet, downstream.URL, nil)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()
		c.Status(resp.StatusCode)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	initialTraceID := uuid.New().String()
	req.Header.Set(TraceIDKey, initialTraceID)

	// When
	r.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, initialTraceID, w.Header().Get(TraceIDKey))
}
