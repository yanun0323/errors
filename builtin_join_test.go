package errors

import (
	"testing"
)

func TestJoin(t *testing.T) {
	err1 := New("test error")
	err2 := New("test error 2")
	errs := Join(err1, err2)

	if errs.Error() != "test error\ntest error 2" {
		t.Errorf("expected 'test error\ntest error 2', got '%s'", errs.Error())
	}

	if !Is(errs, err1) {
		t.Errorf("expected 'test error', got '%s'", errs.Error())
	}

	if !Is(errs, err2) {
		t.Errorf("expected 'test error 2', got '%s'", errs.Error())
	}

	if !Is(errs, errs) {
		t.Errorf("expected 'test error\ntest error 2', got '%s'", errs.Error())
	}
}
