package analyzer

import (
  "github.com/develar/errors"
  "github.com/mcuadros/go-version"
  "github.com/valyala/fastjson"
  "go.uber.org/zap"
  "path/filepath"
  "strings"
)

type CustomReportAnalyzer func(runResult *RunResult, data *fastjson.Value, logger *zap.Logger) error
type InsertStatementWriter func(sb *strings.Builder)

type DatabaseConfiguration struct {
  ReportReader          CustomReportAnalyzer
  insertStatementWriter InsertStatementWriter

  HasProductField   bool
  HasInstallerField bool
  extraFieldCount   int
}

func GetAnalyzer(dbName string) DatabaseConfiguration {
  if dbName == "ij" {
    fieldNames := []string{
      "service.name", "service.start", "service.duration", "service.thread", "service.plugin",
      "classLoadingTime", "classLoadingSearchTime", "classLoadingDefineTime", "classLoadingCount",
      "resourceLoadingTime", "resourceLoadingCount",
      "measure.name", "measure.start", "measure.duration", "measure.thread",
    }
    return DatabaseConfiguration{
      HasProductField:   true,
      HasInstallerField: true,
      extraFieldCount:   len(IjMetricDescriptors) + len(fieldNames),
      ReportReader:      analyzeIjReport,
      insertStatementWriter: func(sb *strings.Builder) {
        for _, metric := range IjMetricDescriptors {
          sb.WriteRune(',')
          sb.WriteString(metric.Name)
        }
        for _, fieldName := range fieldNames {
          sb.WriteRune(',')
          sb.WriteString(fieldName)
        }
      },
    }
  } else if dbName == "sharedIndexes" || strings.HasSuffix(dbName, "perfint") {
    return DatabaseConfiguration{
      ReportReader:      analyzeSharedIndexesReport,
      HasInstallerField: true,
      extraFieldCount:   2,
      insertStatementWriter: func(sb *strings.Builder) {
        sb.WriteString(", measures.name, measures.value")
      },
    }
  } else if dbName == "fleet" {
    return DatabaseConfiguration{
      ReportReader:      analyzeFleetReport,
      HasInstallerField: true,
      extraFieldCount:   4,
      insertStatementWriter: func(sb *strings.Builder) {
        sb.WriteString(", measures.name, measures.value, measures.start, measures.thread")
      },
    }

  } else if dbName == "perffleet" {
    return DatabaseConfiguration{
      ReportReader:    analyzePerfFleetReport,
      extraFieldCount: 4,
      insertStatementWriter: func(sb *strings.Builder) {
        sb.WriteString(", measures.name, measures.value")
      },
    }
  } else {
    panic("unknown db: " + dbName)
  }
}

func analyzeSharedIndexesReport(runResult *RunResult, data *fastjson.Value, logger *zap.Logger) error {
  measureNames := make([]string, 0)
  measureValues := make([]int, 0)
  for _, measure := range data.GetArray("metrics") {
    measureName := string(measure.GetStringBytes("n"))

    // in milliseconds
    value := measure.Get("d")
    if value == nil {
      value = measure.Get("c")
      if value == nil {
        return nil
      }
    }

    floatValue := value.GetFloat64()
    intValue := int(floatValue)
    if floatValue != float64(intValue) {
      return errors.WithMessagef(nil, "int expected, but got float %f", floatValue)
    }

    measureNames = append(measureNames, measureName)
    measureValues = append(measureValues, intValue)
  }

  if len(measureNames) == 0 {
    logger.Warn("invalid report - no measures, report will be skipped", zap.Int("id", runResult.TcBuildId))
    runResult.Report = nil
    return nil
  }

  runResult.extraFieldData = []interface{}{measureNames, measureValues}
  return nil
}

func analyzePerfFleetReport(runResult *RunResult, data *fastjson.Value, logger *zap.Logger) error {
  measureNames := make([]string, 0)
  measureValues := make([]float64, 0)
  for _, measure := range data.GetArray("data") {
    measureName := strings.TrimSuffix(filepath.Base(runResult.ReportFileName), ".json")

    // in milliseconds
    value := measure.GetFloat64("value")

    measureNames = append(measureNames, measureName)
    measureValues = append(measureValues, value)
  }

  if len(measureNames) == 0 {
    logger.Warn("invalid report - no measures, report will be skipped", zap.Int("id", runResult.TcBuildId))
    runResult.Report = nil
    return nil
  }

  runResult.extraFieldData = []interface{}{measureNames, measureValues}
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

  runResult.extraFieldData = []interface{}{names, values, starts, threads}
  return nil
}
