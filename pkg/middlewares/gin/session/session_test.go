package session

import (
	"encoding/gob"
	"net/http/httptest"
	"testing"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func init() {
	// Register types for gob encoding (needed for cookie sessions)
	gob.Register(map[string]string{})
	gob.Register(map[string]interface{}{})
}

// setupTestContext creates a gin context with cookie-based session for testing
func setupTestContext() (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest("GET", "/", nil)
	c.Request = req

	// Use cookie store for testing instead of Redis
	store := cookie.NewStore([]byte("test-session-secret"))
	sessions.Sessions(SessionName, store)(c)

	return c, w
}

// setupTestContextWithUser creates a context with an authenticated user session
func setupTestContextWithUser(userID string) (*gin.Context, *httptest.ResponseRecorder) {
	c, w := setupTestContext()
	session := sessions.Default(c)
	session.Set(UserKey, userID)
	session.Save()
	return c, w
}

func TestCreateSession(t *testing.T) {
	t.Run("Given a context and user ID, When CreateSession is called, Then it should store the user ID in session", func(t *testing.T) {
		// Given
		c, _ := setupTestContext()
		userID := uuid.New().String()

		// When
		err := CreateSession(c, userID)

		// Then
		assert.NoError(t, err)

		// Verify session has user ID
		session := sessions.Default(c)
		storedID := session.Get(UserKey)
		assert.Equal(t, userID, storedID)
	})
}

func TestDestroySession(t *testing.T) {
	t.Run("Given a context with active session, When DestroySession is called, Then it should clear the session", func(t *testing.T) {
		// Given
		c, _ := setupTestContextWithUser("test-user-id")

		// When
		err := DestroySession(c)

		// Then
		assert.NoError(t, err)

		// Verify session is cleared
		session := sessions.Default(c)
		userID := session.Get(UserKey)
		assert.Nil(t, userID)
	})
}

func TestRotateSessionID(t *testing.T) {
	t.Run("Given a context with active session, When RotateSessionID is called, Then it should create a new session with same user ID", func(t *testing.T) {
		// Given
		c, _ := setupTestContextWithUser("test-user-id")

		// When
		err := RotateSessionID(c)

		// Then
		assert.NoError(t, err)

		// Verify user ID is preserved
		session := sessions.Default(c)
		userID := session.Get(UserKey)
		assert.Equal(t, "test-user-id", userID)
	})

	t.Run("Given a context without active session, When RotateSessionID is called, Then it should return error", func(t *testing.T) {
		// Given
		c, _ := setupTestContext()

		// When
		err := RotateSessionID(c)

		// Then
		assert.Error(t, err)
		assert.Equal(t, "no active session to rotate", err.Error())
	})
}

func TestRefreshSession(t *testing.T) {
	t.Run("Given a context with active session, When RefreshSession is called, Then it should update the session", func(t *testing.T) {
		// Given
		c, _ := setupTestContextWithUser("test-user-id")

		// When
		err := RefreshSession(c)

		// Then
		assert.NoError(t, err)
	})

	t.Run("Given a context without active session, When RefreshSession is called, Then it should return error", func(t *testing.T) {
		// Given
		c, _ := setupTestContext()

		// When
		err := RefreshSession(c)

		// Then
		assert.Error(t, err)
		assert.Equal(t, "no active session to refresh", err.Error())
	})
}

func TestSetSessionData(t *testing.T) {
	t.Run("Given a context, When SetSessionData is called, Then it should store the data in session", func(t *testing.T) {
		// Given
		c, _ := setupTestContext()

		// Use a simple string instead of a map for easier serialization
		testData := "test-value"

		// When
		err := SetSessionData(c, "preferences", testData)

		// Then
		assert.NoError(t, err)

		// Verify data is stored
		session := sessions.Default(c)
		storedData := session.Get("preferences")
		assert.Equal(t, testData, storedData)
	})
}

func TestGetSessionData(t *testing.T) {
	t.Run("Given a context with session data, When GetSessionData is called, Then it should return the data", func(t *testing.T) {
		// Given
		c, _ := setupTestContext()
		testData := map[string]string{"theme": "dark", "language": "en"}
		session := sessions.Default(c)
		session.Set("preferences", testData)
		session.Save()

		// When
		data, exists := GetSessionData(c, "preferences")

		// Then
		assert.True(t, exists)
		assert.Equal(t, testData, data)
	})

	t.Run("Given a context without session data, When GetSessionData is called, Then it should return false", func(t *testing.T) {
		// Given
		c, _ := setupTestContext()

		// When
		_, exists := GetSessionData(c, "nonexistent")

		// Then
		assert.False(t, exists)
	})
}
