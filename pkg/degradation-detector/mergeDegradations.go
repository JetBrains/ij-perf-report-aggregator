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
	if other, ok := settings.(PerformanceSettings); ok {
		projects := s.Projects()
		for _, project := range other.Projects() {
			if !slices.Contains(projects, project) {
				projects = append(projects, project)
			}
		}
		s.mergedProjects = projects
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

// Projects returns the tests of the settings: Project, and the tests that MergeDegradations merged into it.
// A test name can contain a comma, so the merged tests are kept as a list, not joined into Project.
func (s PerformanceSettings) Projects() []string {
	if len(s.mergedProjects) == 0 {
		return []string{s.Project}
	}
	return slices.Clone(s.mergedProjects)
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
			// The whole degradation goes on, so its previous build and the analysis that found it reach the message.
			c <- DegradationWithSettings{Details: v.Details[0], Settings: v.Settings}
		}
	}()

	return c
}
