package setting

import (
	"log/slog"
	"net/http"
	"strings"

	detector "github.com/JetBrains/ij-perf-report-aggregator/pkg/degradation-detector"
)

// teamConfig represents the configuration for a team
type teamConfig struct {
	Team             string
	Packages         []string
	SlackChannel     string
	AnalysisSettings *detector.AnalysisSettings
	// AdditionalTestMetrics maps a test class name to extra metrics to assert alongside the
	// default attempt.mean.ms. Values may be SQL LIKE patterns (e.g. "%.expected.%").
	AdditionalTestMetrics map[string][]string
	// Mention is an optional Slack mention prepended to this team's degradation messages.
	// It only fires when the alert lands in this team's own SlackChannel (owner routing can send an
	// alert elsewhere, in which case the ping is dropped). Use the raw Slack syntax, e.g.
	// "<@U01ABC2DEF>" for a user or "<!subteam^S01ABC2DEF>" for a user group.
	Mention string
}

var defaultUnitTestAnalysisSettings = detector.AnalysisSettings{
	MinimumSegmentLength:      30,
	MedianDifferenceThreshold: 20,
	EffectSizeThreshold:       2,
}

func degradationOnlyAnalysisSettings() *detector.AnalysisSettings {
	s := defaultUnitTestAnalysisSettings
	s.ReportType = detector.DegradationEvent
	return &s
}

var teamConfigs = []teamConfig{
	{
		Team:         "ultimate",
		SlackChannel: "ij-u-team-performance-issues-check",
		Packages: []string{
			"com.intellij.freemarker", "com.intellij.reactivestreams", "com.intellij.httpClient", "com.intellij.spring",
			"com.intellij.swagger", "com.intellij.util.xml", "com.intellij.grpc", "com.intellij.guice",
			"com.intellij.helidon", "com.intellij.hibernate", "com.intellij.jdbi", "com.intellij.ktor",
			"com.intellij.micronaut", "com.intellij.openRewrite", "com.intellij.quarkus", "com.intellij.frameworks.thymeleaf",
			"com.intellij.wiremock", "com.intellij.ee", "com.intellij.jsp", "com.jetbrains.jsonSchema",
		},
	},
	{
		Team:         "vcs",
		SlackChannel: "vcs-perf-tests",
		Packages: []string{
			"com.intellij.diff", "com.intellij.openapi.vcs", "git4idea",
		},
	},
	{
		Team:         "java",
		SlackChannel: "idea-java-alerts",
		Packages: []string{
			"com.intellij.java", "org.jetbrains.plugins.groovy", "org.jetbrains.uast.test.java",
			"com.intellij.lang.properties", "com.intellij.structuralsearch", "org.jetbrains.java",
			"org.jetbrains.plugins.cucumber.java",
		},
	},
	{
		Team:         "cloudsAndDeployment",
		SlackChannel: "kubernetes-ui-tests-failures",
		Packages: []string{
			"com.intellij.kubernetes", "org.jetbrains.yaml",
		},
	},
	{
		Team:         "cloudsAndDeployment",
		SlackChannel: "ij-clouds-team-performance-issues-check",
		Packages: []string{
			"com.intellij.performance.json",
		},
	},
	{
		Team:         "datagrip",
		SlackChannel: "dbe-failed-tests",
		Packages: []string{
			"com.intellij.sql",
		},
	},
	{
		Team:         "webstorm",
		SlackChannel: "webstorm-test-degradations",
		Packages: []string{
			"com.intellij.lang.javascript", "com.intellij.html", "org.angular2.lang", "css", "org.jetbrains.vuejs",
			"org.angularjs.performance", "com.intellij.flex.completion", "com.intellij.htmltools",
			"org.jetbrains.plugins.sass", "org.jetbrains.plugins.scss", "org.jetbrains.astro.lang", "com.intellij.psi.html",
			"com.intellij.plugins.watcher", "com.intellij.javascript",
		},
	},
	{
		Team:         "buildTools",
		SlackChannel: "ij-build-tools-perf-tests",
		Packages: []string{
			"org.jetbrains.osgi.maven", "org.jetbrains.plugins.gradle",
		},
	},
	{
		Team:         "rubymine",
		SlackChannel: "rubymine-performance-alerts",
		Packages: []string{
			"org.jetbrains.plugins.ruby",
		},
		AnalysisSettings: degradationOnlyAnalysisSettings(),
	},
	{
		Team:         "phpstorm",
		SlackChannel: "phpstorm-performance-degradations",
		Packages: []string{
			"com.jetbrains.php",
		},
		AnalysisSettings: degradationOnlyAnalysisSettings(),
		Mention:          "<!subteam^S0BH689GW9G>", // @phpstorm-dev-duty
	},
	{
		Team:         "lsp",
		SlackChannel: "kotlin-lsp-alerts",
		Packages: []string{
			"com.jetbrains.ls",
		},
	},
	{
		Team:         "debugger",
		SlackChannel: "debugger-perf-tests",
		Packages: []string{
			"com.intellij.debugger", "org.jetbrains.kotlin.idea.k2.debugger",
		},
		AdditionalTestMetrics: map[string][]string{
			"com.intellij.debugger.impl.PacketsNumberTest":                                                      {"%.expected.%"},
			"org.jetbrains.kotlin.idea.k2.debugger.test.performance.K2IdeK2CodeKotlinSteppingPacketsNumberTest": {"%.expected.%"},
		},
	},
}

const catchAllUnitTestsChannel = "ij-perf-unit-tests-alerts"

