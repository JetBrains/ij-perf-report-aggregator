// Package machine maps a raw perf-agent name to its hardware-class group.
//
// A raw `machine` value is an ephemeral agent name (usually "<class>-i-<aws-instance-id>").
// Comparing branches requires grouping runs by hardware *class*, never by the raw instance,
// otherwise runs on different agents get pooled into one median and produce phantom
// regressions (AT-4930) — e.g. the broad "%linux%" filter matches both the powerful
// "intellij-linux-performance-aws-*" and the weak "intellij-linux-performance-tiny-aws-*".
//
// This is the single source of truth for the grouping: the backend groups agents (for the
// comparison query and for the machine-selector endpoint) and returns them already grouped,
// so the frontend holds no grouping logic. First matching rule wins; unmapped names fall
// back to "Unknown".
package machine

import (
	"regexp"
	"slices"
	"strings"

	"github.com/JetBrains/ij-perf-report-aggregator/pkg/sql-util"
)

const unknownGroup = "Unknown"

type rule struct {
	prefixes []string // match if the name starts with any of these
	names    []string // match if the name equals any of these (legacy fixed-name agents)
	regex    *regexp.Regexp
	group    string
}

// Order matters: first matching rule wins. Prefix/regex rules come first, then the legacy
// fixed-name agents (kept in sync with what used to be MachineConfigurator's valueToGroup —
// only entries not already covered by a prefix rule are listed here).
//
// The frontend queries a selected group by GroupPredicate — an OR of per-prefix LIKEs rendered
// from the rule itself (deliberate: ClickHouse is much slower on `machine IN (...)` with
// hundreds of agents). Keep prefixes selective, and keep them non-overlapping across groups —
// TestGroupPredicateSoundness enforces that no rule's prefix subsumes another group's prefix
// or name, which is what makes the per-group predicate equivalent to first-match-wins
// GroupName.
var rules = []rule{
	{prefixes: []string{"intellij-linux-hw-blade-"}, group: "linux-blade"},
	{prefixes: []string{"ij-linux-x64-perf-hw-blade-"}, group: "linux-unit-perf-blade"},
	{prefixes: []string{"intellij-linux-test-hw-blade-"}, group: "linux-blade-test"},
	{prefixes: []string{"intellij-windows-hw-blade-"}, group: "windows-blade"},
	{prefixes: []string{"intellij-windows-hw-munit-"}, group: "Windows Munich i7-3770, 32 Gb"},
	{prefixes: []string{
		"intellij-linux-aws-amd-lt", "intellij-linux-aws-amd-2-lt",
		"intellij-linux-aws-3-lt", "intellij-linux-aws-lt",
	}, group: "Linux C5ad.xlarge or M5ad.xlarge or M5d.xlarge or C5d.xlarge"},
	{prefixes: []string{"intellij-macos-unit-2200-large-"}, group: "mac large"},
	{prefixes: []string{"intellij-linux-performance-aws-i-", "intellij-linux-performance-aws-lt"}, group: "Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)"},
	{prefixes: []string{"intellij-linux-performance-tiny-aws-i-", "intellij-linux-performance-tiny-aws-on-demand-i-"}, group: "Linux EC2 C6id.xlarge (4 vCPU Xeon, 8 GB)"},
	{prefixes: []string{"default-linux-aws-large-disk-"}, group: "Linux EC2 M5ad.2xlarge (8 vCPU Xeon, 32 GB)"},
	{prefixes: []string{"intellij-windows-performance-aws-i-", "intellij-windows-performance-mem-aws-i"}, group: "Windows EC2 C6id.4xlarge or i4i.4xlarge (16 vCPU Xeon, 32 or 128 GB)"},
	// Same 4 vCPU / 8 GB hardware as "Linux EC2 c5.xlarge (4 vCPU, 8 GB)" (metric values within
	// ~1% for the same test), but kept as its own group: merged, the members' common prefix
	// would collapse to "qodana-" and the group's LIKE query would match every qodana class.
	{prefixes: []string{"qodana-fleet-linux-amd64-heavy"}, group: "Linux EC2 c5.xlarge fleet (4 vCPU, 8 GB)"},
	{prefixes: []string{
		"intellij-linux-2004-aws-m5d-lt", "intellij-linux-2204-aws-m5d-lt",
		"intellij-linux-2004-aws-m5dn-lt", "intellij-linux-2204-aws-m5dn-lt",
		"intellij-linux-2204-large-disk-aws-1", "intellij-linux-2004-large-disk-aws-1",
		"intellij-linux-2204-large-disk-aws-2-i", "intellij-linux-2204-large-disk-aws-3-i",
		"intellij-linux-2204-large-disk-aws-4-i", "intellij-linux-2204-aws-2-i",
		"intellij-linux-2204-aws-1-i", "intellij-linux-2204-aws-4-i-", "intellij-linux-2204-aws-3-i",
	}, group: "Linux EC2 m5d.xlarge (4 vCPU Xeon, 16 GB)"},
	{prefixes: []string{"intellij-linux-hw-munit-"}, group: "Linux Munich i7-3770, 32 Gb"},
	{prefixes: []string{"intellij-linux-hw-EXC"}, group: "Linux JB Expo AMS i7-3770, 32 Gb"},
	{prefixes: []string{"intellij-linux-hw-hetzner", "intellij-linux-agg-hw-hetzner-agent"}, group: "linux-blade-hetzner"},
	{prefixes: []string{"intellij-windows-hw-hetzner"}, group: "windows-blade-hetzner"},
	{prefixes: []string{
		"intellij-macos-munit-741-large", "intellij-macos-de-unit-1219",
		"intellij-macos-munit-739-large", "intellij-macos-munit-738-large",
		"intellij-macos-munit-676-large",
	}, group: "Mac Pro Intel Xeon E5-2697v2 (4x2.7GHz), 24 RAM"},
	{prefixes: []string{"intellij-linux-performance-huge-aws-i"}, group: "Linux EC2 C6id.metal (128 CPU Xeon, 256 GB)"},
	{prefixes: []string{"qodana-aws-cpu-x64"}, group: "Linux EC2 c5a(d).xlarge (4 vCPU, 8 GB)"},
	{prefixes: []string{"qodana-linux-amd64-large"}, group: "Linux EC2 c5.large (2 vCPU, 4 GB)"},
	{prefixes: []string{
		"qodana-linux-amd64-xl", "qodana-linux-amd64-heavy", "intellij-linux-2004-aws-i",
		"intellij-linux-2004-aws-c5d", "intellij-linux-2004-aws-c5ad-lt", "intellij-linux-2004-aws-m5ad-lt",
	}, group: "Linux EC2 c5.xlarge (4 vCPU, 8 GB)"},
	{prefixes: []string{"intellij-linux-2204-aws-c5ad-lt"}, group: "Linux EC2 (2204) c5.xlarge (4 vCPU, 8 GB)"},
	{prefixes: []string{"intellij-linux-2004-aws-r5dn"}, group: "Linux EC2 r5dn.xlarge (4 vCPU, 32 GB)"},
	{prefixes: []string{"intellij-macos-perf-eqx"}, group: "Mac Mini M2 Pro (10 vCPU, 32 GB)"},
	{prefixes: []string{"intellij-windows-aws-i"}, group: "windows aws"},
	{regex: regexp.MustCompile(`ij-w.*-azr.*`), group: "windows-azure"},
	{prefixes: []string{"intellij-windows-hw-de-unit"}, group: "Windows Munich i7-13700, 64 Gb"},
	{prefixes: []string{"intellij-linux-hw-de-unit"}, group: "Linux Munich i7-13700, 64 Gb"},
	{prefixes: []string{"fleet-linux-aws-ui"}, group: "Linux Fleet AWS UI"},
	{prefixes: []string{"fleet-windows-aws-r5d", "fleet-windows-aws-m5d"}, group: "Windows Fleet AWS UI"},
	{prefixes: []string{"fleet-icri-ui-agent"}, group: "Mac Fleet AWS UI"},
	{prefixes: []string{"qodana-linux-arm64-memory-optimised"}, group: "Linux EC R7g.xlarge (4 vCPU ARM, 32 GB)"},
	{prefixes: []string{"cidr.performance."}, group: "Mac Cidr Performance"},
	{prefixes: []string{"intellij-linux-2204-aws-i4i"}, group: "Linux EC2 i4i.xlarge (4 vCPU Xeon, 32 GB)"},
	{prefixes: []string{"intellij-linux-2204-aws-r5d"}, group: "Linux EC2 r5d.xlarge (4 vCPU Xeon, 32 GB)"},
	{prefixes: []string{"intellij-linux-2004-aws-4-i-"}, group: "Linux EC2 c5ad.xlarge (4 vCPU EPYC, 8 GB)"},
	{prefixes: []string{"intellij-linux-2204-aws-c5d"}, group: "Linux EC2 c5d.xlarge (4 vCPU Xeon, 8 GB)"},
	{prefixes: []string{"intellij-macos-docker-hw-de-unit"}, group: "Mac Mini M2 Pro 12 CPU, 32 GB"},

	// Legacy fixed-name agents (migrated from MachineConfigurator.valueToGroup). Entries that a
	// prefix rule above already covers (blade-NNN, *-de-unit-NNNN, docker-*, unit-2200-large-*)
	// are intentionally omitted — they resolve identically through the prefixes.
	{names: []string{"intellij-macos-hw-unit-1550", "intellij-macos-hw-unit-1551", "intellij-macos-hw-unit-1772", "intellij-macos-hw-unit-1773"}, group: "macMini 2018"},
	{names: []string{"intellij-macos-hw-munit-716", "intellij-macos-hw-munit-721", "intellij-macos-hw-munit-722", "intellij-macos-hw-munit-723", "intellij-macos-hw-munit-724"}, group: "macMini Intel 3.2, 16GB"},
	{names: []string{
		"intellij-macos-hw-munit-608", "intellij-macos-hw-munit-689", "intellij-macos-hw-munit-690",
		"intellij-macos-hw-munit-691", "intellij-macos-hw-munit-692", "intellij-macos-hw-munit-693",
		"intellij-macos-hw-munit-694", "intellij-macos-hw-munit-695", "intellij-macos-hw-munit-696",
		"intellij-macos-hw-munit-697", "intellij-macos-hw-munit-698",
	}, group: "macMini M1, 16 Gb"},
	{names: []string{"intellij-macos-hw-unit-2204", "intellij-macos-hw-unit-2205", "intellij-macos-hw-unit-2206", "intellij-macos-hw-unit-2207"}, group: "macMini M1 2020"},
	{names: []string{
		"intellij-windows-hw-unit-498", "intellij-windows-hw-unit-499", "intellij-windows-hw-unit-449",
		"intellij-windows-hw-unit-463", "intellij-windows-hw-unit-493", "intellij-windows-hw-unit-504",
	}, group: "Windows Space i7-3770, 16Gb"},
	{names: []string{
		"intellij-linux-hw-unit-449", "intellij-linux-hw-unit-499", "intellij-linux-hw-unit-450",
		"intellij-linux-hw-unit-484", "intellij-linux-hw-unit-493", "intellij-linux-hw-unit-504",
		"intellij-linux-hw-unit-531", "intellij-linux-hw-unit-534", "intellij-linux-hw-unit-556",
		"intellij-linux-hw-unit-558",
	}, group: "Linux Space i7-3770, 16Gb"},
}

