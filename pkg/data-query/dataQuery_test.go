package data_query

import (
  "github.com/stretchr/testify/assert"
  "testing"
)

// test that config grabs env
func TestAdvancedFilter(t *testing.T) {
  queries, err := readQuery([]byte(`
{
  "filters": [
    {"f": "generated_time", "sql": "> subtractMonths(now(), 1)"}
  ],
  "order": "generated_time"
}
`))
  if err != nil {
    t.Error(err)
  }

  //noinspection GoImportUsedAsName
  assert := assert.New(t)
  assert.NotEmpty(queries)

  sql, _, err := buildSql(queries[0], "test")
  if err != nil {
    t.Error(err)
  }
  assert.Equal("select from test where generated_time > subtractMonths(now(), 1) order by generated_time", sql)
}

func TestAverageAggregate(t *testing.T) {
  queries, err := readQuery([]byte(`{
  "db": "perfint",
  "table": "phpstorm",
  "fields": [
    {
      "n": "t",
      "sql": "toYYYYMMDD(generated_time)"
    },
    {
      "n": "measures",
      "subName": "value"
    }
  ],
  "filters": [
    {
      "f": "measures.name",
      "v": [
        "responsiveness_time"
      ]
    },
    {
      "f": "branch",
      "v": "master"
    }
  ],
  "flat": false,
  "order": [
    "t"
  ],
  "aggregator": "avg",
  "dimensions": [
    {
      "n": "t"
    }
  ]
}
`))
  if err != nil {
    t.Error(err)
  }

  //noinspection GoImportUsedAsName
  assert := assert.New(t)
  assert.NotEmpty(queries)

  sql, _, err := buildSql(queries[0], "test")
  if err != nil {
    t.Error(err)
  }
  assert.Equal("select toYYYYMMDD(generated_time) as `t`, avg(measures.value) as measure_value from test array join measures where measures.name in ('responsiveness_time') and branch='master' group by t order by t", sql)
}

func TestDecode(t *testing.T) {
  query, err := decodeQuery("KLUv_SAMYQAASGVsbG8genN0ZCEh")

  a := assert.New(t)
  a.NoError(err)
  a.Equal("Hello zstd!!", string(query))
}
