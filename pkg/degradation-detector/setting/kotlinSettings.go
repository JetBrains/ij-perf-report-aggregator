package setting

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"

	detector "github.com/JetBrains/ij-perf-report-aggregator/pkg/degradation-detector"
)

func extractStrings(data any) []string {
	var result []string
	switch v := data.(type) {
	case string:
		result = append(result, v)
	case []any:
		for _, item := range v {
			result = append(result, extractStrings(item)...)
		}
	case map[string]any:
		for _, item := range v {
			result = append(result, extractStrings(item)...)
		}
	}
	return result
}

func koDataPath() string {
	if env := os.Getenv("KO_DATA_PATH"); env != "" {
		return env
	}
	_, thisFile, _, _ := runtime.Caller(0)
	repoRoot := filepath.Join(filepath.Dir(thisFile), "..", "..", "..")
	return filepath.Join(repoRoot, "cmd", "degradation-analyzer", "kodata")
}

func getKotlinProjects() ([]string, error) {
	var result any
	kotlinProjectsFilePath := filepath.Join(koDataPath(), "projects", "kotlin_projects.json")
	content, err := os.ReadFile(kotlinProjectsFilePath)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(content, &result); err != nil {
		return nil, fmt.Errorf("failed to parse %s: %w", kotlinProjectsFilePath, err)
	}
	names := extractStrings(result)
	if len(names) == 0 {
		return nil, fmt.Errorf("no project found in %s", kotlinProjectsFilePath)
	}
	return names, nil
}

func GenerateKotlinSettings() []detector.PerformanceSettings {
	testNames, err := getKotlinProjects()
	if err != nil {
		slog.Error("cannot load kotlin projects, no kotlin test will be analyzed", "error", err)
		return nil
	}
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
		"codeTyping#mean_value",
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
		"codeTyping#mean_value":                          "codeTyping",
	}
	settings := make([]detector.PerformanceSettings, 0, len(testNames)*len(metrics)*2)

	for _, test := range tests {
		daysToCheck := -1
		for _, metric := range metrics {
			alias := getAlias(metric, aliases)
			settings = append(settings, detector.PerformanceSettings{
				Db:                 "perfintDev",
				Table:              "kotlin",
				Project:            test,
				MetricAlias:        alias,
				Machine:            "intellij-linux-hw-hetzner%",
				Metric:             metric,
				Branch:             "master",
				Channel:            "kotlin-plugin-perf-tests",
				ProductLink:        "kotlin",
				ReportType:         detector.DegradationEvent,
				DaysToCheckMissing: daysToCheck,
			})
		}
	}
	for _, test := range tests {
		for _, metric := range metrics {
			alias := getAlias(metric, aliases)
			settings = append(settings, detector.PerformanceSettings{
				Db:          "perfintDev",
				Table:       "kotlin",
				Project:     test,
				MetricAlias: alias,
				Machine:     "intellij-linux-hw-hetzner%",
				Metric:      metric,
				Branch:      "kt-master",
				Channel:     "kotlin-plugin-perf-tests-kt-master",
				ProductLink: "kotlin",
				ReportType:  detector.DegradationEvent,
			})
		}
	}
	return settings
}

func generateKotlinTests(tests []string) []string {
	k2tests := make([]string, 0, len(tests))
	for _, test := range tests {
		k2tests = append(k2tests, test+"_k2")
	}
	return k2tests
}

func getAlias(metric string, aliases map[string]string) string {
	alias, ok := aliases[metric]
	if !ok {
		alias = metric
	}
	return alias
}
