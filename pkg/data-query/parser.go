package data_query

import (
  "fmt"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/http-error"
  "github.com/develar/errors"
  "github.com/valyala/fastjson"
  "math"
  "regexp"
  "strings"
)

// https://clickhouse.yandex/docs/en/query_language/syntax/#syntax-identifiers
var reFieldName = regexp.MustCompile(`^[a-zA-Z_]\w*$`)

// opposite to reFieldName, dot is supported for nested fields
var reNestedFieldName = regexp.MustCompile(`^[a-zA-Z_][.\da-zA-Z_]*$`)
var reMetricName = regexp.MustCompile(`^[a-zA-Z\d _]+$`)

// add ().space,'*
var reAggregator = regexp.MustCompile(`^[a-zA-Z_(][\da-zA-Z_(). ,'*<>/+]*$`)

// for db name the same validation rules as for field name
var reDbName = reFieldName

var queryParsers fastjson.ParserPool

func isValidFieldName(v string) bool {
  return reFieldName.MatchString(v)
}

func isValidFilterFieldName(v string) bool {
  return reNestedFieldName.MatchString(v)
}

func readQuery(s []byte) ([]DataQuery, error) {
  parser := queryParsers.Get()
  defer queryParsers.Put(parser)

  value, err := parser.ParseBytes(s)
  if err != nil {
    return nil, errors.WithStack(err)
  }

  var queries []DataQuery

  if value.Type() == fastjson.TypeArray {
    for _, v := range value.GetArray() {
      query, err := readQueryValue(v)
      if err != nil {
        return queries, err
      }

      queries = append(queries, *query)
    }
  } else {
    query, err := readQueryValue(value)
    if err != nil {
      return queries, err
    }

    queries = append(queries, *query)
  }

  return queries, nil
}

func readQueryValue(value *fastjson.Value) (*DataQuery, error) {
  query := &DataQuery{
    Database: string(value.GetStringBytes("db")),
    Table:    string(value.GetStringBytes("table")),
    Flat:     value.GetBool("flat"),
  }

  if len(query.Database) == 0 {
    query.Database = "default"
  } else if !reDbName.MatchString(query.Database) {
    return nil, http_error.NewHttpError(400, fmt.Sprintf("Database name %s contains illegal chars", query.Database))
  } else if len(query.Table) > 0 && !reDbName.MatchString(query.Table) {
    return nil, http_error.NewHttpError(400, fmt.Sprintf("Table name %s contains illegal chars", query.Table))
  }

  err := readDimensions(value.GetArray("fields"), &query.Fields)
  if err != nil {
    return nil, err
  }

  err = readFilters(value.GetArray("filters"), query)
  if err != nil {
    return nil, err
  }

  orderValue := value.Get("order")
  if orderValue == nil {
    return nil, http_error.NewHttpError(400, fmt.Sprintf("order is missing"))
  }
  if orderValue.Type() == fastjson.TypeString {
    field := string(orderValue.GetStringBytes())
    if !reNestedFieldName.MatchString(field) {
      return nil, http_error.NewHttpError(400, fmt.Sprintf("Order %s is not a valid field name", field))
    }
    query.Order = []string{field}
  } else {
    for _, v := range value.GetArray("order") {
      field := string(v.GetStringBytes())
      if !reNestedFieldName.MatchString(field) {
        return nil, http_error.NewHttpError(400, fmt.Sprintf("Order %s is not a valid field name", field))
      }
      query.Order = append(query.Order, field)
    }
  }

  query.Aggregator = string(value.GetStringBytes("aggregator"))
  if len(query.Aggregator) != 0 && !reAggregator.MatchString(query.Aggregator) {
    return nil, http_error.NewHttpError(400, fmt.Sprintf("Aggregator %s contains illegal chars", query.Aggregator))
  }

  err = readDimensions(value.GetArray("dimensions"), &query.Dimensions)
  if err != nil {
    return nil, err
  }

  query.TimeDimensionFormat = string(value.GetStringBytes("timeDimensionFormat"))
  if len(query.Aggregator) != 0 && !reAggregator.MatchString(query.Aggregator) {
    return nil, http_error.NewHttpError(400, fmt.Sprintf("timeDimensionFormat %s contains illegal chars", query.TimeDimensionFormat))
  }
  return query, nil
}

func readDimensions(list []*fastjson.Value, result *[]DataQueryDimension) error {
  for _, v := range list {
    t, err := readDimension(v)
    if err != nil {
      return err
    }

    if len(t.Sql) != 0 && !reAggregator.MatchString(t.Sql) {
      return http_error.NewHttpError(400, fmt.Sprintf("Dimension SQL %s contains illegal chars", t.Sql))
    }
    if len(t.ResultPropertyName) != 0 && !isValidFieldName(t.ResultPropertyName) {
      return http_error.NewHttpError(400, fmt.Sprintf("ResultPropertyName %s is not a valid field name", t.Name))
    }
    *result = append(*result, *t)
  }
  return nil
}

