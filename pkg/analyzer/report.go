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

func ReadReport(runResult *RunResult, config DatabaseConfiguration, logger *zap.Logger) error {
  parser := structParsers.Get()
  defer structParsers.Put(parser)

  report, err := parser.ParseBytes(runResult.RawReport)
  if err != nil {
    logger.Warn("invalid report - corrupted JSON, report will be skipped", zap.Int("id", runResult.TcBuildId), zap.ByteString("rawReport", runResult.RawReport))
    runResult.Report = nil
    return nil
  }

  runResult.Report = &model.Report{
    Version:   string(report.GetStringBytes("version")),
    Generated: string(report.GetStringBytes("generated")),
    Project:   string(report.GetStringBytes("project")),

    Os:          string(report.GetStringBytes("os")),
    ProductCode: string(report.GetStringBytes("productCode")),
    Runtime:     string(report.GetStringBytes("runtime")),
  }

  if config.HasInstallerField {
    runResult.Report.Build = string(report.GetStringBytes("build"))
    runResult.Report.BuildDate = string(report.GetStringBytes("buildDate"))
  }

  err = config.ReportReader(runResult, report, logger)
  if err != nil {
    return nil
  }

  runResult.RawReport = report.MarshalTo(nil)
  return nil
}

func getBuildTimeFromReport(report *model.Report, dbName string) (time.Time, error) {
  var buildTimeUnix time.Time
  if dbName != "ij" || version.Compare(report.Version, "13", ">=") {
    buildTime, err := ParseTime(report.BuildDate)
    if err != nil {
      return time.Time{}, err
    }
    buildTimeUnix = buildTime
  } else {
    buildTimeUnix = time.Time{}
  }
  return buildTimeUnix, nil
}

func ParseTime(s string) (time.Time, error) {
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
    parsedTime, err = time.Parse("20060102T150405+0000", s)
  }

  if err != nil {
    return time.Time{}, errors.WithStack(err)
  }
  return parsedTime, nil
}
