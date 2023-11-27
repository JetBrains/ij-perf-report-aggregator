package degradation_detector

import (
  "testing"
)

func TestMergeDegradations(t *testing.T) {
  inputChan := make(chan DegradationWithSettings)
  go func() {
    inputChan <- DegradationWithSettings{
      Details: Degradation{Build: "123", medianValues: MedianValues{
        previousValue: 10,
        newValue:      20,
      }},
      Settings: PerformanceSettings{Project: "a", Metric: "metric", SlackSettings: SlackSettings{Channel: "slack"}},
    }
    inputChan <- DegradationWithSettings{
      Details: Degradation{Build: "123", medianValues: MedianValues{
        previousValue: 15,
        newValue:      20,
      }},
      Settings: PerformanceSettings{Project: "b", Metric: "metric", SlackSettings: SlackSettings{Channel: "slack"}},
    }
    close(inputChan)
  }()

  outputChan := MergeDegradations(inputChan)
  total := 0
  for r := range outputChan {
    sM := r.Settings.CreateSlackMessage(r.Details)
    eM := SlackMessage{
      Text: ":chart_with_upwards_trend:Test(s): a\nb\n" +
        "Metric: metric\n" +
        "Build: 123\n" +
        "Branch: \n" +
        "Date: 01-01-1970 00:00:00\n" +
        "Reason: Degradation detected. Median changed by: 100.00%. Median was 10.00 and now it is 20.00.\n" +
        "Link: https://ij-perf.labs.jb.gg//tests?machine=&branch=&project=a&project=b&measure=metric&timeRange=1M",
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

func TestSomeDegradationsNotMerged(t *testing.T) {
  inputChan := make(chan DegradationWithSettings)
  go func() {
    inputChan <- DegradationWithSettings{
      Details: Degradation{Build: "123", medianValues: MedianValues{
        previousValue: 10,
        newValue:      20,
      }},
      Settings: PerformanceSettings{Project: "a", Metric: "metric", SlackSettings: SlackSettings{Channel: "slack"}},
    }
    inputChan <- DegradationWithSettings{
      Details: Degradation{Build: "1234", medianValues: MedianValues{
        previousValue: 15,
        newValue:      20,
      }},
      Settings: PerformanceSettings{Project: "b", Metric: "metric", SlackSettings: SlackSettings{Channel: "slack"}},
    }
    inputChan <- DegradationWithSettings{
      Details: Degradation{Build: "123", medianValues: MedianValues{
        previousValue: 15,
        newValue:      20,
      }},
      Settings: PerformanceSettings{Project: "b", Metric: "metric", SlackSettings: SlackSettings{Channel: "slack"}},
    }
    close(inputChan)
  }()

  outputChan := MergeDegradations(inputChan)
  total := 0
  for range outputChan {
    total++
  }
  if total != 2 {
    t.Errorf("Incorrect merge")
  }
}
