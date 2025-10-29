package setting

import (
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"strings"

	detector "github.com/JetBrains/ij-perf-report-aggregator/pkg/degradation-detector"
)

func getIndexingProjects() ([]string, error) {
	var result map[string]interface{}
	env := os.Getenv("KO_DATA_PATH")
	projectsFile := "projects/idea_indexing_projects.json"
	var indexingProjectsFilePath string
	if env == "" {
		indexingProjectsFilePath = filepath.Join("cmd", "degradation-detector", "kodata", projectsFile)
	} else {
		indexingProjectsFilePath = filepath.Join(env, projectsFile)
	}
	content, err := os.ReadFile(indexingProjectsFilePath)
	if err != nil {
		return nil, err
	}
	_ = json.Unmarshal(content, &result)
	return extractStrings(result), nil
}

func GenerateIdeaIndexingSettings() []detector.PerformanceSettings {
	return slices.Concat(
		generateIdeaIndexingSettings(),
	)
}

func generateIdeaIndexingSettings() []detector.PerformanceSettings {
	tests, _ := getIndexingProjects()
	baseSettings := detector.PerformanceSettings{
		Db:    "perfintDev",
		Table: "idea",
		BaseSettings: detector.BaseSettings{
			Branch:  "master",
			Machine: "intellij-linux-performance-aws-%",
		},
	}
	settings := make([]detector.PerformanceSettings, 0, 100)
	machines := []string{"intellij-linux-performance-aws-%", "intellij-windows-performance-%"}
	for _, machine := range machines {
		for _, test := range tests {
			metrics := getIndexingMetricFromTestNameForIDEA(test)
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
							Channel:     "ij-indexes-perf-alerts",
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

func getIndexingMetricFromTestNameForIDEA(test string) []string {
	if strings.Contains(test, "/indexing") {
		return []string{"scanningTimeWithoutPauses", "indexingTimeWithoutPauses"}
	}
	if strings.Contains(test, "-scanning") {
		return []string{"scanningTimeWithoutPauses"}
	}
	return []string{}
}
