package setting

import (
	"errors"
	"testing"

	detector "github.com/JetBrains/ij-perf-report-aggregator/pkg/degradation-detector"
	"github.com/JetBrains/ij-perf-report-aggregator/pkg/degradation-detector/alerts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// fakeBackend answers from maps. The key of tests is "table/branch/machine", of metrics the name of the test.
type fakeBackend struct {
	tests   map[string][]string
	metrics map[string][]string
	owners  map[string]string
	groups  []detector.CodeOwnerGroup
}

func (b fakeBackend) Tests(_, table, branch, machine string) ([]string, error) {
	tests, ok := b.tests[table+"/"+branch+"/"+machine]
	if !ok {
		return nil, errors.New("no data")
	}
	return tests, nil
}

func (b fakeBackend) MetricNames(test detector.PerformanceSettings, pattern string) ([]string, error) {
	var result []string
	for _, metric := range b.metrics[test.Project] {
		if alerts.Matches([]string{toGlob(pattern)}, metric) {
			result = append(result, metric)
		}
	}
	return result, nil
}

func toGlob(likePattern string) string {
	glob := []rune(likePattern)
	for i, r := range glob {
		if r == '%' {
			glob[i] = '*'
		}
	}
	return string(glob)
}

func (b fakeBackend) ProjectOwners(string, string) (map[string]string, error) {
	if b.owners == nil {
		return nil, errors.New("no owners")
	}
	return b.owners, nil
}

func (b fakeBackend) OwnerGroups() ([]detector.CodeOwnerGroup, error) {
	if b.groups == nil {
		return nil, errors.New("no groups")
	}
	return b.groups, nil
}

func parseConfig(t *testing.T, notifications string) *alerts.Config {
	t.Helper()
	config, err := alerts.Parse([]byte(`{"formatVersion": 1, "service": "ij-perf", "version": "test", "notifications": ` + notifications + `}`))
	require.NoError(t, err)
	return config
}

func performanceSettings(t *testing.T, settings []detector.Settings) []detector.PerformanceSettings {
	t.Helper()
	result := make([]detector.PerformanceSettings, 0, len(settings))
	for _, s := range settings {
		p, ok := s.(detector.PerformanceSettings)
		require.True(t, ok)
		result = append(result, p)
	}
	return result
}

func generate(t *testing.T, config *alerts.Config, backend Backend, generator string) []detector.PerformanceSettings {
	t.Helper()
	for _, g := range Generators(config, backend) {
		if g.Name == generator {
			return performanceSettings(t, g.Generate())
		}
	}
	require.Fail(t, "no generator "+generator)
	return nil
}

func TestMetricDegradationBlock(t *testing.T) {
	t.Parallel()
	config := parseConfig(t, `{"metric-degradation": {"subscribers": {"team": {"enabled": true, "match": [{
		"fields": {
			"db": ["perfintDev"], "table": ["idea"], "branch": ["master", "253"],
			"machine": ["intellij-linux-performance-aws-*"], "mode": ["default", "split"],
			"project": ["kotlin/*", "literal/test"], "metric": ["completion", "gc*"]
		},
		"target": {"channel": "#team-alerts", "mention": ["<@U1>", "<@U2>"]},
		"spec": {"reportType": "degradation", "minimumSegmentLength": 8, "medianDifferenceThreshold": 10, "effectSizeThreshold": 2,
			"missingDataDays": 5, "metricRules": "common", "medianDifferenceThresholdByMetric": {"gcPause": 20},
			"metricAliases": {"completion": "typing"}}
	}]}}}}`)
	backend := fakeBackend{
		// The tests of a pattern are the ones that ran on the branch, so the branch 253 without runs has only the literal test.
		tests:   map[string][]string{"idea/master/intellij-linux-performance-aws-%": {"kotlin/completion", "kotlin_coroutines/completion"}},
		metrics: map[string][]string{"kotlin/completion": {"gcPause", "gcCount", "other"}, "literal/test": {"gcPause"}},
	}

	settings := generate(t, config, backend, "metric-degradation/team")

	type series struct{ branch, mode, test, metric string }
	got := map[series]detector.PerformanceSettings{}
	for _, s := range settings {
		got[series{s.Branch, s.Mode, s.Project, s.Metric}] = s
	}
	want := make([]series, 0, 14)
	for _, mode := range []string{"", "split"} {
		for _, metric := range []string{"completion", "gcPause", "gcCount"} {
			want = append(want, series{"master", mode, "kotlin/completion", metric})
		}
		for _, branch := range []string{"master", "253"} {
			want = append(want, series{branch, mode, "literal/test", "completion"}, series{branch, mode, "literal/test", "gcPause"})
		}
	}
	require.Len(t, settings, len(want))
	for _, w := range want {
		assert.Contains(t, got, w)
	}

	s := got[series{"master", "split", "kotlin/completion", "gcPause"}]
	assert.Equal(t, "perfintDev", s.Db)
	assert.Equal(t, "idea", s.Table)
	assert.Equal(t, "intellij-linux-performance-aws-%", s.Machine)
	assert.Equal(t, "team-alerts", s.Channel)
	assert.Equal(t, "<@U1> <@U2>", s.Mention)
	assert.Equal(t, "intellij", s.ProductLink)
	assert.Equal(t, detector.DegradationEvent, s.ReportType)
	assert.Equal(t, 8, s.MinimumSegmentLength)
	assert.InDelta(t, 20, s.MedianDifferenceThreshold, 0)
	assert.Equal(t, -5, s.GetDaysToCheckMissing())
	assert.Equal(t, detector.ChangePointAnalysis, s.AnalysisKind)

	s = got[series{"master", "", "kotlin/completion", "completion"}]
	assert.InDelta(t, 10, s.MedianDifferenceThreshold, 0)
	assert.Equal(t, "typing", s.MetricAlias)
}

func TestMetricsComeFromRulesWhenBlockListsNone(t *testing.T) {
	t.Parallel()
	config := parseConfig(t, `{"metric-degradation": {"subscribers": {"team": {"enabled": true, "match": [{
		"fields": {"db": ["perfintDev"], "table": ["rust"], "branch": ["master"], "machine": ["m*"], "project": ["*"]},
		"target": {"channel": "#c"}, "spec": {"metricRules": "rust"}
	}]}}}}`)
	backend := fakeBackend{tests: map[string][]string{"rust/master/m%": {"cargo/completion/a", "cargo/unknown"}}}

	settings := generate(t, config, backend, "metric-degradation/team")

	require.Len(t, settings, 1)
	assert.Equal(t, "cargo/completion/a", settings[0].Project)
	assert.Equal(t, "completion#mean_value", settings[0].Metric)
	assert.Empty(t, settings[0].Mode)
	assert.Equal(t, "rust", settings[0].ProductLink)
}

func TestFirstBlockOfSubscriberWins(t *testing.T) {
	t.Parallel()
	block := func(project, channel string, onlyAnalyze string) string {
		return `{"fields": {"db": ["d"], "table": ["idea"], "branch": ["master"], "machine": ["m"], "project": ["` + project + `"], "metric": ["x"]},
			"target": {"channel": "` + channel + `"}, "spec": {"onlyAnalyze": ` + onlyAnalyze + `}}`
	}
	config := parseConfig(t, `{"metric-degradation": {"subscribers": {"team": {"enabled": true, "match": [`+
		block("a/quiet", "#first", "true")+`,`+block("a/*", "#second", "false")+`]}}}}`)
	backend := fakeBackend{tests: map[string][]string{"idea/master/m": {"a/quiet", "a/loud"}}}

	settings := generate(t, config, backend, "metric-degradation/team")

	channels := map[string]string{}
	for _, s := range settings {
		channels[s.Project] = s.Channel
	}
	// onlyAnalyze keeps the analysis and sends no message.
	assert.Equal(t, map[string]string{"a/quiet": "", "a/loud": "second"}, channels)
}

func TestMetricThreshold(t *testing.T) {
	t.Parallel()
	config := parseConfig(t, `{"metric-threshold": {"subscribers": {"team": {"enabled": true, "match": [{
		"fields": {"db": ["perfintDev"], "table": ["ml"], "branch": ["master"], "machine": ["m"], "project": ["prod: tokenQuota"], "metric": ["currentPercents"]},
		"target": {"channel": "#quota"}, "spec": {"comparison": "lessThan", "value": 95, "missingDataDays": 3}
	}]}}}}`)

	settings := generate(t, config, fakeBackend{}, "metric-threshold/team")

	require.Len(t, settings, 1)
	assert.Equal(t, detector.ThresholdAnalysis, settings[0].AnalysisKind)
	assert.Equal(t, detector.ThresholdLessThan, settings[0].ThresholdMode)
	assert.InDelta(t, 95, settings[0].ThresholdValue, 0)
	assert.Equal(t, "ml/dev", settings[0].ProductLink)
}

func TestMissingNotificationIsReportedAsEmpty(t *testing.T) {
	t.Parallel()
	config := parseConfig(t, `{"renamed": {}}`)
	for _, g := range Generators(config, fakeBackend{}) {
		if g.Name != "air" {
			assert.Empty(t, g.Generate(), g.Name)
		}
	}
	assert.Len(t, Generators(config, fakeBackend{}), 4)
}

func TestCheckUnitTestOwners(t *testing.T) {
	t.Parallel()
	require.Error(t, CheckUnitTestOwners(fakeBackend{owners: map[string]string{}}))
	require.Error(t, CheckUnitTestOwners(fakeBackend{groups: []detector.CodeOwnerGroup{}}))
	require.NoError(t, CheckUnitTestOwners(fakeBackend{owners: map[string]string{}, groups: []detector.CodeOwnerGroup{}}))
}

func TestDisabledSubscriberHasNoGenerator(t *testing.T) {
	t.Parallel()
	config := parseConfig(t, `{"metric-degradation": {"subscribers": {"off": {"enabled": false}, "on": {"enabled": true}}}}`)
	generators := Generators(config, fakeBackend{})
	names := make([]string, 0, len(generators))
	for _, g := range generators {
		names = append(names, g.Name)
	}
	// metric-threshold is not in this config, so it has a generator that reports it as empty.
	assert.Equal(t, []string{"air", "metric-degradation/on", "metric-threshold", "unit-test-degradation"}, names)
}

// The config that is built into the image must be readable by this code: a typo in a spec or new metric rules would
// otherwise show up only as a subscriber without alerts.
func TestBuiltInConfigIsUnderstood(t *testing.T) {
	t.Setenv(alerts.GithubTokenEnv, "")
	t.Setenv(alerts.ConfigFileEnv, "")
	config, _ := alerts.Load(t.Context(), nil)
	require.NotNil(t, config)
	for id, notification := range config.Notifications {
		_, err := parseSpec(notification.Spec)
		require.NoError(t, err, id)
		for name, subscriber := range notification.Subscribers {
			_, err := parseSpec(subscriber.Spec)
			require.NoError(t, err, "%s/%s", id, name)
			for _, block := range subscriber.Match {
				s, err := parseSpec(block.Spec)
				require.NoError(t, err, "%s/%s", id, name)
				if id == metricDegradation && len(block.Fields["metric"]) == 0 {
					assert.Contains(t, metricRules, s.MetricRules, "%s/%s", id, name)
				}
			}
		}
	}
	_, err := newUnitTestRouter(config.Notifications[unitTestDegradation], nil)
	require.NoError(t, err)
}
