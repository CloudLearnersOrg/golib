package session

import (
	"fmt"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
)

func NewMiddleware(cfg SessionConfig) (gin.HandlerFunc, error) {
	// Initialize Redis store with conditional password
	redisStore, err := redis.NewStore(
		cfg.RedisConnectionPoolSize, // pool size
		"tcp",
		fmt.Sprintf("%s:%d", cfg.RedisHost, cfg.RedisPort),
		"",                        // username (empty for most Redis versions)
		cfg.RedisPassword,         // password as string
		[]byte(cfg.SessionSecret), // key pairs as []byte
	)
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
