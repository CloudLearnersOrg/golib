/*
Package sentinel provides security middleware for Go web applications built with the Gin framework.

The sentinel middleware offers two main security features:
 1. SSRF protection by validating Host headers against a list of expected hosts
 2. Security header management for modern web application protection

Basic Usage:

To use the middleware with default security settings:

	router := gin.Default()
	router.Use(sentinel.Middleware())

This will apply all security headers with sensible defaults, but without host validation.

Host Validation:

To enable Host header validation for SSRF protection:

	config := sentinel.DefaultConfig()
	config.ExpectedHosts = []string{"api.example.com", "localhost:3000"}
	router.Use(sentinel.New(config))

Loading from Environment Variables:

You can also configure the middleware from environment variables:

	hostsString := os.Getenv("EXPECTED_HOSTS")
	config := sentinel.DefaultConfig()
	config.ExpectedHosts = sentinel.ParseStringToHosts(hostsString)
	router.Use(sentinel.New(config))

Configuration Options:

The middleware is highly configurable through the Config struct:

	config := sentinel.DefaultConfig()
	config.XFrameOptions = "SAMEORIGIN"
	config.ReferrerPolicy = "strict-origin"
	config.PermissionsPolicy = "geolocation=(), camera=(self)"
	router.Use(sentinel.New(config))

Security Headers:

The middleware supports setting the following security headers:
  - X-Frame-Options
  - X-XSS-Protection
  - Strict-Transport-Security
  - Referrer-Policy
  - X-Content-Type-Options
  - Permissions-Policy
  - Content-Security-Policy

Best Practices:

1. In production, always define specific ExpectedHosts rather than allowing any host.
2. For Internet-facing applications, consider enabling HSTS with includeSubDomains and preload.
3. Review Content Security Policy regularly to ensure it meets your application's needs.
4. Use ParseStringToHosts helper to load configuration from environment variables.
*/
package sentinel
