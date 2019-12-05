package main

import (
  "context"
  "fmt"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/analyzer"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/util"
  "github.com/develar/errors"
  "github.com/jmoiron/sqlx"
  "go.uber.org/zap"
  "time"
)

/*

1. Download backup, unpack to b and execute:
rm -rf ~/ij-perf-db/clickhouse
mkdir -p ~/ij-perf-db/clickhouse/metadata && mkdir ~/ij-perf-db/clickhouse/data && mv ~/Downloads/b/metadata/default ~/ij-perf-db/clickhouse/metadata/default && mv ~/Downloads/b/shadow/default ~/ij-perf-db/clickhouse/data/default

2. change `migrate/report.sql` as needed and execute.

*/
func main() {
  logger := util.CreateLogger()
  defer func() {
    _ = logger.Sync()
  }()

  err := transform("localhost:9000", logger)
  if err != nil {
    logger.Fatal(fmt.Sprintf("%+v", err))
  }
}

type ReportRow struct {
  Product uint8
  Machine uint8
  Branch  uint8

  GeneratedTime int64 `db:"generated_time"`
  BuildTime     int64     `db:"build_time"`

  RawReport []byte `db:"raw_report"`

  TcBuildId          int    `db:"tc_build_id"`
  TcInstallerBuildId int    `db:"tc_installer_build_id"`
  TcBuildProperties  []byte `db:"tc_build_properties"`

  BuildC1 int `db:"build_c1"`
  BuildC2 int `db:"build_c2"`
  BuildC3 int `db:"build_c3"`

  Project string `db:"project"`
}

type TimeRange struct {
  Min time.Time
  Max time.Time
}

// set insertWorkerCount to 1 if not enough memory
const insertWorkerCount = 2

func transform(clickHouseUrl string, logger *zap.Logger) error {
  db, err := sqlx.Open("clickhouse", "tcp://"+clickHouseUrl+"?read_timeout=600&write_timeout=600&debug=0&compress=1&send_timeout=30000&receive_timeout=3000")
  if err != nil {
    return errors.WithStack(err)
  }

  taskContext, cancel := util.CreateCommandContext()
  defer cancel()

  insertReportManager, err := analyzer.NewInsertReportManager(db, taskContext, "report2", insertWorkerCount, logger)
  if err != nil {
    return err
  }

  // reduce batch size if not enough memory
  //insertReportManager.InsertManager.BatchSize = 2000

  //var products []string
  //db.SelectContext(taskContext, &products, "select distinct product from report")

  // the whole select result in memory - so, limit
  var timeRange TimeRange
  err = db.GetContext(taskContext, &timeRange, "select min(generated_time) as min, max(generated_time) as max from report")
  if err != nil {
    return errors.WithStack(err)
  }

  logger.Info("time range", zap.Time("start", timeRange.Min), zap.Time("end", timeRange.Max))

  // 4 weeks
  selectDuration := time.Hour * 24 * 7 * 4
  for current := timeRange.Min; current.Before(timeRange.Max); {
    next := current.Add(selectDuration)
    err = process(db, current, next, insertReportManager, taskContext, logger)
    if err != nil {
      return err
    }

    current = next
  }

  err = insertReportManager.InsertManager.CommitAndWait()
  if err != nil {
    return err
  }

  err = insertReportManager.Close()
  if err != nil {
    return err
  }

  return nil
}

func process(db *sqlx.DB, startTime time.Time, endTime time.Time, insertReportManager *analyzer.InsertReportManager, taskContext context.Context, logger *zap.Logger) error {
  logger.Info("process", zap.Time("start", startTime), zap.Time("end", endTime))
  // don't forget to update order clause if differs - better to insert data in an expected order
  rows, err := db.QueryxContext(taskContext, `
    select cast(product, 'UInt8') as product, cast(machine, 'UInt8') as machine, cast(branch, 'UInt8') as branch, 
           toUnixTimestamp(generated_time) as generated_time, toUnixTimestamp(build_time) as build_time, raw_report, 
           tc_build_id, tc_installer_build_id, tc_build_properties,
           build_c1, build_c2, build_c3, project
    from report
    where generated_time >= ? and generated_time < ?
    order by product, machine, branch, build_c1, build_c2, build_c3, project, build_time, generated_time
  `, startTime, endTime)
  if err != nil {
    return errors.WithStack(err)
  }

  defer util.Close(rows, logger)

  var row ReportRow
  for rows.Next() {
    err = rows.StructScan(&row)
    if err != nil {
      return errors.WithStack(err)
    }

    reportRow := &analyzer.MetricResult{
      RawReport: row.RawReport,

      Machine: row.Machine,

      GeneratedTime: row.GeneratedTime,
      BuildTime:     row.BuildTime,

      TcBuildId:          row.TcBuildId,
      TcInstallerBuildId: row.TcInstallerBuildId,
      TcBuildProperties:  row.TcBuildProperties,

      BuildC1: row.BuildC1,
      BuildC2: row.BuildC2,
      BuildC3: row.BuildC3,
    }

    err = insertReportManager.WriteMetrics(row.Product, reportRow, row.Branch, row.Project, logger)
    if err != nil {
      return err
    }
  }

  err = rows.Err()
  if err != nil {
    return errors.WithStack(err)
  }

  return nil
}
