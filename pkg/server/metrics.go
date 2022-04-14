package server

import (
  "context"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/data-query"
  "github.com/valyala/quicktemplate"
  "net/http"
)

func (t *StatsServer) handleLoadRequest(request *http.Request) ([]byte, error) {
  dataQueries, wrappedAsArray, err := data_query.ReadQuery(request)
  if err != nil {
    return nil, err
  }

  buffer := byteBufferPool.Get()
  defer byteBufferPool.Put(buffer)

  templateWriter := quicktemplate.AcquireWriter(buffer)
  defer quicktemplate.ReleaseWriter(templateWriter)
  jsonWriter := templateWriter.N()

  if len(dataQueries) > 1 || wrappedAsArray {
    jsonWriter.S("[")
  }

  for index, dataQuery := range dataQueries {
    if index != 0 {
      jsonWriter.S(",")
    }

    err = t.computeMeasureResponse(dataQuery, jsonWriter, request.Context())
    if err != nil {
      return nil, err
    }
  }

  if len(dataQueries) > 1 || wrappedAsArray {
    jsonWriter.S("]")
  }
  return CopyBuffer(buffer), nil
}

func (t *StatsServer) computeMeasureResponse(query data_query.DataQuery, jsonWriter *quicktemplate.QWriter, context context.Context) error {
  table := query.Table
  if len(table) == 0 {
    table = "report"
  }

  err := data_query.SelectRows(query, table, t, jsonWriter, context)
  if err != nil {
    return err
  }
  return nil
}
