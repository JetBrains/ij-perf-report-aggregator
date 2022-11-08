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
  runResult.Report.Project = data.TcBuildType

  measureNames := make([]string, 0)
  measureValues := make([]float64, 0)
  measureTypes := make([]string, 0)

  scanner := bufio.NewScanner(bytes.NewReader(runResult.RawReport))
  for scanner.Scan() {
    text := scanner.Text()
    split := strings.Split(text, "\t")
    if len(split) == 2 {
      name := split[0]
      value, err := strconv.ParseFloat(split[1], 64)
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
