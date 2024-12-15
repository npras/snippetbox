package assert

import (
	"strings"
	"testing"
)

func Equal[T comparable](t *testing.T, got, want T) {
	t.Helper()
	if got != want {
		t.Errorf("got: %v, want: %v", got, want)
	}
}

func StringContains(t *testing.T, str, substr string) {
	t.Helper()
	if !strings.Contains(str, substr) {
		t.Errorf("got: %q, expected to contain: %q", str, substr)
	}
}
