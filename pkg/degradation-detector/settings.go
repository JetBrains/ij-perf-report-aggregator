package degradation_detector

type Settings interface {
  analysisSettings
  dataQueryProducer
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
  Channel     string
  ProductLink string
  CommonAnalysisSettings
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
