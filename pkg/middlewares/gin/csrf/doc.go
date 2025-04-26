// Package csrf provides middleware for Cross-Site Request Forgery protection in Gin applications.
// It leverages the core CSRF package from github.com/CloudLearnersOrg/golib/pkg/csrf
// and integrates with the session management middleware to provide a complete CSRF
// protection solution.
//
// Features:
//   - Middleware for automatic CSRF protection on state-changing requests (non-GET)
//   - Integration with session management for storing CSRF secrets
//   - Token generation for client-side use
//   - Protection for POST, PUT, PATCH, DELETE operations
//   - Multiple token delivery options (header or form field)
//
// Example Usage:
//
// 1. Basic CSRF protection for a Gin application:
//
//	import (
//	    "time"
//
//	    "github.com/CloudLearnersOrg/golib/pkg/middlewares/gin/csrf"
//	    "github.com/CloudLearnersOrg/golib/pkg/middlewares/gin/session"
//	    "github.com/gin-gonic/gin"
//	)
//
//	func main() {
//	    router := gin.Default()
//
//	    // Configure and apply session middleware first
//	    sessionConfig := session.SessionConfig{
//	        RedisHost:     "localhost",
//	        RedisPort:     6379,
//	        SessionSecret: "your-session-secret",
//	        // Additional session configuration...
//	    }
//	    sessionMiddleware, _ := session.NewMiddleware(sessionConfig)
//	    router.Use(sessionMiddleware)
//
//	    // Apply CSRF middleware with 1-hour token expiration
//	    router.Use(csrf.New(1 * time.Hour))
//
//	    // Routes...
//	    router.Run(":8080")
//	}
//
// 2. Creating an endpoint to get CSRF tokens for single-page applications:
//
//	func setupRoutes(router *gin.Engine) {
//	    // Public routes
//	    router.GET("/", publicHandler)
//
//	    // Protected routes
//	    auth := router.Group("/api")
//	    auth.Use(session.ValidateSession())
//	    {
//	        // Endpoint to get a fresh CSRF token
//	        auth.GET("/csrf-token", getCsrfToken)
//
//	        // Protected actions that need CSRF validation
//	        auth.POST("/update-profile", updateProfile)
//	    }
//	}
//
//	func getCsrfToken(c *gin.Context) {
//	    token, err := csrf.GetToken(c)
//	    if err != nil {
//	        c.JSON(500, gin.H{"error": "Failed to generate token"})
//	        return
//	    }
//	    c.JSON(200, gin.H{"token": token})
//	}
//
// 3. Client-side usage with JavaScript fetch API:
//
//	// JavaScript example (fetch API)
//	async function fetchCsrfToken() {
//	    const response = await fetch('/api/csrf-token', {
//	        credentials: 'include' // Important for cookies
//	    });
//	    const data = await response.json();
//	    return data.token;
//	}
//
//	async function updateProfile(profileData) {
//	    const token = await fetchCsrfToken();
//
//	    const response = await fetch('/api/update-profile', {
//	        method: 'POST',
//	        headers: {
//	            'Content-Type': 'application/json',
//	            'X-CSRF-Token': token
//	        },
//	        credentials: 'include',
//	        body: JSON.stringify(profileData)
//	    });
//
//	    return await response.json();
//	}
//
// Security Considerations:
//   - Always use HTTPS in production environments
//   - CSRF protection works with session-based authentication
//   - The middleware automatically skips GET, HEAD, and OPTIONS requests
//   - Tokens expire after the configured maxAge duration
//   - Tokens can be provided in either an X-CSRF-Token header or csrf_token form field
//
// Constants:
//   - TokenHeader: "X-CSRF-Token" - The HTTP header for the CSRF token
//   - TokenFormField: "csrf_token" - The form field name for the CSRF token
//   - TokenKey: "csrf_secret" - The session key for storing the CSRF secret
package csrf