func GenerateAllUnitTestsSettings(backendUrl string, client *http.Client) []detector.PerformanceSettings {
	settings := make([]detector.PerformanceSettings, 0, 1000)

	mainSettings := detector.PerformanceSettings{
		Db:    "perfUnitTests",
		Table: "report",
		BaseSettings: detector.BaseSettings{
			Branch:  "master",
			Machine: "intellij-linux-%-hetzner-%",
			Metric:  "attempt.mean.ms",
		},
	}

	tests, err := detector.FetchAllTests(backendUrl, client, mainSettings)
	if err != nil {
		slog.Error("error while getting tests", "error", err)
		return settings
	}

	// Code owner is the primary attribution key: a test's owner (from project_owner metadata)
	// is resolved to a Slack channel via the CodeOwners service. Either fetch failing degrades
	// to package-based routing rather than dropping notifications entirely.
	projectOwners, err := detector.FetchProjectOwners(backendUrl, client, mainSettings.Db, mainSettings.Table)
	if err != nil {
		slog.Error("error while fetching project owners; falling back to package-based routing", "error", err)
		projectOwners = map[string]string{}
	}
	ownerChannels, err := detector.FetchCodeOwnerChannels(backendUrl, client)
	if err != nil {
		slog.Error("error while fetching code-owner channels; falling back to package-based routing", "error", err)
		ownerChannels = map[string]string{}
	}

	for _, test := range tests {
		route := resolveRoute(test, projectOwners, ownerChannels, teamConfigs)
		metrics := append([]string{mainSettings.Metric}, expandAdditionalMetrics(backendUrl, client, mainSettings, test, route.additionalTestMetrics)...)
		for _, metric := range metrics {
			settings = append(settings, detector.PerformanceSettings{
				Project: test,
				Db:      mainSettings.Db,
				Table:   mainSettings.Table,
				BaseSettings: detector.BaseSettings{
					Branch:           mainSettings.Branch,
					Machine:          mainSettings.Machine,
					Metric:           metric,
					SlackSettings:    detector.SlackSettings{Channel: route.channel, ProductLink: "perfUnit", Mention: route.mention},
					AnalysisSettings: route.analysisSettings,
				},
			})
		}
	}

	return settings
}

// resolvedRoute is the notification destination chosen for a single test.
type resolvedRoute struct {
	channel               string
	analysisSettings      detector.AnalysisSettings
	additionalTestMetrics map[string][]string
	mention               string
}

// resolveRoute determines where a test's notifications go and how it is analyzed.
//
// Channel precedence: the test's code owner resolved via the CodeOwners service (primary) ->
// the first package-prefix match in teamConfigs (fallback) -> the catch-all channel.
//
// Analysis settings and additional metrics always come from the matching package teamConfig
// (so e.g. RubyMine/PhpStorm degradation-only analysis and debugger packet metrics are kept
// regardless of which channel wins); they are independent of the channel decision.
//
// The mention also originates from the package teamConfig, but is a team-specific ping, so it is
// dropped when owner routing redirects the alert away from that team's own channel.
func resolveRoute(test string, projectOwners, ownerChannels map[string]string, configs []teamConfig) resolvedRoute {
	teamCfg := matchTeamConfig(test, configs)

	route := resolvedRoute{
		channel:          catchAllUnitTestsChannel,
		analysisSettings: defaultUnitTestAnalysisSettings,
	}
	if teamCfg != nil {
		route.channel = teamCfg.SlackChannel
		route.analysisSettings = analysisSettingsOrDefault(teamCfg.AnalysisSettings)
		route.additionalTestMetrics = teamCfg.AdditionalTestMetrics
		route.mention = teamCfg.Mention
	}

	// Owner-based channel takes precedence over the package fallback.
	if owner := projectOwners[test]; owner != "" {
		if channel := ownerChannels[owner]; channel != "" {
			route.channel = channel
		}
	}

	// A mention pings a team-specific group, so only keep it when the alert actually lands in that
	// team's own channel. If owner routing redirected the alert elsewhere, drop the ping.
	if teamCfg != nil && route.channel != teamCfg.SlackChannel {
		route.mention = ""
	}

	return route
}

// matchTeamConfig returns the first teamConfig whose packages match the test, or nil.
func matchTeamConfig(test string, configs []teamConfig) *teamConfig {
	for i := range configs {
		if testMatchesPackages(test, configs[i].Packages) {
			return &configs[i]
		}
	}
	return nil
}

func analysisSettingsOrDefault(s *detector.AnalysisSettings) detector.AnalysisSettings {
	if s == nil {
		return defaultUnitTestAnalysisSettings
	}
	return *s
}

func expandAdditionalMetrics(backendUrl string, client *http.Client, mainSettings detector.PerformanceSettings, test string, additionalMetrics map[string][]string) []string {
	probeSettings := mainSettings
	probeSettings.Project = test
	var expanded []string
	for testPrefix, patterns := range additionalMetrics {
		if !strings.HasPrefix(test, testPrefix) {
			continue
		}
		for _, pattern := range patterns {
			metrics, err := detector.FetchMetricNamesByPattern(backendUrl, client, probeSettings, pattern)
			if err != nil {
				slog.Error("error while fetching metric names by pattern", "error", err, "test", test, "pattern", pattern)
				continue
			}
			expanded = append(expanded, metrics...)
		}
	}
	return expanded
}

// testMatchesPackages reports whether the test belongs to one of the given packages,
// matched by test-class-name prefix.
func testMatchesPackages(test string, packages []string) bool {
	for _, pkg := range packages {
		if strings.HasPrefix(test, pkg+".") {
			return true
		}
	}
	return false
}
