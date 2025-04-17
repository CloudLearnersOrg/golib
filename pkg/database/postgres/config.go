package postgres

import (
	"time"

	"github.com/CloudLearnersOrg/golib/pkg/log"
)

func setDefaults(c *Connection) {
	initializeConnectionPool(c)
	setConnectionDefaults(c)
	setPoolDefaults(c.ConnectionPool)
}

func initializeConnectionPool(c *Connection) {
	if c.ConnectionPool == nil {
		log.Warnf("ConnectionPool is nil, using default settings", nil)
		c.ConnectionPool = &ConnectionPool{}
	}
}

func setConnectionDefaults(c *Connection) {
	if c.SSLMode == "" {
		c.SSLMode = "disable"
	}
}

func setPoolDefaults(pool *ConnectionPool) {
	if pool.ValidationQuery == "" {
		pool.ValidationQuery = "SELECT 1"
	}

	if pool.MinPoolSize == 0 {
		pool.MinPoolSize = 2
	}

	if pool.MaxPoolSize == 0 {
		pool.MaxPoolSize = 10
	}

	if pool.MaxConnectionIdleTime == 0 {
		pool.MaxConnectionIdleTime = 30 * time.Second
	}

	if pool.MaxConnectionLifetime == 0 {
		pool.MaxConnectionLifetime = 90 * time.Second
	}

	if pool.ConnectionTimeout == 0 {
		pool.ConnectionTimeout = 5 * time.Second
	}

	if pool.RetryAttempts == 0 {
		pool.RetryAttempts = 3
	}

	if pool.RetryInterval == 0 {
		pool.RetryInterval = 3 * time.Second
	}
}
