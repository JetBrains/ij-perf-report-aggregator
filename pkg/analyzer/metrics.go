package analyzer

import (
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/model"
  "github.com/mcuadros/go-version"
  "go.uber.org/zap"
)

func ComputeMetrics(report *model.Report, logger *zap.Logger) (*model.DurationEventMetrics, *model.InstantEventMetrics) {
  durationMetrics := &model.DurationEventMetrics{
    Bootstrap: -1,

    AppInitPreparation:      -1,
    AppInit:                 -1,
    PluginDescriptorLoading: -1,

    AppComponentCreation:     -1,
    ProjectComponentCreation: -1,

    ProjectDumbAware: -1,

    ModuleLoading:   -1,
    EditorRestoring: -1,
  }

  instantMetrics := &model.InstantEventMetrics{
    Splash:           -1,
    StartUpCompleted: report.TotalDurationActual,
  }

  if version.Compare(report.Version, "12", ">=") && len(report.TraceEvents) == 0 {
    logger.Warn("invalid report (due to opening second project?), report will be skipped")
    return nil, nil
  }

  // v < 12: PluginDescriptorLoading can be or in MainActivities, or in PrepareAppInitActivities

  for _, activity := range report.MainActivities {
    switch activity.Name {
    case "bootstrap":
      durationMetrics.Bootstrap = activity.Duration

    case "app initialization preparation":
      durationMetrics.AppInitPreparation = activity.Duration
    case "app initialization":
      durationMetrics.AppInit = activity.Duration

    case "plugin descriptor loading":
      durationMetrics.PluginDescriptorLoading = activity.Duration
    // old name
    case "plugin descriptors loading":
      durationMetrics.PluginDescriptorLoading = activity.Duration

    case "app component creation":
      durationMetrics.AppComponentCreation = activity.Duration
    case "app components creation":
      durationMetrics.AppComponentCreation = activity.Duration

    case "project component creation":
      durationMetrics.ProjectComponentCreation = activity.Duration
    case "project components creation":
      durationMetrics.ProjectComponentCreation = activity.Duration

    case "project post-startup dumb-aware activities":
      durationMetrics.ProjectDumbAware = activity.Duration

    case "module loading":
      durationMetrics.ModuleLoading = activity.Duration

    case "editor restoring":
      durationMetrics.EditorRestoring = activity.Duration
    // old name
    case "restoring editors":
      durationMetrics.EditorRestoring = activity.Duration
    }
  }

  if version.Compare(report.Version, "11", "<") {
    for _, activity := range report.PrepareAppInitActivities {
      switch activity.Name {
      case "plugin descriptors loading":
        durationMetrics.PluginDescriptorLoading = activity.Start
      case "splash initialization":
        instantMetrics.Splash = activity.Start
      }
    }
  } else {
    for _, activity := range report.TraceEvents {
      if activity.Phase == "i" && (activity.Name == "splash" || activity.Name == "splash shown") {
        instantMetrics.Splash = activity.Timestamp / 1000
      }
    }
  }

  if instantMetrics.StartUpCompleted == -1 {
    logRequiredMetricNotFound(logger, "startUpCompleted")
  }
  if instantMetrics.Splash == -1 && version.Compare(report.Version, "6", ">=") {
    logger.Info("metric 'splash' not found")
  }

  if durationMetrics.Bootstrap == -1 {
    if version.Compare(report.Version, "6", ">=") {
      logRequiredMetricNotFound(logger, "bootstrap")
    }
  }
  if durationMetrics.PluginDescriptorLoading == -1 {
    logRequiredMetricNotFound(logger, "pluginDescriptorsLoading")
    return nil, nil
  }
  if durationMetrics.AppComponentCreation == -1 {
    logRequiredMetricNotFound(logger, "AppComponentCreation")
    return nil, nil
  }
  if durationMetrics.ModuleLoading == -1 {
    logRequiredMetricNotFound(logger, "ModuleLoading")
    return nil, nil
  }
  if durationMetrics.ProjectDumbAware == -1 {
    // not reported because ended too late?
    logger.Info("metric 'ProjectDumbAware' not found")
  }
  if durationMetrics.EditorRestoring == -1 {
    // not reported because ended too late?
    logger.Info("metric 'EditorRestoring' not found")
  }
  return durationMetrics, instantMetrics
}

func logRequiredMetricNotFound(logger *zap.Logger, metricName string) {
  logger.Error("metric is required, but not found, report will be skipped", zap.String("metric", metricName))
}
