package degradation_detector

import (
	"fmt"
	"slices"
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
	c := s
	c.Project = fmt.Sprintf("%s,%s", s.Project, settings.GetProject())
	if s.MetricAlias != "" && s.Metric != settings.GetMetric() {
		c.Metric = fmt.Sprintf("%s,%s", s.Metric, settings.GetMetric())
	}
	return c
}

func (s PerformanceSettings) GetMetricOrAlias() string {
	if s.MetricAlias != "" {
		return s.MetricAlias
	}
	return s.Metric
}

func (s PerformanceSettings) GetMetric() string {
	return s.Metric
}

func (s PerformanceSettings) GetProject() string {
	return s.Project
}

func (s StartupSettings) GetMetric() string {
	return s.Metric
}

func (s StartupSettings) GetMetricOrAlias() string {
	return s.Metric
}

func (s StartupSettings) GetProject() string {
	return s.Project
}

func (s StartupSettings) MergeAnother(settings Settings) Settings {
	c := s
	c.Project = fmt.Sprintf("%s,%s", s.Project, settings.GetProject())
	return c
}

func MergeDegradations(degradations <-chan DegradationWithSettings) chan DegradationWithSettings {
	c := make(chan DegradationWithSettings)
	go func() {
		type degradationKey struct {
			slackChannel string
			metric       string
			build        string
		}
		m := make(map[degradationKey]multipleDegradationWithSettings, 100)
		for degradation := range degradations {
			key := degradationKey{
				slackChannel: degradation.Settings.SlackChannel(),
				metric:       degradation.Settings.GetMetricOrAlias(),
				build:        degradation.Details.Build,
			}
			if existing, found := m[key]; found {
				d := []Degradation{degradation.Details}
				d = append(d, existing.Details...)
				m[key] = multipleDegradationWithSettings{
					Details:  d,
					Settings: existing.Settings.MergeAnother(degradation.Settings),
				}
			} else {
				d := make([]Degradation, 0)
				d = append(d, degradation.Details)
				m[key] = multipleDegradationWithSettings{
					Details:  d,
					Settings: degradation.Settings,
				}
			}
		}
		for _, v := range m {
			merged := v
			slices.SortFunc(merged.Details, func(a, b Degradation) int {
				return int(b.medianValues.PercentageChange() - a.medianValues.PercentageChange())
			})
			d := merged.Details[0]
			c <- DegradationWithSettings{
				Details: Degradation{
					Build:         d.Build,
					timestamp:     d.timestamp,
					medianValues:  d.medianValues,
					IsDegradation: d.IsDegradation,
				},
				Settings: v.Settings,
			}
		}
		close(c)
	}()
	return c
}
