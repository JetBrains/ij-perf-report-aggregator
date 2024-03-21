package degradation_detector

type Settings interface {
  analysisSettings
  queryProducer
  accidentWriter
  mergeInfoProvider
  slackSettings
}

type PerformanceSettings struct {
  Db          string
  Table       string
  Project     string
  Metric      string
  Branch      string
  Machine     string
  MetricAlias string
  SlackSettings
  AnalysisSettings
}

type StartupSettings struct {
  Product string
  Project string
  Metric  string
  Branch  string
  Machine string
  SlackSettings
  AnalysisSettings
}

type FleetStartupSettings struct {
  Metric  string
  Branch  string
  Machine string
  SlackSettings
  AnalysisSettings
}
