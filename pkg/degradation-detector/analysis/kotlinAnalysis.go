package analysis

import detector "github.com/JetBrains/ij-perf-report-aggregator/pkg/degradation-detector"

func GenerateKotlinSettings() []detector.PerformanceSettings {
  testNames := []string{"intellij_commit/completion/empty_place_with_library_cache",
    "intellij_commit/completion/after_parameter_with_library_cache",
    "intellij_commit/completion/empty_place_typing_with_library_cache",
    "intellij_commit/completion/KotlinHighLevelFunctionParameterInfoHandler_emptyPlace_updateUIOrFail_typing_with_library_cache",
    "intellij_commit/completion/KotlinHighLevelFunctionParameterInfoHandler_emptyPlace_updateUIOrFail_with_library_cache",
    "intellij_commit/completion/KtOCSwiftChangeSignatureTest_emptyPlace_changeReturnType_typing_with_library_cache",
    "intellij_commit/completion/KtOCSwiftChangeSignatureTest_emptyPlace_changeReturnType_with_library_cache",
    "intellij_commit/completion/IdeMenuBar_emptyPlace_sout_typing_with_library_cache",
    "intellij_commit/completion/TestModelParser_emptyPlace_if_typing_with_library_cache",
    "intellij_commit/completion/AndroidModuleSystem_emptyPlace_get_typing_with_library_cache",
    "kotlin_lang/completion/after_parameter_with_library_cache",
    "kotlin_lang/completion/empty_place_with_library_cache",
    "kotlin_lang/completion/empty_place_typing_with_library_cache",
    "kotlin_language_server/completion/Completions_emptyPlace_completions_typing_with_library_cache",
    "kotlin_language_server/completion/QuickFixesTest_emptyPlace_completions_typing_with_library_cache",
    "kotlin_language_server/completion/Completions_emptyPlace_completions_with_library_cache",
    "kotlin_language_server/completion/QuickFixesTest_emptyPlace_completions_with_library_cache",
  }
  tests := generateKotlinTests(testNames)
  settings := make([]detector.PerformanceSettings, 0, 50)
  metrics := []string{"completion#mean_value"}
  dbs := []string{"perfint", "perfintDev"}
  branches := []string{"master", "kt-master"}
  for _, test := range tests {
    for _, metric := range metrics {
      for _, branch := range branches {
        for _, db := range dbs {
          settings = append(settings, detector.PerformanceSettings{
            Db:          db,
            Table:       "kotlin",
            Channel:     "kotlin-plugin-perf-tests",
            Machine:     "intellij-linux-hw-hetzner%",
            Project:     test,
            Metric:      metric,
            Branch:      branch,
            ProductLink: "kotlin",
          })
        }
      }
    }

  }
  return settings
}

func generateKotlinTests(tests []string) []string {
  k1K2tests := make([]string, 0, len(tests)*2)
  for _, test := range tests {
    k1K2tests = append(k1K2tests, test+"_k1")
    k1K2tests = append(k1K2tests, test+"_k2")
  }
  return k1K2tests
}
