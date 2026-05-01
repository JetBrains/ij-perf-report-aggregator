package mcp

import (
	"context"
	"fmt"
	"strings"

	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

type listProjectsInput struct {
	Branch         string `json:"branch,omitempty"          jsonschema:"Branch filter (default: master)"`
	Machine        string `json:"machine,omitempty"         jsonschema:"Optional machine LIKE pattern"`
	ProjectPattern string `json:"project_pattern,omitempty" jsonschema:"Optional SQL LIKE pattern to narrow project names"`
	Database       string `json:"database,omitempty"        jsonschema:"Optional database to restrict the scan to"`
	Table          string `json:"table,omitempty"           jsonschema:"Optional table to restrict the scan to"`
	Days           int    `json:"days,omitempty"            jsonschema:"Lookback window in days (default 30, max 365)"`
	Limit          int    `json:"limit,omitempty"           jsonschema:"Max number of (database, table, project) tuples to return (default 500, max 5000)"`
}

type projectRow struct {
	Database string `json:"database"`
	Table    string `json:"table"`
	Project  string `json:"project"`
}

type listProjectsOutput struct {
	Rows  []projectRow `json:"rows"`
	Count int          `json:"count"`
}

func (s *service) listProjects(ctx context.Context, _ *sdk.CallToolRequest, in listProjectsInput) (*sdk.CallToolResult, listProjectsOutput, error) {
	if in.Branch == "" {
		in.Branch = "master"
	}
	tables, err := s.resolveTables(ctx, in.Database, in.Table)
	if err != nil {
		return nil, listProjectsOutput{}, err
	}
	days := clamp(in.Days, 365, 30)
	limit := clamp(in.Limit, 5000, 500)

	perTable := func(r tableRef) (string, []any) {
		var sb strings.Builder
		fmt.Fprintf(&sb, "select distinct '%s' as db_name, '%s' as table_name, project as project_name from %s.%s where generated_time > subtractDays(now(), ?)",
			r.Database, r.Table, quoteIdentifier(r.Database), quoteIdentifier(r.Table))
		args := []any{days}
		if in.Branch != "" {
			sb.WriteString(" and branch = ?")
			args = append(args, in.Branch)
		}
		if in.Machine != "" {
			sb.WriteString(" and machine like ?")
			args = append(args, in.Machine)
		}
		if in.ProjectPattern != "" {
			sb.WriteString(" and project like ?")
			args = append(args, in.ProjectPattern)
		}
		return sb.String(), args
	}

	innerSQL, args := buildUnion(tables, perTable)
	sql := "select db_name, table_name, project_name from (" + innerSQL + ") as u order by db_name, table_name, project_name limit ?"
	args = append(args, limit)

	rows, conn, err := s.query(ctx, sql, args)
	if err != nil {
		return nil, listProjectsOutput{}, err
	}
	defer conn.Close()
	defer rows.Close()

	out := listProjectsOutput{Rows: make([]projectRow, 0, 64)}
	for rows.Next() {
		var r projectRow
		if err := rows.Scan(&r.Database, &r.Table, &r.Project); err != nil {
			return nil, listProjectsOutput{}, fmt.Errorf("scan: %w", err)
		}
		out.Rows = append(out.Rows, r)
	}
	if err := rows.Err(); err != nil {
		return nil, listProjectsOutput{}, fmt.Errorf("rows: %w", err)
	}
	out.Count = len(out.Rows)
	return nil, out, nil
}
