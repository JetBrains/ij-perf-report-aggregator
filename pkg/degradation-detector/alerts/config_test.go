package alerts

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseBuiltInConfig(t *testing.T) {
	t.Parallel()
	config, err := Parse(snapshot)
	require.NoError(t, err)
	assert.NotEmpty(t, config.Version)
	for _, id := range []string{"metric-degradation", "metric-threshold", "unit-test-degradation"} {
		assert.NotEmpty(t, config.Notifications[id].Subscribers, id)
	}
	assert.NotEmpty(t, config.Notifications["unit-test-degradation"].Fallback)
}

func TestParseRejectsWhatItCannotRead(t *testing.T) {
	t.Parallel()
	tests := map[string]string{
		"not json":         `{`,
		"newer format":     `{"formatVersion": 2, "service": "ij-perf", "notifications": {"a": {}}}`,
		"another service":  `{"formatVersion": 1, "service": "socrat", "notifications": {"a": {}}}`,
		"no notifications": `{"formatVersion": 1, "service": "ij-perf"}`,
	}
	for name, data := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			_, err := Parse([]byte(data))
			assert.Error(t, err)
		})
	}
}

func TestParseIgnoresUnknownFields(t *testing.T) {
	t.Parallel()
	config, err := Parse([]byte(`{"formatVersion": 1, "service": "ij-perf", "version": "v1", "future": 1,
		"notifications": {"n": {"future": true, "subscribers": {"s": {"enabled": true, "target": {"channel": "#c", "mention": ["<@U1>"]}}}}}}`))
	require.NoError(t, err)
	target := config.Notifications["n"].Subscribers["s"].Target
	assert.Equal(t, "c", target.SlackChannel())
	assert.Equal(t, []string{"<@U1>"}, target.Mention)
}

func TestMatches(t *testing.T) {
	t.Parallel()
	tests := []struct {
		pattern string
		value   string
		matches bool
	}{
		{"idea/measureStartup", "idea/measureStartup", true},
		{"idea/measureStartup", "idea/measureStartup2", false},
		{"*", "anything", true},
		{"*", "", true},
		{"kotlin/*", "kotlin/completion", true},
		{"kotlin/*", "kotlin_coroutines/completion", false},
		{"*/measureStartup", "idea/measureStartup", true},
		{"*/measureStartup", "idea/measureStartup/extra", false},
		{"*_k2", "arrow/completion_k2", true},
		{"intellij-linux-*-aws-*", "intellij-linux-performance-aws-12", true},
		{"intellij-linux-*-aws-*", "intellij-linux-aws", false},
		{"a*b*a", "abxa", true},
		{"a*b*a", "axa", false},
		// A character of the value is not used twice.
		{"a*a", "a", false},
		{"a*a", "aa", true},
		// Only `*` is special.
		{"com.intellij.java.*", "com.intellij.java.Test", true},
		{"com.intellij.java.*", "com.intellij.javascript.Test", false},
		{"a?c", "abc", false},
		{"a_c", "abc", false},
	}
	for _, test := range tests {
		assert.Equal(t, test.matches, Matches([]string{"no match", test.pattern}, test.value), "%s ~ %s", test.pattern, test.value)
	}
	assert.False(t, Matches(nil, "value"))
}
