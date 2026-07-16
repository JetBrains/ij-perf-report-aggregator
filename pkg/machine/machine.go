// Package machine maps a raw perf-agent name to its hardware-class group.
//
// A raw `machine` value is an ephemeral agent name (usually "<class>-i-<aws-instance-id>").
// Comparing branches requires grouping runs by hardware *class*, never by the raw instance,
// otherwise runs on different agents get pooled into one median and produce phantom
// regressions (AT-4930) — e.g. the broad "%linux%" filter matches both the powerful
// "intellij-linux-performance-aws-*" and the weak "intellij-linux-performance-tiny-aws-*".
//
// machine-groups.json is the single source of truth, shared with the frontend
// (dashboard/new-dashboard/src/configurators/MachineConfigurator.ts:getMachineGroupName).
// The one intentional difference: the UI falls back to "Unknown" for unmapped names, while
// GroupName falls back to the machine stem (instance-id/index stripped). That keeps distinct
// unmapped classes apart instead of pooling them into one bucket, so comparisons stay correct
// even if a machine class is not yet listed here.
package machine

import (
	_ "embed"
	"encoding/json"
	"regexp"
	"strings"
)

//go:embed machine-groups.json
var groupsJSON []byte

type groupRule struct {
	Group    string   `json:"group"`
	Prefixes []string `json:"prefixes"`
	Regex    string   `json:"regex"`
}

type compiledRule struct {
	prefixes []string
	regex    *regexp.Regexp
	group    string
}

var rules []compiledRule

var (
	reInstanceID = regexp.MustCompile(`-i-[0-9a-f]{8,}$`)
	reTrailIndex = regexp.MustCompile(`-\d+$`)
)

func init() {
	var raw []groupRule
	if err := json.Unmarshal(groupsJSON, &raw); err != nil {
		panic("machine: invalid groups.json: " + err.Error())
	}
	rules = make([]compiledRule, len(raw))
	for i, r := range raw {
		cr := compiledRule{prefixes: r.Prefixes, group: r.Group}
		if r.Regex != "" {
			cr.regex = regexp.MustCompile(r.Regex)
		}
		rules[i] = cr
	}
}

// Stem strips the ephemeral AWS instance-id token and a trailing numeric index so that every
// VM of one class collapses to a single stable key. Used as the fallback for names not
// covered by the rules.
func Stem(name string) string {
	s := reInstanceID.ReplaceAllString(name, "")
	s = reTrailIndex.ReplaceAllString(s, "")
	return s
}

// GroupName returns the hardware-class group for a raw machine name (first matching rule
// wins), falling back to the stem for unmapped names.
func GroupName(name string) string {
	for _, rule := range rules {
		if rule.regex != nil {
			if rule.regex.MatchString(name) {
				return rule.group
			}
			continue
		}
		for _, p := range rule.prefixes {
			if strings.HasPrefix(name, p) {
				return rule.group
			}
		}
	}
	return Stem(name)
}

// GroupSQLExpr renders the rules as a ClickHouse multiIf() over the given machine column, so
// the database can GROUP BY the hardware class directly (which keeps the "latest N runs per
// group" pooling correct). Falls back to the stem expression.
func GroupSQLExpr(col string) string {
	var b strings.Builder
	b.WriteString("multiIf(")
	for _, rule := range rules {
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
