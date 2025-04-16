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
	server := httptest.NewServer(http.HandlerFunc(func(body http.ResponseWriter, req *http.Request) {
		assert.Equal(t, http.MethodGet, req.Method)
		body.Write([]byte("test response"))
	}))
	defer server.Close()

	body := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(body)
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
	server := httptest.NewServer(http.HandlerFunc(func(body http.ResponseWriter, req *http.Request) {
		assert.Equal(t, http.MethodPost, req.Method)
		resp, err := io.ReadAll(req.Body)
		if err != nil {
			body.WriteHeader(http.StatusInternalServerError)
			return
		}
		assert.Equal(t, "test body", string(resp))
		body.Write([]byte("test response"))
	}))
	defer server.Close()

	body := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(body)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/test", nil)
	client := NewClient(nil)

	// When
	resp, err := client.OutgoingRequest(ctx, http.MethodPost, server.URL, strings.NewReader("test body"))

	// Then
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestMiddlewareWithoutTraceID(t *testing.T) {
	// Given
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(Middleware())
	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	resp := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)

	// When
	router.ServeHTTP(resp, req)

	// Then
	assert.Equal(t, http.StatusOK, resp.Code)
	assert.NotEmpty(t, resp.Header().Get("X-Trace-ID"))
}

func TestMiddlewareWithExistingTraceID(t *testing.T) {
	// Given
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(Middleware())
	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	resp := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	existingTraceID := uuid.New().String()
	req.Header.Set("X-Trace-ID", existingTraceID)

	// When
	router.ServeHTTP(resp, req)

	// Then
	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, existingTraceID, resp.Header().Get("X-Trace-ID"))
}

func TestTraceIDPropagationThroughChain(t *testing.T) {
	// Given
	downstream := httptest.NewServer(http.HandlerFunc(func(body http.ResponseWriter, req *http.Request) {
		traceID := req.Header.Get("X-Trace-ID")
		body.Header().Set("X-Trace-ID", traceID)
		body.WriteHeader(http.StatusOK)
	}))
	defer downstream.Close()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(Middleware())
	client := NewClient(nil)

	router.GET("/test", func(c *gin.Context) {
		resp, err := client.OutgoingRequest(c, http.MethodGet, downstream.URL, nil)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()
		c.Status(resp.StatusCode)
	})

	resp := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	initialTraceID := uuid.New().String()
	req.Header.Set("X-Trace-ID", initialTraceID)

	// When
	router.ServeHTTP(resp, req)

	// Then
	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, initialTraceID, resp.Header().Get("X-Trace-ID"))
}
