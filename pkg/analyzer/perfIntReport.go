package analyzer

import (
	"github.com/valyala/fastjson"
	"golang.org/x/exp/constraints"
	"log/slog"
	"math"
)

// Numeric is a constraint that permits any numeric type
type Numeric interface {
	constraints.Integer | constraints.Float
}

func analyzePerfReport[T Numeric](runResult *RunResult, data *fastjson.Value) error {
	measureNames := make([]string, 0)
	measureTypes := make([]string, 0)
	measureValues := make([]T, 0)
	for _, measure := range data.GetArray("metrics") {
		measureName := string(measure.GetStringBytes("n"))

		// in milliseconds
		value := measure.Get("d")
		measureType := "d"
		if value == nil {
			value = measure.Get("c")
			measureType = "c"
			if value == nil {
				slog.Error("metric doesn't contain 'd' or 'c', skipping metric", "measureName", measureName, "reportURL", runResult.ReportFileName)
				continue
			}
		}

		floatValue := value.GetFloat64()
		if math.IsNaN(floatValue) {
			slog.Error("invalid value, skipping metric", "measureName", measureName, "value", value, "reportURL", runResult.ReportFileName)
			continue
		}

		var numValue T
		var ok bool
		switch any(numValue).(type) {
		case int32:
			intValue := int32(floatValue)
			if floatValue != float64(intValue) {
				slog.Warn("int expected, but got float, setting metric value to zero", "measureName", measureName, "intValue", intValue, "floatValue", floatValue, "reportURL", runResult.ReportFileName)
				intValue = 0
			}
			numValue, ok = any(intValue).(T)
		case float64:
			numValue, ok = any(floatValue).(T)
		default:
			slog.Error("unexpected type, skipping metric", "type", any(numValue), "measureName", measureName, "reportURL", runResult.ReportFileName)
			continue
		}
		if !ok {
			slog.Warn("unexpected type, skipping metric", "type", any(numValue), "measureName", measureName, "reportURL", runResult.ReportFileName)
			continue
		}

		measureNames = append(measureNames, measureName)
		measureValues = append(measureValues, numValue)
		measureTypes = append(measureTypes, measureType)
	}

	if len(measureNames) == 0 {
		slog.Warn("invalid report - no measures, report will be skipped", "id", runResult.TcBuildId)
		runResult.Report = nil
		return nil
	}

	runResult.ExtraFieldData = []interface{}{measureNames, measureValues, measureTypes}
	return nil
}
