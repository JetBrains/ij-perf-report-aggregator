package sql_util

import "strings"

// https://clickhouse.com/docs/en/sql-reference/syntax/#syntax-string-literal
var StringEscaper = strings.NewReplacer("\\", "\\\\", "'", "''")
