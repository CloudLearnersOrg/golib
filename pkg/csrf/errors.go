package csrf

import "errors"

var (
	ErrInvalidFormat    = errors.New("invalid CSRF token format")
	ErrExpiredToken     = errors.New("CSRF token has expired")
	ErrInvalidSignature = errors.New("CSRF token signature is invalid")
	ErrInvalidTimestamp = errors.New("invalid timestamp in CSRF token")
)
