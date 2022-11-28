package analyzer

import (
  "bufio"
  "bytes"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/model"
  "strconv"
  "strings"
)

func analyzePerfJbrReport(runResult *RunResult, data model.ExtraData) bool {
  runResult.Report = &model.Report{}
  buildType := data.TcBuildType
  if strings.Contains(buildType, "Performance") {
    configName := buildType[strings.LastIndex(buildType, "Performance")+12:]
    runResult.Report.Project = configName
  }

  reportURL := runResult.ReportFileName
  fileName := reportURL[strings.LastIndex(reportURL, "/")+1:]
  if strings.Contains(fileName, "regression") {
    return true
  }

  measureNames := make([]string, 0)
  measureValues := make([]float64, 0)
  measureTypes := make([]string, 0)

  scanner := bufio.NewScanner(bytes.NewReader(runResult.RawReport))
  for scanner.Scan() {
    text := scanner.Text()
    split := strings.Split(text, "\t")
    if len(split) == 2 {
      name := strings.TrimSpace(split[0])
      value, err := strconv.ParseFloat(strings.TrimSpace(split[1]), 64)
      if err != nil {
        //for some reason float is written as "54,14" in some tests so we have to convert it back to the normal one "54.14"
        normalizedValue := strings.Replace(split[1], ",", ".", 1)
        value, err = strconv.ParseFloat(normalizedValue, 64)
      }
      if err == nil {
        measureNames = append(measureNames, name)
        measureValues = append(measureValues, value)
        measureTypes = append(measureTypes, "c")
      }
    }
  }

  if len(measureNames) == 0 && len(measureValues) == 0 {
    return true
  }
  runResult.ExtraFieldData = []interface{}{measureNames, measureValues, measureTypes}
  return false
}
