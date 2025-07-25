// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package errors

import (
	"strings"
)

// Join returns an error that wraps the given errors.
// Any nil error values are discarded.
// Join returns nil if every value in errs is nil.
// The error formats as the concatenation of the strings obtained
// by calling the Error method of each element of errs, with a newline
// between each string.
//
// A non-nil error returned by Join implements the Unwrap() []error method.
func Join(errs ...error) error {
	n := 0
	for _, err := range errs {
		if err != nil {
			n++
		}
	}
	if n == 0 {
		return nil
	}
	e := &joinError{
		errs: make([]error, 0, n),
	}
	for _, err := range errs {
		if err != nil {
			e.errs = append(e.errs, err)
		}
	}
	return e
}

type joinError struct {
	errs []error
}

func (e *joinError) Error() string {
	switch len(e.errs) {
	case 0:
		return ""
	case 1:
		return e.errs[0].Error()
	default:
		buf := stringBuilderPool.Get().(*strings.Builder)
		defer stringBuilderPool.Put(buf)
		buf.Reset()

		buf.WriteString(e.errs[0].Error())
		for _, err := range e.errs[1:] {
			buf.WriteByte('\n')
			buf.WriteString(err.Error())
		}
		return buf.String()
	}
}

func (e *joinError) Unwrap() []error {
	return e.errs
}
