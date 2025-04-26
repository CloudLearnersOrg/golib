package session

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMiddleware(t *testing.T) {
	t.Run("Given valid config, When NewMiddleware is called, Then it should return a middleware function", func(t *testing.T) {
		// Skip this test if we can't connect to Redis
		t.Skip("Skipping Redis-dependent test - mock this in real tests")

		// Given
		config := SessionConfig{
			RedisHost:     "localhost",
			RedisPort:     6379,
			RedisPassword: "",
			SessionSecret: "test-secret",
			CookieDomain:  "example.com",
			CookieSecure:  true,
			SessionMaxAge: 3600,
			SameSite:      http.SameSiteLaxMode,
		}

		// When
		middleware, err := NewMiddleware(config)

		// Then
		assert.NoError(t, err)
		assert.NotNil(t, middleware)
	})

	t.Run("Given invalid Redis config, When NewMiddleware is called, Then it should return an error", func(t *testing.T) {
		// Given
		config := SessionConfig{
			RedisHost:     "nonexistent-host",
			RedisPort:     12345,
			RedisPassword: "wrong-password",
			SessionSecret: "test-secret",
		}

		// When
		middleware, err := NewMiddleware(config)

		// Then
		assert.Error(t, err)
		assert.Nil(t, middleware)
	})
}
