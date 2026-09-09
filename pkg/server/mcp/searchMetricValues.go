package mcp

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

type searchMetricValuesInput struct {
	Project    string `json:"project"             jsonschema:"Project name (exact match) to query values for"`
	MetricName string `json:"metric_name"         jsonschema:"Metric name (exact match in measures.name) to retrieve values for"`
	Branch     string `json:"branch,omitempty"    jsonschema:"Branch filter (default: master)"`
	Machine    string `json:"machine,omitempty"   jsonschema:"Optional machine LIKE pattern. Runs from different machine pools are not comparable; filter to one pool before comparing values."`
	Database   string `json:"database,omitempty"  jsonschema:"Optional database to restrict the scan to"`
	Table      string `json:"table,omitempty"     jsonschema:"Optional table to restrict the scan to"`
	Days       int    `json:"days,omitempty"      jsonschema:"Lookback window in days (default 30, max 365). Reached in full only when limit does not cut the answer short — the covered field says what came back."`
	Limit      int    `json:"limit,omitempty"     jsonschema:"Max rows returned, ordered by generated_time desc (default 200, max 5000). Rows are newest-first, so a limit smaller than the window narrows it to the most recent days rather than sampling across it."`
	Aggregate  string `json:"aggregate,omitempty" jsonschema:"Set to \"daily\" for one row per calendar day (build count, median, min, max) instead of raw per-build rows. Prefer it for any trend or baseline question: a daily series spans the whole lookback window instead of collapsing to the most recent days, and median together with min is what separates a real level shift from noise."`
}

type metricValueRow struct {
	GeneratedTime    string  `json:"generated_time"`
	BuildID          uint32  `json:"tc_build_id"`
	Value            float64 `json:"value"`
	BuildNumber      string  `json:"build_number,omitempty"          jsonschema:"Marketing build number, e.g. 261.27258.48 — FUS product_build without the product prefix. Absent for Dev Server runs."`
	InstallerBuildID uint32  `json:"tc_installer_build_id,omitempty" jsonschema:"TeamCity build id of the installer tested; tc_build_id identifies the perf-test run itself."`
}

// dailyBucket is one calendar day of a metric, aggregated over distinct builds. Median and min
// are the pair that answers "did the level shift": a real regression moves the median (and often
// the floor) and keeps it moved, while a noisy metric only throws occasional high outliers that
// leave both alone.
type dailyBucket struct {
	Date   string  `json:"date"   jsonschema:"Calendar day, YYYY-MM-DD"`
	Builds int     `json:"builds" jsonschema:"Distinct builds measured that day — the sample size behind median/min/max"`
	Median float64 `json:"median" jsonschema:"Median across that day's builds; the level to compare across days"`
	Min    float64 `json:"min"    jsonschema:"Fastest run that day (the floor)"`
	Max    float64 `json:"max"    jsonschema:"Slowest run that day; single-run outliers land here and say little on their own"`
}

type metricValueGroup struct {
	Database string           `json:"database"`
	Table    string           `json:"table"`
	Rows     []metricValueRow `json:"rows,omitempty"  jsonschema:"Per-build measurements from this (database, table) ordered by generated_time desc. Absent when aggregate=daily."`
	Daily    []dailyBucket    `json:"daily,omitempty" jsonschema:"Per-day aggregates from this (database, table) ordered by date desc. Present only when aggregate=daily."`
}

type searchMetricValuesOutput struct {
	Project    string             `json:"project"`
	MetricName string             `json:"metric_name"`
	Branch     string             `json:"branch"`
	Groups     []metricValueGroup `json:"groups"            jsonschema:"Results grouped by source (database, table). Empty if no data found."`
	Count      int                `json:"count"             jsonschema:"Rows returned across all groups — per-build measurements, or day buckets when aggregate=daily"`
	Covered    string             `json:"covered,omitempty" jsonschema:"The time span actually returned, and how much of the requested lookback it spans. When it is narrower than days, limit cut the answer and older data exists that you have not seen — never read it as the full history."`
	Notes      []string           `json:"notes,omitempty"   jsonschema:"Read these before drawing conclusions: they say why the result is empty or partial. Missing when the answer is complete."`
}

