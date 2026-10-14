package setting

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHooksAreComplete(t *testing.T) {
	t.Parallel()
	names := map[string]bool{}
	envs := map[string]bool{}
	for _, hook := range Hooks() {
		assert.NotEmpty(t, hook.Name)
		assert.False(t, names[hook.Name], "hook %s is declared twice", hook.Name)
		names[hook.Name] = true
		assert.NotEmpty(t, hook.Db, hook.Name)
		assert.NotEmpty(t, hook.Table, hook.Name)
		// An empty prefix would send every degradation of the table.
		assert.NotEmpty(t, hook.ProjectPrefix, hook.Name)
		for _, env := range []string{hook.URLEnv, hook.TokenEnv} {
			assert.NotEmpty(t, env, hook.Name)
			assert.False(t, envs[env], "variable %s is used twice", env)
			envs[env] = true
		}
	}
}
