package degradation_detector

import (
  "testing"
)

// TestMergeDegradations tests the mergeDegradations function
func TestMergeDegradations(t *testing.T) {
  inputChan := make(chan DegradationWithContext)
  go func() {
    inputChan <- DegradationWithContext{
      Details: Degradation{Build: "123", medianValues: MedianValues{
        previousValue: 10,
        newValue:      20,
      }},
      Settings: PerformanceSettings{Project: "a", Channel: "slack", Metric: "metric"},
    }
    inputChan <- DegradationWithContext{
      Details: Degradation{Build: "123", medianValues: MedianValues{
        previousValue: 15,
        newValue:      20,
      }},
      Settings: PerformanceSettings{Project: "b", Channel: "slack", Metric: "metric"},
    }
    close(inputChan)
  }()

  outputChan := MergeDegradations(inputChan)
  total := 0
  for r := range outputChan {
    sM := r.Settings.CreateSlackMessage(r.Details)
    eM := SlackMessage{
      Text: ":chart_with_upwards_trend:Test: a,b\n" +
        "Metric: metric\n" +
        "Build: 123\n" +
        "Branch: \n" +
        "Date: 01-01-1970 00:00:00\n" +
        "Reason: Degradation detected. Median changed by: 100.00%. Median was 10.00 and now it is 20.00.\n" +
        "Link: https://ij-perf.labs.jb.gg//tests?machine=&branch=&project=a%2Cb&measure=metric&timeRange=1M",
      Channel: r.Settings.SlackChannel(),
    }
    if sM != eM {
      t.Errorf("Incorrect slack message: %v", sM)
    }
    total++
  }
  if total != 1 {
    t.Errorf("Too many degradations, they were not merged")
  }
}
