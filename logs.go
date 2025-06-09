package errors

import (
	"strings"

	"github.com/yanun0323/errors/internal/logs"
)

// make errorStack implements tsf interface
var (
	_ logs.Error = (*errorStack)(nil)
	_ logs.Frame = (*frame)(nil)
	_ logs.Attr  = (*attr)(nil)
)

func (e *errorStack) Message() string {
	return e.message
}

func (e *errorStack) Cause() error {
	return e.cause
}

func (e *errorStack) Stack() []any {
	frames := make([]any, 0, len(e.stack))
	for _, frame := range e.stack {
		frames = append(frames, frame)
	}

	return frames
}

func (e *errorStack) Attributes() []any {
	attrs := make([]any, 0, len(e.attr))
	for _, attr := range e.attr {
		attrs = append(attrs, attr)
	}

	return attrs
}

func (f frame) Parameters() (file, function, line string) {
	return f.File, f.Function, f.Line
}

func (a attr) Parameters() (key string, value any) {
	buf := stringBuilderPool.Get().(*strings.Builder)
	defer stringBuilderPool.Put(buf)
	buf.Reset()
	buf.Grow(len(a.Function) + len(a.Key) + 1)

	buf.WriteString(a.Function)
	buf.WriteByte('.')
	buf.WriteString(a.Key)

	return buf.String(), a.Value
}
