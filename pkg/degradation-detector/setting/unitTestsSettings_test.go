package setting

import (
	"testing"

	detector "github.com/JetBrains/ij-perf-report-aggregator/pkg/degradation-detector"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const unitTestNotification = `{"unit-test-degradation": {
	"fallback": "#fallback",
	"spec": {"reportType": "all", "minimumSegmentLength": 30, "medianDifferenceThreshold": 20, "missingDataDays": 3},
	"subscribers": {
		"java": {
			"enabled": true, "target": {"channel": "#java"}, "spec": {"minimumSegmentLength": 30},
			"owners": [{"name": "Java", "subgroups": true}],
			"match": [{"fields": {"test": ["com.intellij.java.*"]}, "target": {"channel": "#java-tests", "mention": ["<@U1>"]},
				"spec": {"reportType": "degradation", "minimumSegmentLength": 30}}]
		},
		"java-twin": {
			"enabled": true, "target": {"channel": "#java"}, "spec": {"minimumSegmentLength": 30},
			"owners": [{"name": "Java", "subgroups": true}], "match": []
		},
		"rust": {
			"enabled": true, "target": {"channel": "#rust"}, "spec": {"minimumSegmentLength": 30},
			"owners": [{"name": "Rust", "subgroups": false}], "match": []
		},
		"debugger": {
			"enabled": true, "target": {"channel": "#debugger"},
			"spec": {"minimumSegmentLength": 30, "additionalMetrics": {"com.intellij.debugger.PacketsTest": ["*.expected.*"]}},
			"owners": [], "match": [{"fields": {"test": ["com.intellij.debugger.*"]}, "target": {"channel": "#debugger"},
				"spec": {"minimumSegmentLength": 30, "additionalMetrics": {"com.intellij.debugger.PacketsTest": ["*.expected.*"]}}}]
		},
		"muted": {
			"enabled": false, "target": {"channel": "#todo"}, "spec": {"minimumSegmentLength": 30},
			"owners": [{"name": "Muted", "subgroups": true}], "match": []
		},
		"tickets": {
			"enabled": true, "target": {"youtrack": {"project": "IJPL"}}, "spec": {"minimumSegmentLength": 30},
			"owners": [{"name": "Tickets", "subgroups": true}], "match": []
		}
	}
}}`

var unitTestGroups = []detector.CodeOwnerGroup{
	{Name: "Java"},
	{Name: "Java Debugger", Parent: "Java"},
	{Name: "Java Stepping", Parent: "Java Debugger", Aliases: []string{"Stepping"}},
	{Name: "Rust", Aliases: []string{"Rust Engineering"}},
	{Name: "Rust Toml", Parent: "Rust"},
	{Name: "Muted"},
	{Name: "Tickets"},
}

func TestUnitTestRoutes(t *testing.T) {
	t.Parallel()
	router, err := newUnitTestRouter(parseConfig(t, unitTestNotification).Notifications[unitTestDegradation], unitTestGroups)
	require.NoError(t, err)

	tests := []struct {
		name     string
		test     string
		owner    string
		channels []string
	}{
		{"nobody selects the test", "org.other.Test", "", []string{"fallback"}},
		{"an unknown owner", "org.other.Test", "Unknown", []string{"fallback"}},
		{"the claimed owner", "org.other.Test", "Java", []string{"java"}},
		{"a subgroup of the claimed owner", "org.other.Test", "Java Debugger", []string{"java"}},
		{"a subgroup of a subgroup under its old name", "org.other.Test", "Stepping", []string{"java"}},
		{"the claimed owner under its old name", "org.other.Test", "Rust Engineering", []string{"rust"}},
		{"a subgroup of an owner claimed without subgroups", "org.other.Test", "Rust Toml", []string{"fallback"}},
		{"a block wins over the owner claim of its subscriber", "com.intellij.java.Test", "Java", []string{"java-tests", "java"}},
		{"a block and the claim of another subscriber", "com.intellij.debugger.Test", "Rust", []string{"debugger", "rust"}},
		{"a disabled subscriber keeps its test from the fallback", "org.other.Test", "Muted", []string{""}},
		{"a YouTrack target gets no message", "org.other.Test", "Tickets", []string{""}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			var channels []string
			for _, route := range router.resolve(test.test, test.owner) {
				channels = append(channels, route.slack.Channel)
			}
			assert.Equal(t, test.channels, channels)
		})
	}

	routes := router.resolve("com.intellij.java.Test", "")
	require.Len(t, routes, 1)
	assert.Equal(t, "<@U1>", routes[0].slack.Mention)
	assert.Equal(t, "degradation", routes[0].spec.ReportType)
	assert.Equal(t, "perfUnit", routes[0].slack.ProductLink)
}

