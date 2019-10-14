package filling

import "github.com/bvinc/go-sqlite-lite/sqlite3"

type MetricResult struct {
  productCode string
  machine     string

  generatedTime int64
  tcBuildId     int

  durationMetricsJson string
  instantMetricsJson  string

  rawReport string

  buildC1 int
  buildC2 int
  buildC3 int
}

func scanMetricResult(statement *sqlite3.Stmt, row *MetricResult) error {
  var err error
  i := 0

  row.productCode, _, err = statement.ColumnText(i)
  i++
  if err != nil {
    return err
  }
  row.machine, _, err = statement.ColumnText(i)
  i++
  if err != nil {
    return err
  }

  row.generatedTime, _, err = statement.ColumnInt64(i)
  i++
  if err != nil {
    return err
  }

  row.tcBuildId, _, err = statement.ColumnInt(i)
  i++
  if err != nil {
    return err
  }

  row.durationMetricsJson, _, err = statement.ColumnText(i)
  i++
  if err != nil {
    return err
  }
  row.instantMetricsJson, _, err = statement.ColumnText(i)
  i++
  if err != nil {
    return err
  }

  row.rawReport, _, err = statement.ColumnText(i)
  i++
  if err != nil {
    return err
  }

  row.buildC1, _, err = statement.ColumnInt(i)
  i++
  if err != nil {
    return err
  }
  row.buildC2, _, err = statement.ColumnInt(i)
  i++
  if err != nil {
    return err
  }
  row.buildC3, _, err = statement.ColumnInt(i)
  if err != nil {
    return err
  }

  return nil
}
