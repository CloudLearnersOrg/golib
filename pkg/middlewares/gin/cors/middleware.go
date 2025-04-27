package cors

import (
	"net/http"
	"strings"

	"slices"

	"github.com/gin-gonic/gin"
)

// New returns a new CORS middleware with the provided config
func New(config Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Set origin header
		if len(config.AllowOrigins) > 0 {
			if config.AllowOrigins[0] == "*" {
				c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			}

			if config.AllowOrigins[0] != "*" {
				origin := c.Request.Header.Get("Origin")
				if slices.Contains(config.AllowOrigins, origin) {
					c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
				}
			}
		}

		// Set allow methods header
		if len(config.AllowMethods) > 0 {
			c.Writer.Header().Set("Access-Control-Allow-Methods", strings.Join(config.AllowMethods, ", "))
		}

		// Set allow headers header
		if len(config.AllowHeaders) > 0 {
			c.Writer.Header().Set("Access-Control-Allow-Headers", strings.Join(config.AllowHeaders, ", "))
		}

		// Set expose headers header
		if len(config.ExposeHeaders) > 0 {
			c.Writer.Header().Set("Access-Control-Expose-Headers", strings.Join(config.ExposeHeaders, ", "))
		}

		// Set allow credentials header
		if config.AllowCredentials {
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		// Handle preflight requests
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// Middleware creates a CORS middleware with default configuration
func Middleware() gin.HandlerFunc {
	return New(DefaultConfig())
}
