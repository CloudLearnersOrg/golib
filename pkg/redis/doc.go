// Package redis provides a Redis client wrapper with enhanced configuration options
// and connection pool management. It simplifies Redis connection setup and provides
// reasonable defaults for connection pooling.
//
// Features:
//   - Simple connection configuration
//   - Connection pool management
//   - Automatic connection testing
//   - Configurable timeouts and retries
//
// Basic Usage:
//
//	config := redis.Connection{
//	    Host:     "localhost",
//	    Port:     6379,
//	    Password: "secret",
//	    Database: 0,
//	}
//
//	client, err := redis.NewRedisClient(config)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer client.Close()
//
// Connection Pool Configuration:
// The package provides extensive configuration options for the connection pool:
//   - PoolSize: Maximum number of connections
//   - MinIdleConns: Minimum number of idle connections
//   - MaxRetries: Maximum number of retries
//   - ConnectTimeout: Connection timeout
//   - ReadTimeout: Read timeout
//   - WriteTimeout: Write timeout
//   - PoolTimeout: Pool timeout
//   - IdleTimeout: Connection idle timeout
//   - MaxConnAge: Maximum connection age
//
// Example with Connection Pool:
//
//	config := redis.Connection{
//	    Host:     "localhost",
//	    Port:     6379,
//	    Password: "secret",
//	    Database: 0,
//	    ConnectionPool: &redis.ConnectionPool{
//	        PoolSize:       10,
//	        MinIdleConns:   2,
//	        MaxRetries:     3,
//	        ConnectTimeout: 5 * time.Second,
//	        ReadTimeout:    3 * time.Second,
//	        WriteTimeout:   3 * time.Second,
//	        PoolTimeout:    4 * time.Second,
//	        IdleTimeout:    300 * time.Second,
//	        MaxConnAge:     3600 * time.Second,
//	    },
//	}
//
// Connection Testing:
// The package automatically tests the connection during initialization
// by sending a PING command to Redis. This ensures the connection is
// valid before returning the client.
package redis
