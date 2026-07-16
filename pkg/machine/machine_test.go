package machine

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestGroupName locks the mapping against representative real agent names, including the
// powerful-vs-weak split that AT-4930 hinges on.
func TestGroupName(t *testing.T) {
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
		// legacy fixed-name agents (migrated from valueToGroup)
		"intellij-macos-hw-unit-1550": "macMini 2018",
		"intellij-linux-hw-unit-531":  "Linux Space i7-3770, 16Gb",
	}
	for name, want := range cases {
		assert.Equalf(t, want, GroupName(name), "group for %q", name)
	}

	// Powerful and weak linux classes must never collapse to the same group.
	assert.NotEqual(t,
		GroupName("intellij-linux-performance-aws-i-08aec6c8ee5a71bba"),
		GroupName("intellij-linux-performance-tiny-aws-on-demand-i-0abc12345"),
	)
}

// TestGroupNameFallbackUnknown verifies unmapped names fall back to "Unknown".
func TestGroupNameFallbackUnknown(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "Unknown", GroupName("brand-new-class-aws-i-00ff00ff00ff00ff0"))
	assert.Equal(t, "Unknown", GroupName(""))
}

// TestGroupSQLExpr sanity-checks the generated ClickHouse expression.
func TestGroupSQLExpr(t *testing.T) {
	t.Parallel()
	expr := GroupSQLExpr("machine")
	assert.True(t, strings.HasPrefix(expr, "multiIf("))
	assert.Contains(t, expr, "startsWith(machine, 'intellij-linux-performance-aws-i-')")
	assert.Contains(t, expr, "match(machine, 'ij-w.*-azr.*')")
	assert.Contains(t, expr, "machine IN ('intellij-macos-hw-unit-1550'") // legacy exact-name rule
	assert.True(t, strings.HasSuffix(expr, ", 'Unknown')"))               // fallback
}
