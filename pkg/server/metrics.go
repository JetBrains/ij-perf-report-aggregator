package server

import (
  "context"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/data-query"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/util"
  "github.com/pkg/errors"
  "github.com/valyala/quicktemplate"
  "math"
  "net/http"
  "time"
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

// 2d-table for ECharts
func (t *StatsServer) computeMeasureResponse(query data_query.DataQuery, jsonWriter *quicktemplate.QWriter, context context.Context) error {
  table := query.Table
  if len(table) == 0 {
    table = "report"
  }

  rows, fieldCount, err := data_query.SelectRows(query, table, t, context)
  if err != nil {
    return errors.WithStack(err)
  }

  defer util.Close(rows, t.logger)

  jsonWriter.S("[")

  columnPointers := make([]interface{}, fieldCount)
  for i := range columnPointers {
    columnPointers[i] = new(interface{})
  }

  dataItems := [][]data_query.DataQueryDimension{query.Dimensions, query.Fields}

  isFirstRow := true
  for rows.Next() {
    err := rows.Scan(columnPointers...)
    if err != nil {
      return errors.WithStack(err)
    }

    err = rows.Scan(columnPointers...)
    if err != nil {
      return errors.WithStack(err)
    }

    if isFirstRow {
      isFirstRow = false
    } else {
      jsonWriter.S(",")
    }

    if !query.Flat{
      jsonWriter.S(`[`)
    }

    index := 0
    isFirstColumn := true
    for _, fields := range dataItems {
      for _, field := range fields {
        v := *(columnPointers[index].(*interface{}))
        index++

        if !query.Flat {
          if isFirstColumn {
            isFirstColumn = false
          } else {
            jsonWriter.S(",")
          }
        }

        switch untypedValue := v.(type) {
        case float64:
          jsonWriter.F(math.Round(untypedValue))
        case float32:
          jsonWriter.F(float64(untypedValue))
        case int32:
          jsonWriter.D(int(untypedValue))
        case uint8:
          jsonWriter.D(int(untypedValue))
        case uint16:
          jsonWriter.D(int(untypedValue))
        case uint32:
          jsonWriter.D(int(untypedValue))
        case uint64:
          jsonWriter.DL(int64(untypedValue))
        case int64:
          jsonWriter.DL(untypedValue)
        case string:
          jsonWriter.Q(untypedValue)
        case time.Time:
          jsonWriter.Q(untypedValue.Format(query.TimeDimensionFormat))
        default:
          return errors.Errorf("unknown type: %T for field %s", untypedValue, field.Name)
        }
      }
    }

    if !query.Flat {
      jsonWriter.S("]")
    }
  }

  jsonWriter.S("]")

  return nil
}