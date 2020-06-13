package server

import (
  "context"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/data-query"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/util"
  "github.com/pkg/errors"
  "github.com/valyala/quicktemplate"
  "io"
  "math"
  "net/http"
  "strconv"
  "strings"
)

func (t *StatsServer) handleMetricsRequest(request *http.Request) ([]byte, error) {
  dataQuery, err := data_query.ReadQuery(request)
  if err != nil {
    return nil, err
  }

  buffer := byteBufferPool.Get()
  defer byteBufferPool.Put(buffer)
  err = t.computeMetricsResponse(dataQuery, buffer, request.Context())
  if err != nil {
    return nil, err
  }
  return CopyBuffer(buffer), nil
}

func (t *StatsServer) computeMetricsResponse(query data_query.DataQuery, writer io.Writer, context context.Context) error {
  rows, err := data_query.SelectRows(query, "report", t, context)
  if err != nil {
    return errors.WithStack(err)
  }

  defer util.Close(rows, t.logger)

  templateWriter := quicktemplate.AcquireWriter(writer)
  defer quicktemplate.ReleaseWriter(templateWriter)
  jsonWriter := templateWriter.N()
  jsonWriter.S("[")

  dimensionCount := len(query.Dimensions)
  columnPointers := make([]interface{}, dimensionCount+ len(query.Fields))
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

    // timestamp
    jsonWriter.S(`{`)

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
    //if isSortedByBuildNumber {
    //  // https://www.amcharts.com/docs/v4/tutorials/handling-repeating-categories-on-category-axis/
    //  if lastBuildWithoutUniqueSuffix == buildAsString {
    //    // not unique - add time
    //    sb.WriteRune(' ')
    //    sb.WriteRune('(')
    //    sb.WriteString(strconv.FormatInt(generatedTime, 10))
    //    sb.WriteRune(')')
    //    buildAsString = sb.String()
    //  } else {
    //    lastBuildWithoutUniqueSuffix = buildAsString
    //  }
    //}

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

        jsonWriter.S(`,"`)
        jsonWriter.S(field.Name)
        jsonWriter.S(`":`)

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
        case uint64:
          jsonWriter.DL(int64(untypedValue))
        default:
          return errors.Errorf("unknown type: %T for field %s", untypedValue, field.Name)
        }
      }
    }

    jsonWriter.S("}")
  }

  jsonWriter.S("]")

  return nil
}

const uint16Zero = uint16(0)
const float32Zero = float32(0)