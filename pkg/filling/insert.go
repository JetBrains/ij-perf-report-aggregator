package filling

import (
  "database/sql"
  "github.com/alecthomas/kingpin"
  "github.com/develar/errors"
  "github.com/iancoleman/strcase"
  "github.com/jmoiron/sqlx"
  "go.uber.org/zap"
  "report-aggregator/pkg/analyzer"
  "report-aggregator/pkg/model"
  "report-aggregator/pkg/sql-util"
  "report-aggregator/pkg/util"
  "strconv"
  "strings"
  "time"

  _ "github.com/kshvakov/clickhouse"
  _ "github.com/mattn/go-sqlite3"
)

type MetricResult struct {
  Product string
  Machine uint8

  GeneratedTime int64
  BuildTime     int64

  TcBuildId          int
  TcInstallerBuildId int
  TcBuildProperties  []byte

  Branch sql.RawBytes

  RawReport string

  BuildC1 int `db:"build_c1"`
  BuildC2 int `db:"build_c2"`
  BuildC3 int `db:"build_c3"`
}

func ConfigureFillCommand(app *kingpin.Application, logger *zap.Logger) {
  command := app.Command("fill", "Fill ClickHouse database using SQLite database.")
  dbPath := command.Flag("db", "The SQLite database file.").Required().String()
  clickHouseUrl := command.Flag("clickHouse", "The ClickHouse server URL.").Required().String()
  command.Action(func(context *kingpin.ParseContext) error {
    return fill(*dbPath, *clickHouseUrl, logger)
  })
}

func fill(dbPath string, clickHouseUrl string, logger *zap.Logger) error {
  mainDb, err := sqlx.Open("sqlite3", "file:"+dbPath+"?mode=ro")
  if err != nil {
    return errors.WithStack(err)
  }

  mainDb.MapperFunc(strcase.ToSnake)

  defer util.Close(mainDb, logger)

  // ZSTD 19 is used, read/write timeout should be quite large (10 minutes)
  db, err := sqlx.Open("clickhouse", "tcp://"+clickHouseUrl+"?read_timeout=600&write_timeout=600&compress=1")
  if err != nil {
    return errors.WithStack(err)
  }

  defer util.Close(db, logger)

  var products []model.IdAndName
  err = mainDb.Select(&products, "select row_number() over (order by product) AS id, product as name from report group by product")
  if err != nil {
    return errors.WithStack(err)
  }

  var machines []model.IdAndName
  err = mainDb.Select(&machines, "select ROWID as id, name from machine order by id")
  if err != nil {
    return errors.WithStack(err)
  }

  err = copyInstallers(mainDb, db, logger)
  if err != nil {
    return err
  }

  err = model.CreateTable(db, machines, products)
  if err != nil {
    return errors.WithStack(err)
  }

  var lastMaxGeneratedTime time.Time
  err = db.QueryRow("select max(generated_time) from report").Scan(&lastMaxGeneratedTime)
  if err != nil {
    return errors.WithStack(err)
  }

  sourceRows, err := mainDb.Queryx(`
    select 
      product, machine, generated_time, build_time,
      tc_build_id, tc_installer_build_id, tc_build_properties,
      json_extract(tc_build_properties, '$."vcsroot.branch"') as branch,
      raw_report,
      build_c1, build_c2, build_c3
    from report 
    where generated_time > ?
    order by product, machine, build_c1, build_c2, build_c3, generated_time`, strconv.FormatInt(lastMaxGeneratedTime.Unix(), 10))
  if err != nil {
    return errors.WithStack(err)
  }

  defer util.Close(sourceRows, logger)

  var sb strings.Builder
  sb.WriteString(`insert into report values (`)

  for i := 0; i < 11; i++ {
    if i != 0 {
      sb.WriteRune(',')
    }
    sb.WriteRune('?')
  }
  model.ProcessMetricName(func(name string, isInstant bool) {
    sb.WriteString(", ?")
  })
  sb.WriteRune(')')

  insertManager, err := sql_util.NewBulkInsertManager(db, sb.String(), logger)
  if err != nil {
    return err
  }

  logger.Info("start inserting", zap.String("clickHouseUrl", clickHouseUrl))

  defer util.Close(insertManager, logger)
  err = writeReports(insertManager, sourceRows, lastMaxGeneratedTime.Unix(), logger)
  if err != nil {
    return err
  }

  logger.Info("waiting inserting", zap.Int("transactions", insertManager.GetUncommittedTransactionCount()))
  insertManager.WaitGroup.Wait()
  return insertManager.Error
}

