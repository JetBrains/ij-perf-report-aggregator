package setting

import (
	"slices"

	detector "github.com/JetBrains/ij-perf-report-aggregator/pkg/degradation-detector"
)

func getMavenTests() []string {
	return []string{
		"project-import-maven-quarkus",
		"project-reimport-maven-quarkus",
		"project-import-from-cache-maven-quarkus",
		"project-import-maven-1000-modules",
		"project-import-maven-5000-modules",
		"project-import-maven-keycloak",
		"project-import-maven-javaee7",
		"project-import-maven-javaee8",
		"project-import-maven-jersey",
		"project-import-maven-flink",
		"project-import-maven-drill",
		"project-import-maven-azure-sdk-java",
		"project-import-maven-hive",
		"project-import-maven-quarkus-to-legacy-model",
		"project-import-maven-1000-modules-to-legacy-model",
	}
}

func getMavenMetrics() []string {
	return []string{
		// Main flow
		"maven.sync.duration",
		"maven.projects.processor.resolving.task",
		"maven.projects.processor.reading.task",
		"maven.import.stats.plugins.resolving.task",
		"maven.import.stats.applying.model.task",
		"maven.import.stats.sync.project.task",
		"maven.import.after.import.configuration",
		"maven.import.configure.post.process",
		"maven.project.importer.base.refreshing.files.task",
		"maven.project.importer.post.importing.task.marker",
		"post_import_tasks_run.total_duration_ms",

		// Workspace commit
		"workspace_commit.attempts",
		"workspace_commit.duration_in_background_ms",
		"workspace_commit.duration_in_write_action_ms",
		"workspace_commit.duration_of_workspace_update_call_ms",

		// Workspace import
		"workspace_import.commit.duration_ms",
		"workspace_import.duration_ms",
		"workspace_import.legacy_importers.duration_ms",
		"workspace_import.legacy_importers.stats.duration_of_bridges_creation_ms",
		"workspace_import.legacy_importers.stats.duration_of_bridges_commit_ms",
		"workspace_import.populate.duration_ms",
	}
}

func GenerateMavenSettings() []detector.PerformanceSettings {
	return slices.Concat(
		generateMavenSettingsOnFastInstaller(),
	)
}

func generateMavenSettingsOnFastInstaller() []detector.PerformanceSettings {
	tests := make([]string, 0, len(getMavenTests()))
	for _, test := range getMavenTests() {
		tests = append(tests, test+"/fastInstaller")
	}
	metrics := getMavenMetrics()

	settings := make([]detector.PerformanceSettings, 0, 200)
	for _, test := range tests {
		for _, metric := range metrics {
			settings = append(settings, detector.PerformanceSettings{
				Db:      "perfintDev",
				Table:   "idea",
				Project: test,
				BaseSettings: detector.BaseSettings{
					Machine: "intellij-linux-hw-hetzner%",
					Metric:  metric,
					Branch:  "master",
					SlackSettings: detector.SlackSettings{
						Channel:     "maven-perf-tests-notifications",
						ProductLink: "intellij",
					},
					AnalysisSettings: detector.AnalysisSettings{
						MinimumSegmentLength: 10,
						ReportType:           detector.DegradationEvent,
					},
				},
			})
		}
	}
	return settings
}
