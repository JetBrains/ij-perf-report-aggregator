package filling

import (
  "database/sql"
  "github.com/alecthomas/kingpin"
  "github.com/bvinc/go-sqlite-lite/sqlite3"
  "github.com/develar/errors"
  "github.com/json-iterator/go"
  "go.uber.org/zap"
  "report-aggregator/pkg/analyzer"
  "report-aggregator/pkg/model"
  "report-aggregator/pkg/sqlx"
  "report-aggregator/pkg/util"
  "strconv"
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

  // ZSTD 19 is used, read/write timeout should be quite large (10 minutes)
  db, err := sql.Open("clickhouse", "tcp://"+clickHouseUrl+"?read_timeout=600&write_timeout=600&compress=1")
  if err != nil {
    return errors.WithStack(err)
  }

  defer util.Close(db, logger)

  var lastMaxGeneratedTime time.Time
  err = db.QueryRow("select max(generated_time) from report").Scan(&lastMaxGeneratedTime)
  if err != nil {
    return errors.WithStack(err)
  }

  selectStatement, err := mainDb.Prepare(`
    select 
      product, machine.name as machine, generated_time, tc_build_id,
      duration_metrics, instant_metrics, 
      raw_report,
      build_c1, build_c2, build_c3
    from report 
    inner join machine on machine.rowid = report.machine
    where generated_time > ` + strconv.FormatInt(lastMaxGeneratedTime.Unix(), 10) + `
    order by product, machine, build_c1, build_c2, build_c3, generated_time`)

  if err != nil {
    return errors.WithStack(err)
  }

  defer util.Close(selectStatement, logger)

  err = model.CreateTable(db)
  if err != nil {
    return errors.WithStack(err)
  }

  var sb strings.Builder
  sb.WriteString(`insert into report values (`)

  for i := 0; i < 8; i++ {
    if i != 0 {
      sb.WriteRune(',')
    }
    sb.WriteRune('?')
  }
  model.ProcessMetricName(func(name string, isInstant bool) {
    sb.WriteString(", ?")
  })
  sb.WriteRune(')')

  insertManager := sqlx.NewBulkInsertManager(db, sb.String(), logger)
  defer util.Close(insertManager, logger)
  err = writeReports(insertManager, selectStatement, lastMaxGeneratedTime.Unix(), logger)
  if err != nil {
    return err
  }
  return nil
}

func writeReports(insertManager *sqlx.BulkInsertManager, selectFromOldStatement *sqlite3.Stmt, lastMaxGeneratedTime int64, logger *zap.Logger) error {
  selectStatement, err := insertManager.Db.Prepare("select 1 from report where product = ? and machine = ? and generated_time = ? limit 1")
  if err != nil {
    return errors.WithStack(err)
  }

  row := &MetricResult{}
  for {
    hasRow, err := selectFromOldStatement.Step()
    if err != nil {
      logger.Error("cannot step", zap.Error(err))
      return errors.WithStack(err)
    }

    if !hasRow {
      break
    }

    if row.generatedTime <= lastMaxGeneratedTime {
      fakeResult := -1
      err = selectStatement.QueryRow(row.productCode, row.machine, row.generatedTime).Scan(&fakeResult)
      if err != nil && err != sql.ErrNoRows {
        return errors.WithStack(err)
      }

      if err != sql.ErrNoRows {
        logger.Debug("report already processed", zap.Int64("generatedTime", row.generatedTime))
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

  args := []interface{}{row.productCode, row.machine, row.generatedTime, row.tcBuildId, row.rawReport, row.buildC1, row.buildC2, row.buildC3}
  model.ProcessMetricName(func(name string, isInstant bool) {
    var v int
    if isInstant {
      v = instantMetrics[name]
    } else {
      v = durationMetrics[name]
    }

    if !isInstant && v > 65535 {
      //if _, ok := model.MetricToUint16DataType[name]; ok {
      //if _, ok := model.MetricToUint16DataType[name]; ok {
        err = errors.Errorf("value outside of uint16 range (generatedTime: %d, value: %d)", row.generatedTime, v)
      //}
    }
    args = append(args, v)
  })

  if err != nil {
    return err
  }

  _, err = insertStatement.Exec(args...)
  return err
}
