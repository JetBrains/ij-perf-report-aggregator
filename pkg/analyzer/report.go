package analyzer

import (
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/model"
  "github.com/develar/errors"
  "github.com/mcuadros/go-version"
  "github.com/valyala/fastjson"
  "go.uber.org/zap"
  "time"
)

var structParsers fastjson.ParserPool

func ReadReport(runResult *RunResult, reader CustomReportAnalyzer, logger *zap.Logger) error {
  parser := structParsers.Get()
  defer structParsers.Put(parser)

  report, err := parser.ParseBytes(runResult.RawReport)
  if err != nil {
    return errors.WithStack(err)
  }

  totalDuration := report.GetInt("totalDuration")
  if totalDuration == 0 {
    totalDuration = report.GetInt("totalDurationActual")
  }

  runResult.Report = &model.Report{
    Version:                  string(report.GetStringBytes("version")),
    Generated:                string(report.GetStringBytes("generated")),
    Project:                  string(report.GetStringBytes("project")),

    Build:                    string(report.GetStringBytes("build")),
    BuildDate:                string(report.GetStringBytes("buildDate")),

    Os:                       string(report.GetStringBytes("os")),
    ProductCode:              string(report.GetStringBytes("productCode")),
    Runtime:                  string(report.GetStringBytes("runtime")),

    TotalDuration: totalDuration,
  }

  err = reader(runResult, report, logger)
  if err != nil {
    return err
  }
  return nil
}

func getBuildTimeFromReport(report *model.Report, dbName string) (int64, error) {
  var buildTimeUnix int64
  if dbName != "ij" || version.Compare(report.Version, "13", ">=") {
    buildTime, err := parseTime(report.BuildDate)
    if err != nil {
      return 0, err
    }
    buildTimeUnix = buildTime.Unix()
  } else {
    buildTimeUnix = 0
  }
  return buildTimeUnix, nil
}

func parseTime(s string) (*time.Time, error) {
  parsedTime, err := time.Parse(time.RFC1123Z, s)
  if err != nil {
    parsedTime, err = time.Parse(time.RFC1123, s)
  }

  if err != nil {
    parsedTime, err = time.Parse("Jan 2, 2006, 3:04:05 PM MST", s)
  }

  if err != nil {
    parsedTime, err = time.Parse("Mon, 2 Jan 2006 15:04:05 -0700", s)
  }

  if err != nil {
    parsedTime, err = time.Parse("Mon, 2 Jan 2006 15:04:05 MST", s)
  }

  if err != nil {
    return nil, errors.WithStack(err)
  }
  return &parsedTime, nil
}