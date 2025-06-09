package errors

import (
	"errors"
	"fmt"
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
	return newError(text, 1)
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

	return wrap(err, message, 1, false)
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
		return wrap(err, fmt.Sprintf(format, args...), 1, false)
	}

	return wrap(err, format, 1, false)
}

// Errorf creates a formatted error, supporting '%w' verb for error wrapping
func Errorf(format string, args ...any) Error {
	return errorf(NewTemplate(), format, args...)
}

func errorf(template Template, format string, args ...any) Error {
	if len(args) == 0 {
		return newError(format, 2, template)
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
				return newError(fmt.Sprintf(format, args...), 2, template)
			} else {
				return wrap(err, fmt.Sprintf(format, args...), 2, true, template)
			}
		} else {
			return nil
		}
	}

	return newError(fmt.Sprintf(format, args...), 2, template)
}

func newError(text string, ignoreCallStackCount int, tp ...Template) Error {
	stack := getStack(ignoreCallStackCount)
	lastCaller := frame{}
	if len(stack) != 0 {
		lastCaller = stack[0]
	}

	template := NewTemplate()
	if len(tp) != 0 {
		template = tp[0]
	}

	return &errorStack{
		message:    text,
		cause:      errorString{message: text},
		lastCaller: lastCaller,
		stack:      stack,
		attr:       template.Attrs(lastCaller),
	}
}

func wrap(err error, message string, ignoreCallStackCount int, combineStack bool, tp ...Template) Error {
	if err == nil {
		return nil
	}

	var (
		msg        string
		attrs      []attr
		tempAttrs  []attr
		lastCaller frame
		cause      = err
		stack      = getStack(ignoreCallStackCount)
		ignore     = combineStack
		template   = NewTemplate()
	)

	if len(stack) != 0 {
		lastCaller = stack[0]
	}

	if len(tp) != 0 {
		template = tp[0]
		tempAttrs = template.Attrs(lastCaller)
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
		attrs = make([]attr, 0, len(err.attr)+len(tempAttrs))
		attrs = append(attrs, err.attr...)

		if len(err.stack) > len(stack) {
			stack = err.stack[:]
		}
	}

	attrs = append(attrs, tempAttrs...)

	return &errorStack{
		message:    msg,
		cause:      cause,
		lastCaller: lastCaller,
		stack:      stack[:],
		attr:       attrs,
	}
}
