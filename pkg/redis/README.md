# Redis Client Package

A Go package that provides a simple way to create and configure Redis clients with connection pooling support.

## Installation

```go
go get github.com/CloudLearnersOrg/golib/pkg/redis
```

## Usage

### Configuration Options

#### Connection
- `Host`: Redis server hostname or IP address
- `Port`: Redis server port
- `Password`: Authentication password (if required)
- `Database`: Redis database number
- `ConnectionPool`: Optional connection pool configuration

#### ConnectionPool
- `PoolSize`: Maximum number of connections in the pool
- `MinIdleConns`: Minimum number of idle connections maintained in the pool
- `MaxRetries`: Maximum number of retries before giving up
- `ConnectTimeout`: Timeout for establishing new connections
- `ReadTimeout`: Timeout for socket reads
- `WriteTimeout`: Timeout for socket writes
- `PoolTimeout`: Timeout for getting connection from the pool
- `IdleTimeout`: Maximum amount of time a connection can be idle
- `MaxConnAge`: Maximum age of a connection


### Basic Connection

Create a basic Redis client with minimal configuration:

```go
config := redis.Connection{
    Host:     "localhost",
    Port:     6379,
    Password: "",  // No password
    Database: 0,   // Default database
}

client, err := redis.NewRedisClient(config)
if err != nil {
    log.Fatalf("Failed to connect to Redis: %v", err)
}
defer client.Close()

// Use the client
err = client.Set("key", "value", 0).Err()
```

### Connection with Pool Configuration

Create a Redis client with connection pooling:

```go
config := redis.Connection{
    Host:     "redis.example.com",
    Port:     6379,
    Password: "secret",
    Database: 1,
    ConnectionPool: &redis.ConnectionPool{
        PoolSize:       10,
        MinIdleConns:   2,
        MaxRetries:     3,
        ConnectTimeout: 5 * time.Second,
        ReadTimeout:    3 * time.Second,
        WriteTimeout:   3 * time.Second,
        PoolTimeout:    4 * time.Second,
        IdleTimeout:    5 * time.Minute,
        MaxConnAge:     30 * time.Minute,
    },
}

client, err := redis.NewRedisClient(config)
if err != nil {
    log.Fatalf("Failed to connect to Redis: %v", err)
}
defer client.Close()
```

## Example with Error Handling

```go
config := redis.Connection{
    Host: "localhost",
    Port: 6379,
    ConnectionPool: &redis.ConnectionPool{
        PoolSize:       5,
        ConnectTimeout: 2 * time.Second,
        ReadTimeout:    2 * time.Second,
    },
}

client, err := redis.NewRedisClient(config)
if err != nil {
    log.Fatalf("Failed to create Redis client: %v", err)
}
defer client.Close()

// Perform Redis operations
err = client.Set("mykey", "myvalue", 24*time.Hour).Err()
if err != nil {
    log.Printf("Failed to set key: %v", err)
}

val, err := client.Get("mykey").Result()
if err != nil {
    if err == redis.Nil {
        log.Println("Key does not exist")
    } else {
        log.Printf("Failed to get key: %v", err)
    }
}
```

## Note

This package uses the [go-redis/redis](https://github.com/go-redis/redis) library internally. For detailed Redis operations and API documentation, please refer to the [go-redis documentation](https://pkg.go.dev/github.com/go-redis/redis).