# Middleware loggers for incoming and outgoing requests

Middleware loggers for incoming and outgoing requests are used to capture the requests and responses and log them down in JSON format.

## Usage

### `net/http` based example

```
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	middleware "github.com/CloudLearnersOrg/golib/pkg/middlewares/logger/http"
)

// Sample response struct
type Response struct {
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
	Status    string    `json:"status"`
}

// Sample request struct
type Request struct {
	Name string `json:"name"`
}

// Define a handler that processes incoming requests and makes outgoing requests
func apiHandler(w http.ResponseWriter, r *http.Request) {
	// Parse incoming request body if available
	if r.Method == "POST" {
		var req Request
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		log.Printf("Received request with name: %s", req.Name)
	}

	// Create an outgoing logger for external API calls
	outgoingLogger := middleware.NewOutgoingLogger(nil).
		OutgoingWithRequestBody().
		OutgoingWithResponseBody().
		OutgoingWithHeaders()

	// Make an outgoing request to JSONPlaceholder API
	// Use the request context to propagate the trace ID
	outgoingReq, _ := http.NewRequestWithContext(
		r.Context(),
		"GET",
		"https://jsonplaceholder.typicode.com/todos/1",
		nil,
	)

	// Perform the outgoing request with logging
	resp, err := outgoingLogger.Do(outgoingReq)
	if err != nil {
		http.Error(w, "Error calling external API", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Prepare our own response
	response := Response{
		Message:   "Hello from the middleware test server!",
		Timestamp: time.Now(),
		Status:    "success",
	}

	// Set response headers
	w.Header().Set("Content-Type", "application/json")

	// Use different status codes based on the request path to test status code logging
	path := r.URL.Path
	switch path {
	case "/api/error":
		w.WriteHeader(http.StatusInternalServerError)
		response.Status = "error"
	case "/api/notfound":
		w.WriteHeader(http.StatusNotFound)
		response.Status = "not_found"
	default:
		w.WriteHeader(http.StatusOK)
	}

	// Send the response
	json.NewEncoder(w).Encode(response)
}

func main() {
	// Create and configure the incoming logger middleware
	incomingLogger := middleware.IncomingLogger().
		IncomingWithRequestBody().          // Log request bodies
		IncomingWithResponseBody().         // Log response bodies
		IncomingWithHeaders().              // Log HTTP headers
		WithField("service", "logger-test") // Add a custom field

	// Create a router (using standard http.ServeMux for simplicity)
	mux := http.NewServeMux()

	// Register our test handler wrapped with the logger middleware
	// This will log all requests to these endpoints
	mux.Handle("/api/test", incomingLogger.Handler(http.HandlerFunc(apiHandler)))
	mux.Handle("/api/error", incomingLogger.Handler(http.HandlerFunc(apiHandler)))
	mux.Handle("/api/notfound", incomingLogger.Handler(http.HandlerFunc(apiHandler)))

	// Start the server
	serverAddr := ":8080"
	fmt.Printf("Starting server on %s...\n", serverAddr)
	fmt.Println("Try the following endpoints:")
	fmt.Println("- GET http://localhost:8080/api/test")
	fmt.Println("- GET http://localhost:8080/api/error")
	fmt.Println("- GET http://localhost:8080/api/notfound")
	fmt.Println("- POST http://localhost:8080/api/test with JSON body: {\"name\": \"YourName\"}")

	// Start the HTTP server
	log.Fatal(http.ListenAndServe(serverAddr, mux))
}
```

### `gin` based example

```
package main

import (
    "fmt"
    "log"
    "net/http"
    "time"

    middleware "github.com/CloudLearnersOrg/golib/pkg/middlewares/logger/gin"
    "github.com/gin-gonic/gin"
)

// Response represents the API response structure
type Response struct {
    Message   string    `json:"message"`
    Timestamp time.Time `json:"timestamp"`
    Status    string    `json:"status"`
}

// Request represents the API request structure
type Request struct {
    Name string `json:"name"`
}

// apiHandler processes incoming requests and makes outgoing requests
func apiHandler(c *gin.Context) {
    // Parse incoming request body if POST
    if c.Request.Method == "POST" {
        var req Request
        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
            return
        }
        log.Printf("Received request with name: %s", req.Name)
    }

    // Create an outgoing logger for external API calls
    outgoingLogger := middleware.NewOutgoingLogger(nil).
        OutgoingWithRequestBody().
        OutgoingWithResponseBody().
        OutgoingWithHeaders()

    // Make an outgoing request to JSONPlaceholder API
    outgoingReq, _ := http.NewRequest(
        "GET",
        "https://jsonplaceholder.typicode.com/todos/1",
        nil,
    )

    // Perform the outgoing request with logging
    resp, err := outgoingLogger.Do(c, outgoingReq)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error calling external API"})
        return
    }
    defer resp.Body.Close()

    // Prepare response
    response := Response{
        Message:   "Hello from the Gin middleware test server!",
        Timestamp: time.Now(),
        Status:    "success",
    }

    // Use different status codes based on the request path
    switch c.FullPath() {
    case "/api/error":
        response.Status = "error"
        c.JSON(http.StatusInternalServerError, response)
    case "/api/notfound":
        response.Status = "not_found"
        c.JSON(http.StatusNotFound, response)
    default:
        c.JSON(http.StatusOK, response)
    }
}

func main() {
    // Create and configure the incoming logger middleware
    incomingLogger := middleware.IncomingLogger().
        WithRequestBody().
        WithResponseBody().
        WithHeaders().
        WithField("service", "gin-logger-test")

    // Create a new Gin router with default middleware
    r := gin.New()

    // Use the recovery middleware
    r.Use(gin.Recovery())

    // Add our logger middleware
    r.Use(incomingLogger.Handler())

    // Register API routes
    api := r.Group("/api")
    {
        api.GET("/test", apiHandler)
        api.POST("/test", apiHandler)
        api.GET("/error", apiHandler)
        api.GET("/notfound", apiHandler)
    }

    // Start the server
    serverAddr := ":8080"
    fmt.Printf("Starting Gin server on %s...\n", serverAddr)
    fmt.Println("Try the following endpoints:")
    fmt.Println("- GET http://localhost:8080/api/test")
    fmt.Println("- GET http://localhost:8080/api/error")
    fmt.Println("- GET http://localhost:8080/api/notfound")
    fmt.Println("- POST http://localhost:8080/api/test with JSON body: {\"name\": \"YourName\"}")

    // Start the server
    log.Fatal(r.Run(serverAddr))
}
```