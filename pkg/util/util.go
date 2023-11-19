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
  taskContext, cancel := context.WithCancel(context.Background())
  signals := make(chan os.Signal, 1)
  signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
  go func() {
    <-signals
    println("cancel on signal")
    cancel()
  }()
  return taskContext, cancel
}

func GetEnv(name string, defaultValue string) string {
  value := os.Getenv(name)
  if len(value) == 0 {
    return defaultValue
  }
  return value
}

func GetEnvOrFile(envName string, file string) (string, error) {
  v := os.Getenv(envName)
  if len(v) == 0 {
    b, err := os.ReadFile(file)
    if err != nil {
      return "", err
    }
    return string(b), err
  }
  return v, nil
}

func GetEnvOrFileOrPanic(envName string, file string) string {
  v := os.Getenv(envName)
  if len(v) == 0 {
    b, err := os.ReadFile(file)
    if err != nil {
      panic(err)
    }
    return string(b)
  }
  return v
}
