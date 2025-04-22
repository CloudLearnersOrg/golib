// Package postgres provides a set of methods to simplify PostgreSQL database connection management
// and query execution. It includes features for retrying connections, validating queries, and
// managing connection pool settings.
//
// Key Features:
//   - Connection pooling with configurable settings
//   - Automatic retry mechanism for connection establishment
//   - Connection validation through health check queries
//   - Structured logging for connection events
//   - Default configuration management
//
// Basic Usage:
//
//	config := postgres.Connection{
//	    Host:     "localhost",
//	    Port:     5432,
//	    Username: "user",
//	    Password: "pass",
//	    Database: "mydb",
//	    ConnectionPool: &postgres.ConnectionPool{
//	        MinPoolSize: 2,
//	        MaxPoolSize: 10,
//	    },
//	}
//
//	db, err := postgres.NewDatabase(config)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer db.Close()
//
// Connection Pool Configuration:
// The package provides extensive configuration options for the connection pool:
//   - MinPoolSize: Minimum number of connections (default: 2)
//   - MaxPoolSize: Maximum number of connections (default: 10)
//   - MaxConnectionIdleTime: Maximum idle time (default: 30s)
//   - MaxConnectionLifetime: Maximum connection lifetime (default: 90s)
//   - ConnectionTimeout: Connection timeout (default: 5s)
//   - ValidationQuery: Query to validate connections (default: "SELECT 1")
//   - RetryAttempts: Number of connection attempts (default: 3)
//   - RetryInterval: Time between retries (default: 3s)
//
// Retry Mechanism:
// The package implements automatic retry logic for establishing database connections:
//  1. Attempts to establish initial connection
//  2. Validates connection with ping
//  3. Executes validation query
//  4. Retries on failure with configured interval
//  5. Logs each attempt with detailed error information
//
// Example with Custom Pool Settings:
//
//	config := postgres.Connection{
//	    Host:     "localhost",
//	    Port:     5432,
//	    Username: "user",
//	    Password: "pass",
//	    Database: "mydb",
//	    ConnectionPool: &postgres.ConnectionPool{
//	        MinPoolSize:           5,
//	        MaxPoolSize:           20,
//	        MaxConnectionIdleTime: 1 * time.Minute,
//	        ValidationQuery:       "SELECT NOW()",
//	        RetryAttempts:        5,
//	        RetryInterval:        5 * time.Second,
//	    },
//	}
package postgres
