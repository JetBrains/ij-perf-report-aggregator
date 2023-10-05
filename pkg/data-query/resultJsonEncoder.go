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

type SplitParameters struct {
  numberOfSplits int
  splitField     string
  values         map[string]int
}

//nolint:gocyclo
func writeResult(result *proto.Results, columnNameToIndex map[string]int, columnBuffers [][]*bytebufferpool.ByteBuffer, query DataQuery, splitParameters *SplitParameters) error {
  var rowToSplitIndex []int

  for _, column := range *result {
    rowToSplitIndex = make([]int, column.Data.Rows())
  }

  for _, column := range *result {
    if column.Name == splitParameters.splitField {
      //nolint:gocritic
      switch data := column.Data.(type) {
      case *proto.ColLowCardinality[string]:
        for i, t := range data.Values {
          rowToSplitIndex[i] = splitParameters.values[t]
        }
      }
    }
  }
  for i := 0; i < splitParameters.numberOfSplits; i++ {
    if columnBuffers[i] == nil {
      columnBuffers[i] = make([]*bytebufferpool.ByteBuffer, len(*result))
    }
    for j := range *result {
      columnBuffer := columnBuffers[i][j]
      if columnBuffer == nil {
        columnBuffers[i][j] = byteBufferPool.Get()
      }
    }
  }

  for _, column := range *result {

    columnIndex := columnNameToIndex[column.Name]

    switch data := column.Data.(type) {
    case *proto.ColInt32:
      for i, v := range *data {
        buffer := columnBuffers[rowToSplitIndex[i]][columnIndex]
        if buffer.Len() != 0 {
          _ = buffer.WriteByte(',')
        }
        buffer.B = strconv.AppendInt(buffer.B, int64(v), 10)
      }
    case *proto.ColInt64:
      for i, v := range *data {
        buffer := columnBuffers[rowToSplitIndex[i]][columnIndex]
        if buffer.Len() != 0 {
          _ = buffer.WriteByte(',')
        }
        buffer.B = strconv.AppendInt(buffer.B, v, 10)
      }
    case *proto.ColUInt8:
      for i, v := range *data {
        buffer := columnBuffers[rowToSplitIndex[i]][columnIndex]
        if buffer.Len() != 0 {
          _ = buffer.WriteByte(',')
        }
        buffer.B = strconv.AppendUint(buffer.B, uint64(v), 10)
      }

    case *proto.ColBool:
      for i, v := range *data {
        buffer := columnBuffers[rowToSplitIndex[i]][columnIndex]
        if buffer.Len() != 0 {
          _ = buffer.WriteByte(',')
        }
        buffer.B = strconv.AppendBool(buffer.B, v)
      }

    case *proto.ColUInt16:
      for i, v := range *data {
        buffer := columnBuffers[rowToSplitIndex[i]][columnIndex]
        if buffer.Len() != 0 {
          _ = buffer.WriteByte(',')
        }
        buffer.B = strconv.AppendUint(buffer.B, uint64(v), 10)
      }
    case *proto.ColUInt32:
      for i, v := range *data {
        buffer := columnBuffers[rowToSplitIndex[i]][columnIndex]
        if buffer.Len() != 0 {
          _ = buffer.WriteByte(',')
        }
        buffer.B = strconv.AppendUint(buffer.B, uint64(v), 10)
      }
    case *proto.ColUInt64:
      for i, v := range *data {
        buffer := columnBuffers[rowToSplitIndex[i]][columnIndex]
        if buffer.Len() != 0 {
          _ = buffer.WriteByte(',')
        }
        buffer.B = strconv.AppendUint(buffer.B, v, 10)
      }

    case *proto.ColFloat32:
      for i, v := range *data {
        buffer := columnBuffers[rowToSplitIndex[i]][columnIndex]
        if buffer.Len() != 0 {
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
        buffer := columnBuffers[rowToSplitIndex[i]][columnIndex]
        if buffer.Len() != 0 {
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
      templateWriters := make(map[int]*quicktemplate.Writer)
      jsonWriters := make(map[int]*quicktemplate.QWriter)
      for i := 0; i < splitParameters.numberOfSplits; i++ {
        templateWriters[i] = quicktemplate.AcquireWriter(columnBuffers[i][columnIndex])
        jsonWriters[i] = templateWriters[i].N()
      }
      for i, position := range data.Pos {
        buffer := columnBuffers[rowToSplitIndex[i]][columnNameToIndex[column.Name]]
        jsonWriter := jsonWriters[rowToSplitIndex[i]]
        if buffer.Len() != 0 {
          _ = buffer.WriteByte(',')
        }
        jsonWriter.Q(string(data.Buf[position.Start:position.End]))
      }
      for i := 0; i < splitParameters.numberOfSplits; i++ {
        quicktemplate.ReleaseWriter(templateWriters[i])
      }

    case *proto.ColDateTime:
      templateWriters := make(map[int]*quicktemplate.Writer)
      jsonWriters := make(map[int]*quicktemplate.QWriter)
      for i := 0; i < splitParameters.numberOfSplits; i++ {
        templateWriters[i] = quicktemplate.AcquireWriter(columnBuffers[i][columnIndex])
        jsonWriters[i] = templateWriters[i].N()
      }
      for i, v := range data.Data {
        buffer := columnBuffers[rowToSplitIndex[i]][columnNameToIndex[column.Name]]
        jsonWriter := jsonWriters[rowToSplitIndex[i]]
        if buffer.Len() != 0 {
          _ = buffer.WriteByte(',')
        }
        jsonWriter.Q(v.Time().Format(query.TimeDimensionFormat))
      }
      for i := 0; i < splitParameters.numberOfSplits; i++ {
        quicktemplate.ReleaseWriter(templateWriters[i])
      }

    case *proto.ColLowCardinality[string]:
      templateWriters := make(map[int]*quicktemplate.Writer)
      jsonWriters := make(map[int]*quicktemplate.QWriter)
      for i := 0; i < splitParameters.numberOfSplits; i++ {
        templateWriters[i] = quicktemplate.AcquireWriter(columnBuffers[i][columnIndex])
        jsonWriters[i] = templateWriters[i].N()
      }
      for i, t := range data.Values {
        buffer := columnBuffers[rowToSplitIndex[i]][columnNameToIndex[column.Name]]
        jsonWriter := jsonWriters[rowToSplitIndex[i]]
        if buffer.Len() != 0 {
          _ = buffer.WriteByte(',')
        }
        jsonWriter.Q(t)
      }
      for i := 0; i < splitParameters.numberOfSplits; i++ {
        quicktemplate.ReleaseWriter(templateWriters[i])
      }
    default:
      return errors.Errorf("unsupported column type %T", data)
    }

  }
  return nil
}
