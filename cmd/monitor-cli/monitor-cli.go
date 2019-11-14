package main

import (
  "fmt"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/server"
  "github.com/develar/errors"
  "github.com/jmoiron/sqlx"
  "log"
  "os"
)

// kubectl port-forward svc/clickhouse 9900:9000

func main() {
  err := analyzeTotal("127.0.0.1:9900", "master", "2019-11-04 00:00:00")
  if err != nil {
    log.Fatal(fmt.Sprintf("%+v", err))
  }
}

func analyzeTotal(dbUrl string, branch string, goldWeekStart string) error {
  if len(dbUrl) == 0 {
    dbUrl = "127.0.0.1:9000"
  }

  db, err := sqlx.Open("clickhouse", "tcp://"+dbUrl+"?compress=1")
  if err != nil {
    return errors.WithStack(err)
  }

  result, err := server.CompareMetrics(db, branch, goldWeekStart)
  if err != nil {
    return err
  }

  server.PrintResult(*result, os.Stdout)
  return nil
}
