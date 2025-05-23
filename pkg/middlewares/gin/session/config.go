package session

import (
	"net/http"
)

const (
	UserKey     = "user_id"
	SessionName = "auth_session"
)

// SessionConfig holds configuration for session middleware
type SessionConfig struct {
	RedisHost               string
	RedisPort               int
	RedisPassword           string
	RedisConnectionPoolSize int
	SessionSecret           string
	CookiePath              string
	CookieDomain            string
	CookieSecure            bool
	CookieHttpOnly          bool
	SessionMaxAge           int
	SameSite                http.SameSite
	UserKeySecret           string
}
