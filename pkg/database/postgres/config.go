package postgres

import (
	"time"

	"github.com/CloudLearnersOrg/golib/pkg/logger"
)

func setDefaults(c *Connection) {
	defaultValidationQuery := "SELECT 1"
	defaultConfigSSLMode := "disable"

	// Initialize ConnectionPool if nil to avoid panic
	if c.ConnectionPool == nil {
		logger.Warnf("ConnectionPool is nil, using default settings", nil)
		c.ConnectionPool = &ConnectionPool{}
	}

	if c.ConnectionPool.ValidationQuery == "" {
		c.ConnectionPool.ValidationQuery = defaultValidationQuery
	}

	if c.SSLMode == "" {
		c.SSLMode = defaultConfigSSLMode
	}

	if c.ConnectionPool.MinPoolSize == 0 {
		c.ConnectionPool.MinPoolSize = 2
	}

	if c.ConnectionPool.MaxPoolSize == 0 {
		c.ConnectionPool.MaxPoolSize = 10
	}

	if c.ConnectionPool.MaxConnectionIdleTime == 0 {
		c.ConnectionPool.MaxConnectionIdleTime = 30 * time.Second
	}

	if c.ConnectionPool.MaxConnectionLifetime == 0 {
		c.ConnectionPool.MaxConnectionLifetime = 90 * time.Second
	}

	if c.ConnectionPool.ConnectionTimeout == 0 {
		c.ConnectionPool.ConnectionTimeout = 5 * time.Second
	}

	if c.ConnectionPool.RetryAttempts == 0 {
		c.ConnectionPool.RetryAttempts = 3
	}

	if c.ConnectionPool.RetryInterval == 0 {
		c.ConnectionPool.RetryInterval = 3 * time.Second
	}
}
