/*
Package cors provides Cross-Origin Resource Sharing (CORS) middleware for the Gin framework.

CORS is a mechanism that allows restricted resources on a web page to be requested from
another domain outside the domain from which the first resource was served. This package
helps implement proper CORS headers for your Gin-based API server.

Basic Usage:

To use the default CORS configuration:

	router := gin.Default()
	router.Use(cors.Middleware())

Custom Configuration:

To customize CORS settings:

	router := gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://example.com", "https://yourdomain.com"}
	corsConfig.AllowCredentials = true
	router.Use(cors.New(corsConfig))

Environment-Based Configuration:

You can also load CORS configuration from environment variables:

	router := gin.Default()
	corsConfig := cors.Config{
	    AllowOrigins:     strings.Split(os.Getenv("CORS_ALLOW_ORIGINS"), ","),
	    AllowMethods:     strings.Split(os.Getenv("CORS_ALLOW_METHODS"), ","),
	    AllowHeaders:     strings.Split(os.Getenv("CORS_ALLOW_HEADERS"), ","),
	    ExposeHeaders:    strings.Split(os.Getenv("CORS_EXPOSE_HEADERS"), ","),
	    AllowCredentials: os.Getenv("CORS_ALLOW_CREDENTIALS") == "true",
	}
	router.Use(cors.New(corsConfig))

Configuration Options:

- AllowOrigins: Which origins can make requests. Use ["*"] to allow any origin.
- AllowMethods: HTTP methods that can be used.
- AllowHeaders: HTTP headers that can be used in requests.
- ExposeHeaders: HTTP headers that can be exposed to browsers.
- AllowCredentials: Whether requests can include credentials (cookies, auth headers).

Security Considerations:

 1. When AllowCredentials is true, you cannot use the wildcard "*" for AllowOrigins.
    You must specify exact origins for security reasons.

2. For production APIs, always specify exact origins rather than using the wildcard "*".

3. Include "X-CSRF-Token" in AllowHeaders if you're implementing CSRF protection.
*/
package cors
