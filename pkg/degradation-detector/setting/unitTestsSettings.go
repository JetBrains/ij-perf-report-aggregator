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
			"com.intellij.diff", "com.intellij.openapi.vcs",
		},
	},
	{
		Team:         "java",
		SlackChannel: "idea-java-alerts",
		Packages: []string{
			"com.intellij.java", "org.jetbrains.plugins.groovy", "org.jetbrains.uast.test.java",
			"com.intellij.lang.properties", "com.intellij.structuralsearch", "org.jetbrains.java",
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
		Team:         "datagrip",
		SlackChannel: "datagrip",
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
	},
}

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

	// Iterate over team configurations
	for _, config := range teamConfigs {
		teamSettings := generateProductTestsSettings(
			tests,
			mainSettings,
			config,
		)
		settings = append(settings, teamSettings...)
	}

	// Collect all packages to exclude from default settings
	allPackages := collectAllPackages(teamConfigs)
	// Generate default settings (excluding specified packages)
	defaultSlackSettings := detector.SlackSettings{
		Channel:     "ij-perf-unit-tests-alerts",
		ProductLink: "perfUnit",
	}
	defaultTests := filterTests(tests, allPackages, false)
	for _, test := range defaultTests {
		settings = append(settings, detector.PerformanceSettings{
			Project: test,
			Db:      mainSettings.Db,
			Table:   mainSettings.Table,
			BaseSettings: detector.BaseSettings{
				Branch:        mainSettings.Branch,
				Machine:       mainSettings.Machine,
				Metric:        mainSettings.Metric,
				SlackSettings: defaultSlackSettings,
				AnalysisSettings: detector.AnalysisSettings{
					MinimumSegmentLength:      30,
					MedianDifferenceThreshold: 20,
					EffectSizeThreshold:       2,
				},
			},
		})
	}

	return settings
}

func generateProductTestsSettings(
	allTests []string,
	mainSettings detector.PerformanceSettings,
	config teamConfig,
) []detector.PerformanceSettings {
	// Filter tests for the team
	teamTests := filterTests(allTests, config.Packages, true)

	// Set default AnalysisSettings if nil
	if config.AnalysisSettings == nil {
		config.AnalysisSettings = &detector.AnalysisSettings{
			MinimumSegmentLength:      30,
			MedianDifferenceThreshold: 20,
			EffectSizeThreshold:       2,
		}
	}

	slackSettings := detector.SlackSettings{
		Channel:     config.SlackChannel,
		ProductLink: "perfUnit",
	}

	settings := make([]detector.PerformanceSettings, 0, len(teamTests))
	for _, test := range teamTests {
		settings = append(settings, detector.PerformanceSettings{
			Project: test,
			Db:      mainSettings.Db,
			Table:   mainSettings.Table,
			BaseSettings: detector.BaseSettings{
				Branch:           mainSettings.Branch,
				Machine:          mainSettings.Machine,
				Metric:           mainSettings.Metric,
				SlackSettings:    slackSettings,
				AnalysisSettings: *config.AnalysisSettings,
			},
		})
	}

	return settings
}

// filterTests returns a new slice of tests based on inclusion or exclusion of packages.
// If include is true, we keep tests that match the packages
// If include is false, we remove tests that match the packages
func filterTests(tests []string, packages []string, include bool) []string {
	var filtered []string
	for _, test := range tests {
		matches := false
		for _, pkg := range packages {
			if strings.HasPrefix(test, pkg+".") {
				matches = true
				break
			}
		}
		if (include && matches) || (!include && !matches) {
			filtered = append(filtered, test)
		}
	}
	return filtered
}

func collectAllPackages(teamConfigs []teamConfig) []string {
	var allPackages []string
	for _, config := range teamConfigs {
		allPackages = append(allPackages, config.Packages...)
	}
	return allPackages
}
