package setting

import (
	"net/http"
	"slices"
	"strings"

	detector "github.com/JetBrains/ij-perf-report-aggregator/pkg/degradation-detector"
)

func GenerateJavaSettings(backendUrl string, client *http.Client) []detector.PerformanceSettings {
	return slices.Concat(
		generateJavaDevAnalysisSettings(backendUrl, client),
	)
}

func generateJavaDevAnalysisSettings(backendUrl string, client *http.Client) []detector.PerformanceSettings {
	tests := []string{
		"grails/%", "java/%", "spring_boot/%", "spring_boot_maven/%", "intellij_commit/%", "hadoop_commit/%",
	}

	baseSettings := detector.PerformanceSettings{
		Db:    "perfintDev",
		Table: "idea",
		BaseSettings: detector.BaseSettings{
			Branch:  "master",
			Machine: "intellij-linux-performance-aws-%",
		},
	}
	testsExpanded := detector.ExpandTestsByPattern(backendUrl, client, tests, baseSettings)
	testsWithoutIndexingScanning := filterIndexingScanningTests(testsExpanded)
	settings := make([]detector.PerformanceSettings, 0, 100)
	machines := []string{"intellij-linux-performance-aws-%", "intellij-windows-performance-%"}
	for _, machine := range machines {
		for _, test := range testsWithoutIndexingScanning {
			metrics := getJavaMetricsFromTestsNames(test)
			for _, metric := range metrics {
				settings = append(settings, detector.PerformanceSettings{
					Db:      baseSettings.Db,
					Table:   baseSettings.Table,
					Project: test,
					BaseSettings: detector.BaseSettings{
						Branch:  baseSettings.Branch,
						Machine: machine,
						Metric:  metric,
						SlackSettings: detector.SlackSettings{
							Channel:     "idea-java-alerts",
							ProductLink: "intellij",
						},
						AnalysisSettings: detector.AnalysisSettings{MinimumSegmentLength: 8},
					},
				})
			}
		}
	}
	return settings
}

func getJavaMetricsFromTestsNames(test string) []string {
	if strings.Contains(test, "/rebuild") {
		return []string{"build_compilation_duration"}
	}
	if strings.Contains(test, "/inspection") {
		return []string{"globalInspections"}
	}
	if strings.Contains(test, "/showIntentions") {
		return []string{"Test#average_awt_delay", "showQuickFixes"}
	}
	if strings.Contains(test, "/localInspection/java_file") {
		return []string{"localInspections", "firstCodeAnalysis"}
	}
	if strings.Contains(test, "/completion") {
		return []string{"completion"}
	}
	if strings.Contains(test, "/showIntentions") {
		return []string{"Test#average_awt_delay", "showQuickFixes"}
	}
	if strings.Contains(test, "/createJavaClass") {
		return []string{"createJavaFile"}
	}
	if strings.Contains(test, "/rename-method") {
		return []string{"performInlineRename"}
	}
	if strings.Contains(test, "/rename-class") {
		return []string{"performInlineRename"}
	}
	if strings.Contains(test, "/change-signature") {
		return []string{"changeJavaSignature: add parameter"}
	}
	if strings.Contains(test, "/move-class") {
		return []string{"moveClassToPackage"}
	}
	if strings.Contains(test, "/inline-method") {
		return []string{"inlineJavaMethod"}
	}
	if strings.Contains(test, "/rename-package") {
		return []string{"renameDirectoryAsPackage"}
	}

	// for the future cases
	if strings.Contains(test, "/findUsages") {
		return []string{"findUsages", "fus_find_usages_all", "fus_find_usages_first"}
	}
	return []string{}
}
