package main

import (
  "fmt"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/server"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/util"
)

func main() {
  logger := util.CreateLogger()
  defer func() {
    _ = logger.Sync()
  }()

  err := server.Serve(util.GetEnv("CLICKHOUSE", server.DefaultDbUrl), util.GetEnv("NATS", ""), logger)
  if err != nil {
    logger.Fatal(fmt.Sprintf("%+v", err))
  }
}
