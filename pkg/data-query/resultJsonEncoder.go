package data_query

import (
  "github.com/develar/errors"
  "github.com/go-faster/ch/proto"
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
    buffer := columnBuffers[columnIndex]
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
        //if v {
        //  buffer.B = strconv.AppendUint(buffer.B, 1, 10)
        //} else {
        //  buffer.B = strconv.AppendUint(buffer.B, 0, 10)
        //}
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
      for i, v := range *data {
        if i != 0 {
          _ = buffer.WriteByte(',')
        }
        jsonWriter.Q(v.Time().Format(query.TimeDimensionFormat))
      }

    case *proto.ColLowCardinality:
      str := data.Index.(*proto.ColStr)

      switch data.Key {
      case proto.KeyUInt8:
        for i, key := range data.Keys8 {
          if i != 0 {
            _ = buffer.WriteByte(',')
          }
          position := str.Pos[key]
          jsonWriter.Q(string(str.Buf[position.Start:position.End]))
        }
      case proto.KeyUInt16:
        for i, key := range data.Keys16 {
          if i != 0 {
            _ = buffer.WriteByte(',')
          }
          position := str.Pos[key]
          jsonWriter.Q(string(str.Buf[position.Start:position.End]))
        }
      case proto.KeyUInt32:
        for i, key := range data.Keys32 {
          if i != 0 {
            _ = buffer.WriteByte(',')
          }
          position := str.Pos[key]
          jsonWriter.Q(string(str.Buf[position.Start:position.End]))
        }
      case proto.KeyUInt64:
        for i, key := range data.Keys64 {
          if i != 0 {
            _ = buffer.WriteByte(',')
          }
          position := str.Pos[key]
          jsonWriter.Q(string(str.Buf[position.Start:position.End]))
        }
      default:
        return errors.Errorf("unsupported key %d", data.Key)
      }

    default:
      return errors.Errorf("unsupported column type %T", data)
    }

    quicktemplate.ReleaseWriter(templateWriter)
  }
  return nil
}
