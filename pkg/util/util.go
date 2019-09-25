package util

import (
	"go.uber.org/zap"
	"io"
	"os"
)

func Close(c io.Closer, log *zap.Logger) {
	err := c.Close()
	if err != nil && err != os.ErrClosed && err != io.ErrClosedPipe {
		if e, ok := err.(*os.PathError); ok && e.Err == os.ErrClosed {
			return
		}
		log.Error("cannot close", zap.Error(err))
	}
}
