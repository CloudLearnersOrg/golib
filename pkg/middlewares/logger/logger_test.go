package logger

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

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
