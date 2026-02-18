package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"math"
	"net/http"
	"slices"
	"strings"
	"sync"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/JetBrains/ij-perf-report-aggregator/pkg/degradation-detector/statistic"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var mainMetrics = []string{
	"indexingTimeWithoutPauses",
	"scanningTimeWithoutPauses",
	"dumbModeWithPauses",
	"vfs_initial_refresh",
	"build_compilation_duration",
	"globalInspections",
	"findUsages",
	"localInspections",
	"firstCodeAnalysis",
	"completion",
	"searchEverywhere",
	"showFileHistory",
	"%expandMainMenu",
	"%expandProjectMenu",
	"%expandEditorMenu",
	"FileStructurePopup",
	"createKotlinFile",
	"highlighting",
	"vcs-log-indexing",
	"startInlineRename",
	"debugRunConfiguration",
	"debugStep_into",
	"searchEverywhere_dialog_shown",
	"showQuickFixes",
	"createJavaFile",
	"typingCodeAnalyzing",
	"test#max_awt_delay",
	"MatchedRatio",
	"SyntaxErrorsSessionRatio",
	"EditSimilarity",
	"attempt.mean.ms",
}

type compareByOwnerRequest struct {
	Owner             string   `json:"owner"`
	BaseBranch        string   `json:"baseBranch"`
	CompareBranch     string   `json:"compareBranch"`
	Machine           string   `json:"machine"`
	Mode              string   `json:"mode"`
	AdditionalMetrics []string `json:"additionalMetrics"`
}

type comparisonResponseItem struct {
	Project            string  `json:"project"`
	Metric             string  `json:"metric"`
	BaseBranchValue    float64 `json:"baseBranchValue"`
	CompareBranchValue float64 `json:"compareBranchValue"`
	Diff               float64 `json:"diff"`
	Link               string  `json:"link"`
}

type branchMedianItem struct {
	Branch      string
	Project     string
	MeasureName string
	Median      float64
	DbName      string
	TableName   string
}

type projectOwnerEntry struct {
	Project   string
	DbName    string
	TableName string
}

type dbTableKey struct {
	DbName    string
	TableName string
}

func (t *StatsServer) CreateCompareByOwnerHandler(metaDb *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		var params compareByOwnerRequest
		decoder := json.NewDecoder(request.Body)
		defer request.Body.Close()
		if err := decoder.Decode(&params); err != nil {
			http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
			return
		}

		if params.Owner == "" || params.BaseBranch == "" || params.CompareBranch == "" || params.Machine == "" {
			http.Error(w, "owner, baseBranch, compareBranch, and machine are required", http.StatusBadRequest)
			return
		}

		if params.Mode == "default" {
			params.Mode = ""
		}

		entries, err := getProjectsByOwner(request.Context(), metaDb, params.Owner)
		if err != nil {
			slog.Error("unable to get projects by owner", "error", err)
			http.Error(w, "Failed to get projects for owner", http.StatusInternalServerError)
			return
		}
		if len(entries) == 0 {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte("[]"))
			return
		}

		machineLike := "%" + params.Machine + "%"

		metrics := buildMetricsList(params.AdditionalMetrics)

		quotedMetrics := make([]string, len(metrics))
		for i, m := range metrics {
			quotedMetrics[i] = "'" + m + "'"
		}
		metricsStr := strings.Join(quotedMetrics, ",")

		db, err := t.openDatabaseConnection()
		if err != nil {
			slog.Error("unable to open database connection", "error", err)
			http.Error(w, "Failed to open database connection", http.StatusInternalServerError)
			return
		}
		defer func(db driver.Conn) {
			_ = db.Close()
		}(db)

		var mu sync.Mutex
		var allItems []branchMedianItem
		var wg sync.WaitGroup

		// Group projects by db+table to query only relevant combinations
		dbTableProjects := make(map[dbTableKey][]string)
		for _, e := range entries {
			key := dbTableKey{DbName: e.DbName, TableName: e.TableName}
			dbTableProjects[key] = append(dbTableProjects[key], e.Project)
		}

		for key, projects := range dbTableProjects {
			projectsStr := quoteAndJoin(projects)
			wg.Go(func() {
				items, queryErr := queryTableForComparison(request.Context(), db, key.DbName, key.TableName, params.BaseBranch, params.CompareBranch, metricsStr, machineLike, params.Mode, projectsStr)
				if queryErr != nil {
					slog.Error("failed to query table", "db", key.DbName, "table", key.TableName, "error", queryErr)
					return
				}
				for i := range items {
					items[i].DbName = key.DbName
					items[i].TableName = key.TableName
				}
				mu.Lock()
				allItems = append(allItems, items...)
				mu.Unlock()
			})
		}
		wg.Wait()

		response := buildComparisonResponse(allItems, params.BaseBranch, params.CompareBranch, params.Machine)

		jsonData, err := json.Marshal(response)
		if err != nil {
			http.Error(w, "Failed to marshal response: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(jsonData)
	}
}

