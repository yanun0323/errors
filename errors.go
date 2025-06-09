package errors

import (
	"encoding/json"
	"fmt"
	"runtime"
	"strconv"
	"strings"

	"github.com/yanun0323/errors/internal/colorize"
)

const (
	_defaultSkip = 3
	_tab         = "    "
)

type attr struct {
	Function string `json:"function"`
	Key      string `json:"key"`
	Value    any    `json:"value"`
}

// errorStack the custom error type
type errorStack struct {
	message    string
	cause      error
	lastCaller frame
	stack      []frame
	attr       []attr
}

/*
	########  ##     ## ########  ##       ####  ######
	##     ## ##     ## ##     ## ##        ##  ##    ##
	##     ## ##     ## ##     ## ##        ##  ##
	########  ##     ## ########  ##        ##  ##
	##        ##     ## ##     ## ##        ##  ##
	##        ##     ## ##     ## ##        ##  ##    ##
	##         #######  ########  ######## ####  ######
*/

// Error implements the error interface
func (e *errorStack) Error() string {
	if e == nil {
		return ""
	}
	return e.message
}

// Unwrap implements the unwrap interface
func (e *errorStack) Unwrap() error {
	if e == nil {
		return nil
	}

	if err, ok := e.cause.(unwrap); ok {
		return err.Unwrap()
	}

	return e.cause
}

// Format implements the fmt.Formatter interface
//
// '%s' - error message
// '%v' - text format
// '%+v' - colorized format
// '%#v' - json format
func (e *errorStack) Format(f fmt.State, c rune) {
	if e == nil {
		return
	}

	switch c {
	case 'v':
		if f.Flag('+') {
			f.Write([]byte(e.formatColorized()))
			return
		}

		if f.Flag('#') {
			f.Write([]byte(e.formatJson()))
			return
		}

		f.Write([]byte(e.formatText()))
		return
	}

	f.Write([]byte(e.Error()))
}

// With adds additional fields, supporting method chaining
func (e *errorStack) With(args ...any) Error {
	if e == nil {
		return nil
	}

	attrs := make([]attr, 0, len(e.attr)+len(args)/2)
	attrs = append(attrs, e.attr...)
	attrs = append(attrs, makeArgs(e.lastCaller.Function, args...)...)

	return &errorStack{
		message:    e.message,
		cause:      e.cause,
		lastCaller: e.lastCaller,
		stack:      e.stack,
		attr:       attrs,
	}
}

// String returns basic string format
func (e *errorStack) String() string {
	return e.Error()
}

/*
	########  ########  #### ##     ##    ###    ######## ########
	##     ## ##     ##  ##  ##     ##   ## ##      ##    ##
	##     ## ##     ##  ##  ##     ##  ##   ##     ##    ##
	########  ########   ##  ##     ## ##     ##    ##    ######
	##        ##   ##    ##   ##   ##  #########    ##    ##
	##        ##    ##   ##    ## ##   ##     ##    ##    ##
	##        ##     ## ####    ###    ##     ##    ##    ########
*/

// formatText returns text formatted error information
func (e *errorStack) formatText() string {
	if e == nil {
		return ""
	}

	buf := stringBuilderPool.Get().(*strings.Builder)
	defer stringBuilderPool.Put(buf)
	buf.Reset()

	buf.Grow(1024)

	buf.WriteByte('\n')
	buf.WriteString("error:\n")
	buf.WriteString(_tab)
	buf.WriteString(e.message)
	buf.WriteByte('\n')

	if e.cause != nil {
		buf.WriteString("cause:\n")
		buf.WriteString(_tab)
		buf.WriteString(e.cause.Error())
		buf.WriteByte('\n')
	}

	if len(e.attr) != 0 {
		var (
			attrMap       = make(map[string][]attr, 32)
			attrFunctions = make([]string, 0, 32)
		)

		for _, a := range e.attr {
			if _, ok := attrMap[a.Function]; !ok {
				attrFunctions = append(attrFunctions, a.Function)
			}
			attrMap[a.Function] = append(attrMap[a.Function], a)
		}

		buf.WriteString("field:\n")

		for _, key := range attrFunctions {
			funcName := key
			if funcName == "" {
				funcName = "unknown"
			}
			buf.WriteString(_tab)
			buf.WriteString(funcName)
			buf.WriteByte(':')
			buf.WriteByte(' ')
			buf.WriteByte('\n')

			for _, a := range attrMap[key] {
				buf.WriteString(_tab)
				buf.WriteString(_tab)
				buf.WriteString(a.Key)
				buf.WriteByte(':')
				buf.WriteByte(' ')
				buf.WriteString(fmt.Sprintf("%+v", a.Value))
				buf.WriteByte('\n')
			}
		}
	}

	if len(e.stack) != 0 {
		buf.WriteString("stack:\n")
		for _, f := range e.stack {
			buf.WriteString(_tab)
			buf.WriteString(f.Function)
			buf.WriteByte(':')
			buf.WriteByte('\n')
			buf.WriteString(_tab)
			buf.WriteString(_tab)
			buf.WriteString(f.FormatText())
			buf.WriteByte('\n')
		}
	}

	return buf.String()
}

