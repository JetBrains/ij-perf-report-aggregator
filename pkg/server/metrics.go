package server

import (
  "context"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/data-query"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/util"
  "github.com/pkg/errors"
  "github.com/valyala/quicktemplate"
  "math"
  "net/http"
  "strconv"
  "strings"
)

func (t *StatsServer) handleLoadRequest(request *http.Request) ([]byte, error) {
  return t.doHandleMetricsRequest(request, 2)
}

func (t *StatsServer) handleMetricsRequest(request *http.Request) ([]byte, error) {
  return t.doHandleMetricsRequest(request, 1)
}

func (t *StatsServer) doHandleMetricsRequest(request *http.Request, version int) ([]byte, error) {
  dataQueries, err := data_query.ReadQuery(request)
  if err != nil {
    return nil, err
  }

  buffer := byteBufferPool.Get()
  defer byteBufferPool.Put(buffer)

  templateWriter := quicktemplate.AcquireWriter(buffer)
  defer quicktemplate.ReleaseWriter(templateWriter)
  jsonWriter := templateWriter.N()

  if len(dataQueries) > 1 {
    jsonWriter.S("[")
  }

  for index, dataQuery := range dataQueries {
    if index != 0 {
      jsonWriter.S(",")
    }

    if version == 1 {
      err = t.computeMetricsResponse(dataQuery, jsonWriter, request.Context())
    } else {
      err = t.computeMetricsResponse2(dataQuery, jsonWriter, request.Context())
    }
    if err != nil {
      return nil, err
    }
  }

  if len(dataQueries) > 1 {
    jsonWriter.S("]")
  }
  return CopyBuffer(buffer), nil
}

func (t *StatsServer) computeMetricsResponse(query data_query.DataQuery, jsonWriter *quicktemplate.QWriter, context context.Context) error {
  table := "report"
  if query.Database == "sharedIndexes" {
    table = "metrics"
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

  //isSortedByBuildNumber := len(query.Order) > 1 && query.Order[0] == "build_c1"
  buildNumberColumnOffset := 1
  if len(query.Order) > 1 {
    buildNumberColumnOffset = 0
  }

  var sb strings.Builder

  dataItems := [][]data_query.DataQueryDimension{query.Dimensions, query.Fields}

  isFirst := true
  //lastBuildWithoutUniqueSuffix := ""
  for rows.Next() {
    err := rows.Scan(columnPointers...)
    if err != nil {
      return errors.WithStack(err)
    }

    err = rows.Scan(columnPointers...)
    if err != nil {
      return errors.WithStack(err)
    }

    if isFirst {
      isFirst = false
    } else {
      jsonWriter.S(",")
    }

    if !query.Flat{
      jsonWriter.S(`{`)
    }

    // build number with addition if not unique (build as x coordinate)
    sb.Reset()
    sb.WriteString(strconv.Itoa(int((*(columnPointers[buildNumberColumnOffset].(*interface{}))).(uint8))))
    sb.WriteRune('.')
    sb.WriteString(strconv.Itoa(int((*(columnPointers[buildNumberColumnOffset+1].(*interface{}))).(uint16))))
    buildC3 := int((*(columnPointers[buildNumberColumnOffset+2].(*interface{}))).(uint16))
    if buildC3 != 0 {
      sb.WriteRune('.')
      sb.WriteString(strconv.Itoa(buildC3))
    }

    buildAsString := sb.String()

    jsonWriter.S(`"build":`)
    jsonWriter.Q(buildAsString)

    index := 0
    for _, fields := range dataItems {
      for _, field := range fields {
        if strings.HasPrefix(field.Name, "build_") {
          index++
          continue
        }

        v := *(columnPointers[index].(*interface{}))
        index++

        if v == uint16Zero || v == float32Zero {
          // skip 0 values (0 as null - not existent)
          continue
        }

        if !query.Flat {
          jsonWriter.S(`,"`)
          if len(field.ResultPropertyName) == 0 {
            jsonWriter.S(field.Name)
          } else {
            jsonWriter.S(field.ResultPropertyName)
          }
          jsonWriter.S(`":`)
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
        default:
          return errors.Errorf("unknown type: %T for field %s", untypedValue, field.Name)
        }
      }
    }

    if !query.Flat {
      jsonWriter.S("}")
    }
  }

  jsonWriter.S("]")

  return nil
}

// 2d-table for ECharts
func (t *StatsServer) computeMetricsResponse2(query data_query.DataQuery, jsonWriter *quicktemplate.QWriter, context context.Context) error {
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

        //if v == int32Zero || v == uint16Zero || v == uint32Zero || v == float32Zero {
        //  // skip 0 values (0 as null - not existent)
        //  continue
        //}

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

const uint16Zero = uint16(0)
const uint32Zero = uint32(0)
const float32Zero = float32(0)
const int32Zero = int32(0)
