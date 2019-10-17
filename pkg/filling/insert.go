package filling

import (
  "database/sql"
  "github.com/alecthomas/kingpin"
  "github.com/bvinc/go-sqlite-lite/sqlite3"
  "github.com/develar/errors"
  "github.com/mcuadros/go-version"
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
  command.Action(func(context *kingpin.ParseContext) error {
    return fill(*dbPath, *clickHouseUrl, logger)
  })
}

func fill(dbPath string, clickHouseUrl string, logger *zap.Logger) error {
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

  err = model.CreateTable(db)
  if err != nil {
    return errors.WithStack(err)
  }

  var lastMaxGeneratedTime time.Time
  err = db.QueryRow("select max(generated_time) from report").Scan(&lastMaxGeneratedTime)
  if err != nil {
    return errors.WithStack(err)
  }

  selectStatement, err := mainDb.Prepare(`
    select 
      product, machine.name as machine, generated_time, tc_build_id,
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

  insertManager, err := sqlx.NewBulkInsertManager(db, sb.String(), logger)
  if err != nil {
    return err
  }

  defer util.Close(insertManager, logger)
  err = writeReports(insertManager, selectStatement, lastMaxGeneratedTime.Unix(), logger)
  if err != nil {
    return err
  }

  logger.Info("waiting inserting", zap.Int("transactions", insertManager.GetUncommittedTransactionCount()))
  insertManager.WaitGroup.Wait()
  return insertManager.Error
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

    err = scanMetricResult(selectFromOldStatement, row)
    if err != nil {
      logger.Error("cannot scan", zap.Error(err))
      return err
    }

    err = writeMetrics(row, insertManager.InsertStatement, logger)
    if err != nil {
      return err
    }
  }
  return insertManager.Commit()
}

func writeMetrics(row *MetricResult, insertStatement *sql.Stmt, logger *zap.Logger) error {
  report, err := analyzer.ReadReport([]byte(row.rawReport))
  if err != nil {
    return errors.WithStack(err)
  }

  durationMetrics, instantMetrics := analyzer.ComputeMetrics(report, logger)
  // or both null, or not - no need to check each one
  if durationMetrics == nil || instantMetrics == nil {
    return errors.New("metrics cannot be computed")
  }

  var buildTimeUnix int64
  if version.Compare(report.Version, "13", ">=") {
    buildTime, err := analyzer.ParseTime(report.BuildDate)
    if err != nil {
      return err
    }
    buildTimeUnix = buildTime.Unix()
  } else {
    buildTimeUnix = 0
  }

  args := []interface{}{row.productCode, row.machine, buildTimeUnix, row.generatedTime, row.tcBuildId, row.rawReport, row.buildC1, row.buildC2, row.buildC3}
  for _, name := range model.DurationMetricNames {
    var v int
    switch name {
    case "bootstrap":
      v = durationMetrics.Bootstrap
    case "appInitPreparation":
      v = durationMetrics.AppInitPreparation
    case "appInit":
      v = durationMetrics.AppInit
    case "pluginDescriptorLoading":
      v = durationMetrics.PluginDescriptorLoading
    case "appComponentCreation":
      v = durationMetrics.AppComponentCreation
    case "projectComponentCreation":
      v = durationMetrics.ProjectComponentCreation
    case "moduleLoading":
      v = durationMetrics.ModuleLoading
    default:
      return errors.New("unknown metric " + name)
    }

    if v > 65535 {
      return errors.Errorf("value outside of uint16 range (generatedTime: %d, value: %d)", row.generatedTime, v)
    }

    args = append(args, v)
  }

  for _, name := range model.InstantMetricNames {
    var v int
    switch name {
    case "splash":
      v = instantMetrics.Splash
      if v < 0 {
        continue
      }
    case "startUpCompleted":
      v = instantMetrics.StartUpCompleted
    default:
      return errors.New("unknown metric " + name)
    }

    args = append(args, v)
  }

  _, err = insertStatement.Exec(args...)
  return err
}
