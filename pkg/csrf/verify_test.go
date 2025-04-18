package csrf

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"testing"
	"time"
)

func TestVerifyCSRFToken(t *testing.T) {
	secret := "test-secret"
	maxAge := time.Minute * 5

	t.Run("Given a valid CSRF token, When it is verified, Then it should succeed", func(t *testing.T) {
		// Given
		token := GenerateToken(secret)

		// When
		err := VerifyCSRFToken(token, secret, maxAge)

		// Then
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("Given an expired CSRF token, When it is verified, Then it should fail", func(t *testing.T) {
		// Given
		timestamp := time.Now().Add(-maxAge - time.Minute).Unix()
		message := fmt.Sprintf("%d", timestamp)
		h := hmac.New(sha256.New, []byte(secret))
		h.Write([]byte(message))
		signature := hex.EncodeToString(h.Sum(nil))
		token := fmt.Sprintf("%s:%s", message, signature)

		// When
		err := VerifyCSRFToken(token, secret, maxAge)

		// Then
		if err == nil || err.Error() != "CSRF token has expired" {
			t.Errorf("expected 'CSRF token has expired', got %v", err)
		}
	})

	t.Run("Given a token with an invalid signature, When it is verified, Then it should fail", func(t *testing.T) {
		// Given
		token := GenerateToken(secret)
		invalidToken := token[:len(token)-1] + "x" // Corrupt the token

		// When
		err := VerifyCSRFToken(invalidToken, secret, maxAge)

		// Then
		if err == nil || err.Error() != "CSRF token signature is invalid" {
			t.Errorf("expected 'CSRF token signature is invalid', got %v", err)
		}
	})

	t.Run("Given a token with an invalid format, When it is verified, Then it should fail", func(t *testing.T) {
		// Given
		token := "invalid-token-format"

		// When
		err := VerifyCSRFToken(token, secret, maxAge)

		// Then
		if err == nil || err.Error() != "invalid CSRF token format" {
			t.Errorf("expected 'invalid CSRF token format', got %v", err)
		}
	})

	t.Run("Given a token with an invalid timestamp, When it is verified, Then it should fail", func(t *testing.T) {
		// Given
		token := "invalid-timestamp:signature"

		// When
		err := VerifyCSRFToken(token, secret, maxAge)

		// Then
		if err == nil || err.Error() != "invalid timestamp in CSRF token: strconv.ParseInt: parsing \"invalid-timestamp\": invalid syntax" {
			t.Errorf("expected 'invalid timestamp in CSRF token', got %v", err)
		}
	})
}
