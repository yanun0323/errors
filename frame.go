package errors

import (
	"strings"

	"github.com/yanun0323/errors/internal/colorize"
)

// frame represents a single frame in the stack trace
type frame struct {
	File     string `json:"file"`
	Function string `json:"function"`
	Line     string `json:"line"`
}

func (f frame) FormatText() string {
	buf := stringBuilderPool.Get().(*strings.Builder)
	defer stringBuilderPool.Put(buf)
	buf.Reset()

	buf.Grow(len(f.File) + len(f.Line) + len(f.Function) + 4)

	buf.WriteString(f.File)
	buf.WriteByte(':')
	buf.WriteString(f.Line)
	buf.WriteString(" in ")
	buf.WriteString(f.Function)

	return buf.String()
}

func (f frame) FormatColorized(funcColor, fileColor string) string {
	buf := stringBuilderPool.Get().(*strings.Builder)
	defer stringBuilderPool.Put(buf)
	buf.Reset()

	buf.Grow(len(f.File) + len(f.Line) + len(f.Function) + 4)

	colorize.WriteString(buf, funcColor, "[", f.Function, "] ")
	colorize.WriteString(buf, fileColor, f.File+":"+f.Line)

	return buf.String()
}
