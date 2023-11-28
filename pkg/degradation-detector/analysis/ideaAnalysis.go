package analysis

import (
  detector "github.com/JetBrains/ij-perf-report-aggregator/pkg/degradation-detector"
  "net/http"
  "strings"
)

func GenerateIdeaSettings(backendUrl string, client *http.Client) []detector.PerformanceSettings {
  settings := make([]detector.PerformanceSettings, 0, 1000)
  settings = append(settings, generateIdeaOnInstallerAnalysisSettings(backendUrl, client)...)
  settings = append(settings, generateIdeaDevAnalysisSettings(backendUrl, client)...)
  return settings
}

func generateIdeaOnInstallerAnalysisSettings(backendUrl string, client *http.Client) []detector.PerformanceSettings {
  tests := []string{
    "intellij_sources/%", "grails/%", "java/%", "spring_boot/%",
    "spring_boot_maven/%", "spring_boot/%", "kotlin/%", "kotlin_coroutines/%",
    "kotlin_petclinic/%", "community/%", "empty_project/%", "keycloak_release_20/%",
    "space/%", "toolbox_enterprise/%", "train-ticket/%",
  }
  baseSettings := detector.PerformanceSettings{
    Db:      "perfint",
    Table:   "idea",
    Branch:  "master",
    Machine: "intellij-linux-performance-aws-%",
  }
  testsExpanded := detector.ExpandTestsByPattern(backendUrl, client, tests, baseSettings)
  settings := make([]detector.PerformanceSettings, 0, 100)
  for _, test := range testsExpanded {
    metrics := getMetricFromTestName(test)
    for _, metric := range metrics {
      settings = append(settings, detector.PerformanceSettings{
        Db:      baseSettings.Db,
        Table:   baseSettings.Table,
        Branch:  baseSettings.Branch,
        Machine: baseSettings.Machine,
        Project: test,
        Metric:  metric,
        SlackSettings: detector.SlackSettings{
          Channel:     "ij-perf-report-aggregator",
          ProductLink: "intellij",
        },
      })
    }

  }
  return settings
}

func generateIdeaDevAnalysisSettings(backendUrl string, client *http.Client) []detector.PerformanceSettings {
  tests := []string{"intellij_commit/%"}
  baseSettings := detector.PerformanceSettings{
    Db:      "perfintDev",
    Table:   "idea",
    Branch:  "master",
    Machine: "intellij-linux-performance-aws-%",
  }
  testsExpanded := detector.ExpandTestsByPattern(backendUrl, client, tests, baseSettings)
  settings := make([]detector.PerformanceSettings, 0, 100)
  for _, test := range testsExpanded {
    metrics := getMetricFromTestName(test)
    for _, metric := range metrics {
      settings = append(settings, detector.PerformanceSettings{
        Db:      baseSettings.Db,
        Table:   baseSettings.Table,
        Branch:  baseSettings.Branch,
        Machine: baseSettings.Machine,
        Project: test,
        Metric:  metric,
        SlackSettings: detector.SlackSettings{
          Channel:     "ij-perf-report-aggregator",
          ProductLink: "intellij",
        },
      })
    }

  }
  return settings
}

func getMetricFromTestName(test string) []string {
  if strings.Contains(test, "/vfsRefresh") {
    return []string{"vfs_initial_refresh"}
  }
  if strings.Contains(test, "/rebuild") {
    return []string{"build_compilation_duration"}
  }
  if strings.Contains(test, "/inspection") {
    return []string{"globalInspections"}
  }
  if strings.Contains(test, "/localInspection") {
    return []string{"localInspections", "firstCodeAnalysis"}
  }
  if strings.Contains(test, "/completion") {
    return []string{"completion"}
  }
  if strings.Contains(test, "/debug") {
    return []string{"debugRunConfiguration", "debugStep_into"}
  }
  if strings.Contains(test, "/showIntentions") {
    return []string{"Test#average_awt_delay", "showQuickFixes"}
  }
  if strings.Contains(test, "/showFileHistory") {
    return []string{"showFileHistory"}
  }
  if strings.Contains(test, "/expandProjectMenu") {
    return []string{"%expandProjectMenu"}
  }
  if strings.Contains(test, "/expandMainMenu") {
    return []string{"%expandMainMenu"}
  }
  if strings.Contains(test, "/expandEditorMenu") {
    return []string{"%expandEditorMenu"}
  }
  if strings.Contains(test, "/highlight") {
    return []string{"highlighting"}
  }
  if strings.Contains(test, "/FileStructureDialog") {
    return []string{"FileStructurePopup"}
  }
  if strings.Contains(test, "/createJavaClass") {
    return []string{"createJavaFile"}
  }
  if strings.Contains(test, "/createKotlinClass") {
    return []string{"createKotlinFile"}
  }
  if strings.Contains(test, "/indexing") {
    return []string{"scanningTimeWithoutPauses", "indexingTimeWithoutPauses"}
  }
  if strings.Contains(test, "/inlineRename") {
    return []string{"startInlineRename"}
  }
  if strings.Contains(test, "/-scanning") {
    return []string{"scanningTimeWithoutPauses"}
  }
  if strings.Contains(test, "/findUsages") {
    return []string{"findUsages", "fus_find_usages_all", "fus_find_usages_first"}
  }
  if strings.Contains(test, "/go-to-") {
    return []string{"searchEverywhere"}
  }
  if strings.Contains(test, "/ultimate") {
    return []string{"localInspections", "firstCodeAnalysis", "typingCodeAnalyzing", "completion"}
  }
  if strings.Contains(test, "/typing") {
    return []string{"typingCodeAnalyzing", "typing"}
  }
  return []string{}
}
