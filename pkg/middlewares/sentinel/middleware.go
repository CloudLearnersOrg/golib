package sentinel

import (
	"strings"

	ginhttp "github.com/CloudLearnersOrg/golib/pkg/ginhttp/gin/statuses"
	"github.com/CloudLearnersOrg/golib/pkg/log"
	"github.com/gin-gonic/gin"
)

// Middleware returns a security middleware with default configuration.
// This is a convenience function for quickly adding security headers
// without custom configuration.
func Middleware() gin.HandlerFunc {
	return New(DefaultConfig())
}

// New creates a security middleware with the provided configuration.
// It performs host validation if ExpectedHosts is set and applies
// security headers according to the configuration.
func New(config *Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// SSRF Protection - Host Validation
		if len(config.ExpectedHosts) > 0 {
			// Validate the host
			if !isValidHost(c.Request.Host, config.ExpectedHosts) {
				log.Warnf("Host validation failed", map[string]any{
					"requestHost":   c.Request.Host,
					"expectedHosts": config.ExpectedHosts,
				})

				ginhttp.StatusForbidden(c, "Invalid Host Header", map[string]any{
					"error": "Host header does not match expected hosts",
				})
				return
			}
		}

		// Apply security headers
		applySecurityHeaders(c, config)
		c.Next()
	}
}

// isValidHost checks if the request host matches any of the allowed hosts
// considering both with and without port versions
func isValidHost(requestHost string, allowedHosts []string) bool {
	// Extract hostname without port if present
	requestHostname := extractHostWithoutPort(requestHost)

	for _, allowedHost := range allowedHosts {
		// Check for exact match
		if allowedHost == requestHost {
			return true
		}

		// Check hostname without port
		if allowedHost == requestHostname {
			return true
		}

		// Check if the allowed host's hostname (without port) matches
		allowedHostname := extractHostWithoutPort(allowedHost)
		if allowedHostname == requestHostname {
			return true
		}
	}

	return false
}

// extractHostWithoutPort returns the hostname part without the port
func extractHostWithoutPort(host string) string {
	if idx := strings.Index(host, ":"); idx > 0 {
		return host[:idx]
	}
	return host
}

// applySecurityHeaders adds security headers based on the configuration
func applySecurityHeaders(c *gin.Context, config *Config) {
	if config.XFrameOptions != "" {
		c.Header("X-Frame-Options", config.XFrameOptions)
	}

	if config.ContentSecurityPolicy != "" {
		c.Header("Content-Security-Policy", config.ContentSecurityPolicy)
	}

	if config.EnableXSSProtection {
		c.Header("X-XSS-Protection", "1; mode=block")
	}

	if config.StrictTransportSecurityPolicy != "" {
		c.Header("Strict-Transport-Security", config.StrictTransportSecurityPolicy)
	}

	if config.ReferrerPolicy != "" {
		c.Header("Referrer-Policy", config.ReferrerPolicy)
	}

	if config.XContentTypeOptions != "" {
		c.Header("X-Content-Type-Options", config.XContentTypeOptions)
	}

	if config.PermissionsPolicy != "" {
		c.Header("Permissions-Policy", config.PermissionsPolicy)
	}
}
