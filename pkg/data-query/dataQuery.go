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
  Database string
  Table    string
  Flat     bool

  Fields  []DataQueryDimension
  Filters []DataQueryFilter
  Order   []string

  // used only for grouped query
  Aggregator          string
  Dimensions          []DataQueryDimension
  TimeDimensionFormat string
}

type DataQueryFilter struct {
  Field    string      `json:"field"`
  Value    interface{} `json:"value"`
  Sql      string      `json:"sql"`
  Operator string      `json:"operator"`
}

type DataQueryDimension struct {
  Name string
  Sql  string

  metricPath      string
  metricName      string
  metricValueName rune

  ResultPropertyName string

  arrayJoin string
}

type DatabaseConnectionSupplier interface {
  GetDatabase(name string) (*sqlx.DB, error)
}

func ReadQuery(request *http.Request) ([]DataQuery, bool, error) {
  path := request.URL.Path
  // rison doesn't escape /, so, client should use object notation (i.e. wrap into ())
  // array?
  arrayStart := strings.IndexRune(path, '!')
  objectStart := strings.IndexRune(path, '(')
  var index int
  wrappedAsArray := arrayStart < objectStart
  if wrappedAsArray {
    index = arrayStart
  } else {
    index = objectStart
  }
  if index == -1 {
    return nil, false, errors.New("query not found")
  }

  jsonData, err := rison.ToJSON([]byte(path[index:]), rison.Rison)
  if err != nil {
    return nil, false, errors.WithStack(err)
  }

  list, err := readQuery(jsonData)
  if err != nil {
    return nil, false, err
  }
  return list, wrappedAsArray, nil
}

func SelectRows(query DataQuery, table string, dbSupplier DatabaseConnectionSupplier, queryContext context.Context) (*sql.Rows, int, error) {
  sqlQuery, args, fieldCount, err := buildSql(query, table)
  if err != nil {
    return nil, -1, err
  }

  db, err := dbSupplier.GetDatabase(query.Database)
  if err != nil {
    return nil, -1, err
  }

  rows, err := db.QueryContext(queryContext, sqlQuery, args...)
  if err != nil {
    if err == context.Canceled {
      return nil, -1, err
    } else {
      return nil, -1, errors.WithMessage(err, "cannot execute SQL:\n"+sqlQuery+"\n")
    }
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

  // the only array join is supported for now
  arrayJoin := ""
  for _, dimension := range query.Fields {
    if len(dimension.arrayJoin) != 0 {
      // the only array join is supported for now
      arrayJoin = dimension.arrayJoin
      // for field add distinct to filter duplicates out
      //sb.WriteString(" distinct ")
      break
    }
  }

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
      if len(field.arrayJoin) == 0 {
        sb.WriteString(field.Name)
      } else {
        // measures.values is not a valid field name
        sb.WriteString("measure_value")
      }
    }
  }

  sb.WriteString(" from ")
  sb.WriteString(table)

  if len(arrayJoin) == 0 {
    for _, dimension := range query.Dimensions {
      if len(dimension.arrayJoin) != 0 {
        arrayJoin = dimension.arrayJoin
        break
      }
    }
  }

  if len(arrayJoin) != 0 {
    sb.WriteString(" array join ")
    sb.WriteString(arrayJoin)
  }

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
    // escape - maybe nested name with dot
    sb.WriteRune('`')
    sb.WriteString(dimension.Name)
    sb.WriteRune('`')
  } else {
    sb.WriteString(dimension.Name)
  }
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
