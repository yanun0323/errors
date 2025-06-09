package errors

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/yanun0323/errors/internal/colorize"
)

var errCause = New("root")

func TestNew(t *testing.T) {
	err := New("test error").With("k1", "v1").With("k2", 2).(*errorStack)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	expected := `
[error] test error
[cause] test error
[field]
    [TestNew] 
        [k1] v1
        [k2] 2
[stack]
    [TestNew] /Users/Shared/Project/personal/go/errors/errors_test.go:16
`

	f := FormatColorized(err)
	f = colorize.ResetString(f)
	if !strings.EqualFold(f, expected) {
		t.Errorf("Expected colorized output '%s', but got '%s'", expected, f)
	}
}

func TestErrorfBasic(t *testing.T) {
	err := Errorf("formatted error %d", 123).With(
		"k1", "v1",
	).With(
		"k2", 2,
		"k3", 3,
	).(*errorStack)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	expected := `
[error] formatted error 123
[cause] formatted error 123
[field]
    [TestErrorfBasic] 
        [k1] v1
        [k2] 2
        [k3] 3
[stack]
    [TestErrorfBasic] /Users/Shared/Project/personal/go/errors/errors_test.go:40
`

	f := FormatColorized(err)
	f = colorize.ResetString(f)
	if !strings.EqualFold(f, expected) {
		t.Errorf("Expected colorized output '%s', but got '%s'", expected, f)
	}
}

func causeError() error {
	return Wrap(errCause).With("r1", "rv1").With("r2", 22)
}

func TestErrorfWrap(t *testing.T) {
	err := Errorf("formatted error, err: %w", causeError()).
		With("user_id", 123).
		With(
			"k1", "v1",
			"k2", 2,
		).(*errorStack)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	expected := `
[error] formatted error, err: root
[cause] root
[field]
    [causeError] 
        [r1] rv1
        [r2] 22
    [TestErrorfWrap] 
        [user_id] 123
        [k1] v1
        [k2] 2
[stack]
    [causeError] /Users/Shared/Project/personal/go/errors/errors_test.go:70
    [TestErrorfWrap] /Users/Shared/Project/personal/go/errors/errors_test.go:74
`

	f := FormatColorized(err)
	f = colorize.ResetString(f)
	if !strings.EqualFold(f, expected) {
		t.Errorf("Expected colorized output '%s', but got '%s'", expected, f)
	}
}

func TestWrap(t *testing.T) {
	wrappedErr := Wrap(causeError(), "wrapped").With(
		"k1", "v1",
		"k2", 2,
	).With(
		"k3", 3,
	)

	if unwrapped := Unwrap(wrappedErr); !Is(unwrapped, errCause) {
		t.Error("Unwrap failed")
		t.Errorf("unwrapped: %+v", Unwrap(unwrapped))
		t.Errorf("baseErr: %+v", Unwrap(errCause))
	}

	expected := `
[error] wrapped, err: root
[cause] root
[field]
    [causeError] 
        [r1] rv1
        [r2] 22
    [TestWrap] 
        [k1] v1
        [k2] 2
        [k3] 3
[stack]
    [causeError] /Users/Shared/Project/personal/go/errors/errors_test.go:70
    [TestWrap] /Users/Shared/Project/personal/go/errors/errors_test.go:108
`

	f := FormatColorized(wrappedErr)
	f = colorize.ResetString(f)
	if !strings.EqualFold(f, expected) {
		t.Errorf("Expected colorized output '%s', but got '%s'", expected, f)
	}
}

func TestWrapf(t *testing.T) {
	wrappedErr := Wrapf(causeError(), "wrapped %s", "world").With(
		"k1", "v1",
		"k2", 2,
	).With(
		"k3", 3,
	)

	if unwrapped := Unwrap(wrappedErr); !Is(unwrapped, errCause) {
		t.Error("Unwrap failed")
		t.Errorf("unwrapped: %+v", Unwrap(unwrapped))
		t.Errorf("baseErr: %+v", Unwrap(errCause))
	}

	expected := `
[error] wrapped world, err: root
[cause] root
[field]
    [causeError] 
        [r1] rv1
        [r2] 22
    [TestWrapf] 
        [k1] v1
        [k2] 2
        [k3] 3
[stack]
    [causeError] /Users/Shared/Project/personal/go/errors/errors_test.go:70
    [TestWrapf] /Users/Shared/Project/personal/go/errors/errors_test.go:145
`

	f := FormatColorized(wrappedErr)
	f = colorize.ResetString(f)
	if !strings.EqualFold(f, expected) {
		t.Errorf("Expected colorized output '%s', but got '%s'", expected, f)
	}
}

func TestWith(t *testing.T) {
	err := New("test error").With("user_id", 123).With("action", "create")

	attrs := err.(*errorStack).attr
	if len(attrs) != 2 {
		t.Fatalf("Expected 2 attributes, got %d", len(attrs))
	}

	if attrs[0] != (attr{"TestWith", "user_id", 123}) {
		t.Errorf("Expected user_id=123, got %v", attrs[0])
	}

	if attrs[1] != (attr{"TestWith", "action", "create"}) {
		t.Errorf("Expected action=create, got %v", attrs[1])
	}
}

