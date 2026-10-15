package setting

import detector "github.com/JetBrains/ij-perf-report-aggregator/pkg/degradation-detector"

// Hooks are the webhooks that get the new degradations of a team besides its Slack messages.
//
// To add one:
//   - add an entry here with the Slack channel of the team and the names of its URL and token variables;
//   - add its secret name to degradationDetector.hookSecrets in deployment/report-aggregator/values.yaml;
//   - create that Kubernetes secret with the two keys.
//
// Each hook has its own secret, so a team creates and rotates its keys without the keys of the other hooks. A hook
// whose secret does not exist sends nothing, so an entry can be merged before its secret exists.
func Hooks() []detector.Hook {
	return []detector.Hook{
		{
			// The Air Automation of the Air performance tests.
			Name:          "air-perf",
			Channel:       "ij-air-perf",
			Db:            "perfintDev",
			Table:         "air",
			ProjectPrefix: "air/",
			Secret:        "degradation-detector-hook-air-perf",
			URLEnv:        "AIR_PERF_HOOK_URL",
			TokenEnv:      "AIR_PERF_HOOK_TOKEN",
			AuthScheme:    "ApiKey",
		},
	}
}
