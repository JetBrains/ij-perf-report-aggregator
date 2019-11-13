package main

import (
  "flag"
  "fmt"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/server"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/util"
)

func main() {
  logger := util.CreateLogger()
  defer func() {
    _ = logger.Sync()
  }()

  dbUrl := flag.String("db", "127.0.0.1:9000", "The ClickHouse database URL.")
  natsUrl := flag.String("nats", "", "The NATS URL.")
  flag.Parse()

  err := server.Serve(*dbUrl, *natsUrl, logger)
  if err != nil {
    logger.Fatal(fmt.Sprintf("%+v", err))
  }
}
