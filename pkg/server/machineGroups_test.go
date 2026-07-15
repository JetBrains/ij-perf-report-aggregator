package server

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestMachineGroupName locks the Go port of getMachineGroupName against representative
// real agent names, including the powerful-vs-weak split that AT-4930 hinges on.
func TestMachineGroupName(t *testing.T) {
	t.Parallel()
	cases := map[string]string{
		"intellij-linux-performance-aws-i-08aec6c8ee5a71bba":        "Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)",
		"intellij-linux-performance-aws-lt-a-i-045667485579a157a":   "Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)",
		"intellij-linux-performance-tiny-aws-on-demand-i-0abc12345": "Linux EC2 C6id.xlarge (4 vCPU Xeon, 8 GB)",
		"intellij-windows-performance-aws-i-0deadbeef0011":          "Windows EC2 C6id.4xlarge or i4i.4xlarge (16 vCPU Xeon, 32 or 128 GB)",
		"intellij-windows-performance-mem-aws-i-0deadbeef0022":      "Windows EC2 C6id.4xlarge or i4i.4xlarge (16 vCPU Xeon, 32 or 128 GB)",
		"intellij-linux-hw-hetzner-agent-42":                        "linux-blade-hetzner",
		"intellij-macos-perf-eqx-143291":                            "Mac Mini M2 Pro (10 vCPU, 32 GB)",
		"ij-w11u-azr7":                                              "windows-azure",
	}
	for machine, want := range cases {
		assert.Equalf(t, want, machineGroupName(machine), "group for %q", machine)
	}

	// Powerful and weak linux classes must never collapse to the same group.
	assert.NotEqual(t,
		machineGroupName("intellij-linux-performance-aws-i-08aec6c8ee5a71bba"),
		machineGroupName("intellij-linux-performance-tiny-aws-on-demand-i-0abc12345"),
	)
}

// TestMachineGroupNameFallbackStem verifies unmapped names fall back to a stem (never one
// shared "Unknown" bucket), so distinct unmapped classes stay separate.
func TestMachineGroupNameFallbackStem(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "brand-new-class-aws", machineGroupName("brand-new-class-aws-i-00ff00ff00ff00ff0"))
	assert.NotEqual(t,
		machineGroupName("brand-new-class-a-aws-i-00ff00ff00ff00ff0"),
		machineGroupName("brand-new-class-b-aws-i-11ff11ff11ff11ff1"),
	)
}

// TestMachineGroupSQLExpr sanity-checks the generated ClickHouse expression.
func TestMachineGroupSQLExpr(t *testing.T) {
	t.Parallel()
	expr := machineGroupSQLExpr("machine")
	assert.True(t, strings.HasPrefix(expr, "multiIf("))
	assert.Contains(t, expr, "startsWith(machine, 'intellij-linux-performance-aws-i-')")
	assert.Contains(t, expr, "match(machine, 'ij-w.*-azr.*')")
	assert.Contains(t, expr, "replaceRegexpOne(") // stem fallback
}
