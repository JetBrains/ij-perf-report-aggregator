package server

import (
  "github.com/asaskevich/govalidator"
  "github.com/develar/errors"
  "go.uber.org/zap"
  "math"
  "net/http"
  "net/url"
  "report-aggregator/pkg/model"
  "report-aggregator/pkg/util"
  "strconv"
  "strings"
)

func (t *StatsServer) handleGroupedMetricsRequest(request *http.Request) ([]byte, error) {
  urlQuery, err := parseQuery(request)
  if err != nil {
    return nil, err
  }

  var query Query

  query.product, query.machine, query.eventType, err = getProductAndMachine(urlQuery)
  if err != nil {
    return nil, err
  }

  bytes, err := validateAndConfigureOperator(&query, urlQuery)
  if err != nil {
    return bytes, err
  }

  var metricNames []string
  if query.eventType == 'd' {
    metricNames = model.EssentialDurationMetricNames
  } else {
    metricNames = model.InstantMetricNames
  }

  results, err := t.getAggregatedResults(metricNames, query)
  if err != nil {
    return nil, err
  }

  buffer := byteBufferPool.Get()
  defer byteBufferPool.Put(buffer)
  WriteGroupedMetricList(buffer, results)
  return CopyBuffer(buffer), nil
}

func validateAndConfigureOperator(query *Query, urlQuery url.Values) ([]byte, error) {
  query.operator = urlQuery.Get("operator")
  if len(query.operator) == 0 {
    query.operator = "quantile"
  } else if !govalidator.IsAlpha(query.operator) {
    return nil, NewHttpError(400, "The operator parameter must contain only letters a-zA-Z")
  }

  operatorArg := urlQuery.Get("operatorArg")
  if query.operator == "quantile" {
    if len(operatorArg) == 0 {
      return nil, NewHttpError(400, "The operatorArg parameter is required")
    } else if !govalidator.IsNumeric(operatorArg) {
      return nil, NewHttpError(400, "The operatorArg parameter must be numeric")
    }

    v, err := strconv.ParseInt(operatorArg, 10, 8)
    if err != nil {
      return nil, NewHttpError(400, "quantile is not correct")
    }
    query.quantile = float64(v) / 100
  }
  return nil, nil
}

type Query struct {
  product   string
  machine   string
  eventType rune

  operator string
  quantile float64
}

func (t *StatsServer) getAggregatedResults(metricNames []string, query Query) ([]MedianResult, error) {
  var sb strings.Builder
  sb.WriteString("select ")
  sb.WriteString("build_c1")

  operator := query.operator
  if operator == "quantile" {
    operator = "quantileTDigest"
  }
  metricNameToValues := make(map[string][]Value)
  for _, name := range metricNames {
    sb.WriteString(", ")
    sb.WriteString(operator)
    if operator == "quantileTDigest" {
      sb.WriteRune('(')
      sb.WriteString(strconv.FormatFloat(query.quantile, 'f', 1, 32))
      sb.WriteRune(')')
    }
    sb.WriteRune('(')
    sb.WriteString(name)
    sb.WriteRune('_')
    sb.WriteRune(query.eventType)
    sb.WriteRune(')')
    sb.WriteString(" as ")
    sb.WriteString(name)
  }
  sb.WriteString(" from report group by build_c1 order by build_c1")

  rows, err := t.db.Query(sb.String())
  if err != nil {
    t.logger.Error("cannot execute SQL", zap.String("query", sb.String()))
    return nil, errors.WithStack(err)
  }
  defer util.Close(rows, t.logger)

  columnPointers := make([]interface{}, 1+len(metricNames))

  for i := range columnPointers {
    columnPointers[i] = new(interface{})
  }

  for rows.Next() {
    err := rows.Scan(columnPointers...)
    if err != nil {
      return nil, err
    }

    groupName := int((*(columnPointers[0].(*interface{}))).(uint8))

    for index, name := range metricNames {
      values, ok := metricNameToValues[name]
      var v int
      switch untypedValue := (*(columnPointers[index+1].(*interface{}))).(type) {
      case float64:
        v = int(math.Round(untypedValue))
      case float32:
        v = int(math.Round(float64(untypedValue)))
      case int32:
        v = int(untypedValue)
      }

      value := Value{buildC1: groupName, value: v}
      if ok {
        metricNameToValues[name] = append(values, value)
      } else {
        metricNameToValues[name] = []Value{value}
      }
    }
  }

  result := make([]MedianResult, 0, len(metricNameToValues))
  for _, name := range metricNames {
    values, ok := metricNameToValues[name]
    if !ok {
      continue
    }

    result = append(result, MedianResult{
      metricName:   name,
      buildToValue: values,
    })
  }

  return result, nil
}
