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
}

const appInitCategory = 1

var metricNameToDescriptor map[string]*Metric
var MetricDescriptors []*Metric

const nonMetricFieldCount = 13

func init() {
  index := 0
  createMetric := func(name string) *Metric {
    result := &Metric{
      Name:  name,
      index: index,
    }
    index++

    MetricDescriptors = append(MetricDescriptors, result)
    return result
  }

  createRequiredMetric := func(name string) *Metric {
    result := createMetric(name)
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

  pluginDescriptorLoading := createRequiredMetric("pluginDescriptorLoading_d")

  projectProfileLoading := createMetricWithCategory("projectProfileLoading_d", appInitCategory)
  editorRestoring := createMetric("editorRestoring_d")

  appComponentCreation := createMetric("appComponentCreation_d")
  projectComponentCreation := createMetric("projectComponentCreation_d")

  metricNameToDescriptor = map[string]*Metric{
    "bootstrap":                      createRequiredMetric("bootstrap_d"),
    "app initialization preparation": createRequiredMetric("appInitPreparation_d"),
    "app initialization":             createRequiredMetric("appInit_d"),

    "plugin descriptor loading": pluginDescriptorLoading,
    // old name
    "plugin descriptors loading": pluginDescriptorLoading,

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
    "project post-startup dumb-aware activities": createMetric("projectDumbAware_d"),

    "editor restoring":            editorRestoring,
    "editor restoring till paint": createMetricWithCategory("editorRestoringTillPaint_d", appInitCategory),
    // old name
    "restoring editors": editorRestoring,

    // instant
    "splash initialization": createInstantMetric("splash_i"),
    "startUpCompleted":      createInstantMetric("startUpCompleted_i"),
  }
}

func ComputeMetrics(report *model.Report, result *[]interface{}, logger *zap.Logger) error {
  if version.Compare(report.Version, "12", ">=") && len(report.TraceEvents) == 0 {
    logger.Warn("invalid report (due to opening second project?), report will be skipped")
    return nil
  }

  for range MetricDescriptors {
    *result = append(*result, -1)
  }

  (*result)[nonMetricFieldCount+metricNameToDescriptor["startUpCompleted"].index] = report.TotalDurationActual

  // v < 12: PluginDescriptorLoading can be or in MainActivities, or in PrepareAppInitActivities

  for _, activity := range report.MainActivities {
    err := setMetric(activity, report, result)
    if err != nil {
      return err
    }
  }

  for _, activity := range report.PrepareAppInitActivities {
    switch activity.Name {
    case "plugin descriptors loading":
      (*result)[nonMetricFieldCount+metricNameToDescriptor["plugin descriptor loading"].index] = activity.Duration
    default:
      err := setMetric(activity, report, result)
      if err != nil {
        return err
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
  for _, metric := range MetricDescriptors {
    if (*result)[nonMetricFieldCount+metric.index] == -1 {
      if metric.isRequired {
        if metric.Name != "bootstrap_d" || version.Compare(report.Version, "6", ">=") {
          logRequiredMetricNotFound(logger, metric.Name)
          return nil
        }
      }

      // undefined
      (*result)[nonMetricFieldCount+metric.index] = 0
      if is14orGreater || (metric.Name != "editorRestoringTillPaint_d" && metric.Name != "projectProfileLoading_d") {
        notFoundMetrics = append(notFoundMetrics, metric.Name)
      }
    }
  }

  if len(notFoundMetrics) > 0 {
    logger.Debug("metrics not found", zap.Strings("name", notFoundMetrics))
  }

  return nil
}

func setMetric(activity model.Activity, report *model.Report, result *[]interface{}) error {
  info, ok := metricNameToDescriptor[activity.Name]
  if ok {
    var v int
    if info.IsInstant {
      v = activity.Start
    } else {
      v = activity.Duration
      if v > 65535 {
        return errors.Errorf("value outside of uint16 range (generatedTime: %s, value: %v)", report.Generated, v)
      }
    }

    if v < 0 {
      return errors.Errorf("value must be positive (generatedTime: %s, value: %v)", report.Generated, v)
    }
    (*result)[nonMetricFieldCount+info.index] = v
  }
  return nil
}

func logRequiredMetricNotFound(logger *zap.Logger, metricName string) {
  logger.Error("metric is required, but not found, report will be skipped", zap.String("metric", metricName))
}