func copyInstallers(sourceDb *sqlx.DB, db *sqlx.DB, logger *zap.Logger) error {
  err := model.CreateInstallerTable(db)
  if err != nil {
    return err
  }

  //noinspection SqlResolve
  rows, err := sourceDb.Query("select id, changes from installer order by id")
  if err != nil {
    return errors.WithStack(err)
  }

  insertManager, err := sql_util.NewInstallerManager(db, logger)
  if err != nil {
    return err
  }
  defer util.Close(insertManager, logger)

  for rows.Next() {
    var id int
    var changes sql.RawBytes
    err = rows.Scan(&id, &changes)
    if err != nil {
      return errors.WithStack(err)
    }

    err = insertManager.Insert(id, string(changes))
    if err != nil {
      return err
    }
  }

  err = rows.Err()
  if err != nil {
    return errors.WithStack(err)
  }

  err = insertManager.Commit()
  if err != nil {
    return errors.WithStack(err)
  }
  return nil
}

func writeReports(insertManager *sql_util.BulkInsertManager, sourceRows *sqlx.Rows, lastMaxGeneratedTime int64, logger *zap.Logger) error {
  selectStatement, err := insertManager.Db.Prepare("select 1 from report where product = ? and machine = ? and generated_time = ? limit 1")
  if err != nil {
    return errors.WithStack(err)
  }

  row := &MetricResult{}
  for sourceRows.Next() {
    err = sourceRows.StructScan(row)
    if err != nil {
      logger.Error("cannot scan", zap.Error(err))
      return err
    }

    if row.GeneratedTime <= lastMaxGeneratedTime {
      fakeResult := -1
      err = selectStatement.QueryRow(row.Product, row.Machine, row.GeneratedTime).Scan(&fakeResult)
      if err != nil && err != sql.ErrNoRows {
        return errors.WithStack(err)
      }

      if err != sql.ErrNoRows {
        logger.Debug("report already processed", zap.Int64("generatedTime", row.GeneratedTime))
        continue
      }
    }

    insertStatement, err := insertManager.PrepareForInsert()
    if err != nil {
      return err
    }

    err = writeMetrics(row, insertStatement, logger)
    if err != nil {
      return err
    }
  }

  err = sourceRows.Err()
  if err != nil {
    return errors.WithStack(err)
  }
  return insertManager.Commit()
}

func writeMetrics(row *MetricResult, insertStatement *sql.Stmt, logger *zap.Logger) error {
  report, err := analyzer.ReadReport([]byte(row.RawReport))
  if err != nil {
    return errors.WithStack(err)
  }

  durationMetrics, instantMetrics := analyzer.ComputeMetrics(report, logger)
  // or both null, or not - no need to check each one
  if durationMetrics == nil || instantMetrics == nil {
    return errors.New("metrics cannot be computed")
  }

  buildTimeUnix, err := analyzer.GetBuildTimeFromReport(report)
  if err != nil {
    return err
  }

  if buildTimeUnix <= 0 {
    buildTimeUnix = row.BuildTime
  }

  var branch string
  if len(row.Branch) == 0 {
    branch = "master"
  } else {
    branch = string(row.Branch)
  }

  args := []interface{}{row.Product, row.Machine, buildTimeUnix, row.GeneratedTime,
    row.TcBuildId, row.TcInstallerBuildId, row.TcBuildProperties,
    branch,
    row.RawReport, row.BuildC1, row.BuildC2, row.BuildC3}
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
      return errors.Errorf("value outside of uint16 range (generatedTime: %d, value: %d)", row.GeneratedTime, v)
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
  return errors.WithStack(err)
}
