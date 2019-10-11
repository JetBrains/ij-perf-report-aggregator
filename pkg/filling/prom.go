package filling

import (
  "database/sql"
  "github.com/alecthomas/kingpin"
  "github.com/bvinc/go-sqlite-lite/sqlite3"
  "github.com/develar/errors"
  "github.com/json-iterator/go"
  "go.uber.org/zap"
  "report-aggregator/pkg/analyzer"
  "report-aggregator/pkg/server"
  "report-aggregator/pkg/util"
  "strings"

  _ "github.com/kshvakov/clickhouse"
)

func ConfigureFillCommand(app *kingpin.Application, logger *zap.Logger) {
  command := app.Command("fill", "Fill VictoriaMetrics database using SQLite database.")
  dbPath := command.Flag("db", "The SQLite database file.").Required().String()
  clickHouseUrl := command.Flag("clickHouse", "The ClickHouse server URL.").Required().String()
  updateMetrics := command.Flag("update-metrics", "Whether to update computed metrics if outdated. Think about backup.").Bool()
  command.Action(func(context *kingpin.ParseContext) error {
    return fill(*dbPath, *updateMetrics, *clickHouseUrl, logger)
  })
}

func fill(dbPath string, updateMetrics bool, clickHouseUrl string, logger *zap.Logger) error {
  if updateMetrics {
    err := analyzer.UpdateMetrics(dbPath, logger)
    if err != nil {
      return err
    }
  }

  mainDb, err := sqlite3.Open(dbPath, sqlite3.OPEN_READONLY)
  if err != nil {
    return errors.WithStack(err)
  }

  defer util.Close(mainDb, logger)

  selectStatement, err := mainDb.Prepare(`
    select 
      id, product, machine.name as machine, generated_time, 
      duration_metrics, instant_metrics, 
      raw_report,
      build_c1, build_c2, build_c3
    from report 
    inner join machine on machine.rowid = report.machine 
    order by generated_time
  `)

  if err != nil {
    return errors.WithStack(err)
  }

  defer util.Close(selectStatement, logger)

  // ZSTD 22 is used, read/write timeout should be quite large (10 minutes)
  db, err := sql.Open("clickhouse", "tcp://"+clickHouseUrl+"?compress=1")
  if err != nil {
    return errors.WithStack(err)
  }

  defer util.Close(db, logger)

  tx, err := db.Begin()
  if err != nil {
    return errors.WithStack(err)
  }

  // https://www.altinity.com/blog/2019/7/new-encodings-to-improve-clickhouse
  var sb strings.Builder
  sb.WriteString(`
    CREATE TABLE report (
      id FixedString(27) Codec(ZSTD(22)),
      product FixedString(2) Codec(ZSTD(22)),
      machine String Codec(ZSTD(22)),
      generated_time DateTime Codec(Delta, ZSTD(22)),
      
      raw_report String Codec(ZSTD(22)),
      
      build_c1 UInt8 Codec(DoubleDelta, ZSTD(22)),
      build_c2 UInt16 Codec(DoubleDelta, ZSTD(22)),
      build_c3 UInt16 Codec(DoubleDelta, ZSTD(22))`)
  processMetricName(func(name string, isInstant bool) {
    sb.WriteRune(',')
    sb.WriteRune(' ')
    sb.WriteString(name)
    sb.WriteRune('_')
    if isInstant {
      sb.WriteRune('i')
    } else {
      sb.WriteRune('d')
    }
    sb.WriteString(" Int32 Codec(Gorilla, ZSTD(22))")
  })
  sb.WriteString(") engine MergeTree partition by toYYYYMM(generated_time) order by (generated_time, product, machine, build_c1, build_c2, build_c3)")

  _, err = db.Exec(sb.String())
  if err != nil {
    return errors.WithStack(err)
  }

  sb.Reset()
  sb.WriteString(`INSERT INTO report VALUES (`)
  for i := 0; i < 8; i++ {
    if i != 0 {
      sb.WriteRune(',')
    }
    sb.WriteRune('?')
  }
  processMetricName(func(name string, isInstant bool) {
    sb.WriteString(", ?")
  })
  sb.WriteRune(')')
  insertStatement, err := tx.Prepare(sb.String())
  if err != nil {
    return err
  }

  defer util.Close(insertStatement, logger)

  err = writeReports(selectStatement, insertStatement, logger)
  if err != nil {
    return err
  }

  err = tx.Commit()
  if err != nil {
    return err
  }
  return nil
}

func writeReports(selectStatement *sqlite3.Stmt, insertStatement *sql.Stmt, logger *zap.Logger) error {
  row := &MetricResult{}
  for {
    hasRow, err := selectStatement.Step()
    if !hasRow {
      return nil
    }

    if err != nil {
      logger.Error("cannot step", zap.Error(err))
      return err
    }

    err = writeMetrics(selectStatement, row, insertStatement, logger)
    if err != nil {
      return err
    }
  }
}

func writeMetrics(selectStatement *sqlite3.Stmt, row *MetricResult, insertStatement *sql.Stmt, logger *zap.Logger) error {
  err := scanMetricResult(selectStatement, row)
  if err != nil {
    logger.Error("cannot scan", zap.Error(err))
    return err
  }

  var durationMetrics map[string]int
  err = jsoniter.ConfigFastest.Unmarshal([]byte(row.durationMetricsJson), &durationMetrics)
  if err != nil {
    return errors.WithStack(err)
  }

  var instantMetrics map[string]int
  err = jsoniter.ConfigFastest.Unmarshal([]byte(row.instantMetricsJson), &instantMetrics)
  if err != nil {
    return errors.WithStack(err)
  }

  args := []interface{}{row.id, row.productCode, row.machine, row.generatedTime, row.rawReport, row.buildC1, row.buildC2, row.buildC3}
  processMetricName(func(name string, isInstant bool) {
    if isInstant {
      args = append(args, instantMetrics[name])
    } else {
      args = append(args, durationMetrics[name])
    }
  })
  _, err = insertStatement.Exec(args...)
  return err
}

func processMetricName(handler func(string, bool)) {
  for _, name := range server.EssentialDurationMetricNames {
    handler(name, false)
  }
  handler("moduleLoading", false)
  for _, name := range server.InstantMetricNames {
    handler(name, true)
  }
}