// formatJson returns formatJson formatted error information
func (e *errorStack) formatJson() string {
	if e == nil {
		return _emptyJSONString
	}

	data := dataPool.Get().(map[string]any)
	defer dataPool.Put(data)

	data["error"] = e.message
	data["field"] = e.attr

	if e.cause != nil {
		data["cause"] = e.cause.Error()
	}

	if len(e.stack) > 0 {
		data["stack"] = e.stack
	}

	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Sprintf(`{"error": "marshal error: %s"}`, err.Error())
	}

	return string(jsonBytes)
}

// formatColorized returns colorized readable format (ANSI color codes)
func (e *errorStack) formatColorized() string {
	if e == nil {
		return ""
	}

	buf := stringBuilderPool.Get().(*strings.Builder)
	defer stringBuilderPool.Put(buf)
	buf.Reset()
	buf.Grow(1024)

	buf.WriteByte('\n')
	colorize.WriteString(buf, colorize.Red, "[error] ")
	buf.WriteString(e.message)
	buf.WriteByte('\n')

	if e.cause != nil {
		colorize.WriteString(buf, colorize.Yellow, "[cause] ")
		buf.WriteString(e.cause.Error())
		buf.WriteByte('\n')
	}

	if len(e.attr) > 0 {
		var (
			attrMap       = make(map[string][]attr, 32)
			attrFunctions = make([]string, 0, 32)
		)
		for _, a := range e.attr {
			if _, ok := attrMap[a.Function]; !ok {
				attrFunctions = append(attrFunctions, a.Function)
			}
			attrMap[a.Function] = append(attrMap[a.Function], a)
		}

		colorize.WriteString(buf, colorize.Cyan, "[field]")
		buf.WriteByte('\n')
		for _, key := range attrFunctions {
			hasFuncName := key != ""
			if hasFuncName {
				buf.WriteString(_tab)
				colorize.WriteString(buf, colorize.Blue, "[", key, "] ")
				buf.WriteByte('\n')
			}

			for _, a := range attrMap[key] {
				if hasFuncName {
					buf.WriteString(_tab)
				}
				buf.WriteString(_tab)
				colorize.WriteString(buf, colorize.Magenta, "[", a.Key, "] ")
				colorize.WriteString(buf, colorize.Black, fmt.Sprintf("%+v\n", a.Value))
			}
		}
	}

	if len(e.stack) > 0 {
		colorize.WriteString(buf, colorize.Cyan, "[stack]")
		buf.WriteByte('\n')
		for _, f := range e.stack {
			if strings.HasPrefix(f.Function, "runtime") {
				continue
			}

			buf.WriteString(_tab)
			buf.WriteString(f.FormatColorized(colorize.Blue, colorize.Black))
			buf.WriteByte('\n')
		}
	}

	return buf.String()
}

// getStack captures the current call stack
func getStack(additionalSkip ...int) []frame {
	skip := _defaultSkip
	if len(additionalSkip) != 0 {
		skip += additionalSkip[0]
	}

	var frames []frame

	pc := make([]uintptr, 32)
	n := runtime.Callers(skip, pc)

	if n == 0 {
		return frames
	}

	callersFrames := runtime.CallersFrames(pc)

	for {
		f, more := callersFrames.Next()
		if SkipRuntimeStackTrace && canSkip(f) {
			continue
		}

		if f.Function != "" {
			funcName := f.Function
			span := strings.Split(funcName, "/")
			funcName = span[len(span)-1]
			span = strings.Split(funcName, ".")
			funcName = span[len(span)-1]

			frames = append(frames, frame{
				File:     f.File,
				Function: funcName,
				Line:     strconv.Itoa(f.Line),
			})
		}

		if !more {
			break
		}
	}

	return frames
}

func canSkip(f runtime.Frame) bool {
	if strings.HasSuffix(f.Function, ".init") {
		return true
	}

	if strings.HasSuffix(f.Function, "tRunner") {
		return true
	}

	if strings.Contains(f.File, "/src/runtime/") {
		return true
	}

	return false
}
