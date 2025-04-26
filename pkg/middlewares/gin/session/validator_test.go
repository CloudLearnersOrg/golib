package session

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestValidateSession(t *testing.T) {
	t.Run("Given a request with valid session, When ValidateSession middleware is applied, Then it should continue", func(t *testing.T) {
		// Given
		gin.SetMode(gin.TestMode)
		router := gin.New()

		// Configure session storage
		store := cookie.NewStore([]byte("test-secret"))
		router.Use(sessions.Sessions(SessionName, store))

		// Create a route that sets up a valid session before testing
		router.GET("/setup", func(c *gin.Context) {
			session := sessions.Default(c)
			session.Set(UserKey, "test-user-id")
			err := session.Save()
			assert.NoError(t, err)
			c.String(http.StatusOK, "Session created")
		})

		// Create a protected route with the ValidateSession middleware
		protected := router.Group("/")
		protected.Use(ValidateSession())
		protected.GET("/protected", func(c *gin.Context) {
			c.String(http.StatusOK, "Access granted")
		})

		// First, set up the session
		w1 := httptest.NewRecorder()
		req1 := httptest.NewRequest(http.MethodGet, "/setup", nil)
		router.ServeHTTP(w1, req1)
		assert.Equal(t, http.StatusOK, w1.Code)

		// Extract the session cookie
		cookies := w1.Result().Cookies()
		assert.NotEmpty(t, cookies, "No cookies were set")

		// Now make the request to protected route with the cookie
		req2 := httptest.NewRequest(http.MethodGet, "/protected", nil)
		for _, cookie := range cookies {
			req2.AddCookie(cookie)
		}
		w2 := httptest.NewRecorder()

		// When
		router.ServeHTTP(w2, req2)

		// Then
		assert.Equal(t, http.StatusOK, w2.Code)
	})

	t.Run("Given a request without valid session, When ValidateSession middleware is applied, Then it should return 401", func(t *testing.T) {
		// Given
		gin.SetMode(gin.TestMode)
		router := gin.New()

		// Configure session storage
		store := cookie.NewStore([]byte("test-secret"))
		router.Use(sessions.Sessions(SessionName, store))

		// Apply the middleware being tested
		router.Use(ValidateSession())

		// Add a handler that shouldn't be called
		router.GET("/protected", func(c *gin.Context) {
			c.String(http.StatusOK, "This shouldn't be called")
		})

		// Request without session cookie
		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		w := httptest.NewRecorder()

		// When
		router.ServeHTTP(w, req)

		// Then
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Given a request with valid session, When ValidateSession is applied, Then it should set UserKey in context", func(t *testing.T) {
		// Given
		gin.SetMode(gin.TestMode)
		router := gin.New()

		// Configure session storage
		store := cookie.NewStore([]byte("test-secret"))
		router.Use(sessions.Sessions(SessionName, store))

		// Create a route that sets up a valid session before testing
		router.GET("/setup", func(c *gin.Context) {
			session := sessions.Default(c)
			session.Set(UserKey, "test-user-id")
			err := session.Save()
			assert.NoError(t, err)
			c.String(http.StatusOK, "Session created")
		})

		// Create a protected route with the ValidateSession middleware
		// that checks for user ID in context
		protected := router.Group("/")
		protected.Use(ValidateSession())

		var userIDInContext string
		protected.GET("/protected", func(c *gin.Context) {
			userID, exists := c.Get(UserKey)
			if exists {
				userIDInContext = userID.(string)
			}
			c.String(http.StatusOK, "Access granted")
		})

		// First, set up the session
		w1 := httptest.NewRecorder()
		req1 := httptest.NewRequest(http.MethodGet, "/setup", nil)
		router.ServeHTTP(w1, req1)

		// Extract the session cookie
		cookies := w1.Result().Cookies()

		// Now make the request to protected route with the cookie
		req2 := httptest.NewRequest(http.MethodGet, "/protected", nil)
		for _, cookie := range cookies {
			req2.AddCookie(cookie)
		}
		w2 := httptest.NewRecorder()

		// When
		router.ServeHTTP(w2, req2)

		// Then
		assert.Equal(t, "test-user-id", userIDInContext)
	})
}
