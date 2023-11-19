package analyzer

import (
  "github.com/valyala/fastjson"
  "log/slog"
)

func analyzePerfReport(runResult *RunResult, data *fastjson.Value) error {
  measureNames := make([]string, 0)
  measureTypes := make([]string, 0)
  measureValues := make([]int32, 0)
  for _, measure := range data.GetArray("metrics") {
    measureName := string(measure.GetStringBytes("n"))

    // in milliseconds
    value := measure.Get("d")
    measureType := "d"
    if value == nil {
      value = measure.Get("c")
      measureType = "c"
      if value == nil {
        return nil
      }
    }

    floatValue := value.GetFloat64()
    intValue := int32(floatValue)
    if floatValue != float64(intValue) {
      slog.Warn("int expected, but got float, setting metric value to zero", "measureName", measureName, "intValue", intValue, "floatValue", floatValue, "reportURL", runResult.ReportFileName)
      intValue = 0
    }

    measureNames = append(measureNames, measureName)
    measureValues = append(measureValues, intValue)
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
