package analyzer

import (
  "context"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/model"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/sql-util"
  "github.com/deanishe/go-env"
  "github.com/develar/errors"
  "github.com/jmoiron/sqlx"
  "go.uber.org/zap"
  "strings"
  "time"
)

// use `select distinct cast(machine, 'Uint16') as id, machine as name FROM report order by id` to get current enum values
// for now machine enum should be updated manually if a new machine will be added

var ErrMetricsCannotBeComputed = errors.New("metrics cannot be computed")

type MetricResult struct {
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
}

type InsertReportManager struct {
  sql_util.InsertDataManager

  context                  context.Context
  MaxGeneratedTime                 int64
  IsCheckThatNotAlreadyAddedNeeded bool

  insertInstallerManager *InsertInstallerManager
  tableName              string
}

func NewInsertReportManager(db *sqlx.DB, context context.Context, tableName string, insertWorkerCount int, logger *zap.Logger) (*InsertReportManager, error) {
  selectStatement, err := db.Prepare("select 1 from " + tableName + " where product = ? and machine = ? and project = ? and generated_time = ? limit 1")
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
  sb.WriteString(" (product, machine, build_time, generated_time, project, tc_build_id, tc_installer_build_id, tc_build_properties, branch, raw_report, build_c1, build_c2, build_c3")
  for _, metric := range MetricDescriptors {
    sb.WriteRune(',')
    sb.WriteString(metric.Name)
  }
  sb.WriteString(") values (")

  for i, n := 0, nonMetricFieldCount+len(MetricDescriptors); i < n; i++ {
    if i != 0 {
      sb.WriteRune(',')
    }
    sb.WriteRune('?')
  }
  sb.WriteRune(')')

  sql := sb.String()
  insertManager, err := sql_util.NewBulkInsertManager(db, context, sql, insertWorkerCount, logger.Named("report"))
  if err != nil {
    return nil, errors.WithStack(err)
  }

  // large inserts leads to large memory usage, so, allow to override INSERT_BATCH_SIZE via env
  insertManager.BatchSize = env.GetInt("INSERT_BATCH_SIZE", 1000)

  installerManager, err := NewInstallerInsertManager(db, context, logger)
  if err != nil {
    return nil, err
  }

  manager := &InsertReportManager{
    tableName: tableName,
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
func (t *InsertReportManager) Insert(row *MetricResult, branch string, providedProject string) error {
  logger := t.Logger.With(zap.String("product", row.Product), zap.String("generatedTime", time.Unix(row.GeneratedTime, 0).Format(time.RFC1123)))

  // tc collector uses tc build id to avoid duplicates, so, IsCheckThatNotAlreadyAddedNeeded is set to false by default
  if t.IsCheckThatNotAlreadyAddedNeeded && row.GeneratedTime <= t.MaxGeneratedTime {
    exists, err := t.CheckExists(t.SelectStatement.QueryRow(row.Product, row.Machine, providedProject, row.GeneratedTime))
    if err != nil {
      return err
    }

    if exists {
      logger.Debug("report already processed")
      return nil
    }
  }

  err := t.WriteMetrics(row.Product, row, branch, providedProject, logger)
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

func (t *InsertReportManager) WriteMetrics(product interface{}, row *MetricResult, branch interface{}, providedProject string, logger *zap.Logger) error {
  insertStatement, err := t.InsertManager.PrepareForInsert()
  if err != nil {
    return err
  }

  report := row.Report
  if report == nil {
    report, err = ReadReport(row.RawReport)
    if err != nil {
      return err
    }
  }

  project := report.Project
  if len(project) == 0 {
    project = providedProject
    if len(project) == 0 {
      return errors.New("unknown project")
    }
  }

  buildTimeUnix, err := GetBuildTimeFromReport(report)
  if err != nil {
    return err
  }

  if buildTimeUnix <= 0 {
    buildTimeUnix = row.BuildTime
  }

  if strings.HasPrefix(row.Machine.(string), "intellij-linux-hw-blade-") {
    return nil
  }

  args := make([]interface{}, 0, nonMetricFieldCount+len(MetricDescriptors))
  args = append(args, product, row.Machine, buildTimeUnix, row.GeneratedTime, project,
    row.TcBuildId, row.TcInstallerBuildId, row.TcBuildProperties,
    branch,
    row.RawReport, row.BuildC1, row.BuildC2, row.BuildC3)

  err = ComputeMetrics(report, &args, logger)
  if err != nil {
    return err
  }

  _, err = insertStatement.ExecContext(t.context, args...)
  return errors.WithStack(err)
}