func readDimension(v *fastjson.Value) (*DataQueryDimension, error) {
  var t DataQueryDimension
  if v.Type() == fastjson.TypeString {
    t = DataQueryDimension{
      Name: string(v.GetStringBytes()),
    }
  } else {
    subNameValue := v.Get("subName")
    if subNameValue == nil {
      t = DataQueryDimension{
        Name: string(v.GetStringBytes("n")),
        Sql:  string(v.GetStringBytes("sql")),
      }
    } else {
      arrayJoin := string(v.GetStringBytes("n"))
      t = DataQueryDimension{
        Name:               arrayJoin + "." + string(subNameValue.GetStringBytes()),
        arrayJoin:          arrayJoin,
        Sql:                string(v.GetStringBytes("sql")),
        ResultPropertyName: string(v.GetStringBytes("resultKey")),
      }

      if !reNestedFieldName.MatchString(t.Name) {
        return nil, http_error.NewHttpError(400, fmt.Sprintf("Name %s is not a valid field name", t.Name))
      }
      if !isValidFieldName(t.arrayJoin) {
        return nil, http_error.NewHttpError(400, fmt.Sprintf("subName %s is not a valid field name", t.Name))
      }

      return &t, nil
    }
  }

  t.ResultPropertyName = string(v.GetStringBytes("resultKey"))

  qualifierDotIndex := strings.IndexRune(t.Name, '.')
  if qualifierDotIndex != -1 {
    t.metricPath = t.Name[0:qualifierDotIndex]
    t.metricName = t.Name[qualifierDotIndex+1:]
    t.metricValueName = 'd'

    metricNameLength := len(t.metricName)
    if metricNameLength > 2 && t.metricName[metricNameLength-2] == '.' {
      t.metricValueName = rune(t.metricName[metricNameLength-1])
      t.metricName = t.metricName[:metricNameLength-2]
    }

    if len(t.ResultPropertyName) == 0 {
      t.ResultPropertyName = strings.ReplaceAll(t.metricName, " ", "_")
    }

    if !isValidFieldName(t.metricPath) || !reMetricName.MatchString(t.metricName) {
      return nil, http_error.NewHttpError(400, fmt.Sprintf("Name %s is not a valid field name", t.Name))
    }
  }
  return &t, nil
}

func readFilters(list []*fastjson.Value, query *DataQuery) error {
  for _, v := range list {
    t := DataQueryFilter{
      Field:    string(v.GetStringBytes("f")),
      Sql:      string(v.GetStringBytes("q")),
      Operator: string(v.GetStringBytes("o")),
    }

    if len(t.Sql) == 0 {
      t.Sql = string(v.GetStringBytes("sql"))
    }

    value := v.Get("v")
    if value == nil {
      value = v.Get("value")
    }

    if !isValidFilterFieldName(t.Field) {
      return http_error.NewHttpError(400, fmt.Sprintf("%s is not a valid filter field name", t.Field))
    }

    if len(t.Sql) == 0 {
      if value == nil {
        return errors.New("Filter value is not specified")
      } else if value.Type() == fastjson.TypeString {
        t.Value = string(value.GetStringBytes())
      } else if value.Type() == fastjson.TypeNumber {
        number, err := value.Float64()
        if err != nil {
          return errors.WithStack(err)
        }

        if number == math.Trunc(number) {
          // convert to int (to be able to use time unix timestamps from client side)
          t.Value = int(number)
        } else {
          t.Value = number
        }
      } else if value.Type() == fastjson.TypeArray {
        t.Value = readArray(value)
      } else if value.Type() == fastjson.TypeFalse {
        t.Value = value.GetBool()
      } else if value.Type() == fastjson.TypeTrue {
        t.Value = value.GetBool()
      } else {
        return errors.Errorf("Filter value %v is not supported", value)
      }

      if len(t.Operator) == 0 {
        t.Operator = "="
      } else if t.Operator != ">" && t.Operator != "<" && t.Operator != "=" && t.Operator != "!=" && t.Operator != "like" {
        return errors.Errorf("Operator %s is not supported", t.Operator)
      }
    } else {
      // by intention sql string is not validated
      if len(t.Operator) != 0 {
        return http_error.NewHttpError(400, fmt.Sprintf("sql and operator are mutually exclusive (filter=%s)", t.Field))
      }
      if t.Value != nil {
        return http_error.NewHttpError(400, fmt.Sprintf("sql and value are mutually exclusive (filter=%s)", t.Field))
      }
    }

    query.Filters = append(query.Filters, t)
  }
  return nil
}

func readArray(parentValue *fastjson.Value) []interface{} {
  list := make([]interface{}, 0)
  for _, v := range parentValue.GetArray() {
    if v.Type() == fastjson.TypeFalse {
      list = append(list, false)
    } else if v.Type() == fastjson.TypeTrue {
      list = append(list, true)
    } else {
      list = append(list, string(v.GetStringBytes()))
    }
  }
  return list
}
