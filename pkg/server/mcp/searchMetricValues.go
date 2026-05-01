package mcp

import (
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
	GeneratedTime string  `json:"generated_time"`
	BuildID       uint32  `json:"tc_build_id"`
	Value         float64 `json:"value"`
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
		in.Branch = "master"
	}
	tables, err := s.resolveTables(ctx, in.Database, in.Table)
	if err != nil {
		return nil, searchMetricValuesOutput{}, err
	}
	days := clamp(in.Days, 365, 30)
	limit := clamp(in.Limit, 5000, 200)

	perTable := func(r tableRef) (string, []any) {
		var sb strings.Builder
		fmt.Fprintf(&sb,
			"select '%s' as db_name, '%s' as table_name, "+
				"toString(generated_time) as gen_time, tc_build_id as build_id, "+
				"toFloat64(`measures.value`[idx]) as value "+
				"from %s.%s array join arrayEnumerate(`measures.name`) as idx "+
				"where project = ? and `measures.name`[idx] = ? "+
				"and generated_time > subtractDays(now(), ?)",
			r.Database, r.Table, quoteIdentifier(r.Database), quoteIdentifier(r.Table))
		args := []any{in.Project, in.MetricName, days}
		if in.Branch != "" {
			sb.WriteString(" and branch = ?")
			args = append(args, in.Branch)
		}
		if in.Machine != "" {
			sb.WriteString(" and machine like ?")
			args = append(args, in.Machine)
		}
		return sb.String(), args
	}

	innerSQL, args := buildUnion(tables, perTable)
	sql := "select db_name, table_name, gen_time, build_id, value from (" +
		innerSQL + ") as u order by gen_time desc limit ?"
	args = append(args, limit)

	rows, conn, err := s.query(ctx, sql, args)
	if err != nil {
		return nil, searchMetricValuesOutput{}, err
	}
	defer conn.Close()
	defer rows.Close()

	type groupKey struct{ database, table string }
	out := searchMetricValuesOutput{Project: in.Project, MetricName: in.MetricName, Branch: in.Branch, Groups: []metricValueGroup{}}
	groupIndex := make(map[groupKey]int)
	for rows.Next() {
		var key groupKey
		var r metricValueRow
		if err := rows.Scan(&key.database, &key.table, &r.GeneratedTime, &r.BuildID, &r.Value); err != nil {
			return nil, searchMetricValuesOutput{}, fmt.Errorf("scan: %w", err)
		}
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
