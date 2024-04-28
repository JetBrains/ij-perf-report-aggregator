package server

import (
	"context"
	"encoding/json"
	"github.com/JetBrains/ij-perf-report-aggregator/pkg/data-query"
	"github.com/valyala/bytebufferpool"
	"github.com/valyala/quicktemplate"
	"net/http"
)

type queryResultType []struct {
	Time               *uint64 `ch:"t"`
	Machine            string  `ch:"machine"`
	BuildTime          string  `ch:"build_time"`
	GeneratedTime      *uint64 `ch:"generated_time"`
	Project            string  `ch:"project"`
	TcBuildId          *uint32 `ch:"tc_build_id"`
	TcInstallerBuildId *uint32 `ch:"tc_installer_build_id"`
	Branch             string  `ch:"branch"`
	TcBuildType        string  `ch:"tc_build_type"`
	MeasureName        string  `ch:"measures.name"`
	MeasureValue       *int32  `ch:"measures.value"`
	MeasureType        string  `ch:"measures.type"`
	BuildC1            *uint8  `ch:"build_c1"`
	BuildC2            *uint16 `ch:"build_c2"`
	BuildC3            *uint16 `ch:"build_c3"`
	TriggeredBy        string  `ch:"triggeredBy"`
}

func (t *StatsServer) handleLoadRequestV2(request *http.Request) (*bytebufferpool.ByteBuffer, bool, error) {
	dataQueries, _, err := data_query.ReadQueryV2(request)
	if err != nil {
		return nil, false, err
	}
	return t.load(request.Context(), dataQueries)
}

func (t *StatsServer) handleLoadRequest(request *http.Request) (*bytebufferpool.ByteBuffer, bool, error) {
	dataQueries, _, err := data_query.ReadQuery(request)
	if err != nil {
		return nil, false, err
	}

	return t.load(request.Context(), dataQueries)
}

func (t *StatsServer) load(context context.Context, dataQueries []data_query.Query) (*bytebufferpool.ByteBuffer, bool, error) {
	buffer := byteBufferPool.Get()
	isOk := false
	defer func() {
		if !isOk {
			byteBufferPool.Put(buffer)
		}
	}()

	var result []queryResultType

	for _, dataQuery := range dataQueries {
		var queryResult queryResultType

		err := t.computeMeasureResponse(context, dataQuery, &queryResult)
		if err != nil {
			return nil, false, err
		}

		result = append(result, queryResult)
	}

	bytes, _ := json.Marshal(result)
	_, err := buffer.Write(bytes)
	if err != nil {
		return nil, false, err
	}
	isOk = true
	return buffer, true, nil
}

func (t *StatsServer) computeMeasureResponse(ctx context.Context, query data_query.Query, queryResults *queryResultType) error {
	table := query.Table
	if table == "" {
		table = "report"
	}
	connection, err := t.openDatabaseConnection(query.Database)
	if err != nil {
		return err
	}
	defer connection.Close()
	sqlQuery, _, _ := data_query.BuildSql(query, table)
	err = connection.Select(ctx, queryResults, sqlQuery)
	if err != nil {
		return err
	}
	return nil
}

func (t *StatsServer) oldLoad(request *http.Request, dataQueries []data_query.Query, wrappedAsArray bool) (*bytebufferpool.ByteBuffer, bool, error) {
	buffer := byteBufferPool.Get()
	isOk := false
	defer func() {
		if !isOk {
			byteBufferPool.Put(buffer)
		}
	}()

	templateWriter := quicktemplate.AcquireWriter(buffer)
	defer quicktemplate.ReleaseWriter(templateWriter)
	jsonWriter := templateWriter.N()

	if len(dataQueries) > 1 || wrappedAsArray {
		jsonWriter.S("[")
	}

	for index, dataQuery := range dataQueries {
		if index != 0 {
			jsonWriter.S(",")
		}

		err := t.oldComputeMeasureResponse(request.Context(), dataQuery, jsonWriter)
		if err != nil {
			return nil, false, err
		}
	}

	if len(dataQueries) > 1 || wrappedAsArray {
		jsonWriter.S("]")
	}
	isOk = true
	if len(buffer.B) == 0 {
		jsonWriter.S("[]")
	}
	return buffer, true, nil
}

func (t *StatsServer) oldComputeMeasureResponse(ctx context.Context, query data_query.Query, jsonWriter *quicktemplate.QWriter) error {
	table := query.Table
	if table == "" {
		table = "report"
	}

	err := data_query.SelectRows(ctx, query, table, t, jsonWriter)
	if err != nil {
		return err
	}
	return nil
}
