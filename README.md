# Golib

Golib is a golang library that provides common packages for microservices. The goal is to provide a set of packages to simplify the development of microservices in Golang. The library includes packages for authentication, authorization, logging and other common functionalities.

## Installation

To install Golib, you can use the following command:

```bash
go get -u github.com/CloudLearnersOrg/golib
```

## Example Usage

Import the library in your Go application:

```go
import "github.com/CloudLearnersOrg/golib/pkg/log"

func main() {
    log.Infof("Hello, World!", map[string]any{
        "key": "value",
    })
}
```

The documentation for the library is available at [https://pkg.go.dev/github.com/CloudLearnersOrg/golib](https://pkg.go.dev/github.com/CloudLearnersOrg/golib).