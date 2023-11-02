package main

import "strings"

func generateIdeaAnalysisSettings() []AnalysisSettings {
  tests := []string{"intellij_sources/vfsRefresh/default", "intellij_sources/vfsRefresh/with-1-thread(s)", "intellij_sources/vfsRefresh/git-status",
    "community/rebuild", "intellij_sources/rebuild", "grails/rebuild", "java/rebuild", "spring_boot/rebuild", "java/inspection", "grails/inspection",
    "spring_boot_maven/inspection", "spring_boot/inspection", "kotlin/inspection", "kotlin_coroutines/inspection",
    "intellij_sources/localInspection/java_file", "intellij_sources/localInspection/kotlin_file", "kotlin/localInspection",
    "kotlin_coroutines/localInspection", "intellij_sources/localInspection/java_file", "intellij_sources/localInspection/kotlin_file",
    "kotlin/localInspection", "kotlin_coroutines/localInspection", "community/completion/kotlin_file", "grails/completion/groovy_file",
    "grails/completion/java_file", "kotlin_petclinic/debug", "grails/showIntentions/Find cause", "kotlin/showIntention/Import", "spring_boot/showIntentions",
    "intellij_sources/showFileHistory/EditorImpl", "intellij_sources/expandProjectMenu", "intellij_sources/expandMainMenu", "intellij_sources/expandEditorMenu",
    "kotlin/highlight", "kotlin_coroutines/highlight", "intellij_sources/FileStructureDialog/java_file", "intellij_sources/FileStructureDialog/kotlin_file",
    "intellij_sources/createJavaClass", "intellij_sources/createKotlinClass", "typingInJavaFile_16Threads/typing", "typingInJavaFile_4Threads/typing", "typingInKotlinFile_16Threads/typing",
    "typingInKotlinFile_4Threads/typing",
  }
  settings := make([]AnalysisSettings, 0, 100)
  for _, test := range tests {
    metrics := getMetricFromTestName(test)
    for _, metric := range metrics {
      settings = append(settings, AnalysisSettings{
        db:      "perfint",
        table:   "idea",
        channel: "ij-perf-report-aggregator",
        branch:  "master",
        machine: "intellij-linux-performance-aws-%",
        test:    test,
        metric:  metric,
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
    return []string{"test#average_awt_delay", "showQuickFixes", "AWTEventQueue.dispatchTimeTotal"}
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
  return []string{}
}
