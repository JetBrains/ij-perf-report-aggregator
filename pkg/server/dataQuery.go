package server

import (
  "context"
  "database/sql"
  "github.com/jmoiron/sqlx"
  "github.com/json-iterator/go"
  "github.com/pkg/errors"
  "math"
  "net/http"
  "regexp"
  "strings"
)

type DataQuery struct {
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
  Operator string      `json:"operator"`
}

type DataQueryDimension struct {
  Name string `json:"name"`
  Sql  string `json:"sql"`
}

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

  if !reFieldName.MatchString(t.Name) {
    return errors.Errorf("Name %s is not a valid field name", t.Name)
  }
  if len(t.Sql) != 0 && !reAggregator.MatchString(t.Sql) {
    return errors.Errorf("Dimension SQL %s contains illegal chars", t.Name)
  }

  return nil
}

func ReadQuery(request *http.Request) (DataQuery, error) {
  path := request.URL.Path
  index := strings.LastIndexByte(path, '/')
  var result DataQuery
  if index != -1 {
    err := jsoniter.ConfigFastest.UnmarshalFromString(path[index+1:], &result)
    if err != nil {
      return result, errors.WithStack(err)
    }
  }
  return result, nil
}

// https://clickhouse.yandex/docs/en/query_language/syntax/#syntax-identifiers
var reFieldName = regexp.MustCompile("^[a-zA-Z_][0-9a-zA-Z_]*$")

// add ().space,'*
var reAggregator = regexp.MustCompile("^[a-zA-Z_][0-9a-zA-Z_(). ,'*]*$")

func SelectData(query DataQuery, table string, db *sqlx.DB, context context.Context) (*sql.Rows, error) {
  var sb strings.Builder

  aggregator := query.Aggregator
  if len(aggregator) != 0 && !reAggregator.MatchString(aggregator) {
    return nil, errors.Errorf("Aggregator %s contains illegal chars", aggregator)
  }

  sb.WriteString("select")

  if len(query.Dimensions) != 0 {
    err := writeDimensions(query, &sb)
    if err != nil {
      return nil, err
    }
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
    sb.WriteString(field.Name)
    if len(aggregator) != 0 {
      sb.WriteString(") as ")
      sb.WriteString(field.Name)
    }
  }

  sb.WriteString(" from ")
  sb.WriteString(table)

  var args []interface{}
  if len(query.Filters) != 0 {
    err := writeWhereClause(&sb, query, &args)
    if err != nil {
      return nil, err
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
      if !reFieldName.MatchString(field) {
        return nil, errors.Errorf("Name %s is not a valid field name", field)
      }

      if i != 0 {
        sb.WriteRune(',')
      }
      sb.WriteRune(' ')
      sb.WriteString(field)
    }
  }

  generatedSql := sb.String()
  return db.QueryContext(context, generatedSql, args...)
}

func writeDimensions(query DataQuery, sb *strings.Builder) error {
  for i, dimension := range query.Dimensions {
    if i != 0 {
      sb.WriteRune(',')
    }
    sb.WriteRune(' ')
    writeDimension(dimension, sb)
  }
  return nil
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
    if !reFieldName.MatchString(filter.Field) {
      return errors.Errorf("Name %s is not a valid field name", filter.Field)
    }

    if i != 0 {
      sb.WriteString(" and")
    }
    sb.WriteRune(' ')
    sb.WriteString(filter.Field)

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
