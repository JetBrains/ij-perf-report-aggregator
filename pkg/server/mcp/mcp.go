// Package mcp exposes a Streamable HTTP MCP endpoint backed by the project's ClickHouse instance.
// The Register function mounts the handler on a chi-style router along with a small set of tools
// for discovering tables, projects, metric names, and metric values.
package mcp

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

// Register builds the MCP server, registers the tools, and mounts the Streamable HTTP handler
// at /api/mcp (and /api/mcp/* for session continuation) on the given handle function.
// The /api/ prefix matches the existing ingress routing used by the deployed backend.
func Register(dbUrl string, handle func(string, http.Handler)) {
	s := &service{dbUrl: dbUrl}

	server := sdk.NewServer(&sdk.Implementation{
		Name:    "ij-perf-report-aggregator",
		Version: "v1.0.0",
	}, nil)

	sdk.AddTool(server, &sdk.Tool{
		Name: "search_metric_names",
		Description: "Search distinct metric names recorded for the given project across every known table. " +
			"project is required; database/table/branch/machine/name_pattern are optional further filters. " +
			"Each result row is tagged with its source database and table.",
	}, s.searchMetricNames)

	sdk.AddTool(server, &sdk.Tool{
		Name: "search_metric_values",
		Description: "Fetch the most recent metric values for a given project and metric_name. " +
			"database/table are optional — when omitted the server scans every table tagged with measures " +
			"and returns rows from whichever ones contain the data. Each row is tagged with its source.",
	}, s.searchMetricValues)

	sdk.AddTool(server, &sdk.Tool{
		Name: "list_projects",
		Description: "List distinct projects that produced data in the lookback window. " +
			"All filters are optional. Each row is tagged with database/table so callers can pick the right one.",
	}, s.listProjects)

	sdk.AddTool(server, &sdk.Tool{
		Name:        "list_tables",
		Description: "List all (database, table) pairs the server can query. Useful as a discovery starting point.",
	}, s.listTablesTool)

	// Stateless mode: each HTTP request is independent so the handler works behind a
	// load balancer with multiple replicas (the in-memory session map is otherwise
	// per-pod and breaks across round-robin requests).
	handler := sdk.NewStreamableHTTPHandler(
		func(*http.Request) *sdk.Server { return server },
		&sdk.StreamableHTTPOptions{Stateless: true},
	)
	handle("/api/mcp", handler)
	handle("/api/mcp/*", handler)
}

type service struct {
	dbUrl string

	tablesMu     sync.Mutex
	tablesCache  []tableRef
	tablesCached time.Time
}

const tablesTTL = 10 * time.Minute

type tableRef struct {
	Database string `json:"database"`
	Table    string `json:"table"`
}

func (s *service) openConnection(database string) (driver.Conn, error) {
	if database == "" {
		database = "system"
	}
	return clickhouse.Open(&clickhouse.Options{
		Addr: []string{s.dbUrl},
		Auth: clickhouse.Auth{Database: database},
		Settings: map[string]any{
			"readonly":         1,
			"max_query_size":   1000000,
			"max_memory_usage": 3221225472,
		},
		DialTimeout:     5 * time.Second,
		ConnMaxLifetime: time.Hour,
	})
}

// listTables returns every (database, table) pair that has both `project` and `measures.name` columns.
func (s *service) listTables(ctx context.Context) ([]tableRef, error) {
	s.tablesMu.Lock()
	defer s.tablesMu.Unlock()
	if time.Since(s.tablesCached) < tablesTTL && s.tablesCache != nil {
		return s.tablesCache, nil
	}

	conn, err := s.openConnection("system")
	if err != nil {
		return nil, fmt.Errorf("open connection: %w", err)
	}
	defer conn.Close()

	rows, err := conn.Query(ctx, `
		select database, table
		from system.columns
		where name in ('measures.name', 'project')
		group by database, table
		having sum(name = 'measures.name') > 0 and sum(name = 'project') > 0
		order by database, table
	`)
	if err != nil {
		return nil, fmt.Errorf("discover tables: %w", err)
	}
	defer rows.Close()

	out := make([]tableRef, 0, 64)
	for rows.Next() {
		var r tableRef
		if err := rows.Scan(&r.Database, &r.Table); err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		out = append(out, r)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows: %w", err)
	}
	s.tablesCache = out
	s.tablesCached = time.Now()
	return out, nil
}

// resolveTables returns the candidate (database, table) pairs for a tool call.
// Empty db/table → all known tables; partial → filtered subset; both → exact match (must exist).
func (s *service) resolveTables(ctx context.Context, db, table string) ([]tableRef, error) {
	if db != "" {
		if err := validateIdentifier("database", db); err != nil {
			return nil, err
		}
	}
	if table != "" {
		if err := validateIdentifier("table", table); err != nil {
			return nil, err
		}
	}
	all, err := s.listTables(ctx)
	if err != nil {
		return nil, err
	}
	if db == "" && table == "" {
		if len(all) == 0 {
			return nil, fmt.Errorf("no known tables available")
		}
		return all, nil
	}
	out := make([]tableRef, 0, 4)
	for _, r := range all {
		if db != "" && r.Database != db {
			continue
		}
		if table != "" && r.Table != table {
			continue
		}
		out = append(out, r)
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("no known table matches database=%q table=%q", db, table)
	}
	return out, nil
}

type listTablesOutput struct {
	Tables []tableRef `json:"tables" jsonschema:"All (database, table) pairs that store performance measurements"`
	Count  int        `json:"count"`
}

func (s *service) listTablesTool(ctx context.Context, _ *sdk.CallToolRequest, _ struct{}) (*sdk.CallToolResult, listTablesOutput, error) {
	all, err := s.listTables(ctx)
	if err != nil {
		return nil, listTablesOutput{}, err
	}
	return nil, listTablesOutput{Tables: all, Count: len(all)}, nil
}

func (s *service) query(ctx context.Context, sql string, args []any) (driver.Rows, driver.Conn, error) {
	conn, err := s.openConnection("system")
	if err != nil {
		return nil, nil, fmt.Errorf("open connection: %w", err)
	}
	rows, err := conn.Query(ctx, sql, args...)
	if err != nil {
		conn.Close()
		return nil, nil, fmt.Errorf("query failed: %w", err)
	}
	return rows, conn, nil
}

func buildUnion(tables []tableRef, perTable func(tableRef) (string, []any)) (string, []any) {
	var sb strings.Builder
	args := make([]any, 0, len(tables)*4)
	for i, r := range tables {
		if i > 0 {
			sb.WriteString(" union all ")
		}
		sb.WriteString("(")
		sql, a := perTable(r)
		sb.WriteString(sql)
		args = append(args, a...)
		sb.WriteString(")")
	}
	return sb.String(), args
}

func validateIdentifier(field, value string) error {
	if value == "" {
		return fmt.Errorf("%s is required", field)
	}
	for _, r := range value {
		if !(r == '_' || (r >= '0' && r <= '9') || (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')) {
			return fmt.Errorf("%s contains invalid character %q", field, r)
		}
	}
	return nil
}

func quoteIdentifier(name string) string {
	return "`" + name + "`"
}

func clamp(v, maxVal, defaultVal int) int {
	if v == 0 {
		return defaultVal
	}
	if v < 1 {
		return 1
	}
	if v > maxVal {
		return maxVal
	}
	return v
}
