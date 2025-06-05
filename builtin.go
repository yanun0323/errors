package errors

import (
	"errors"
	"fmt"
	"slices"
	"strings"
)

var (
	// SkipRuntimeStackTrace is a flag to skip the runtime stack trace
	//
	// It is useful to skip the runtime stack trace when you want to get the error message
	// without the runtime stack trace
	//
	// It is true by default
	SkipRuntimeStackTrace = true
)

// Is represents builtin errors.Is
func Is(err, target error) bool {
	return errors.Is(Unwrap(err), Unwrap(target))
}

// As represents builtin errors.As
func As(err error, target any) bool {
	return errors.As(err, target)
}

// Unwrap represents builtin errors.Unwrap
func Unwrap(err error) error {
	if err == nil {
		return nil
	}

	u, ok := err.(unwrap)
	if ok {
		return u.Unwrap()
	}
	return err
}

type errorString struct {
	message string
}

func (e errorString) Error() string {
	return e.message
}

// New creates a new error with stack trace
func New(text string) Error {
	return newError(text)
}

// Wrap wraps an error and formats using the default formats for its operands and returns the resulting string. Spaces are added between operands when neither is a string.
//
// # It has better performance than Errorf with '%w' verb
//
// is the same as
//
//	errors.Errorf("%w", err)
func Wrap(err error, args ...any) Error {
	var message string
	if len(args) != 0 {
		message = fmt.Sprint(args...)
	}

	return wrap(err, message)
}

// Wrapf wraps an error and formats according to a format specifier and returns the resulting string.
//
// # It has better performance than Errorf with '%w' verb
//
// is the same as
//
//	errors.Errorf("%w", err)
func Wrapf(err error, format string, args ...any) Error {
	if len(args) != 0 {
		return wrap(err, fmt.Sprintf(format, args...))
	}

	return wrap(err, format)
}

// Errorf creates a formatted error, supporting '%w' verb for error wrapping
func Errorf(format string, args ...any) Error {
	if len(args) == 0 {
		return newError(format)
	}

	// Check if args contains error types and format string contains %w
	if strings.Contains(format, "%w") {
		before, _, _ := strings.Cut(format, "%w")
		idx := strings.Count(before, "%")
		if len(args) < idx {
			return New("errors: Errorf format contains more than one '%w' verb")
		}

		if err, ok := args[idx].(error); ok {
			// Replace %w with %v and don't pass error to formatting arguments
			newFormat := strings.Replace(format, "%w", "%s", 1)
			format = newFormat

			if err == nil {
				return newError(fmt.Sprintf(format, args...))
			} else {
				return wrap(err, fmt.Sprintf(format, args...), true)
			}
		}
	}

	return newError(fmt.Sprintf(format, args...))
}

func newError(text string) Error {
	stack := getStack(1)
	lastCaller := frame{}
	if len(stack) != 0 {
		lastCaller = stack[0]
	}

	return &errorStack{
		message:    text,
		cause:      errorString{message: text},
		lastCaller: lastCaller,
		stack:      stack,
		attr:       []attr{},
	}
}

func wrap(err error, message string, ignoreErrorMessage ...bool) Error {
	if err == nil {
		return nil
	}

	var (
		msg        string
		attrs      []attr
		lastCaller frame
		cause      = err
		stack      = getStack(1)
		ignore     = len(ignoreErrorMessage) != 0 && ignoreErrorMessage[0]
	)

	if len(stack) != 0 {
		lastCaller = stack[0]
	}

	if message == "" {
		msg = err.Error()
	} else {
		if ignore {
			msg = message
		} else {
			msg = message + ", err: " + err.Error()
		}
	}

	if err, ok := err.(*errorStack); ok {
		cause = err.cause
		attrs = slices.Clone(err.attr)
		if len(err.stack) > len(stack) {
			stack = err.stack[:]
		}
	}

	return &errorStack{
		message:    msg,
		cause:      cause,
		lastCaller: lastCaller,
		stack:      stack[:],
		attr:       attrs,
	}
}
