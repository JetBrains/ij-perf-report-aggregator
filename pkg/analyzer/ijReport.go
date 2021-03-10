package analyzer

import (
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/model"
  "github.com/mcuadros/go-version"
  "github.com/valyala/fastjson"
)

func analyzeIJReport(runResult *RunResult, data *fastjson.Value) error {
  report := runResult.Report
  for _, v := range data.GetArray("traceEvents") {
    report.TraceEvents = append(report.TraceEvents, model.TraceEvent{
      Name:      string(v.GetStringBytes("name")),
      Phase:     string(v.GetStringBytes("ph")),
      Timestamp: v.GetInt("ts"),
      Category:  string(v.GetStringBytes("cat")),
    })
  }
  report.MainActivities = readActivitiesInOldFormat("items", data)
  if version.Compare(report.Version, "20", ">=") {
    report.PrepareAppInitActivities = readActivities("prepareAppInitActivities", data)
    if version.Compare(report.Version, "27", ">=") {
      classLoading := data.Get("classLoading")
      resourceLoading := data.Get("resourceLoading")
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
    report.PrepareAppInitActivities = readActivitiesInOldFormat("prepareAppInitActivities", data)
  }
  return nil
}

func readActivitiesInOldFormat(key string, data *fastjson.Value) []model.Activity {
  array := data.GetArray(key)
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