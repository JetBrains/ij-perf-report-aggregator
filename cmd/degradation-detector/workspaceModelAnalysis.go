package main

func generateWorkspaceAnalysisSettings() []AnalysisSettings {
  tests := []string{"project-import-jps-kotlin-50_000-modules/measureStartup"}
  metrics := []string{"project.opening", "jps.apply.loaded.storage.ms", "jps.load.project.to.empty.storage.ms", "jps.project.serializers.load.ms"}
  settings := make([]AnalysisSettings, 0, 10)
  for _, test := range tests {
    for _, metric := range metrics {
      settings = append(settings, AnalysisSettings{
        db:          "perfint",
        table:       "idea",
        branch:      "master",
        channel:     "ij-workspace-model-degradations",
        test:        test,
        machine:     "intellij-linux-hw-hetzner%",
        metric:      metric,
        productLink: "intellij",
      })
    }

  }
  return settings
}
