package csrf

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// setupSessionContext creates a testing context with a cookie session
func setupSessionContext() (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest("GET", "/", nil)
	c.Request = req

	// Set up the session
	store := cookie.NewStore([]byte("test-secret"))
	sessions.Sessions("csrf-test", store)(c)

	return c, w
}

// TestInitializeWithExistingSecret tests Initialize when a secret already exists
func TestInitializeWithExistingSecret(t *testing.T) {
	// Given
	c, _ := setupSessionContext()
	session := sessions.Default(c) // Initialize the session
	expectedSecret := "existing-secret"
	session.Set(TokenKey, expectedSecret)
	session.Save()

	// When
	secret, err := Initialize(c)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, expectedSecret, secret)
}

// Since we can't mock the GenerateSecret function directly, we'll test the initialize
// function's behavior with a real session
func TestInitializeWithoutSecret(t *testing.T) {
	// Given
	c, _ := setupSessionContext()
	session := sessions.Default(c)

	// Make sure no secret exists
	session.Delete(TokenKey)
	session.Save()

	// When
	secret, err := Initialize(c)

	// Then
	assert.NoError(t, err)
	assert.NotEmpty(t, secret, "Generated secret should not be empty")

	// Verify it was stored in session
	storedSecret := sessions.Default(c).Get(TokenKey)
	assert.Equal(t, secret, storedSecret)
}

// Test that we can get the token successfully
func TestGetTokenSuccess(t *testing.T) {
	// Given
	c, _ := setupSessionContext()
	session := sessions.Default(c)
	session.Set(TokenKey, "test-secret")
	session.Save()

	// When
	token, err := GetToken(c)

	// Then
	assert.NoError(t, err)
	assert.NotEmpty(t, token, "Generated token should not be empty")
}

// Test that GetToken properly initializes when no secret exists
func TestGetTokenInitializeSecret(t *testing.T) {
	// Given
	c, _ := setupSessionContext()
	session := sessions.Default(c)

	// Make sure no secret exists
	session.Delete(TokenKey)
	session.Save()

	// When
	token, err := GetToken(c)

	// Then
	assert.NoError(t, err)
	assert.NotEmpty(t, token, "Generated token should not be empty")

	// Verify secret was stored in session
	storedSecret := sessions.Default(c).Get(TokenKey)
	assert.NotNil(t, storedSecret, "Secret should be stored in session")
}

// Test that GetToken returns proper values
func TestGetTokenReturnValues(t *testing.T) {
	// Given
	c, _ := setupSessionContext()
	session := sessions.Default(c)
	session.Set(TokenKey, "fixed-test-secret")
	session.Save()

	// When
	token1, err1 := GetToken(c)
	token2, err2 := GetToken(c)

	// Then
	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.NotEmpty(t, token1)
	assert.NotEmpty(t, token2)

	// The tokens should be deterministic for the same secret
	assert.Equal(t, token1, token2, "Tokens should be consistent for the same secret")
}

// Test token format
func TestGetTokenFormat(t *testing.T) {
	// Given
	c, _ := setupSessionContext()
	session := sessions.Default(c)
	session.Set(TokenKey, "test-secret")
	session.Save()

	// When
	token, err := GetToken(c)

	// Then
	assert.NoError(t, err)
	assert.Contains(t, token, ":", "Token should contain a colon separator")

	// Basic format check - tokens from csrf package are timestamp:signature
	parts := len(token)
	assert.True(t, parts > 10, "Token should have sufficient length")
}

func TestMultipleTokensWithSameSecret(t *testing.T) {
	// Given
	c, _ := setupSessionContext()
	session := sessions.Default(c)
	session.Set(TokenKey, "multi-test-secret")
	session.Save()

	// When
	token1, _ := GetToken(c)
	token2, _ := GetToken(c)
	token3, _ := GetToken(c)

	// Then - tokens should be deterministic for the same secret
	assert.Equal(t, token1, token2)
	assert.Equal(t, token2, token3)
}

func TestGetTokenAfterInitialize(t *testing.T) {
	// Given
	c, _ := setupSessionContext()

	// When - first initialize
	secret, err1 := Initialize(c)
	assert.NoError(t, err1)

	// Then get token
	token, err2 := GetToken(c)

	// Then
	assert.NoError(t, err2)
	assert.NotEmpty(t, token)

	// The session should still have the same secret
	storedSecret := sessions.Default(c).Get(TokenKey)
	assert.Equal(t, secret, storedSecret)
}
