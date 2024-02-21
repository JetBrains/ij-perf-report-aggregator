package data_query

import (
  "errors"
  "fmt"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/http-error"
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
var reAggregator = regexp.MustCompile(`^[a-zA-Z_'(][\da-zA-Z_(). ,'*@<>\-/+]*$`)

// for db name the same validation rules as for field name
var reDbName = reFieldName

var queryParsers fastjson.ParserPool

func isValidFieldName(v string) bool {
  return reFieldName.MatchString(v)
}

func isValidFilterFieldName(v string) bool {
  return reNestedFieldName.MatchString(v)
}

func readQuery(s []byte) ([]Query, error) {
  parser := queryParsers.Get()
  defer queryParsers.Put(parser)

  value, err := parser.ParseBytes(s)
  if err != nil {
    return nil, fmt.Errorf("cannot parse query: %w", err)
  }

  var queries []Query

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

func readQueryValue(value *fastjson.Value) (*Query, error) {
  query := &Query{
    Database: string(value.GetStringBytes("db")),
    Table:    string(value.GetStringBytes("table")),
    Flat:     value.GetBool("flat"),
  }

  switch {
  case len(query.Database) == 0:
    query.Database = "default"
  case !reDbName.MatchString(query.Database):
    return nil, http_error.NewHttpError(400, fmt.Sprintf("Database name %s contains illegal chars", query.Database))
  case len(query.Table) > 0 && !reDbName.MatchString(query.Table):
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
    return nil, http_error.NewHttpError(400, "order is missing")
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

func readDimensions(list []*fastjson.Value, result *[]QueryDimension) error {
  for _, v := range list {
    t, err := readDimension(v)
    if err != nil {
      return err
    }

    if len(t.Sql) != 0 && !reAggregator.MatchString(t.Sql) {
      return http_error.NewHttpError(400, fmt.Sprintf("Dimension SQL %s contains illegal chars", t.Sql))
    }
    if len(t.resultPropertyName) != 0 && !isValidFieldName(t.resultPropertyName) {
      return http_error.NewHttpError(400, fmt.Sprintf("resultPropertyName %s is not a valid field name", t.Name))
    }
    *result = append(*result, *t)
  }
  return nil
}

func readDimension(v *fastjson.Value) (*QueryDimension, error) {
  var t QueryDimension
  if v.Type() == fastjson.TypeString {
    t = QueryDimension{
      Name: string(v.GetStringBytes()),
    }
  } else {
    subNameValue := v.Get("subName")
    if subNameValue == nil {
      t = QueryDimension{
        Name: string(v.GetStringBytes("n")),
        Sql:  string(v.GetStringBytes("sql")),
      }
    } else {
      arrayJoin := string(v.GetStringBytes("n"))
      t = QueryDimension{
        Name:               arrayJoin + "." + string(subNameValue.GetStringBytes()),
        arrayJoin:          arrayJoin,
        Sql:                string(v.GetStringBytes("sql")),
        resultPropertyName: string(v.GetStringBytes("resultKey")),
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

  t.resultPropertyName = string(v.GetStringBytes("resultKey"))

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

    if len(t.resultPropertyName) == 0 {
      t.resultPropertyName = strings.ReplaceAll(t.metricName, " ", "_")
    }

    if !isValidFieldName(t.metricPath) || !reMetricName.MatchString(t.metricName) {
      return nil, http_error.NewHttpError(400, fmt.Sprintf("Name %s is not a valid field name", t.Name))
    }
  }
  return &t, nil
}

func readFilters(list []*fastjson.Value, query *Query) error {
  for _, v := range list {
    t := QueryFilter{
      Field:    string(v.GetStringBytes("f")),
      Sql:      string(v.GetStringBytes("q")),
      Operator: string(v.GetStringBytes("o")),
      Split:    false,
    }

    if len(t.Sql) == 0 {
      t.Sql = string(v.GetStringBytes("sql"))
    }

    value := v.Get("v")
    t.Split = v.GetBool("s")
    if value == nil {
      value = v.Get("value")
    }

    if !isValidFilterFieldName(t.Field) && len(t.Sql) == 0 {
      return http_error.NewHttpError(400, t.Field+" is not a valid filter field name")
    }

    if len(t.Sql) == 0 && value == nil {
      return errors.New("filter value is not specified")
    }
    if len(t.Sql) == 0 {
      switch value.Type() {
      case fastjson.TypeString:
        t.Value = string(value.GetStringBytes())
      case fastjson.TypeNumber:
        number, err := value.Float64()
        if err != nil {
          return fmt.Errorf("cannot parse filter value %s: %w", value, err)
        }
        if number == math.Trunc(number) {
          // convert to int (to be able to use time unix timestamps from client side)
          t.Value = int(number)
        } else {
          t.Value = number
        }
      case fastjson.TypeArray:
        t.Value = readArray(value)
      case fastjson.TypeFalse:
        t.Value = value.GetBool()
      case fastjson.TypeTrue:
        t.Value = value.GetBool()
      default:
        return fmt.Errorf("filter value %v is not supported", value)
      }

      if len(t.Operator) == 0 {
        t.Operator = "="
      } else if t.Operator != ">" && t.Operator != "<" && t.Operator != "=" && t.Operator != "!=" && t.Operator != "like" {
        return fmt.Errorf("operator %s is not supported", t.Operator)
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
    switch v.Type() {
    case fastjson.TypeFalse:
      list = append(list, false)
    case fastjson.TypeTrue:
      list = append(list, true)
    default:
      list = append(list, string(v.GetStringBytes()))
    }
  }
  return list
}
