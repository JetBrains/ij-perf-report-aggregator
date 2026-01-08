package analyzer

import (
	"bytes"
	"log/slog"
	"sort"

	"github.com/JetBrains/ij-perf-report-aggregator/pkg/model"
	"github.com/valyala/fastjson"
)

type measureItem struct {
	name     string
	start    uint32
	duration uint32
	thread   string
}

func analyzeIjReport(runResult *RunResult, data *fastjson.Value) error {
	report := runResult.Report

	report.TotalDuration = data.GetInt("totalDuration")
	if report.TotalDuration == 0 {
		report.TotalDuration = data.GetInt("totalDurationActual")
	}

	traceEvents := data.GetArray("traceEvents")

	for _, v := range traceEvents {
		report.TraceEvents = append(report.TraceEvents, model.TraceEvent{
			Name:      string(v.GetStringBytes("name")),
			Phase:     string(v.GetStringBytes("ph")),
			Timestamp: v.GetInt("ts"),
			Category:  string(v.GetStringBytes("cat")),
		})
	}

	var clTotal int32
	var clSearch int32
	var clDefine int32
	var clCount int32
	var clPreparedCount int32
	var clLoadedCount int32

	var rlTime int32
	var rlCount int32

	report.Activities = readActivities("items", data)

	measures := make([]measureItem, 0, len(report.Activities)+1)

	for _, activity := range report.Activities {
		measures = append(measures, measureItem{
			name:     activity.Name,
			start:    uint32(activity.Start),
			duration: uint32(activity.Duration),
			thread:   activity.Thread,
		})
	}

	classLoading := data.Get("classLoading")
	resourceLoading := data.Get("resourceLoading")
	if classLoading != nil && resourceLoading != nil {
		clTotal = int32(classLoading.GetInt("time"))
		clSearch = int32(classLoading.GetInt("searchTime"))
		clDefine = int32(classLoading.GetInt("defineTime"))
		clCount = int32(classLoading.GetInt("count"))

		clPreparedCount = int32(classLoading.GetInt("preparedCount"))
		clLoadedCount = int32(classLoading.GetInt("loadedCount"))

		rlTime = int32(resourceLoading.GetInt("time"))
		rlCount = int32(resourceLoading.GetInt("count"))
	}

	measures = append(measures, measureItem{
		name:     "elementTypeCount",
		start:    0,
		duration: uint32(data.GetInt("langLoading", "elementTypeCount")),
		thread:   "",
	})

	// Sort for better compression (same data pattern across column values). It is confirmed by experiment.
	sort.Slice(measures, func(i, j int) bool {
		return measures[i].name < measures[j].name
	})

	measureCount := len(measures)
	measureName := make([]string, measureCount)
	measureStart := make([]uint32, measureCount)
	measureDuration := make([]uint32, measureCount)
	measureThread := make([]string, measureCount)
	for i, info := range measures {
		measureName[i] = info.name
		measureStart[i] = info.start
		measureDuration[i] = info.duration
		measureThread[i] = info.thread
	}

	metricNames := make([]string, 0)
	metricValues := make([]uint32, 0)
	additionalMetrics := data.GetObject("additionalMetrics")
	if additionalMetrics != nil {
		additionalMetrics.Visit(func(groupName []byte, v *fastjson.Value) {
			v.GetObject().Visit(func(metricName []byte, v *fastjson.Value) {
				value, err := v.Int()
				if err != nil {
					slog.Warn("invalid value", "id", runResult.TcBuildId, "generated", report.Generated, "metricName", string(metricName))
					return
				}
				metricNames = append(metricNames, string(groupName)+"/"+string(metricName))
				metricValues = append(metricValues, uint32(value))
			})
		})
	}

	runResult.ExtraFieldData = []any{
		clTotal, clSearch, clDefine, clCount, clPreparedCount, clLoadedCount, rlTime, rlCount,
		measureName, measureStart, measureDuration, measureThread, metricNames, metricValues,
	}
	return nil
}

func readActivities(key string, value *fastjson.Value) []model.Activity {
	array := value.GetArray(key)
	result := make([]model.Activity, 0, len(array))
	scheduledSuffix := []byte(": scheduled")
	completingSuffix := []byte(": completing")
	for _, v := range array {
		start := v.GetInt("s")
		duration := v.GetInt("d")

		ownDuration := v.GetInt("od")
		if ownDuration == 0 {
			ownDuration = duration
		}

		name := v.GetStringBytes("n")
		if bytes.HasSuffix(name, scheduledSuffix) || bytes.HasSuffix(name, completingSuffix) {
			continue
		}

		result = append(result, model.Activity{
			Name:     string(v.GetStringBytes("n")),
			Thread:   string(v.GetStringBytes("t")),
			Start:    start,
			End:      start + duration,
			Duration: ownDuration,
		})
	}
	return result
}
