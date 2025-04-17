package postgres

import "time"

// Connection holds PostgreSQL connection parameters
type Connection struct {
	Host           string
	Port           int
	Username       string
	Password       string
	Database       string
	SSLMode        string
	ConnectionPool *ConnectionPool
}

// ConnectionPool holds PostgreSQL connection pool configuration
type ConnectionPool struct {
	MinPoolSize           int32
	MaxPoolSize           int32
	MaxConnectionIdleTime time.Duration
	MaxConnectionLifetime time.Duration
	ConnectionTimeout     time.Duration
	ValidationQuery       string
	RetryAttempts         int
	RetryInterval         time.Duration
}
