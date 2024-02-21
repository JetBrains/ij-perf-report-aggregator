package analyzer

import (
  "github.com/valyala/fastjson"
  "strings"
)

type CustomReportAnalyzer func(runResult *RunResult, data *fastjson.Value) error
type InsertStatementWriter func(sb *strings.Builder)

type DatabaseConfiguration struct {
  DbName    string
  TableName string

  ReportReader          CustomReportAnalyzer
  insertStatementWriter InsertStatementWriter

  HasProductField             bool
  HasInstallerField           bool
  HasBuildTypeField           bool
  HasRawReport                bool
  HasBuildNumber              bool
  HasNoInstallerButHasChanges bool
  HasMetaDB                   bool
  extraFieldCount             int
}

func GetAnalyzer(id string) DatabaseConfiguration {
  switch {
  case id == "ij":
    fieldNames := []string{
      "classLoadingTime", "classLoadingSearchTime", "classLoadingDefineTime", "classLoadingCount", "classLoadingPreparedCount", "classLoadingLoadedCount",
      "resourceLoadingTime", "resourceLoadingCount",
      "measure.name", "measure.start", "measure.duration", "measure.thread", "metrics.name", "metrics.value",
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
  case id == "ijDev":
    fieldNames := []string{
      "classLoadingTime", "classLoadingSearchTime", "classLoadingDefineTime", "classLoadingCount", "classLoadingPreparedCount", "classLoadingLoadedCount",
      "resourceLoadingTime", "resourceLoadingCount",
      "measure.name", "measure.start", "measure.duration", "measure.thread", "metrics.name", "metrics.value",
    }
    return DatabaseConfiguration{
      DbName:            id,
      HasProductField:   true,
      HasInstallerField: false,
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
  case strings.HasPrefix(id, "perfintDev"):
    dbName, tableName := splitId(id)
    return DatabaseConfiguration{
      DbName:                      dbName,
      TableName:                   tableName,
      ReportReader:                analyzePerfReport,
      HasBuildTypeField:           true,
      HasMetaDB:                   true,
      HasNoInstallerButHasChanges: true,
      extraFieldCount:             3,
      insertStatementWriter: func(sb *strings.Builder) {
        sb.WriteString(", measures.name, measures.value, measures.type")
      },
    }
  case strings.HasPrefix(id, "perfint"):
    dbName, tableName := splitId(id)
    return DatabaseConfiguration{
      DbName:            dbName,
      TableName:         tableName,
      ReportReader:      analyzePerfReport,
      HasInstallerField: true,
      HasBuildTypeField: true,
      HasMetaDB:         true,
      extraFieldCount:   3,
      insertStatementWriter: func(sb *strings.Builder) {
        sb.WriteString(", measures.name, measures.value, measures.type")
      },
    }
  case id == "fleet":
    return DatabaseConfiguration{
      DbName:            "fleet",
      ReportReader:      analyzeFleetReport,
      HasInstallerField: true,
      extraFieldCount:   4,
      insertStatementWriter: func(sb *strings.Builder) {
        sb.WriteString(", measures.name, measures.value, measures.start, measures.thread")
      },
    }
  case id == "perf_fleet":
    return DatabaseConfiguration{
      DbName:          "fleet",
      TableName:       "measure",
      ReportReader:    analyzePerfFleetReport,
      extraFieldCount: 2,
      insertStatementWriter: func(sb *strings.Builder) {
        sb.WriteString(", name, value")
      },
    }
  case id == "jbr":
    return DatabaseConfiguration{
      DbName:          "jbr",
      TableName:       "report",
      extraFieldCount: 3,
      HasBuildNumber:  true,
      insertStatementWriter: func(sb *strings.Builder) {
        sb.WriteString(", measures.name, measures.value, measures.type")
      },
    }
  case id == "bazel":
    return DatabaseConfiguration{
      DbName:          "bazel",
      TableName:       "report",
      extraFieldCount: 3,
      insertStatementWriter: func(sb *strings.Builder) {
        sb.WriteString(", measures.name, measures.value, measures.type")
      },
    }
  case id == "qodana":
    return DatabaseConfiguration{
      DbName:          "qodana",
      TableName:       "report",
      extraFieldCount: 3,
      insertStatementWriter: func(sb *strings.Builder) {
        sb.WriteString(", measures.name, measures.value, measures.type")
      },
    }
  case id == "perfUnitTests":
    return DatabaseConfiguration{
      DbName:                      "perfUnitTests",
      TableName:                   "report",
      ReportReader:                analyzePerfReport,
      HasBuildTypeField:           true,
      HasMetaDB:                   false,
      HasNoInstallerButHasChanges: true,
      extraFieldCount:             3,
      insertStatementWriter: func(sb *strings.Builder) {
        sb.WriteString(", measures.name, measures.value, measures.type")
      },
    }
  default:
    panic("unknown project: " + id)
  }
}

func splitId(id string) (string, string) {
  result := strings.SplitN(id, "_", 2)
  return result[0], result[1]
}
