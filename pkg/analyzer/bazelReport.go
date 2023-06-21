package analyzer

import (
  "bufio"
  "bytes"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/model"
  "strconv"
  "strings"
)

func analyzePerfBazelReport(runResult *RunResult) bool {
  runResult.Report = &model.Report{}

  measureNames := make([]string, 0)
  measureValues := make([]float64, 0)
  measureTypes := make([]string, 0)

  scanner := bufio.NewScanner(bytes.NewReader(runResult.RawReport))
  first := true
  for scanner.Scan() {
    text := scanner.Text()
    if first {
      runResult.Report.Project = text
      first = false
      continue
    }
    split := strings.Fields(text)
    if len(split) == 2 {
      name := split[0]
      value, err := strconv.ParseFloat(strings.Replace(split[1], ",", ".", 1), 64)
      if err == nil {
        measureNames = append(measureNames, name)
        measureValues = append(measureValues, value)
        if strings.HasSuffix(name, "ms") {
          measureTypes = append(measureTypes, "d")
        } else {
          measureTypes = append(measureTypes, "c")
        }
      }
    }
  }

  if len(measureNames) == 0 && len(measureValues) == 0 {
    return true
  }
  runResult.ExtraFieldData = []interface{}{measureNames, measureValues, measureTypes}
  return false
}
