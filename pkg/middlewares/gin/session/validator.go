package session

import (
	ginhttp "github.com/CloudLearnersOrg/golib/pkg/ginhttp/gin/statuses"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func ValidateSession() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userID := session.Get(UserKey)
		if userID == nil {
			ginhttp.StatusUnauthorized(c, "Authentication required.", nil)
			return
		}

		c.Set(UserKey, userID)
		c.Next()
	}
}
