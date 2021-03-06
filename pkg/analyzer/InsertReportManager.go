package analyzer

import (
  "context"
  "database/sql"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/model"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/sql-util"
  "github.com/develar/errors"
  "github.com/jmoiron/sqlx"
  "go.deanishe.net/env"
  "go.uber.org/zap"
  "strings"
  "time"
)

// use `select distinct cast(machine, 'Uint16') as id, machine as name FROM report order by id` to get current enum values
// for now machine enum should be updated manually if a new machine will be added

var ErrMetricsCannotBeComputed = errors.New("metrics cannot be computed")

type RunResult struct {
  Product string

  // or uint8 as enum id, or string as enum value
  Machine interface{}

  BuildTime     int64
  GeneratedTime int64

  TcBuildId          int
  TcInstallerBuildId int
  TcBuildProperties  []byte

  RawReport []byte

  // maybe null
  Report *model.Report

  BuildC1 int `db:"build_c1"`
  BuildC2 int `db:"build_c2"`
  BuildC3 int `db:"build_c3"`

  extraFieldData  []interface{}

  branch string
}

type InsertReportManager struct {
  sql_util.InsertDataManager

  context                          context.Context
  MaxGeneratedTime                 int64
  IsCheckThatNotAlreadyAddedNeeded bool
  config                           DatabaseConfiguration
  dbName                           string
  nonMetricFieldCount              int
  insertInstallerManager           *InsertInstallerManager
  tableName                        string
}

func NewInsertReportManager(db *sqlx.DB, dbName string, context context.Context, tableName string, insertWorkerCount int, logger *zap.Logger) (*InsertReportManager, error) {
  var selectStatement *sql.Stmt
  var err error
  config := GetAnalyzer(dbName)
  if config.HasProductField {
    selectStatement, err = db.Prepare("select 1 from " + tableName + " where product = ? and machine = ? and project = ? and generated_time = ? limit 1")
  } else {
    selectStatement, err = db.Prepare("select 1 from " + tableName + " where machine = ? and project = ? and generated_time = ? limit 1")
  }
  if err != nil {
    return nil, errors.WithStack(err)
  }

  // product, row.Machine, buildTimeUnix, row.GeneratedTime, project,
  //    row.TcBuildId, row.TcInstallerBuildId, row.TcBuildProperties,
  //    branch,
  //    row.RawReport, row.BuildC1, row.BuildC2, row.BuildC3
  var sb strings.Builder
  sb.WriteString("insert into ")
  sb.WriteString(tableName)
  sb.WriteString(" (")
  if config.HasProductField {
    sb.WriteString("product, ")
  }
  sb.WriteString("machine, build_time, generated_time, project, tc_build_id, tc_installer_build_id, tc_build_properties, branch, raw_report, build_c1, build_c2, build_c3")

  config.insertStatementWriter(&sb)
  sb.WriteString(") values (")

  nonMetricFieldCount := 13
  if !config.HasProductField {
    nonMetricFieldCount -= 1
  }

  for i, n := 0, nonMetricFieldCount+config.extraFieldCount; i < n; i++ {
    if i != 0 {
      sb.WriteRune(',')
    }
    sb.WriteRune('?')
  }
  sb.WriteRune(')')

  effectiveSql := sb.String()
  insertManager, err := sql_util.NewBulkInsertManager(db, context, effectiveSql, insertWorkerCount, logger.Named("report"))
  if err != nil {
    return nil, errors.WithStack(err)
  }

  // large inserts leads to large memory usage, so, allow to override INSERT_BATCH_SIZE via env
  insertManager.BatchSize = env.GetInt("INSERT_BATCH_SIZE", 1_000)

  installerManager, err := NewInstallerInsertManager(db, context, logger)
  if err != nil {
    return nil, err
  }

  manager := &InsertReportManager{
    nonMetricFieldCount: nonMetricFieldCount,
    config:              config,
    dbName:              dbName,
    tableName:           tableName,
    InsertDataManager: sql_util.InsertDataManager{
      Db: db,

      SelectStatement: selectStatement,
      InsertManager:   insertManager,

      Logger: logger,
    },

    context:                context,
    insertInstallerManager: installerManager,
  }

  insertManager.AddDependency(installerManager.InsertManager)
  return manager, nil
}

