package csrf

import (
	"encoding/hex"
	"strings"
	"testing"
)

func TestGenerateToken(t *testing.T) {
	secret := "test-secret"

	t.Run("Given a valid secret, When GenerateToken is called, Then it should return a valid token", func(t *testing.T) {
		// Given
		// Secret is already defined

		// When
		token := GenerateToken(secret)

		// Then
		parts := strings.Split(token, ":")
		if len(parts) != 2 {
			t.Errorf("expected token to have 2 parts separated by ':', got %v", parts)
		}

		timestamp := parts[0]
		signature := parts[1]

		if timestamp == "" {
			t.Errorf("expected timestamp to be non-empty, got empty string")
		}

		if signature == "" {
			t.Errorf("expected signature to be non-empty, got empty string")
		}
	})
}

func TestGenerateSecret(t *testing.T) {
	t.Run("When GenerateSecret is called, Then it should return a valid secret", func(t *testing.T) {
		// When
		secret, err := GenerateSecret()

		// Then
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		decoded, err := hex.DecodeString(secret)
		if err != nil {
			t.Errorf("expected secret to be a valid hex string, got error: %v", err)
		}

		if len(decoded) != 32 {
			t.Errorf("expected secret to be 32 bytes, got %d bytes", len(decoded))
		}
	})
}
