package csrf

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

func GenerateToken(secret string) string {
	timestamp := time.Now().Unix()
	message := fmt.Sprintf("%d", timestamp)

	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(message))

	signature := hex.EncodeToString(h.Sum(nil))
	return fmt.Sprintf("%s:%s", message, signature)
}

func GenerateSecret() (string, error) {
	secret := make([]byte, 32)
	_, err := rand.Read(secret)
	if err != nil {
		return "", fmt.Errorf("failed to generate CSRF token secret: %w", err)
	}

	return hex.EncodeToString(secret), nil
}
