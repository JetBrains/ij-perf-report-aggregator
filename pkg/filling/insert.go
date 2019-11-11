package filling

import (
  "context"
  "database/sql"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/analyzer"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/model"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/util"
  "github.com/alecthomas/kingpin"
  "github.com/develar/errors"
  "github.com/iancoleman/strcase"
  "github.com/jmoiron/sqlx"
  "go.uber.org/zap"

  _ "github.com/ClickHouse/clickhouse-go"
  _ "github.com/mattn/go-sqlite3"
)

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
  //noinspection SqlResolve
  err = mainDb.Select(&products, "select row_number() over (order by product) AS id, product as name from report group by product")
  if err != nil {
    return errors.WithStack(err)
  }

  var machines []model.IdAndName
  //noinspection SqlResolve
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

  insertManager, err := analyzer.NewInsertReportManager(db, context.Background(), logger)
  if err != nil {
    return err
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
    order by product, machine, build_c1, build_c2, build_c3, generated_time`, insertManager.MaxGeneratedTime)
  if err != nil {
    return errors.WithStack(err)
  }

  defer util.Close(sourceRows, logger)

  logger.Info("start inserting", zap.String("clickHouseUrl", clickHouseUrl))

  defer util.Close(insertManager, logger)
  err = writeReports(insertManager, sourceRows)
  if err != nil {
    return err
  }

  logger.Info("waiting inserting", zap.Int("transactions", insertManager.GetUncommittedTransactionCount()))
  err = insertManager.CommitAndWait()
  if err != nil {
    return err
  }
  return nil
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

  insertManager, err := analyzer.NewInstallerInsertManager(db, context.Background(), logger)
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

    err = insertManager.Insert(id, changes)
    if err != nil {
      return err
    }
  }

  err = rows.Err()
  if err != nil {
    return errors.WithStack(err)
  }

  err = insertManager.CommitAndWait()
  if err != nil {
    return errors.WithStack(err)
  }
  return nil
}

type metricResultFromDb struct {
  analyzer.MetricResult

  Branch sql.NullString
}

func writeReports(insertManager *analyzer.InsertReportManager, sourceRows *sqlx.Rows) error {
  row := &metricResultFromDb{}
  for sourceRows.Next() {
    err := sourceRows.StructScan(row)
    if err != nil {
      return errors.WithStack(err)
    }


    row.Machine = uint8((row.Machine.(interface{})).(int64))

    var branch string
    if row.Branch.Valid {
      branch = row.Branch.String
    } else {
      branch = "master"
    }

    err = insertManager.Insert(&row.MetricResult, branch)
    if err != nil {
      return err
    }
  }

  err := sourceRows.Err()
  if err != nil {
    return errors.WithStack(err)
  }
  return nil
}