// GroupName returns the hardware-class group for a raw machine name (first matching rule
// wins), falling back to "Unknown" for unmapped names.
func GroupName(name string) string {
	for _, r := range rules {
		if r.regex != nil {
			if r.regex.MatchString(name) {
				return r.group
			}
			continue
		}
		if slices.Contains(r.names, name) {
			return r.group
		}
		for _, p := range r.prefixes {
			if strings.HasPrefix(name, p) {
				return r.group
			}
		}
	}
	return unknownGroup
}

// GroupPredicate renders the WHERE-clause suffix (the part after the column name) that selects
// exactly the given group's agents, e.g. "like 'a%' or machine like 'b%'". The suffix form is
// the /api/q filter contract: {f: col, q: suffix} renders as "(col suffix)", and since every
// term is a disjunct, callers may chain several predicates with " or <col> ". Returns "" for
// groups no prefix/name rule defines (Unknown, regex-matched groups) — callers fall back to
// filtering by the member list. Unlike a common prefix inferred from live members, the rule's
// own prefixes stay exact for groups mixing name stems and stable as agents churn.
func GroupPredicate(col, group string) string {
	var parts []string
	for _, r := range rules {
		if r.group != group {
			continue
		}
		for _, p := range r.prefixes {
			parts = append(parts, "like '"+sql_util.StringEscaper.Replace(p)+"%'")
		}
		if len(r.names) != 0 {
			quoted := make([]string, len(r.names))
			for i, n := range r.names {
				quoted[i] = "'" + sql_util.StringEscaper.Replace(n) + "'"
			}
			parts = append(parts, "in ("+strings.Join(quoted, ", ")+")")
		}
	}
	return strings.Join(parts, " or "+col+" ")
}

