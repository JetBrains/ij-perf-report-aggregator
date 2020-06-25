package data_query

import (
  "context"
  "database/sql"
  "fmt"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/http-error"
  "github.com/jmoiron/sqlx"
  "github.com/json-iterator/go"
  "github.com/pkg/errors"
  "gopkg.in/sakura-internet/go-rison.v3"
  "math"
  "net/http"
  "regexp"
  "strings"
)

type DataQuery struct {
  Database string `json:"db"`

  Fields  []DataQueryDimension `json:"fields"`
  Filters []DataQueryFilter    `json:"filters"`
  Order   []string             `json:"order"`

  // used only for grouped query
  Aggregator          string               `json:"aggregator"`
  Dimensions          []DataQueryDimension `json:"dimensions"`
  TimeDimensionFormat string               `json:"timeDimensionFormat"`
}

type DataQueryFilter struct {
  Field    string      `json:"field"`
  Value    interface{} `json:"value"`
  Sql      string      `json:"sql"`
  Operator string      `json:"operator"`
}

type DataQueryDimension struct {
  Name string `json:"name"`
  Sql  string `json:"sql"`

  metricPath string
  metricName string
  // tuple index - 3 for duration, 2 for start
  metricValueIndex rune

  ResultPropertyName string
}

type DatabaseConnectionSupplier interface {
  GetDatabase(name string) (*sqlx.DB, error)
}

// https://clickhouse.yandex/docs/en/query_language/syntax/#syntax-identifiers
var reFieldName = regexp.MustCompile("^[a-zA-Z_][0-9a-zA-Z_]*$")
var reMetricName = regexp.MustCompile("^[a-zA-Z0-9 _]+$")
var reDbName = reFieldName

// add ().space,'*
var reAggregator = regexp.MustCompile("^[a-zA-Z_][0-9a-zA-Z_(). ,'*<>/]*$")

func (t *DataQueryDimension) UnmarshalJSON(b []byte) error {
  if b[0] == '"' {
    var s string
    err := jsoniter.ConfigFastest.Unmarshal(b, &s)
    if err != nil {
      return err
    }
    t.Name = s
  } else {
    iterator := jsoniter.ConfigFastest.BorrowIterator(b)
    defer jsoniter.ConfigFastest.ReturnIterator(iterator)
    iterator.ReadObjectCB(func(iterator *jsoniter.Iterator, s string) bool {
      if s == "name" {
        t.Name = iterator.ReadString()
      } else if s == "sql" {
        t.Sql = iterator.ReadString()
      }
      return true
    })
  }

  qualifierDotIndex := strings.IndexRune(t.Name, '.')
  if qualifierDotIndex == -1 {
    if !isValidFieldName(t.Name) {
      return http_error.NewHttpError(400, fmt.Sprintf("Name %s is not a valid field name", t.Name))
    }
  } else {
    t.metricPath = t.Name[0:qualifierDotIndex]
    t.metricName = t.Name[qualifierDotIndex+1:]
    t.metricValueIndex = '3'

    metricNameLength := len(t.metricName)
    if metricNameLength > 2 && t.metricName[metricNameLength - 2] == '.' {
      if t.metricName[metricNameLength - 1] == 's' {
        t.metricValueIndex = '2'
      }
      t.metricName = t.metricName[:metricNameLength - 2]
    }

    t.ResultPropertyName = strings.ReplaceAll(t.metricName, " ", "_")

    if !isValidFieldName(t.metricPath) || !reMetricName.MatchString(t.metricName) {
      return http_error.NewHttpError(400, fmt.Sprintf("Name %s is not a valid field name", t.Name))
    }
  }

  if len(t.Sql) != 0 && !reAggregator.MatchString(t.Sql) {
    return http_error.NewHttpError(400, fmt.Sprintf("Dimension SQL %s contains illegal chars", t.Name))
  }

  return nil
}

func isValidFieldName(v string) bool {
  return reFieldName.MatchString(v)
}

func ReadQuery(request *http.Request) (DataQuery, error) {
  var result DataQuery
  err := readQueryFromRequest(request, &result)
  if err != nil {
    return result, err
  }

  result.Database, err = ValidateDatabaseName(result.Database)
  if err != nil {
    return DataQuery{}, err
  }
  return result, nil
}

func ValidateDatabaseName(db string) (string, error) {
  // for db name the same validation rules as for field name
  if len(db) == 0 {
    return "default", nil
  } else if !reDbName.MatchString(db) {
    return "", errors.Errorf("Database name %s contains illegal chars", db)
  }
  return db, nil
}

func readQueryFromRequest(request *http.Request, v interface{}) error {
  path := request.URL.Path
  // rison doesn't escape /, so, client should use object notation (i.e. wrap into ())
  index := strings.IndexRune(path, '(')
  if index == -1 {
    return errors.New("query not found")
  }

  jsonData, err := rison.ToJSON([]byte(path[index:]), rison.Rison)
  if err != nil {
    return errors.WithStack(err)
  }
  return readQuery(jsonData, v)
}

func readQuery(s []byte, v interface{}) error {
  err := jsoniter.ConfigFastest.Unmarshal(s, v)
  if err != nil {
    return http_error.NewHttpError(400, err.Error())
  }
  return nil
}

func SelectRows(query DataQuery, table string, dbSupplier DatabaseConnectionSupplier, context context.Context) (*sql.Rows, error) {
  sqlQuery, args, err := buildSql(query, table)

  //fmt.Println(sqlQuery)

  if err != nil {
    return nil, err
  }

  db, err := dbSupplier.GetDatabase(query.Database)
  if err != nil {
    return nil, err
  }
  return db.QueryContext(context, sqlQuery, args...)
}