func TestIs(t *testing.T) {
	baseErr := New("base error")
	wrappedErr := &errorStack{
		message: "wrapped error",
		cause:   baseErr,
		stack:   getStack(0),
		attr:    []attr{},
	}

	if !Is(wrappedErr, baseErr) {
		t.Error("Is should return true for wrapped error")
	}

	otherErr := New("other error")
	if Is(wrappedErr, otherErr) {
		t.Error("Is should return false for different error")
	}
}

func TestAs(t *testing.T) {
	customErr := &errorStack{
		message: "custom error",
		stack:   getStack(0),
		attr:    []attr{},
	}
	wrappedErr := &errorStack{
		message: "wrapped error",
		cause:   customErr,
		stack:   getStack(0),
		attr:    []attr{},
	}

	var target *errorStack
	if !As(wrappedErr, &target) {
		t.Error("As should return true for correct type")
	}

	if target != wrappedErr {
		t.Errorf("As should set target to wrappedErr, got %v, expected %v", target, wrappedErr)
	}

	var otherTarget *OtherError
	if As(wrappedErr, &otherTarget) {
		t.Error("As should return false for non-matching type")
	}
}

type OtherError struct {
	msg string
}

func (e *OtherError) Error() string {
	return e.msg
}

func TestString(t *testing.T) {
	err := New("test error").With("key", "value")
	str := err.(*errorStack).String()

	if str != "test error" {
		t.Errorf("String should contain error message, got: %s", str)
	}
}

func TestJSON(t *testing.T) {
	err := New("test error").With("key", "value")
	jsonStr := FormatJson(err)

	if !containsString(jsonStr, "test error") {
		t.Error("JSON should contain error message")
	}

	if !containsString(jsonStr, "key") {
		t.Error("JSON should contain fields")
	}
}

func TestColorizedString(t *testing.T) {
	err := New("test error").With("key", "value")
	colorStr := FormatColorized(err)

	if !containsString(colorStr, "test error") {
		t.Error("ColorizedString should contain error message")
	}

	if !containsString(colorStr, "key") {
		t.Error("ColorizedString should contain fields")
	}
}

func TestStackTrace(t *testing.T) {
	err := New("test error")
	stack := err.(*errorStack).stack

	if len(stack) == 0 {
		t.Error("Expected stack trace")
	}

	if stack[0].Function == "" {
		t.Error("Expected function name in stack trace")
	}
}

// containsString checks if a string contains a substring
func containsString(s, substr string) bool {
	return len(s) >= len(substr) && s != "" && findInString(s, substr)
}

func findInString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// Example demonstrates how to use the errors package
func Example() {
	err1 := New("something went wrong")
	fmt.Println(err1.Error())

	err2 := &errorStack{
		message: "failed to process user: something went wrong",
		cause:   err1,
		stack:   getStack(0),
		attr:    []attr{},
	}
	fmt.Println(err2.Error())

	err3 := New("database error").
		With("table", "users").
		With("operation", "insert")

	fmt.Println(err3.(*errorStack).String())

	fmt.Println("JSON:", FormatJson(err3))
	fmt.Println("Colorized:", FormatColorized(err3))
}

func TestBasicCreation(t *testing.T) {
	err := New("database connection")
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if err.Error() != "database connection" {
		t.Errorf("Expected 'database connection', got '%s'", err.Error())
	}
}

func TestErrorWrapping(t *testing.T) {
	err1 := New("database connection")
	err2 := Errorf("failed to create user: %w", err1)
	if err2 == nil {
		t.Fatal("Expected error, got nil")
	}

	if err2.Error() != "failed to create user: database connection" {
		t.Errorf("Expected 'failed to create user: database connection', got '%s'", err2.Error())
	}
}

func TestAdditionalFields(t *testing.T) {
	err := New("user validation failed").
		With("user_id", 12345).
		With("email", "user@example.com").
		With("attempt", 3).(*errorStack)

	if err.Error() != "user validation failed" {
		t.Errorf("Expected 'user validation failed', got '%s'", err.Error())
	}

	if len(err.attr) != 3 {
		t.Fatalf("Expected 3 attributes, got %d", len(err.attr))
	}

	if err.attr[0] != (attr{"TestAdditionalFields", "user_id", 12345}) {
		t.Errorf("Expected user_id=12345, got %v", err.attr[0])
	}

	if err.attr[1] != (attr{"TestAdditionalFields", "email", "user@example.com"}) {
		t.Errorf("Expected email=user@example.com, got %v", err.attr[1])
	}

	if err.attr[2] != (attr{"TestAdditionalFields", "attempt", 3}) {
		t.Errorf("Expected attempt=3, got %v", err.attr[2])
	}
}

