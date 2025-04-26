package csrf

import (
	"github.com/CloudLearnersOrg/golib/pkg/csrf"
	"github.com/CloudLearnersOrg/golib/pkg/middlewares/gin/session"
	"github.com/gin-gonic/gin"
)

// Initialize creates a CSRF secret and stores it in the session
func Initialize(c *gin.Context) (string, error) {
	// Check if we already have a CSRF secret in the session
	secretValue, exists := session.GetSessionData(c, TokenKey)

	if exists {
		return secretValue.(string), nil
	}

	// Generate a new secret
	secret, err := csrf.GenerateSecret()
	if err != nil {
		return "", err
	}

	// Store it in the session
	if err := session.SetSessionData(c, TokenKey, secret); err != nil {
		return "", err
	}

	return secret, nil
}

// GetToken generates a new CSRF token for the current session
func GetToken(c *gin.Context) (string, error) {
	secret, err := Initialize(c)
	if err != nil {
		return "", err
	}

	return csrf.GenerateToken(secret), nil
}
