package analyzer

import (
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/model"
  "github.com/mcuadros/go-version"
  "github.com/valyala/fastjson"
  "go.uber.org/zap"
  "sort"
)

type measureItem struct {
  name     string
  start    uint32
  duration uint32
  thread   string
}

type serviceItem struct {
  name     string
  start    uint32
  duration uint32
  thread   string
  plugin   string
}

func analyzeIjReport(runResult *RunResult, data *fastjson.Value, logger *zap.Logger) error {
  report := runResult.Report

  report.TotalDuration = data.GetInt("totalDuration")
  if report.TotalDuration == 0 {
    report.TotalDuration = data.GetInt("totalDurationActual")
  }

  traceEvents := data.GetArray("traceEvents")

  if version.Compare(report.Version, "12", ">=") && len(traceEvents) == 0 {
    logger.Warn("invalid report (due to opening second project?), report will be skipped", zap.Int("id", runResult.TcBuildId), zap.String("generated", report.Generated))
    runResult.Report = nil
    return nil
  }

  for _, v := range traceEvents {
    report.TraceEvents = append(report.TraceEvents, model.TraceEvent{
      Name:      string(v.GetStringBytes("name")),
      Phase:     string(v.GetStringBytes("ph")),
      Timestamp: v.GetInt("ts"),
      Category:  string(v.GetStringBytes("cat")),
    })
  }

  var clTotal int32
  var clSearch int32
  var clDefine int32
  var clCount int32

  var rlTime int32
  var rlCount int32

  services := make([]serviceItem, 0)
  measures := make([]measureItem, 0)

  if version.Compare(report.Version, "20", ">=") {
    if version.Compare(report.Version, "32", ">=") {
      report.Activities = readActivities("items", data)

      for _, activity := range report.Activities {
        measures = append(measures, measureItem{
          name:     activity.Name,
          start:    uint32(activity.Start),
          duration: uint32(activity.Duration),
          thread:   activity.Thread,
        })
      }
    } else {
      report.Activities = readActivitiesInOldFormat("items", data)
      report.PrepareAppInitActivities = readActivities("prepareAppInitActivities", data)
    }

    if version.Compare(report.Version, "27", ">=") {
      classLoading := data.Get("classLoading")
      resourceLoading := data.Get("resourceLoading")
      if classLoading != nil && resourceLoading != nil {
        clTotal = int32(classLoading.GetInt("time"))
        clSearch = int32(classLoading.GetInt("searchTime"))
        clDefine = int32(classLoading.GetInt("defineTime"))
        clCount = int32(classLoading.GetInt("count"))

        rlTime = int32(resourceLoading.GetInt("time"))
        rlCount = int32(resourceLoading.GetInt("count"))
      }
    }

    readServices(data, "appComponents", &services)
    readServices(data, "appServices", &services)
    readServices(data, "projectComponents", &services)
    readServices(data, "projectServices", &services)
  } else {
    report.Activities = readActivitiesInOldFormat("items", data)
    report.PrepareAppInitActivities = readActivitiesInOldFormat("prepareAppInitActivities", data)
  }

  // Sort for better compression (same data pattern across column values). It is confirmed by experiment.
  sort.Slice(measures, func(i, j int) bool {
    return measures[i].name < measures[j].name
  })

  measureCount := len(measures)
  measureName := make([]string, measureCount, measureCount)
  measureStart := make([]uint32, measureCount, measureCount)
  measureDuration := make([]uint32, measureCount, measureCount)
  measureThread := make([]string, measureCount, measureCount)
  for i, info := range measures {
    measureName[i] = info.name
    measureStart[i] = info.start
    measureDuration[i] = info.duration
    measureThread[i] = info.thread
  }

  serviceCount := len(services)
  serviceName := make([]string, serviceCount, serviceCount)
  serviceStart := make([]uint32, serviceCount, serviceCount)
  serviceDuration := make([]uint32, serviceCount, serviceCount)
  serviceThread := make([]string, serviceCount, serviceCount)
  servicePlugin := make([]string, serviceCount, serviceCount)
  for i, info := range services {
    serviceName[i] = info.name
    serviceStart[i] = info.start
    serviceDuration[i] = info.duration
    serviceThread[i] = info.thread
    servicePlugin[i] = info.plugin
  }

  runResult.ExtraFieldData = []interface{}{
    serviceName, serviceStart, serviceDuration, serviceThread, servicePlugin,
    clTotal, clSearch, clDefine, clCount, rlTime, rlCount,
    measureName, measureStart, measureDuration, measureThread,
  }
  return nil
}

func readServices(
  data *fastjson.Value,
  category string,
  services *[]serviceItem,
) {
  for _, measure := range data.GetArray(category) {
    *services = append(*services, serviceItem{
      name:     string(measure.GetStringBytes("n")),
      start:    uint32(measure.GetInt("s")),
      duration: uint32(measure.GetInt("d")),
      thread:   string(measure.GetStringBytes("t")),
      plugin:   string(measure.GetStringBytes("p")),
    })
  }

  // remove to reduce size of raw report
  data.Del(category)
}

func readActivitiesInOldFormat(key string, data *fastjson.Value) []model.Activity {
  array := data.GetArray(key)
  result := make([]model.Activity, 0, len(array))
  for _, v := range array {
    result = append(result, model.Activity{
      Name:     string(v.GetStringBytes("name")),
      Thread:   string(v.GetStringBytes("thread")),
      Start:    v.GetInt("start"),
      End:      v.GetInt("end"),
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
