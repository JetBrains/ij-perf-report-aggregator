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
  "strings"
)

func ConfigureFillCommand(app *kingpin.Application, logger *zap.Logger) {
  command := app.Command("fill", "Fill VictoriaMetrics database using SQLite database.")
  dbPath := command.Flag("db", "The SQLite database file.").Required().String()
  promServer := command.Flag("prom", "The VictoriaMetrics/Influx server.").Required().String()
  command.Action(func(context *kingpin.ParseContext) error {
    err := fill(*dbPath, *promServer, logger)
    if err != nil {
      return err
    }

    return nil
  })
}

func fill(dbPath string, promServer string, logger *zap.Logger) error {
  db, err := sqlite3.Open(dbPath, sqlite3.OPEN_READONLY)
  if err != nil {
    return errors.WithStack(err)
  }

  defer util.Close(db, logger)

  statement, err := db.Prepare(`
 select id, machine, generated_time, metrics, 
        product,
        build_c1, build_c2, build_c3
 from report order by generated_time
 	`)

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
  for {
    hasRow, err := statement.Step()
    if !hasRow {
      return nil
    }

    if err != nil {
      logger.Error("cannot step", zap.Error(err))
      return err
    }

    err = writeMetrics(statement, pw, logger)
    if err != nil {
      return err
    }
  }
}

func writeMetrics(statement *sqlite3.Stmt, pw *io.PipeWriter, logger *zap.Logger) error {
  var id sqlite3.RawString
  var machine sqlite3.RawString
  var generatedTime int64
  var metricsJson sqlite3.RawString
  var productCode sqlite3.RawString
  var buildC1 int
  var buildC2 int
  var buildC3 int
  err := statement.Scan(&id, &machine, &generatedTime, &metricsJson, &productCode, &buildC1, &buildC2, &buildC3)
  if err != nil {
    logger.Error("cannot scan", zap.Error(err))
    return err
  }

  var metrics map[string]int
  err = jsoniter.ConfigFastest.Unmarshal([]byte(metricsJson), &metrics)
  if err != nil {
    return errors.WithStack(err)
  }

  var buf bytes.Buffer
  for key, value := range metrics {
    isInstant := strings.HasPrefix(key, "i_")

    if isInstant {
      buf.WriteString(key[2:])
    } else {
      buf.WriteString(key)
    }

    buf.WriteByte(',')

    buf.WriteString("id=")
    buf.WriteString(string(id))
    buf.WriteByte(',')

    buf.WriteString("machine=")
    buf.WriteString(string(machine))
    buf.WriteByte(',')

    buf.WriteString("product=")
    buf.WriteString(string(productCode))
    buf.WriteByte(',')

    buf.WriteString("buildC1=")
    buf.WriteString(strconv.Itoa(buildC1))
    buf.WriteByte(',')

    buf.WriteString("buildC2=")
    buf.WriteString(strconv.Itoa(buildC2))
    buf.WriteByte(',')

    buf.WriteString("buildC3=")
    buf.WriteString(strconv.Itoa(buildC3))

    buf.WriteByte(' ')
    if isInstant {
      buf.WriteByte('i')
    } else {
      buf.WriteByte('d')
    }
    buf.WriteByte('=')
    buf.WriteString(strconv.Itoa(value))
    buf.WriteByte(' ')

    buf.WriteString(strconv.FormatInt(generatedTime, 10))

    logger.Info(buf.String())

    buf.WriteByte('\n')

    _, err := buf.WriteTo(pw)
    if err != nil {
      logger.Error("cannot write", zap.Error(err))
      return err
    }

    buf.Reset()
  }

  return nil
}