func SelectRow(query DataQuery, table string, dbSupplier DatabaseConnectionSupplier, context context.Context) (*sql.Row, error) {
  sqlQuery, args, err := buildSql(query, table)
  if err != nil {
    return nil, err
  }

  db, err := dbSupplier.GetDatabase(query.Database)
  if err != nil {
    return nil, err
  }
  return db.QueryRowContext(context, sqlQuery, args...), nil
}

func buildSql(query DataQuery, table string) (string, []interface{}, error) {
  var sb strings.Builder
  var args []interface{}

  aggregator := query.Aggregator
  if len(aggregator) != 0 && !reAggregator.MatchString(aggregator) {
    return "", args, errors.Errorf("Aggregator %s contains illegal chars", aggregator)
  }

  sb.WriteString("select")

  if len(query.Dimensions) != 0 {
    writeDimensions(query, &sb)
  }

  for i, field := range query.Fields {
    if i != 0 || len(query.Dimensions) != 0 {
      sb.WriteRune(',')
    }
    sb.WriteRune(' ')

    if len(field.Sql) != 0 {
      writeDimension(field, &sb)
      continue
    }

    if len(aggregator) != 0 {
      sb.WriteString(aggregator)
      sb.WriteRune('(')
    }

    if len(field.metricPath) == 0 {
      sb.WriteString(field.Name)
    } else {
      // select arrayElement(arrayFilter(it -> it.1 = 'start main frontend', JSONExtract(raw_report, 'prepareAppInitActivities', 'Array(Tuple(String, Int, Int))')), 1).3 as v
      // from report;
      sb.WriteString("arrayElement(arrayFilter(it -> it.1 = '")
      sb.WriteString(field.metricName)
      sb.WriteString("', JSONExtract(raw_report, '")
      sb.WriteString(field.metricPath)
      sb.WriteString("', 'Array(Tuple(String, Int, Int))')), 1).")
      sb.WriteRune(field.metricValueIndex)
    }

    if len(aggregator) != 0 {
      sb.WriteRune(')')
    }

    if len(field.ResultPropertyName) != 0 {
      sb.WriteString(" as ")
      sb.WriteString(field.ResultPropertyName)
    } else if len(aggregator) != 0 {
      sb.WriteString(" as ")
      sb.WriteString(field.Name)
    }
  }

  sb.WriteString(" from ")
  sb.WriteString(table)

  if len(query.Filters) != 0 {
    err := writeWhereClause(&sb, query, &args)
    if err != nil {
      return "", args, err
    }
  }

  if len(query.Dimensions) != 0 {
    sb.WriteString(" group by")
    for i, dimension := range query.Dimensions {
      if i != 0 {
        sb.WriteRune(',')
      }
      sb.WriteRune(' ')
      sb.WriteString(dimension.Name)
    }
  }

  if len(query.Order) != 0 {
    sb.WriteString(" order by")
    for i, field := range query.Order {
      if !isValidFieldName(field) {
        return "", args, http_error.NewHttpError(400, fmt.Sprintf("Name %s is not a valid field name", field))
      }

      if i != 0 {
        sb.WriteRune(',')
      }
      sb.WriteRune(' ')
      sb.WriteString(field)
    }
  }

  return sb.String(), args, nil
}

func writeDimensions(query DataQuery, sb *strings.Builder) {
  for i, dimension := range query.Dimensions {
    if i != 0 {
      sb.WriteRune(',')
    }
    sb.WriteRune(' ')
    writeDimension(dimension, sb)
  }
}

func writeDimension(dimension DataQueryDimension, sb *strings.Builder) {
  if len(dimension.Sql) != 0 {
    sb.WriteString(dimension.Sql)
    sb.WriteString(" as ")
  }
  sb.WriteString(dimension.Name)
}

func writeWhereClause(sb *strings.Builder, query DataQuery, args *[]interface{}) error {
  sb.WriteString(" where")
loop:
  for i, filter := range query.Filters {
    if !isValidFieldName(filter.Field) {
      return http_error.NewHttpError(400, fmt.Sprintf("Name %s is not a valid field name", filter.Field))
    }

    if i != 0 {
      sb.WriteString(" and")
    }
    sb.WriteRune(' ')
    sb.WriteString(filter.Field)

    if len(filter.Sql) != 0 {
      if len(filter.Operator) != 0 {
        return errors.Errorf("sql and operator are mutually exclusive")
      }
      if filter.Value != nil {
        return errors.Errorf("sql and value are mutually exclusive")
      }

      sb.WriteRune(' ')
      sb.WriteString(filter.Sql)
      continue loop
    }

    operator := filter.Operator
    if len(operator) == 0 {
      operator = "="
    } else if len(operator) > 2 || (operator != ">" && operator != "<" && operator != "=" && operator != "!=") {
      return errors.Errorf("Operator %s is not supported", operator)
    }

    switch v := filter.Value.(type) {
    case int:
      // default
    case float64:
      if v == math.Trunc(v) {
        // convert to int (to be able to use time unix timestamps from client side)
        filter.Value = int(v)
      }
    case string:
      // default
    case []interface{}:
      sb.WriteString(" in (")
      *args = append(*args, v...)
      for i := 0; i < len(v); i++ {
        if i != 0 {
          sb.WriteString(", ")
        }
        sb.WriteRune('?')
      }
      sb.WriteRune(')')
      continue loop
    default:
      return errors.Errorf("Filter value type %T is not supported", v)
    }

    sb.WriteRune(' ')
    sb.WriteString(operator)
    sb.WriteString(" ?")
    *args = append(*args, filter.Value)
  }
  return nil
}