const aggregateDaily = "daily"

func (s *service) searchMetricValues(ctx context.Context, _ *sdk.CallToolRequest, in searchMetricValuesInput) (*sdk.CallToolResult, searchMetricValuesOutput, error) {
	if in.Project == "" {
		return nil, searchMetricValuesOutput{}, errors.New("project is required")
	}
	if in.MetricName == "" {
		return nil, searchMetricValuesOutput{}, errors.New("metric_name is required")
	}
	if in.Aggregate != "" && in.Aggregate != aggregateDaily {
		return nil, searchMetricValuesOutput{}, fmt.Errorf("aggregate must be %q or omitted, got %q", aggregateDaily, in.Aggregate)
	}
	if in.Branch == "" {
		in.Branch = defaultBranch
	}
	tables, err := s.resolveTables(ctx, in.Database, in.Table)
	if err != nil {
		return nil, searchMetricValuesOutput{}, err
	}
	days := min(max(cmp.Or(in.Days, 30), 1), 365)
	limit := min(max(cmp.Or(in.Limit, 200), 1), 5000)
	daily := in.Aggregate == aggregateDaily

	perTable := func(r tableRef) (string, []any) {
		buildComponentsExpr := absentBuildComponentsSQL
		if r.HasBuildComponents {
			buildComponentsExpr = correctedBuildComponentsSQL
		}
		installerIDExpr := "toUInt32(0) as inst_id"
		if r.HasInstallerID {
			installerIDExpr = "toUInt32(tc_installer_build_id) as inst_id"
		}
		var sb strings.Builder
		fmt.Fprintf(&sb,
			"select ? as db_name, ? as table_name, "+
				"generated_time as gen_time, tc_build_id as build_id, "+
				"toFloat64(`measures.value`[idx]) as value, "+
				"%s, %s "+
				"from %s.%s array join arrayEnumerate(`measures.name`) as idx "+
				"where project = ? and `measures.name`[idx] = ? "+
				"and generated_time > subtractDays(now(), ?)",
			buildComponentsExpr, installerIDExpr, r.Database, r.Table)
		args := []any{r.Database, r.Table, in.Project, in.MetricName, days}
		args = appendBranchMachine(&sb, args, in.Branch, in.Machine)
		return sb.String(), args
	}

	innerSQL, args := buildUnion(tables, perTable)
	var sql string
	if daily {
		// `distinct` first: a build can appear more than once per (day, table), and duplicates would
		// double-weight it in the median.
		sql = "select db_name, table_name, toString(toDate(gen_time)) as day, " +
			"toUInt32(count()) as builds, quantileExact(0.5)(value) as med, " +
			"min(value) as lo, max(value) as hi from (" +
			"select distinct db_name, table_name, gen_time, build_id, value from (" + innerSQL + ") as d" +
			") as u group by db_name, table_name, day order by day desc limit ?"
	} else {
		sql = "select db_name, table_name, toString(gen_time) as gen_time, build_id, value, " +
			"bc1, bc2, bc3, inst_id from (" +
			innerSQL + ") as u order by gen_time desc limit ?"
	}
	args = append(args, limit)

	rows, err := s.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, searchMetricValuesOutput{}, fmt.Errorf("search_metric_values: %w", err)
	}
	defer rows.Close()

	type groupKey struct{ database, table string }
	out := searchMetricValuesOutput{Project: in.Project, MetricName: in.MetricName, Branch: in.Branch, Groups: []metricValueGroup{}}
	groupIndex := make(map[groupKey]int)
	// Rows arrive newest-first, so the first one seen is the newest and the last is the oldest.
	var newest, oldest string
	group := func(key groupKey) *metricValueGroup {
		idx, ok := groupIndex[key]
		if !ok {
			idx = len(out.Groups)
			groupIndex[key] = idx
			out.Groups = append(out.Groups, metricValueGroup{Database: key.database, Table: key.table})
		}
		return &out.Groups[idx]
	}
	for rows.Next() {
		var key groupKey
		var stamp string
		if daily {
			var b dailyBucket
			var builds uint32
			if err := rows.Scan(&key.database, &key.table, &b.Date, &builds, &b.Median, &b.Min, &b.Max); err != nil {
				return nil, searchMetricValuesOutput{}, fmt.Errorf("scan: %w", err)
			}
			b.Builds = int(builds)
			stamp = b.Date
			g := group(key)
			g.Daily = append(g.Daily, b)
		} else {
			var r metricValueRow
			var bc1, bc2, bc3 uint16
			if err := rows.Scan(&key.database, &key.table, &r.GeneratedTime, &r.BuildID, &r.Value,
				&bc1, &bc2, &bc3, &r.InstallerBuildID); err != nil {
				return nil, searchMetricValuesOutput{}, fmt.Errorf("scan: %w", err)
			}
			r.BuildNumber = formatBuildNumber(bc1, bc2, bc3)
			stamp = r.GeneratedTime
			g := group(key)
			g.Rows = append(g.Rows, r)
		}
		if newest == "" {
			newest = stamp
		}
		oldest = stamp
		out.Count++
	}
	if err := rows.Err(); err != nil {
		return nil, searchMetricValuesOutput{}, fmt.Errorf("rows: %w", err)
	}
	out.Covered = coveredSpan(oldest, newest, days)
	if out.Count == 0 {
		filters := []string{
			fmt.Sprintf("project=%q", in.Project),
			fmt.Sprintf("metric_name=%q", in.MetricName),
			fmt.Sprintf("branch=%q", in.Branch),
			fmt.Sprintf("last %dd", days),
		}
		if in.Machine != "" {
			filters = append(filters, fmt.Sprintf("machine=%q", in.Machine))
		}
		out.Notes = append(out.Notes, noRowsNote("measurement", filters, tables),
			"metric_name and project must match exactly — list the real ones with search_metric_names / list_projects before concluding the metric has no data")
	}
	if note := narrowedWindowNote(out.Count, limit, days, out.Covered, daily); note != "" {
		out.Notes = append(out.Notes, note)
	}
	return nil, out, nil
}

