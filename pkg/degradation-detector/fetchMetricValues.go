package degradation_detector

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"

	dataQuery "github.com/JetBrains/ij-perf-report-aggregator/pkg/data-query"
	"github.com/JetBrains/ij-perf-report-aggregator/pkg/util"
	"github.com/alitto/pond"
)

type queryResult struct {
	timestamps []int64
	values     []int
	builds     []string
	buildTypes []string
}

type QueryResultWithSettings struct {
	queryResult
	Settings
}

type queryProducer interface {
	query() dataQuery.Query
}

func FetchMetricsFromClickhouse(settings []Settings, client *http.Client, backendUrl string) chan QueryResultWithSettings {
	dataChan := make(chan QueryResultWithSettings, 5)
	go func() {
		defer close(dataChan)
		var wg sync.WaitGroup
		pool := pond.New(5, 1000)
		for _, setting := range settings {
			wg.Add(1)
			pool.Submit(func() {
				defer wg.Done()
				ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
				defer cancel()
				data, err := getDataFromClickhouse(ctx, client, backendUrl, setting.query())
				if err != nil {
					slog.Error("error while getting queryResult from clickhouse", "error", err, "settings", setting)
					return
				}
				slog.Debug("fetched from clickhouse", "settings", setting)
				dataChan <- QueryResultWithSettings{
					queryResult: data,
					Settings:    setting,
				}
			})
		}
		wg.Wait()
	}()
	return dataChan
}

func getDataFromClickhouse(ctx context.Context, client *http.Client, backendUrl string, query dataQuery.Query) (queryResult, error) {
	response, err := getValuesFromServer(ctx, client, backendUrl, query)
	if err != nil {
		return queryResult{}, err
	}
	data, err := extractDataFromRequest(response)
	return data, err
}

func getValuesFromServer(ctx context.Context, client *http.Client, backendURL string, query dataQuery.Query) ([]byte, error) {
	url := backendURL + "/api/q/"
	queries := []dataQuery.Query{query}
	jsonQuery, err := json.Marshal(queries)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal query: %w", err)
	}

	encoded, err := util.EncodeQuery(jsonQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to encode query: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url+encoded, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send GET request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get data: %v", resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	return body, err
}

func (s FleetStartupSettings) query() dataQuery.Query {
	fields := []dataQuery.QueryDimension{
		{Name: "t", Sql: "toUnixTimestamp(generated_time)*1000"},
	}
	filters := []dataQuery.QueryFilter{
		{Field: "branch", Value: s.Branch},
		{Field: "generated_time", Sql: ">subtractDays(now(),100)"},
		{Field: "project", Value: "fleet"},
		{Field: "machine", Value: s.Machine, Operator: "like"},
		{Field: "triggeredBy", Value: ""},
	}
	if strings.HasSuffix(s.Metric, ".end") {
		metricName, _ := strings.CutSuffix(s.Metric, ".end")
		filters = append(filters, dataQuery.QueryFilter{Field: "measures.name", Value: metricName})
		fields = append(fields, dataQuery.QueryDimension{Name: "measures", SubName: "end", Sql: "(measures.start+measures.value)"})
	}
	fields = append(fields, dataQuery.QueryDimension{Name: "Build", Sql: "concat(toString(build_c1),'.',toString(build_c2),'.',toString(build_c3))"}, dataQuery.QueryDimension{Name: "tc_build_type"})

	query := dataQuery.Query{
		Database: "fleet",
		Table:    "report",
		Fields:   fields,
		Filters:  filters,
		Order:    []string{"t"},
	}
	return query
}

