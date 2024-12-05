package degradation_detector

import (
	"fmt"
	"slices"
	"strings"
)

type multipleDegradationWithSettings struct {
	Details  []Degradation
	Settings Settings
}

type mergeInfoProvider interface {
	GetMetric() string
	GetMetricOrAlias() string
	GetProject() string
	MergeAnother(settings Settings) Settings
}

func (s PerformanceSettings) MergeAnother(settings Settings) Settings {
	if !strings.Contains(s.Project, settings.GetProject()) {
		s.Project = fmt.Sprintf("%s,%s", s.Project, settings.GetProject())
	}
	if s.MetricAlias != "" && s.Metric != settings.GetMetric() && !strings.Contains(s.Metric, settings.GetMetric()) {
		s.Metric = fmt.Sprintf("%s,%s", s.Metric, settings.GetMetric())
	}
	return s
}

func (s BaseSettings) GetMetric() string {
	return s.Metric
}

func (s BaseSettings) GetMetricOrAlias() string {
	return s.Metric
}

func (s PerformanceSettings) GetMetricOrAlias() string {
	if s.MetricAlias != "" {
		return s.MetricAlias
	}
	return s.Metric
}

func (s PerformanceSettings) GetProject() string {
	return s.Project
}

func (s StartupSettings) GetProject() string {
	return s.Project
}

func (s StartupSettings) MergeAnother(settings Settings) Settings {
	s.Project = fmt.Sprintf("%s,%s", s.Project, settings.GetProject())
	return s
}

func (s FleetStartupSettings) GetProject() string {
	return "fleet"
}

func (s FleetStartupSettings) MergeAnother(settings Settings) Settings {
	s.Metric = fmt.Sprintf("%s,%s", s.Metric, settings.GetMetric())
	return s
}

func MergeDegradations(degradations <-chan DegradationWithSettings) chan DegradationWithSettings {
	c := make(chan DegradationWithSettings, 100)
	go func() {
		defer close(c)

		type degradationKey struct {
			slackChannel string
			metric       string
			build        string
		}
		m := make(map[degradationKey]multipleDegradationWithSettings, 100)
		// Collect and merge degradations
		for d := range degradations {
			key := degradationKey{
				slackChannel: d.Settings.SlackChannel(),
				metric:       d.Settings.GetMetricOrAlias(),
				build:        d.Details.Build,
			}
			if existing, found := m[key]; found {
				m[key] = multipleDegradationWithSettings{
					Details:  append([]Degradation{d.Details}, existing.Details...),
					Settings: existing.Settings.MergeAnother(d.Settings),
				}
			} else {
				m[key] = multipleDegradationWithSettings{
					Details:  []Degradation{d.Details},
					Settings: d.Settings,
				}
			}
		}
		// Find the largest degradation
		for _, v := range m {
			slices.SortFunc(v.Details, func(a, b Degradation) int {
				return int(b.medianValues.PercentageChange() - a.medianValues.PercentageChange())
			})
			highest := v.Details[0]
			c <- DegradationWithSettings{
				Details: Degradation{
					Build:         highest.Build,
					timestamp:     highest.timestamp,
					medianValues:  highest.medianValues,
					IsDegradation: highest.IsDegradation,
				},
				Settings: v.Settings,
			}
		}
	}()

	return c
}
