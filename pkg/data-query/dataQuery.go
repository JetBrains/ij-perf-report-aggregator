package data_query

import (
  "context"
  _ "embed"
  "encoding/base64"
  "github.com/develar/errors"
  "github.com/go-faster/ch"
  "github.com/go-faster/ch/proto"
  "github.com/klauspost/compress/zstd"
  "github.com/sakura-internet/go-rison/v4"
  "github.com/valyala/bytebufferpool"
  "github.com/valyala/quicktemplate"
  "math"
  "net/http"
  "strconv"
  "strings"
)

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
  Field    string
  Value    interface{}
  Sql      string
  Operator string
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

//go:embed zstd.dictionary
var ZstdDictionary []byte

func decodeQuery(encoded string) ([]byte, error) {
  compressed, err := base64.RawURLEncoding.DecodeString(encoded)
  if err != nil {
    return nil, errors.WithStack(err)
  }

  reader, err := zstd.NewReader(nil, zstd.WithDecoderConcurrency(0), zstd.WithDecoderDicts(ZstdDictionary))
  if err != nil {
    return nil, errors.WithStack(err)
  }
  defer reader.Close()

  decompressed, err := reader.DecodeAll(compressed, nil)
  if err != nil {
    return nil, errors.WithStack(err)
  }
  return decompressed, nil
}

func ReadQueryV2(request *http.Request) ([]DataQuery, bool, error) {
  decompressed, err := decodeQuery(request.URL.Path[len("/api/q/"):])
  if err != nil {
    return nil, false, errors.WithStack(err)
  }

  if len(decompressed) == 0 {
    rawPath := request.URL.RawPath
    return nil, false, errors.New("query not found: " + rawPath)
  }

  wrappedAsArray := decompressed[0] == '['
  parser := queryParsers.Get()
  defer queryParsers.Put(parser)

  //fileName := strconv.FormatUint(xxh3.HashString(request.URL.Path), 36) + ".json"
  //_ = os.WriteFile("/Volumes/data/queries/"+fileName, decompressed, 0644)

  list, err := readQuery(decompressed)
  if err != nil {
    return nil, false, err
  }
  return list, wrappedAsArray, nil
}

func ReadQuery(request *http.Request) ([]DataQuery, bool, error) {
  payload := request.URL.Path

  // array?
  arrayStart := strings.IndexRune(payload, '!')
  objectStart := strings.IndexRune(payload, '(')
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

  jsonData, err := rison.ToJSON([]byte(payload[index:]), rison.Rison)
  if err != nil {
    return nil, false, errors.WithStack(err)
  }

  list, err := readQuery(jsonData)
  if err != nil {
    return nil, false, err
  }
  return list, wrappedAsArray, nil
}

func SelectRows(query DataQuery, table string, dbSupplier DatabaseConnectionSupplier, totalWriter *quicktemplate.QWriter, ctx context.Context) error {
  sqlQuery, columnNameToIndex, err := buildSql(query, table)
  if err != nil {
    return err
  }

  columnBuffers := make([]*bytebufferpool.ByteBuffer, len(columnNameToIndex))

  err = executeQuery(sqlQuery, query, dbSupplier, ctx, func(ctx context.Context, block proto.Block, result *proto.Results) error {
    if block.Rows == 0 {
      return nil
    }
    return writeResult(result, columnNameToIndex, columnBuffers, query)
  })
  if err != nil {
    return err
  }

  if !query.Flat {
    totalWriter.S("[")
  }
  for columnIndex, buffer := range columnBuffers {
    if columnIndex != 0 {
      totalWriter.S(",")
    }

    totalWriter.S("[")

    if buffer != nil {
      _, _ = buffer.WriteTo(totalWriter)
      byteBufferPool.Put(buffer)
    }

    totalWriter.S("]")
  }
  if !query.Flat {
    totalWriter.S("]")
  }
  return nil
}

func SelectString(query DataQuery, table string, dbSupplier DatabaseConnectionSupplier, context context.Context) ([]byte, error) {
  sqlQuery, _, err := buildSql(query, table)
  if err != nil {
    return nil, err
  }

  dbResource, err := dbSupplier.AcquireDatabase(query.Database, context)
  if err != nil {
    return nil, err
  }

  var res proto.ColStr
  err = dbResource.Value().Do(context, ch.Query{Body: sqlQuery, Result: proto.ResultColumn{Data: &res}})
  if err != nil {
    return nil, err
  }

  defer dbResource.Release()
  return res.Buf[res.Pos[0].Start:res.Pos[0].End], nil
}

