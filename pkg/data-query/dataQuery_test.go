package data_query

import (
  "github.com/stretchr/testify/assert"
  "testing"
)

// test that config grabs env
func TestAdvancedFilter(t *testing.T) {
  var query DataQuery
  err := readQuery([]byte(`
{
  "filters": [
    {"field": "generated_time", "sql": "> subtractMonths(now(), 1)"}
  ]
}
`), &query)
  if err != nil {
    t.Error(err)
  }

  //noinspection GoImportUsedAsName
  assert := assert.New(t)

  sql, args, _, err := buildSql(query, "test")
  if err != nil {
    t.Error(err)
  }
  assert.Equal("select from test where generated_time > subtractMonths(now(), 1)", sql)
  assert.Empty(args)
}