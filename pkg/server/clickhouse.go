package server

import (
	"encoding/json"
	"fmt"
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/JetBrains/ij-perf-report-aggregator/pkg/degradation-detector/statistic"
	"github.com/JetBrains/ij-perf-report-aggregator/pkg/util"
	"github.com/valyala/bytebufferpool"
	"net/http"
	"slices"
	"strings"
	"sync"
)

func (t *StatsServer) openDatabaseConnection() (driver.Conn, error) {
	return clickhouse.Open(&clickhouse.Options{
		Addr: []string{t.dbUrl},
		Auth: clickhouse.Auth{
			Database: "ij",
		},
		Settings: map[string]interface{}{
			"readonly":         1,
			"max_query_size":   1000000,
			"max_memory_usage": 3221225472,
		},
	})
}

func toJSONBuffer(data interface{}) (*bytebufferpool.ByteBuffer, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	buffer := bytebufferpool.Get()
	_, err = buffer.Write(jsonData)
	if err != nil {
		return nil, err
	}

	return buffer, nil
}

type responseItem struct {
	Project     string
	MeasureName string
	Median      float64
}

func (t *StatsServer) getBranchComparison(request *http.Request) (*bytebufferpool.ByteBuffer, bool, error) {

	type requestParams struct {
		Table        string   `json:"table"`
		MeasureNames []string `json:"measure_names"`
		Branch       string   `json:"branch"`
		Machine      string   `json:"machine"`
	}

	var params requestParams
	data, err := util.DecodeQuery(request.URL.Path[len("/api/compareBranches/"):])
	if err != nil {
		return nil, false, err
	}
	err = json.Unmarshal(data, &params)
	if err != nil {
		return nil, false, err
	}

	quotedMeasureNames := make([]string, len(params.MeasureNames))
	for i, name := range params.MeasureNames {
		quotedMeasureNames[i] = "'" + name + "'"
	}
	measureNamesString := strings.Join(quotedMeasureNames, ",")

	sql := fmt.Sprintf("SELECT project as Project, measure_name as MeasureName, arraySlice(groupArray(measure_value), 1, 50) AS MeasureValues FROM (SELECT project, measures.name as measure_name, measures.value as measure_value FROM %s ARRAY JOIN measures WHERE branch = '%s' AND measure_name in (%s) AND machine like '%s' ORDER BY generated_time DESC)GROUP BY project, measure_name;", params.Table, params.Branch, measureNamesString, params.Machine)
	db, err := t.openDatabaseConnection()
	defer func(db driver.Conn) {
		_ = db.Close()
	}(db)
	if err != nil {
		return nil, false, err
	}

	var queryResults []struct {
		Project       string
		MeasureName   string
		MeasureValues []int
	}

	err = db.Select(request.Context(), &queryResults, sql)
	if err != nil {
		return nil, false, err
	}

	response := getMedianValues(queryResults)
	buffer, err := toJSONBuffer(response)
	return buffer, true, err
}

func getMedianValues(queryResults []struct {
	Project       string
	MeasureName   string
	MeasureValues []int
}) []responseItem {

	responseChan := make(chan responseItem, len(queryResults))
	var wg sync.WaitGroup
	for _, result := range queryResults {
		wg.Add(1)
		go func() {
			defer wg.Done()
			values := result.MeasureValues
			slices.Reverse(values)
			indexes := statistic.GetChangePointIndexes(values, 1)
			var valuesAfterLastChangePoint []int
			if len(indexes) == 0 {
				valuesAfterLastChangePoint = values
			} else {
				lastIndex := indexes[len(indexes)-1]
				valuesAfterLastChangePoint = values[lastIndex:]
			}
			median := statistic.Median(valuesAfterLastChangePoint)

			responseChan <- responseItem{
				Project:     result.Project,
				MeasureName: result.MeasureName,
				Median:      median,
			}
		}()
	}

	go func() {
		wg.Wait()
		close(responseChan)
	}()

	response := make([]responseItem, 0, len(queryResults))
	for item := range responseChan {
		response = append(response, item)
	}
	return response
}

func (t *StatsServer) getDistinctHighlightingPasses(request *http.Request) (*bytebufferpool.ByteBuffer, bool, error) {
	db, err := t.openDatabaseConnection()
	if err != nil {
		return nil, false, err
	}
	defer func(db driver.Conn) {
		_ = db.Close()
	}(db)

	var queryResult []struct {
		PassName string
	}

	sql := "SELECT DISTINCT arrayJoin((arrayFilter(x-> x LIKE 'highlighting/%', `metrics.name`))) as PassName from report where generated_time >subtractMonths(now(),12)"
	err = db.Select(request.Context(), &queryResult, sql)
	if err != nil {
		return nil, false, err
	}

	passes := make([]string, len(queryResult))
	for i, v := range queryResult {
		passes[i] = v.PassName
	}

	buffer, err := toJSONBuffer(passes)
	return buffer, true, err
}
