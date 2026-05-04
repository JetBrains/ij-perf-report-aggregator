package mcp

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"strings"

	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

type searchMetricNamesInput struct {
	Project     string `json:"project"                jsonschema:"Project to look up metric names for (exact match)"`
	NamePattern string `json:"name_pattern,omitempty" jsonschema:"Optional SQL LIKE pattern (e.g. 'startup%') to narrow metric names"`
	Branch      string `json:"branch,omitempty"       jsonschema:"Branch filter (default: master)"`
	Machine     string `json:"machine,omitempty"      jsonschema:"Optional machine LIKE pattern (e.g. 'intellij-linux-hw-hetzner%')"`
	Database    string `json:"database,omitempty"     jsonschema:"Optional database to restrict the scan to"`
	Table       string `json:"table,omitempty"        jsonschema:"Optional table to restrict the scan to"`
	Days        int    `json:"days,omitempty"         jsonschema:"Lookback window in days (default 30, max 365)"`
	Limit       int    `json:"limit,omitempty"        jsonschema:"Max number of (database, table, name) tuples to return (default 500, max 5000)"`
}

type metricNameRow struct {
	Database string `json:"database"`
	Table    string `json:"table"`
	Name     string `json:"name"`
}

type searchMetricNamesOutput struct {
	Rows  []metricNameRow `json:"rows"  jsonschema:"Distinct (database, table, metric name) tuples matching the filters"`
	Count int             `json:"count"`
}

func (s *service) searchMetricNames(ctx context.Context, _ *sdk.CallToolRequest, in searchMetricNamesInput) (*sdk.CallToolResult, searchMetricNamesOutput, error) {
	if in.Project == "" {
		return nil, searchMetricNamesOutput{}, errors.New("project is required")
	}
	if in.Branch == "" {
		in.Branch = defaultBranch
	}
	tables, err := s.resolveTables(ctx, in.Database, in.Table)
	if err != nil {
		return nil, searchMetricNamesOutput{}, err
	}
	days := min(max(cmp.Or(in.Days, 30), 1), 365)
	limit := min(max(cmp.Or(in.Limit, 500), 1), 5000)

	perTable := func(r tableRef) (string, []any) {
		var sb strings.Builder
		fmt.Fprintf(&sb, "select ? as db_name, ? as table_name, arrayJoin(`measures.name`) as metric_name from %s.%s where generated_time > subtractDays(now(), ?) and project = ?",
			r.Database, r.Table)
		args := []any{r.Database, r.Table, days, in.Project}
		args = appendBranchMachine(&sb, args, in.Branch, in.Machine)
		return sb.String(), args
	}

	innerSQL, args := buildUnion(tables, perTable)
	sql := "select distinct db_name, table_name, metric_name from (" + innerSQL + ") as u"
	if in.NamePattern != "" {
		sql += " where metric_name like ?"
		args = append(args, in.NamePattern)
	}
	sql += " order by db_name, table_name, metric_name limit ?"
	args = append(args, limit)

	rows, err := s.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, searchMetricNamesOutput{}, fmt.Errorf("search_metric_names: %w", err)
	}
	defer rows.Close()

	out := searchMetricNamesOutput{Rows: make([]metricNameRow, 0, 64)}
	for rows.Next() {
		var r metricNameRow
		if err := rows.Scan(&r.Database, &r.Table, &r.Name); err != nil {
			return nil, searchMetricNamesOutput{}, fmt.Errorf("scan: %w", err)
		}
		out.Rows = append(out.Rows, r)
	}
	if err := rows.Err(); err != nil {
		return nil, searchMetricNamesOutput{}, fmt.Errorf("rows: %w", err)
	}
	out.Count = len(out.Rows)
	return nil, out, nil
}
