package setting

import (
	"fmt"

	detector "github.com/JetBrains/ij-perf-report-aggregator/pkg/degradation-detector"
)

// The Air performance tests run each scenario 5 times a night and report the median of the passes, so one value is one
// stable measurement. airStep judges each night against the 10 nights before it and reports a step on the first slow
// night, with the last fast build beside it. Its band is the largest of a relative change, an absolute change and 5
// robust standard deviations of the baseline, so a noisy metric needs a larger change. Over the nights from 2026-10-12
// to 2026-11-21 it reported each step that a person saw in the charts, the step of 2026-11-19 on its first night, and
// no step that the charts do not show.
func airStep(relative float64, absolute float64) detector.AnalysisSettings {
	return detector.AnalysisSettings{
		ReportType:   detector.DegradationEvent,
		AnalysisKind: detector.StepAnalysis,
		Step:         detector.StepSettings{Relative: relative, Absolute: absolute, Noise: 5},
	}
}

var (
	airTime = airStep(25, 30)
	airHeap = airStep(20, 50)
)

// airOverWindow is the number of nights whose median airOver compares to the limit. One slow night of a run that
// measures once a night is noise, so it alone does not cross the limit.
const airOverWindow = 3

// airOver alerts when the median of the latest airOverWindow nights crosses limitMs. A metric of a few frames moves by
// one frame from night to night, which is 30-80 % of its value, so a step of its median says nothing. A value over the
// limit is what a user notices. The alert comes once, on the night of the crossing, not on each night over the limit.
func airOver(limitMs float64) detector.AnalysisSettings {
	return detector.AnalysisSettings{
		ReportType:          detector.DegradationEvent,
		AnalysisKind:        detector.ThresholdAnalysis,
		ThresholdMode:       detector.ThresholdGreaterThan,
		ThresholdValue:      limitMs,
		ThresholdWindow:     airOverWindow,
		ThresholdOnCrossing: true,
	}
}

func GenerateAirSettings() []detector.PerformanceSettings {
	type metric struct {
		name     string
		analysis detector.AnalysisSettings
	}
	longSession := []metric{
		{"open.%s.gestureToTabMs", airTime},
		{"open.%s.coverRemovedMs", airTime},
		{"checkpoint.%s.frontend_transcript_rows_applied.elapsedMs", airTime},
		{"continue.%s.loadAnsweredToPromptMs", airOver(300)},
		{"answer.%s.streamLagP95Ms", airTime},
		{"answer.%s.streamEndToSettledMs", airTime},
		{"heap.%s.heapLive.MiB", airHeap},
		{"prompt.%s.codeP95Ms", airOver(50)},
		{"scroll.%s.wheelP95Ms", airOver(50)},
		// A jump to the start and to 0.4 of the history lands on rows that are not built yet. A blank screen of 200 ms
		// after a drag is what a user notices.
		{"scroll.%s.jump0.settledMs", airOver(200)},
		{"scroll.%s.jump40.settledMs", airOver(200)},
	}
	metrics := map[string][]metric{
		"air/many-agents": {
			{"manyAgents.heapLive.MiB", airHeap},
			{"manyAgents.focusedContentLagP95Ms", airTime},
			// The longest EDT wait of the run is about 50 ms on a usual night, so the limit is twice that: 6 frames.
			{"manyAgents.edtMaxMs", airOver(100)},
		},
		"air/session-list": {
			{"sessionList.heapLive.MiB", airHeap},
			{"sessionList.refreshMedianMs", airTime},
			{"sessionList.firstShowMs", airTime},
		},
	}
	for _, size := range []string{"M", "XL", "XXL"} {
		project := "air/long-session/" + size
		for _, m := range longSession {
			metrics[project] = append(metrics[project], metric{fmt.Sprintf(m.name, size), m.analysis})
		}
	}

	settings := make([]detector.PerformanceSettings, 0, 37)
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
