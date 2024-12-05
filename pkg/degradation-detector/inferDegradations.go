package degradation_detector

import (
	"net/http"
)

type DegradationWithSettings struct {
	Details  Degradation
	Settings Settings
}

func InferDegradations(data <-chan queryResultWithSettings) <-chan DegradationWithSettings {
	degradationChan := make(chan DegradationWithSettings, 100)
	go func() {
		for datum := range data {
			for _, degradation := range detectDegradations(datum.values, datum.builds, datum.timestamps, datum.Settings) {
				degradationChan <- DegradationWithSettings{Details: degradation, Settings: datum.Settings}
			}
		}
		close(degradationChan)
	}()
	return degradationChan
}

func GetDegradations(settings []Settings, client *http.Client, backendUrl string) <-chan DegradationWithSettings {
	return InferDegradations(fetchMetricsFromClickhouse(settings, client, backendUrl))
}
