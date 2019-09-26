package analyzer

import (
	"github.com/mcuadros/go-version"
	"go.uber.org/zap"
	"report-aggregator/pkg/model"
)

func (t *ReportAnalyzer) computeMetrics(report *model.Report, logger *zap.Logger) *model.Metrics {
	metrics := &model.Metrics{
		Bootstrap: -1,
		Splash:    -1,

		AppInitPreparation:       -1,
		AppInit:                  -1,
		PluginDescriptorsLoading: -1,

		AppComponentCreation:     -1,
		ProjectComponentCreation: -1,
		ModuleLoading:            -1,
	}

	if version.Compare(report.Version, "12", ">=") && len(report.TraceEvents) == 0 {
		logger.Warn("invalid report (due to opening second project?), report will be skipped")
		return nil
	}

	// v < 12: PluginDescriptorsLoading can be or in MainActivities, or in PrepareAppInitActivities

	for _, activity := range report.MainActivities {
		switch activity.Name {
		case "bootstrap":
			metrics.Bootstrap = activity.Duration

		case "app initialization preparation":
			metrics.AppInitPreparation = activity.Duration
		case "app initialization":
			metrics.AppInit = activity.Duration
		case "plugin descriptors loading":
			metrics.PluginDescriptorsLoading = activity.Duration

		case "app component creation":
			metrics.AppComponentCreation = activity.Duration
		case "project component creation":
			metrics.ProjectComponentCreation = activity.Duration
		case "module loading":
			metrics.ModuleLoading = activity.Duration
		}
	}

	if version.Compare(report.Version, "11", "<") {
		for _, activity := range report.PrepareAppInitActivities {
			switch activity.Name {
			case "plugin descriptors loading":
				metrics.PluginDescriptorsLoading = activity.Start
			case "splash initialization":
				metrics.Splash = activity.Start
			}
		}
	} else {
		for _, activity := range report.TraceEvents {
			if activity.Phase == "i" && (activity.Name == "splash" || activity.Name == "splash shown") {
				metrics.Splash = activity.Timestamp / 1000
			}
		}
	}

	if metrics.Bootstrap == -1 {
		logRequiredMetricNotFound(logger, "bootstrap")
		return nil
	}
	if metrics.PluginDescriptorsLoading == -1 {
		logRequiredMetricNotFound(logger, "pluginDescriptorsLoading")
		return nil
	}
	if metrics.AppComponentCreation == -1 {
		logRequiredMetricNotFound(logger, "AppComponentCreation")
		return nil
	}
	if metrics.ModuleLoading == -1 {
		logRequiredMetricNotFound(logger, "ModuleLoading")
		return nil
	}
	return metrics
}

func logRequiredMetricNotFound(logger *zap.Logger, metricName string) {
	logger.Error("metric is required, but not found, report will be skipped", zap.String("metric", metricName))
}
