package setting

import (
	"encoding/json"
	detector "github.com/JetBrains/ij-perf-report-aggregator/pkg/degradation-detector"
	"os"
	"path/filepath"
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
		"semanticHighlighting#mean_value", "localInspections#mean_value",
		"completion#firstElementShown#mean_value", "evaluateExpression#mean_value",
		"performInlineRename#mean_value", "startInlineRename#mean_value",
		"prepareForRename#mean_value", "execute_editor_optimizeimports",
		"localInspections_cold#mean_value", "localInspections_hot#mean_value",
		"execute_editor_gotodeclaration_cold#mean_value", "execute_editor_gotodeclaration_hot#mean_value",
		"convertJavaToKotlin", "moveFiles#mean_value", "moveFiles_back#mean_value"}
	aliases := map[string]string{
		"completion#mean_value":                          "completion",
		"completion#firstElementShown#mean_value":        "completion",
		"findUsages#mean_value":                          "findUsages",
		"semanticHighlighting#mean_value":                "highlighting",
		"localInspections#mean_value":                    "highlighting",
		"performInlineRename#mean_value":                 "rename",
		"prepareForRename#mean_value":                    "rename",
		"startInlineRename#mean_value":                   "rename",
		"execute_editor_optimizeimports":                 "optimizeimports",
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
		for _, metric := range metrics {
			alias := getAlias(metric, aliases)
			settings = append(settings, detector.PerformanceSettings{
				Db:          "perfintDev",
				Table:       "kotlin",
				Machine:     "intellij-linux-hw-hetzner%",
				Project:     test,
				Metric:      metric,
				Branch:      "master",
				MetricAlias: alias,
				SlackSettings: detector.SlackSettings{
					Channel:     "kotlin-plugin-perf-tests",
					ProductLink: "kotlin",
				},
				AnalysisSettings: detector.AnalysisSettings{
					ReportType: detector.DegradationEvent,
				},
			})
		}
	}
	for _, test := range tests {
		for _, metric := range metrics {
			settings = append(settings, detector.PerformanceSettings{
				Db:          "perfintDev",
				Table:       "kotlin",
				Machine:     "intellij-linux-hw-hetzner%",
				Project:     test,
				Metric:      metric,
				Branch:      "master",
				MetricAlias: "perf",
				SlackSettings: detector.SlackSettings{
					Channel:     "kotlin-plugin-perf-test-merged",
					ProductLink: "kotlin",
				},
				AnalysisSettings: detector.AnalysisSettings{
					ReportType: detector.DegradationEvent,
				},
			})
		}
	}
	for _, test := range tests {
		for _, metric := range metrics {
			alias := getAlias(metric, aliases)
			settings = append(settings, detector.PerformanceSettings{
				Db:          "perfintDev",
				Table:       "kotlin",
				Machine:     "intellij-linux-hw-hetzner%",
				Project:     test,
				Metric:      metric,
				Branch:      "master",
				MetricAlias: alias,
				SlackSettings: detector.SlackSettings{
					Channel:     "kotlin-plugin-perf-tests-optimization",
					ProductLink: "kotlin",
				},
				AnalysisSettings: detector.AnalysisSettings{
					ReportType: detector.ImprovementEvent,
				},
			})
		}
	}
	for _, test := range tests {
		for _, metric := range metrics {
			alias := getAlias(metric, aliases)
			settings = append(settings, detector.PerformanceSettings{
				Db:          "perfintDev",
				Table:       "kotlin",
				Machine:     "intellij-linux-hw-hetzner%",
				Project:     test,
				Metric:      metric,
				Branch:      "kt-master",
				MetricAlias: alias,
				SlackSettings: detector.SlackSettings{
					Channel:     "kotlin-plugin-perf-tests-kt-master",
					ProductLink: "kotlin",
				},
				AnalysisSettings: detector.AnalysisSettings{
					ReportType: detector.DegradationEvent,
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
