# Errors

A lightweight Go errors package with stack tracing and structured fields.

[![Go Reference](https://pkg.go.dev/badge/github.com/yanun0323/errors.svg)](https://pkg.go.dev/github.com/yanun0323/errors)
[![Go Report Card](https://goreportcard.com/badge/github.com/yanun0323/errors)](https://goreportcard.com/report/github.com/yanun0323/errors)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.21-blue)](https://golang.org/dl/)

## Features

- ✅ **Standard library compatible**: Drop-in replacement for `errors` package
- 🔍 **Automatic stack tracing**: Captures call stack when errors are created
- 📊 **Structured fields**: Add key-value pairs with `With()` method
- 🎨 **Multiple formats**: Text, JSON, and colorized output
- 🔗 **Error wrapping**: Full support for `%w` verb and error chains
- ⚡ **High performance**: Efficient implementation with object pooling
- 🔌 **Logs integration**: Native support for [github.com/yanun0323/logs](https://github.com/yanun0323/logs) package

> ⚠️ **Caution**: using `fmt.Errorf` to wrap errors is not compatible with `errors.Is` and `errors.As` methods.

## Installation

```bash
go get github.com/yanun0323/errors
```

## Requirements

- Go 1.21+

## Quick Start

```go
import "github.com/yanun0323/errors"

// Create error with fields
err := errors.New("connection failed").
    With(
        "host", "localhost",
        "port", 5432,
    )

// Error wrapping
wrapped := errors.Errorf("database error: %w", err)
// or
wrapped = errors.Wrap(err)
wrapped = errors.Wrap(err, "database error")
wrapped = errors.Wrapf(err, "database %s error", database)

// Format output
println(err.Error())                        // Basic message

errors.Format(err)                          // Text with stack trace
errors.FormatJson(err)                      // JSON text with stack trace
errors.FormatColorized(err)                 // Colorized text with stack trace

fmt.Printf("%s\n", err)                     // Basic message
fmt.Printf("%v\n", err)                     // Formatted text with stack trace
fmt.Printf("%+v\n", err)                    // Formatted & colorized text with stack trace
fmt.Printf("%#v\n", err)                    // JSON text with stack trace

// Error Template
errTmp := errors.NewTemplate("service", "user-service")
errTmpDB := errTmp.With("component", "database")
err1 := errTmpDB.New("connection established")
err2 := errTmpDB.Errorf("query timeout after %d seconds", 30)
err3 := errTmpDB.Wrapf(originalErr, "database operation %d times", 3)

// Logs Package Integration
import "github.com/yanun0323/logs"
logs.WithError(err1).Error("connection established")
logs.WithError(err2).Error("query timeout")
logs.WithError(err3).Error("database operation")
```

## API

### Creating Errors

```go
errors.New(text string) Error
errors.Errorf(format string, args ...any) Error
errors.Wrap(err error, args ...any) Error
errors.Wrapf(err error, format string, args ...any) Error
```

### Template

Create error templates with predefined attributes for reuse:

```go
errors.NewTemplate(args ...any) Template
```

#### Template Methods

```go
template.With(args ...any) Template             // Add more attributes (chainable)
template.New(text string) Error                 // Create error with template attributes
template.Wrap(err error, args ...any) Error     // Wrap error with template attributes
template.Wrapf(err error, format string, args ...any) Error  // Wrap error with formatted message
template.Errorf(format string, args ...any) Error           // Create formatted error
```

### Error Methods

```go
err.Error() string                          // Standard error message
err.With(args ...any) Error                 // Add fields (chainable)
```

### Standard Functions

```go
errors.Is(err, target error) bool
errors.As(err error, target any) bool
errors.Unwrap(err error) error
```

### Formatting Functions

```go
errors.Format(err error) string             // Text with stack trace
errors.FormatColorized(err error) string    // Colorized text with stack trace
errors.FormatJson(err error) string         // JSON text with stack trace
```

### Logs Package Integration

This package interoperates with the [github.com/yanun0323/logs](https://github.com/yanun0323/logs) package.

```go
logger := logs.Default()

err := errors.New("database connection failed").
    With(
        "host", "localhost",
        "port", 5432,
    )

logger.WithError(err).Error("Operation error")
```

When using with the `logs` package, errors created by this package can be directly passed to log functions and will automatically extract structured fields and stack traces.

## Examples

### Basic Usage

```go
err := errors.New("validation failed").
    With(
        "field", "email",
        "value", "invalid@",
        "rule", "email_format",
    )
```

### Error Wrapping

```go
original := errors.New("network timeout")
wrapped := errors.Errorf("failed to fetch user: %w", original)
```

### Template Usage

```go
// Create a template with common attributes
errTmp := errors.NewTemplate("service", "user-service", "version", "1.0.0")

// Add more attributes to the template
errTmpDB := errTmp.With("component", "database", "host", "localhost")

// Create errors using the template
err1 := errTmpDB.New("connection failed")
err2 := errTmpDB.Errorf("query timeout after %d seconds", 30)
err3 := errTmpDB.Wrap(originalErr, "database operation failed")
```

### JSON Output

```go
err := errors.New("process failed").With("pid", 1234)
fmt.Printf("%#v\n", err)
// Outputs structured JSON with message, fields, and stack trace
```

### Logs Package Integration

```go
import (
    "github.com/yanun0323/errors"
    "github.com/yanun0323/logs"
)

// Create an error with structured fields
err := errors.New("database connection failed").
    With(
        "host", "localhost",
        "port", 5432,
        "timeout", "30s",
    )

// Pass directly to logs package
logs.Error("Operation failed", err)
// The logs package will automatically extract:
// - Error message
// - Structured attributes (host, port, timeout)
// - Stack trace information
```

### Output Formats

#### Text

```
error:
    process user, err: user validation failed, err: root: user not found
cause:
    user not found
field:
    validateUser:
        user_id: 0
        table: users
        func: validateUser
    processUser:
        func: processUser
    handleRequest:
        host: db.example.com
        port: 5432
        timeout: 30s
        func: handleRequest
stack:
    validateUser:
        /Users/Shared/Project/personal/go/errors/example/main.go:58 in validateUser
    processUser:
        /Users/Shared/Project/personal/go/errors/example/main.go:47 in processUser
    handleRequest:
        /Users/Shared/Project/personal/go/errors/example/main.go:34 in handleRequest
    main:
        /Users/Shared/Project/personal/go/errors/example/main.go:18 in main
```

#### Colorized

![Colorized](https://raw.githubusercontent.com/yanun0323/assets/refs/heads/master/errors.colorized.png)

#### JSON

```json
{
  "cause": "user not found",
  "error": "process user, err: user validation failed, err: root: user not found",
  "field": [
    {
      "function": "validateUser",
      "key": "user_id",
      "value": 0
    }
    // ...
  ],
  "stack": [
    {
      "file": "/Users/Shared/Project/personal/go/errors/example/main.go",
      "function": "validateUser",
      "line": "58"
    }
    // ...
  ]
}
```

## Important Notes

⚠️ **Do not use `fmt.Errorf`**

> Use `errors.New` or `errors.Errorf` instead for proper compatibility with `errors.Is` and `errors.As`.

## License

[MIT License](LICENSE)
