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

func (t Template) With(args ...any) Template {
	attrs := make([]attr, 0, len(t.attr)+len(args)/2)
	attrs = append(attrs, t.attr...)
	attrs = append(attrs, makeArgs("", args...)...)

	return Template{
		attr: attrs,
	}
}

func (t Template) New(text string) Error {
	return newError(text, 1, t)
}

func (t Template) Wrap(err error, args ...any) Error {
	var message string
	if len(args) != 0 {
		message = fmt.Sprint(args...)
	}

	return wrap(err, message, 1, false, t)
}

func (t Template) Wrapf(err error, format string, args ...any) Error {
	if len(args) != 0 {
		return wrap(err, fmt.Sprintf(format, args...), 1, false, t)
	}

	return wrap(err, format, 1, false, t)
}

func (t Template) Errorf(format string, args ...any) Error {
	return errorf(t, format, args...)
}

func (t Template) Attrs(lastCaller frame) []attr {
	attrs := slices.Clone(t.attr)
	for i := range attrs {
		attrs[i].Function = lastCaller.Function
	}

	return attrs
}
