package analyzer

// IjMetricNames are the startup metric columns of the ij/ijDev databases, in column order, as served by
// the /meta/measure endpoint. Collection into those databases stopped in August 2025, so the list is frozen.
var IjMetricNames = []string{
	"pluginDescriptorLoading_d",
	"projectProfileLoading_d",
	"editorRestoring",
	"appComponentCreation_d",
	"projectComponentCreation_d",
	"bootstrap_d",
	"appInitPreparation_d",
	"appInit_d",
	"pluginDescriptorInitV18_d",
	"projectFrameInit_d",
	"projectDumbAware",
	"editorRestoringTillPaint",
	"splash_i",
	"startUpCompleted",
	"appStarter_d",
	"euaShowing_d",
	"serviceSyncPreloading_d",
	"serviceAsyncPreloading_d",
	"projectServiceSyncPreloading_d",
	"projectServiceAsyncPreloading_d",
}
