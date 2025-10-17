package errors

// Error is an interface that wraps the error interface and provides additional methods
//
// It also implements the error interface, so it can be used as an error
type Error interface {
	error

	With(args ...any) Error
	WithMap(map[string]any) Error
}

type unwrap interface {
	Unwrap() error
}
