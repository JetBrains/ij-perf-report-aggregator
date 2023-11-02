package main

func generateWorkspaceAnalysisSettings() []AnalysisSettings {
  tests := []string{"project-import-jps-kotlin-50_000-modules/measureStartup"}
  metrics := []string{"project.opening", "jps.apply.loaded.storage.ms", "jps.load.project.to.empty.storage.ms", "jps.project.serializers.load.ms"}
  ideaSettings := make([]AnalysisSettings, 0, 10)
  for _, test := range tests {
    for _, metric := range metrics {
      ideaSettings = append(ideaSettings, AnalysisSettings{
        db:      "perfint",
        table:   "idea",
        channel: "ij-workspace-model-degradations",
        test:    test,
        metric:  metric,
      })
    }

  }
  return ideaSettings
}
