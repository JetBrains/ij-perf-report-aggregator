package setting

import (
	"testing"

	detector "github.com/JetBrains/ij-perf-report-aggregator/pkg/degradation-detector"
	"github.com/stretchr/testify/assert"
)

func TestTestMatchesPackages(t *testing.T) {
	t.Parallel()
	packages := []string{"com.intellij.java", "org.jetbrains.plugins.ruby"}

	assert.True(t, testMatchesPackages("com.intellij.java.Foo", packages))
	assert.True(t, testMatchesPackages("org.jetbrains.plugins.ruby.Bar", packages))
	// The prefix must end on a dot boundary: a longer package that merely starts with
	// the same characters must not match.
	assert.False(t, testMatchesPackages("com.intellij.javaScript.Baz", packages))
	// The bare package name (no trailing class) is not a member of the package.
	assert.False(t, testMatchesPackages("com.intellij.java", packages))
	assert.False(t, testMatchesPackages("com.example.Other", packages))
}

func TestResolveRoute(t *testing.T) {
	t.Parallel()
	// Synthetic routing table: one plain team, one with custom (degradation-only) analysis, and
	// one with additional metrics — enough to exercise every branch without depending on the real
	// team/channel/package mappings. Passing it explicitly to resolveRoute keeps the test isolated
	// from the production mapping data (and free of shared global state, so it can run in parallel).
	const deltaMention = "<!subteam^SDELTA>"
	configs := []teamConfig{
		{Team: "alpha", SlackChannel: "alpha-channel", Packages: []string{"com.test.alpha"}},
		{Team: "beta", SlackChannel: "beta-channel", Packages: []string{"com.test.beta"}, AnalysisSettings: degradationOnlyAnalysisSettings()},
		{Team: "gamma", SlackChannel: "gamma-channel", Packages: []string{"com.test.gamma"}, AdditionalTestMetrics: map[string][]string{"com.test.gamma.PacketsTest": {"%.expected.%"}}},
		{Team: "delta", SlackChannel: "delta-channel", Packages: []string{"com.test.delta"}, Mention: deltaMention},
	}

	const (
		plainTest     = "com.test.alpha.SomeTest"
		degradedTest  = "com.test.beta.SomeTest"
		metricsTest   = "com.test.gamma.PacketsTest"
		mentionedTest = "com.test.delta.SomeTest"
		unknownTest   = "com.test.unknown.SomeTest"
	)

	tests := []struct {
		name           string
		test           string
		projectOwners  map[string]string
		ownerChannels  map[string]string
		wantChannel    string
		wantAnalysis   detector.AnalysisSettings
		wantHasMetrics bool
		wantMention    string
	}{
		{
			name:         "no package match and no owner falls back to the catch-all channel",
			test:         unknownTest,
			wantChannel:  catchAllUnitTestsChannel,
			wantAnalysis: defaultUnitTestAnalysisSettings,
		},
		{
			name:         "package match routes to the team channel with default analysis",
			test:         plainTest,
			wantChannel:  "alpha-channel",
			wantAnalysis: defaultUnitTestAnalysisSettings,
		},
		{
			name:         "package with custom analysis keeps its analysis settings",
			test:         degradedTest,
			wantChannel:  "beta-channel",
			wantAnalysis: *degradationOnlyAnalysisSettings(),
		},
		{
			name:          "owner channel overrides the package channel but analysis stays from the package",
			test:          degradedTest,
			projectOwners: map[string]string{degradedTest: "beta-owner"},
			ownerChannels: map[string]string{"beta-owner": "owner-specific-channel"},
			wantChannel:   "owner-specific-channel",
			wantAnalysis:  *degradationOnlyAnalysisSettings(),
		},
		{
			name:          "owner channel applies to a test that matches no package",
			test:          unknownTest,
			projectOwners: map[string]string{unknownTest: "some-owner"},
			ownerChannels: map[string]string{"some-owner": "owner-channel"},
			wantChannel:   "owner-channel",
			wantAnalysis:  defaultUnitTestAnalysisSettings,
		},
		{
			name:          "owner without a channel mapping keeps the package channel",
			test:          plainTest,
			projectOwners: map[string]string{plainTest: "unmapped-owner"},
			ownerChannels: map[string]string{},
			wantChannel:   "alpha-channel",
			wantAnalysis:  defaultUnitTestAnalysisSettings,
		},
		{
			name:          "empty owner string is ignored and keeps the package channel",
			test:          plainTest,
			projectOwners: map[string]string{plainTest: ""},
			ownerChannels: map[string]string{"": "should-not-be-used"},
			wantChannel:   "alpha-channel",
			wantAnalysis:  defaultUnitTestAnalysisSettings,
		},
		{
			name:           "additional test metrics are preserved even when the owner overrides the channel",
			test:           metricsTest,
			projectOwners:  map[string]string{metricsTest: "gamma-owner"},
			ownerChannels:  map[string]string{"gamma-owner": "gamma-owner-channel"},
			wantChannel:    "gamma-owner-channel",
			wantAnalysis:   defaultUnitTestAnalysisSettings,
			wantHasMetrics: true,
		},
		{
			name:         "mention is kept when the alert lands in the team's own channel via the package fallback",
			test:         mentionedTest,
			wantChannel:  "delta-channel",
			wantAnalysis: defaultUnitTestAnalysisSettings,
			wantMention:  deltaMention,
		},
		{
			name:          "mention is kept when owner routing resolves to the team's own channel",
			test:          mentionedTest,
			projectOwners: map[string]string{mentionedTest: "delta-owner"},
			ownerChannels: map[string]string{"delta-owner": "delta-channel"},
			wantChannel:   "delta-channel",
			wantAnalysis:  defaultUnitTestAnalysisSettings,
			wantMention:   deltaMention,
		},
		{
			name:          "mention is dropped when owner routing redirects the alert to another channel",
			test:          mentionedTest,
			projectOwners: map[string]string{mentionedTest: "delta-owner"},
			ownerChannels: map[string]string{"delta-owner": "some-other-channel"},
			wantChannel:   "some-other-channel",
			wantAnalysis:  defaultUnitTestAnalysisSettings,
			wantMention:   "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			route := resolveRoute(tc.test, tc.projectOwners, tc.ownerChannels, configs)
			assert.Equal(t, tc.wantChannel, route.channel)
			assert.Equal(t, tc.wantAnalysis, route.analysisSettings)
			assert.Equal(t, tc.wantMention, route.mention)
			if tc.wantHasMetrics {
				assert.NotEmpty(t, route.additionalTestMetrics)
			} else {
				assert.Empty(t, route.additionalTestMetrics)
			}
		})
	}
}
