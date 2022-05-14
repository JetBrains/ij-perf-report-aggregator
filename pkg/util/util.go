package util

import (
  "context"
  "go.uber.org/zap"
  "io"
  "log"
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
    println("cancel on signal")
    cancel()
  }()
  return taskContext, cancel
}

func CreateLogger() *zap.Logger {
  config := zap.NewProductionConfig()
  config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
  config.DisableCaller = true
  config.DisableStacktrace = true
  // https://www.outcoldsolutions.com/blog/2018-08-10-timestamps-in-container-logs/
  config.EncoderConfig.TimeKey = ""
  logger, err := config.Build()
  if err != nil {
    log.Fatal(err)
  }
  return logger
}

func GetEnvOrPanic(name string) string {
  value := os.Getenv(name)
  if len(value) == 0 {
    panic("env " + name + " is not set")
  }
  return value
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
