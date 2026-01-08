package degradation_detector

import (
	"runtime"
	"sync"

	"github.com/alitto/pond"
)

type DegradationWithSettings struct {
	Details  Degradation
	Settings Settings
}

func InferDegradations(data <-chan QueryResultWithSettings) <-chan DegradationWithSettings {
	degradationChan := make(chan DegradationWithSettings, 100)
	go func() {
		var wg sync.WaitGroup
		pool := pond.New(runtime.NumCPU(), 1000)
		defer pool.StopAndWait()
		for datum := range data {
			wg.Add(1)
			pool.Submit(func() {
				defer wg.Done()
				for _, degradation := range detectDegradations(datum.values, datum.builds, datum.timestamps, datum.Settings) {
					degradationChan <- DegradationWithSettings{Details: degradation, Settings: datum.Settings}
				}
			})
		}
		wg.Wait()
		close(degradationChan)
	}()
	return degradationChan
}