// checks that not duplicated, warn if metrics cannot be computed
func (t *InsertReportManager) Insert(runResult *RunResult) error {
  logger := t.Logger.With(zap.String("product", runResult.Product), zap.String("generatedTime", time.Unix(runResult.GeneratedTime, 0).Format(time.RFC1123)))

  // tc collector uses tc build id to avoid duplicates, so, IsCheckThatNotAlreadyAddedNeeded is set to false by default
  if t.IsCheckThatNotAlreadyAddedNeeded && runResult.GeneratedTime <= t.MaxGeneratedTime {
    exists, err := t.CheckExists(t.SelectStatement.QueryRow(runResult.Product, runResult.Machine, runResult.Report.Project, runResult.GeneratedTime))
    if err != nil {
      return err
    }

    if exists {
      logger.Debug("report already processed")
      return nil
    }
  }

  err := t.WriteMetrics(runResult.Product, runResult, runResult.branch, runResult.Report.Project, logger)
  if err != nil {
    if err == ErrMetricsCannotBeComputed {
      logger.Warn(err.Error())
      return nil
    }
    return err
  }

  logger.Debug("new report added")
  return nil
}

//goland:noinspection SpellCheckingInspection
var projectIdToName = map[string]string{
  "Fleet": "fleet",

  // IJ simple project - project v3
  "Mc92Qmj3NY0xxdIiX9ayVbbEZ7s": "simple for IJ",
  // IJ simple project - project v2
  "73YWaW9bytiPDGuKvwNIYMK5CKI": "simple for IJ",

  // idea project (v2)
  "26hfTKDRtXpJ6U7ivgfKthtyU0A": "idea",
  // idea project (v3)
  "nC4MRRFMVYUSQLNIvPgDt+B3JqA": "idea",

  // light edit
  "6hglkyp/cmAi7ntjrg7dHwd5NG4": "light edit (IJ)",
  "1PbxeQ044EEghMOG9hNEFee05kM": "light edit (IJ)",
}

func (t *InsertReportManager) WriteMetrics(product string, row *RunResult, branch string, providedProject string, logger *zap.Logger) error {
  insertStatement, err := t.InsertManager.PrepareForInsert()
  if err != nil {
    return err
  }

  var project string
  if len(row.Report.Project) == 0 {
    project = providedProject
  } else {
    project = projectIdToName[row.Report.Project]
    if len(project) == 0 {
      project = row.Report.Project
    }
  }

  buildTimeUnix, err := getBuildTimeFromReport(row.Report, t.dbName)
  if err != nil {
    return err
  }

  if buildTimeUnix <= 0 {
    buildTimeUnix = row.BuildTime
  }

  if m, ok := row.Machine.(string); ok {
    if strings.HasPrefix(m, "intellij-linux-hw-compile-hp-blade-") {
      return nil
    }
  }

  args := make([]interface{}, 0, t.nonMetricFieldCount+t.config.extraFieldCount)
  if t.config.HasProductField {
    args = append(args, product)
  }
  args = append(args, row.Machine, buildTimeUnix, row.GeneratedTime, project,
    row.TcBuildId, row.TcInstallerBuildId, row.TcBuildProperties,
    branch,
    row.RawReport, row.BuildC1, row.BuildC2, row.BuildC3)

  if t.dbName == "ij" {
    err = ComputeIjMetrics(t.nonMetricFieldCount, row.Report, &args, logger)
    if err != nil {
      return err
    }
  }

  args = append(args, row.extraFieldData...)

  _, err = insertStatement.ExecContext(t.context, args...)
  if err != nil {
    return errors.WithStack(err)
  }
  return nil
}