// GroupSQLExpr renders the rules as a ClickHouse multiIf() over the given machine column, so
// the database can GROUP BY the hardware class directly (which keeps the "latest N runs per
// group" pooling correct). Falls back to "Unknown".
func GroupSQLExpr(col string) string {
	var b strings.Builder
	b.WriteString("multiIf(")
	for _, r := range rules {
		switch {
		case r.regex != nil:
			b.WriteString("match(" + col + ", '" + sql_util.StringEscaper.Replace(r.regex.String()) + "')")
		case len(r.names) != 0:
			quoted := make([]string, len(r.names))
			for i, n := range r.names {
				quoted[i] = "'" + sql_util.StringEscaper.Replace(n) + "'"
			}
			b.WriteString(col + " IN (" + strings.Join(quoted, ", ") + ")")
		default:
			conds := make([]string, len(r.prefixes))
			for i, p := range r.prefixes {
				conds[i] = "startsWith(" + col + ", '" + sql_util.StringEscaper.Replace(p) + "')"
			}
			b.WriteString("(" + strings.Join(conds, " OR ") + ")")
		}
		b.WriteString(", '" + sql_util.StringEscaper.Replace(r.group) + "', ")
	}
	b.WriteString("'" + unknownGroup + "')")
	return b.String()
}
