package setting

import detector "github.com/JetBrains/ij-perf-report-aggregator/pkg/degradation-detector"

// Hooks are the webhooks that get the new degradations of a dashboard besides its Slack messages.
//
// To add one, add an entry here and the two keys URLEnv and TokenEnv to the Kubernetes secret
// degradation-detector-hooks: the degradation-detector CronJob reads every key of that secret as an environment
// variable. A hook whose keys are missing sends nothing, so an entry can be merged before its secret exists.
func Hooks() []detector.Hook {
	return []detector.Hook{
		{
			// The Air Automation of the Air performance tests.
			Name:          "air-perf",
			Db:            "perfintDev",
			Table:         "air",
			ProjectPrefix: "air/",
			URLEnv:        "AIR_PERF_HOOK_URL",
			TokenEnv:      "AIR_PERF_HOOK_TOKEN",
			AuthScheme:    "ApiKey",
		},
	}
}