func buildSql(query DataQuery, table string) (string, map[string]int, error) {
  var sb strings.Builder

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

  columnNameToIndex := make(map[string]int, len(query.Dimensions)+len(query.Fields))
  columnIndex := 0

  if len(query.Dimensions) != 0 {
    sb.WriteRune(' ')
    for i, dimension := range query.Dimensions {
      if i != 0 {
        sb.WriteRune(',')
      }

      columnNameToIndex[dimension.Name] = columnIndex
      columnIndex++

      writeDimension(dimension, &sb)
    }
  }

  // write extra fields to the end, so, it maybe skipped during serialization
  for i, field := range query.Fields {
    if i != 0 || len(query.Dimensions) != 0 {
      sb.WriteRune(',')
    }
    sb.WriteRune(' ')

    if len(field.Sql) != 0 {
      columnNameToIndex[field.Name] = columnIndex
      columnIndex++
      writeDimension(field, &sb)
      continue
    }

    var effectiveColumnName = ""

    if len(query.Aggregator) != 0 {
      sb.WriteString(query.Aggregator)
      sb.WriteRune('(')
    }

    if len(field.metricPath) == 0 {
      sb.WriteString(field.Name)
      effectiveColumnName = field.Name
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
      effectiveColumnName = field.ResultPropertyName
    } else if len(query.Aggregator) != 0 {
      sb.WriteString(" as ")
      if len(field.arrayJoin) == 0 {
        effectiveColumnName = field.Name
      } else {
        // measures.values is not a valid field name
        effectiveColumnName = "measure_value"
      }
      sb.WriteString(effectiveColumnName)
    }

    columnNameToIndex[effectiveColumnName] = columnIndex
    columnIndex++
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
    err := writeWhereClause(&sb, query)
    if err != nil {
      return "", nil, err
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

  return sb.String(), columnNameToIndex, nil
}

func writeExtractJsonObject(sb *strings.Builder, field DataQueryDimension) {
  sb.WriteString("arrayFirst(it -> JSONExtractString(it, 'n') = '")
  sb.WriteString(field.metricName)
  sb.WriteString("', JSONExtractArrayRaw(raw_report, '")
  sb.WriteString(field.metricPath)
  sb.WriteString("'))")
}

func writeDimension(dimension DataQueryDimension, sb *strings.Builder) {
  if len(dimension.Sql) == 0 {
    sb.WriteString(dimension.Name)
  } else {
    sb.WriteString(dimension.Sql)
    sb.WriteString(" as ")
    // escape - maybe nested name with dot
    sb.WriteRune('`')
    sb.WriteString(dimension.Name)
    sb.WriteRune('`')
  }
}

// https://clickhouse.com/docs/en/sql-reference/syntax/#syntax-string-literal
var stringEscaper = strings.NewReplacer("\\", "\\\\", "'", "''")

func writeWhereClause(sb *strings.Builder, query DataQuery) error {
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
      sb.WriteString(filter.Operator)
      sb.WriteString(strconv.Itoa(filter.Value.(int)))
    case float64:
      sb.WriteString(filter.Operator)
      if v == math.Trunc(v) {
        sb.WriteString(strconv.Itoa(int(v)))
      } else {
        sb.WriteString(strconv.FormatFloat(v, 'f', -1, 64))
      }
    case string:
      sb.WriteString(filter.Operator)
      writeString(sb, v)
    case []string:
      sb.WriteString(" in (")
      for j := 0; j < len(v); j++ {
        if j != 0 {
          sb.WriteRune(',')
        }
        writeString(sb, v[j])
      }
      sb.WriteRune(')')
    case []interface{}:
      sb.WriteString(" in (")
      for j := 0; j < len(v); j++ {
        if j != 0 {
          sb.WriteRune(',')
        }
        writeString(sb, v[j].(string))
      }
      sb.WriteRune(')')
    default:
      return errors.Errorf("Filter value type %T is not supported", v)
    }
  }
  return nil
}

func writeString(sb *strings.Builder, s string) {
  sb.WriteByte('\'')
  _, _ = stringEscaper.WriteString(sb, s)
  sb.WriteByte('\'')
}

//var fileLogger *zap.Logger
//
//func init() {
//  var cfg = zap.NewDevelopmentConfig()
//  cfg.EncoderConfig.EncodeTime = func(time time.Time, encoder zapcore.PrimitiveArrayEncoder) {
//  }
//  cfg.EncoderConfig.EncodeLevel = func(level zapcore.Level, encoder zapcore.PrimitiveArrayEncoder) {
//  }
//  cfg.DisableCaller = true
//  cfg.OutputPaths = []string{
//    "",
//  }
//  fileLogger, _ = cfg.Build()
//}
