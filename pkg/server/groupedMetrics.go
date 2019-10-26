package server

import (
  "github.com/asaskevich/govalidator"
  "github.com/develar/errors"
  "github.com/jmoiron/sqlx"
  "go.uber.org/zap"
  "math"
  "net/http"
  "net/url"
  "github.com/JetBrains/ij-perf-report-aggregator/common/model"
  "github.com/JetBrains/ij-perf-report-aggregator/common/util"
  "strconv"
  "strings"
  "time"
)

func (t *StatsServer) handleGroupedMetricsRequest(request *http.Request) ([]byte, error) {
  urlQuery, err := parseQuery(request)
  if err != nil {
    return nil, err
  }

  var query GroupedMetricQuery

  err = validateAndConfigureOperator(&query, urlQuery)
  if err != nil {
    return nil, err
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

func validateAndConfigureOperator(query *GroupedMetricQuery, urlQuery url.Values) error {
  var err error

  query.product, query.machines, query.eventType, err = getProductAndMachine(urlQuery)
  if err != nil {
    return err
  }

  query.operator = urlQuery.Get("operator")
  if len(query.operator) == 0 {
    query.operator = "quantile"
  } else if !govalidator.IsAlpha(query.operator) {
    return NewHttpError(400, "The operator parameter must contain only letters a-zA-Z")
  }

  operatorArg := urlQuery.Get("operatorArg")
  if query.operator == "quantile" {
    if len(operatorArg) == 0 {
      return NewHttpError(400, "The operatorArg parameter is required")
    } else if !govalidator.IsNumeric(operatorArg) {
      return NewHttpError(400, "The operatorArg parameter must be numeric")
    }

    v, err := strconv.ParseInt(operatorArg, 10, 8)
    if err != nil {
      return NewHttpError(400, "quantile is not correct")
    }
    query.quantile = float64(v) / 100
  }
  return nil
}

type GroupedMetricQuery struct {
  BaseMetricQuery

  operator string
  quantile float64
}

func (t *StatsServer) getAggregatedResults(metricNames []string, query GroupedMetricQuery) ([]MedianResult, error) {
  whereStatement, whereArgs, err := sqlx.In(" where product = ? and machine in(?)", query.product, query.machines)
  if err != nil {
    return nil, errors.WithStack(err)
  }

  // assume that for one branch reports generated in sequence (later, when reports from logs will be supported again, some settings can be added)
  var countOfBranches int
  err = t.db.QueryRow("select uniq(branch) from report "+whereStatement, whereArgs...).Scan(&countOfBranches)
  if err != nil {
    return nil, errors.WithStack(err)
  }

  groupByMonth := countOfBranches == 1
  sql := buildSql(query, whereStatement, metricNames, groupByMonth)

  rows, err := t.db.Query(sql, whereArgs...)
  if err != nil {
    t.logger.Error("cannot execute", zap.String("query", sql))
    return nil, errors.WithStack(err)
  }
  defer util.Close(rows, t.logger)

  columnPointers := make([]interface{}, 1+len(metricNames))

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
    if groupByMonth {
      // do not use "Jan 06" because not clear - 06 here it is month or year
      groupName = ((*(columnPointers[0].(*interface{}))).(time.Time)).Format("Jan")
    } else {
      groupName = strconv.FormatInt(int64((*(columnPointers[0].(*interface{}))).(uint8)), 10)
    }

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
      case uint16:
        v = int(untypedValue)
      case int:
        v = untypedValue
      default:
        return nil, errors.Errorf("unknown type: %v", untypedValue)
      }

      value := Value{group: groupName, value: v}
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
      metricName:    name,
      groupedValues: values,
    })
  }

  return result, nil
}

func buildSql(query GroupedMetricQuery, whereStatement string, metricNames []string, groupByMonth bool) string {
  var sb strings.Builder
  sb.WriteString("select ")

  var groupField string
  if groupByMonth {
    groupField = "yearAndMonth"
    sb.WriteString("toStartOfMonth(generated_time) as ")
    sb.WriteString(groupField)
  } else {
    groupField = "build_c1"
    sb.WriteString(groupField)
  }

  operator := query.operator
  if operator == "quantile" {
    operator = "quantileTDigest"
  }
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
  sb.WriteString(" from report ")
  sb.WriteString(whereStatement)
  sb.WriteString(" group by ")
  sb.WriteString(groupField)
  sb.WriteString(" order by ")
  sb.WriteString(groupField)

  sql := sb.String()
  return sql
}
