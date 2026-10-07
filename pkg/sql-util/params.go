package sql_util

import (
	"fmt"
	"strings"
)

// StringEscaper follows https://clickhouse.com/docs/en/sql-reference/syntax/#syntax-string-literal
var StringEscaper = strings.NewReplacer("\\", "\\\\", "'", "''")

// ValidateIdentifier accepts a plain database or table name, so it can be put into SQL without quoting.
func ValidateIdentifier(field, value string) error {
	if value == "" {
		return fmt.Errorf("%s is required", field)
	}
	for _, r := range value {
		if !(r == '_' || (r >= '0' && r <= '9') || (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')) {
			return fmt.Errorf("%s contains invalid character %q", field, r)
		}
	}
	return nil
}
