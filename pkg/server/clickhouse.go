package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strings"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/JetBrains/ij-perf-report-aggregator/pkg/degradation-detector/statistic"
	"github.com/JetBrains/ij-perf-report-aggregator/pkg/machine"
	"github.com/JetBrains/ij-perf-report-aggregator/pkg/outlier-detection"
	"github.com/JetBrains/ij-perf-report-aggregator/pkg/sql-util"
	"github.com/JetBrains/ij-perf-report-aggregator/pkg/util"
	"github.com/valyala/bytebufferpool"
)

func (t *StatsServer) openDatabaseConnection() (driver.Conn, error) {
	return clickhouse.Open(&clickhouse.Options{
		Addr: []string{t.dbUrl},
		Auth: clickhouse.Auth{
			Database: "ij",
		},
		Settings: map[string]any{
			"readonly":         1,
			"max_query_size":   1000000,
			"max_memory_usage": 3221225472,
		},
	})
}

func toJSONBuffer(data any) (*bytebufferpool.ByteBuffer, error) {
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

// machineCondition renders the WHERE condition for a machine selection mixing hardware-class
// group names and raw agent names: groups match by the computed machine_group alias (ClickHouse
// allows aliases in WHERE), raw names by exact machine. An empty selection matches everything.
func machineCondition(selection []string) string {
	var groups, names []string
	for _, value := range selection {
		quoted := "'" + sql_util.StringEscaper.Replace(value) + "'"
		if machine.IsGroup(value) {
			groups = append(groups, quoted)
		} else {
			names = append(names, quoted)
		}
	}

	var conditions []string
	if len(groups) > 0 {
		conditions = append(conditions, "machine_group IN ("+strings.Join(groups, ", ")+")")
	}
	if len(names) > 0 {
		conditions = append(conditions, "machine IN ("+strings.Join(names, ", ")+")")
	}
	if len(conditions) == 0 {
		return "1"
	}
	return strings.Join(conditions, " OR ")
}

func (t *StatsServer) getBranchComparison(request *http.Request) (*bytebufferpool.ByteBuffer, bool, error) {
	type requestParams struct {
		Table        string   `json:"table"`
		MeasureNames []string `json:"measure_names"`
		Branch1      string   `json:"branch1"`
		Branch2      string   `json:"branch2"`
		Machines     []string `json:"machines"`
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

	mode := params.Mode
	if params.Mode == "default" {
		mode = ""
	}

	sql := groupedComparisonSQL("branch", params.Table, fmt.Sprintf(
		"branch IN ('%s', '%s') AND measure_name IN (%s) AND (%s) AND mode = '%s'",
		params.Branch1, params.Branch2, quoteAndJoin(params.MeasureNames), machineCondition(params.Machines), mode))
	db, err := t.openDatabaseConnection()
	defer func(db driver.Conn) {
		_ = db.Close()
	}(db)
	if err != nil {
		return nil, false, err
	}

	var queryResults []ownerQueryResult

	err = db.Select(request.Context(), &queryResults, sql)
	if err != nil {
		return nil, false, err
	}

	response := buildGroupedComparisonResponse(queryResults, params.Branch1, params.Branch2)
	buffer, err := toJSONBuffer(response)
	return buffer, true, err
}

func (t *StatsServer) getModeComparison(request *http.Request) (*bytebufferpool.ByteBuffer, bool, error) {
	type requestParams struct {
		Table        string   `json:"table"`
		MeasureNames []string `json:"measure_names"`
		Branch       string   `json:"branch"`
		Machines     []string `json:"machines"`
		Mode1        string   `json:"mode1"`
		Mode2        string   `json:"mode2"`
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

	mode1 := params.Mode1
	if mode1 == "default" {
		mode1 = ""
	}
	mode2 := params.Mode2
	if mode2 == "default" {
		mode2 = ""
	}

	sql := groupedComparisonSQL("mode", params.Table, fmt.Sprintf(
		"mode IN ('%s', '%s') AND branch = '%s' AND measure_name IN (%s) AND (%s)",
		mode1, mode2, params.Branch, quoteAndJoin(params.MeasureNames), machineCondition(params.Machines)))
	db, err := t.openDatabaseConnection()
	defer func(db driver.Conn) {
		_ = db.Close()
	}(db)
	if err != nil {
		return nil, false, err
	}

	var queryResults []ownerQueryResult

	err = db.Select(request.Context(), &queryResults, sql)
	if err != nil {
		return nil, false, err
	}

	response := buildGroupedComparisonResponse(queryResults, mode1, mode2)
	buffer, err := toJSONBuffer(response)
	return buffer, true, err
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
		machineLike := removeLastPart(params.Machine) + "%"
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
	`, table, params.Branch, params.MetricName, machineLike, params.TestName, params.Mode)

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
