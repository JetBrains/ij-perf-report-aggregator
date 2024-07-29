package setting

import (
  detector "github.com/JetBrains/ij-perf-report-aggregator/pkg/degradation-detector"
  "net/http"
  "slices"
  "strings"
)

func GenerateIdeaIndexingSettings(backendUrl string, client *http.Client) []detector.PerformanceSettings {
  return slices.Concat(
    generateIdeaIndexingSettings(backendUrl, client),
  )
}

func generateIdeaIndexingSettings(backendUrl string, client *http.Client) []detector.PerformanceSettings {
  tests := []string{"%/indexing", "%/%-scanning"}
  baseSettings := detector.PerformanceSettings{
    Db:      "perfintDev",
    Table:   "idea",
    Branch:  "master",
    Machine: "intellij-linux-performance-aws-%",
  }
  testsExpanded := detector.ExpandTestsByPattern(backendUrl, client, tests, baseSettings)
  settings := make([]detector.PerformanceSettings, 0, 100)
  machines := []string{"intellij-linux-performance-aws-%", "intellij-windows-performance-aws-%"}
  for _, machine := range machines {
    for _, test := range testsExpanded {
      metrics := getIndexingMetricFromTestNameForIDEA(test)
      for _, metric := range metrics {
        settings = append(settings, detector.PerformanceSettings{
          Db:      baseSettings.Db,
          Table:   baseSettings.Table,
          Branch:  baseSettings.Branch,
          Machine: machine,
          Project: test,
          Metric:  metric,
          SlackSettings: detector.SlackSettings{
            Channel:     "ij-indexes-perf-alerts",
            ProductLink: "intellij",
          },
          AnalysisSettings: detector.AnalysisSettings{MinimumSegmentLength: 8},
        })
      }
    }
  }
  return settings
}

func getIndexingMetricFromTestNameForIDEA(test string) []string {
  if strings.Contains(test, "/indexing") {
    return []string{"scanningTimeWithoutPauses", "indexingTimeWithoutPauses", "processingTime#Kotlin", "processingTime#JAVA"}
  }
  if strings.Contains(test, "-scanning") {
    return []string{"scanningTimeWithoutPauses"}
  }
  return []string{}
}
