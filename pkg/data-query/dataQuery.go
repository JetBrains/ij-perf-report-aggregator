package data_query

import (
	"context"
	"errors"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/ClickHouse/ch-go/proto"
	sqlutil "github.com/JetBrains/ij-perf-report-aggregator/pkg/sql-util"
	"github.com/JetBrains/ij-perf-report-aggregator/pkg/util"
	"github.com/sakura-internet/go-rison/v4"
	"github.com/valyala/bytebufferpool"
	"github.com/valyala/quicktemplate"
)

type Query struct {
	Database string           `json:"db"`
	Table    string           `json:"table"`
	Flat     bool             `json:"flat"`
	Fields   []QueryDimension `json:"fields,omitempty"`
	Filters  []QueryFilter    `json:"filters,omitempty"`
	Order    []string         `json:"order,omitempty"`

	Aggregator          string           `json:"-"`
	Dimensions          []QueryDimension `json:"-"`
	TimeDimensionFormat string           `json:"-"`
}

type QueryFilter struct {
	Field    string      `json:"f"`
	Value    interface{} `json:"v,omitempty"`
	Sql      string      `json:"q,omitempty"`
	Operator string      `json:"o,omitempty"`
	Split    bool        `json:"s"`
}

type QueryDimension struct {
	Name    string `json:"n"`
	Sql     string `json:"sql"`
	SubName string `json:"subName,omitempty"`

	metricPath      string
	metricName      string
	metricValueName rune

	resultPropertyName string

	arrayJoin string
}

func ReadQueryV2(request *http.Request) ([]Query, bool, error) {
	decompressed, err := util.DecodeQuery(request.URL.Path[len("/api/q/"):])
	if err != nil {
		return nil, false, fmt.Errorf("cannot decode query: %w", err)
	}

	if len(decompressed) == 0 {
		rawPath := request.URL.RawPath
		return nil, false, errors.New("query not found: " + rawPath)
	}

	wrappedAsArray := decompressed[0] == '['
	parser := queryParsers.Get()
	defer queryParsers.Put(parser)

	// fileName := strconv.FormatUint(xxh3.HashString(request.URL.Path), 36) + ".json"
	// _ = os.WriteFile("/Volumes/data/queries/"+fileName, decompressed, 0644)

	list, err := readQuery(decompressed)
	if err != nil {
		return nil, false, err
	}
	return list, wrappedAsArray, nil
}

func getSplitParameters(query Query) (*SplitParameters, error) {
	splitParameters := SplitParameters{
		numberOfSplits: 1,
	}

	if len(query.Filters) == 0 {
		// Handle case when there are no filters.
		return &splitParameters, nil
	}

	for _, filter := range query.Filters {
		if !filter.Split {
			continue
		}

		valueSlice, err := assertValueSlice(filter.Value)
		if err != nil {
			return nil, err
		}

		values, err := convertValuesToMap(valueSlice)
		if err != nil {
			return nil, err
		}

		splitParameters.numberOfSplits = len(values)
		splitParameters.splitField = filter.Field
		splitParameters.values = values
	}

	return &splitParameters, nil
}

// Helper function to assert filter.Value to a slice of empty interfaces.
func assertValueSlice(value interface{}) ([]interface{}, error) {
	valueSlice, ok := value.([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid filter.Value type %T, expected array", value)
	}
	return valueSlice, nil
}

// Helper function to convert values in the slice to a map with string keys.
func convertValuesToMap(valueSlice []interface{}) (map[string]int, error) {
	values := make(map[string]int)
	for i, value := range valueSlice {
		strValue, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf("invalid filter.Value type %T, expected string", value)
		}
		values[strValue] = i
	}
	return values, nil
}

func ReadQuery(request *http.Request) ([]Query, bool, error) {
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
		return nil, false, fmt.Errorf("cannot decode query: %w", err)
	}

	list, err := readQuery(jsonData)
	if err != nil {
		return nil, false, err
	}
	return list, wrappedAsArray, nil
}

func SelectRows(ctx context.Context, query Query, table string, dbSupplier DatabaseConnectionSupplier, totalWriter *quicktemplate.QWriter) error {
	sqlQuery, columnNameToIndex, err := buildSql(query, table)
	if err != nil {
		return err
	}

	splitParameters, err := getSplitParameters(query)
	if err != nil {
		return nil
	}

	columnBuffers := make([][]*bytebufferpool.ByteBuffer, splitParameters.numberOfSplits)
	err = executeQuery(ctx, sqlQuery, query, dbSupplier, func(_ context.Context, block proto.Block, result *proto.Results) error {
		if block.Rows == 0 {
			return nil
		}
		return writeResult(result, columnNameToIndex, columnBuffers, query, splitParameters)
	})
	if err != nil {
		return err
	}

	writeBuffers(columnBuffers, totalWriter, query.Flat)
	return nil
}

