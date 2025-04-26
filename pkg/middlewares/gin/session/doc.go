// Package session provides middleware for secure session management in Gin applications.
// It offers Redis-backed sessions with configurable options and utility functions
// for creating, validating, and managing user sessions across microservices.
//
// Features:
//   - Redis-backed session storage for distributed environments
//   - Session validation middleware for protected routes
//   - Session creation, destruction, and rotation
//   - UUID compatibility for user identification
//   - Additional session data storage beyond authentication
//   - Configurable cookie settings (domain, secure, SameSite, etc.)
//
// Example Usage:
//
// 1. Basic session setup in a Gin application:
//
//	import (
//	    "net/http"
//
//	    "github.com/CloudLearnersOrg/golib/pkg/middlewares/gin/session"
//	    "github.com/gin-gonic/gin"
//	)
//
//	func main() {
//	    router := gin.Default()
//
//	    // Configure session middleware
//	    sessionConfig := session.SessionConfig{
//	        RedisHost:     "localhost",
//	        RedisPort:     6379,
//	        RedisPassword: "",
//	        SessionSecret: "your-session-secret",
//	        CookieDomain:  ".yourdomain.com",
//	        CookieSecure:  true,
//	        SessionMaxAge: 86400,  // 24 hours
//	        SameSite:      http.SameSiteLaxMode,
//	    }
//
//	    // Apply session middleware
//	    sessionMiddleware, err := session.NewMiddleware(sessionConfig)
//	    if err != nil {
//	        panic("Failed to setup session store: " + err.Error())
//	    }
//	    router.Use(sessionMiddleware)
//
//	    // Routes...
//	    router.Run(":8080")
//	}
//
// 2. Implementing login, logout, and protected routes:
//
//	func setupRoutes(router *gin.Engine) {
//	    // Public routes
//	    router.POST("/login", loginHandler)
//
//	    // Protected routes
//	    auth := router.Group("/api")
//	    auth.Use(session.ValidateSession())
//	    {
//	        auth.GET("/profile", profileHandler)
//	        auth.POST("/logout", logoutHandler)
//	    }
//	}
//
//	func loginHandler(c *gin.Context) {
//	    // Authenticate user (example)
//	    userID := "550e8400-e29b-41d4-a716-446655440000" // Example UUID
//
//	    // Create session
//	    if err := session.CreateSession(c, userID); err != nil {
//	        c.JSON(500, gin.H{"error": "Failed to create session"})
//	        return
//	    }
//
//	    c.JSON(200, gin.H{"message": "Login successful"})
//	}
//
//	func logoutHandler(c *gin.Context) {
//	    if err := session.DestroySession(c); err != nil {
//	        c.JSON(500, gin.H{"error": "Failed to logout"})
//	        return
//	    }
//
//	    c.JSON(200, gin.H{"message": "Logout successful"})
//	}
//
//	func profileHandler(c *gin.Context) {
//	    // Get the UUID of the current user
//	    userID, err := session.GetCurrentUserID(c)
//	    if err != nil {
//	        c.JSON(500, gin.H{"error": "Failed to get user ID"})
//	        return
//	    }
//
//	    // Use userID to fetch and return profile data
//	    c.JSON(200, gin.H{
//	        "user_id": userID,
//	        "profile": "User profile data here",
//	    })
//	}
//
// 3. Storing and retrieving additional session data:
//
//	func storePreferences(c *gin.Context) {
//	    // Store user preferences in session
//	    prefs := map[string]string{
//	        "theme": "dark",
//	        "language": "en",
//	    }
//
//	    if err := session.SetSessionData(c, "preferences", prefs); err != nil {
//	        c.JSON(500, gin.H{"error": "Failed to save preferences"})
//	        return
//	    }
//
//	    c.JSON(200, gin.H{"message": "Preferences saved"})
//	}
//
//	func getPreferences(c *gin.Context) {
//	    // Retrieve user preferences from session
//	    prefs, exists := session.GetSessionData(c, "preferences")
//	    if !exists {
//	        c.JSON(404, gin.H{"error": "No preferences found"})
//	        return
//	    }
//
//	    c.JSON(200, gin.H{"preferences": prefs})
//	}
//
// 4. Using in a microservices architecture:
//
//	// All microservices should use the same Redis configuration
//	// and session secret to share sessions across services
//
//	// Auth Service - handles login
//	func loginHandler(c *gin.Context) {
//	    // ... authenticate user
//	    session.CreateSession(c, userID)
//	}
//
//	// Profile Service - protected endpoint
//	func setupProfileService() {
//	    router := gin.Default()
//
//	    // Same session config as auth service
//	    router.Use(sessionMiddleware)
//
//	    // All routes require valid session
//	    router.Use(session.ValidateSession())
//
//	    router.GET("/profile", getProfileHandler)
//	}
//
// Security Considerations:
//   - Always use HTTPS in production (set CookieSecure: true)
//   - Use a strong, randomly generated session secret
//   - Set appropriate SameSite cookie policy
//   - Consider session rotation for enhanced security
//   - Store minimal data in sessions
//   - Sessions are stored in Redis with authentication data shared across services
//
// Session Helper Functions:
//   - CreateSession: Creates a new authenticated session
//   - DestroySession: Ends a session (logout)
//   - ValidateSession: Middleware to enforce authentication
//   - GetCurrentUserID: Extracts and parses UUID from session
//   - RefreshSession: Updates session expiry time
//   - RotateSessionID: Changes session ID while preserving data
//   - SetSessionData: Stores additional data in session
//   - GetSessionData: Retrieves additional data from session
package session
