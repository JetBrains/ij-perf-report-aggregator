package analyzer

import (
  "github.com/valyala/fastjson"
  "go.uber.org/zap"
  "strings"
)

type CustomReportAnalyzer func(runResult *RunResult, data *fastjson.Value, logger *zap.Logger) error
type InsertStatementWriter func(sb *strings.Builder)

type DatabaseConfiguration struct {
  DbName    string
  TableName string

  ReportReader          CustomReportAnalyzer
  insertStatementWriter InsertStatementWriter

  HasProductField   bool
  HasInstallerField bool
  HasBuildTypeField bool
  extraFieldCount   int
}

func GetAnalyzer(id string) DatabaseConfiguration {
  if id == "ij" {
    fieldNames := []string{
      "service.name", "service.start", "service.duration", "service.thread", "service.plugin",
      "classLoadingTime", "classLoadingSearchTime", "classLoadingDefineTime", "classLoadingCount",
      "resourceLoadingTime", "resourceLoadingCount",
      "measure.name", "measure.start", "measure.duration", "measure.thread",
    }
    return DatabaseConfiguration{
      DbName:            id,
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
  } else if strings.HasPrefix(id, "perfintDev") {
    dbName, tableName := splitId(id)
    return DatabaseConfiguration{
      DbName:            dbName,
      TableName:         tableName,
      ReportReader:      analyzePerfReport,
      HasBuildTypeField: true,
      extraFieldCount:   3,
      insertStatementWriter: func(sb *strings.Builder) {
        sb.WriteString(", measures.name, measures.value, measures.type")
      },
    }
  } else if strings.HasPrefix(id, "perfint") {
    dbName, tableName := splitId(id)
    return DatabaseConfiguration{
      DbName:            dbName,
      TableName:         tableName,
      ReportReader:      analyzePerfReport,
      HasInstallerField: true,
      HasBuildTypeField: true,
      extraFieldCount:   3,
      insertStatementWriter: func(sb *strings.Builder) {
        sb.WriteString(", measures.name, measures.value, measures.type")
      },
    }
  } else if id == "fleet" {
    return DatabaseConfiguration{
      DbName:            "fleet",
      ReportReader:      analyzeFleetReport,
      HasInstallerField: true,
      extraFieldCount:   4,
      insertStatementWriter: func(sb *strings.Builder) {
        sb.WriteString(", measures.name, measures.value, measures.start, measures.thread")
      },
    }

  } else if id == "perf_fleet" {
    return DatabaseConfiguration{
      DbName:          "fleet",
      TableName:       "measure",
      ReportReader:    analyzePerfFleetReport,
      extraFieldCount: 2,
      insertStatementWriter: func(sb *strings.Builder) {
        sb.WriteString(", name, value")
      },
    }
  } else {
    panic("unknown project: " + id)
  }
}

func splitId(id string) (string, string) {
  result := strings.SplitN(id, "_", 2)
  return result[0], result[1]
}

func analyzePerfReport(runResult *RunResult, data *fastjson.Value, logger *zap.Logger) error {
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
      logger.Warn("int expected, but got float, setting metric value to zero",
        zap.String("measureName", measureName), zap.Int32("intValue", intValue), zap.Float64("floatValue", floatValue),
        zap.String("reportURL", runResult.ReportFileName))
      intValue = 0
    }

    measureNames = append(measureNames, measureName)
    measureValues = append(measureValues, intValue)
    measureTypes = append(measureTypes, measureType)
  }

  if len(measureNames) == 0 {
    logger.Warn("invalid report - no measures, report will be skipped", zap.Int("id", runResult.TcBuildId))
    runResult.Report = nil
    return nil
  }

  runResult.ExtraFieldData = []interface{}{measureNames, measureValues, measureTypes}
  return nil
}
