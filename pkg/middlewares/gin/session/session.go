package session

import (
	"errors"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// RotateSessionID generates a new session ID and updates the session store.
func RotateSessionID(c *gin.Context) error {
	session := sessions.Default(c)
	userID := session.Get(UserKey)
	if userID == nil {
		return errors.New("no active session to rotate")
	}

	// Clear the current session
	session.Clear()

	// Create a new session
	newSession := sessions.Default(c)
	newSession.Set(UserKey, userID)
	return newSession.Save()
}

// CreateSession creates a new session for the user and stores the user ID in the session.
func CreateSession(c *gin.Context, userID string) error {
	session := sessions.Default(c)
	session.Set(UserKey, userID)
	return session.Save()
}

// DestroySession clears the session for the user, effectively logging them out.
func DestroySession(c *gin.Context) error {
	session := sessions.Default(c)
	session.Clear()
	session.Options(sessions.Options{MaxAge: -1})
	return session.Save()
}

// RefreshSession updates the session expiry time
func RefreshSession(c *gin.Context) error {
	session := sessions.Default(c)
	userID := session.Get(UserKey)
	if userID == nil {
		return errors.New("no active session to refresh")
	}
	return session.Save()
}

// SetSessionData stores additional data in the session
func SetSessionData(c *gin.Context, key string, value interface{}) error {
	session := sessions.Default(c)
	session.Set(key, value)
	return session.Save()
}

// GetSessionData retrieves additional data from the session
func GetSessionData(c *gin.Context, key string) (interface{}, bool) {
	session := sessions.Default(c)
	value := session.Get(key)
	return value, value != nil
}
