package setting

import (
	"encoding/json"
	"log/slog"
	"strings"

	detector "github.com/JetBrains/ij-perf-report-aggregator/pkg/degradation-detector"
	"github.com/JetBrains/ij-perf-report-aggregator/pkg/degradation-detector/alerts"
)

// The notifications of IJ Perf in the alerts repository.
const (
	metricDegradation   = "metric-degradation"
	metricThreshold     = "metric-threshold"
	unitTestDegradation = "unit-test-degradation"
)

// spec is the settings of a notification. The three notifications share most of the fields, and a field that a
// notification does not have keeps its zero value.
type spec struct {
	ReportType                        string              `json:"reportType"`
	MinimumSegmentLength              int                 `json:"minimumSegmentLength"`
	MedianDifferenceThreshold         float64             `json:"medianDifferenceThreshold"`
	EffectSizeThreshold               float64             `json:"effectSizeThreshold"`
	MissingDataDays                   int                 `json:"missingDataDays"`
	MetricRules                       string              `json:"metricRules"`
	OnlyAnalyze                       bool                `json:"onlyAnalyze"`
	MedianDifferenceThresholdByMetric map[string]float64  `json:"medianDifferenceThresholdByMetric"`
	MetricAliases                     map[string]string   `json:"metricAliases"`
	AdditionalMetrics                 map[string][]string `json:"additionalMetrics"`
	Comparison                        string              `json:"comparison"`
	Value                             float64             `json:"value"`
}

func parseSpec(raw json.RawMessage) (spec, error) {
	var s spec
	err := json.Unmarshal(raw, &s)
	return s, err
}

// analysis returns the settings of the analysis of one metric.
func (s spec) analysis(notification string, metric string) detector.AnalysisSettings {
	result := detector.AnalysisSettings{
		MinimumSegmentLength:      s.MinimumSegmentLength,
		MedianDifferenceThreshold: s.MedianDifferenceThreshold,
		EffectSizeThreshold:       s.EffectSizeThreshold,
		DaysToCheckMissing:        -s.MissingDataDays,
	}
	if threshold, ok := s.MedianDifferenceThresholdByMetric[metric]; ok {
		result.MedianDifferenceThreshold = threshold
	}
	switch s.ReportType {
	case "degradation":
		result.ReportType = detector.DegradationEvent
	case "improvement":
		result.ReportType = detector.ImprovementEvent
	}
	if notification == metricThreshold {
		result.AnalysisKind = detector.ThresholdAnalysis
		result.ThresholdValue = s.Value
		if s.Comparison == "lessThan" {
			result.ThresholdMode = detector.ThresholdLessThan
		}
	}
	return result
}

// slack returns where the messages go. IJ Perf sends only to Slack, so a target without a channel gets no message:
// its tests are analyzed, and the dashboards show their degradations.
func slack(target alerts.Target, onlyAnalyze bool, productLink string) detector.SlackSettings {
	if onlyAnalyze {
		return detector.SlackSettings{ProductLink: productLink}
	}
	return detector.SlackSettings{
		Channel:     target.SlackChannel(),
		ProductLink: productLink,
		Mention:     strings.Join(target.Mention, " "),
	}
}

// productLinks is the dashboard of a table: the part of a chart link after https://ij-perf.labs.jb.gg/.
// A table that is not here has a dashboard with the name of the table.
var productLinks = map[string]string{
	"idea":   "intellij",
	"ml":     "ml/dev",
	"ruby":   "rubymine",
	"report": "perfUnit",
}

func productLink(table string) string {
	if link, ok := productLinks[table]; ok {
		return link
	}
	return table
}

// like turns a pattern of the alerts config into a pattern of SQL LIKE.
func like(pattern string) string {
	return strings.ReplaceAll(pattern, "*", "%")
}

// generateMetricSettings returns the settings of one subscriber of metric-degradation or metric-threshold.
//
// A block of the subscriber lists tests, metrics, branches, machines and modes, and each combination is one series
// to analyze. The first block that selects a series wins, as for any notification of the alerts repository.
func generateMetricSettings(notification string, name string, subscriber alerts.Subscriber, backend Backend) []detector.Settings {
	var result []detector.Settings
	seen := map[seriesKey]bool{}
	for index, block := range subscriber.Match {
		log := slog.With("notification", notification, "subscriber", name, "block", index)
		s, err := parseSpec(block.Spec)
		if err != nil {
			log.Error("invalid spec in the alerts config", "error", err)
			continue
		}
		fields := block.Fields
		modes := fields["mode"]
		if len(modes) == 0 {
			modes = []string{"default"}
		}
		for _, db := range fields["db"] {
			for _, table := range fields["table"] {
				for _, branch := range fields["branch"] {
					for _, machine := range fields["machine"] {
						base := detector.PerformanceSettings{
							Db: db, Table: table, Branch: branch, Machine: like(machine),
						}
						for _, test := range selectTests(fields["project"], base, backend, log) {
							base.Project = test
							for _, metric := range selectMetrics(fields["metric"], s.MetricRules, base, backend, log) {
								for _, mode := range modes {
									settings := base
									settings.Metric = metric
									settings.MetricAlias = s.MetricAliases[metric]
									if mode != "default" {
										settings.Mode = mode
									}
									key := seriesKey{db, table, branch, machine, settings.Mode, test, metric}
									if seen[key] {
										continue
									}
									seen[key] = true
									settings.AnalysisSettings = s.analysis(notification, metric)
									settings.SlackSettings = slack(block.Target, s.OnlyAnalyze, productLink(table))
									result = append(result, settings)
								}
							}
						}
					}
				}
			}
		}
	}
	return result
}

type seriesKey struct {
	db, table, branch, machine, mode, test, metric string
}

// selectTests returns the tests of a block: a name as it is, and for a pattern the tests that ran on the branch and
// the machines of the settings.
func selectTests(projects []string, base detector.PerformanceSettings, backend Backend, log *slog.Logger) []string {
	var result []string
	var patterns []string
	for _, project := range projects {
		if alerts.IsPattern(project) {
			patterns = append(patterns, project)
		} else {
			result = append(result, project)
		}
	}
	if len(patterns) == 0 {
		return result
	}
	tests, err := backend.Tests(base.Db, base.Table, base.Branch, base.Machine)
	if err != nil {
		log.Error("error while getting tests", "error", err, "table", base.Table, "branch", base.Branch, "machine", base.Machine)
		return result
	}
	for _, test := range tests {
		if alerts.Matches(patterns, test) {
			result = append(result, test)
		}
	}
	return result
}

// selectMetrics returns the metrics of a test: the ones that the block lists, or the ones that the rules pick by the
// name of the test when the block lists none.
func selectMetrics(metrics []string, rules string, test detector.PerformanceSettings, backend Backend, log *slog.Logger) []string {
	if len(metrics) == 0 {
		rule, ok := metricRules[rules]
		if !ok {
			log.Error("unknown metric rules in the alerts config", "metricRules", rules)
			return nil
		}
		return rule(test.Project)
	}
	var result []string
	for _, metric := range metrics {
		if !alerts.IsPattern(metric) {
			result = append(result, metric)
			continue
		}
		names, err := backend.MetricNames(test, like(metric))
		if err != nil {
			log.Error("error while fetching metric names by pattern", "error", err, "test", test.Project, "pattern", metric)
			continue
		}
		result = append(result, names...)
	}
	return result
}
