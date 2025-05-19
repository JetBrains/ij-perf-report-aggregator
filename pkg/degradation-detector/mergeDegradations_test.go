package degradation_detector

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMergeDegradations(t *testing.T) {
	t.Parallel()
	inputChan := make(chan DegradationWithSettings)
	go func() {
		inputChan <- DegradationWithSettings{
			Details: Degradation{Build: "123", medianValues: MedianValues{
				previousValue: 10,
				newValue:      20,
			}},
			Settings: PerformanceSettings{Project: "a", BaseSettings: BaseSettings{Metric: "metric", SlackSettings: SlackSettings{Channel: "slack"}}},
		}
		inputChan <- DegradationWithSettings{
			Details: Degradation{Build: "123", medianValues: MedianValues{
				previousValue: 15,
				newValue:      20,
			}},
			Settings: PerformanceSettings{Project: "b", BaseSettings: BaseSettings{Metric: "metric", SlackSettings: SlackSettings{Channel: "slack"}}},
		}
		close(inputChan)
	}()

	outputChan := MergeDegradations(inputChan)
	total := 0
	now := time.Now()
	currentDate := fmt.Sprintf("%d-%02d-%d", now.Year(), now.Month(), now.Day())
	for r := range outputChan {
		sM := r.Settings.CreateSlackMessage(r.Details)
		eM := SlackMessage{
			Text: ":chart_with_upwards_trend:Test(s): a\nb\n" +
				"Metric: metric\n" +
				"Mode: default\n" +
				"Build: 123\n" +
				"Branch: \n" +
				"Date: 01-01-1970 00:00:00\n" +
				"Reason: Degradation detected. Median changed by: 100.00%. Median was 10.00 and now it is 20.00.\n" +
				"<https://ij-perf.labs.jb.gg//tests?mode=default&machine=&branch=&project=a&project=b&measure=metric&timeRange=custom&customRange=1969-12-25:" + currentDate + "&point=|See charts>\n" +
				"<https://ij-perf.labs.jb.gg/degradations/report?tests=a,b&build=123&date=01-01-1970|Report event>",
			Channel: r.Settings.SlackChannel(),
		}
		assert.Equal(t, eM, sM, "Incorrect slack message")
		total++
	}
	assert.Equal(t, 1, total, "Incorrect merge")
}

func TestSomeDegradationsNotMerged(t *testing.T) {
	t.Parallel()
	inputChan := make(chan DegradationWithSettings)
	go func() {
		inputChan <- DegradationWithSettings{
			Details: Degradation{Build: "123", medianValues: MedianValues{
				previousValue: 10,
				newValue:      20,
			}},
			Settings: PerformanceSettings{Project: "a", BaseSettings: BaseSettings{Metric: "metric", SlackSettings: SlackSettings{Channel: "slack"}}},
		}
		inputChan <- DegradationWithSettings{
			Details: Degradation{Build: "1234", medianValues: MedianValues{
				previousValue: 15,
				newValue:      20,
			}},
			Settings: PerformanceSettings{Project: "b", BaseSettings: BaseSettings{Metric: "metric", SlackSettings: SlackSettings{Channel: "slack"}}},
		}
		inputChan <- DegradationWithSettings{
			Details: Degradation{Build: "123", medianValues: MedianValues{
				previousValue: 15,
				newValue:      20,
			}},
			Settings: PerformanceSettings{Project: "b", BaseSettings: BaseSettings{Metric: "metric", SlackSettings: SlackSettings{Channel: "slack"}}},
		}
		close(inputChan)
	}()

	outputChan := MergeDegradations(inputChan)
	total := 0
	for range outputChan {
		total++
	}
	assert.Equal(t, 2, total, "Incorrect merge")
}

func TestMetricAlias(t *testing.T) {
	t.Parallel()
	inputChan := make(chan DegradationWithSettings)
	go func() {
		inputChan <- DegradationWithSettings{
			Details: Degradation{Build: "123", medianValues: MedianValues{
				previousValue: 10,
				newValue:      20,
			}},
			Settings: PerformanceSettings{Project: "a", MetricAlias: "metric", BaseSettings: BaseSettings{Metric: "metricBetta", SlackSettings: SlackSettings{Channel: "slack"}}},
		}
		inputChan <- DegradationWithSettings{
			Details: Degradation{Build: "123", medianValues: MedianValues{
				previousValue: 10,
				newValue:      20,
			}},
			Settings: PerformanceSettings{Project: "a", MetricAlias: "metric", BaseSettings: BaseSettings{Metric: "metricBetta", SlackSettings: SlackSettings{Channel: "slack"}}},
		}
		inputChan <- DegradationWithSettings{
			Details: Degradation{Build: "123", medianValues: MedianValues{
				previousValue: 15,
				newValue:      20,
			}},
			Settings: PerformanceSettings{Project: "b", MetricAlias: "metric", BaseSettings: BaseSettings{Metric: "metricAlpha", SlackSettings: SlackSettings{Channel: "slack"}}},
		}
		close(inputChan)
	}()

	outputChan := MergeDegradations(inputChan)
	total := 0
	now := time.Now()
	currentDate := fmt.Sprintf("%d-%02d-%d", now.Year(), now.Month(), now.Day())
	for r := range outputChan {
		sM := r.Settings.CreateSlackMessage(r.Details)
		eM := SlackMessage{
			Text: ":chart_with_upwards_trend:Test(s): a\nb\n" +
				"Metric: metricBetta,metricAlpha\n" +
				"Mode: default\n" +
				"Build: 123\n" +
				"Branch: \n" +
				"Date: 01-01-1970 00:00:00\n" +
				"Reason: Degradation detected. Median changed by: 100.00%. Median was 10.00 and now it is 20.00.\n" +
				"<https://ij-perf.labs.jb.gg//tests?mode=default&machine=&branch=&project=a&project=b&measure=metricBetta&measure=metricAlpha&timeRange=custom&customRange=1969-12-25:" + currentDate + "&point=|See charts>\n" +
				"<https://ij-perf.labs.jb.gg/degradations/report?tests=a,b&build=123&date=01-01-1970|Report event>",
			Channel: r.Settings.SlackChannel(),
		}
		assert.Equal(t, eM, sM, "Incorrect slack message")
		total++
	}
	assert.Equal(t, 1, total, "Incorrect merge")
}
