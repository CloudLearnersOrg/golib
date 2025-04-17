package redis

import (
	"time"
)

// Connection holds Redis connection configuration
type Connection struct {
	Host           string
	Port           int
	Password       string
	Database       int
	ConnectionPool *ConnectionPool
}

// ConnectionPool holds Redis connection pool configuration
type ConnectionPool struct {
	PoolSize       int
	MinIdleConns   int
	MaxRetries     int
	ConnectTimeout time.Duration
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	PoolTimeout    time.Duration
	IdleTimeout    time.Duration
	MaxConnAge     time.Duration
}
