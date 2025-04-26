package session

import (
	"fmt"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
)

func NewMiddleware(cfg SessionConfig) (gin.HandlerFunc, error) {
	// Set password to empty string if not provided
	password := cfg.RedisPassword
	if password == "" {
		password = ""
	}

	// Initialize Redis store with conditional password
	redisStore, err := redis.NewStore(10, "tcp",
		fmt.Sprintf("%s:%d", cfg.RedisHost, cfg.RedisPort),
		password,
		cfg.SessionSecret)

	if err != nil {
		return nil, err
	}

	// Configure session cookie
	redisStore.Options(sessions.Options{
		Path:     cfg.CookiePath,
		Domain:   cfg.CookieDomain,
		MaxAge:   cfg.SessionMaxAge,
		Secure:   cfg.CookieSecure,
		HttpOnly: cfg.CookieHttpOnly,
		SameSite: cfg.SameSite,
	})

	return sessions.Sessions(SessionName, redisStore), nil
}
