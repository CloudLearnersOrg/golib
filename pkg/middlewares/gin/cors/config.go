package cors

// Config represents the configuration for CORS middleware
type Config struct {
	// AllowOrigins is a list of origins a cross-domain request can be executed from.
	// If the special "*" value is present in the list, all origins will be allowed.
	// Default value is ["*"]
	AllowOrigins []string

	// AllowMethods is a list of methods the client is allowed to use with
	// cross-domain requests. Default value is simple methods (GET, POST, PUT, PATCH, DELETE, OPTIONS)
	AllowMethods []string

	// AllowHeaders is list of non-simple headers the client is allowed to use with
	// cross-domain requests.
	AllowHeaders []string

	// ExposeHeaders indicates which headers are safe to expose to the API
	ExposeHeaders []string

	// AllowCredentials indicates whether the request can include user credentials like
	// cookies, HTTP authentication or client side SSL certificates.
	AllowCredentials bool
}

// DefaultConfig returns a generic default configuration
func DefaultConfig() Config {
	return Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		AllowCredentials: false,
	}
}
