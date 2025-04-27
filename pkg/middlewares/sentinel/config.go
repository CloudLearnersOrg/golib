package sentinel

// Config defines the security settings used by the sentinel middleware.
type Config struct {
	// ExpectedHosts is a list of allowed hosts (including ports if needed).
	// If empty, host validation is disabled.
	// Example: []string{"api.example.com", "localhost:3000"}
	ExpectedHosts []string

	// XFrameOptions controls the X-Frame-Options header.
	// Common values: "DENY", "SAMEORIGIN"
	XFrameOptions string

	// ContentSecurityPolicy controls the Content-Security-Policy header.
	ContentSecurityPolicy string

	// EnableXSSProtection enables the X-XSS-Protection header.
	// When true, sets "1; mode=block"
	EnableXSSProtection bool

	// StrictTransportSecurityPolicy controls the Strict-Transport-Security header.
	StrictTransportSecurityPolicy string

	// ReferrerPolicy controls the Referrer-Policy header.
	ReferrerPolicy string

	// XContentTypeOptions controls the X-Content-Type-Options header.
	XContentTypeOptions string

	// PermissionsPolicy controls the Permissions-Policy header.
	PermissionsPolicy string
}

// DefaultConfig returns a configuration with secure default settings.
// This provides a baseline of security suitable for most applications,
// but should be customized for specific requirements.
func DefaultConfig() *Config {
	return &Config{
		ExpectedHosts:                 []string{},
		XFrameOptions:                 "DENY",
		ContentSecurityPolicy:         "default-src 'self'; connect-src *; font-src *; script-src-elem * 'unsafe-inline'; img-src * data:; style-src * 'unsafe-inline';",
		EnableXSSProtection:           true,
		StrictTransportSecurityPolicy: "max-age=31536000; includeSubDomains; preload",
		ReferrerPolicy:                "strict-origin",
		XContentTypeOptions:           "nosniff",
		PermissionsPolicy:             "geolocation=(),midi=(),sync-xhr=(),microphone=(),camera=(),magnetometer=(),gyroscope=(),fullscreen=(self),payment=()",
	}
}
