package util

import "sync"

func Broadcast[T any](input <-chan T, outputs ...chan T) {
	var wg sync.WaitGroup
	go func() {
		for value := range input {
			for _, out := range outputs {
				wg.Add(1)
				go func() {
					defer wg.Done()
					out <- value
				}()
			}
		}
		wg.Wait()
		// Close all output channels when input is exhausted
		for _, out := range outputs {
			close(out)
		}
	}()
}
