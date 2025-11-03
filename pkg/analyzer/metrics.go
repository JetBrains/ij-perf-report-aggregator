package analyzer

import (
	"fmt"
	"log/slog"

	"github.com/JetBrains/ij-perf-report-aggregator/pkg/model"
	"github.com/mcuadros/go-version"
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

var (
	metricNameToDescriptor map[string]*Metric
	IjMetricDescriptors    []*Metric
)

func init() {
	index := 0
	createMetric := func(name string) *Metric {
		result := &Metric{
			Name:     name,
			index:    index,
			maxValue: 65535,
		}
		index++

		IjMetricDescriptors = append(IjMetricDescriptors, result)
		return result
	}

	createVersionedUint16Metric := func(name string, sinceVersion string) *Metric {
		result := createMetric(name)
		result.sinceVersion = sinceVersion
		result.maxValue = 65535
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

	createInt32Metric := func(name string) *Metric {
		result := createMetric(name)
		result.maxValue = 2147483647
		return result
	}

	createUint16Metric := func(name string) *Metric {
		result := createMetric(name)
		result.maxValue = 65535
		return result
	}

	createUint32RequiredMetric := func(name string) *Metric {
		result := createMetric(name)
		result.maxValue = 4294967295
		result.isRequired = true
		return result
	}

	createUint16MetricWithCategory := func(name string, category int) *Metric {
		result := createMetric(name)
		result.category = category
		result.maxValue = 65535
		return result
	}

	createInt32MetricWithCategory := func(name string, category int) *Metric {
		result := createMetric(name)
		result.maxValue = 2147483647
		result.category = category
		return result
	}

	createInstantMetric := func(name string) *Metric {
		result := createMetric(name)
		result.IsInstant = true
		result.maxValue = 2147483647
		return result
	}

	pluginDescriptorLoading := createUint16Metric("pluginDescriptorLoading_d")
	projectProfileLoading := createUint16MetricWithCategory("projectProfileLoading_d", appInitCategory)
	editorRestoring := createInt32Metric("editorRestoring")

	appComponentCreation := createUint16Metric("appComponentCreation_d")
	projectComponentCreation := createUint16Metric("projectComponentCreation_d")

	metricNameToDescriptor = map[string]*Metric{
		"bootstrap":                      createUint32RequiredMetric("bootstrap_d"),
		"app initialization preparation": createRequiredMetric("appInitPreparation_d"),
		"app initialization":             createRequiredMetric("appInit_d"),

		"plugin descriptor loading": pluginDescriptorLoading,
		// old name
		"plugin descriptors loading": pluginDescriptorLoading,
		"plugin initialization":      createVersionedUint16Metric("pluginDescriptorInitV18_d", "18"),

		"app component creation":  appComponentCreation,
		"app components creation": appComponentCreation,

		"project component creation":  projectComponentCreation,
		"project components creation": projectComponentCreation,

		"project frame initialization": createUint16MetricWithCategory("projectFrameInit_d", appInitCategory),

		"project inspection profile loading": projectProfileLoading,
		// old name
		"project inspection profiles loading": projectProfileLoading,

		"project post-startup dumb-aware activities": createInt32Metric("projectDumbAware"),

		"editor restoring":            editorRestoring,
		"editor restoring till paint": createInt32MetricWithCategory("editorRestoringTillPaint", appInitCategory),
		// old name
		"restoring editors": editorRestoring,

		// instant
		"splash initialization": createInstantMetric("splash_i"),
		"startUpCompleted":      createInstantMetric("startUpCompleted"),

		"appStarter": createUint16Metric("appStarter_d"),
		// v19+
		"eua showing": createVersionedUint16Metric("euaShowing_d", "19"),

		"service sync preloading":          createUint32Metric("serviceSyncPreloading_d"),
		"service async preloading":         createUint32Metric("serviceAsyncPreloading_d"),
		"project service sync preloading":  createUint32Metric("projectServiceSyncPreloading_d"),
		"project service async preloading": createUint32Metric("projectServiceAsyncPreloading_d"),
	}
}

func ComputeIjMetrics(nonMetricFieldCount int, report *model.Report, result *[]any, logger *slog.Logger) error {
	for _, info := range IjMetricDescriptors {
		switch info.maxValue {
		case 65535:
			*result = append(*result, uint16(0))
		case 4294967295:
			*result = append(*result, uint32(0))
		case 2147483647:
			*result = append(*result, int32(-1))
		default:
			*result = append(*result, -1)
		}
	}

	(*result)[nonMetricFieldCount+metricNameToDescriptor["startUpCompleted"].index] = int32(report.TotalDuration)

	for _, activity := range report.Activities {
		err := setMetric(nonMetricFieldCount, activity, report, result)
		if err != nil {
			return err
		}
	}

	for _, activity := range report.PrepareAppInitActivities {
		switch activity.Name {
		case "plugin descriptors loading":
			(*result)[nonMetricFieldCount+metricNameToDescriptor["plugin descriptor loading"].index] = uint16(activity.Duration)
		default:
			err := setMetric(nonMetricFieldCount, activity, report, result)
			if err != nil {
				return err
			}
		}
	}

	for _, activity := range report.TraceEvents {
		if activity.Phase == "i" && (activity.Name == "splash" || activity.Name == "splash shown") {
			(*result)[nonMetricFieldCount+metricNameToDescriptor["splash initialization"].index] = int32(activity.Timestamp / 1000)
		}
	}

	var notFoundMetrics []string
	for _, metric := range IjMetricDescriptors {
		if (*result)[nonMetricFieldCount+metric.index] != -1 {
			continue
		}

		if metric.isRequired {
			if metric.Name != "bootstrap_d" {
				logger.Error("metric is required, but not found, report will be skipped", "metric", metric.Name)
				return nil
			}
		}

		// undefined
		(*result)[nonMetricFieldCount+metric.index] = 0
		if metric.sinceVersion != "" && version.Compare(report.Version, metric.sinceVersion, ">=") {
			notFoundMetrics = append(notFoundMetrics, metric.Name)
		}
	}

	if len(notFoundMetrics) > 0 {
		logger.Info("metrics not found", "name", notFoundMetrics)
	}

	return nil
}

func setMetric(nonMetricFieldCount int, activity model.Activity, report *model.Report, result *[]any) error {
	info, ok := metricNameToDescriptor[activity.Name]
	if !ok {
		return nil
	}

	if info.sinceVersion != "" && version.Compare(report.Version, info.sinceVersion, "<") {
		return nil
	}

	var v int
	if info.IsInstant {
		v = activity.Start
	} else {
		v = activity.Duration
		if v > info.maxValue {
			return fmt.Errorf("value outside of 0-%d range (generatedTime=%s, value=%v, activity=%s)", info.maxValue, report.Generated, v, activity.Name)
		}
	}

	if v < 0 {
		return fmt.Errorf("value must be positive (generatedTime: %s, value: %v)", report.Generated, v)
	}

	switch info.maxValue {
	case 65535:
		(*result)[nonMetricFieldCount+info.index] = uint16(v)
	case 4294967295:
		(*result)[nonMetricFieldCount+info.index] = uint32(v)
	case 2147483647:
		(*result)[nonMetricFieldCount+info.index] = int32(v)
	default:
		(*result)[nonMetricFieldCount+info.index] = v
	}

	return nil
}
