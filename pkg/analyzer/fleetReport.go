package analyzer

import (
  "github.com/mcuadros/go-version"
  "github.com/valyala/fastjson"
  "go.uber.org/zap"
  "math"
  "path/filepath"
  "strings"
  "time"
)

func analyzePerfFleetReport(runResult *RunResult, data *fastjson.Value, logger *zap.Logger) error {
  values := data.GetArray("data")

  if len(values) == 0 {
    logger.Warn("invalid report - no measures, report will be skipped", zap.Int("id", runResult.TcBuildId))
    runResult.Report = nil
    return nil
  }

  first := values[0]
  runResult.GeneratedTime = time.Unix(0, first.GetInt64("epochNanos")).Unix()
  runResult.Report.Project = string(first.GetStringBytes("attributes", "test.name"))

  measureName := strings.TrimSuffix(filepath.Base(runResult.ReportFileName), ".json")
  // convert float milliseconds to nanoseconds
  value := uint64(math.Round(first.GetFloat64("value") * 1_000_000))

  runResult.ExtraFieldData = []interface{}{measureName, value}
  return nil
}

func analyzeFleetReport(runResult *RunResult, data *fastjson.Value, _ *zap.Logger) error {
  names := make([]string, 0)
  values := make([]int, 0)
  starts := make([]int, 0)
  threads := make([]string, 0)
  items := data.GetArray("items")
  for _, measure := range items {
    name := string(measure.GetStringBytes("name"))
    // in milliseconds
    names = append(names, name)
    values = append(values, measure.GetInt("duration"))
    starts = append(starts, measure.GetInt("start"))
    threads = append(threads, string(measure.GetStringBytes("thread")))
  }

  mapNameV22 := version.Compare(runResult.Report.Version, "22", "<=")
  isLessThan36 := version.Compare(runResult.Report.Version, "36", "<")

  for _, groupField := range []string{"items", "prepareAppInitActivities"} {
    for _, measure := range data.GetArray(groupField) {
      name := string(measure.GetStringBytes("n"))
      if len(name) == 0 {
        continue
      }

      if mapNameV22 {
        if name == "create window" {
          name = "editor appeared"
        } else if name == "render" {
          name = "window appeared"
        }
      } else if isLessThan36 {
        if name == "render editor" {
          name = "editor appeared"
        } else if name == "render real panels" {
          name = "window appeared"
        }
      }

      // in milliseconds
      names = append(names, name)
      values = append(values, measure.GetInt("d"))
      starts = append(starts, measure.GetInt("s"))
      threads = append(threads, string(measure.GetStringBytes("t")))
    }
  }

  if len(names) == 0 {
    return nil
  }

  runResult.ExtraFieldData = []interface{}{names, values, starts, threads}
  return nil
}
