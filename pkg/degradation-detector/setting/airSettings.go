package setting

import (
	"fmt"

	detector "github.com/JetBrains/ij-perf-report-aggregator/pkg/degradation-detector"
)

// The Air performance tests run once a night on intellij-linux-performance-aws agents. A stable metric varies by up to
// about 7 % from night to night, a noisy one by 10 % to 40 %. The thresholds keep a false alarm to a few percent of the
// metrics a month: a stable metric reports a step of about 20 %, a noisy one a step of about 50 %, after 5 nights.
var (
	airStable = detector.AnalysisSettings{ReportType: detector.DegradationEvent, MedianDifferenceThreshold: 15, EffectSizeThreshold: 2}
	airNoisy  = detector.AnalysisSettings{ReportType: detector.DegradationEvent, MedianDifferenceThreshold: 25, EffectSizeThreshold: 3}
)

func GenerateAirSettings() []detector.PerformanceSettings {
	type metric struct {
		name     string
		analysis detector.AnalysisSettings
	}
	longSession := []metric{
		{"open.%s.coverRemovedMs", airStable},
		{"open.%s.gestureToSettledMs", airStable},
		{"continue.%s.sendToPromptMs", airStable},
		{"answer.%s.streamEndToSettledMs", airNoisy},
		{"heap.%s.heapLive.MiB", airStable},
		{"prompt.%s.codeP95Ms", airNoisy},
		{"scroll.%s.wheelP95Ms", airNoisy},
		// A jump to the start and to 0.4 of the history shows rows far from the bottom, whose text the transcript builds
		// only near the visible area. The time runs to the frame that shows the built text, so it is noisy until the
		// nights show otherwise.
		{"scroll.%s.jump0.settledMs", airNoisy},
		{"scroll.%s.jump40.settledMs", airNoisy},
	}
	metrics := map[string][]metric{
		"air/many-agents": {
			{"manyAgents.heapLive.MiB", airStable},
			{"manyAgents.openTabLagMaxMs", airStable},
			{"manyAgents.focusedTabLagMs", airNoisy},
			{"manyAgents.edtP95Ms", airNoisy},
		},
		"air/session-list": {
			{"sessionList.heapLive.MiB", airStable},
			{"sessionList.refreshMedianMs", airStable},
			{"sessionList.firstShowMs", airNoisy},
		},
	}
	for _, size := range []string{"M", "XL", "XXL"} {
		project := "air/long-session/" + size
		for _, m := range longSession {
			metrics[project] = append(metrics[project], metric{fmt.Sprintf(m.name, size), m.analysis})
		}
	}

	settings := make([]detector.PerformanceSettings, 0, 34)
	for test, metrics := range metrics {
		for _, m := range metrics {
			settings = append(settings, detector.PerformanceSettings{
				Db:               "perfintDev",
				Table:            "air",
				Project:          test,
				Machine:          "intellij-linux-performance-aws-%",
				Metric:           m.name,
				Branch:           "master",
				Channel:          "ij-air-perf",
				ProductLink:      "air",
				AnalysisSettings: m.analysis,
			})
		}
	}
	return settings
}
