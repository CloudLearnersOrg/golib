package cors

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter(config Config) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(New(config))

	r.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "GET OK")
	})

	r.POST("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "POST OK")
	})

	r.OPTIONS("/test", func(c *gin.Context) {
		// Should never reach here for OPTIONS due to middleware handling
		c.String(http.StatusOK, "OPTIONS OK")
	})

	return r
}

func TestWildcardOrigin(t *testing.T) {
	// Given
	config := DefaultConfig()
	config.AllowOrigins = []string{"*"}
	r := setupRouter(config)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Origin", "http://example.com")
	w := httptest.NewRecorder()

	// When
	r.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
}

func TestSpecificOriginAllowed(t *testing.T) {
	// Given
	config := DefaultConfig()
	config.AllowOrigins = []string{"http://allowed-origin.com", "http://another-allowed.com"}
	r := setupRouter(config)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Origin", "http://allowed-origin.com")
	w := httptest.NewRecorder()

	// When
	r.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "http://allowed-origin.com", w.Header().Get("Access-Control-Allow-Origin"))
}

func TestSpecificOriginNotAllowed(t *testing.T) {
	// Given
	config := DefaultConfig()
	config.AllowOrigins = []string{"http://allowed-origin.com"}
	r := setupRouter(config)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Origin", "http://not-allowed-origin.com")
	w := httptest.NewRecorder()

	// When
	r.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "", w.Header().Get("Access-Control-Allow-Origin"))
}

func TestAllowMethodsHeader(t *testing.T) {
	// Given
	config := DefaultConfig()
	config.AllowMethods = []string{"GET", "POST", "PUT"}
	r := setupRouter(config)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	// When
	r.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "GET, POST, PUT", w.Header().Get("Access-Control-Allow-Methods"))
}

func TestAllowHeadersHeader(t *testing.T) {
	// Given
	config := DefaultConfig()
	config.AllowHeaders = []string{"Content-Type", "Authorization", "X-CSRF-Token"}
	r := setupRouter(config)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	// When
	r.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "Content-Type, Authorization, X-CSRF-Token", w.Header().Get("Access-Control-Allow-Headers"))
}

func TestExposeHeadersHeader(t *testing.T) {
	// Given
	config := DefaultConfig()
	config.ExposeHeaders = []string{"Content-Length", "X-Custom-Header"}
	r := setupRouter(config)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	// When
	r.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "Content-Length, X-Custom-Header", w.Header().Get("Access-Control-Expose-Headers"))
}

func TestAllowCredentialsHeader(t *testing.T) {
	// Given
	config := DefaultConfig()
	config.AllowCredentials = true
	r := setupRouter(config)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	// When
	r.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "true", w.Header().Get("Access-Control-Allow-Credentials"))
}

func TestPreflightRequests(t *testing.T) {
	// Given
	config := DefaultConfig()
	config.AllowOrigins = []string{"http://example.com"}
	config.AllowMethods = []string{"GET", "POST"}
	config.AllowHeaders = []string{"Content-Type"}
	r := setupRouter(config)

	req := httptest.NewRequest(http.MethodOptions, "/test", nil)
	req.Header.Set("Origin", "http://example.com")
	w := httptest.NewRecorder()

	// When
	r.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Equal(t, "http://example.com", w.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "GET, POST", w.Header().Get("Access-Control-Allow-Methods"))
	assert.Equal(t, "Content-Type", w.Header().Get("Access-Control-Allow-Headers"))
}

func TestEmptyOriginConfiguration(t *testing.T) {
	// Given
	config := Config{
		// AllowOrigins intentionally empty
		AllowMethods: []string{"GET"},
	}
	r := setupRouter(config)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Origin", "http://example.com")
	w := httptest.NewRecorder()

	// When
	r.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "", w.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "GET", w.Header().Get("Access-Control-Allow-Methods"))
}

func TestDefaultConfigAndMiddleware(t *testing.T) {
	// Given
	defaultConfig := DefaultConfig()

	// When + Then
	assert.Equal(t, []string{"*"}, defaultConfig.AllowOrigins)
	assert.Equal(t, []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}, defaultConfig.AllowMethods)
	assert.Equal(t, []string{"Origin", "Content-Type", "Accept", "Authorization"}, defaultConfig.AllowHeaders)
	assert.Equal(t, []string{"Content-Length", "Content-Type"}, defaultConfig.ExposeHeaders)
	assert.False(t, defaultConfig.AllowCredentials)

	// Test that Middleware() returns a handler using the default config
	handler := Middleware()
	assert.NotNil(t, handler)
}

func TestCompleteConfig(t *testing.T) {
	// Given
	config := Config{
		AllowOrigins:     []string{"http://example.com", "https://api.example.com"},
		AllowMethods:     []string{"GET", "POST", "PUT"},
		AllowHeaders:     []string{"Content-Type", "Authorization", "X-CSRF-Token"},
		ExposeHeaders:    []string{"Content-Length", "X-Request-ID"},
		AllowCredentials: true,
	}
	r := setupRouter(config)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Origin", "http://example.com")
	w := httptest.NewRecorder()

	// When
	r.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "http://example.com", w.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "GET, POST, PUT", w.Header().Get("Access-Control-Allow-Methods"))
	assert.Equal(t, "Content-Type, Authorization, X-CSRF-Token", w.Header().Get("Access-Control-Allow-Headers"))
	assert.Equal(t, "Content-Length, X-Request-ID", w.Header().Get("Access-Control-Expose-Headers"))
	assert.Equal(t, "true", w.Header().Get("Access-Control-Allow-Credentials"))
}
