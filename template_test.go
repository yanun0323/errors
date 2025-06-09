package errors

import (
	"strings"
	"testing"

	"github.com/yanun0323/errors/internal/colorize"
)

func TestTemplateNew(t *testing.T) {
	tpl := NewTemplate(
		"k1", "v1",
		"k2", 2,
	).With(
		"k3", "v3",
		"k4", 4,
	)

	err := tpl.New("test").With(
		"k5", "v5",
		"k6", 6,
	)

	expected := `
[error] test
[cause] test
[field]
    [TestTemplateNew] 
        [k1] v1
        [k2] 2
        [k3] v3
        [k4] 4
        [k5] v5
        [k6] 6
[stack]
    [TestTemplateNew] /Users/Shared/Project/personal/go/errors/template_test.go:19
`

	f := FormatColorized(err)
	f = colorize.ResetString(f)
	if !strings.EqualFold(f, expected) {
		t.Errorf("Expected colorized output '%s', but got '%s'", expected, f)
	}
}

func TestTemplateWrap(t *testing.T) {
	tpl := NewTemplate(
		"k1", "v1",
		"k2", 2,
	).With(
		"k3", "v3",
		"k4", 4,
	)

	err := tpl.Wrap(causeError(), "hello").With(
		"k5", "v5",
		"k6", 6,
	)

	expected := `
[error] hello, err: root
[cause] root
[field]
    [causeError] 
        [r1] rv1
        [r2] 22
    [TestTemplateWrap] 
        [k1] v1
        [k2] 2
        [k3] v3
        [k4] 4
        [k5] v5
        [k6] 6
[stack]
    [causeError] /Users/Shared/Project/personal/go/errors/errors_test.go:70
    [TestTemplateWrap] /Users/Shared/Project/personal/go/errors/template_test.go:55
`

	f := FormatColorized(err)
	f = colorize.ResetString(f)
	if !strings.EqualFold(f, expected) {
		t.Errorf("Expected colorized output '%s', but got '%s'", expected, f)
	}
}

func TestTemplateWrapf(t *testing.T) {
	tpl := NewTemplate(
		"k1", "v1",
		"k2", 2,
	).With(
		"k3", "v3",
		"k4", 4,
	)

	err := tpl.Wrapf(causeError(), "hello %s", "world").With(
		"k5", "v5",
		"k6", 6,
	)

	expected := `
[error] hello world, err: root
[cause] root
[field]
    [causeError] 
        [r1] rv1
        [r2] 22
    [TestTemplateWrapf] 
        [k1] v1
        [k2] 2
        [k3] v3
        [k4] 4
        [k5] v5
        [k6] 6
[stack]
    [causeError] /Users/Shared/Project/personal/go/errors/errors_test.go:70
    [TestTemplateWrapf] /Users/Shared/Project/personal/go/errors/template_test.go:95
`

	f := FormatColorized(err)
	f = colorize.ResetString(f)
	if !strings.EqualFold(f, expected) {
		t.Errorf("Expected colorized output '%s', but got '%s'", expected, f)
	}
}

func TestTemplateErrorf(t *testing.T) {
	tpl := NewTemplate(
		"k1", "v1",
		"k2", 2,
	).With(
		"k3", "v3",
		"k4", 4,
	)

	err := tpl.Errorf("hello %s", "error").With(
		"k5", "v5",
		"k6", 6,
	)

	expected := `
[error] hello error
[cause] hello error
[field]
    [TestTemplateErrorf] 
        [k1] v1
        [k2] 2
        [k3] v3
        [k4] 4
        [k5] v5
        [k6] 6
[stack]
    [TestTemplateErrorf] /Users/Shared/Project/personal/go/errors/template_test.go:135
`

	f := FormatColorized(err)
	f = colorize.ResetString(f)
	if !strings.EqualFold(f, expected) {
		t.Errorf("Expected colorized output '%s', but got '%s'", expected, f)
	}
}

func TestTemplateErrorWrap(t *testing.T) {
	tpl := NewTemplate(
		"k1", "v1",
		"k2", 2,
	).With(
		"k3", "v3",
		"k4", 4,
	)

	err := tpl.Errorf("hello %s, err: %w", "world", causeError()).With(
		"k5", "v5",
		"k6", 6,
	)

	expected := `
[error] hello world, err: root
[cause] root
[field]
    [causeError] 
        [r1] rv1
        [r2] 22
    [TestTemplateErrorWrap] 
        [k1] v1
        [k2] 2
        [k3] v3
        [k4] 4
        [k5] v5
        [k6] 6
[stack]
    [causeError] /Users/Shared/Project/personal/go/errors/errors_test.go:70
    [TestTemplateErrorWrap] /Users/Shared/Project/personal/go/errors/template_test.go:171
`

	f := FormatColorized(err)
	f = colorize.ResetString(f)
	if !strings.EqualFold(f, expected) {
		t.Errorf("Expected colorized output '%s', but got '%s'", expected, f)
	}
}

func TestTemplateNil(t *testing.T) {
	temp := NewTemplate(
		"k1", "v1",
		"k2", 2,
	)

	err := nilError()

	if temp.Wrap(err, "hello") != nil {
		t.Fatal("nilError should be nil")
	}

	if temp.Wrapf(err, "hello %s", "world") != nil {
		t.Fatal("nilError should be nil")
	}

	if temp.Errorf("hello %w", err) != nil {
		t.Fatal("nilError should be nil")
	}

	{
		var err error
		if temp.Errorf("hello %w", err) != nil {
			t.Fatal("nilError should be nil")
		}
	}

	{
		var err any
		if temp.Errorf("hello %w", err) != nil {
			t.Fatal("nilError should be nil")
		}
	}

	{
		var err any = 123
		if temp.Errorf("hello %w", err) != nil {
			t.Fatal("nilError should be nil")
		}
	}
}