// coveredSpan renders the span the answer actually covers, e.g.
// "2026-09-03..2026-09-09 (7 of 90 requested days)". The row list alone cannot show this: an
// answer cut short by `limit` looks complete, just shorter.
func coveredSpan(oldest, newest string, days int) string {
	if oldest == "" || newest == "" {
		return ""
	}
	from, to := dayOf(oldest), dayOf(newest)
	span := from + ".." + to
	o, errO := time.Parse(time.DateOnly, from)
	n, errN := time.Parse(time.DateOnly, to)
	if errO != nil || errN != nil {
		return span
	}
	return fmt.Sprintf("%s (%d of %d requested days)", span, int(n.Sub(o).Hours()/24)+1, days)
}

func dayOf(stamp string) string {
	if len(stamp) < len(time.DateOnly) {
		return stamp
	}
	return stamp[:len(time.DateOnly)]
}

// narrowedWindowNote fires when `limit` truncated the answer. Because rows come back newest-first,
// truncation does not thin the window evenly — it discards the OLD end. Asking for 90 days of a
// project that runs ~36 builds/day and taking 200 rows returns about six days, all of them recent,
// while looking like a full history. A caller that misses this computes its "long-run baseline"
// from a window that may lie entirely after the change it was asked to explain, and then reports
// the change as noise. Say so loudly, and name the way out.
func narrowedWindowNote(rows, limit, days int, covered string, daily bool) string {
	if rows < limit {
		return ""
	}
	unit, escape := "row", `re-ask with aggregate="daily" to get the whole window as one row per day, or raise limit`
	if daily {
		unit, escape = "day", "raise limit to cover the whole window"
	}
	return fmt.Sprintf("TRUNCATED at limit=%d %ss, newest-first: this answer covers %s of the %d requested — older data exists and is NOT included. "+
		"Do not read it as the full history, and do not compute a baseline or call a change noise from it: the window may lie entirely on one side of the change you are investigating. %s.",
		limit, unit, cmp.Or(covered, "an unknown span"), days, escape)
}
