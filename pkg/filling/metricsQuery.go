package filling

import "github.com/bvinc/go-sqlite-lite/sqlite3"

type MetricResult struct {
	id sqlite3.RawString

	productCode sqlite3.RawString
	machine     sqlite3.RawString

	generatedTime int64

	durationMetricsJson sqlite3.RawString
	instantMetricsJson  sqlite3.RawString

	buildC1 int
	buildC2 int
	buildC3 int
}

func scanMetricResult(statement *sqlite3.Stmt, row *MetricResult) error {
	var err error
	i := 0

	row.id, _, err = statement.ColumnRawString(i)
	i++
	if err != nil {
		return err
	}

	row.productCode, _, err = statement.ColumnRawString(i)
	i++
	if err != nil {
		return err
	}
	row.machine, _, err = statement.ColumnRawString(i)
	i++
	if err != nil {
		return err
	}

	row.generatedTime, _, err = statement.ColumnInt64(i)
	i++
	if err != nil {
		return err
	}

	row.durationMetricsJson, _, err = statement.ColumnRawString(i)
	i++
	if err != nil {
		return err
	}
	row.instantMetricsJson, _, err = statement.ColumnRawString(i)
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
