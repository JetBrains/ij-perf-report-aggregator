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

// airOver alerts when the latest value is over limitMs. A metric of a few frames moves by one frame from night to night,
// which is 30-80 % of its value, so a step of its median says nothing. A value over the limit is what a user notices.
func airOver(limitMs float64) detector.AnalysisSettings {
	return detector.AnalysisSettings{
		ReportType:     detector.DegradationEvent,
		AnalysisKind:   detector.ThresholdAnalysis,
		ThresholdMode:  detector.ThresholdGreaterThan,
		ThresholdValue: limitMs,
	}
}

func GenerateAirSettings() []detector.PerformanceSettings {
	type metric struct {
		name     string
		analysis detector.AnalysisSettings
	}
	longSession := []metric{
		{"open.%s.coverRemovedMs", airStable},
		{"checkpoint.%s.frontend_transcript_rows_applied.elapsedMs", airStable},
		{"continue.%s.loadAnsweredToPromptMs", airOver(300)},
		{"answer.%s.streamLagP95Ms", airNoisy},
		{"answer.%s.streamEndToSettledMs", airOver(300)},
		{"heap.%s.heapLive.MiB", airStable},
		{"prompt.%s.codeP95Ms", airOver(50)},
		{"scroll.%s.wheelP95Ms", airOver(50)},
		// A jump to the start and to 0.4 of the history lands on rows that are not built yet. A blank screen of 200 ms
		// after a drag is what a user notices.
		{"scroll.%s.jump0.settledMs", airOver(200)},
		{"scroll.%s.jump40.settledMs", airOver(200)},
	}
	metrics := map[string][]metric{
		"air/many-agents": {
			{"manyAgents.heapLive.MiB", airStable},
			{"manyAgents.focusedContentLagP95Ms", airNoisy},
			{"manyAgents.edtMaxMs", airOver(50)},
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
