package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// NewRedisClient creates a new Redis client with the given configuration
func NewRedisClient(config Connection) (*redis.Client, error) {
	opts := &redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password: config.Password,
		DB:       config.Database,
	}

	if config.ConnectionPool != nil {
		opts.PoolSize = config.ConnectionPool.PoolSize
		opts.MinIdleConns = config.ConnectionPool.MinIdleConns
		opts.MaxRetries = config.ConnectionPool.MaxRetries
		opts.DialTimeout = config.ConnectionPool.ConnectTimeout
		opts.ReadTimeout = config.ConnectionPool.ReadTimeout
		opts.WriteTimeout = config.ConnectionPool.WriteTimeout
		opts.PoolTimeout = config.ConnectionPool.PoolTimeout
	}

	client := redis.NewClient(opts)

	// Test the connection using context
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return client, nil
}
