package degradation_detector

import (
	"fmt"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMergeDegradations(t *testing.T) {
	t.Parallel()
	inputChan := make(chan DegradationWithSettings)
	go func() {
		inputChan <- DegradationWithSettings{
			Details: Degradation{Build: "123", medianValues: CenterValues{
				previousValue: 10,
				newValue:      20,
			}},
			Settings: PerformanceSettings{Project: "a", BaseSettings: BaseSettings{Metric: "metric", SlackSettings: SlackSettings{Channel: "slack"}}},
		}
		inputChan <- DegradationWithSettings{
			Details: Degradation{Build: "123", medianValues: CenterValues{
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
				"Metric: `metric`\n" +
				"Mode: default\n" +
				"Build: 123\n" +
				"Branch: \n" +
				"Date: 01-01-1970 00:00:00\n" +
				"Reason: Degradation detected. Median changed by: 100.00%. Median was 10.00 and now it is 20.00.\n" +
				"<https://ij-perf.labs.jb.gg//tests?mode=default&machine=&branch=&project=a&project=b&measure=metric&timeRange=custom&customRange=1969-12-1:" + currentDate + "&point=123|See charts>\n" +
				"<https://ij-perf.labs.jb.gg/degradations/report?test=a&test=b&build=123&date=01-01-1970|Report event>",
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
			Details: Degradation{Build: "123", medianValues: CenterValues{
				previousValue: 10,
				newValue:      20,
			}},
			Settings: PerformanceSettings{Project: "a", BaseSettings: BaseSettings{Metric: "metric", SlackSettings: SlackSettings{Channel: "slack"}}},
		}
		inputChan <- DegradationWithSettings{
			Details: Degradation{Build: "1234", medianValues: CenterValues{
				previousValue: 15,
				newValue:      20,
			}},
			Settings: PerformanceSettings{Project: "b", BaseSettings: BaseSettings{Metric: "metric", SlackSettings: SlackSettings{Channel: "slack"}}},
		}
		inputChan <- DegradationWithSettings{
			Details: Degradation{Build: "123", medianValues: CenterValues{
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
			Details: Degradation{Build: "123", medianValues: CenterValues{
				previousValue: 10,
				newValue:      20,
			}},
			Settings: PerformanceSettings{Project: "a", MetricAlias: "metric", BaseSettings: BaseSettings{Metric: "metricBetta", SlackSettings: SlackSettings{Channel: "slack"}}},
		}
		inputChan <- DegradationWithSettings{
			Details: Degradation{Build: "123", medianValues: CenterValues{
				previousValue: 10,
				newValue:      20,
			}},
			Settings: PerformanceSettings{Project: "a", MetricAlias: "metric", BaseSettings: BaseSettings{Metric: "metricBetta", SlackSettings: SlackSettings{Channel: "slack"}}},
		}
		inputChan <- DegradationWithSettings{
			Details: Degradation{Build: "123", medianValues: CenterValues{
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
				"Metric: `metricBetta,metricAlpha`\n" +
				"Mode: default\n" +
				"Build: 123\n" +
				"Branch: \n" +
				"Date: 01-01-1970 00:00:00\n" +
				"Reason: Degradation detected. Median changed by: 100.00%. Median was 10.00 and now it is 20.00.\n" +
				"<https://ij-perf.labs.jb.gg//tests?mode=default&machine=&branch=&project=a&project=b&measure=metricBetta&measure=metricAlpha&timeRange=custom&customRange=1969-12-1:" + currentDate + "&point=123|See charts>\n" +
				"<https://ij-perf.labs.jb.gg/degradations/report?test=a&test=b&build=123&date=01-01-1970|Report event>",
			Channel: r.Settings.SlackChannel(),
		}
		assert.Equal(t, eM, sM, "Incorrect slack message")
		total++
	}
	assert.Equal(t, 1, total, "Incorrect merge")
}

func TestMergeDegradationsKeepsCommaInTestName(t *testing.T) {
	t.Parallel()
	const withComma = "com.intellij.java.find.FindInCommentsAndLiteralsPerformanceTest.testLivePreviewUpdateWhenAllOccurrencesShareOneToken - live preview update except literals, all occurrences in one token"
	inputChan := make(chan DegradationWithSettings)
	go func() {
		for _, project := range []string{withComma, "b", withComma} {
			inputChan <- DegradationWithSettings{
				Details:  Degradation{Build: "123", medianValues: CenterValues{previousValue: 10, newValue: 20}},
				Settings: PerformanceSettings{Project: project, BaseSettings: BaseSettings{Metric: "metric", SlackSettings: SlackSettings{Channel: "slack"}}},
			}
		}
		close(inputChan)
	}()

	total := 0
	for r := range MergeDegradations(inputChan) {
		settings, ok := r.Settings.(PerformanceSettings)
		assert.True(t, ok)
		assert.Equal(t, []string{withComma, "b"}, settings.Projects())
		escaped := url.QueryEscape(withComma)
		link := settings.ChartLink(r.Details)
		assert.Contains(t, link, "&project="+escaped+"&project=b&measure=metric&")
		message := settings.CreateSlackMessage(r.Details).Text
		assert.Contains(t, message, "Test(s): "+withComma+"\nb\n")
		assert.Contains(t, message, "/degradations/report?test="+escaped+"&test=b&build=123&")
		total++
	}
	assert.Equal(t, 1, total)
}
