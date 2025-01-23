package degradation_detector

import "log/slog"

type DegradationWithSettings struct {
	Details  Degradation
	Settings Settings
}

func InferDegradations(data <-chan QueryResultWithSettings) <-chan DegradationWithSettings {
	degradationChan := make(chan DegradationWithSettings, 100)
	go func() {
		for datum := range data {
			slog.Info("infer degradations", "settings", datum.Settings)
			for _, degradation := range detectDegradations(datum.values, datum.builds, datum.timestamps, datum.Settings) {
				degradationChan <- DegradationWithSettings{Details: degradation, Settings: datum.Settings}
			}
		}
		close(degradationChan)
	}()
	return degradationChan
}
