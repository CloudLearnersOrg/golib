package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Database struct {
	poolDriver *pgxpool.Pool
	config     Connection
}

// NewDatabase creates a new database connection pool using the provided configuration.
func NewDatabase(config Connection) (*Database, error) {
	setDefaults(&config)

	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Host,
		config.Port,
		config.Username,
		config.Password,
		config.Database,
		config.SSLMode,
	)

	pgxConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("parse pgx config: %w", err)
	}

	ctx := context.Background()
	connection, err := RetryConnection(ctx, pgxConfig, config.ConnectionPool.ValidationQuery, config.ConnectionPool.RetryAttempts, config.ConnectionPool.RetryInterval)
	if err != nil {
		return nil, fmt.Errorf("connect to database: %w", err)
	}

	return &Database{
		poolDriver: connection,
		config:     config,
	}, nil
}

func (db *Database) Close() {
	if db.poolDriver != nil {
		db.poolDriver.Close()
	}
}

func (db *Database) Pool() *pgxpool.Pool {
	return db.poolDriver
}
