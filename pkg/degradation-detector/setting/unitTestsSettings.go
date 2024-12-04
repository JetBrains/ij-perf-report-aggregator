package setting

import (
	"log/slog"
	"net/http"
	"slices"
	"strings"

	detector "github.com/JetBrains/ij-perf-report-aggregator/pkg/degradation-detector"
)

func GenerateUnitTestsSettings(backendUrl string, client *http.Client) []detector.PerformanceSettings {
	settings := make([]detector.PerformanceSettings, 0, 1000)
	mainSettings := detector.PerformanceSettings{
		Db:      "perfUnitTests",
		Table:   "report",
		Branch:  "master",
		Machine: "intellij-linux-%-hetzner-%",
		Metric:  "attempt.mean.ms",
	}
	slackSettings := detector.SlackSettings{
		Channel:     "ij-perf-unit-tests-alerts",
		ProductLink: "perfUnit",
	}
	tests, err := detector.FetchAllTests(backendUrl, client, mainSettings)
	if err != nil {
		slog.Error("error while getting tests", "error", err)
		return settings
	}

	tests = slices.DeleteFunc(tests, func(test string) bool {
		return slices.ContainsFunc(ultimatePackages, func(pkg string) bool {
			return strings.HasPrefix(test, pkg)
		})
	})

	for _, test := range tests {
		settings = append(settings, detector.PerformanceSettings{
			Project:       test,
			Db:            mainSettings.Db,
			Table:         mainSettings.Table,
			Branch:        mainSettings.Branch,
			Machine:       mainSettings.Machine,
			Metric:        mainSettings.Metric,
			SlackSettings: slackSettings,
			AnalysisSettings: detector.AnalysisSettings{
				MinimumSegmentLength:      30,
				MedianDifferenceThreshold: 20,
				EffectSizeThreshold:       2,
			},
		})
	}
	return settings
}

var ultimatePackages = []string{
	"com.intellij.freemarker", "com.intellij.reactivestreams", "com.intellij.httpClient", "com.intellij.spring", "com.intellij.swagger", "com.intellij.util.xml",
	"com.intellij.grpc", "com.intellij.guice", "com.intellij.helidon", "com.intellij.hibernate", "com.intellij.jdbi", "com.intellij.ktor", "com.intellij.micronaut", "com.intellij.openRewrite",
	"com.intellij.quarkus", "com.intellij.frameworks.thymeleaf", "com.intellij.wiremock", "com.intellij.ee",
}

func GenerateUltimateUnitTestsSettings(backendUrl string, client *http.Client) []detector.PerformanceSettings {
	settings := make([]detector.PerformanceSettings, 0, 1000)
	mainSettings := detector.PerformanceSettings{
		Db:      "perfUnitTests",
		Table:   "report",
		Branch:  "master",
		Machine: "intellij-linux-%-hetzner-%",
		Metric:  "attempt.mean.ms",
	}
	slackSettings := detector.SlackSettings{
		Channel:     "ij-u-team-performance-issues-check",
		ProductLink: "perfUnit",
	}

	tests, err := detector.FetchAllTests(backendUrl, client, mainSettings)
	if err != nil {
		slog.Error("error while getting tests", "error", err)
		return settings
	}

	tests = slices.DeleteFunc(tests, func(test string) bool {
		return !slices.ContainsFunc(ultimatePackages, func(pkg string) bool {
			return strings.HasPrefix(test, pkg)
		})
	})

	for _, test := range tests {
		settings = append(settings, detector.PerformanceSettings{
			Project:       test,
			Db:            mainSettings.Db,
			Table:         mainSettings.Table,
			Branch:        mainSettings.Branch,
			Machine:       mainSettings.Machine,
			Metric:        mainSettings.Metric,
			SlackSettings: slackSettings,
			AnalysisSettings: detector.AnalysisSettings{
				MinimumSegmentLength:      30,
				MedianDifferenceThreshold: 20,
				EffectSizeThreshold:       2,
			},
		})
	}
	return settings
}
