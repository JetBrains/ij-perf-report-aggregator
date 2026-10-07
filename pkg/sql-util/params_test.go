package sql_util

import (
	"strings"
	"testing"
)

func TestValidateIdentifier(t *testing.T) {
	t.Parallel()
	good := []string{"foo", "FOO", "foo_bar", "f00", "_underscore", "perfintDev", "a"}
	for _, s := range good {
		if err := ValidateIdentifier("table", s); err != nil {
			t.Errorf("ValidateIdentifier(%q) unexpected error: %v", s, err)
		}
	}

	bad := []struct {
		v       string
		wantSub string
	}{
		{"", "is required"},
		{"foo bar", "invalid character"},
		{"foo;drop", "invalid character"},
		{"foo-bar", "invalid character"},
		{"foo.bar", "invalid character"},
		{"foo`bar", "invalid character"},
		{"foo'bar", "invalid character"},
		{"foo\nbar", "invalid character"},
	}
	for _, tc := range bad {
		err := ValidateIdentifier("table", tc.v)
		if err == nil {
			t.Errorf("ValidateIdentifier(%q) expected error", tc.v)
			continue
		}
		if !strings.Contains(err.Error(), tc.wantSub) {
			t.Errorf("ValidateIdentifier(%q) = %v, want substring %q", tc.v, err, tc.wantSub)
		}
	}
}
