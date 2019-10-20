package analyzer

import (
  "github.com/develar/errors"
  "github.com/json-iterator/go"
  "github.com/mcuadros/go-version"
  "report-aggregator/pkg/model"
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