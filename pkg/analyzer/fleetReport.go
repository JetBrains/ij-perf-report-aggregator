package analyzer

import (
  "github.com/valyala/fastjson"
  "log/slog"
  "math"
  "path/filepath"
  "strings"
  "time"
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
  runResult.Report.Project = string(first.GetStringBytes("attributes", "test.name"))

  measureName := strings.TrimSuffix(filepath.Base(runResult.ReportFileName), ".json")
  // convert float milliseconds to nanoseconds
  value := uint64(math.Round(first.GetFloat64("value") * 1_000_000))

  runResult.ExtraFieldData = []interface{}{measureName, value}
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
      if len(name) == 0 {
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

  runResult.ExtraFieldData = []interface{}{names, values, starts, threads}
  return nil
}
