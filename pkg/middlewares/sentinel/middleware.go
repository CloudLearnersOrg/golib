package sentinel

import (
	"slices"

	ginhttp "github.com/CloudLearnersOrg/golib/pkg/ginhttp/gin/statuses"
	"github.com/gin-gonic/gin"
)

// New creates a security middleware with the provided configuration.
// It performs host validation if ExpectedHosts is set and applies
// security headers according to the configuration.
func New(config *Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// SSRF Protection - Host Validation
		if len(config.ExpectedHosts) > 0 {
			requestHost := c.Request.Host
			hostValid := slices.Contains(config.ExpectedHosts, requestHost)

			if !hostValid {
				ginhttp.StatusForbidden(c, "Invalid Host Header", map[string]any{
					"error": "Host header does not match expected hosts",
				})
				return
			}
		}

		// Security Headers
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

		c.Next()
	}
}

// Middleware returns a security middleware with default configuration.
// This is a convenience function for quickly adding security headers
// without custom configuration.
func Middleware() gin.HandlerFunc {
	return New(DefaultConfig())
}
