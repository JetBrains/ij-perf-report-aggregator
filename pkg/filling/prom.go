package filling

import (
  "bytes"
  "github.com/alecthomas/kingpin"
  "github.com/bvinc/go-sqlite-lite/sqlite3"
  "github.com/develar/errors"
  "github.com/json-iterator/go"
  "go.uber.org/zap"
  "io"
  "net/http"
  "report-aggregator/pkg/util"
  "strconv"
)

func ConfigureFillCommand(app *kingpin.Application, logger *zap.Logger) {
	command := app.Command("fill", "Fill VictoriaMetrics database using SQLite database.")
	dbPath := command.Flag("db", "The SQLite database file.").Required().String()
	promServer := command.Flag("prom", "The VictoriaMetrics/Influx server.").Required().String()
	command.Action(func(context *kingpin.ParseContext) error {
		return fill(*dbPath, *promServer, logger)
	})
}

func fill(dbPath string, promServer string, logger *zap.Logger) error {
	db, err := sqlite3.Open(dbPath, sqlite3.OPEN_READONLY)
	if err != nil {
		return errors.WithStack(err)
	}

	defer util.Close(db, logger)

	statement, err := createMetricsQuery(db)
	if err != nil {
		return errors.WithStack(err)
	}

	defer util.Close(statement, logger)

	pr, pw := io.Pipe()
	go func() {
		_ = pw.CloseWithError(writeReports(statement, pw, logger))
	}()

	// it is not data (e.g. duration or start) precision, but timestamp (report generated time)
	response, err := http.Post(promServer+"/write?precision=s", "text/plain", pr)
	if err != nil {
		return errors.WithStack(err)
	}

	if response != nil {
		defer util.Close(response.Body, logger)
	}

	if response.StatusCode >= 300 {
		return errors.New("Failed to send: " + response.Status)
	}
	return nil
}

func writeReports(statement *sqlite3.Stmt, pw *io.PipeWriter, logger *zap.Logger) error {
  row := &MetricResult{}
	for {
		hasRow, err := statement.Step()
		if !hasRow {
			return nil
		}

		if err != nil {
			logger.Error("cannot step", zap.Error(err))
			return err
		}

		err = writeMetrics(statement, row, pw, logger)
		if err != nil {
			return err
		}
	}
}

func writeMetrics(statement *sqlite3.Stmt, row *MetricResult, pw *io.PipeWriter, logger *zap.Logger) error {
	err := scanMetricResult(statement, row)
	if err != nil {
		logger.Error("cannot scan", zap.Error(err))
		return err
	}

	var instantMetricsJson map[string]int
	err = jsoniter.ConfigFastest.Unmarshal([]byte(row.instantMetricsJson), &instantMetricsJson)
	if err != nil {
		return errors.WithStack(err)
	}

	var buf bytes.Buffer
  err = doWriteMetrics(row.durationMetricsJson, false, &buf, row, pw, logger)
  if err != nil {
    return err
  }

  err = doWriteMetrics(row.instantMetricsJson, true, &buf, row, pw, logger)
  if err != nil {
    return err
  }

  return nil
}

func doWriteMetrics(metricsJson sqlite3.RawString, isInstant bool, buf *bytes.Buffer, row *MetricResult, pw *io.PipeWriter, logger *zap.Logger) error {
  var metrics map[string]int
 	err := jsoniter.ConfigFastest.Unmarshal([]byte(metricsJson), &metrics)
 	if err != nil {
 		return errors.WithStack(err)
 	}

  for key, value := range metrics {
    // skip missing values (-1 means null (undefined))
    if value == -1 {
      continue
    }

    buf.Reset()

    buf.WriteString(key)

    buf.WriteByte(',')

    // write tags

    buf.WriteString("id=")
    buf.WriteString(string(row.id))
    buf.WriteByte(',')

    buf.WriteString("machine=")
    buf.WriteString(string(row.machine))
    buf.WriteByte(',')

    buf.WriteString("product=")
    buf.WriteString(string(row.productCode))
    buf.WriteByte(',')

    buf.WriteString("buildC1=")
    buf.WriteString(strconv.Itoa(row.buildC1))
    buf.WriteByte(',')

    buf.WriteString("buildC2=")
    buf.WriteString(strconv.Itoa(row.buildC2))
    buf.WriteByte(',')

    buf.WriteString("buildC3=")
    buf.WriteString(strconv.Itoa(row.buildC3))

    // write fields
    buf.WriteByte(' ')
    if isInstant {
      buf.WriteByte('i')
    } else {
      buf.WriteByte('d')
    }
    buf.WriteByte('=')
    buf.WriteString(strconv.Itoa(value))
    // https://docs.influxdata.com/influxdb/v1.7/write_protocols/line_protocol_tutorial/#data-types
    // integer type
    buf.WriteByte('i')
    buf.WriteByte(' ')

    buf.WriteString(strconv.FormatInt(row.generatedTime, 10))

    logger.Info(buf.String())

    buf.WriteByte('\n')

    _, err := buf.WriteTo(pw)
    if err != nil {
      logger.Error("cannot write", zap.Error(err))
      return err
    }
  }
  return nil
}
