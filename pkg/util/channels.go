package util

func Broadcast[T any](input <-chan T, outputs ...chan T) {
	go func() {
		for value := range input {
			for _, out := range outputs {
				go func() {
					out <- value
				}()
			}
		}
		// Close all output channels when input is exhausted
		for _, out := range outputs {
			close(out)
		}
	}()
}
