package setting

import (
	"os"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHooksAreComplete(t *testing.T) {
	t.Parallel()
	values, err := os.ReadFile("../../../deployment/report-aggregator/values.yaml")
	require.NoError(t, err)
	names := map[string]bool{}
	envs := map[string]bool{}
	for _, hook := range Hooks() {
		assert.NotEmpty(t, hook.Name)
		assert.False(t, names[hook.Name], "hook %s is declared twice", hook.Name)
		names[hook.Name] = true
		assert.NotEmpty(t, hook.Channel, "hook %s selects by the Slack channel of its team", hook.Name)
		assert.NotEmpty(t, hook.Secret, hook.Name)
		assert.Regexp(t, `(?m)^\s+- `+regexp.QuoteMeta(hook.Secret)+`\s*$`, string(values),
			"the secret of hook %s is not in degradationDetector.hookSecrets", hook.Name)
		for _, env := range []string{hook.URLEnv, hook.TokenEnv} {
			assert.NotEmpty(t, env, hook.Name)
			assert.False(t, envs[env], "variable %s is used twice", env)
			envs[env] = true
		}
	}
}
