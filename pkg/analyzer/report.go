package analyzer

import (
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/model"
  "github.com/develar/errors"
  "github.com/mcuadros/go-version"
  "github.com/valyala/fastjson"
  "time"
)

var structParsers fastjson.ParserPool

func ReadReport(data []byte) (*model.Report, error) {
  parser := structParsers.Get()
  defer structParsers.Put(parser)

  value, err := parser.ParseBytes(data)
  if err != nil {
    //fmt.Print(string(data))
    return nil, errors.WithStack(err)
  }

  report := model.Report{
    Version:                  string(value.GetStringBytes("version")),
    Generated:                string(value.GetStringBytes("generated")),
    Project:                  string(value.GetStringBytes("project")),

    Build:                    string(value.GetStringBytes("build")),
    BuildDate:                string(value.GetStringBytes("buildDate")),

    Os:                       string(value.GetStringBytes("os")),
    ProductCode:              string(value.GetStringBytes("productCode")),
    Runtime:                  string(value.GetStringBytes("runtime")),

    TotalDurationActual:      value.GetInt("totalDurationActual"),
  }

  for _, v := range value.GetArray("traceEvents") {
    report.TraceEvents = append(report.TraceEvents, model.TraceEvent{
      Name:      string(v.GetStringBytes("name")),
      Phase:     string(v.GetStringBytes("ph")),
      Timestamp: v.GetInt("ts"),
      Category:  string(v.GetStringBytes("cat")),
    })
  }
  report.MainActivities = readActivitiesInOldFormat("items", value)
  if version.Compare(report.Version, "20", ">=") {
    report.PrepareAppInitActivities = readActivities("prepareAppInitActivities", value)
    if version.Compare(report.Version, "27", ">=") {
      classLoading := value.Get("classLoading")
      resourceLoading := value.Get("resourceLoading")
      if classLoading != nil && resourceLoading != nil {
        report.PrepareAppInitActivities = append(report.PrepareAppInitActivities,
          model.Activity{Name: clTotal, Start: classLoading.GetInt("time")},
          model.Activity{Name: clSearch, Start: classLoading.GetInt("searchTime")},
          model.Activity{Name: clDefine, Start: classLoading.GetInt("defineTime")},
          model.Activity{Name: clCount, Start: classLoading.GetInt("count")},

          model.Activity{Name: rlTime, Start: resourceLoading.GetInt("time")},
          model.Activity{Name: rlCount, Start: resourceLoading.GetInt("count")},
        )
      }
    }
  } else {
    report.PrepareAppInitActivities = readActivitiesInOldFormat("prepareAppInitActivities", value)
  }
  return &report, nil
}

func readActivitiesInOldFormat(key string, value *fastjson.Value) []model.Activity {
  array := value.GetArray(key)
  result := make([]model.Activity, 0, len(array))
  for _, v := range array {
    result = append(result, model.Activity{
      Name:      string(v.GetStringBytes("name")),
      Thread:     string(v.GetStringBytes("thread")),
      Start: v.GetInt("start"),
      End: v.GetInt("end"),
      Duration: v.GetInt("duration"),
    })
  }
  return result
}

func readActivities(key string, value *fastjson.Value) []model.Activity {
  array := value.GetArray(key)
  result := make([]model.Activity, 0, len(array))
  for _, v := range array {
    start := v.GetInt("s")
    duration := v.GetInt("d")

    ownDuration := v.GetInt("od")
    if ownDuration == 0 {
      ownDuration = duration
    }

    result = append(result, model.Activity{
      Name:     string(v.GetStringBytes("n")),
      Thread:   string(v.GetStringBytes("t")),
      Start:    start,
      End:      start + duration,
      Duration: ownDuration,
    })
  }
  return result
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