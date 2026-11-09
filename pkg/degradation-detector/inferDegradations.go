package degradation_detector

import (
	"runtime"

	"golang.org/x/sync/errgroup"
)

type DegradationWithSettings struct {
	Details  Degradation
	Settings Settings
}

func InferDegradations(data <-chan QueryResultWithSettings) <-chan DegradationWithSettings {
	degradationChan := make(chan DegradationWithSettings, 100)
	go func() {
		defer close(degradationChan)
		var group errgroup.Group
		group.SetLimit(runtime.GOMAXPROCS(0))
		for datum := range data {
			group.Go(func() error {
				ignoreBefore := datum.Settings.GetIgnoreChangesBefore()
				for _, degradation := range detectDegradations(datum.values, datum.builds, datum.timestamps, datum.Settings) {
					if !ignoreBefore.IsZero() && degradation.timestamp < ignoreBefore.UnixMilli() {
						continue
					}
					degradationChan <- DegradationWithSettings{Details: degradation, Settings: datum.Settings}
				}
				return nil
			})
		}
		_ = group.Wait()
	}()
	return degradationChan
}
