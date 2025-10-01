package server

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"slices"
	"strings"
	"sync"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	degradation_detector "github.com/JetBrains/ij-perf-report-aggregator/pkg/degradation-detector"
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
	awsIndex := strings.LastIndex(s, "aws")
	if awsIndex != -1 {
		return s[:awsIndex+3]
	}

	lastIndex := strings.LastIndex(s, "-")
	if lastIndex == -1 {
		return s
	}
	return s[:lastIndex]
}

func (t *StatsServer) CreateProcessMetricDataHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
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
		defer request.Body.Close()
		if err := decoder.Decode(&params); err != nil {
			http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
			return
		}

		// Replace all matches with %
		machine := removeLastPart(params.Machine) + "%"
		// Map product to table - this is a simplified mapping, adjust as needed
		table := mapProductToTable(params.Product)
		if table == "" {
			http.Error(w, "Unknown product: "+params.Product, http.StatusBadRequest)
			return
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
			http.Error(w, "Failed to open database: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer func(db driver.Conn) {
			_ = db.Close()
		}(db)

		var queryResult struct {
			MetricValues []int
		}

		err = db.QueryRow(request.Context(), sql).Scan(&queryResult.MetricValues)
		if err != nil {
			http.Error(w, "Database query failed: "+err.Error(), http.StatusInternalServerError)
			return
		}

		if len(queryResult.MetricValues) == 0 {
			http.Error(w, "No data found for the specified parameters", http.StatusNotFound)
			return
		}

		// Run Change Point algorithm
		changePoints := statistic.GetChangePointIndexes(queryResult.MetricValues, 3)

		// Median difference and effect size thresholds (same as degradation detector)
		medianDifferenceThreshold := 10.0
		effectSizeThreshold := 2.0

		// Filter change points based on median difference and effect size
		validChangePoints := filterValidChangePoints(queryResult.MetricValues, changePoints, medianDifferenceThreshold, effectSizeThreshold)

		// Get the segment to analyze for max value
		// Strategy: use data after the last significant behavior change
		var segmentForAnalysis []int
		if len(validChangePoints) == 0 {
			// No valid change points - use all data
			segmentForAnalysis = queryResult.MetricValues
		} else {
			// Use data after the last valid change point
			lastChangePoint := validChangePoints[len(validChangePoints)-1]
			segmentForAnalysis = queryResult.MetricValues[lastChangePoint:]

		}

		// Remove outliers using MAD-based detection (windowSize=5, threshold=3)
		processedData := outlier_detection.RemoveOutliers(segmentForAnalysis, 5, 3.0)

		type processMetricResponse struct {
			MaxValue int `json:"maxValue"`
		}

		response := processMetricResponse{
			MaxValue: slices.Max(processedData),
		}

		jsonData, err := json.Marshal(response)
		if err != nil {
			http.Error(w, "Failed to marshal response: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(jsonData)
	}
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

// filterValidChangePoints filters change points based on median difference and effect size thresholds
// This logic is similar to the degradation detector's approach
func filterValidChangePoints(values []int, changePoints []int, medianDifferenceThreshold float64, effectSizeThreshold float64) []int {
	if len(changePoints) == 0 {
		return changePoints
	}

	// Split values into segments
	segments := degradation_detector.GetSegmentsBetweenChangePoints(changePoints, values)
	if len(segments) < 2 {
		return []int{}
	}

	validChangePoints := make([]int, 0)

	// Iterate through segments and validate change points
	for i := 1; i < len(segments); i++ {
		prevSegment := segments[i-1]
		currentSegment := segments[i]

		// Calculate medians
		prevMedian := statistic.Median(prevSegment)
		currentMedian := statistic.Median(currentSegment)

		// Calculate percentage change and absolute change
		percentageChange := math.Abs((currentMedian - prevMedian) / prevMedian * 100)
		absoluteChange := math.Abs(currentMedian - prevMedian)

		// Check if change is significant enough
		if absoluteChange < 10 || percentageChange < medianDifferenceThreshold {
			continue
		}

		// Calculate effect size
		effectSize := statistic.EffectSize(currentSegment, prevSegment)
		if effectSize < effectSizeThreshold {
			continue
		}

		// This is a valid change point
		validChangePoints = append(validChangePoints, changePoints[i-1])
	}

	return validChangePoints
}
