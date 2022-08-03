package analyzer

import (
  "github.com/valyala/fastjson"
  "go.uber.org/zap"
  "strconv"
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
  HasRawReport      bool
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
      HasRawReport:      true,
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
      HasRawReport:      true,
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
      HasRawReport:      true,
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
      HasRawReport:      true,
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

func MapPerfMeasureName(measureName string, measureNames []string, logger *zap.Logger) string {
  mappedMeasureName := measureName

  if measureName == "CPU | Load | 75th pctl" {
    mappedMeasureName = "cpuLoad75th"
  }

  isMetricNameSatisfies := func(f func(string) bool) bool {
    changedMetrics := []string{"delayType", "doComplete", "findUsages", "doLocalInspection", "altEnter"}
    for _, changedMetric := range changedMetrics {
      if f(changedMetric) {
        return true
      }
    }
    return false
  }

  //increase number for metrics like doCompletion_1
  if isMetricNameSatisfies(func(changedMetric string) bool { return strings.HasPrefix(measureName, changedMetric+"_") }) {
    split := strings.SplitN(mappedMeasureName, "_", 2)
    index, err := strconv.Atoi(split[1])
    if err == nil {
      mappedMeasureName = split[0] + "_" + strconv.Itoa(index+1)
    } else {
      spiltAttributes := strings.SplitN(measureName, "#", 2)
      mappedMeasureName = MapPerfMeasureName(spiltAttributes[0], measureNames, logger) + "#" + spiltAttributes[1]
    }
  }

  //add _1 if there are multiple metrics named the same
  var isSingle = true
  measureNameForAttribute := strings.SplitN(measureName, "#", 2)[0]
  for _, name := range measureNames {
    if name == measureNameForAttribute+"_1" {
      isSingle = false
      break
    }
  }
  if !isSingle && isMetricNameSatisfies(func(changedMetric string) bool { return measureName == changedMetric }) {
    mappedMeasureName = mappedMeasureName + "_1"
  }
  //update metric name for attribute as well
  if isMetricNameSatisfies(func(changedMetric string) bool { return strings.HasPrefix(measureName, changedMetric+"#") }) {
    split := strings.SplitN(mappedMeasureName, "#", 2)
    mappedMeasureName = MapPerfMeasureName(split[0], measureNames, logger) + "#" + split[1]
  }
  if strings.Contains(mappedMeasureName, "#number") {
    mappedMeasureName = strings.SplitN(mappedMeasureName, "#", 2)[0] + "#number"
  }

  updateMeasureName := func(measureName string, oldName string, newName string) string {
    if strings.HasPrefix(measureName, oldName) {
      return newName + strings.TrimPrefix(measureName, oldName)
    } else {
      return measureName
    }
  }
  mappedMeasureName = updateMeasureName(mappedMeasureName, "completion_execution_time", "completion")
  mappedMeasureName = updateMeasureName(mappedMeasureName, "doComplete", "completion")
  mappedMeasureName = updateMeasureName(mappedMeasureName, "typing_total_time", "typing")
  mappedMeasureName = updateMeasureName(mappedMeasureName, "delayType", "typing")
  mappedMeasureName = updateMeasureName(mappedMeasureName, "find_usages_execution_time", "findUsages")
  mappedMeasureName = updateMeasureName(mappedMeasureName, "doLocalInspection", "localInspections")
  mappedMeasureName = updateMeasureName(mappedMeasureName, "local_inspection_execution_time", "localInspections")
  mappedMeasureName = updateMeasureName(mappedMeasureName, "inspection_execution_time", "globalInspections")
  mappedMeasureName = updateMeasureName(mappedMeasureName, "altEnter", "showQuickFixes")
  mappedMeasureName = updateMeasureName(mappedMeasureName, "show_intention_execution_time", "showQuickFixes")
  mappedMeasureName = updateMeasureName(mappedMeasureName, "responsiveness_time", "typing#max_awt_delay")
  mappedMeasureName = updateMeasureName(mappedMeasureName, "average_responsiveness_time", "typing#average_awt_delay")
  if strings.HasPrefix(mappedMeasureName, "find_usages_number_of_found_usages") {
    mappedMeasureName = "findUsages#number"
  }
  if measureName != mappedMeasureName {
    logger.Info("Converting " + measureName + " to " + mappedMeasureName)
  }
  return mappedMeasureName
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

  mappedMeasureNames := make([]string, 0)
  if runResult.branch <= "222" {
    for _, measure := range measureNames {
      mappedMeasureNames = append(mappedMeasureNames, MapPerfMeasureName(measure, measureNames, logger))
    }
  } else {
    mappedMeasureNames = measureNames
  }

  if len(mappedMeasureNames) == 0 {
    logger.Warn("invalid report - no measures, report will be skipped", zap.Int("id", runResult.TcBuildId))
    runResult.Report = nil
    return nil
  }

  runResult.ExtraFieldData = []interface{}{mappedMeasureNames, measureValues, measureTypes}
  return nil
}
