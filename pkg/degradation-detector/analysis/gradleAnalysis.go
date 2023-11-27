package analysis

import detector "github.com/JetBrains/ij-perf-report-aggregator/pkg/degradation-detector"

func GenerateGradleSettings() []detector.PerformanceSettings {
  tests := []string{
    "grazie-platform-project-import-gradle/measureStartup",
    "project-import-gradle-monolith-51-modules-4000-dependencies-2000000-files/measureStartup",
    "project-import-gradle-micronaut/measureStartup",
    "project-import-gradle-hibernate-orm/measureStartup",
    "project-import-gradle-cas/measureStartup",
    "project-reimport-gradle-cas/measureStartup",
    "project-import-from-cache-gradle-cas/measureStartup",
    "project-import-gradle-1000-modules/measureStartup",
    "project-import-gradle-1000-modules-limited-ram/measureStartup",
    "project-import-gradle-5000-modules/measureStartup",
    "project-import-gradle-android-extra-large/measureStartup",
    "project-reimport-space/measureStartup",
    "project-import-space/measureStartup",
    "project-import-open-telemetry/measureStartup",
  }
  settings := make([]detector.PerformanceSettings, 0, 100)
  metrics := []string{"gradle.sync.duration",
    "GRADLE_CALL",
    "PROJECT_RESOLVERS",
    "DATA_SERVICES",
    "WORKSPACE_MODEL_APPLY",
    "fus_gradle.sync"}
  for _, test := range tests {
    for _, metric := range metrics {
      settings = append(settings, detector.PerformanceSettings{
        Db:          "perfint",
        Table:       "idea",
        Channel:     "build-tools-perf-tests-notifications",
        Machine:     "intellij-linux-hw-hetzner%",
        Project:     test,
        Metric:      metric,
        Branch:      "master",
        ProductLink: "intellij",
      })
    }

  }
  return settings
}
