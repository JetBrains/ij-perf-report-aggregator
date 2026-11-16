package setting

import detector "github.com/JetBrains/ij-perf-report-aggregator/pkg/degradation-detector"

func GenerateTODOSettings() []detector.PerformanceSettings {
	testMetrics := []testMetricDef{
		{test: []string{"todo_many_files/todo/open", "todo_many_items/todo/open"}, metric: []string{"todo.open.treeReady"}},
		{test: []string{"todo_many_files/todo/edit", "todo_many_items/todo/edit"}, metric: []string{"todo.edit.add.treeReady"}},
		{test: []string{"todo_many_files/todo/currentFile", "todo_many_items/todo/currentFile"}, metric: []string{"todo.currentFile.switch.items500"}},
		{test: []string{"todo_many_files/todo/scope", "todo_many_items/todo/scope"}, metric: []string{"todo.scope.narrow.treeReady"}},
		{test: []string{"todo_many_files/todo/filter", "todo_many_items/todo/filter"}, metric: []string{"todo.filter.apply.treeReady"}},
		{test: []string{"todo_many_files/todo/preview", "todo_many_items/todo/preview"}, metric: []string{"todo.preview.enable", "todo.navigate.nextFile"}},
	}

	settings := make([]detector.PerformanceSettings, 0, 28)
	for _, mode := range []string{"", "split"} {
		for _, testMetric := range testMetrics {
			for _, test := range testMetric.test {
				for _, metric := range testMetric.metric {
					settings = append(settings, detector.PerformanceSettings{
						BaseSettings: detector.BaseSettings{
							SlackSettings: detector.SlackSettings{
								Channel:     "rd-qa-perf-test-alerts",
								ProductLink: "intellij",
							},
							AnalysisSettings: detector.AnalysisSettings{
								ReportType:                detector.DegradationEvent,
								MinimumSegmentLength:      8,
								MedianDifferenceThreshold: 1,
								EffectSizeThreshold:       2,
							},
							Metric:  metric,
							Branch:  "master",
							Machine: "intellij-linux-performance-aws-%",
						},
						Db:      "perfintDev",
						Table:   "idea",
						Project: test,
						Mode:    mode,
					})
				}
			}
		}
	}
	return settings
}
