# EnvLoader

[![GitHub](https://img.shields.io/github/license/unsafe-risk/envloader?style=for-the-badge)](https://github.com/unsafe-risk/envloader/blob/main/LICENSE)
[![Go Reference](https://img.shields.io/badge/go-reference-%23007d9c?style=for-the-badge&logo=go)](https://pkg.go.dev/gopkg.eu.org/envloader)

EnvLoader is a Go package that loads environment variables from a `.env` file and optionally binds them to a struct.

## Installation

```bash
go get gopkg.eu.org/envloader
```

## Usage

### Loading Environment Variables

```go
package main

import (
    "fmt"

    "gopkg.eu.org/envloader"
)

func main() {
    // Load environment variables from a file.
    err := envloader.LoadEnvFile(".env")
    if err != nil {
        panic(err)
    }

    // Access the environment variables.
    fmt.Println(os.Getenv("API_KEY"))
    fmt.Println(os.Getenv("DATABASE_URL"))
}
```

### Binding Environment Variables to a Struct

```go
package main

import (
    "fmt"

    "gopkg.eu.org/envloader"
)

type Config struct {
    APIKey     string `env:"API_KEY"`
    DatabaseURL string `env:"DATABASE_URL"`
}

func main() {
    // Create a config struct.
    var config Config

    // Load and bind environment variables from a file.
    err := envloader.LoadAndBindEnvFile(".env", &config)
    if err != nil {
        panic(err)
    }

    // Access the config values.
    fmt.Println(config.APIKey)
    fmt.Println(config.DatabaseURL)
}
```

## Struct Tag Options

You can use the `env` struct tag to specify the environment variable name to bind to a struct field:

```go
type Config struct {
    APIKey     string `env:"API_KEY"`
    DatabaseURL string `env:"DATABASE_URL"`
}
```

Additionally, you can mark a field as required by adding the `required` option to the `env` tag. This will raise an error if the environment variable is not set:

```go
type Config struct {
    APIKey     string `env:"API_KEY,required"`
    DatabaseURL string `env:"DATABASE_URL,required"`
}
```

## Contributing

Contributions are welcome! Please open an issue or pull request if you have any suggestions or bug reports.

## License

This package is licensed under the [Unlicense](https://github.com/unsafe-risk/envloader/blob/main/LICENSE).
