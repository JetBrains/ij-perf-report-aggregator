package analyzer

import (
  "github.com/alecthomas/kingpin"
  "github.com/bvinc/go-sqlite-lite/sqlite3"
  "github.com/develar/errors"
  "go.uber.org/zap"
  "report-aggregator/pkg/util"
)

func ConfigureUpdateMetricsCommand(app *kingpin.Application, logger *zap.Logger) {
  command := app.Command("update-computed-metrics", "Update computed metrics.")
  dbPath := command.Flag("db", "The SQLite database file.").Required().String()
  command.Action(func(context *kingpin.ParseContext) error {
    return UpdateMetrics(*dbPath, logger)
  })
}

func UpdateMetrics(dbPath string, logger *zap.Logger) error {
  db, err := sqlite3.Open(dbPath, sqlite3.OPEN_READWRITE)
  if err != nil {
    return errors.WithStack(err)
  }

  defer util.Close(db, logger)

  selectStatement, err := db.Prepare(`
	   select id, raw_report
	   from report
	   where metrics_version < ?
  `)
  if err != nil {
    return errors.WithStack(err)
  }

  defer util.Close(selectStatement, logger)

  updateStatement, err := db.Prepare(`
	   update report
	   set metrics_version = ?, duration_metrics = ?, instant_metrics = ?
	   where id = ?
  `)
  if err != nil {
    return errors.WithStack(err)
  }

  defer util.Close(updateStatement, logger)

  err = selectStatement.Bind(MetricsVersion)
  if err != nil {
    return errors.WithStack(err)
  }

  updatedCount := 0

  for {
    hasRow, err := selectStatement.Step()
    if err != nil {
      return errors.WithStack(err)
    }

    if !hasRow {
      break
    }

    id, _, err := selectStatement.ColumnRawString(0)
    if err != nil {
      return errors.WithStack(err)
    }

    rawJson, err := selectStatement.ColumnRawBytes(1)
    if err != nil {
      return errors.WithStack(err)
    }

    report, err := readReport(rawJson)
    if err != nil {
      return err
    }

    serializedDurationMetrics, serializedInstantMetrics, err := computeAndSerializeMetrics(report, logger)
    if err != nil {
      return err
    }

    if len(serializedDurationMetrics) == 0 || len(serializedInstantMetrics) == 0 {
      // it is not warn for metric updater because on update metrics must be computed
      return errors.New("metrics cannot be computed")
    }

    err = updateStatement.Exec(MetricsVersion, serializedDurationMetrics, serializedInstantMetrics, id)
    if err != nil {
      return errors.WithStack(err)
    }

    updatedCount++
  }

  if updatedCount == 0 {
    logger.Info("all metrics up to date, nothing to update", zap.Int("currentVersion", MetricsVersion))
  } else {
    logger.Info("metrics updated", zap.Int("count", updatedCount), zap.Int("currentVersion", MetricsVersion))
  }
  return nil
}
