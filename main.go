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
