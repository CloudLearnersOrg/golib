// Package csrf provides Cross-Site Request Forgery (CSRF) token generation and verification.
// It implements secure token creation using HMAC-SHA256 and includes timestamp-based
// expiration for enhanced security.
//
// Features:
//   - Secure token generation using HMAC-SHA256
//   - Random secret generation for token signing
//   - Token verification with expiration checking
//   - Built-in error types for specific failure cases
//
// Basic Usage:
//
//	// Generate a secret key (store this securely)
//	secret, err := csrf.GenerateSecret()
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Generate a token
//	token := csrf.GenerateToken(secret)
//
//	// Later, verify the token
//	err = csrf.VerifyCSRFToken(token, secret, 1*time.Hour)
//	if err != nil {
//	    switch err {
//	    case csrf.ErrExpiredToken:
//	        // Handle expired token
//	    case csrf.ErrInvalidSignature:
//	        // Handle invalid signature
//	    case csrf.ErrInvalidFormat:
//	        // Handle invalid format
//	    }
//	}
//
// Token Format:
// The CSRF token consists of two parts separated by a colon:
//   - Unix timestamp (seconds since epoch)
//   - HMAC-SHA256 signature of the timestamp
//
// Example token: "1650644857:a1b2c3d4e5f6..."
//
// Security Features:
//   - Timing attack resistant comparison
//   - Cryptographically secure random secret generation
//   - Timestamp-based expiration
//   - HMAC-SHA256 for token signing
//
// Error Handling:
// The package provides specific error types for different failure cases:
//   - ErrInvalidFormat: Token format is incorrect
//   - ErrExpiredToken: Token has exceeded its maximum age
//   - ErrInvalidSignature: Token signature verification failed
//   - ErrInvalidTimestamp: Timestamp parsing failed
//
// Best Practices:
//  1. Generate and store the secret securely
//  2. Use HTTPS to transmit tokens
//  3. Set appropriate token expiration times
//  4. Validate tokens on all state-changing operations
//  5. Include tokens in custom headers or hidden form fields
package csrf
