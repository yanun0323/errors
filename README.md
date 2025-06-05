# Errors

A lightweight Go errors package with stack tracing and structured fields.

## Features

- ✅ **Standard library compatible**: Drop-in replacement for `errors` package
- 🔍 **Automatic stack tracing**: Captures call stack when errors are created
- 📊 **Structured fields**: Add key-value pairs with `With()` method
- 🎨 **Multiple formats**: Text, JSON, and colorized output
- 🔗 **Error wrapping**: Full support for `%w` verb and error chains
- ⚡ **High performance**: Efficient implementation with object pooling

> ⚠️ **Caution**: using `fmt.Errorf` to wrap errors is not compatible with `errors.Is` and `errors.As` methods.

## Installation

```bash
go get github.com/yanun0323/errors
```

## Quick Start

```go
import "github.com/yanun0323/errors"

// Create error with fields
err := errors.New("connection failed").
    With("host", "localhost").
    With("port", 5432)

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
fmt.Printf("%v\n", err)                     // Text with stack trace
fmt.Printf("%+v\n", err)                    // Colorized text with stack trace
fmt.Printf("%#v\n", err)                    // JSON text with stack trace
```

## API

### Creating Errors

```go
errors.New(text string) Error
errors.Errorf(format string, args ...any) Error
errors.Wrap(err error, args ...any) Error
errors.Wrapf(err error, format string, args ...any) Error
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

## Examples

### Basic Usage

```go
err := errors.New("validation failed").
    With("field", "email").
    With("value", "invalid@").
    With("rule", "email_format")
```

### Error Wrapping

```go
original := errors.New("network timeout")
wrapped := errors.Errorf("failed to fetch user: %w", original)
```

### JSON Output

```go
err := errors.New("process failed").With("pid", 1234)
fmt.Printf("%#v\n", err)
// Outputs structured JSON with message, fields, and stack trace
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
