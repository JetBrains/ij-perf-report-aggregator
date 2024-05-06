package main

import (
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/server"
  "go.deanishe.net/env"
  "log/slog"
  "os"
)

func main() {
  err := server.Serve(env.Get("CLICKHOUSE", server.DefaultDbUrl), env.Get("NATS", ""))
  if err != nil {
    slog.Error("error on starting backend", "error", err)
    os.Exit(1)
  }
}
