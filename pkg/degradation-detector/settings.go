package degradation_detector

type Settings interface {
	analysisSettings
	queryProducer
	accidentWriter
	mergeInfoProvider
	slackSettings
	metricsMerger
}

type BaseSettings struct {
	SlackSettings
	AnalysisSettings

	Metric  string
	Branch  string
	Machine string
}

type PerformanceSettings struct {
	BaseSettings

	Db          string
	Table       string
	Project     string
	MetricAlias string
	Mode        string
}

type StartupSettings struct {
	BaseSettings

	Db      string
	Table   string
	Product string
	Project string
}

type FleetStartupSettings struct {
	BaseSettings
}
