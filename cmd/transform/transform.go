package main

import (
  "context"
  "fmt"
  "github.com/ClickHouse/clickhouse-go/v2"
  "github.com/ClickHouse/clickhouse-go/v2/lib/driver"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/analyzer"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/util"
  "github.com/develar/errors"
  "go.deanishe.net/env"
  "go.uber.org/zap"
  "log"
  "strings"
  "time"
)

/*
1. run restore-backup RC
2. change `migrate/report.sql` as needed and execute.
*/
func main() {
  config := zap.NewDevelopmentConfig()
  config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
  config.DisableCaller = true
  config.DisableStacktrace = true
  logger, err := config.Build()
  if err != nil {
    log.Fatal(err)
  }

  defer func() {
    _ = logger.Sync()
  }()

  err = transform("localhost:9000", env.Get("DB"), logger)
  if err != nil {
    logger.Fatal(fmt.Sprintf("%+v", err))
  }
}

type ReportRow struct {
  Product string `ch:"product"`
  Machine string `ch:"machine"`
  Branch  string `ch:"branch"`

  Project string `ch:"project"`

  GeneratedTime time.Time `ch:"generated_time"`
  BuildTime     time.Time `ch:"build_time"`

  RawReport string `ch:"raw_report"`

  TcBuildId          uint32 `ch:"tc_build_id"`
  TcInstallerBuildId uint32 `ch:"tc_installer_build_id"`

  BuildC1 uint8  `ch:"build_c1"`
  BuildC2 uint16 `ch:"build_c2"`
  BuildC3 uint16 `ch:"build_c3"`

  ServiceName     []string `ch:"service.name"`
  ServiceStart    []uint32 `ch:"service.start"`
  ServiceDuration []uint32 `ch:"service.duration"`
  ServiceThread   []string `ch:"service.thread"`
  ServicePlugin   []string `ch:"service.plugin"`
}

// set insertWorkerCount to 1 if not enough memory
const insertWorkerCount = 4

func transform(clickHouseUrl string, dbName string, logger *zap.Logger) error {
  logger.Info("start transforming", zap.String("db", dbName))

  db, err := clickhouse.Open(&clickhouse.Options{
    Addr: []string{clickHouseUrl},
    Auth: clickhouse.Auth{
      Database: dbName,
    },
    DialTimeout:     time.Second,
    ConnMaxLifetime: time.Hour,
    Settings: map[string]interface{}{
      // https://github.com/ClickHouse/ClickHouse/issues/2833
      // ZSTD 19+ is used, read/write timeout should be quite large (10 minutes)
      "send_timeout":     30_000,
      "receive_timeout":  3000,
      "max_memory_usage": 100000000000,
    },
  })
  if err != nil {
    return errors.WithStack(err)
  }

  defer util.Close(db, logger)

  taskContext, cancel := util.CreateCommandContext()
  defer cancel()

  config := analyzer.GetAnalyzer(dbName)

  insertReportManager, err := analyzer.NewInsertReportManager(db, config, taskContext, "report2", insertWorkerCount, logger)
  if err != nil {
    return err
  }

  insertManager := insertReportManager.InsertManager
  // we send batch in the end of each iteration
  insertManager.BatchSize = 50_000

  // the whole select result in memory - so, limit
  var minTime time.Time
  var maxTime time.Time
  // use something like (now() - toIntervalMonth(1)) to test the transformer on a fresh data
  err = db.QueryRow(taskContext, "select min(generated_time) as min, max(generated_time) as max from report").Scan(&minTime, &maxTime)
  if err != nil {
    return errors.WithStack(err)
  }

  logger.Info("time range", zap.Time("start", minTime), zap.Time("end", maxTime))

  // round to the start of the month
  minTime = time.Date(minTime.Year(), minTime.Month(), 1, 0, 0, 0, 0, minTime.Location())
  // round to the end of the month
  if maxTime.Month() == 12 {
    maxTime = time.Date(maxTime.Year()+1, 1, 1, 0, 0, 0, 0, maxTime.Location())
  } else {
    maxTime = time.Date(maxTime.Year(), maxTime.Month()+1, 1, 0, 0, 0, 0, maxTime.Location())
  }

  for current := minTime; current.Before(maxTime); {
    // 1 month
    next := current.AddDate(0, 1, 0)
    err = process(db, config, current, next, insertReportManager, taskContext, logger)
    if err != nil {
      return err
    }

    current = next

    if insertManager.GetQueuedItemCount() > 10_000 {
      insertManager.ScheduleSendBatch()
    }
  }

  err = insertReportManager.InsertManager.Close()
  if err != nil {
    return err
  }

  logger.Info("transforming finished")
  return nil
}

