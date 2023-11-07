package analysis

import detector "github.com/JetBrains/ij-perf-report-aggregator/pkg/degradation-detector"

func GenerateMavenSettings() []detector.Settings {
  tests := []string{
    "project-import-maven-quarkus/measureStartup",
    "project-reimport-maven-quarkus/measureStartup",
    "project-import-from-cache-maven-quarkus/measureStartup",
    "project-import-maven-1000-modules/measureStartup",
    "project-import-maven-5000-modules/measureStartup",
    "project-import-maven-keycloak/measureStartup",
    "project-import-maven-javaee7/measureStartup",
    "project-import-maven-javaee8/measureStartup",
    "project-import-maven-jersey/measureStartup",
    "project-import-maven-flink/measureStartup",
    "project-import-maven-drill/measureStartup",
    "project-import-maven-azure-sdk-java/measureStartup",
    "project-import-maven-hive/measureStartup",
    "project-import-maven-quarkus-to-legacy-model/measureStartup",
    "project-import-maven-1000-modules-to-legacy-model/measureStartup",
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

  settings := make([]detector.Settings, 0, 200)
  for _, test := range tests {
    for _, metric := range metrics {
      settings = append(settings, detector.Settings{
        Db:          "perfint",
        Table:       "idea",
        Channel:     "build-tools-perf-tests-notifications",
        Machine:     "intellij-linux-hw-hetzner%",
        Test:        test,
        Metric:      metric,
        Branch:      "master",
        ProductLink: "intellij",
      })
    }

  }
  return settings
}
