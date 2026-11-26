package setting

import (
	detector "github.com/JetBrains/ij-perf-report-aggregator/pkg/degradation-detector"
	"github.com/JetBrains/ij-perf-report-aggregator/pkg/degradation-detector/alerts"
)

type Generator struct {
	Name     string
	Generate func() []detector.Settings
}

// Generators returns one generator per subscriber of a metric notification and one for the unit tests, so
// that a subscriber whose settings select nothing is visible. Air is not in the alerts config: its settings are in Go.
func Generators(config *alerts.Config, backend Backend) []Generator {
	generators := []Generator{{Name: "air", Generate: func() []detector.Settings {
		air := GenerateAirSettings()
		settings := make([]detector.Settings, 0, len(air))
		for _, s := range air {
			settings = append(settings, s)
		}
		return settings
	}}}
	for _, id := range []string{metricDegradation, metricThreshold} {
		notification, ok := config.Notifications[id]
		if !ok {
			// A notification that the config lost or renamed generates nothing, so it is reported as empty.
			generators = append(generators, Generator{Name: id, Generate: func() []detector.Settings { return nil }})
			continue
		}
		for _, name := range notification.SubscriberNames() {
			subscriber := notification.Subscribers[name]
			if !subscriber.Enabled {
				continue
			}
			generators = append(generators, Generator{
				Name:     id + "/" + name,
				Generate: func() []detector.Settings { return generateMetricSettings(id, name, subscriber, backend) },
			})
		}
	}
	return append(generators, Generator{
		Name: unitTestDegradation,
		Generate: func() []detector.Settings {
			return generateUnitTestSettings(config.Notifications[unitTestDegradation], backend)
		},
	})
}
