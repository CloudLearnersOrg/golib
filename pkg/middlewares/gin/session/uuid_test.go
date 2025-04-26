package session

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetCurrentUserID(t *testing.T) {
	t.Run("Given a context with valid UUID, When GetCurrentUserID is called, Then it should return the UUID", func(t *testing.T) {
		// Given
		c, _ := gin.CreateTestContext(nil)
		validUUID := uuid.New()
		c.Set(UserKey, validUUID.String())

		// When
		userID, err := GetCurrentUserID(c)

		// Then
		assert.NoError(t, err)
		assert.Equal(t, validUUID, userID)
	})

	t.Run("Given a context without user ID, When GetCurrentUserID is called, Then it should return error", func(t *testing.T) {
		// Given
		c, _ := gin.CreateTestContext(nil)

		// When
		userID, err := GetCurrentUserID(c)

		// Then
		assert.Error(t, err)
		assert.Equal(t, "user ID not found in context", err.Error())
		assert.Equal(t, uuid.Nil, userID)
	})

	t.Run("Given a context with non-string user ID, When GetCurrentUserID is called, Then it should return error", func(t *testing.T) {
		// Given
		c, _ := gin.CreateTestContext(nil)
		c.Set(UserKey, 12345) // Not a string

		// When
		userID, err := GetCurrentUserID(c)

		// Then
		assert.Error(t, err)
		assert.Equal(t, "user ID is not a string", err.Error())
		assert.Equal(t, uuid.Nil, userID)
	})

	t.Run("Given a context with invalid UUID string, When GetCurrentUserID is called, Then it should return error", func(t *testing.T) {
		// Given
		c, _ := gin.CreateTestContext(nil)
		c.Set(UserKey, "not-a-valid-uuid")

		// When
		userID, err := GetCurrentUserID(c)

		// Then
		assert.Error(t, err)
		assert.Equal(t, "invalid user ID format", err.Error())
		assert.Equal(t, uuid.Nil, userID)
	})
}
