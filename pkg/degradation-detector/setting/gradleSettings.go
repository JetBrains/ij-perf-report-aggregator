package setting

import (
	"slices"

	detector "github.com/JetBrains/ij-perf-report-aggregator/pkg/degradation-detector"
)

func getGradleTests() []string {
	return []string{
		"grazie-platform-project-import-gradle/",
		"project-import-gradle-monolith-51-modules-4000-dependencies-2000000-files/",
		"project-import-gradle-micronaut/",
		"project-import-gradle-hibernate-orm/",
		"project-import-gradle-cas/",
		"project-reimport-gradle-cas/",
		"project-import-from-cache-gradle-cas/",
		"project-import-gradle-1000-modules/",
		"project-import-gradle-1000-modules-limited-ram/",
		"project-import-gradle-5000-modules/",
		"project-import-gradle-android-extra-large/",
		"project-import-android-500-modules/",
		"project-reimport-space/",
		"project-import-space/",
		"project-import-open-telemetry/",
	}
}

func getGradleMetrics() []string {
	return []string{
		// total sync time
		"ExternalSystemSyncProjectTask",
		// full time of the sink operation, with all our overhead for preparation
		"GradleExecution",
		// work inside Gradle connection, operations that are performed inside connection
		"GradleConnection",
		// resolving models from daemon
		"GradleCall",
		// processing the data we received from Gradle
		"ExternalSystemSyncResultProcessing",
		// work of data services
		"ProjectDataServices",
		// project resolve
		"GradleProjectResolverDataProcessing",
		// total sync time from fus
		"fus_gradle.sync",
	}
}

func GenerateGradleSettings() []detector.PerformanceSettings {
	return slices.Concat(
		generateGradleSettingsOnInstaller(),
		generateGradleSettingsOnFastInstaller(),
	)
}

func generateGradleSettingsOnInstaller() []detector.PerformanceSettings {
	tests := make([]string, 0, len(getGradleTests()))
	for _, test := range getGradleTests() {
		tests = append(tests, test+"measureStartup")
	}
	metrics := getGradleMetrics()

	settings := make([]detector.PerformanceSettings, 0, 200)
	for _, test := range tests {
		for _, metric := range metrics {
			settings = append(settings, detector.PerformanceSettings{
				Db:      "perfint",
				Table:   "idea",
				Project: test,
				BaseSettings: detector.BaseSettings{
					Machine: "intellij-linux-hw-hetzner%",
					Metric:  metric,
					Branch:  "master",
					SlackSettings: detector.SlackSettings{
						Channel:     "gradle-perf-tests-notifications",
						ProductLink: "intellij",
					},
				},
			})
		}
	}
	return settings
}

func generateGradleSettingsOnFastInstaller() []detector.PerformanceSettings {
	tests := make([]string, 0, len(getGradleTests()))
	for _, test := range getGradleTests() {
		tests = append(tests, test+"fastInstaller")
	}
	metrics := getGradleMetrics()

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
						Channel:     "gradle-perf-tests-notifications",
						ProductLink: "intellij",
					},
					AnalysisSettings: detector.AnalysisSettings{
						MinimumSegmentLength: 10,
					},
				},
			})
		}
	}
	return settings
}
