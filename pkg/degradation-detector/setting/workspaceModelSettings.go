package setting

import detector "github.com/JetBrains/ij-perf-report-aggregator/pkg/degradation-detector"

func GenerateWorkspaceSettings() []detector.PerformanceSettings {
	tests := []string{"project-import-jps-kotlin-50_000-modules/measureStartup"}
	metrics := []string{"project.opening", "jps.apply.loaded.storage.ms", "jps.load.project.to.empty.storage.ms", "jps.project.serializers.load.ms"}
	settings := make([]detector.PerformanceSettings, 0, 10)
	for _, test := range tests {
		for _, metric := range metrics {
			settings = append(settings, detector.PerformanceSettings{
				Db:      "perfint",
				Table:   "idea",
				Branch:  "master",
				Project: test,
				Machine: "intellij-linux-hw-hetzner%",
				Metric:  metric,
				SlackSettings: detector.SlackSettings{
					Channel:     "ij-workspace-model-perf-tests",
					ProductLink: "intellij",
				},
			})
		}
	}
	return settings
}
