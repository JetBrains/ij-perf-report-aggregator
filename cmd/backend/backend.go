package main

import (
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/server"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/util"
  "log/slog"
  "os"
)

func main() {
  err := server.Serve(util.GetEnv("CLICKHOUSE", server.DefaultDbUrl), util.GetEnv("NATS", ""))
  if err != nil {
    slog.Error("error", err)
    os.Exit(1)
  }
}
