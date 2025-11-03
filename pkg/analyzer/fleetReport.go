package analyzer

import (
	"log/slog"
	"path/filepath"
	"strings"
	"time"

	"github.com/valyala/fastjson"
)

func analyzePerfFleetReport(runResult *RunResult, data *fastjson.Value) error {
	values := data.GetArray("data")

	if len(values) == 0 {
		slog.Warn("invalid report - no measures, report will be skipped", "id", runResult.TcBuildId)
		runResult.Report = nil
		return nil
	}

	first := values[0]
	runResult.GeneratedTime = time.Unix(0, first.GetInt64("epochNanos"))
	runResult.Report.Project = strings.ReplaceAll(filepath.Base(filepath.Dir(runResult.ReportFileName)), "%20", " ")

	fileName := filepath.Base(runResult.ReportFileName)
	metricNames := make([]string, 0)
	metricValues := make([]float64, 0)
	metricTypes := make([]string, 0)
	if strings.Contains(fileName, "histogram") {
		pmetrics := histogramToMetrics(data)
		for _, pmetric := range pmetrics {
			metricNames = append(metricNames, pmetric.Key)
			metricValues = append(metricValues, pmetric.Value)
			metricTypes = append(metricTypes, "d")
		}
	} else {
		measureName := strings.TrimSuffix(fileName, ".json")
		value := first.GetFloat64("value")
		metricNames = append(metricNames, measureName)
		metricValues = append(metricValues, value)
		metricTypes = append(metricTypes, "d")
	}

	runResult.ExtraFieldData = []any{metricNames, metricValues, metricTypes}
	return nil
}

func analyzeFleetReport(runResult *RunResult, data *fastjson.Value) error {
	names := make([]string, 0)
	values := make([]int32, 0)
	starts := make([]int32, 0)
	threads := make([]string, 0)
	items := data.GetArray("items")
	for _, measure := range items {
		name := string(measure.GetStringBytes("n"))
		// in milliseconds
		names = append(names, name)
		values = append(values, int32(measure.GetInt("d")))
		starts = append(starts, int32(measure.GetInt("s")))
		threads = append(threads, string(measure.GetStringBytes("t")))
	}

	for _, groupField := range []string{"items", "prepareAppInitActivities"} {
		for _, measure := range data.GetArray(groupField) {
			name := string(measure.GetStringBytes("n"))
			if name == "" {
				continue
			}

			// in milliseconds
			names = append(names, name)
			values = append(values, int32(measure.GetInt("d")))
			starts = append(starts, int32(measure.GetInt("s")))
			threads = append(threads, string(measure.GetStringBytes("t")))
		}
	}

	if len(names) == 0 {
		return nil
	}

	runResult.ExtraFieldData = []any{names, values, starts, threads}
	return nil
}
