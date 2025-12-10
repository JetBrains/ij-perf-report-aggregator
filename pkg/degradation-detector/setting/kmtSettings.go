package setting

import (
	"log/slog"
	"net/http"

	detector "github.com/JetBrains/ij-perf-report-aggregator/pkg/degradation-detector"
)

func GenerateKotlinMultiplatformToolingSettings(backendUrl string, client *http.Client) []detector.PerformanceSettings {
	baseSettings := detector.PerformanceSettings{
		Db:    "perfintDev",
		Table: "kmt",
		BaseSettings: detector.BaseSettings{
			Branch:  "master",
			Machine: "cidr.performance.appcode.osx%",
		},
	}
	tests, err := detector.FetchAllTests(backendUrl, client, baseSettings)
	settings := make([]detector.PerformanceSettings, 0, 300)
	if err != nil {
		slog.Error("error while getting tests", "error", err)
		return settings
	}
	for _, test := range tests {
		metrics := []string{
			"Create KMP Run Configurations",
			"Progress: Generating Xcode files…",
			"globalInspections",
			"SourceKitDiagnosticsPass#mean_value",
			"SourceKitSemanticHighlightingPass#mean_value",
			"completion",
			"XCodeBuild",
			"IosAppStartup",
			"IosAppStartupDebug",
			"KmpIosConfigurationRun",
		}
		for _, metric := range metrics {
			settings = append(settings, detector.PerformanceSettings{
				Db:      baseSettings.Db,
				Table:   baseSettings.Table,
				Project: test,
				Mode:    "intellij-idea",
				BaseSettings: detector.BaseSettings{
					Branch:  baseSettings.Branch,
					Machine: baseSettings.Machine,
					Metric:  metric,
					SlackSettings: detector.SlackSettings{
						Channel:     "kmt-infrastructure-alerts",
						ProductLink: "kmt",
					},
				},
			})
		}
	}
	return settings
}
