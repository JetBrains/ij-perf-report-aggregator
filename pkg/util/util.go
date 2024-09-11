package util

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
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

func CreateCommandContext() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signals
		println("cancel on signal")
		cancel()
	}()
	return ctx, cancel
}
