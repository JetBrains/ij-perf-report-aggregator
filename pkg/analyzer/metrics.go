package analyzer

import (
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/model"
  "github.com/develar/errors"
  "github.com/mcuadros/go-version"
  "go.uber.org/zap"
)

type Metric struct {
  Name     string
  index    int
  category int

  isRequired bool
  IsInstant  bool

  sinceVersion string

  maxValue int
}

const appInitCategory = 1

var metricNameToDescriptor map[string]*Metric
var IjMetricDescriptors []*Metric

const clTotal = "class loading: total"
const clSearch = "class loading: search"
const clDefine = "class loading: define"
const clCount = "class loading: count"

const rlTime = "resource loading: time"
const rlCount = "resource loading: count"

func init() {
  index := 0
  createMetric := func(name string) *Metric {
    result := &Metric{
      Name:  name,
      index: index,
      maxValue: 65535,
    }
    index++

    IjMetricDescriptors = append(IjMetricDescriptors, result)
    return result
  }

  createVersionedMetric := func(name string, sinceVersion string) *Metric {
    result := createMetric(name)
    result.sinceVersion = sinceVersion
    return result
  }

  createRequiredMetric := func(name string) *Metric {
    result := createMetric(name)
    result.maxValue = 2147483647
    result.isRequired = true
    return result
  }

  createUint32Metric := func(name string) *Metric {
    result := createMetric(name)
    result.maxValue = 4294967295
    return result
  }

  createUint32RequiredMetric := func(name string) *Metric {
    result := createMetric(name)
    result.maxValue = 4294967295
    result.isRequired = true
    return result
  }

  createMetricWithCategory := func(name string, category int) *Metric {
    result := createMetric(name)
    result.category = category
    return result
  }

  createInstantMetric := func(name string) *Metric {
    result := createMetric(name)
    result.IsInstant = true
    return result
  }

  pluginDescriptorLoading := createMetric("pluginDescriptorLoading_d")
  projectProfileLoading := createMetricWithCategory("projectProfileLoading_d", appInitCategory)
  editorRestoring := createMetric("editorRestoring_d")

  appComponentCreation := createMetric("appComponentCreation_d")
  projectComponentCreation := createMetric("projectComponentCreation_d")

  metricNameToDescriptor = map[string]*Metric{
    "bootstrap":                      createUint32RequiredMetric("bootstrap_d"),
    "app initialization preparation": createRequiredMetric("appInitPreparation_d"),
    "app initialization":             createRequiredMetric("appInit_d"),

    "plugin descriptor loading": pluginDescriptorLoading,
    // old name
    "plugin descriptors loading": pluginDescriptorLoading,
    "plugin initialization": createVersionedMetric("pluginDescriptorInitV18_d", "18"),

    "app component creation":  appComponentCreation,
    "app components creation": appComponentCreation,

    "project component creation":  projectComponentCreation,
    "project components creation": projectComponentCreation,

    "project frame initialization": createMetricWithCategory("projectFrameInit_d", appInitCategory),

    "project inspection profile loading": projectProfileLoading,
    // old name
    "project inspection profiles loading": projectProfileLoading,

    // light edit mode doesn't have moduleLoading phase
    "module loading": createMetric("moduleLoading_d"),
    "project post-startup dumb-aware activities": createUint32Metric("projectDumbAware_d"),

    "editor restoring":            editorRestoring,
    "editor restoring till paint": createMetricWithCategory("editorRestoringTillPaint_d", appInitCategory),
    // old name
    "restoring editors": editorRestoring,

    // instant
    "splash initialization": createInstantMetric("splash_i"),
    "startUpCompleted":      createInstantMetric("startUpCompleted_i"),

    "appStarter": createMetric("appStarter_d"),
    // v19+
    "eua showing": createVersionedMetric("euaShowing_d", "19"),

    "service sync preloading":          createUint32Metric("serviceSyncPreloading_d"),
    "service async preloading":         createUint32Metric("serviceAsyncPreloading_d"),
    "project service sync preloading":  createUint32Metric("projectServiceSyncPreloading_d"),
    "project service async preloading": createUint32Metric("projectServiceAsyncPreloading_d"),

    // v27+
    clTotal:  createInstantMetric("classLoadingTime_i"),
    clSearch: createInstantMetric("classLoadingSearchTime_i"),
    clDefine: createInstantMetric("classLoadingDefineTime_i"),
    clCount:  createInstantMetric("classLoadingCount_i"),

    rlTime:  createInstantMetric("resourceLoadingTime_i"),
    rlCount: createInstantMetric("resourceLoadingCount_i"),
  }
}

