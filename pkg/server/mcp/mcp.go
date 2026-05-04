// Package mcp exposes a Streamable HTTP MCP endpoint backed by the project's ClickHouse instance.
// The Register function mounts the handler on a chi-style router along with a small set of tools
// for discovering tables, projects, metric names, and metric values.
package mcp

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

// chConn is the narrow slice of clickhouse-go that the MCP service depends on.
// Production wires a real driver.Conn via driverConnAdapter; tests pass a fake.
type chConn interface {
	Query(ctx context.Context, query string, args ...any) (chRows, error)
	Close() error
}

// chRows is the subset of driver.Rows used by tools (Next/Scan/Err/Close only).
type chRows interface {
	Next() bool
	Scan(dest ...any) error
	Err() error
	Close() error
}

// Register opens a long-lived ClickHouse connection pool, builds the MCP server,
// and mounts the Streamable HTTP handler at /api/mcp (and /api/mcp/* for session
// continuation) on the given handle function. The /api/ prefix matches the
// existing ingress routing used by the deployed backend.
func Register(dbUrl string, handle func(string, http.Handler)) error {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{dbUrl},
		Auth: clickhouse.Auth{Database: "system"},
		Settings: map[string]any{
			"readonly":         1,
			"max_query_size":   1000000,
			"max_memory_usage": 3221225472,
		},
		DialTimeout:     5 * time.Second,
		ConnMaxLifetime: time.Hour,
	})
	if err != nil {
		return fmt.Errorf("mcp: open clickhouse: %w", err)
	}
	s := newService(driverConnAdapter{conn: conn})
	server := s.buildServer()

	// Stateless mode: each HTTP request is independent so the handler works behind a
	// load balancer with multiple replicas (the in-memory session map is otherwise
	// per-pod and breaks across round-robin requests).
	handler := sdk.NewStreamableHTTPHandler(
		func(*http.Request) *sdk.Server { return server },
		&sdk.StreamableHTTPOptions{Stateless: true},
	)
	handle("/api/mcp", handler)
	handle("/api/mcp/*", handler)
	return nil
}

type service struct {
	db chConn

	tablesMu     sync.Mutex
	tablesCache  []tableRef
	tablesCached time.Time
}

const tablesTTL = 10 * time.Minute

type tableRef struct {
	Database       string `json:"database"`
	Table          string `json:"table"`
	HasBuildTime   bool   `json:"-"`
	HasInstallerID bool   `json:"-"`
}

func newService(db chConn) *service {
	return &service{db: db}
}

func (s *service) buildServer() *sdk.Server {
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

	sdk.AddTool(server, &sdk.Tool{
		Name: "get_build",
		Description: "Look up a TeamCity build by its tc_build_id. Returns build-level metadata at the response root " +
			"(branch, build_time, machine, teamcity_url, first_commit, last_commit) and the list of distinct " +
			"(database, table, project) tuples this build produced data for. Use search_metric_values afterwards " +
			"to fetch actual measurements for any project of interest.",
	}, s.getBuild)

	return server
}

// driverConnAdapter wraps a clickhouse-go driver.Conn into the package-local chConn.
type driverConnAdapter struct {
	conn driver.Conn
}

func (a driverConnAdapter) Query(ctx context.Context, sql string, args ...any) (chRows, error) {
	return a.conn.Query(ctx, sql, args...)
}

func (a driverConnAdapter) Close() error { return a.conn.Close() }

// listTables returns every (database, table) pair that has both `project` and `measures.name` columns.
func (s *service) listTables(ctx context.Context) ([]tableRef, error) {
	s.tablesMu.Lock()
	if time.Since(s.tablesCached) < tablesTTL && s.tablesCache != nil {
		cached := s.tablesCache
		s.tablesMu.Unlock()
		return cached, nil
	}
	s.tablesMu.Unlock()

	// Query without holding the cache mutex so a slow ClickHouse round-trip
	// doesn't serialize concurrent callers behind it. If two callers race to
	// refresh, last writer wins — both get correct data.
	rows, err := s.db.Query(ctx, `
		select database, table,
		       sum(name = 'build_time') > 0 as has_build_time,
		       sum(name = 'tc_installer_build_id') > 0 as has_installer_id
		from system.columns
		where name in ('measures.name', 'project', 'build_time', 'tc_installer_build_id')
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
		if err := rows.Scan(&r.Database, &r.Table, &r.HasBuildTime, &r.HasInstallerID); err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		// Tools interpolate r.Database/r.Table directly into SQL. Drop anything that
		// doesn't fit the safe-identifier shape so we never produce a query that
		// needs quoting; in practice every real perf-report table satisfies this.
		if validateIdentifier("database", r.Database) != nil || validateIdentifier("table", r.Table) != nil {
			slog.Warn("mcp: skipping non-identifier table name", "db", r.Database, "table", r.Table)
			continue
		}
		out = append(out, r)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows: %w", err)
	}

	s.tablesMu.Lock()
	s.tablesCache = out
	s.tablesCached = time.Now()
	s.tablesMu.Unlock()
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
			return nil, errors.New("no known tables available")
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

func buildUnion(tables []tableRef, perTable func(tableRef) (string, []any)) (string, []any) {
	var sb strings.Builder
	sb.Grow(len(tables) * 600)
	args := make([]any, 0, len(tables)*5)
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

// appendBranchMachine writes the optional `and branch = ?` / `and machine like ?` clauses
// shared by every list/search tool, returning the (possibly extended) args slice.
func appendBranchMachine(sb *strings.Builder, args []any, branch, machine string) []any {
	if branch != "" {
		sb.WriteString(" and branch = ?")
		args = append(args, branch)
	}
	if machine != "" {
		sb.WriteString(" and machine like ?")
		args = append(args, machine)
	}
	return args
}

const defaultBranch = "master"

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
