package csrf

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// VerifyCSRFToken verifies the provided CSRF token using the secret.
// It checks the token's signature and ensures it is not expired.
func VerifyCSRFToken(token, secret string, maxAge time.Duration) error {
	parts := strings.Split(token, ":")
	if len(parts) != 2 {
		return errors.New("invalid CSRF token format")
	}

	timestampStr, signature := parts[0], parts[1]

	// Parse the timestamp
	timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid timestamp in CSRF token: %w", err)
	}

	// Check token expiration
	if time.Since(time.Unix(timestamp, 0)) > maxAge {
		return errors.New("CSRF token has expired")
	}

	// Recreate the HMAC signature
	message := fmt.Sprintf("%d", timestamp)
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(message))
	expectedSignature := hex.EncodeToString(h.Sum(nil))

	// Compare the provided signature with the expected signature
	if !hmac.Equal([]byte(signature), []byte(expectedSignature)) {
		return errors.New("CSRF token signature is invalid")
	}

	return nil
}