func ComputeIjMetrics(nonMetricFieldCount int, report *model.Report, result *[]interface{}, logger *zap.Logger) error {
  if version.Compare(report.Version, "12", ">=") && len(report.TraceEvents) == 0 {
    logger.Warn("invalid report (due to opening second project?), report will be skipped")
    return nil
  }

  for range IjMetricDescriptors {
    *result = append(*result, -1)
  }

  (*result)[nonMetricFieldCount+metricNameToDescriptor["startUpCompleted"].index] = report.TotalDurationActual

  for _, activity := range report.MainActivities {
    err := setMetric(nonMetricFieldCount, activity, report, result)
    if err != nil {
      return err
    }
  }

  if version.Compare(report.Version, "18", ">=") {
    for _, activity := range report.PrepareAppInitActivities {
      err := setMetric(nonMetricFieldCount, activity, report, result)
      if err != nil {
        return err
      }
    }
  } else {
    for _, activity := range report.PrepareAppInitActivities {
      switch activity.Name {
      case "plugin descriptors loading":
        (*result)[nonMetricFieldCount+metricNameToDescriptor["plugin descriptor loading"].index] = activity.Duration
      default:
        err := setMetric(nonMetricFieldCount, activity, report, result)
        if err != nil {
          return err
        }
      }
    }
  }

  if version.Compare(report.Version, "11", ">=") {
    for _, activity := range report.TraceEvents {
      if activity.Phase == "i" && (activity.Name == "splash" || activity.Name == "splash shown") {
        (*result)[nonMetricFieldCount+metricNameToDescriptor["splash initialization"].index] = activity.Timestamp / 1000
      }
    }
  }

  is14orGreater := version.Compare(report.Version, "14", ">=")

  var notFoundMetrics []string
  for _, metric := range IjMetricDescriptors {
    if (*result)[nonMetricFieldCount+metric.index] != -1 {
      continue
    }

    if metric.isRequired {
      if metric.Name != "bootstrap_d" || version.Compare(report.Version, "6", ">=") {
        logRequiredMetricNotFound(logger, metric.Name)
        return nil
      }
    }

    // undefined
    (*result)[nonMetricFieldCount+metric.index] = 0
    if is14orGreater || (metric.Name != "editorRestoringTillPaint_d" && metric.Name != "projectProfileLoading_d") {
      if len(metric.sinceVersion) != 0 && version.Compare(report.Version, metric.sinceVersion, ">=") {
        notFoundMetrics = append(notFoundMetrics, metric.Name)
      }
    }
  }

  if len(notFoundMetrics) > 0 {
    logger.Debug("metrics not found", zap.Strings("name", notFoundMetrics))
  }

  return nil
}

func setMetric(nonMetricFieldCount int, activity model.Activity, report *model.Report, result *[]interface{}) error {
  info, ok := metricNameToDescriptor[activity.Name]
  if !ok {
    return nil
  }

  if len(info.sinceVersion) != 0 && version.Compare(report.Version, info.sinceVersion, "<") {
    return nil
  }

  var v int
  if info.IsInstant {
    v = activity.Start
  } else {
    v = activity.Duration
    if v > info.maxValue {
      return errors.Errorf("value outside of 0-%d range (generatedTime=%s, value=%v, activity=%s)", info.maxValue, report.Generated, v, activity.Name)
    }
  }

  if v < 0 {
    return errors.Errorf("value must be positive (generatedTime: %s, value: %v)", report.Generated, v)
  }
  (*result)[nonMetricFieldCount+info.index] = v
  return nil
}

func logRequiredMetricNotFound(logger *zap.Logger, metricName string) {
  logger.Error("metric is required, but not found, report will be skipped", zap.String("metric", metricName))
}
