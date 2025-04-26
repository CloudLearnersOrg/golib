package csrf

import (
	"net/http"
	"time"

	"github.com/CloudLearnersOrg/golib/pkg/csrf"
	ginhttp "github.com/CloudLearnersOrg/golib/pkg/ginhttp/gin/statuses"
	"github.com/CloudLearnersOrg/golib/pkg/middlewares/gin/session"
	"github.com/gin-gonic/gin"
)

// Middleware creates a CSRF protection middleware for Gin
func New(maxAge time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip for GET, HEAD, OPTIONS - they don't modify state
		if c.Request.Method == http.MethodGet ||
			c.Request.Method == http.MethodHead ||
			c.Request.Method == http.MethodOptions {
			c.Next()
			return
		}

		// Get CSRF secret from session
		secretValue, exists := session.GetSessionData(c, TokenKey)
		if !exists {
			ginhttp.StatusForbidden(c, "CSRF protection error: Missing CSRF secret", nil)
			return
		}
		secret := secretValue.(string)

		// Get token from header or form
		token := c.GetHeader(TokenHeader)
		if token == "" {
			token = c.PostForm(TokenFormField)
		}

		if token == "" {
			ginhttp.StatusForbidden(c, "CSRF protection error: Missing CSRF token", nil)
			return
		}

		if err := csrf.VerifyCSRFToken(token, secret, maxAge); err != nil {
			ginhttp.StatusForbidden(c, "CSRF protection error: Invalid CSRF token", err)
			return
		}

		c.Next()
	}
}

