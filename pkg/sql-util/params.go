package sql_util

import "strings"

// StringEscaper follows https://clickhouse.com/docs/en/sql-reference/syntax/#syntax-string-literal
var StringEscaper = strings.NewReplacer("\\", "\\\\", "'", "''")
