package data_query

import (
  "github.com/ClickHouse/ch-go/proto"
  "github.com/develar/errors"
  "github.com/valyala/bytebufferpool"
  "github.com/valyala/quicktemplate"
  "strconv"
)

// separate byte buffer pool - different sizes
var byteBufferPool bytebufferpool.Pool

//nolint:gocyclo
func writeResult(result *proto.Results, columnNameToIndex map[string]int, columnBuffers []*bytebufferpool.ByteBuffer, query DataQuery) error {
  for _, column := range *result {
    columnIndex := columnNameToIndex[column.Name]
    var buffer *bytebufferpool.ByteBuffer
    if columnIndex < len(columnBuffers) {
      buffer = columnBuffers[columnIndex]
    } else {
      return errors.Errorf("invalid columnIndex = %d, it is > len(columnBuffers) = %d", columnIndex, len(columnBuffers))
    }
    if buffer == nil {
      buffer = byteBufferPool.Get()
      columnBuffers[columnIndex] = buffer
    } else {
      _ = buffer.WriteByte(',')
    }

    templateWriter := quicktemplate.AcquireWriter(buffer)
    jsonWriter := templateWriter.N()

    switch data := column.Data.(type) {
    case *proto.ColInt32:
      for i, v := range *data {
        if i != 0 {
          _ = buffer.WriteByte(',')
        }
        buffer.B = strconv.AppendInt(buffer.B, int64(v), 10)
      }
    case *proto.ColInt64:
      for i, v := range *data {
        if i != 0 {
          _ = buffer.WriteByte(',')
        }
        buffer.B = strconv.AppendInt(buffer.B, v, 10)
      }
    case *proto.ColUInt8:
      for i, v := range *data {
        if i != 0 {
          _ = buffer.WriteByte(',')
        }
        buffer.B = strconv.AppendUint(buffer.B, uint64(v), 10)
      }

    case *proto.ColBool:
      for i, v := range *data {
        if i != 0 {
          _ = buffer.WriteByte(',')
        }
        buffer.B = strconv.AppendBool(buffer.B, v)
      }

    case *proto.ColUInt16:
      for i, v := range *data {
        if i != 0 {
          _ = buffer.WriteByte(',')
        }
        buffer.B = strconv.AppendUint(buffer.B, uint64(v), 10)
      }
    case *proto.ColUInt32:
      for i, v := range *data {
        if i != 0 {
          _ = buffer.WriteByte(',')
        }
        buffer.B = strconv.AppendUint(buffer.B, uint64(v), 10)
      }
    case *proto.ColUInt64:
      for i, v := range *data {
        if i != 0 {
          _ = buffer.WriteByte(',')
        }
        buffer.B = strconv.AppendUint(buffer.B, v, 10)
      }

    case *proto.ColFloat32:
      for i, v := range *data {
        if i != 0 {
          _ = buffer.WriteByte(',')
        }
        n := int64(v)
        if float32(n) == v {
          // fast path - just int
          buffer.B = strconv.AppendInt(buffer.B, n, 10)
        } else {
          buffer.B = strconv.AppendFloat(buffer.B, float64(v), 'f', -1, 32)
        }
      }

    case *proto.ColFloat64:
      for i, v := range *data {
        if i != 0 {
          _ = buffer.WriteByte(',')
        }
        n := int64(v)
        if float64(n) == v {
          // fast path - just int
          buffer.B = strconv.AppendInt(buffer.B, n, 10)
        } else {
          buffer.B = strconv.AppendFloat(buffer.B, v, 'f', -1, 64)
        }
      }

    case *proto.ColStr:
      for i, position := range data.Pos {
        if i != 0 {
          _ = buffer.WriteByte(',')
        }
        jsonWriter.Q(string(data.Buf[position.Start:position.End]))
      }

    case *proto.ColDateTime:
      for i, v := range data.Data {
        if i != 0 {
          _ = buffer.WriteByte(',')
        }
        jsonWriter.Q(v.Time().Format(query.TimeDimensionFormat))
      }

    case *proto.ColLowCardinality[string]:
      for i, t := range data.Values {
        if i != 0 {
          _ = buffer.WriteByte(',')
        }
        jsonWriter.Q(t)
      }

    default:
      return errors.Errorf("unsupported column type %T", data)
    }

    quicktemplate.ReleaseWriter(templateWriter)
  }
  return nil
}
