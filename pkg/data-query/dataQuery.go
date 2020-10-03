package data_query

import (
  "context"
  "database/sql"
  "github.com/jmoiron/sqlx"
  "github.com/pkg/errors"
  "github.com/valyala/fastjson"
  "gopkg.in/sakura-internet/go-rison.v3"
  "math"
  "net/http"
  "strings"
)

var queryParsers fastjson.ParserPool

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

  metricPath      string
  metricName      string
  metricValueName rune

  ResultPropertyName string
}

type DatabaseConnectionSupplier interface {
  GetDatabase(name string) (*sqlx.DB, error)
}

func ReadQuery(request *http.Request) (DataQuery, error) {
  var result DataQuery
  err := readQueryFromRequest(request, &result)
  if err != nil {
    return result, err
  }
  return result, nil
}

func readQueryFromRequest(request *http.Request, v *DataQuery) error {
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

func SelectRows(query DataQuery, table string, dbSupplier DatabaseConnectionSupplier, context context.Context) (*sql.Rows, int, error) {
  sqlQuery, args, fieldCount, err := buildSql(query, table)

  if err != nil {
    return nil, -1, err
  }

  db, err := dbSupplier.GetDatabase(query.Database)
  if err != nil {
    return nil, -1, err
  }
  rows, err := db.QueryContext(context, sqlQuery, args...)
  if err != nil {
    return nil, -1, errors.WithStack(err)
  }
  return rows, fieldCount, nil
}

func SelectRow(query DataQuery, table string, dbSupplier DatabaseConnectionSupplier, context context.Context) (*sql.Row, error) {
  sqlQuery, args, _, err := buildSql(query, table)
  if err != nil {
    return nil, err
  }

  db, err := dbSupplier.GetDatabase(query.Database)
  if err != nil {
    return nil, err
  }
  return db.QueryRowContext(context, sqlQuery, args...), nil
}

func buildSql(query DataQuery, table string) (string, []interface{}, int, error) {
  var sb strings.Builder
  var args []interface{}

  sb.WriteString("select")

  fieldCount := len(query.Dimensions)
  if len(query.Dimensions) != 0 {
    writeDimensions(query, &sb)
  }

  // write extra fields to the end, so, it maybe skipped during serialization
  for i, field := range query.Fields {
    if i != 0 || len(query.Dimensions) != 0 {
      sb.WriteRune(',')
    }
    sb.WriteRune(' ')

    fieldCount++

    if len(field.Sql) != 0 {
      writeDimension(field, &sb)
      continue
    }

    if len(query.Aggregator) != 0 {
      sb.WriteString(query.Aggregator)
      sb.WriteRune('(')
    }

    if len(field.metricPath) == 0 {
      sb.WriteString(field.Name)
    } else {
      // select JSONExtractInt(arrayFirst(it -> JSONExtractString(it, 'n') = 'start main frontend', JSONExtractArrayRaw(raw_report, 'prepareAppInitActivities')), 'd') as v
      // from report;
      if field.metricValueName == 'e' {
        // arraySum(it -> it.1 = 's' or it.1 = 'd' ? it.2 : 0, JSONExtractKeysAndValues(arrayFirst(it -> JSONExtractString(it, 'n') = 'render', JSONExtractArrayRaw(raw_report, 'prepareAppInitActivities')), 'Int'))
        sb.WriteString("arraySum(it -> it.1 = 's' or it.1 = 'd' ? it.2 : 0, JSONExtractKeysAndValues(")
        writeExtractJsonObject(&sb, field)
        sb.WriteString(", 'Int'))")
      } else {
        sb.WriteString("JSONExtractInt(")
        writeExtractJsonObject(&sb, field)
        sb.WriteString(", '")
        sb.WriteRune(field.metricValueName)
        sb.WriteString("')")
      }
    }

    if len(query.Aggregator) != 0 {
      sb.WriteRune(')')
    }

    if len(field.ResultPropertyName) != 0 {
      sb.WriteString(" as ")
      sb.WriteString(field.ResultPropertyName)
    } else if len(query.Aggregator) != 0 {
      sb.WriteString(" as ")
      sb.WriteString(field.Name)
    }
  }

  sb.WriteString(" from ")
  sb.WriteString(table)

  if len(query.Filters) != 0 {
    err := writeWhereClause(&sb, query, &args)
    if err != nil {
      return "", args, -1, err
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
      if i != 0 {
        sb.WriteRune(',')
      }
      sb.WriteRune(' ')
      sb.WriteString(field)
    }
  }

  return sb.String(), args, fieldCount, nil
}

func writeExtractJsonObject(sb *strings.Builder, field DataQueryDimension) {
  sb.WriteString("arrayFirst(it -> JSONExtractString(it, 'n') = '")
  sb.WriteString(field.metricName)
  sb.WriteString("', JSONExtractArrayRaw(raw_report, '")
  sb.WriteString(field.metricPath)
  sb.WriteString("'))")
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
    case []string:
      sb.WriteString(" in (")
      for i := 0; i < len(v); i++ {
        *args = append(*args, v)
        if i != 0 {
          sb.WriteString(", ")
        }
        sb.WriteRune('?')
      }
      sb.WriteRune(')')
      continue loop
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
    sb.WriteString(filter.Operator)
    sb.WriteString(" ?")
    *args = append(*args, filter.Value)
  }
  return nil
}
