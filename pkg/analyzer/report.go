package analyzer

import (
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/model"
  "github.com/develar/errors"
  "github.com/json-iterator/go"
  "github.com/mcuadros/go-version"
  "time"
)

func ReadReport(data []byte) (*model.Report, error) {
  var report model.Report
  err := jsoniter.ConfigFastest.Unmarshal(data, &report)
  if err != nil {
    return nil, errors.WithStack(err)
  }
  return &report, nil
}

func GetBuildTimeFromReport(report *model.Report) (int64, error) {
  var buildTimeUnix int64
  if version.Compare(report.Version, "13", ">=") {
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
    return nil, errors.WithStack(err)
  }
  return &parsedTime, nil
}