func (s StartupSettings) query() dataQuery.Query {
	fields := []dataQuery.QueryDimension{
		{Name: "t", Sql: "toUnixTimestamp(generated_time)*1000"},
	}
	filters := []dataQuery.QueryFilter{
		{Field: "branch", Value: s.Branch},
		{Field: "generated_time", Sql: ">subtractDays(now(),100)"},
		{Field: "project", Value: s.Project},
		{Field: "product", Value: s.Product},
		{Field: "machine", Value: s.Machine, Operator: "like"},
		{Field: "triggeredBy", Value: ""},
	}
	if strings.Contains(s.Metric, "/") {
		filters = append(filters, dataQuery.QueryFilter{Field: "metrics.name", Value: s.Metric})
		fields = append(fields, dataQuery.QueryDimension{Name: "metrics", SubName: "value"})
	}
	if strings.HasSuffix(s.Metric, ".end") {
		metricName, _ := strings.CutSuffix(s.Metric, ".end")
		filters = append(filters, dataQuery.QueryFilter{Field: "measure.name", Value: metricName})
		fields = append(fields, dataQuery.QueryDimension{Name: "measure", SubName: "end", Sql: "(measure.start+measure.duration)"})
	}
	if !strings.HasSuffix(s.Metric, ".end") && !strings.Contains(s.Metric, "/") {
		fields = append(fields, dataQuery.QueryDimension{Name: s.Metric})
	}
	fields = append(fields, dataQuery.QueryDimension{Name: "Build", Sql: "concat(toString(build_c1),'.',toString(build_c2))"}, dataQuery.QueryDimension{Name: "tc_build_type"})

	query := dataQuery.Query{
		Database: "ijDev",
		Table:    "report",
		Fields:   fields,
		Filters:  filters,
		Order:    []string{"t"},
	}
	return query
}

func (s PerformanceSettings) query() dataQuery.Query {
	fields := []dataQuery.QueryDimension{
		{Name: "t", Sql: "toUnixTimestamp(generated_time)*1000"},
		{Name: "measures", SubName: "value"},
	}
	if s.Db == "perfint" {
		fields = append(fields, dataQuery.QueryDimension{Name: "Build", Sql: "concat(toString(build_c1),'.',toString(build_c2))"})
	} else {
		fields = append(fields, dataQuery.QueryDimension{Name: "tc_build_id"})
	}
	fields = append(fields, dataQuery.QueryDimension{Name: "tc_build_type"})

	query := dataQuery.Query{
		Database: s.Db,
		Table:    s.Table,
		Fields:   fields,
		Filters: []dataQuery.QueryFilter{
			{Field: "branch", Value: s.Branch},
			{Field: "mode", Value: s.Mode},
			{Field: "generated_time", Sql: ">subtractDays(now(),100)"},
			{Field: "project", Value: s.Project},
			{Field: "measures.name", Value: s.Metric},
			{Field: "machine", Value: s.Machine, Operator: "like"},
			{Field: "triggeredBy", Value: ""},
		},
		Order: []string{"t"},
	}
	return query
}

func extractDataFromRequest(response []byte) (queryResult, error) {
	var data [][][]interface{}

	err := json.Unmarshal(response, &data)
	if err != nil {
		return queryResult{}, fmt.Errorf("failed to decode JSON: %w", err)
	}
	if len(data) == 0 {
		return queryResult{}, errors.New("no data")
	}
	if len(data[0]) < 3 {
		return queryResult{}, errors.New("not enough data")
	}
	timestamps, err := SliceToSliceInt64(data[0][0])
	if err != nil {
		return queryResult{}, fmt.Errorf("failed to convert values: %w", err)
	}
	values, err := SliceToSliceOfInt(data[0][1])
	if err != nil {
		return queryResult{}, fmt.Errorf("failed to convert values: %w", err)
	}
	builds, err := SliceToSliceOfString(data[0][2])
	if err != nil {
		return queryResult{}, fmt.Errorf("failed to convert values: %w", err)
	}
	buildTypes, err := SliceToSliceOfString(data[0][3])
	if err != nil {
		return queryResult{}, fmt.Errorf("failed to convert values: %w", err)
	}
	return queryResult{
		timestamps: timestamps,
		values:     values,
		builds:     builds,
		buildTypes: buildTypes,
	}, err
}
