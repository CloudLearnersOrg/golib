package session

import (
	"fmt"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
)

func NewMiddleware(cfg SessionConfig) (gin.HandlerFunc, error) {
	redisStore, err := redis.NewStore(10, "tcp",
		fmt.Sprintf("%s:%d", cfg.RedisHost, cfg.RedisPort),
		cfg.RedisPassword,
		cfg.SessionSecret)
	if err != nil {
		return nil, err
	}

	// Configure session cookie
	redisStore.Options(sessions.Options{
		Path:     "/",
		Domain:   cfg.CookieDomain,
		MaxAge:   cfg.SessionMaxAge,
		Secure:   cfg.CookieSecure,
		HttpOnly: true,
		SameSite: cfg.SameSite,
	})

	return sessions.Sessions(SessionName, redisStore), nil
}
