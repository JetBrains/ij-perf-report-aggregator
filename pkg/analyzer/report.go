package analyzer

import (
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/model"
  "github.com/develar/errors"
  "github.com/mcuadros/go-version"
  "github.com/valyala/fastjson"
  "time"
)

var structParsers fastjson.ParserPool

func ReadReport(runResult *RunResult) error {
  parser := structParsers.Get()
  defer structParsers.Put(parser)

  report, err := parser.ParseBytes(runResult.RawReport)
  if err != nil {
    //fmt.Print(string(data))
    return  errors.WithStack(err)
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

    TotalDurationActual: report.GetInt("totalDurationActual"),
  }

  measureNames := make([]string, 0)
  measureValues := make([]int, 0)
  for _, measure := range report.GetArray("metrics") {
    measureName := string(measure.GetStringBytes("n"))

    value := measure.Get("d")
    if value == nil {
      value = measure.Get("c")
      if value == nil {
        return nil
      }
    }

    floatValue := value.GetFloat64()
    intValue := int(floatValue)
    if floatValue != float64(intValue) {
      return errors.WithMessagef(nil, "int expected, but got float %f", floatValue)
    }

    measureNames = append(measureNames, measureName)
    measureValues = append(measureValues, intValue)
  }

  if len(measureNames) == 0 {
    return nil
  }

  runResult.measureNames = measureNames
  runResult.measureValues = measureValues
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