func getProjectsByOwner(ctx context.Context, metaDb *pgxpool.Pool, owner string) ([]projectOwnerEntry, error) {
	rows, err := metaDb.Query(ctx, "SELECT project, db_name, table_name FROM project_owner WHERE owner=$1", owner)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	entries, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (projectOwnerEntry, error) {
		var e projectOwnerEntry
		err := row.Scan(&e.Project, &e.DbName, &e.TableName)
		return e, err
	})
	if err != nil {
		return nil, err
	}
	if entries == nil {
		entries = []projectOwnerEntry{}
	}
	return entries, nil
}

func quoteAndJoin(items []string) string {
	quoted := make([]string, len(items))
	for i, s := range items {
		quoted[i] = "'" + s + "'"
	}
	return strings.Join(quoted, ",")
}

func buildMetricsList(additionalMetrics []string) []string {
	seen := make(map[string]bool, len(mainMetrics)+len(additionalMetrics))
	result := make([]string, 0, len(mainMetrics)+len(additionalMetrics))
	for _, m := range mainMetrics {
		if !seen[m] {
			seen[m] = true
			result = append(result, m)
		}
	}
	for _, m := range additionalMetrics {
		if !seen[m] {
			seen[m] = true
			result = append(result, m)
		}
	}
	return result
}

func queryTableForComparison(ctx context.Context, db driver.Conn, dbName, table, baseBranch, compareBranch, metricsStr, machine, mode, projectsStr string) ([]branchMedianItem, error) {
	sql := fmt.Sprintf(
		"SELECT branch AS Branch, project AS Project, measure_name AS MeasureName, "+
			"arraySlice(groupArray(measure_value), 1, 50) AS MeasureValues "+
			"FROM ("+
			"SELECT branch, project, measures.name AS measure_name, measures.value AS measure_value "+
			"FROM %s.%s ARRAY JOIN measures "+
			"WHERE branch IN ('%s', '%s') "+
			"AND measure_name IN (%s) "+
			"AND machine LIKE '%s' "+
			"AND mode = '%s' "+
			"AND project IN (%s) "+
			"AND generated_time > now() - interval 1 month "+
			"ORDER BY generated_time DESC"+
			") "+
			"GROUP BY branch, project, measure_name",
		dbName, table, baseBranch, compareBranch, metricsStr, machine, mode, projectsStr,
	)

	var queryResults []struct {
		Branch        string
		Project       string
		MeasureName   string
		MeasureValues []int
	}

	if err := db.Select(ctx, &queryResults, sql); err != nil {
		return nil, err
	}

	resultChan := make(chan branchMedianItem, len(queryResults))
	var wg sync.WaitGroup
	for _, result := range queryResults {
		wg.Go(func() {
			values := make([]int, len(result.MeasureValues))
			copy(values, result.MeasureValues)
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

			resultChan <- branchMedianItem{
				Branch:      result.Branch,
				Project:     result.Project,
				MeasureName: result.MeasureName,
				Median:      median,
			}
		})
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	items := make([]branchMedianItem, 0, len(queryResults))
	for item := range resultChan {
		items = append(items, item)
	}
	return items, nil
}

func buildComparisonResponse(items []branchMedianItem, baseBranch, compareBranch, machine string) []comparisonResponseItem {
	type key struct {
		Project string
		Metric  string
	}

	baseMap := make(map[key]float64)
	compareMap := make(map[key]float64)
	projectDbTable := make(map[string]dbTableKey)

	for _, item := range items {
		k := key{Project: item.Project, Metric: item.MeasureName}
		if item.Branch == baseBranch {
			baseMap[k] = item.Median
		} else if item.Branch == compareBranch {
			compareMap[k] = item.Median
		}
		if _, exists := projectDbTable[item.Project]; !exists {
			projectDbTable[item.Project] = dbTableKey{DbName: item.DbName, TableName: item.TableName}
		}
	}

	response := make([]comparisonResponseItem, 0)
	for k, baseVal := range baseMap {
		compareVal, ok := compareMap[k]
		if !ok {
			continue
		}

		var diff float64
		if baseVal != 0 {
			diff = math.Round(((compareVal-baseVal)/baseVal)*1000) / 10
		}

		dt := projectDbTable[k.Project]
		link := buildTestLink(dt.DbName, dt.TableName, machine, baseBranch, compareBranch, k.Project, k.Metric)

		response = append(response, comparisonResponseItem{
			Project:            k.Project,
			Metric:             k.Metric,
			BaseBranchValue:    baseVal,
			CompareBranchValue: compareVal,
			Diff:               diff,
			Link:               link,
		})
	}

	return response
}

func buildTestLink(dbName, table, machine, baseBranch, compareBranch, project, metric string) string {
	return fmt.Sprintf("/owners/test?dbName=%s&table=%s&machine=%s&branch=%s&branch=%s&project=%s&measure=%s",
		dbName, table, machine, baseBranch, compareBranch, project, metric)
}