func TestChainChecking(t *testing.T) {
	err := New("API call failed").
		With("user_id", 12345).
		With("email", "user@example.com").
		With("attempt", 3).(*errorStack)

	if err.Error() != "API call failed" {
		t.Errorf("Expected 'API call failed', got '%s'", err.Error())
	}

	if len(err.attr) != 3 {
		t.Fatalf("Expected 3 attributes, got %d", len(err.attr))
	}

	if err.attr[0] != (attr{"TestChainChecking", "user_id", 12345}) {
		t.Errorf("Expected user_id=12345, got %v", err.attr[0])
	}

	if err.attr[1] != (attr{"TestChainChecking", "email", "user@example.com"}) {
		t.Errorf("Expected email=user@example.com, got %v", err.attr[1])
	}

	if err.attr[2] != (attr{"TestChainChecking", "attempt", 3}) {
		t.Errorf("Expected attempt=3, got %v", err.attr[2])
	}
}

func TestFormat(t *testing.T) {
	err := New("test error").With("key", "value")
	formatted := Format(err)

	expected := `
error:
    test error
cause:
    test error
field:
    TestFormat: 
        key: value
stack:
    TestFormat:
        /Users/Shared/Project/personal/go/errors/errors_test.go:416 in TestFormat
`

	if formatted != expected {
		t.Errorf("Expected '%s', got '%s'", expected, formatted)
	}

	os.WriteFile("./test/format.txt", []byte(formatted), 0644)
}

func TestFormatJson(t *testing.T) {
	err := New("user validation failed").
		With("user_id", 12345).
		With("email", "user@example.com").
		With("attempt", 3)

	expected := `{
  "cause": "user validation failed",
  "error": "user validation failed",
  "field": [
    {
      "function": "TestFormatJson",
      "key": "user_id",
      "value": 12345
    },
    {
      "function": "TestFormatJson",
      "key": "email",
      "value": "user@example.com"
    },
    {
      "function": "TestFormatJson",
      "key": "attempt",
      "value": 3
    }
  ],
  "stack": [
    {
      "file": "/Users/Shared/Project/personal/go/errors/errors_test.go",
      "function": "TestFormatJson",
      "line": "440"
    }
  ]
}`

	if FormatJson(err) != expected {
		t.Errorf("Expected JSON output, got '%s'", FormatJson(err))
	}

	os.WriteFile("./test/json.json", []byte(FormatJson(err)), 0644)
}

func TestFormatColorized(t *testing.T) {
	err := New("user validation failed").
		With("user_id", 12345).
		With("email", "user@example.com").
		With("attempt", 3)

	expected := `
[error] user validation failed
[cause] user validation failed
[field]
    [TestFormatColorized] 
        [user_id] 12345
        [email] user@example.com
        [attempt] 3
[stack]
    [TestFormatColorized] /Users/Shared/Project/personal/go/errors/errors_test.go:482
`

	f := FormatColorized(err)
	f = colorize.ResetString(f)
	if !strings.EqualFold(f, expected) {
		t.Errorf("Expected colorized output, got '%s'", f)
	}

	os.WriteFile("./test/colorized.txt", []byte(f), 0644)
}

var errConnectionTimeout = New("database connection timeout")

func TestRealWorld(t *testing.T) {
	err := processUser(123)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	println(FormatColorized(err))

	if !Is(err, errConnectionTimeout) {
		t.Fatalf("Expected error to be a database connection timeout, got: %+v", err)
	}

	expected := "user validation, err: check time, user(123), err: database connection timeout"
	if err.Error() != expected {
		t.Fatalf("Expected '%s', got '%s'", expected, err.Error())
	}

	sErr := err.(*errorStack)

	if len(sErr.attr) != 3 {
		t.Fatalf("Expected 3 attributes, got %d", len(sErr.attr))
	}

	if sErr.attr[0] != (attr{"checkTime", "now", "2025-06-04 16:47:09"}) {
		t.Errorf("Expected now=2025-06-04 16:47:09, got %v", sErr.attr[0])
	}

	if sErr.attr[1] != (attr{"validateUser", "user_id", 123}) {
		t.Errorf("Expected user_id=123, got %v", sErr.attr[1])
	}

	if sErr.attr[2] != (attr{"validateUser", "table", "users"}) {
		t.Errorf("Expected table=users, got %v", sErr.attr[2])
	}
}

// processUser handles user-related operations
func processUser(userID int) error {
	if err := validateUser(userID); err != nil {
		return Errorf("user validation, err: %w", err)
	}

	return nil
}

// validateUser validates if user exists
func validateUser(userID int) error {
	if userID == 0 {
		return New("user not found").
			With("user_id", userID).
			With("table", "users")
	}

	if err := checkTime(); err != nil {
		return Wrapf(err, "check time, user(%d)", userID).
			With("user_id", userID).
			With("table", "users")
	}

	return nil
}

func checkTime() error {
	now := time.Now()
	if now.Unix() > 0 {
		err := Wrap(errConnectionTimeout).
			With("now", "2025-06-04 16:47:09")
		return err
	}

	return nil
}
