// Package http from http/gin/statuses provides a set of utility functions for handling HTTP responses in a Gin web application.
// It offers a standardized way to send JSON responses with proper HTTP status codes, messages, and data.
//
// The package organizes HTTP status handlers into the following categories:
//   - 2xx (Success): OK, Created
//   - 3xx (Redirection): TemporaryRedirect, PermanentRedirect, Found, MovedPermanently
//   - 4xx (Client Errors): BadRequest, Unauthorized, Forbidden, NotFound, etc.
//   - 5xx (Server Errors): InternalServerError, BadGateway, ServiceUnavailable, etc.
//
// All responses follow a consistent JSON structure:
//
//	{
//	    "code": 200,          // HTTP status code
//	    "message": "string",  // Human-readable message
//	    "data": any,         // Optional response data (success responses)
//	    "error": "string"    // Optional error message (error responses)
//	}
//
// Success Response Example (2xx):
//
//	http.StatusOK(ctx, "User retrieved successfully", user)
//	// Returns:
//	// {
//	//     "code": 200,
//	//     "message": "User retrieved successfully",
//	//     "data": {"id": 1, "name": "John"}
//	// }
//
// Error Response Example (4xx/5xx):
//
//	http.StatusNotFound(ctx, "User not found", err)
//	// Returns:
//	// {
//	//     "code": 404,
//	//     "message": "User not found",
//	//     "error": "user with id 1 does not exist"
//	// }
//
// Redirect Example (3xx):
//
//	http.StatusTemporaryRedirect(ctx, "/new-location")
//	// Redirects to the specified location with 307 status code
//
// The package aims to provide:
//   - Consistent response structure across all endpoints
//   - Type-safe status code handling
//   - Separation of concerns by status code category
//   - Easy-to-use interface for common HTTP responses
package http
