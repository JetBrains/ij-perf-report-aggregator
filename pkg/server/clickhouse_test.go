package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestMachineCondition locks the WHERE condition rendered for a machine selection mixing
// hardware-class group names and raw agent names.
func TestMachineCondition(t *testing.T) {
	t.Parallel()
	// group names filter by the computed machine_group
	assert.Equal(t, "machine_group IN ('linux-blade')", machineCondition([]string{"linux-blade"}))
	assert.Equal(t, "machine_group IN ('linux-blade', 'linux-blade-hetzner')", machineCondition([]string{"linux-blade", "linux-blade-hetzner"}))
	// the Unknown bucket is a group too
	assert.Equal(t, "machine_group IN ('Unknown')", machineCondition([]string{"Unknown"}))
	// raw agent names filter by exact machine
	assert.Equal(t, "machine IN ('intellij-linux-hw-blade-023')", machineCondition([]string{"intellij-linux-hw-blade-023"}))
	// mixed selections combine as a disjunction
	assert.Equal(t,
		"machine_group IN ('linux-blade') OR machine IN ('intellij-linux-hw-hetzner-agent-42')",
		machineCondition([]string{"linux-blade", "intellij-linux-hw-hetzner-agent-42"}))
	// an empty selection matches everything
	assert.Equal(t, "1", machineCondition(nil))
	// values are escaped before being embedded into SQL
	assert.Equal(t, "machine IN ('it''s')", machineCondition([]string{"it's"}))
}
