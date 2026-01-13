package setting

import detector "github.com/JetBrains/ij-perf-report-aggregator/pkg/degradation-detector"

type testMetricDef struct {
	test   []string
	metric []string
}

func GenerateVCSSettings() []detector.PerformanceSettings {
	testMetrics := []testMetricDef{
		{test: []string{"intellij_clone_specific_commit/gitLogIndexing", "intellij_sources/gitLogIndexing"}, metric: []string{"vcs-log-indexing"}},
		{test: []string{
			"intellij_sources/EditorImpl-phm",
			"intellij_sources/EditorImpl-noindex",
			"intellij_sources/showFileHistory/EditorImpl",
			"intellij_sources/showFileHistory/EditorImpl-instant-git",
			"intellij_commit/showFileHistory/EditorImpl",
			"intellij_commit/showFileHistory/EditorImpl-instant-git",
		}, metric: []string{"showFileHistory", "showFirstPack"}},
		{test: []string{"intellij_sources/filterVcsLogTab-phm", "intellij_sources/filterVcsLogTab-noindex"}, metric: []string{"vcs-log-filtering"}},
		{test: []string{"intellij_sources/filterVcsLogTab-path-phm", "intellij_sources/filterVcsLogTab-path-noindex"}, metric: []string{"vcs-log-filtering"}},
		{test: []string{"intellij_sources/filterVcsLogTab-date-phm", "intellij_sources/filterVcsLogTab-date-noindex"}, metric: []string{"vcs-log-filtering"}},
		{test: []string{"intellij_sources/git-commit"}, metric: []string{"fus_vcs_commit_duration", "vcs-log-refreshing"}},
		{test: []string{"intellij_sources/git-branch-widget", "vcs_100k_branches/git-branch-widget"}, metric: []string{"gitShowBranchWidget"}},
		{test: []string{"intellij_sources/vcs-annotate-instant-git", "intellij_sources/vcs-annotate"}, metric: []string{"showFileAnnotation", "git-open-annotation"}},
	}

	machines := []string{"intellij-linux-performance-aws-%"} // uncomment latter to cover all os
	// "intellij-windows-performance-aws-%",
	// "intellij-macos-perf-eqx-%",

	settings := make([]detector.PerformanceSettings, 0, 100)
	for _, testMetric := range testMetrics {
		for _, test := range testMetric.test {
			for _, metric := range testMetric.metric {
				for _, machine := range machines {
					settings = append(settings, detector.PerformanceSettings{
						Db:      "perfintDev",
						Table:   "idea",
						Project: test,
						BaseSettings: detector.BaseSettings{
							Machine: machine,
							Metric:  metric,
							Branch:  "master",
							SlackSettings: detector.SlackSettings{
								Channel:     "vcs-perf-tests",
								ProductLink: "intellij",
							},
							AnalysisSettings: detector.AnalysisSettings{
								MinimumSegmentLength:      7,
								MedianDifferenceThreshold: 10,
							},
						},
					})
				}
			}
		}
	}
	return settings
}
