package setting

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	detector "github.com/JetBrains/ij-perf-report-aggregator/pkg/degradation-detector"
)

func extractStrings(data interface{}) []string {
	var result []string
	switch v := data.(type) {
	case string:
		result = append(result, v)
	case []interface{}:
		for _, item := range v {
			result = append(result, extractStrings(item)...)
		}
	case map[string]interface{}:
		for _, item := range v {
			result = append(result, extractStrings(item)...)
		}
	}
	return result
}

func getKotlinProjects() ([]string, error) {
	var result interface{}
	env := os.Getenv("KO_DATA_PATH")
	projectsFile := "projects/kotlin_projects.json"
	var kotlinProjectsFilePath string
	if env == "" {
		kotlinProjectsFilePath = filepath.Join("..", "..", "..", "cmd", "degradation-detector", "kodata", projectsFile)
	} else {
		kotlinProjectsFilePath = filepath.Join(env, projectsFile)
	}
	content, err := os.ReadFile(kotlinProjectsFilePath)
	if err != nil {
		return nil, err
	}
	_ = json.Unmarshal(content, &result)
	return extractStrings(result), nil
}

func GenerateKotlinSettings() []detector.PerformanceSettings {
	testNames, _ := getKotlinProjects()
	tests := generateKotlinTests(testNames)
	metrics := []string{
		"completion#mean_value", "findUsages#mean_value",
		"localInspections#mean_value",
		"completion#firstElementShown#mean_value", "evaluateExpression#mean_value",
		"performInlineRename#mean_value", "startInlineRename#mean_value",
		"prepareForRename#mean_value", "fus_refactoring_usages_searched", "execute_editor_optimizeimports", "execute_editor_optimizeimports#mean_value",
		"localInspections_cold#mean_value", "localInspections_hot#mean_value",
		"execute_editor_gotodeclaration_cold#mean_value", "execute_editor_gotodeclaration_hot#mean_value",
		"convertJavaToKotlin", "moveFiles#mean_value", "moveFiles_back#mean_value", "moveDeclarations#mean_value", "moveDeclarations_back#mean_value",
	}
	aliases := map[string]string{
		"completion#mean_value":                          "completion",
		"completion#firstElementShown#mean_value":        "completion",
		"findUsages#mean_value":                          "findUsages",
		"localInspections#mean_value":                    "highlighting",
		"performInlineRename#mean_value":                 "rename",
		"prepareForRename#mean_value":                    "rename",
		"startInlineRename#mean_value":                   "rename",
		"fus_refactoring_usages_searched":                "rename",
		"execute_editor_optimizeimports":                 "optimizeimports",
		"execute_editor_optimizeimports#mean_value":      "optimizeimports",
		"evaluateExpression#mean_value":                  "debugger",
		"execute_editor_gotodeclaration_hot#mean_value":  "gotodeclaration_hot_cache",
		"execute_editor_gotodeclaration_cold#mean_value": "gotodeclaration_cold_cache",
		"localInspections_cold#mean_value":               "highlighting_cold_cache",
		"localInspections_hot#mean_value":                "highlighting_hot_cache",
		"convertJavaToKotlin":                            "J2K",
		"moveFiles#mean_value":                           "moveFiles",
	}
	settings := make([]detector.PerformanceSettings, 0, len(testNames)*len(metrics)*2)

	for _, test := range tests {
		threshold := 10
		if strings.HasSuffix(test, "_k1") {
			threshold = 20
		}
		for _, metric := range metrics {
			alias := getAlias(metric, aliases)
			settings = append(settings, detector.PerformanceSettings{
				Db:          "perfintDev",
				Table:       "kotlin",
				Project:     test,
				MetricAlias: alias,
				BaseSettings: detector.BaseSettings{
					Machine: "intellij-linux-hw-hetzner%",
					Metric:  metric,
					Branch:  "master",
					SlackSettings: detector.SlackSettings{
						Channel:     "kotlin-plugin-perf-tests",
						ProductLink: "kotlin",
					},
					AnalysisSettings: detector.AnalysisSettings{
						ReportType:         detector.DegradationEvent,
						MedianDifferenceThreshold: threshold,
						DaysToCheckMissing: -1,
					},
				},
			})
		}
	}
	for _, test := range tests {
	  threshold := 10
  	if strings.HasSuffix(test, "_k1") {
  		threshold = 20
  	}
		for _, metric := range metrics {
			alias := getAlias(metric, aliases)
			settings = append(settings, detector.PerformanceSettings{
				Db:          "perfintDev",
				Table:       "kotlin",
				Project:     test,
				MetricAlias: alias,
				BaseSettings: detector.BaseSettings{
					Machine: "intellij-linux-hw-hetzner%",
					Metric:  metric,
					Branch:  "kt-master",
					SlackSettings: detector.SlackSettings{
						Channel:     "kotlin-plugin-perf-tests-kt-master",
						ProductLink: "kotlin",
					},
					AnalysisSettings: detector.AnalysisSettings{
						ReportType: detector.DegradationEvent,
						MedianDifferenceThreshold: threshold,
					},
				},
			})
		}
	}
	return settings
}

func generateKotlinTests(tests []string) []string {
	k1K2tests := make([]string, 0, len(tests)*2)
	for _, test := range tests {
		k1K2tests = append(k1K2tests, test+"_k1", test+"_k2")
	}
	return k1K2tests
}

func getAlias(metric string, aliases map[string]string) string {
	alias, ok := aliases[metric]
	if !ok {
		alias = metric
	}
	return alias
}
