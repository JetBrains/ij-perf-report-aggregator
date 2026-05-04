package mcp

import (
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strings"
	"testing"
	"time"

	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

// callTool invokes the named MCP tool with the given args, decoding StructuredContent into out.
// Fails the test on transport errors. Returns the raw result so callers can also assert on IsError/Content.
func callTool(t *testing.T, cs *sdk.ClientSession, name string, args map[string]any, out any) *sdk.CallToolResult {
	t.Helper()
	res, err := cs.CallTool(t.Context(), &sdk.CallToolParams{Name: name, Arguments: args})
	if err != nil {
		t.Fatalf("CallTool(%s): transport error: %v", name, err)
	}
	if out != nil && res.StructuredContent != nil {
		raw, err := json.Marshal(res.StructuredContent)
		if err != nil {
			t.Fatalf("marshal StructuredContent: %v", err)
		}
		if err := json.Unmarshal(raw, out); err != nil {
			t.Fatalf("unmarshal StructuredContent into %T: %v\nraw: %s", out, err, raw)
		}
	}
	return res
}

// errorText returns the concatenated text content of a tool result reported with IsError=true.
func errorText(t *testing.T, res *sdk.CallToolResult) string {
	t.Helper()
	if !res.IsError {
		t.Fatalf("expected IsError, got success: %+v", res)
	}
	var sb strings.Builder
	for _, c := range res.Content {
		if tc, ok := c.(*sdk.TextContent); ok {
			sb.WriteString(tc.Text)
		}
	}
	return sb.String()
}

// --- list_projects -------------------------------------------------------------------

func TestListProjects_HappyPath(t *testing.T) {
	t.Parallel()
	db := &fakeDriver{}
	db.push(fakeQueryResult{
		rows: [][]any{
			{"perfintDev", "ide", "kotlin"},
			{"perfintDev", "ide", "spring"},
		},
		verify: func(sql string, args []any) error {
			if !strings.Contains(sql, "subtractDays(now(), ?)") {
				return errors.New("expected lookback window in SQL")
			}
			if !strings.Contains(sql, "and branch = ?") {
				return errors.New("expected branch filter")
			}
			if strings.Contains(sql, "and machine like ?") {
				return errors.New("did not expect machine filter (none supplied)")
			}
			if got, want := strings.Count(sql, "?"), len(args); got != want {
				return fmt.Errorf("placeholder/arg count mismatch: %d ? vs %d args", got, want)
			}
			if !slices.Contains(args, "master") {
				return errors.New("expected branch=master to appear in args")
			}
			return nil
		},
	})

	svc := newTestService(db, []tableRef{{Database: "perfintDev", Table: "ide"}})
	cs := connectClient(t, svc)

	var out listProjectsOutput
	res := callTool(t, cs, "list_projects", map[string]any{"branch": "master"}, &out)
	if res.IsError {
		t.Fatalf("unexpected error: %s", errorText(t, res))
	}
	if out.Count != 2 || len(out.Rows) != 2 {
		t.Fatalf("count=%d rows=%v", out.Count, out.Rows)
	}
	if out.Rows[0].Project != "kotlin" || out.Rows[1].Project != "spring" {
		t.Errorf("unexpected projects: %+v", out.Rows)
	}
}

func TestListProjects_BranchDefaultsToMaster(t *testing.T) {
	t.Parallel()
	db := &fakeDriver{}
	db.push(fakeQueryResult{
		verify: func(_ string, args []any) error {
			if !slices.Contains(args, "master") {
				return errors.New("expected branch=master to appear in args")
			}
			return nil
		},
	})

	svc := newTestService(db, []tableRef{{Database: "d", Table: "t"}})
	cs := connectClient(t, svc)

	var out listProjectsOutput
	res := callTool(t, cs, "list_projects", map[string]any{}, &out)
	if res.IsError {
		t.Fatalf("unexpected error: %s", errorText(t, res))
	}
}

// --- search_metric_names -------------------------------------------------------------

func TestSearchMetricNames_RequiresProject(t *testing.T) {
	t.Parallel()
	svc := newTestService(&fakeDriver{}, []tableRef{{Database: "d", Table: "t"}})
	cs := connectClient(t, svc)

	// Missing key → SDK-level schema validation rejects it before the handler runs.
	res := callTool(t, cs, "search_metric_names", map[string]any{}, nil)
	if msg := errorText(t, res); !strings.Contains(msg, "project") {
		t.Errorf("expected error mentioning project, got %q", msg)
	}

	// Empty value → handler's manual check fires.
	res = callTool(t, cs, "search_metric_names", map[string]any{"project": ""}, nil)
	if msg := errorText(t, res); !strings.Contains(msg, "project is required") {
		t.Errorf("expected 'project is required', got %q", msg)
	}
}

func TestSearchMetricNames_AppliesFilters(t *testing.T) {
	t.Parallel()
	db := &fakeDriver{}
	db.push(fakeQueryResult{
		rows: [][]any{{"perfintDev", "ide", "startup_total"}},
		verify: func(sql string, args []any) error {
			if !strings.Contains(sql, "metric_name like ?") {
				return errors.New("expected name_pattern LIKE clause")
			}
			if !slices.Contains(args, "startup%") {
				return errors.New("expected name_pattern in args")
			}
			if !slices.Contains(args, 500) {
				return errors.New("expected default limit 500 in args")
			}
			return nil
		},
	})

	svc := newTestService(db, []tableRef{{Database: "perfintDev", Table: "ide"}})
	cs := connectClient(t, svc)

	var out searchMetricNamesOutput
	res := callTool(t, cs, "search_metric_names", map[string]any{
		"project":      "kotlin",
		"name_pattern": "startup%",
	}, &out)
	if res.IsError {
		t.Fatalf("unexpected error: %s", errorText(t, res))
	}
	if out.Count != 1 || out.Rows[0].Name != "startup_total" {
		t.Errorf("unexpected output: %+v", out)
	}
}

// --- search_metric_values ------------------------------------------------------------

func TestSearchMetricValues_RequiresProjectAndMetric(t *testing.T) {
	t.Parallel()
	svc := newTestService(&fakeDriver{}, []tableRef{{Database: "d", Table: "t"}})
	cs := connectClient(t, svc)

	// Empty values → handler's manual check fires (schema only enforces presence).
	res := callTool(t, cs, "search_metric_values", map[string]any{"project": "", "metric_name": "x"}, nil)
	if msg := errorText(t, res); !strings.Contains(msg, "project is required") {
		t.Errorf("got %q", msg)
	}

	res = callTool(t, cs, "search_metric_values", map[string]any{"project": "kotlin", "metric_name": ""}, nil)
	if msg := errorText(t, res); !strings.Contains(msg, "metric_name is required") {
		t.Errorf("got %q", msg)
	}
}

func TestSearchMetricValues_GroupsByDatabaseTable(t *testing.T) {
	t.Parallel()
	db := &fakeDriver{}
	// Rows arrive ordered by gen_time desc; groups are formed in arrival order.
	db.push(fakeQueryResult{
		rows: [][]any{
			{"perfintDev", "ide", "2026-05-01 10:00:00", uint32(1001), 12.5},
			{"perfintDev", "ide", "2026-05-01 09:00:00", uint32(1000), 12.0},
			{"perfintDev", "kotlin", "2026-05-01 08:00:00", uint32(999), 33.1},
		},
	})

	svc := newTestService(db, []tableRef{
		{Database: "perfintDev", Table: "ide"},
		{Database: "perfintDev", Table: "kotlin"},
	})
	cs := connectClient(t, svc)

	var out searchMetricValuesOutput
	res := callTool(t, cs, "search_metric_values", map[string]any{
		"project":     "kotlin",
		"metric_name": "startup_total",
	}, &out)
	if res.IsError {
		t.Fatalf("unexpected error: %s", errorText(t, res))
	}
	if out.Count != 3 {
		t.Errorf("count = %d, want 3", out.Count)
	}
	if len(out.Groups) != 2 {
		t.Fatalf("groups = %d, want 2", len(out.Groups))
	}
	if out.Groups[0].Table != "ide" || len(out.Groups[0].Rows) != 2 {
		t.Errorf("group[0] = %+v", out.Groups[0])
	}
	if out.Groups[1].Table != "kotlin" || len(out.Groups[1].Rows) != 1 {
		t.Errorf("group[1] = %+v", out.Groups[1])
	}
	if out.Project != "kotlin" || out.MetricName != "startup_total" || out.Branch != "master" {
		t.Errorf("root metadata wrong: %+v", out)
	}
}

// --- get_build -----------------------------------------------------------------------

func TestGetBuild_RequiresPositiveID(t *testing.T) {
	t.Parallel()
	svc := newTestService(&fakeDriver{}, []tableRef{{Database: "d", Table: "t"}})
	cs := connectClient(t, svc)

	for _, id := range []int{0, -1} {
		res := callTool(t, cs, "get_build", map[string]any{"tc_build_id": id}, nil)
		if !strings.Contains(errorText(t, res), "tc_build_id is required") {
			t.Errorf("id=%d: got %q", id, errorText(t, res))
		}
	}
}

func TestGetBuild_HappyPath_WithInstallerCommits(t *testing.T) {
	t.Parallel()

	const newestSHA = "22596363b3de40b06f981fb85d82312e8c0ed511"
	const oldestSHA = "0123456789abcdef0123456789abcdef01234567"
	newestEnc := encodeSHA1(t, newestSHA)
	oldestEnc := encodeSHA1(t, oldestSHA)

	db := &fakeDriver{}
	// Single row from the union; fields per getBuild SQL:
	// db_name, table_name, project_name, branch_name, machine_name, bld_time, installer_changes
	db.push(fakeQueryResult{
		rows: [][]any{
			{"perfintDev", "ide", "kotlin", "master", "linux-hetzner-1", "2026-05-01 10:00:00", []string{newestEnc, oldestEnc}},
			// duplicate (db,table,project) — should be deduped.
			{"perfintDev", "ide", "kotlin", "master", "linux-hetzner-1", "2026-05-01 09:30:00", []string{newestEnc, oldestEnc}},
			{"perfintDev", "ide", "spring", "master", "linux-hetzner-1", "2026-05-01 09:00:00", []string{}},
		},
	})

	svc := newTestService(db, []tableRef{{Database: "perfintDev", Table: "ide", HasBuildTime: true, HasInstallerID: true}})
	cs := connectClient(t, svc)

	var out getBuildOutput
	res := callTool(t, cs, "get_build", map[string]any{"tc_build_id": 12345}, &out)
	if res.IsError {
		t.Fatalf("unexpected error: %s", errorText(t, res))
	}
	if out.BuildID != 12345 {
		t.Errorf("BuildID = %d", out.BuildID)
	}
	if out.Branch != "master" || out.Machine != "linux-hetzner-1" {
		t.Errorf("metadata: %+v", out)
	}
	if !strings.HasPrefix(out.TeamCityURL, teamCityBaseURL) || !strings.Contains(out.TeamCityURL, "12345") {
		t.Errorf("teamcity_url = %q", out.TeamCityURL)
	}
	if out.FirstCommit != oldestSHA[:shortCommitLen] || out.LastCommit != newestSHA[:shortCommitLen] {
		t.Errorf("commits: first=%q last=%q", out.FirstCommit, out.LastCommit)
	}
	if out.Count != 2 {
		t.Errorf("count = %d, want 2 (kotlin + spring, dedup of duplicate kotlin row)", out.Count)
	}
	if len(out.Projects) != 2 {
		t.Fatalf("projects = %d", len(out.Projects))
	}
}

func TestGetBuild_InstallerFallback(t *testing.T) {
	t.Parallel()

	const fallbackSHA = "abcdef1234567890abcdef1234567890abcdef12"
	fallbackEnc := encodeSHA1(t, fallbackSHA)

	db := &fakeDriver{}
	// First call: union query returns one row but with no installer_changes (table has no installer column).
	db.push(fakeQueryResult{
		rows: [][]any{
			{"perfintDev", "kotlin", "kotlin-proj", "master", "linux", "2026-05-01 10:00:00", []string{}},
		},
	})
	// Second call: the fallback `select ... from perfintDev.installer where id = ?` lookup.
	db.push(fakeQueryResult{
		rows: [][]any{{[]string{fallbackEnc}}},
	})

	// Note HasInstallerID=false so the SQL doesn't try to join installer in the union.
	svc := newTestService(db, []tableRef{{Database: "perfintDev", Table: "kotlin"}})
	cs := connectClient(t, svc)

	var out getBuildOutput
	res := callTool(t, cs, "get_build", map[string]any{"tc_build_id": 999}, &out)
	if res.IsError {
		t.Fatalf("unexpected error: %s", errorText(t, res))
	}
	if out.FirstCommit != fallbackSHA[:shortCommitLen] || out.LastCommit != fallbackSHA[:shortCommitLen] {
		t.Errorf("fallback commits not applied: first=%q last=%q", out.FirstCommit, out.LastCommit)
	}
}

// --- list_tables ---------------------------------------------------------------------

func TestListTablesTool_UsesCache(t *testing.T) {
	t.Parallel()
	db := &fakeDriver{}
	// No expectations queued: any query would fail. This proves the cache is hit.

	svc := newTestService(db, []tableRef{
		{Database: "perfintDev", Table: "ide"},
		{Database: "perfintDev", Table: "kotlin"},
	})
	cs := connectClient(t, svc)

	var out listTablesOutput
	res := callTool(t, cs, "list_tables", nil, &out)
	if res.IsError {
		t.Fatalf("unexpected error: %s", errorText(t, res))
	}
	if out.Count != 2 {
		t.Errorf("count = %d", out.Count)
	}
	if len(db.calls) != 0 {
		t.Errorf("expected 0 db calls (cache hit), got %d", len(db.calls))
	}
}

func TestListTables_RefreshesAfterTTL(t *testing.T) {
	t.Parallel()
	db := &fakeDriver{}
	db.push(fakeQueryResult{
		rows: [][]any{
			{"perfintDev", "ide", true, true},
			{"perfintDev", "kotlin", false, false},
		},
	})

	svc := newService(db)
	// Force a stale cache so the next call refreshes via the (queued) DB result.
	svc.tablesCache = []tableRef{{Database: "stale", Table: "stale"}}
	svc.tablesCached = time.Now().Add(-2 * tablesTTL)

	got, err := svc.listTables(t.Context())
	if err != nil {
		t.Fatalf("listTables: %v", err)
	}
	if len(got) != 2 || got[0].Table != "ide" || !got[0].HasInstallerID {
		t.Errorf("refreshed result wrong: %+v", got)
	}
	if len(db.calls) != 1 {
		t.Errorf("expected exactly 1 refresh call, got %d", len(db.calls))
	}

	// Second call within TTL must not re-query.
	if _, err := svc.listTables(t.Context()); err != nil {
		t.Fatalf("listTables (cached): %v", err)
	}
	if len(db.calls) != 1 {
		t.Errorf("cached call triggered another query (calls=%d)", len(db.calls))
	}
}

// --- resolveTables -------------------------------------------------------------------

func TestResolveTables_RejectsInvalidIdentifier(t *testing.T) {
	t.Parallel()
	svc := newTestService(&fakeDriver{}, []tableRef{{Database: "d", Table: "t"}})

	if _, err := svc.resolveTables(t.Context(), "ok; drop", ""); err == nil {
		t.Errorf("expected validation error for bad database identifier")
	}
	if _, err := svc.resolveTables(t.Context(), "", "ok-bad"); err == nil {
		t.Errorf("expected validation error for bad table identifier")
	}
}

func TestResolveTables_NoMatch(t *testing.T) {
	t.Parallel()
	svc := newTestService(&fakeDriver{}, []tableRef{{Database: "perfintDev", Table: "ide"}})

	_, err := svc.resolveTables(t.Context(), "perfintDev", "missing")
	if err == nil || !strings.Contains(err.Error(), "no known table matches") {
		t.Errorf("expected 'no known table matches' error, got: %v", err)
	}
}

func TestResolveTables_FilterByDatabase(t *testing.T) {
	t.Parallel()
	svc := newTestService(&fakeDriver{}, []tableRef{
		{Database: "perfintDev", Table: "ide"},
		{Database: "perfintDev", Table: "kotlin"},
		{Database: "other", Table: "x"},
	})

	got, err := svc.resolveTables(t.Context(), "perfintDev", "")
	if err != nil {
		t.Fatalf("resolveTables: %v", err)
	}
	if len(got) != 2 {
		t.Errorf("len = %d, want 2", len(got))
	}
}
