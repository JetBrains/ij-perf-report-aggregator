package data_query

import (
	"fmt"
	"math"
	"strconv"

	"github.com/ClickHouse/ch-go/proto"
	"github.com/valyala/bytebufferpool"
	"github.com/valyala/quicktemplate"
)

// separate byte buffer pool - different sizes
var byteBufferPool bytebufferpool.Pool

type SplitParameters struct {
	numberOfSplits int
	splitField     string
	values         map[string]int
}

func writeResult(result *proto.Results, columnNameToIndex map[string]int, columnBuffers [][]*bytebufferpool.ByteBuffer, query Query, splitParameters *SplitParameters) error {
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
			case *proto.ColStr:
				for i, position := range data.Pos {
					rowToSplitIndex[i] = splitParameters.values[string(data.Buf[position.Start:position.End])]
				}
			}
		}
	}
	for i := range splitParameters.numberOfSplits {
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
				buffer := getBufferAndWriteSplittingComma(columnBuffers, rowToSplitIndex, i, columnIndex)
				buffer.B = strconv.AppendInt(buffer.B, int64(v), 10)
			}
		case *proto.ColInt64:
			for i, v := range *data {
				buffer := getBufferAndWriteSplittingComma(columnBuffers, rowToSplitIndex, i, columnIndex)
				buffer.B = strconv.AppendInt(buffer.B, v, 10)
			}
		case *proto.ColUInt8:
			for i, v := range *data {
				buffer := getBufferAndWriteSplittingComma(columnBuffers, rowToSplitIndex, i, columnIndex)
				buffer.B = strconv.AppendUint(buffer.B, uint64(v), 10)
			}

		case *proto.ColBool:
			for i, v := range *data {
				buffer := getBufferAndWriteSplittingComma(columnBuffers, rowToSplitIndex, i, columnIndex)
				buffer.B = strconv.AppendBool(buffer.B, v)
			}

		case *proto.ColUInt16:
			for i, v := range *data {
				buffer := getBufferAndWriteSplittingComma(columnBuffers, rowToSplitIndex, i, columnIndex)
				buffer.B = strconv.AppendUint(buffer.B, uint64(v), 10)
			}
		case *proto.ColUInt32:
			for i, v := range *data {
				buffer := getBufferAndWriteSplittingComma(columnBuffers, rowToSplitIndex, i, columnIndex)
				buffer.B = strconv.AppendUint(buffer.B, uint64(v), 10)
			}
		case *proto.ColUInt64:
			for i, v := range *data {
				buffer := getBufferAndWriteSplittingComma(columnBuffers, rowToSplitIndex, i, columnIndex)
				buffer.B = strconv.AppendUint(buffer.B, v, 10)
			}

		case *proto.ColFloat32:
			for i, v := range *data {
				buffer := getBufferAndWriteSplittingComma(columnBuffers, rowToSplitIndex, i, columnIndex)
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
				buffer := getBufferAndWriteSplittingComma(columnBuffers, rowToSplitIndex, i, columnIndex)
				n := int64(v)
				if float64(n) == v {
					// fast path - just int
					buffer.B = strconv.AppendInt(buffer.B, n, 10)
				} else {
					if math.IsNaN(v) {
						buffer.B = strconv.AppendFloat(buffer.B, 0, 'f', -1, 64)
					} else {
						buffer.B = strconv.AppendFloat(buffer.B, v, 'f', -1, 64)
					}
				}
			}

		case *proto.ColStr:
			templateWriters, jsonWriters := getWriters(splitParameters, columnBuffers, columnIndex)
			for i, position := range data.Pos {
				getBufferAndWriteSplittingComma(columnBuffers, rowToSplitIndex, i, columnIndex)
				jsonWriter := jsonWriters[rowToSplitIndex[i]]
				jsonWriter.Q(string(data.Buf[position.Start:position.End]))
			}
			releaseWriters(splitParameters, templateWriters)

		case *proto.ColDateTime:
			templateWriters, jsonWriters := getWriters(splitParameters, columnBuffers, columnIndex)
			for i, v := range data.Data {
				getBufferAndWriteSplittingComma(columnBuffers, rowToSplitIndex, i, columnIndex)
				jsonWriter := jsonWriters[rowToSplitIndex[i]]
				jsonWriter.Q(v.Time().Format(query.TimeDimensionFormat))
			}
			releaseWriters(splitParameters, templateWriters)

		case *proto.ColLowCardinality[string]:
			templateWriters, jsonWriters := getWriters(splitParameters, columnBuffers, columnIndex)
			for i, t := range data.Values {
				getBufferAndWriteSplittingComma(columnBuffers, rowToSplitIndex, i, columnIndex)
				jsonWriter := jsonWriters[rowToSplitIndex[i]]
				jsonWriter.Q(t)
			}
			releaseWriters(splitParameters, templateWriters)
		default:
			return fmt.Errorf("unsupported column type %T", data)
		}
	}
	return nil
}

func getBufferAndWriteSplittingComma(columnBuffers [][]*bytebufferpool.ByteBuffer, rowToSplitIndex []int, i int, columnIndex int) *bytebufferpool.ByteBuffer {
	buffer := columnBuffers[rowToSplitIndex[i]][columnIndex]
	if buffer.Len() != 0 {
		_ = buffer.WriteByte(',')
	}
	return buffer
}

func releaseWriters(splitParameters *SplitParameters, templateWriters map[int]*quicktemplate.Writer) {
	for i := range splitParameters.numberOfSplits {
		quicktemplate.ReleaseWriter(templateWriters[i])
	}
}

func getWriters(splitParameters *SplitParameters, columnBuffers [][]*bytebufferpool.ByteBuffer, columnIndex int) (map[int]*quicktemplate.Writer, map[int]*quicktemplate.QWriter) {
	templateWriters := make(map[int]*quicktemplate.Writer)
	jsonWriters := make(map[int]*quicktemplate.QWriter)
	for i := range splitParameters.numberOfSplits {
		templateWriters[i] = quicktemplate.AcquireWriter(columnBuffers[i][columnIndex])
		jsonWriters[i] = templateWriters[i].N()
	}
	return templateWriters, jsonWriters
}
