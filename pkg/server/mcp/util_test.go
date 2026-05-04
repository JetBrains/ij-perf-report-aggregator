package mcp

import (
	"encoding/base64"
	"encoding/hex"
	"strings"
	"testing"
)

func TestValidateIdentifier(t *testing.T) {
	t.Parallel()
	good := []string{"foo", "FOO", "foo_bar", "f00", "_underscore", "perfintDev", "a"}
	for _, s := range good {
		if err := validateIdentifier("table", s); err != nil {
			t.Errorf("validateIdentifier(%q) unexpected error: %v", s, err)
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
		err := validateIdentifier("table", tc.v)
		if err == nil {
			t.Errorf("validateIdentifier(%q) expected error", tc.v)
			continue
		}
		if !strings.Contains(err.Error(), tc.wantSub) {
			t.Errorf("validateIdentifier(%q) = %v, want substring %q", tc.v, err, tc.wantSub)
		}
	}
}

func TestBuildUnion(t *testing.T) {
	t.Parallel()
	tables := []tableRef{
		{Database: "db1", Table: "t1"},
		{Database: "db2", Table: "t2"},
	}
	perTable := func(r tableRef) (string, []any) {
		return "select 1 from " + r.Database + "." + r.Table + " where x = ?", []any{r.Database}
	}

	sql, args := buildUnion(tables, perTable)

	wantSQL := "(select 1 from db1.t1 where x = ?) union all (select 1 from db2.t2 where x = ?)"
	if sql != wantSQL {
		t.Errorf("sql mismatch:\n got: %s\nwant: %s", sql, wantSQL)
	}
	if len(args) != 2 || args[0] != "db1" || args[1] != "db2" {
		t.Errorf("args = %v, want [db1 db2]", args)
	}
}

func TestBuildUnion_SingleTable(t *testing.T) {
	t.Parallel()
	sql, args := buildUnion(
		[]tableRef{{Database: "d", Table: "t"}},
		func(r tableRef) (string, []any) { return "select " + r.Database, []any{1} },
	)
	if sql != "(select d)" {
		t.Errorf("single-table sql = %q", sql)
	}
	if len(args) != 1 || args[0] != 1 {
		t.Errorf("args = %v", args)
	}
}

func TestBuildUnion_Empty(t *testing.T) {
	t.Parallel()
	sql, args := buildUnion(nil, func(tableRef) (string, []any) { return "x", nil })
	if sql != "" {
		t.Errorf("empty buildUnion sql = %q, want empty", sql)
	}
	if len(args) != 0 {
		t.Errorf("empty buildUnion args = %v, want empty", args)
	}
}

// encodeSHA1 returns a fixed-shape base64 encoding of a 20-byte SHA-1, matching how
// the collector stores `installer.changes` entries (raw std b64, no padding).
func encodeSHA1(t *testing.T, hexSHA string) string {
	t.Helper()
	b, err := hex.DecodeString(hexSHA)
	if err != nil {
		t.Fatalf("encodeSHA1 fixture invalid: %v", err)
	}
	return base64.RawStdEncoding.EncodeToString(b)
}

func TestCommitRange(t *testing.T) {
	t.Parallel()
	// Collector stores commits newest-first. commitRange returns (oldest, newest).
	const newestSHA = "22596363b3de40b06f981fb85d82312e8c0ed511"
	const oldestSHA = "0123456789abcdef0123456789abcdef01234567"
	newestEnc := encodeSHA1(t, newestSHA)
	oldestEnc := encodeSHA1(t, oldestSHA)

	if first, last := commitRange(nil); first != "" || last != "" {
		t.Errorf("commitRange(nil) = %q,%q, want empty", first, last)
	}

	first, last := commitRange([]string{newestEnc, oldestEnc})
	if last != newestSHA[:shortCommitLen] {
		t.Errorf("last = %q, want %q (newest, head of slice)", last, newestSHA[:shortCommitLen])
	}
	if first != oldestSHA[:shortCommitLen] {
		t.Errorf("first = %q, want %q (oldest, tail of slice)", first, oldestSHA[:shortCommitLen])
	}

	// Single commit: first == last.
	f1, l1 := commitRange([]string{newestEnc})
	if f1 != l1 || f1 != newestSHA[:shortCommitLen] {
		t.Errorf("single-commit range = %q,%q, want both %q", f1, l1, newestSHA[:shortCommitLen])
	}

	// Non-base64 entries fall through unchanged so a future format switch is visible.
	if f, l := commitRange([]string{"not base64!"}); f != "not base64!" || l != "not base64!" {
		t.Errorf("non-base64 fallback: first=%q last=%q", f, l)
	}

	// Empty string in the slice produces empty output rather than panicking.
	if f, l := commitRange([]string{""}); f != "" || l != "" {
		t.Errorf("empty entry: first=%q last=%q", f, l)
	}
}
