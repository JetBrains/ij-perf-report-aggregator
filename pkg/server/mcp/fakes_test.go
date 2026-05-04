package mcp

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"sync"
	"testing"
	"time"

	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

// fakeDriver is a chConn that returns canned rows in FIFO order from `queue`.
// Each tool call typically pops one expectation. For tests that exercise the
// table-discovery path, prepend an extra expectation for the listTables query.
type fakeDriver struct {
	mu    sync.Mutex
	queue []fakeQueryResult
	calls []recordedCall
}

type fakeQueryResult struct {
	queryErr error                     // returned from Query()
	rowsErr  error                     // returned by rows.Err() at end of iteration
	rows     [][]any                   // each inner slice is one row's column values
	verify   func(string, []any) error // optional assertion on sql + args
}

type recordedCall struct {
	sql  string
	args []any
}

func (f *fakeDriver) Query(_ context.Context, sql string, args ...any) (chRows, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.calls = append(f.calls, recordedCall{sql: sql, args: args})
	if len(f.queue) == 0 {
		return nil, fmt.Errorf("fakeDriver: unexpected query (queue empty): %s", sql)
	}
	r := f.queue[0]
	f.queue = f.queue[1:]
	if r.verify != nil {
		if err := r.verify(sql, args); err != nil {
			return nil, fmt.Errorf("fakeDriver: verify failed: %w", err)
		}
	}
	if r.queryErr != nil {
		return nil, r.queryErr
	}
	return &fakeRows{rows: r.rows, err: r.rowsErr}, nil
}

func (f *fakeDriver) Close() error { return nil }

func (f *fakeDriver) push(r fakeQueryResult) { f.queue = append(f.queue, r) }

// fakeRows is a chRows backed by a [][]any. Scan copies via reflection so callers
// can pass *string/*int64/*[]string/etc. and have the values land in place.
type fakeRows struct {
	rows [][]any
	idx  int
	err  error
}

func (r *fakeRows) Next() bool {
	if r.idx >= len(r.rows) {
		return false
	}
	r.idx++
	return true
}

func (r *fakeRows) Scan(dest ...any) error {
	if r.idx == 0 {
		return errors.New("fakeRows.Scan called before Next")
	}
	row := r.rows[r.idx-1]
	if len(dest) != len(row) {
		return fmt.Errorf("fakeRows.Scan: %d dest, %d cols", len(dest), len(row))
	}
	for i, d := range dest {
		dv := reflect.ValueOf(d)
		if dv.Kind() != reflect.Pointer || dv.IsNil() {
			return fmt.Errorf("fakeRows.Scan: dest[%d] not a non-nil pointer", i)
		}
		sv := reflect.ValueOf(row[i])
		dt := dv.Elem().Type()
		if !sv.Type().ConvertibleTo(dt) {
			return fmt.Errorf("fakeRows.Scan: col %d type %s not convertible to %s", i, sv.Type(), dt)
		}
		dv.Elem().Set(sv.Convert(dt))
	}
	return nil
}

func (r *fakeRows) Err() error { return r.err }

func (r *fakeRows) Close() error { return nil }

// newTestService returns a service whose tablesCache is pre-seeded with the given
// tables. The discovery query is therefore skipped, so `db` only needs canned
// results for the actual tool query.
func newTestService(db chConn, tables []tableRef) *service {
	s := newService(db)
	s.tablesCache = tables
	s.tablesCached = time.Now()
	return s
}

// connectClient builds an in-memory MCP server from svc, connects a fresh client,
// and returns the client session for use in `cs.CallTool(...)`.
func connectClient(t *testing.T, svc *service) *sdk.ClientSession {
	t.Helper()
	ctx := t.Context()

	server := svc.buildServer()
	st, ct := sdk.NewInMemoryTransports()
	if _, err := server.Connect(ctx, st, nil); err != nil {
		t.Fatalf("server connect: %v", err)
	}
	client := sdk.NewClient(&sdk.Implementation{Name: "test", Version: "v0"}, nil)
	cs, err := client.Connect(ctx, ct, nil)
	if err != nil {
		t.Fatalf("client connect: %v", err)
	}
	t.Cleanup(func() { _ = cs.Close() })
	return cs
}
