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

// TestGroupPredicate locks the WHERE-clause suffix rendered per group.
func TestGroupPredicate(t *testing.T) {
	t.Parallel()
	// single prefix
	assert.Equal(t, "like 'intellij-linux-hw-blade-%'", GroupPredicate("machine", "linux-blade"))
	// two prefixes chain as a disjunction in the suffix form
	assert.Equal(t,
		"like 'intellij-linux-hw-hetzner%' or machine like 'intellij-linux-agg-hw-hetzner-agent%'",
		GroupPredicate("machine", "linux-blade-hetzner"))
	// legacy fixed-name rule renders as IN
	assert.Equal(t,
		"in ('intellij-macos-hw-unit-1550', 'intellij-macos-hw-unit-1551', 'intellij-macos-hw-unit-1772', 'intellij-macos-hw-unit-1773')",
		GroupPredicate("machine", "macMini 2018"))
	// regex-defined and unmapped groups have no predicate — callers fall back to member lists
	assert.Empty(t, GroupPredicate("machine", "windows-azure"))
	assert.Empty(t, GroupPredicate("machine", unknownGroup))
	assert.Empty(t, GroupPredicate("machine", "no-such-group"))
}

// TestGroupPredicateSoundness enforces the invariant GroupPredicate depends on: a group's
// prefixes and names must be claimed by that group alone. If a prefix of one group subsumed a
// prefix or name of another, the per-group predicate would select agents that first-match-wins
// GroupName assigns elsewhere. It also keeps prefixes free of LIKE wildcards, which the
// predicate embeds into LIKE patterns verbatim.
func TestGroupPredicateSoundness(t *testing.T) {
	t.Parallel()
	for _, r := range rules {
		for _, p := range r.prefixes {
			assert.NotContainsf(t, p, "%", "prefix %q of group %q must not contain LIKE wildcards", p, r.group)
			assert.NotContainsf(t, p, "_", "prefix %q of group %q must not contain LIKE wildcards", p, r.group)
		}
	}
	for _, r1 := range rules {
		for _, r2 := range rules {
			if r1.group == r2.group {
				continue
			}
			for _, p1 := range r1.prefixes {
				for _, p2 := range r2.prefixes {
					assert.Falsef(t, strings.HasPrefix(p2, p1),
						"prefix %q of group %q subsumes prefix %q of group %q", p1, r1.group, p2, r2.group)
				}
				for _, n2 := range r2.names {
					assert.Falsef(t, strings.HasPrefix(n2, p1),
						"prefix %q of group %q subsumes name %q of group %q", p1, r1.group, n2, r2.group)
				}
			}
		}
	}
}
