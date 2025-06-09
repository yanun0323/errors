package logs

// Error is the interface for supporting logs (github.com/yanun0323/logs)
type Error interface {
	Message() string
	Cause() error
	Stack() []any
	Attributes() []any
}

// Frame is the interface for supporting logs (github.com/yanun0323/logs)
type Frame interface {
	Parameters() (file, function, line string)
}

// Attr is the interface for supporting logs (github.com/yanun0323/logs)
type Attr interface {
	Parameters() (key string, value any)
}
