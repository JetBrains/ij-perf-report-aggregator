package analysis

func GenerateWorkspaceSettings() []Settings {
  tests := []string{"project-import-jps-kotlin-50_000-modules/measureStartup"}
  metrics := []string{"project.opening", "jps.apply.loaded.storage.ms", "jps.load.project.to.empty.storage.ms", "jps.project.serializers.load.ms"}
  settings := make([]Settings, 0, 10)
  for _, test := range tests {
    for _, metric := range metrics {
      settings = append(settings, Settings{
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