func process(
  db driver.Conn,
  config analyzer.DatabaseConfiguration,
  startTime time.Time,
  endTime time.Time,
  insertReportManager *analyzer.InsertReportManager,
  taskContext context.Context,
  logger *zap.Logger,
) error {
  logger.Info("process", zap.Time("start", startTime), zap.Time("end", endTime))
  // don't forget to update order clause if differs - better to insert data in an expected order

  var err error
  var rows driver.Rows
  if config.HasProductField {
    rows, err = db.Query(taskContext, `
      select product, machine, branch,
             generated_time, build_time, raw_report,
             tc_build_id, tc_installer_build_id,
             build_c1, build_c2, build_c3, project,
             service.name, service.start, service.duration, service.thread, service.plugin
      from report
      where generated_time >= $1 and generated_time < $2
      order by product, machine, branch, project, build_c1, build_c2, build_c3, build_time, generated_time
    `, startTime, endTime)
  } else {
    rows, err = db.Query(taskContext, `
      select machine, branch,
             generated_time, build_time, raw_report,
             tc_build_id, tc_installer_build_id,
             build_c1, build_c2, build_c3, project
      from report
      where generated_time >= $1 and generated_time < $2
      order by machine, branch, project, build_c1, build_c2, build_c3, build_time, generated_time
    `, startTime, endTime)
  }
  if err != nil {
    return errors.WithStack(err)
  }

  defer util.Close(rows, logger)

  var row ReportRow
rowLoop:
  for rows.Next() {
    err = rows.ScanStruct(&row)
    if err != nil {
      return errors.WithStack(err)
    }

    runResult := &analyzer.RunResult{
      RawReport: []byte(row.RawReport),

      Machine: row.Machine,

      GeneratedTime: row.GeneratedTime,
      BuildTime:     row.BuildTime,

      TcBuildId:          int(row.TcBuildId),
      TcInstallerBuildId: int(row.TcInstallerBuildId),

      BuildC1: int(row.BuildC1),
      BuildC2: int(row.BuildC2),
      BuildC3: int(row.BuildC3),
    }

    err = analyzer.ReadReport(runResult, config, logger)
    if err != nil {
      return err
    }

    if runResult.Report == nil {
      // ignore report
      continue rowLoop
    }

    if config.HasProductField {
      runResult.ExtraFieldData[0] = row.ServiceName
      runResult.ExtraFieldData[1] = row.ServiceStart
      runResult.ExtraFieldData[2] = row.ServiceDuration
      runResult.ExtraFieldData[3] = row.ServiceThread
      runResult.ExtraFieldData[4] = row.ServicePlugin
    }

    if strings.HasPrefix(row.Project, "2tI") || strings.HasPrefix(row.Project, "dEQ") {
      continue rowLoop
    }
    if row.Project == "73YWaW9bytiPDGuKvwNIYMK5CKI" {
      runResult.Report.Project = "simple for IJ"
    }
    if strings.Contains(row.RawReport, "modules loading with cache") {
      if row.Project == "A/vsu8PQGaeUtpfB1yz/I4EwVnI" {
        runResult.Report.Project = "open-telemetry - gradle from cache"
      } else if row.Project == "cVqhfTfTzoDOzZx2ZbLSKxC2TpM" {
        runResult.Report.Project = "gradle-500-modules - from cache"
      }
    } else if strings.Contains(row.RawReport, "modules loading without cache") {
      if row.Project == "A/vsu8PQGaeUtpfB1yz/I4EwVnI" {
        runResult.Report.Project = "open-telemetry - gradle without cache"
      } else if row.Project == "cVqhfTfTzoDOzZx2ZbLSKxC2TpM" {
        runResult.Report.Project = "gradle-500-modules - without cache"
      }
    }

    err = insertReportManager.WriteMetrics(row.Product, runResult, row.Branch, row.Project, logger)
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