func writeBuffers(columnBuffers [][]*bytebufferpool.ByteBuffer, totalWriter *quicktemplate.QWriter, isFlat bool) {
	if !isFlat && len(columnBuffers) == 1 {
		totalWriter.S("[")
	}
	for splitNumber, splitColumnBuffers := range columnBuffers {
		if splitNumber != 0 {
			totalWriter.S(",")
		}
		if len(columnBuffers) > 1 {
			totalWriter.S("[")
		}
		for columnIndex, buffer := range splitColumnBuffers {
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
		if len(columnBuffers) > 1 {
			totalWriter.S("]")
		}
	}
	if !isFlat && len(columnBuffers) == 1 {
		totalWriter.S("]")
	}
}

//gocyclo:ignore
func buildSql(query Query, table string) (string, map[string]int, error) {
	var sb strings.Builder

	sb.WriteString("select")

	// the only array join is supported for now
	arrayJoin := ""
	for _, dimension := range query.Fields {
		if dimension.arrayJoin != "" {
			// the only array join is supported for now
			arrayJoin = dimension.arrayJoin
			// for field add distinct to filter duplicates out
			// sb.WriteString(" distinct ")
			break
		}
	}

	columnNameToIndex := make(map[string]int, len(query.Dimensions)+len(query.Fields))
	columnIndex := 0

	dimensionWritten := false
	for _, dimension := range query.Dimensions {
		// check that the field with the same name doesn't exists
		fieldExist := false
		for _, field := range query.Fields {
			if field.Name == dimension.Name {
				fieldExist = true
				break
			}
		}
		if !fieldExist {
			if !dimensionWritten {
				sb.WriteRune(' ')
			} else {
				sb.WriteRune(',')
			}
			columnNameToIndex[dimension.Name] = columnIndex
			columnIndex++

			writeDimension(dimension, &sb)
			dimensionWritten = true
		}
	}

	// write extra fields to the end, so, it maybe skipped during serialization
	for i, field := range query.Fields {
		if i != 0 || dimensionWritten {
			sb.WriteRune(',')
		}
		sb.WriteRune(' ')

		if field.Sql != "" {
			columnNameToIndex[field.Name] = columnIndex
			columnIndex++
			writeDimension(field, &sb)
			continue
		}

		effectiveColumnName := ""

		if query.Aggregator != "" {
			sb.WriteString(query.Aggregator)
			sb.WriteRune('(')
		}

		if field.metricPath == "" {
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

		if query.Aggregator != "" {
			sb.WriteRune(')')
		}

		if field.resultPropertyName != "" {
			sb.WriteString(" as ")
			sb.WriteString(field.resultPropertyName)
			effectiveColumnName = field.resultPropertyName
		} else if query.Aggregator != "" {
			sb.WriteString(" as ")
			if field.arrayJoin == "" {
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

	if arrayJoin == "" {
		for _, dimension := range query.Dimensions {
			if dimension.arrayJoin != "" {
				arrayJoin = dimension.arrayJoin
				break
			}
		}
	}

	if arrayJoin != "" {
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

func writeExtractJsonObject(sb *strings.Builder, field QueryDimension) {
	sb.WriteString("arrayFirst(it -> JSONExtractString(it, 'n') = '")
	sb.WriteString(field.metricName)
	sb.WriteString("', JSONExtractArrayRaw(raw_report, '")
	sb.WriteString(field.metricPath)
	sb.WriteString("'))")
}

func writeDimension(dimension QueryDimension, sb *strings.Builder) {
	if dimension.Sql == "" {
		sb.WriteString(dimension.Name)
	} else {
		sb.WriteString(dimension.Sql)
		sb.WriteString(" as ")
		// escape - maybe nested name with dot
		sb.WriteByte('`')
		sb.WriteString(dimension.Name)
		sb.WriteByte('`')
	}
}

func writeWhereClause(sb *strings.Builder, query Query) error {
	sb.WriteString(" where")
	for i, filter := range query.Filters {
		if i != 0 {
			sb.WriteString(" and")
		}
		sb.WriteString(" (")
		sb.WriteString(filter.Field)

		if filter.Sql != "" {
			if filter.Operator != "" {
				return errors.New("sql and operator are mutually exclusive")
			}
			if filter.Value != nil {
				return errors.New("sql and value are mutually exclusive")
			}

			sb.WriteByte(' ')
			sb.WriteString(filter.Sql)
			sb.WriteByte(')')
			continue
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
		case bool:
			sb.WriteString(filter.Operator)
			sb.WriteString(strconv.FormatBool(v))
		case string:
			sb.WriteString(" ")
			sb.WriteString(filter.Operator)
			sb.WriteString(" ")
			writeString(sb, v)
		case []string:
			sb.WriteString(" in (")
			for j := range v {
				if j != 0 {
					sb.WriteByte(',')
				}
				writeString(sb, v[j])
			}
			sb.WriteByte(')')
		case []interface{}:
			sb.WriteString(" in (")
			for j := range v {
				if j != 0 {
					sb.WriteByte(',')
				}
				switch e := v[j].(type) {
				case string:
					writeString(sb, e)
				case bool:
					sb.WriteString(strconv.FormatBool(e))
				default:
					return fmt.Errorf("filter value type [%T] is not supported", v[j])
				}
			}
			sb.WriteByte(')')
		default:
			return fmt.Errorf("filter value type %T is not supported", v)
		}
		sb.WriteByte(')')
	}
	return nil
}

func writeString(sb *strings.Builder, s string) {
	sb.WriteByte('\'')
	_, _ = sqlutil.StringEscaper.WriteString(sb, s)
	sb.WriteByte('\'')
}
