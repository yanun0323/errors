package errors

import (
	"fmt"
	"slices"
)

// Template is a template for creating errors. It contains args that can be used to create an error.
type Template struct {
	attr []attr
}

// NewTemplate creates a new Template.
func NewTemplate(args ...any) Template {
	return Template{
		attr: makeArgs("", args...),
	}
}

// With creates a new Template by appending additional attributes to the existing ones.
// It returns a new Template instance without modifying the original one.
func (t Template) With(args ...any) Template {
	attrs := make([]attr, 0, len(t.attr)+len(args)/2)
	attrs = append(attrs, t.attr...)
	attrs = append(attrs, makeArgs("", args...)...)

	return Template{
		attr: attrs,
	}
}

// WithMap creates a new Template by appending additional attributes to the existing ones.
// It returns a new Template instance without modifying the original one.
func (t Template) WithMap(m map[string]any) Template {
	attrs := make([]attr, 0, len(t.attr)+len(m))
	attrs = append(attrs, t.attr...)
	for k, v := range m {
		attrs = append(attrs, attr{
			Function: "",
			Key:      k,
			Value:    v,
		})
	}

	return Template{
		attr: attrs,
	}
}

// New creates a new Error with the given text message and the template's attributes.
func (t Template) New(text string) Error {
	return newError(text, 1, t)
}

// Wrap wraps an existing error with optional additional message arguments.
// If args are provided, they will be concatenated as the wrap message.
func (t Template) Wrap(err error, args ...any) Error {
	var message string
	if len(args) != 0 {
		message = fmt.Sprint(args...)
	}

	return wrap(err, message, 1, false, t)
}

// Wrapf wraps an existing error with a formatted message using fmt.Sprintf.
// If no args are provided, the format string is used as-is.
func (t Template) Wrapf(err error, format string, args ...any) Error {
	if len(args) != 0 {
		return wrap(err, fmt.Sprintf(format, args...), 1, false, t)
	}

	return wrap(err, format, 1, false, t)
}

// Errorf creates a new formatted Error using the template's attributes.
// It formats the message using fmt.Sprintf with the provided format and args.
func (t Template) Errorf(format string, args ...any) Error {
	return errorf(t, format, args...)
}

// Attrs returns a copy of the template's attributes with the Function field
// set to the provided lastCaller frame's Function value.
func (t Template) Attrs(lastCaller frame) []attr {
	attrs := slices.Clone(t.attr)
	for i := range attrs {
		attrs[i].Function = lastCaller.Function
	}

	return attrs
}

// Clone creates a new Template with the same attributes.
func (t Template) Clone() Template {
	return Template{
		attr: slices.Clone(t.attr),
	}
}
