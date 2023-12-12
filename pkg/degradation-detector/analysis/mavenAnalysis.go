package analysis

import detector "github.com/JetBrains/ij-perf-report-aggregator/pkg/degradation-detector"

func getTests() []string {
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

func GenerateMavenSettings() []detector.PerformanceSettings {
  settings := make([]detector.PerformanceSettings, 0, 1000)
  settings = append(settings, generateMavenSettingsOnInstaller()...)
  settings = append(settings, generateMavenSettingsOnFastInstaller()...)
  return settings
}

func generateMavenSettingsOnInstaller() []detector.PerformanceSettings {
  tests := make([]string, 0, len(getTests()))
  for _, test := range getTests() {
    tests = append(tests, test+"/measureStartup")
  }
  metrics := []string{"maven.sync.duration",
    "maven.projects.processor.resolving.task",
    "maven.projects.processor.reading.task",
    "maven.import.stats.plugins.resolving.task",
    "maven.import.stats.applying.model.task",
    "maven.import.stats.sync.project.task",
    "maven.import.after.import.configuration",
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

    // Legacy import
    "legacy_import.create_modules.duration_ms",
    "legacy_import.delete_obsolete.duration_ms",
    "legacy_import.duration_ms",
    "legacy_import.importers.duration_ms"}

  settings := make([]detector.PerformanceSettings, 0, 200)
  for _, test := range tests {
    for _, metric := range metrics {
      settings = append(settings, detector.PerformanceSettings{
        Db:      "perfint",
        Table:   "idea",
        Machine: "intellij-linux-hw-hetzner%",
        Project: test,
        Metric:  metric,
        Branch:  "master",
        SlackSettings: detector.SlackSettings{
          Channel:     "build-tools-perf-tests-notifications",
          ProductLink: "intellij",
        },
      })
    }

  }
  return settings
}

func generateMavenSettingsOnFastInstaller() []detector.PerformanceSettings {
  tests := make([]string, 0, len(getTests()))
  for _, test := range getTests() {
    tests = append(tests, test+"/fastInstaller")
  }
  metrics := []string{"maven.sync.duration",
    "maven.projects.processor.resolving.task",
    "maven.projects.processor.reading.task",
    "maven.import.stats.plugins.resolving.task",
    "maven.import.stats.applying.model.task",
    "maven.import.stats.sync.project.task",
    "maven.import.after.import.configuration",
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

    // Legacy import
    "legacy_import.create_modules.duration_ms",
    "legacy_import.delete_obsolete.duration_ms",
    "legacy_import.duration_ms",
    "legacy_import.importers.duration_ms"}

  settings := make([]detector.PerformanceSettings, 0, 200)
  for _, test := range tests {
    for _, metric := range metrics {
      settings = append(settings, detector.PerformanceSettings{
        Db:      "perfintDev",
        Table:   "idea",
        Machine: "intellij-linux-hw-hetzner%",
        Project: test,
        Metric:  metric,
        Branch:  "master",
        SlackSettings: detector.SlackSettings{
          Channel:     "build-tools-perf-tests-notifications",
          ProductLink: "intellij",
        },
        AnalysisSettings: detector.AnalysisSettings{
          MinimumSegmentLength: 10,
        },
      })
    }

  }
  return settings
}