func TestUnitTestSettings(t *testing.T) {
	t.Parallel()
	config := parseConfig(t, unitTestNotification)
	backend := fakeBackend{
		tests: map[string][]string{
			"report/master/" + bladeUnitTestsMachine:   {"com.intellij.debugger.PacketsTest.test", "org.other.Test"},
			"report/master/" + hetznerUnitTestsMachine: {"org.other.Test", "org.other.HetznerTest"},
		},
		metrics: map[string][]string{"com.intellij.debugger.PacketsTest.test": {"attempt.mean.ms", "a.expected.b", "c.actual.d"}},
		owners:  map[string]string{"org.other.Test": "Java Debugger"},
		groups:  unitTestGroups,
	}

	settings := generate(t, config, backend, unitTestDegradation)

	type series struct{ machine, test, metric, channel string }
	got := make([]series, 0, len(settings))
	for _, s := range settings {
		got = append(got, series{s.Machine, s.Project, s.Metric, s.Channel})
		assert.Equal(t, "perfUnitTests", s.Db)
		assert.Equal(t, "report", s.Table)
		assert.Equal(t, "master", s.Branch)
		assert.Equal(t, 30, s.MinimumSegmentLength)
		assert.Equal(t, s.Machine == bladeUnitTestsMachine, !s.IgnoreChangesBefore.IsZero(), "only the Blade series ignore the old changes")
	}
	assert.ElementsMatch(t, []series{
		{bladeUnitTestsMachine, "com.intellij.debugger.PacketsTest.test", "attempt.mean.ms", "debugger"},
		{bladeUnitTestsMachine, "com.intellij.debugger.PacketsTest.test", "a.expected.b", "debugger"},
		{bladeUnitTestsMachine, "org.other.Test", "attempt.mean.ms", "java"},
		{hetznerUnitTestsMachine, "org.other.HetznerTest", "attempt.mean.ms", "fallback"},
	}, got)
}

func TestUnitTestsAreRoutedByNameWithoutOwners(t *testing.T) {
	t.Parallel()
	config := parseConfig(t, unitTestNotification)
	backend := fakeBackend{tests: map[string][]string{
		"report/master/" + bladeUnitTestsMachine:   {"com.intellij.java.Test", "org.other.Test"},
		"report/master/" + hetznerUnitTestsMachine: {},
	}}

	settings := generate(t, config, backend, unitTestDegradation)

	channels := map[string]string{}
	for _, s := range settings {
		channels[s.Project] = s.Channel
	}
	assert.Equal(t, map[string]string{"com.intellij.java.Test": "java-tests", "org.other.Test": "fallback"}, channels)
}

func TestSplitUnitTestsByMachine(t *testing.T) {
	t.Parallel()
	series := splitUnitTestsByMachine(
		[]string{"com.test.Both", "com.test.BladeOnly"},
		[]string{"com.test.Both", "com.test.HetznerOnly"},
	)

	require.Len(t, series, 2)
	assert.Equal(t, bladeUnitTestsMachine, series[0].machine)
	assert.Equal(t, []string{"com.test.Both", "com.test.BladeOnly"}, series[0].tests)
	assert.Equal(t, unitTestsOnBladeSince, series[0].ignoreChangesBefore)
	// A test that runs on both pools is analyzed on Blade only.
	assert.Equal(t, hetznerUnitTestsMachine, series[1].machine)
	assert.Equal(t, []string{"com.test.HetznerOnly"}, series[1].tests)
	assert.True(t, series[1].ignoreChangesBefore.IsZero())
}
