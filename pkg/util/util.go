package util

import (
  "context"
  "go.uber.org/zap"
	"io"
	"os"
  "os/signal"
  "syscall"
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

func CreateCommandContext() (context.Context, context.CancelFunc) {
  taskContext, cancel := context.WithCancel(context.Background())
  signals := make(chan os.Signal, 1)
  signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
  go func() {
    <-signals
    cancel()
  }()
  return taskContext, cancel
}
