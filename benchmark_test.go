package errors

import "testing"

func BenchmarkNew(b *testing.B) {
	for b.Loop() {
		New("test error")
	}
}

func BenchmarkErrorf(b *testing.B) {
	err := New("test error")

	for b.Loop() {
		Errorf("test error: %w", err)
	}
}

func BenchmarkWrap(b *testing.B) {
	err := New("test error")

	for b.Loop() {
		Wrapf(err, "test error")
	}
}

func BenchmarkFormat(b *testing.B) {
	err := New("test error").With("key", "value")

	for b.Loop() {
		Format(err)
	}
}

func BenchmarkFormatJson(b *testing.B) {
	err := New("test error").With("key", "value")

	for b.Loop() {
		FormatJson(err)
	}
}

func BenchmarkFormatColorized(b *testing.B) {
	err := New("test error").With("key", "value")

	for b.Loop() {
		FormatColorized(err)
	}
}
