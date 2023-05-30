package server

import (
  "github.com/ClickHouse/clickhouse-go/v2"
  "github.com/valyala/bytebufferpool"
  "github.com/valyala/quicktemplate"
  "net/http"
  "time"
)

func (t *StatsServer) getDistinctHighlightingPasses(request *http.Request) (*bytebufferpool.ByteBuffer, bool, error) {
  sql := "SELECT DISTINCT arrayJoin((arrayFilter(x-> x LIKE 'highlighting/%', `metrics.name`))) as PassName from ij.report"
  buffer := byteBufferPool.Get()
  defer byteBufferPool.Put(buffer)
  templateWriter := quicktemplate.AcquireWriter(buffer)
  defer quicktemplate.ReleaseWriter(templateWriter)
  jsonWriter := templateWriter.N()
  jsonWriter.S("[")

  db, err := clickhouse.Open(&clickhouse.Options{
    Addr: []string{t.dbUrl},
    Auth: clickhouse.Auth{
      Database: "ij",
    },
    DialTimeout:     10 * time.Second,
    ConnMaxLifetime: time.Hour,
    Settings: map[string]interface{}{
      "readonly":         1,
      "max_query_size":   1000000,
      "max_memory_usage": 3221225472,
    },
  })
  var result []struct {
    PassName string
  }
  err = db.Select(request.Context(), &result, sql)
  if err == nil {
    for i, v := range result {
      if i != 0 {
        jsonWriter.S(",")
      }
      jsonWriter.Q(v.PassName)
    }
  }

  jsonWriter.S("]")
  return buffer, true, err
}
