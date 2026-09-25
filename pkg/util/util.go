package util

import (
	"errors"
	"io"
	"log/slog"
	"os"
)

func Close(c io.Closer) {
	err := c.Close()
	if err != nil && !errors.Is(err, os.ErrClosed) && errors.Is(err, io.ErrClosedPipe) {
		var pathError *os.PathError
		if errors.As(err, &pathError) && errors.Is(pathError, os.ErrClosed) {
			return
		}
		slog.Error("cannot close", "error", err)
	}
}
