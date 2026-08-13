package mcp

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"strings"

	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

type searchMetricValuesInput struct {
	Project    string `json:"project"            jsonschema:"Project name (exact match) to query values for"`
	MetricName string `json:"metric_name"        jsonschema:"Metric name (exact match in measures.name) to retrieve values for"`
	Branch     string `json:"branch,omitempty"   jsonschema:"Branch filter (default: master)"`
	Machine    string `json:"machine,omitempty"  jsonschema:"Optional machine LIKE pattern"`
	Database   string `json:"database,omitempty" jsonschema:"Optional database to restrict the scan to"`
	Table      string `json:"table,omitempty"    jsonschema:"Optional table to restrict the scan to"`
	Days       int    `json:"days,omitempty"     jsonschema:"Lookback window in days (default 30, max 365)"`
	Limit      int    `json:"limit,omitempty"    jsonschema:"Max rows returned, ordered by generated_time desc (default 200, max 5000)"`
}

type metricValueRow struct {
	GeneratedTime    string  `json:"generated_time"`
	BuildID          uint32  `json:"tc_build_id"`
	Value            float64 `json:"value"`
	BuildNumber      string  `json:"build_number,omitempty"          jsonschema:"Marketing build number, e.g. 261.27258.48 — FUS product_build without the product prefix. Absent for Dev Server runs."`
	InstallerBuildID uint32  `json:"tc_installer_build_id,omitempty" jsonschema:"TeamCity build id of the installer tested; tc_build_id identifies the perf-test run itself."`
}

type metricValueGroup struct {
	Database string           `json:"database"`
	Table    string           `json:"table"`
	Rows     []metricValueRow `json:"rows"     jsonschema:"Measurements from this (database, table) ordered by generated_time desc"`
}

type searchMetricValuesOutput struct {
	Project    string             `json:"project"`
	MetricName string             `json:"metric_name"`
	Branch     string             `json:"branch"`
	Groups     []metricValueGroup `json:"groups"      jsonschema:"Results grouped by source (database, table). Empty if no data found."`
	Count      int                `json:"count"       jsonschema:"Total number of measurement rows across all groups"`
}

func (s *service) searchMetricValues(ctx context.Context, _ *sdk.CallToolRequest, in searchMetricValuesInput) (*sdk.CallToolResult, searchMetricValuesOutput, error) {
	if in.Project == "" {
		return nil, searchMetricValuesOutput{}, errors.New("project is required")
	}
	if in.MetricName == "" {
		return nil, searchMetricValuesOutput{}, errors.New("metric_name is required")
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
	sql := "select db_name, table_name, toString(gen_time) as gen_time, build_id, value, " +
		"bc1, bc2, bc3, inst_id from (" +
		innerSQL + ") as u order by gen_time desc limit ?"
	args = append(args, limit)

	rows, err := s.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, searchMetricValuesOutput{}, fmt.Errorf("search_metric_values: %w", err)
	}
	defer rows.Close()

	type groupKey struct{ database, table string }
	out := searchMetricValuesOutput{Project: in.Project, MetricName: in.MetricName, Branch: in.Branch, Groups: []metricValueGroup{}}
	groupIndex := make(map[groupKey]int)
	for rows.Next() {
		var key groupKey
		var r metricValueRow
		var bc1, bc2, bc3 uint16
		if err := rows.Scan(&key.database, &key.table, &r.GeneratedTime, &r.BuildID, &r.Value,
			&bc1, &bc2, &bc3, &r.InstallerBuildID); err != nil {
			return nil, searchMetricValuesOutput{}, fmt.Errorf("scan: %w", err)
		}
		r.BuildNumber = formatBuildNumber(bc1, bc2, bc3)
		idx, ok := groupIndex[key]
		if !ok {
			idx = len(out.Groups)
			groupIndex[key] = idx
			out.Groups = append(out.Groups, metricValueGroup{Database: key.database, Table: key.table})
		}
		out.Groups[idx].Rows = append(out.Groups[idx].Rows, r)
		out.Count++
	}
	if err := rows.Err(); err != nil {
		return nil, searchMetricValuesOutput{}, fmt.Errorf("rows: %w", err)
	}
	return nil, out, nil
}
