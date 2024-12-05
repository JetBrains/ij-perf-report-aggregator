package degradation_detector

type Settings interface {
	analysisSettings
	queryProducer
	accidentWriter
	mergeInfoProvider
	slackSettings
}

type BaseSettings struct {
	Metric  string
	Branch  string
	Machine string
	SlackSettings
	AnalysisSettings
}

type PerformanceSettings struct {
	BaseSettings
	Db          string
	Table       string
	Project     string
	MetricAlias string
}

type StartupSettings struct {
	BaseSettings
	Product string
	Project string
}

type FleetStartupSettings struct {
	BaseSettings
}
