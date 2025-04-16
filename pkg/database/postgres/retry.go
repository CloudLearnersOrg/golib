package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/CloudLearnersOrg/golib/pkg/logger"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// RetryConnection attempts to establish and validate a database connection with retry logic
func RetryConnection(ctx context.Context, pgxConfig *pgxpool.Config, validationQuery string, retryAttempts int, retryInterval time.Duration) (*pgxpool.Pool, error) {
	// Setup validation function within the connection config
	pgxConfig.BeforeAcquire = func(ctx context.Context, conn *pgx.Conn) bool {
		var result int
		err := conn.QueryRow(ctx, validationQuery).Scan(&result)
		if err != nil {
			logger.Errorf("validation query failed", map[string]any{
				"error": err,
				"query": validationQuery,
			})

			return false
		}

		return true
	}

	var pool *pgxpool.Pool
	var err error
	for attempt := 0; attempt < retryAttempts; attempt++ {
		if attempt > 0 {
			fmt.Printf("Retrying database connection (attempt %d/%d) after %v\n",
				attempt+1, retryAttempts, retryInterval)
			time.Sleep(retryInterval)
		}

		pool, err = pgxpool.ConnectConfig(ctx, pgxConfig)
		if err == nil {
			// Test the connection with a ping
			if err := pool.Ping(ctx); err != nil {
				pool.Close()
				fmt.Printf("Ping failed during database connection attempt %d/%d: %v", attempt+1, retryAttempts, err)
				continue
			}

			// Additionally run a validation query directly
			var result int
			if err := pool.QueryRow(ctx, validationQuery).Scan(&result); err != nil {
				pool.Close()
				fmt.Printf("Validation query failed during database connection attempt %d/%d: %v", attempt+1, retryAttempts, err)
				continue
			}

			return pool, nil
		}
	}

	return nil, fmt.Errorf("failed to connect to database after %d attempts: %w", retryAttempts, err)
}
