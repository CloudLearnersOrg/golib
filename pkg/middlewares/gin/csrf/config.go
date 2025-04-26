package csrf

const (
	// TokenHeader is the header name for CSRF token
	TokenHeader = "X-CSRF-Token"
	// TokenFormField is the form field name for CSRF token
	TokenFormField = "csrf_token"
	// TokenKey is the session key for storing CSRF token
	TokenKey = "csrf_secret"
)
