package server

import (
  "context"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/util"
  "github.com/develar/errors"
  "math"
  "net/http"
  "strconv"
  "time"
)

func (t *StatsServer) handleGroupedMetricsRequest(request *http.Request) ([]byte, error) {
  dataQuery, err := ReadQuery(request)
  if err != nil {
    return nil, err
  }

  results, err := t.getAggregatedResults(dataQuery, request.Context())
  if err != nil {
    return nil, err
  }

  buffer := byteBufferPool.Get()
  defer byteBufferPool.Put(buffer)
  WriteGroupedMetricList(buffer, results)
  return CopyBuffer(buffer), nil
}

func (t *StatsServer) getAggregatedResults(dataQuery DataQuery, requestContext context.Context) ([]MedianResult, error) {
  rows, err := SelectData(dataQuery, "report", t.db, requestContext)
  if err != nil {
    return nil, errors.WithStack(err)
  }

  defer util.Close(rows, t.logger)

  valueColumnOffset := len(dataQuery.Dimensions)
  columnPointers := make([]interface{}, valueColumnOffset+len(dataQuery.Fields))

  for i := range columnPointers {
    columnPointers[i] = new(interface{})
  }

  metricNameToValues := make(map[string][]Value)

  for rows.Next() {
    err := rows.Scan(columnPointers...)
    if err != nil {
      return nil, err
    }

    var groupName string
    switch v := (*(columnPointers[0].(*interface{}))).(type) {
    case time.Time:
      groupName = v.Format(dataQuery.TimeDimensionFormat)
    case uint8:
      groupName = strconv.FormatInt(int64(v), 10)
    case uint16:
      groupName = strconv.FormatInt(int64(v), 10)
    case int:
      groupName = strconv.FormatInt(int64(v), 10)
    case string:
      groupName = v
    default:
      return nil, errors.Errorf("unknown type: %T", v)
    }

    for index, field := range dataQuery.Fields {
      values, ok := metricNameToValues[field.Name]
      var v int
      switch untypedValue := (*(columnPointers[valueColumnOffset+index].(*interface{}))).(type) {
      case float64:
        v = int(math.Round(untypedValue))
      case float32:
        v = int(math.Round(float64(untypedValue)))
      case int32:
        v = int(untypedValue)
      case uint16:
        v = int(untypedValue)
      case int:
        v = untypedValue
      default:
        return nil, errors.Errorf("unknown type: %T", untypedValue)
      }

      value := Value{group: groupName, value: v}
      if ok {
        metricNameToValues[field.Name] = append(values, value)
      } else {
        metricNameToValues[field.Name] = []Value{value}
      }
    }
  }

  result := make([]MedianResult, 0, len(metricNameToValues))
  for _, field := range dataQuery.Fields {
    values, ok := metricNameToValues[field.Name]
    if !ok {
      continue
    }

    result = append(result, MedianResult{
      metricName:    field.Name,
      groupedValues: values,
    })
  }

  return result, nil
}