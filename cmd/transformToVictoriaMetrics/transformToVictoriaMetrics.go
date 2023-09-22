package main

import (
  "bytes"
  "context"
  "fmt"
  "github.com/ClickHouse/clickhouse-go/v2"
  "github.com/ClickHouse/clickhouse-go/v2/lib/driver"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/analyzer"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/util"
  "github.com/develar/errors"
  "github.com/golang/protobuf/proto"
  "github.com/klauspost/compress/snappy"
  "github.com/prometheus/prometheus/prompb"
  "go.deanishe.net/env"
  "go.uber.org/zap"
  "log"
  "net/http"
  "strings"
  "time"
)

const victoriaMetricsURL = "http://localhost:8428/api/v1/write"

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

  db := env.Get("DB")
  table := env.Get("TABLE")
  split := strings.Split(env.Get("DB"), "_")
  if len(split) > 1 {
    table = split[1]
  }
  if db == "" || table == "" {
    logger.Fatal("Missing db or table, don't forget to set env variables DB and/or TABLE")
  }
  err = transform("127.0.0.1:9000", db, table, logger)
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
  TcBuildType        string `ch:"tc_build_type"`
  TriggeredBy        string `ch:"triggeredBy"`

  BuildC1 uint8  `ch:"build_c1"`
  BuildC2 uint16 `ch:"build_c2"`
  BuildC3 uint16 `ch:"build_c3"`

  ServiceName     []string `ch:"service.name"`
  ServiceStart    []uint32 `ch:"service.start"`
  ServiceDuration []uint32 `ch:"service.duration"`
  ServiceThread   []string `ch:"service.thread"`
  ServicePlugin   []string `ch:"service.plugin"`

  MeasuresName  []string `ch:"measures.name"`
  MeasuresValue []int32  `ch:"measures.value"`
  MeasuresType  []string `ch:"measures.type"`
}

const numWorkers = 2

func transform(clickHouseUrl string, idName string, tableName string, logger *zap.Logger) error {
  logger.Info("start transforming", zap.String("db", idName))

  split := strings.Split(idName, "_")
  dbName := idName
  if len(split) > 1 {
    dbName = split[0]
  }

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

  config := analyzer.GetAnalyzer(idName)

  // the whole select result in memory - so, limit
  var minTime time.Time
  var maxTime time.Time
  // use something like (now() - toIntervalMonth(1)) to test the transformer on a fresh data
  err = db.QueryRow(taskContext, "select min(generated_time) as min, max(generated_time) as max from "+tableName).Scan(&minTime, &maxTime)
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
    err = process(taskContext, db, config, current, next, tableName, logger)
    if err != nil {
      return err
    }

    current = next
  }

  if err != nil {
    return err
  }

  logger.Info("transforming finished")
  return nil
}

func worker(tasks chan ReportRow, logger *zap.Logger) {
  for row := range tasks {
    // The data processing code comes here
    for i, name := range row.MeasuresName {
      // Create the labels
      labels := []prompb.Label{
        {Name: "machine", Value: row.Machine},
        {Name: "build_time", Value: row.BuildTime.Format(time.RFC3339)},
        {Name: "generated_time", Value: row.GeneratedTime.Format(time.RFC3339)},
        {Name: "project", Value: row.Project},
        {Name: "tc_build_id", Value: fmt.Sprintf("%d", row.TcBuildId)},
        {Name: "tc_installer_build_id", Value: fmt.Sprintf("%d", row.TcInstallerBuildId)},
        {Name: "branch", Value: row.Branch},
        {Name: "tc_build_type", Value: row.TcBuildType},
        {Name: "triggeredBy", Value: row.TriggeredBy},
        {Name: "build_c1", Value: fmt.Sprintf("%d", row.BuildC1)},
        {Name: "build_c2", Value: fmt.Sprintf("%d", row.BuildC2)},
        {Name: "build_c3", Value: fmt.Sprintf("%d", row.BuildC3)},
      }

      // Convert your value to float64
      value := float64(row.MeasuresValue[i])

      // Get the timestamp in milliseconds
      timestamp := row.GeneratedTime.Unix() * 1000

      // Send to VictoriaMetrics
      err := sendToVictoriaMetrics(name, value, labels, timestamp)
      if err != nil {
        log.Printf("Failed to send data to VictoriaMetrics: %v", err)
      }
    }
  }
}

func process(taskContext context.Context, db driver.Conn, config analyzer.DatabaseConfiguration, startTime time.Time, endTime time.Time, tableName string, logger *zap.Logger, ) error {
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
    buildFields := ""
    if config.HasInstallerField {
      buildFields = "build_c1, build_c2, build_c3,"
    }
    installerFields := ""
    if config.HasInstallerField {
      installerFields = "tc_installer_build_id, " + buildFields
    }
    rawReportField := ""
    if config.HasRawReport {
      rawReportField = "raw_report,"
    }
    rows, err = db.Query(taskContext, `
      select machine, branch,
             generated_time, build_time, `+rawReportField+`
             tc_build_id,`+installerFields+` project, measures.name, measures.value, measures.type, triggeredBy
      from `+tableName+`
      where generated_time >= $1 and generated_time < $2
      order by machine, branch, project, `+buildFields+` build_time, generated_time
    `, startTime, endTime)

  }
  if err != nil {
    return errors.WithStack(err)
  }

  defer util.Close(rows, logger)

  tasks := make(chan ReportRow, numWorkers)
  for i := 0; i < numWorkers; i++ {
    go worker(tasks, logger)
  }

  var row ReportRow
  for rows.Next() {
    err = rows.ScanStruct(&row)
    if err != nil {
      return errors.WithStack(err)
    }

    tasks <- row
  }
  close(tasks)

  err = rows.Err()
  if err != nil {
    return errors.WithStack(err)
  }

  return nil
}

func sendToVictoriaMetrics(metricName string, value float64, labels []prompb.Label, timestamp int64) error {
  // Create a sample for the metric
  labels = append([]prompb.Label{{Name: "__name__", Value: metricName}}, labels...)

  sample := &prompb.Sample{
    Value:     value,
    Timestamp: timestamp,
  }

  // Create a TimeSeries with your labels and sample
  ts := &prompb.TimeSeries{
    Labels:  labels,
    Samples: []prompb.Sample{*sample},
  }

  // Create a WriteRequest with the TimeSeries
  req := &prompb.WriteRequest{
    Timeseries: []prompb.TimeSeries{*ts},
  }

  // Serialize WriteRequest to protobuf
  data, err := proto.Marshal(req)
  if err != nil {
    return err
  }

  // Compress data using Snappy
  compressedData := snappy.Encode(nil, data)

  // Send to VictoriaMetrics
  resp, err := http.Post(victoriaMetricsURL, "application/x-protobuf", bytes.NewReader(compressedData))
  if err != nil {
    return err
  }
  defer resp.Body.Close()

  if resp.StatusCode != http.StatusNoContent {
    return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
  }

  return nil
}
