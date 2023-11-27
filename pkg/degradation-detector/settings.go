package degradation_detector

import (
  "fmt"
  dataQuery "github.com/JetBrains/ij-perf-report-aggregator/pkg/data-query"
  "net/url"
  "slices"
  "strings"
  "time"
)

type AnalysisSettings interface {
  GetDoNotReportImprovement() bool
  GetMinimumSegmentLength() int
  GetMedianDifferenceThreshold() float64
}

type CommonAnalysisSettings struct {
  DoNotReportImprovement    bool
  MinimumSegmentLength      int
  MedianDifferenceThreshold float64
}

func (s CommonAnalysisSettings) GetDoNotReportImprovement() bool {
  return s.DoNotReportImprovement
}

func (s CommonAnalysisSettings) GetMinimumSegmentLength() int {
  return s.MinimumSegmentLength
}

func (s CommonAnalysisSettings) GetMedianDifferenceThreshold() float64 {
  return s.MedianDifferenceThreshold
}

type Settings interface {
  DataQuery() []dataQuery.DataQuery
  DBTestName() string
  CreateSlackMessage(degradations []Degradation) SlackMessage
  SlackChannel() string
  GetMetric() string
  GetProject() string
  MergeAnother(settings Settings) Settings
  AnalysisSettings
}

type PerformanceSettings struct {
  Db          string
  Table       string
  Project     string
  Metric      string
  Branch      string
  Machine     string
  Channel     string
  ProductLink string
  CommonAnalysisSettings
}

func (s PerformanceSettings) MergeAnother(settings Settings) Settings {
  c := s
  c.Project = s.Project + "," + settings.GetProject()
  return c
}

func (s PerformanceSettings) DBTestName() string {
  return s.Project + "/" + s.Metric
}

func (s PerformanceSettings) slackLink() string {
  testPage := "tests"
  if strings.HasSuffix(s.Db, "Dev") {
    testPage = "testsDev"
  }
  machineGroup := getMachineGroup(s.Machine)
  return fmt.Sprintf("https://ij-perf.labs.jb.gg/%s/%s?machine=%s&branch=%s&project=%s&measure=%s&timeRange=1M",
    s.ProductLink, testPage, url.QueryEscape(machineGroup), url.QueryEscape(s.Branch), url.QueryEscape(s.Project), url.QueryEscape(s.Metric))
}

func (s PerformanceSettings) CreateSlackMessage(degradations []Degradation) SlackMessage {
  slices.SortFunc(degradations, func(a, b Degradation) int {
    return int(b.medianValues.PercentageChange() - a.medianValues.PercentageChange())
  })
  d := degradations[0]
  reason := getMessageBasedOnMedianChange(d.medianValues)
  date := time.UnixMilli(d.timestamp).UTC().Format("02-01-2006 15:04:05")
  link := s.slackLink()
  text := fmt.Sprintf(
    "%sTest: %s\n"+
      "Metric: %s\n"+
      "Build: %s\n"+
      "Branch: %s\n"+
      "Date: %s\n"+
      "Reason: %s\n"+
      "Link: %s", icon(d.medianValues), s.Project, s.Metric, d.Build, s.Branch, date, reason, link)
  return SlackMessage{
    Text:    text,
    Channel: s.Channel,
  }
}

func (s PerformanceSettings) SlackChannel() string {
  return s.Channel
}

func (s PerformanceSettings) GetMetric() string {
  return s.Metric
}

func (s PerformanceSettings) GetProject() string {
  return s.Project
}

type StartupSettings struct {
  Product     string
  Project     string
  Metric      string
  Branch      string
  Machine     string
  Channel     string
  ProductLink string
  CommonAnalysisSettings
}

func (s StartupSettings) DBTestName() string {
  return s.Product + "/" + s.Project + "/" + s.Metric
}

func (s StartupSettings) slackLink() string {
  machineGroup := getMachineGroup(s.Machine)
  return fmt.Sprintf("https://ij-perf.labs.jb.gg/%s/startup?machine=%s&branch=%s&product=%s&timeRange=1M",
    s.ProductLink, url.QueryEscape(machineGroup), url.QueryEscape(s.Branch), url.QueryEscape(s.Product))
}

func (s StartupSettings) CreateSlackMessage(degradations []Degradation) SlackMessage {
  slices.SortFunc(degradations, func(a, b Degradation) int {
    return int(b.medianValues.PercentageChange() - a.medianValues.PercentageChange())
  })
  d := degradations[0]
  reason := getMessageBasedOnMedianChange(d.medianValues)
  date := time.UnixMilli(d.timestamp).UTC().Format("02-01-2006 15:04:05")
  link := s.slackLink()

  text := fmt.Sprintf(
    "%sProject: %s\n"+
      "Metric: %s\n"+
      "Build: %s\n"+
      "Branch: %s\n"+
      "Date: %s\n"+
      "Reason: %s\n"+
      "Link: %s", icon(d.medianValues), s.Project, s.Metric, d.Build, s.Branch, date, reason, link)
  return SlackMessage{
    Text:    text,
    Channel: s.Channel,
  }
}

func (s StartupSettings) SlackChannel() string {
  return s.Channel
}

func (s StartupSettings) GetMetric() string {
  return s.Metric
}

func (s StartupSettings) GetProject() string {
  return s.Project
}

func (s StartupSettings) MergeAnother(settings Settings) Settings {
  c := s
  c.Project = s.Project + "," + settings.GetProject()
  return c
}

func icon(v MedianValues) string {
  icon := ""
  if v.newValue > v.previousValue {
    icon = ":chart_with_upwards_trend:"
  } else {
    icon = ":chart_with_downwards_trend:"
  }
  return icon
}
