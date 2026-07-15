package server

import (
	"regexp"
	"strings"
)

type machineGroupRule struct {
	prefixes []string
	regex    *regexp.Regexp
	group    string
}

// Order matters: first matching rule wins (same as the UI's if/else chain).
var machineGroupRules = []machineGroupRule{
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
}

var (
	reMachineInstanceID = regexp.MustCompile(`-i-[0-9a-f]{8,}$`)
	reMachineTrailIndex = regexp.MustCompile(`-\d+$`)
)

// machineStem strips the ephemeral AWS instance-id token and a trailing numeric index so
// that every VM of one class collapses to a single stable key. Used as the fallback for
// names not covered by machineGroupRules.
func machineStem(machine string) string {
	s := reMachineInstanceID.ReplaceAllString(machine, "")
	s = reMachineTrailIndex.ReplaceAllString(s, "")
	return s
}

// machineGroupName returns the hardware-class group for a raw machine name.
func machineGroupName(machine string) string {
	for _, rule := range machineGroupRules {
		if rule.regex != nil {
			if rule.regex.MatchString(machine) {
				return rule.group
			}
			continue
		}
		for _, p := range rule.prefixes {
			if strings.HasPrefix(machine, p) {
				return rule.group
			}
		}
	}
	return machineStem(machine)
}

// machineGroupSQLExpr renders the same rules as a ClickHouse multiIf() over the given
// machine column, so the database can GROUP BY the hardware class directly (which keeps the
// "latest 50 runs per group" pooling correct). Falls back to the stem expression.
func machineGroupSQLExpr(col string) string {
	var b strings.Builder
	b.WriteString("multiIf(")
	for _, rule := range machineGroupRules {
		if rule.regex != nil {
			b.WriteString("match(" + col + ", '" + chEscape(rule.regex.String()) + "')")
		} else {
			conds := make([]string, len(rule.prefixes))
			for i, p := range rule.prefixes {
				conds[i] = "startsWith(" + col + ", '" + chEscape(p) + "')"
			}
			b.WriteString("(" + strings.Join(conds, " OR ") + ")")
		}
		b.WriteString(", '" + chEscape(rule.group) + "', ")
	}
	// fallback: stem (instance-id then trailing index stripped)
	b.WriteString("replaceRegexpOne(replaceRegexpOne(" + col + ", '-i-[0-9a-f]{8,}$', ''), '-[0-9]+$', '')")
	b.WriteString(")")
	return b.String()
}

func chEscape(s string) string {
	return strings.ReplaceAll(s, "'", "''")
}
