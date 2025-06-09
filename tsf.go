package errors

import (
	"strings"
)

// tsf is the interface for supporting logs (github.com/yanun0323/logs)
type tsf interface {
	Message() string
	Cause() error
	Stack() []tsfFrame
	Attributes() []tsfAttr
	LastCaller() tsfFrame
}

// tsfFrame is the interface for supporting logs (github.com/yanun0323/logs)
type tsfFrame interface {
	File() string
	Function() string
	Line() string
}

// tsfAttr is the interface for supporting logs (github.com/yanun0323/logs)
type tsfAttr interface {
	Key() string
	Value() any
}

type tsfFrameImpl struct {
	file     string
	function string
	line     string
}

func (f tsfFrameImpl) File() string {
	return f.file
}

func (f tsfFrameImpl) Function() string {
	return f.function
}

func (f tsfFrameImpl) Line() string {
	return f.line
}

type tsfAttrImpl struct {
	key   string
	value any
}

func (a tsfAttrImpl) Key() string {
	return a.key
}

func (a tsfAttrImpl) Value() any {
	return a.value
}

// make errorStack implements tsf interface

func (e *errorStack) Message() string {
	return e.message
}

func (e *errorStack) Cause() error {
	return e.cause
}

func (e *errorStack) Stack() []tsfFrame {
	frames := make([]tsfFrame, len(e.stack))
	for i, frame := range e.stack {
		frames[i] = tsfFrameImpl{
			file:     frame.File,
			function: frame.Function,
			line:     frame.Line,
		}
	}
	return frames
}

func (e *errorStack) Attributes() []tsfAttr {
	attrs := make([]tsfAttr, len(e.attr))
	for i, attr := range e.attr {
		buf := stringBuilderPool.Get().(*strings.Builder)
		buf.Reset()
		buf.Grow(len(attr.Function) + len(attr.Key) + 1)

		buf.WriteString(attr.Function)
		buf.WriteByte('.')
		buf.WriteString(attr.Key)

		attrs[i] = tsfAttrImpl{
			key:   buf.String(),
			value: attr.Value,
		}

		stringBuilderPool.Put(buf)
	}
	return attrs
}
