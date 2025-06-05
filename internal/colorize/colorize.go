package colorize

import (
	"bytes"
	"io"
	"sync"
)

const (
	Reset = "\x1b[0m"
)

const (
	Black   = "\x1b[30m"
	Red     = "\x1b[31m"
	Green   = "\x1b[32m"
	Yellow  = "\x1b[33m"
	Blue    = "\x1b[34m"
	Magenta = "\x1b[35m"
	Cyan    = "\x1b[36m"
	White   = "\x1b[37m"

	BlackReversed   = "\x1b[40m"
	RedReversed     = "\x1b[41m"
	GreenReversed   = "\x1b[42m"
	YellowReversed  = "\x1b[43m"
	BlueReversed    = "\x1b[44m"
	MagentaReversed = "\x1b[45m"
	CyanReversed    = "\x1b[46m"
	WhiteReversed   = "\x1b[47m"

	BrightBlack   = "\x1b[90m"
	BrightRed     = "\x1b[91m"
	BrightGreen   = "\x1b[92m"
	BrightYellow  = "\x1b[93m"
	BrightBlue    = "\x1b[94m"
	BrightMagenta = "\x1b[95m"
	BrightCyan    = "\x1b[96m"
	BrightWhite   = "\x1b[97m"

	BrightBlackReversed   = "\x1b[100m"
	BrightRedReversed     = "\x1b[101m"
	BrightGreenReversed   = "\x1b[102m"
	BrightYellowReversed  = "\x1b[103m"
	BrightBlueReversed    = "\x1b[104m"
	BrightMagentaReversed = "\x1b[105m"
	BrightCyanReversed    = "\x1b[106m"
	BrightWhiteReversed   = "\x1b[107m"
)

var (
	bufferPool = sync.Pool{
		New: func() any {
			return bytes.NewBuffer(make([]byte, 0, 1024))
		},
	}
)

// String colorize string
func String(c string, contents ...string) string {
	buf := bufferPool.Get().(*bytes.Buffer)
	defer bufferPool.Put(buf)
	buf.Reset()
	WriteString(buf, c, contents...)
	return buf.String()
}

// Bytes colorize bytes
func Bytes(color string, contents ...[]byte) []byte {
	buf := bufferPool.Get().(*bytes.Buffer)
	defer bufferPool.Put(buf)
	buf.Reset()
	WriteBytes(buf, color, contents...)
	return buf.Bytes()
}

// Writer is the interface that wraps the basic WriteString method.
type Writer interface {
	io.Writer
	io.StringWriter
	io.ByteWriter
}

// WriteString write colorized string to buffer
func WriteString(buf Writer, color string, contents ...string) {
	buf.WriteString(color)
	for _, s := range contents {
		buf.WriteString(s)
	}
	buf.WriteString(Reset)
}

// WriteBytes write colorized bytes to buffer
func WriteBytes(buf Writer, color string, contents ...[]byte) {
	buf.WriteString(color)
	for _, b := range contents {
		buf.Write(b)
	}
	buf.WriteString(Reset)
}

// ResetString reset string color
func ResetString(s string) string {
	if len(s) == 0 {
		return s
	}

	buf := bufferPool.Get().(*bytes.Buffer)
	defer bufferPool.Put(buf)
	buf.Reset()

	i := 0
	for i < len(s) {
		if s[i] == '\x1b' && i+1 < len(s) && s[i+1] == '[' {
			// Find the end of ANSI escape sequence
			j := i + 2
			for j < len(s) && s[j] != 'm' {
				j++
			}
			if j < len(s) {
				// Skip the entire escape sequence including 'm'
				i = j + 1
			} else {
				// Malformed escape sequence, keep the character
				buf.WriteByte(s[i])
				i++
			}
		} else {
			buf.WriteByte(s[i])
			i++
		}
	}

	return buf.String()
}

// ResetBytes reset bytes color
func ResetBytes(s []byte) []byte {
	if len(s) == 0 {
		return s
	}

	buf := bufferPool.Get().(*bytes.Buffer)
	defer bufferPool.Put(buf)
	buf.Reset()

	i := 0
	for i < len(s) {
		if s[i] == '\x1b' && i+1 < len(s) && s[i+1] == '[' {
			// Find the end of ANSI escape sequence
			j := i + 2
			for j < len(s) && s[j] != 'm' {
				j++
			}
			if j < len(s) {
				// Skip the entire escape sequence including 'm'
				i = j + 1
			} else {
				// Malformed escape sequence, keep the character
				buf.WriteByte(s[i])
				i++
			}
		} else {
			buf.WriteByte(s[i])
			i++
		}
	}

	return buf.Bytes()
}
