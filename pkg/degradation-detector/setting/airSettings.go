package setting

import (
	"fmt"

	detector "github.com/JetBrains/ij-perf-report-aggregator/pkg/degradation-detector"
)

func GenerateAirSettings() []detector.PerformanceSettings {
	longSession := []string{
		"open.%s.coverRemovedMs",
		"open.%s.gestureToSettledMs",
		"continue.%s.sendToPromptMs",
		"answer.%s.streamStartToCompletedMs",
		"prompt.%s.codeP95Ms",
		"scroll.%s.wheelP95Ms",
		"heap.%s.heapLive.MiB",
	}
	metrics := map[string][]string{
		"air/many-agents":  {"manyAgents.focusedTabLagMs", "manyAgents.openTabLagMaxMs", "manyAgents.edtP95Ms", "manyAgents.heapLive.MiB"},
		"air/session-list": {"sessionList.firstShowMs", "sessionList.refreshMedianMs", "sessionList.wheelP95Ms", "sessionList.heapLive.MiB"},
	}
	for _, size := range []string{"M", "XL", "XXL"} {
		project := "air/long-session/" + size
		for _, metric := range longSession {
			metrics[project] = append(metrics[project], fmt.Sprintf(metric, size))
		}
	}

	settings := make([]detector.PerformanceSettings, 0, 30)
	for test, metrics := range metrics {
		for _, metric := range metrics {
			settings = append(settings, detector.PerformanceSettings{
				Db:          "perfintDev",
				Table:       "air",
				Project:     test,
				Machine:     "intellij-linux-performance-aws-%",
				Metric:      metric,
				Branch:      "master",
				Channel:     "ij-air-perf",
				ProductLink: "air",
				ReportType:  detector.DegradationEvent,
			})
		}
	}
	return settings
}
