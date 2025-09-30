package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strings"
	"sync"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/JetBrains/ij-perf-report-aggregator/pkg/degradation-detector/statistic"
	"github.com/JetBrains/ij-perf-report-aggregator/pkg/outlier-detection"
	"github.com/JetBrains/ij-perf-report-aggregator/pkg/util"
	"github.com/valyala/bytebufferpool"
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
		Mode         string   `json:"mode"`
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

	mode := params.Mode
	if params.Mode == "default" {
		mode = ""
	}

	sql := fmt.Sprintf("SELECT project as Project, measure_name as MeasureName, arraySlice(groupArray(measure_value), 1, 50) AS MeasureValues FROM (SELECT project, measures.name as measure_name, measures.value as measure_value FROM %s ARRAY JOIN measures WHERE branch = '%s' AND measure_name in (%s) AND machine like '%s' and mode = '%s' ORDER BY generated_time DESC)GROUP BY project, measure_name;", params.Table, params.Branch, measureNamesString, params.Machine, mode)
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

func (t *StatsServer) getModeComparison(request *http.Request) (*bytebufferpool.ByteBuffer, bool, error) {
	type requestParams struct {
		Table        string   `json:"table"`
		MeasureNames []string `json:"measure_names"`
		Branch       string   `json:"branch"`
		Machine      string   `json:"machine"`
		Mode         string   `json:"mode"`
	}

	var params requestParams
	data, err := util.DecodeQuery(request.URL.Path[len("/api/compareModes/"):])
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

	sql := fmt.Sprintf("SELECT project as Project, measure_name as MeasureName, arraySlice(groupArray(measure_value), 1, 50) AS MeasureValues FROM (SELECT project, measures.name as measure_name, measures.value as measure_value FROM %s ARRAY JOIN measures WHERE mode = '%s' AND branch = '%s' AND measure_name in (%s) AND machine like '%s' AND generated_time >subtractMonths(now(),1) ORDER BY generated_time DESC)GROUP BY project, measure_name;", params.Table, params.Mode, params.Branch, measureNamesString, params.Machine)
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
},
) []responseItem {
	responseChan := make(chan responseItem, len(queryResults))
	var wg sync.WaitGroup
	for _, result := range queryResults {
		wg.Go(func() {
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
		})
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

func removeLastPart(s string) string {
	lastIndex := strings.LastIndex(s, "-")
	if lastIndex == -1 {
		return s
	}
	return s[:lastIndex]
}

func (t *StatsServer) processMetricData(request *http.Request) (*bytebufferpool.ByteBuffer, bool, error) {
	type requestParams struct {
		TestName   string `json:"testName"`
		Branch     string `json:"branch"`
		Machine    string `json:"machine"`
		Product    string `json:"product"`
		MetricName string `json:"metricName"`
		Mode       string `json:"mode"`
	}

	var params requestParams
	decoder := json.NewDecoder(request.Body)
	if err := decoder.Decode(&params); err != nil {
		return nil, false, err
	}

	// Set default minimum points if not specified
	minPoints := 5

	// Map OS to machine field (assuming OS maps to machine in the database)

	// Replace all matches with %
	machine := removeLastPart(params.Machine) + "%"
	// Map product to table - this is a simplified mapping, adjust as needed
	table := mapProductToTable(params.Product)
	if table == "" {
		return nil, false, fmt.Errorf("unknown product: %s", params.Product)
	}

	// Query the database for metric values
	sql := fmt.Sprintf(`
		SELECT groupArray(metric_value) AS MetricValues
		FROM (
			SELECT measures.value as metric_value
			FROM perfintDev.%s ARRAY JOIN measures
			WHERE branch = '%s'
			AND measures.name = '%s'
			AND machine LIKE '%s'
			AND project = '%s'
			AND mode = '%s'
			AND generated_time >= now() - INTERVAL 1 MONTH
			ORDER BY generated_time
		)
	`, table, params.Branch, params.MetricName, machine, params.TestName, params.Mode)

	db, err := t.openDatabaseConnection()
	if err != nil {
		return nil, false, err
	}
	defer func(db driver.Conn) {
		_ = db.Close()
	}(db)

	var queryResult struct {
		MetricValues []int
	}

	err = db.QueryRow(request.Context(), sql).Scan(&queryResult.MetricValues)
	if err != nil {
		return nil, false, err
	}

	if len(queryResult.MetricValues) == 0 {
		return nil, false, errors.New("no data found for the specified parameters")
	}

	// Run Change Point algorithm
	changePoints := statistic.GetChangePointIndexes(queryResult.MetricValues, 3)

	// Get the last segment that's larger than minPoints
	var lastSegment []int
	if len(changePoints) == 0 {
		lastSegment = queryResult.MetricValues
	} else {
		// Find the last segment with sufficient points
		for i := len(changePoints) - 1; i >= 0; i-- {
			// Check if this is the last segment and it has enough points
			if i == len(changePoints)-1 {
				segment := queryResult.MetricValues[changePoints[i]:]
				if len(segment) >= minPoints {
					lastSegment = segment
					break
				}
			}
		}

		// If no suitable segment found from change points, use the entire last segment
		if len(lastSegment) == 0 {
			if len(changePoints) > 0 {
				lastIndex := changePoints[len(changePoints)-1]
				lastSegment = queryResult.MetricValues[lastIndex:]
			} else {
				lastSegment = queryResult.MetricValues
			}
		}
	}

	// If the segment is still too small, use all data
	if len(lastSegment) < minPoints {
		lastSegment = queryResult.MetricValues
	}

	// Remove outliers using MAD-based detection (windowSize=5, threshold=3)
	processedData := outlier_detection.RemoveOutliers(lastSegment, 5, 3.0)

	type processMetricResponse struct {
		MaxValue int `json:"max_value"`
	}

	response := processMetricResponse{
		MaxValue: slices.Max(processedData),
	}

	buffer, err := toJSONBuffer(response)
	return buffer, true, err
}

// mapProductToTable maps product names to database table names
func mapProductToTable(product string) string {
	switch product {
	case "IU":
		return "idea"
	case "GO":
		return "goland"
	case "RM":
		return "ruby"
	case "PS":
		return "phpstorm"
	case "PY":
		return "pycharm"
	case "WS":
		return "webstorm"
	default:
		return ""
	}
}
