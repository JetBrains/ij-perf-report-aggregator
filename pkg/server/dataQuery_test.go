package server

import (
  "github.com/stretchr/testify/assert"
  "testing"
)

// test that config grabs env
func TestAdvancedFilter(t *testing.T) {
  query, err := readQuery(`
{
  "filters": [
    {"field": "generated_time", "sql": "> subtractMonths(now(), 1)"}
  ]
}
`)
  if err != nil {
    t.Error(err)
  }

  //noinspection GoImportUsedAsName
  assert := assert.New(t)

  sql, args, err := buildSql(query, "test")
  assert.Equal("select from test where generated_time > subtractMonths(now(), 1)", sql)
  assert.Empty(args)
}

func TestLimit(t *testing.T) {
  query, err := readQuery(`
{
  "limit": 1
}
`)
  if err != nil {
    t.Error(err)
  }

  //noinspection GoImportUsedAsName
  assert := assert.New(t)

  sql, args, err := buildSql(query, "test")
  assert.Equal("select from test limit 1", sql)
  assert.Empty(args)
}
