package postgres

import (
	"testing"
	"time"
)

func TestSetDefaults(t *testing.T) {
	// Test cases
	tests := []struct {
		name           string
		input          Connection
		expectedOutput Connection
	}{
		{
			name:  "empty connection config with nil ConnectionPool",
			input: Connection{},
			expectedOutput: Connection{
				SSLMode:        "disable",
				ConnectionPool: nil, // Will be initialized in setDefaults
			},
		},
		{
			name: "partial connection config",
			input: Connection{
				Host:    "localhost",
				SSLMode: "require",
				ConnectionPool: &ConnectionPool{
					MinPoolSize: 5,
					// Other fields empty
				},
			},
			expectedOutput: Connection{
				Host:    "localhost",
				SSLMode: "require", // Should not be changed
				ConnectionPool: &ConnectionPool{
					MinPoolSize:           5,  // Should not be changed
					MaxPoolSize:           10, // Default
					MaxConnectionIdleTime: 30 * time.Second,
					MaxConnectionLifetime: 90 * time.Second,
					ConnectionTimeout:     5 * time.Second,
					ValidationQuery:       "SELECT 1",
					RetryAttempts:         3,
					RetryInterval:         3 * time.Second,
				},
			},
		},
		{
			name: "complete connection config",
			input: Connection{
				Host:     "localhost",
				Port:     5432,
				Username: "user",
				Password: "pass",
				Database: "db",
				SSLMode:  "verify-full",
				ConnectionPool: &ConnectionPool{
					MinPoolSize:           10,
					MaxPoolSize:           50,
					MaxConnectionIdleTime: 1 * time.Minute,
					MaxConnectionLifetime: 5 * time.Minute,
					ConnectionTimeout:     10 * time.Second,
					ValidationQuery:       "SELECT now()",
					RetryAttempts:         5,
					RetryInterval:         2 * time.Second,
				},
			},
			// Expected output should be the same as input for complete config
			expectedOutput: Connection{
				Host:     "localhost",
				Port:     5432,
				Username: "user",
				Password: "pass",
				Database: "db",
				SSLMode:  "verify-full",
				ConnectionPool: &ConnectionPool{
					MinPoolSize:           10,
					MaxPoolSize:           50,
					MaxConnectionIdleTime: 1 * time.Minute,
					MaxConnectionLifetime: 5 * time.Minute,
					ConnectionTimeout:     10 * time.Second,
					ValidationQuery:       "SELECT now()",
					RetryAttempts:         5,
					RetryInterval:         2 * time.Second,
				},
			},
		},
	}

	// Run tests
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Copy the input to avoid modifying the test case
			input := tc.input

			// For the case with nil ConnectionPool, we need to handle expected results differently
			isNilConnectionPool := input.ConnectionPool == nil

			// Call the function under test
			setDefaults(&input)

			// Check SSLMode
			if input.SSLMode != tc.expectedOutput.SSLMode {
				t.Errorf("SSLMode: expected %s, got %s", tc.expectedOutput.SSLMode, input.SSLMode)
			}

			// Check if ConnectionPool was properly initialized for nil case
			if isNilConnectionPool {
				if input.ConnectionPool == nil {
					t.Errorf("Expected ConnectionPool to be initialized, but got nil")
					return
				}

				// Check default values for initialized ConnectionPool
				if input.ConnectionPool.MinPoolSize != 2 {
					t.Errorf("MinPoolSize: expected 2, got %d", input.ConnectionPool.MinPoolSize)
				}
				if input.ConnectionPool.MaxPoolSize != 10 {
					t.Errorf("MaxPoolSize: expected 10, got %d", input.ConnectionPool.MaxPoolSize)
				}
				return
			}

			// For non-nil ConnectionPool, check all fields
			cp := input.ConnectionPool
			expected := tc.expectedOutput.ConnectionPool

			if cp.MinPoolSize != expected.MinPoolSize {
				t.Errorf("MinPoolSize: expected %d, got %d", expected.MinPoolSize, cp.MinPoolSize)
			}
			if cp.MaxPoolSize != expected.MaxPoolSize {
				t.Errorf("MaxPoolSize: expected %d, got %d", expected.MaxPoolSize, cp.MaxPoolSize)
			}
			if cp.MaxConnectionIdleTime != expected.MaxConnectionIdleTime {
				t.Errorf("MaxConnectionIdleTime: expected %v, got %v",
					expected.MaxConnectionIdleTime, cp.MaxConnectionIdleTime)
			}
			if cp.MaxConnectionLifetime != expected.MaxConnectionLifetime {
				t.Errorf("MaxConnectionLifetime: expected %v, got %v",
					expected.MaxConnectionLifetime, cp.MaxConnectionLifetime)
			}
			if cp.ConnectionTimeout != expected.ConnectionTimeout {
				t.Errorf("ConnectionTimeout: expected %v, got %v",
					expected.ConnectionTimeout, cp.ConnectionTimeout)
			}
			if cp.ValidationQuery != expected.ValidationQuery {
				t.Errorf("ValidationQuery: expected %s, got %s",
					expected.ValidationQuery, cp.ValidationQuery)
			}
			if cp.RetryAttempts != expected.RetryAttempts {
				t.Errorf("RetryAttempts: expected %d, got %d",
					expected.RetryAttempts, cp.RetryAttempts)
			}
			if cp.RetryInterval != expected.RetryInterval {
				t.Errorf("RetryInterval: expected %v, got %v",
					expected.RetryInterval, cp.RetryInterval)
			}
		})
	}
}
