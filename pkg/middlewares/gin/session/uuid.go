package session

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetCurrentUserID(c *gin.Context) (uuid.UUID, error) {
	userIDKey, exists := c.Get(UserKey)
	if !exists {
		return uuid.Nil, errors.New("user ID not found in context")
	}

	userIDStr, ok := userIDKey.(string)
	if !ok {
		return uuid.Nil, errors.New("user ID is not a string")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return uuid.Nil, errors.New("invalid user ID format")
	}

	return userID, nil
}
