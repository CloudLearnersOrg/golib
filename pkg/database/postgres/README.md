# Postgres

`postgres` package is made to simplify the connection pooling configuration across microservices.

It provides a wrapper around the `pgx` library to initialize the database connection fast.

## Usage

### Configuration Options

#### Connection Parameters
- `Host` - Database server hostname or IP address (default: none, required)
- `Port` - Database server port (default: none, required)
- `Username` - Database user name (default: none, required)
- `Password` - Database password (default: none, required)
- `Database` - Database name to connect to (default: none, required)
- `SSLMode` - SSL mode for connection (default: "disable")

#### Connection Pool Parameters
All pool parameters are optional and have reasonable defaults:

- `MinPoolSize` - Minimum number of connections in the pool (default: 2)
- `MaxPoolSize` - Maximum number of connections in the pool (default: 10)
- `MaxConnectionIdleTime` - Maximum time a connection can be idle (default: 30s)
- `MaxConnectionLifetime` - Maximum lifetime of a connection (default: 90s)
- `ConnectionTimeout` - Maximum time to wait when establishing a connection (default: 5s)
- `ValidationQuery` - Query used to validate connections (default: "SELECT 1")
- `RetryAttempts` - Number of connection retry attempts (default: 3)
- `RetryInterval` - Time to wait between retry attempts (default: 3s)

### Simple setup

```go
func main() {
    // Define basic connection config without explicit connection pooling
    config := postgres.Connection{
        Host:     "localhost",
        Port:     5432,
        Username: "postgres",
        Password: "yourpassword",
        Database: "yourdb",
        SSLMode:  "disable",
        // No ConnectionPool field - will use default values
    }

    // Connect to the database
    db, err := postgres.NewDatabase(config)
    if err != nil {
        log.Fatalf("Failed to connect to PostgreSQL: %v", err)
    }
    defer db.Close()

    // Test the connection
    ctx := context.Background()
    var result int
    err = db.Pool().QueryRow(ctx, "SELECT 1").Scan(&result)
    if err != nil {
        log.Fatalf("Query failed: %v", err)
    }

    fmt.Println("Successfully connected to database with default connection pool settings!")
    fmt.Printf("Query result: %d\n", result)
    
    // Use the connection for a simple query
    var version string
    err = db.Pool().QueryRow(ctx, "SELECT version()").Scan(&version)
    if err != nil {
        log.Fatalf("Error querying PostgreSQL version: %v", err)
    }
    
    fmt.Printf("PostgreSQL server version: %s\n", version)
}
```

### Setup with connection pooling

Setting up a database connection with customized connection pooling:

```go
// Import the package
import (
    "log"
    "time"

    "github.com/your-username/golib/pkg/database/postgresql"
)

// Define your connection configuration
postgresConfig := postgresql.Connection{
    Host:     "localhost",
    Port:     5432,
    Username: "postgres",
    Password: "yourpassword",
    Database: "yourdatabase",
    SSLMode:  "disable",
    ConnectionPool: &postgresql.ConnectionPool{
        MinPoolSize:           5,
        MaxPoolSize:           10,
        MaxConnectionIdleTime: 5 * time.Minute,
        MaxConnectionLifetime: 30 * time.Minute,
        ConnectionTimeout:     5 * time.Second,
        ValidationQuery: "SELECT 1",
        RetryAttempts: 3,
        RetryInterval: 100 * time.Millisecond,
    },
}

// Connect to the database
db, err := postgresql.NewDatabase(postgresConfig)
if err != nil {
    log.Fatalf("Failed to connect to PostgreSQL: %v", err)
}
defer db.Close() // Always close the connection when done
```

## Microservice Example

Using shared pool configuration across different service databases:

```go
// Define common pool settings
sharedPool := &postgresql.ConnectionPool{
    MinPoolSize:           5, 
    MaxPoolSize:           10,
    MaxConnectionIdleTime: 5 * time.Minute,
    ConnectionTimeout:     5 * time.Second,
    ValidationQuery: "SELECT 1",
    RetryAttempts: 3,
    RetryInterval: 100 * time.Millisecond,
}

// Configure user service database
userServiceDB, err := postgresql.NewDatabase(postgresql.Connection{
    Host:           "user-service-db.example.com",
    Port:           5432,
    Username:       "user_service",
    Password:       "user_password",
    Database:       "user_db",
    ConnectionPool: sharedPool,
})
if err != nil {
    log.Fatalf("Failed to connect to user service database: %v", err)
}
defer userServiceDB.Close()

// Configure order service database 
orderServiceDB, err := postgresql.NewDatabase(postgresql.Connection{
    Host:           "order-service-db.example.com",
    Port:           5432,
    Username:       "order_service",
    Password:       "order_password", 
    Database:       "order_db",
    ConnectionPool: sharedPool,
})
if err != nil {
    log.Fatalf("Failed to connect to order service database: %v", err)
}
defer orderServiceDB.Close()
```