package analyzer

import (
  "context"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/model"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/sql-util"
  "github.com/develar/errors"
  "github.com/jmoiron/sqlx"
  "go.uber.org/zap"
  "strings"
  "time"
)

// use `select distinct cast(machine, 'Uint16') as id, machine as name FROM report order by id` to get current enum values
// for now machine enum should be updated manually if a new machine will be added

type MetricResult struct {
  Product string

  // or uint8 as enum id, or string as enum value
  Machine interface{}

  GeneratedTime int64
  BuildTime     int64

  TcBuildId          int
  TcInstallerBuildId int
  TcBuildProperties  []byte

  RawReport []byte

  BuildC1 int `db:"build_c1"`
  BuildC2 int `db:"build_c2"`
  BuildC3 int `db:"build_c3"`
}

type InsertReportManager struct {
  sql_util.InsertDataManager

  context          context.Context
  MaxGeneratedTime int64

  insertInstallerManager *InsertInstallerManager
}

func NewInsertReportManager(db *sqlx.DB, context context.Context, logger *zap.Logger) (*InsertReportManager, error) {
  selectStatement, err := db.Prepare("select 1 from report where product = ? and machine = ? and generated_time = ? limit 1")
  if err != nil {
    return nil, errors.WithStack(err)
  }

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

  insertManager, err := sql_util.NewBulkInsertManager(db, context, sb.String(), logger.Named("report"))
  if err != nil {
    return nil, errors.WithStack(err)
  }

  installerManager, err := NewInstallerInsertManager(db, context, logger)
  if err != nil {
    return nil, err
  }

  manager := &InsertReportManager{
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

  //noinspection SqlResolve
  var maxGeneratedTime time.Time
  err = db.QueryRow("select max(generated_time) from report").Scan(&maxGeneratedTime)
  if err != nil {
    return nil, errors.WithStack(err)
  }
  manager.MaxGeneratedTime = maxGeneratedTime.Unix()

  return manager, nil
}

func (t *InsertReportManager) Insert(row *MetricResult, branch string) error {
  logger := t.Logger.With(zap.String("product", row.Product), zap.String("generatedTime", time.Unix(row.GeneratedTime, 0).Format(time.RFC1123)))

  if row.GeneratedTime <= t.MaxGeneratedTime {
    exists, err := t.CheckExists(t.SelectStatement.QueryRow(row.Product, row.Machine, row.GeneratedTime))
    if err != nil {
      return err
    }

    if exists {
      logger.Debug("report already processed")
      return nil
    }
  }

  err := t.writeMetrics(row, branch, logger)
  if err != nil {
    return err
  }

  //logger.Debug("new report added")
  return nil
}

func (t *InsertReportManager) writeMetrics(row *MetricResult, branch string, logger *zap.Logger) error {
  insertStatement, err := t.InsertManager.PrepareForInsert()
  if err != nil {
    return err
  }

  report, err := ReadReport(row.RawReport)
  if err != nil {
    return errors.WithStack(err)
  }

  durationMetrics, instantMetrics := ComputeMetrics(report, logger)
  // or both null, or not - no need to check each one
  if durationMetrics == nil || instantMetrics == nil {
    return errors.New("metrics cannot be computed")
  }

  buildTimeUnix, err := GetBuildTimeFromReport(report)
  if err != nil {
    return err
  }

  if buildTimeUnix <= 0 {
    buildTimeUnix = row.BuildTime
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

  _, err = insertStatement.ExecContext(t.context, args...)
  return errors.WithStack(err)
}
