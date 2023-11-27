package server

import (
  "context"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/data-query"
  "github.com/valyala/bytebufferpool"
  "github.com/valyala/quicktemplate"
  "net/http"
)

func (t *StatsServer) handleLoadRequestV2(request *http.Request) (*bytebufferpool.ByteBuffer, bool, error) {
  dataQueries, wrappedAsArray, err := data_query.ReadQueryV2(request)
  if err != nil {
    return nil, false, err
  }
  return t.load(request, dataQueries, wrappedAsArray)
}

func (t *StatsServer) handleLoadRequest(request *http.Request) (*bytebufferpool.ByteBuffer, bool, error) {
  dataQueries, wrappedAsArray, err := data_query.ReadQuery(request)
  if err != nil {
    return nil, false, err
  }

  return t.load(request, dataQueries, wrappedAsArray)
}

func (t *StatsServer) load(request *http.Request, dataQueries []data_query.Query, wrappedAsArray bool) (*bytebufferpool.ByteBuffer, bool, error) {
  buffer := byteBufferPool.Get()
  isOk := false
  defer func() {
    if !isOk {
      byteBufferPool.Put(buffer)
    }
  }()

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

    err := t.computeMeasureResponse(request.Context(), dataQuery, jsonWriter)
    if err != nil {
      return nil, false, err
    }
  }

  if len(dataQueries) > 1 || wrappedAsArray {
    jsonWriter.S("]")
  }
  isOk = true
  if len(buffer.B) == 0 {
    jsonWriter.S("[]")
  }
  return buffer, true, nil
}

func (t *StatsServer) computeMeasureResponse(context context.Context, query data_query.Query, jsonWriter *quicktemplate.QWriter) error {
  table := query.Table
  if len(table) == 0 {
    table = "report"
  }

  err := data_query.SelectRows(context, query, table, t, jsonWriter)
  if err != nil {
    return err
  }
  return nil
}
