package analysis

import detector "github.com/JetBrains/ij-perf-report-aggregator/pkg/degradation-detector"

func GenerateWorkspaceSettings() []detector.Settings {
  tests := []string{"project-import-jps-kotlin-50_000-modules/measureStartup"}
  metrics := []string{"project.opening", "jps.apply.loaded.storage.ms", "jps.load.project.to.empty.storage.ms", "jps.project.serializers.load.ms"}
  settings := make([]detector.Settings, 0, 10)
  for _, test := range tests {
    for _, metric := range metrics {
      settings = append(settings, detector.Settings{
        Db:          "perfint",
        Table:       "idea",
        Branch:      "master",
        Channel:     "ij-workspace-model-degradations",
        Test:        test,
        Machine:     "intellij-linux-hw-hetzner%",
        Metric:      metric,
        ProductLink: "intellij",
      })
    }

  }
  return settings
}
