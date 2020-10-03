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
var reFieldName = regexp.MustCompile("^[a-zA-Z_][0-9a-zA-Z_]*$")
var reMetricName = regexp.MustCompile("^[a-zA-Z0-9 _]+$")

// add ().space,'*
var reAggregator = regexp.MustCompile("^[a-zA-Z_][0-9a-zA-Z_(). ,'*<>/]*$")

// for db name the same validation rules as for field name
var reDbName = reFieldName

func ValidateDatabaseName(db string) (string, error) {
  if len(db) == 0 {
    return "default", nil
  } else if !reDbName.MatchString(db) {
    return "", errors.Errorf("Database name %s contains illegal chars", db)
  }
  return db, nil
}

func isValidFieldName(v string) bool {
  return reFieldName.MatchString(v)
}

func readQuery(s []byte, query *DataQuery) error {
  parser := queryParsers.Get()
  defer queryParsers.Put(parser)

  value, err := parser.ParseBytes(s)
  if err != nil {
    return errors.WithStack(err)
  }

  query.Database = string(value.GetStringBytes("db"))

  if len(query.Database) == 0 {
    query.Database = "default"
  } else if !reDbName.MatchString(query.Database) {
    return http_error.NewHttpError(400, fmt.Sprintf("Database name %s contains illegal chars", query.Database))
  }

  err = readDimensions(value.GetArray("fields"), &query.Fields)
  if err != nil {
    return err
  }

  err = readFilters(value.GetArray("filters"), query)
  if err != nil {
    return err
  }

  for _, v := range value.GetArray("order") {
    field := string(v.GetStringBytes())
    if !isValidFieldName(field) {
      return http_error.NewHttpError(400, fmt.Sprintf("Order %s is not a valid field name", field))
    }
    query.Order = append(query.Order, field)
  }

  query.Aggregator = string(value.GetStringBytes("aggregator"))
  if len(query.Aggregator) != 0 && !reAggregator.MatchString(query.Aggregator) {
    return http_error.NewHttpError(400, fmt.Sprintf("Aggregator %s contains illegal chars", query.Aggregator))
  }

  err = readDimensions(value.GetArray("dimensions"), &query.Dimensions)
  if err != nil {
    return err
  }

  query.TimeDimensionFormat = string(value.GetStringBytes("timeDimensionFormat"))
  if len(query.Aggregator) != 0 && !reAggregator.MatchString(query.Aggregator) {
    return http_error.NewHttpError(400, fmt.Sprintf("timeDimensionFormat %s contains illegal chars", query.TimeDimensionFormat))
  }
  return nil
}

func readDimensions(list []*fastjson.Value, result *[]DataQueryDimension) error {
  for _, v := range list {
    var t DataQueryDimension
    if v.Type() == fastjson.TypeString {
      t = DataQueryDimension{
        Name: string(v.GetStringBytes()),
      }
    } else {
      t = DataQueryDimension{
        Name: string(v.GetStringBytes("name")),
        Sql:  string(v.GetStringBytes("sql")),
      }

      if len(t.Sql) != 0 && !reAggregator.MatchString(t.Sql) {
        return http_error.NewHttpError(400, fmt.Sprintf("Dimension SQL %s contains illegal chars", t.Name))
      }
    }

    qualifierDotIndex := strings.IndexRune(t.Name, '.')
    if qualifierDotIndex == -1 {
      if !isValidFieldName(t.Name) {
        return http_error.NewHttpError(400, fmt.Sprintf("Name %s is not a valid field name", t.Name))
      }
    } else {
      t.metricPath = t.Name[0:qualifierDotIndex]
      t.metricName = t.Name[qualifierDotIndex+1:]
      t.metricValueName = 'd'

      metricNameLength := len(t.metricName)
      if metricNameLength > 2 && t.metricName[metricNameLength-2] == '.' {
        t.metricValueName = rune(t.metricName[metricNameLength-1])
        t.metricName = t.metricName[:metricNameLength-2]
      }

      t.ResultPropertyName = strings.ReplaceAll(t.metricName, " ", "_")

      if !isValidFieldName(t.metricPath) || !reMetricName.MatchString(t.metricName) {
        return http_error.NewHttpError(400, fmt.Sprintf("Name %s is not a valid field name", t.Name))
      }
    }

    *result = append(*result, t)
  }
  return nil
}

func readFilters(list []*fastjson.Value, query *DataQuery) error {
  for _, v := range list {
    t := DataQueryFilter{
      Field:    string(v.GetStringBytes("field")),
      Sql:      string(v.GetStringBytes("sql")),
      Operator: string(v.GetStringBytes("operator")),
    }

    value := v.Get("value")
    if value == nil {
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
      t.Value = readStringList(value)
    } else {
      return errors.Errorf("Filter value %v is not supported", value)
    }

    if !isValidFieldName(t.Field) {
      return http_error.NewHttpError(400, fmt.Sprintf("Name %s is not a valid field name", t.Field))
    }

    if len(t.Sql) == 0 {
      if len(t.Operator) == 0 {
        t.Operator = "="
      } else if len(t.Operator) > 2 || (t.Operator != ">" && t.Operator != "<" && t.Operator != "=" && t.Operator != "!=") {
        return errors.Errorf("Operator %s is not supported", t.Operator)
      }
    } else {
      // by intention sql string is not validated
      if len(t.Operator) != 0 {
        return http_error.NewHttpError(400, fmt.Sprintf("sql and operator are mutually exclusive"))
      }
      if t.Value != nil {
        return http_error.NewHttpError(400, fmt.Sprintf("sql and value are mutually exclusive"))
      }
    }

    query.Filters = append(query.Filters, t)
  }
  return nil
}

func readStringList(parentValue *fastjson.Value) []string {
  var list []string
  for _, v := range parentValue.GetArray() {
    list = append(list, string(v.GetStringBytes()))
  }
  return list
}
