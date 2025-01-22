package util

func Broadcast[T any](input <-chan T, outputs ...chan T) {
	go func() {
		defer func() {
			for _, out := range outputs {
				close(out)
			}
		}()

		for value := range input {
			for _, out := range outputs {
				out <- value
			}
		}
	}()
}
