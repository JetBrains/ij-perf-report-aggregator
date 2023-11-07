package analysis

import (
  detector "github.com/JetBrains/ij-perf-report-aggregator/pkg/degradation-detector"
  "strings"
)

func GenerateIdeaSettings() []detector.Settings {
  settings := make([]detector.Settings, 0, 1000)
  settings = append(settings, generateIdeaOnInstallerAnalysisSettings()...)
  settings = append(settings, generateIdeaDevAnalysisSettings()...)
  return settings
}

func generateIdeaOnInstallerAnalysisSettings() []detector.Settings {
  tests := []string{"intellij_sources/vfsRefresh/default", "intellij_sources/vfsRefresh/with-1-thread(s)", "intellij_sources/vfsRefresh/git-status",
    "community/rebuild", "intellij_sources/rebuild", "grails/rebuild", "java/rebuild", "spring_boot/rebuild", "java/inspection", "grails/inspection",
    "spring_boot_maven/inspection", "spring_boot/inspection", "kotlin/inspection", "kotlin_coroutines/inspection",
    "intellij_sources/localInspection/java_file", "intellij_sources/localInspection/kotlin_file", "kotlin/localInspection",
    "kotlin_coroutines/localInspection", "intellij_sources/localInspection/java_file", "intellij_sources/localInspection/kotlin_file",
    "kotlin/localInspection", "kotlin_coroutines/localInspection", "community/completion/kotlin_file", "grails/completion/groovy_file",
    "grails/completion/java_file", "kotlin_petclinic/debug", "grails/showIntentions/Find cause", "kotlin/showIntention/Import", "spring_boot/showIntentions",
    "intellij_sources/showFileHistory/EditorImpl", "intellij_sources/expandProjectMenu", "intellij_sources/expandMainMenu", "intellij_sources/expandEditorMenu",
    "kotlin/highlight", "kotlin_coroutines/highlight", "intellij_sources/FileStructureDialog/java_file", "intellij_sources/FileStructureDialog/kotlin_file",
    "intellij_sources/createJavaClass", "intellij_sources/createKotlinClass",
  }
  settings := make([]detector.Settings, 0, 100)
  for _, test := range tests {
    metrics := getMetricFromTestName(test)
    for _, metric := range metrics {
      settings = append(settings, detector.Settings{
        Db:          "perfint",
        Table:       "idea",
        Channel:     "ij-perf-report-aggregator",
        Branch:      "master",
        Machine:     "intellij-linux-performance-aws-%",
        Test:        test,
        Metric:      metric,
        ProductLink: "intellij",
      })
    }

  }
  return settings
}

func generateIdeaDevAnalysisSettings() []detector.Settings {
  tests := []string{"intellij_commit/indexing", "intellij_commit/second-scanning", "intellij_commit/third-scanning", "intellij_commit/findUsages/Application_runReadAction",
    "intellij_commit/findUsages/Library_getName", "intellij_commit/findUsages/PsiManager_getInstance", "intellij_commit/findUsages/PropertyMapping_value",
    "intellij_commit/findUsages/ActionsKt_runReadAction", "intellij_commit/findUsages/DynamicPluginListener_TOPIC", "intellij_commit/findUsages/Path_div",
    "intellij_commit/findUsages/Persistent_absolutePath", "intellij_sources/localInspection/java_file",
    "intellij_sources/localInspection/kotlin_file",
    "intellij_commit/localInspection/java_file",
    "intellij_commit/localInspection/kotlin_file",
    "intellij_commit/localInspection/kotlin_file_DexInlineTest",
    "intellij_commit/localInspection/java_file_ContentManagerImpl",
    "intellij_sources/completion/java_file",
    "intellij_sources/completion/kotlin_file",
    "intellij_commit/completion/java_file",
    "intellij_commit/completion/kotlin_file"}

  settings := make([]detector.Settings, 0, 100)
  for _, test := range tests {
    metrics := getMetricFromTestName(test)
    for _, metric := range metrics {
      settings = append(settings, detector.Settings{
        Db:          "perfintDev",
        Table:       "idea",
        Channel:     "ij-perf-report-aggregator",
        Branch:      "master",
        Machine:     "intellij-linux-performance-aws-%",
        Test:        test,
        Metric:      metric,
        ProductLink: "intellij",
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
    return []string{"highlight"}
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
  return []string{}
}
