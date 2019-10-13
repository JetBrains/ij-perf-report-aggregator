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
  "report-aggregator/pkg/sqlx"
  "report-aggregator/pkg/util"
  "strings"
  "time"

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
    order by product, machine, build_c1, build_c2, build_c3, generated_time
  `)

  if err != nil {
    return errors.WithStack(err)
  }

  defer util.Close(selectStatement, logger)

  // ZSTD 22 is used, read/write timeout should be quite large (10 minutes)
  db, err := sql.Open("clickhouse", "tcp://"+clickHouseUrl+"?read_timeout=600&write_timeout=600&compress=1")
  if err != nil {
    return errors.WithStack(err)
  }

  defer util.Close(db, logger)

  _, err = db.Exec("SET allow_experimental_data_skipping_indices = 1")
  if err != nil {
    return errors.WithStack(err)
  }

  // https://www.altinity.com/blog/2019/7/new-encodings-to-improve-clickhouse
  // see zstd-compression-level.txt
  var sb strings.Builder
  sb.WriteString(`
    create table if not exists report (
      id FixedString(27) Codec(ZSTD(19)),
      INDEX id_index id type minmax granularity 1,

      product FixedString(2) Codec(ZSTD(19)),
      machine String Codec(ZSTD(19)),
      generated_time DateTime Codec(Delta, ZSTD(19)),
      
      raw_report String Codec(ZSTD(19)),
      
      build_c1 UInt8 Codec(DoubleDelta, ZSTD(19)),
      build_c2 UInt16 Codec(DoubleDelta, ZSTD(19)),
      build_c3 UInt16 Codec(DoubleDelta, ZSTD(19)),

      metrics_version UInt8 Codec(Delta, ZSTD(19)),
      index metrics_version_index metrics_version type minmax granularity 1
`)
  server.ProcessMetricName(func(name string, isInstant bool) {
    sb.WriteRune(',')
    sb.WriteRune(' ')
    sb.WriteString(name)
    sb.WriteRune('_')
    if isInstant {
      sb.WriteRune('i')
    } else {
      sb.WriteRune('d')
    }
    sb.WriteString(" Int32 Codec(Gorilla, ZSTD(19))")
  })

  // https://github.com/ClickHouse/ClickHouse/issues/3758#issuecomment-444490724
  sb.WriteString(") engine MergeTree partition by (product, toYYYYMM(generated_time)) order by (product, machine, build_c1, build_c2, build_c3, generated_time) SETTINGS old_parts_lifetime = 10")

  _, err = db.Exec(sb.String())
  if err != nil {
    return errors.WithStack(err)
  }

  sb.Reset()
  sb.WriteString(`insert into report values (`)

  for i := 0; i < 9; i++ {
    if i != 0 {
      sb.WriteRune(',')
    }
    sb.WriteRune('?')
  }
  server.ProcessMetricName(func(name string, isInstant bool) {
    sb.WriteString(", ?")
  })
  sb.WriteRune(')')

  err = writeReports(sqlx.NewBulkInsertManager(db, sb.String(), logger), selectStatement, logger)
  if err != nil {
    return err
  }
  return nil
}

func writeReports(insertManager *sqlx.BulkInsertManager, selectFromOldStatement *sqlite3.Stmt, logger *zap.Logger) error {
  defer util.Close(insertManager, logger)

  selectStatement, err := insertManager.Db.Prepare("select metrics_version from report where id = ? limit 1")
  if err != nil {
    return errors.WithStack(err)
  }


  var lastMaxGeneratedTime time.Time
  err = insertManager.Db.QueryRow("select max(generated_time) from report").Scan(&lastMaxGeneratedTime)
  if err != nil {
    return errors.WithStack(err)
  }

  // large inserts leads to large memory usage, so, insert by 500 items
  row := &MetricResult{}
  for {
    hasRow, err := selectFromOldStatement.Step()
    if !hasRow {
      break
    }

    if err != nil {
      logger.Error("cannot step", zap.Error(err))
      return err
    }

    currentMetricsVersion := -1
    if row.generatedTime <= lastMaxGeneratedTime.Unix() {
      err = selectStatement.QueryRow(row.id).Scan(&currentMetricsVersion)
      if err != nil && err != sql.ErrNoRows {
        return errors.WithStack(err)
      }

      if currentMetricsVersion == analyzer.MetricsVersion {
        logger.Debug("report already processed")
        continue
      }
    }

    err = insertManager.PrepareForInsert()
    if err != nil {
      return err
    }

    err = writeMetrics(selectFromOldStatement, row, insertManager.InsertStatement, logger)
    if err != nil {
      return err
    }
  }
  return insertManager.Commit()
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

  args := []interface{}{row.id, row.productCode, row.machine, row.generatedTime, row.rawReport, row.buildC1, row.buildC2, row.buildC3, analyzer.MetricsVersion}
  server.ProcessMetricName(func(name string, isInstant bool) {
    if isInstant {
      args = append(args, instantMetrics[name])
    } else {
      args = append(args, durationMetrics[name])
    }
  })
  _, err = insertStatement.Exec(args...)
  return err
}