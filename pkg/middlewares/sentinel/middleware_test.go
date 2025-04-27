package sentinel

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter(config *Config) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(New(config))

	r.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	return r
}

func TestHostValidation(t *testing.T) {
	t.Run("ValidHost", func(t *testing.T) {
		// Given
		config := DefaultConfig()
		config.ExpectedHosts = []string{"example.com", "localhost"}
		r := setupRouter(config)

		// When
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Host = "example.com"
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		// Then
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "OK", w.Body.String())
	})

	t.Run("InvalidHost", func(t *testing.T) {
		// Given
		config := DefaultConfig()
		config.ExpectedHosts = []string{"example.com", "localhost"}
		r := setupRouter(config)

		// When
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Host = "evil.com"
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		// Then
		assert.Equal(t, http.StatusForbidden, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid Host Header")
	})

	t.Run("NoHostValidation", func(t *testing.T) {
		// Given
		config := DefaultConfig()
		// Empty ExpectedHosts array - no host validation
		r := setupRouter(config)

		// When
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Host = "any-host.com"
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		// Then
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "OK", w.Body.String())
	})
}

func TestSecurityHeaders(t *testing.T) {
	// Given - a request with default config
	config := DefaultConfig()
	r := setupRouter(config)

	// When
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Then - all security headers should be set correctly
	assert.Equal(t, http.StatusOK, w.Code)

	headers := w.Header()
	assert.Equal(t, "DENY", headers.Get("X-Frame-Options"))
	assert.Equal(t, "1; mode=block", headers.Get("X-XSS-Protection"))
	assert.Equal(t, "max-age=31536000; includeSubDomains; preload", headers.Get("Strict-Transport-Security"))
	assert.Equal(t, "strict-origin", headers.Get("Referrer-Policy"))
	assert.Equal(t, "nosniff", headers.Get("X-Content-Type-Options"))
	assert.NotEmpty(t, headers.Get("Permissions-Policy"))
	assert.NotEmpty(t, headers.Get("Content-Security-Policy"))
}

func TestCustomHeaders(t *testing.T) {
	// Given - custom config
	config := DefaultConfig()
	config.XFrameOptions = "SAMEORIGIN"
	config.ContentSecurityPolicy = "default-src 'self'"
	config.StrictTransportSecurityPolicy = "max-age=86400"
	r := setupRouter(config)

	// When
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Then - custom headers should be used
	headers := w.Header()
	assert.Equal(t, "SAMEORIGIN", headers.Get("X-Frame-Options"))
	assert.Equal(t, "default-src 'self'", headers.Get("Content-Security-Policy"))
	assert.Equal(t, "max-age=86400", headers.Get("Strict-Transport-Security"))
}

func TestMiddlewareFunction(t *testing.T) {
	// Given
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(Middleware()) // Use the shorthand function

	r.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	// When
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Then - should use default config
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "DENY", w.Header().Get("X-Frame-Options"))
}
