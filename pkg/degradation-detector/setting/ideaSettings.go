package setting

import (
	"net/http"
	"slices"
	"strings"

	detector "github.com/JetBrains/ij-perf-report-aggregator/pkg/degradation-detector"
)

func GenerateIdeaSettings(backendUrl string, client *http.Client) []detector.PerformanceSettings {
	return slices.Concat(
		generateIdeaDevAnalysisSettings(backendUrl, client),
	)
}

func generateIdeaDevAnalysisSettings(backendUrl string, client *http.Client) []detector.PerformanceSettings {
	tests := []string{
		"intellij_commit/%", "grails/%", "java/%", "spring_boot/%",
		"spring_boot_maven/%", "spring_boot/%", "kotlin/%", "kotlin_coroutines/%",
		"kotlin_petclinic/%", "community/%", "empty_project/%", "keycloak_release_20/%",
		"space/%", "toolbox_enterprise/%", "train-ticket/%", "popups-performance-test/%", "%terminal-completion-all%",
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
			metrics := getMetricFromTestName(test)
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
							Channel:     "ij-perf-report-aggregator",
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

func filterIndexingScanningTests(input []string) []string {
	var result []string
	for _, str := range input {
		if !strings.Contains(str, "indexing") && !strings.Contains(str, "scanning") {
			result = append(result, str)
		}
	}
	return result
}

func getMetricFromTestName(test string) []string {
	if strings.Contains(test, "/vfsRefresh") {
		return []string{"vfs_initial_refresh"}
	}
	if strings.Contains(test, "/rebuild") {
		return []string{"build_compilation_duration"}
	}
	if strings.Contains(test, "/inspection") || strings.Contains(test, "/globalInspection") {
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
	if strings.Contains(test, "/scrollEditor") {
		return []string{"scrollEditor#average_awt_delay", "scrollEditor#max_awt_delay", "scrollEditor#average_cpu_load", "scrollEditor#max_cpu_load"}
	}
	return []string{}
}
