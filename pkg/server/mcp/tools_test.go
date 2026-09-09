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
			{"perfintDev", "ide", "2026-05-01 10:00:00", uint32(1001), 12.5, uint16(261), uint16(27258), uint16(48), uint32(777)},
			{"perfintDev", "ide", "2026-05-01 09:00:00", uint32(1000), 12.0, uint16(261), uint16(27258), uint16(48), uint32(777)},
			// Dev Server row: no build components.
			{"perfintDev", "kotlin", "2026-05-01 08:00:00", uint32(999), 33.1, uint16(0), uint16(0), uint16(0), uint32(0)},
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
	if out.Groups[0].Rows[0].BuildNumber != "261.27258.48" || out.Groups[0].Rows[0].InstallerBuildID != 777 {
		t.Errorf("installer row lost its build number: %+v", out.Groups[0].Rows[0])
	}
	if out.Groups[1].Rows[0].BuildNumber != "" || out.Groups[1].Rows[0].InstallerBuildID != 0 {
		t.Errorf("dev-server row should carry no build number: %+v", out.Groups[1].Rows[0])
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
	// db_name, table_name, project_name, branch_name, machine_name, bld_time,
	// bc1, bc2, bc3, inst_id, installer_changes
	db.push(fakeQueryResult{
		rows: [][]any{
			{"perfintDev", "ide", "kotlin", "master", "linux-hetzner-1", "2026-05-01 10:00:00", uint16(261), uint16(27258), uint16(48), uint32(777), []string{newestEnc, oldestEnc}},
			// duplicate (db,table,project) — should be deduped.
			{"perfintDev", "ide", "kotlin", "master", "linux-hetzner-1", "2026-05-01 09:30:00", uint16(261), uint16(27258), uint16(48), uint32(777), []string{newestEnc, oldestEnc}},
			{"perfintDev", "ide", "spring", "master", "linux-hetzner-1", "2026-05-01 09:00:00", uint16(261), uint16(27258), uint16(48), uint32(777), []string{}},
		},
	})

	svc := newTestService(db, []tableRef{{Database: "perfintDev", Table: "ide", HasBuildTime: true, HasInstallerID: true, HasBuildComponents: true}})
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
	if out.BuildNumber != "261.27258.48" {
		t.Errorf("build_number = %q, want 261.27258.48", out.BuildNumber)
	}
	if out.InstallerBuildID != 777 {
		t.Errorf("tc_installer_build_id = %d, want 777", out.InstallerBuildID)
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
			{"perfintDev", "kotlin", "kotlin-proj", "master", "linux", "2026-05-01 10:00:00", uint16(0), uint16(0), uint16(0), uint32(0), []string{}},
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
	// Dev Server builds must report no build number rather than a bogus "0".
	if out.BuildNumber != "" || out.InstallerBuildID != 0 {
		t.Errorf("dev-server build should have no build number: %q / %d", out.BuildNumber, out.InstallerBuildID)
	}
}

// --- empty and unanswered results ----------------------------------------------------
//
// The rule these guard: a tool never answers with a shaped-but-blank result. Either it says why the
// answer is empty, or it fails. A silent blank is indistinguishable from "this does not exist", and
// the LLM above the MCP then reports a conclusion built on data it never received.

func TestGetBuild_UnknownBuildIsError(t *testing.T) {
	t.Parallel()
	db := &fakeDriver{}
	db.push(fakeQueryResult{rows: [][]any{}})

	svc := newTestService(db, []tableRef{{Database: "perfintDev", Table: "ide"}})
	cs := connectClient(t, svc)

	res := callTool(t, cs, "get_build", map[string]any{"tc_build_id": 4242}, nil)
	msg := errorText(t, res)
	if !strings.Contains(msg, "no data for tc_build_id=4242") || !strings.Contains(msg, "perfintDev.ide") {
		t.Errorf("unknown build must fail naming the id and the tables scanned, got %q", msg)
	}
}

func TestGetBuild_FailedCommitLookupIsReportedNotSwallowed(t *testing.T) {
	t.Parallel()
	db := &fakeDriver{}
	// A row with no installer commits, so the fallback lookup runs...
	db.push(fakeQueryResult{
		rows: [][]any{
			{"perfintDev", "kotlin", "kotlin-proj", "master", "linux", "2026-05-01 10:00:00", uint16(0), uint16(0), uint16(0), uint32(0), []string{}},
		},
	})
	// ...and fails. That must not read as "this build has no commits".
	db.push(fakeQueryResult{queryErr: errors.New("clickhouse: table installer doesn't exist")})

	svc := newTestService(db, []tableRef{{Database: "perfintDev", Table: "kotlin"}})
	cs := connectClient(t, svc)

	var out getBuildOutput
	res := callTool(t, cs, "get_build", map[string]any{"tc_build_id": 999}, &out)
	if res.IsError {
		t.Fatalf("unexpected error: %s", errorText(t, res))
	}
	if out.FirstCommit != "" {
		t.Errorf("first_commit = %q, want empty", out.FirstCommit)
	}
	joined := strings.Join(out.Notes, " | ")
	if !strings.Contains(joined, "unknown, not absent") {
		t.Errorf("a failed commit lookup must be reported in notes, got %q", joined)
	}
}

func TestGetBuild_FallbackFailureBeforeSuccessDoesNotLeaveWarning(t *testing.T) {
	t.Parallel()
	db := &fakeDriver{}
	db.push(fakeQueryResult{
		rows: [][]any{
			{"perfintDev", "kotlin", "kotlin-proj", "master", "linux", "2026-05-01 10:00:00", uint16(0), uint16(0), uint16(0), uint32(0), []string{}},
			{"perfint", "kotlin", "kotlin-proj", "master", "linux", "2026-05-01 09:00:00", uint16(0), uint16(0), uint16(0), uint32(0), []string{}},
		},
	})
	db.push(fakeQueryResult{queryErr: errors.New("installer lookup failed")})
	db.push(fakeQueryResult{rows: [][]any{{[]string{"newest-commit", "oldest-commit"}}}})

	svc := newTestService(db, []tableRef{{Database: "perfintDev", Table: "kotlin"}, {Database: "perfint", Table: "kotlin"}})
	cs := connectClient(t, svc)

	var out getBuildOutput
	res := callTool(t, cs, "get_build", map[string]any{"tc_build_id": 999}, &out)
	if res.IsError {
		t.Fatalf("unexpected error: %s", errorText(t, res))
	}
	if out.FirstCommit != "oldest-commit" || out.LastCommit != "newest-commit" {
		t.Fatalf("unexpected commit range: first=%q last=%q", out.FirstCommit, out.LastCommit)
	}
	if len(out.Notes) != 0 {
		t.Errorf("resolved commits must not retain lookup warnings, got %v", out.Notes)
	}
}

func TestGetBuild_LinkedInstallerWithoutCommitsIsNoted(t *testing.T) {
	t.Parallel()
	db := &fakeDriver{}
	db.push(fakeQueryResult{
		rows: [][]any{
			{"perfintDev", "ide", "kotlin-proj", "master", "linux", "2026-05-01 10:00:00", uint16(0), uint16(0), uint16(0), uint32(777), []string{}},
		},
	})
	db.push(fakeQueryResult{rows: [][]any{}})

	svc := newTestService(db, []tableRef{{Database: "perfintDev", Table: "ide", HasInstallerID: true}})
	cs := connectClient(t, svc)

	var out getBuildOutput
	res := callTool(t, cs, "get_build", map[string]any{"tc_build_id": 999}, &out)
	if res.IsError {
		t.Fatalf("unexpected error: %s", errorText(t, res))
	}
	if out.InstallerBuildID != 777 || out.FirstCommit != "" || out.LastCommit != "" {
		t.Fatalf("unexpected installer metadata: %+v", out)
	}
	joined := strings.Join(out.Notes, " | ")
	if !strings.Contains(joined, "no commit range available") || !strings.Contains(joined, "links installer 777") || strings.Contains(joined, "links no installer") {
		t.Errorf("missing commits must not imply an absent installer link, got %q", joined)
	}
}

func TestGetBuild_NoInstallerIsNoted(t *testing.T) {
	t.Parallel()
	db := &fakeDriver{}
	db.push(fakeQueryResult{
		rows: [][]any{
			{"perfintDev", "kotlin", "kotlin-proj", "master", "linux", "2026-05-01 10:00:00", uint16(0), uint16(0), uint16(0), uint32(0), []string{}},
		},
	})
	// The fallback finds no installer row: legitimately absent, but the caller must still be told.
	db.push(fakeQueryResult{rows: [][]any{}})

	svc := newTestService(db, []tableRef{{Database: "perfintDev", Table: "kotlin"}})
	cs := connectClient(t, svc)

	var out getBuildOutput
	callTool(t, cs, "get_build", map[string]any{"tc_build_id": 999}, &out)
	if !strings.Contains(strings.Join(out.Notes, " | "), "links no installer") {
		t.Errorf("an empty commit range must be explained, got notes %v", out.Notes)
	}
}

func TestSearchMetricValues_EmptyResultExplainsItself(t *testing.T) {
	t.Parallel()
	db := &fakeDriver{}
	db.push(fakeQueryResult{rows: [][]any{}})

	svc := newTestService(db, []tableRef{{Database: "perfintDev", Table: "ide"}})
	cs := connectClient(t, svc)

	var out searchMetricValuesOutput
	callTool(t, cs, "search_metric_values", map[string]any{
		"project":     "kotlin",
		"metric_name": "typo_total",
	}, &out)
	if out.Count != 0 || len(out.Notes) == 0 {
		t.Fatalf("empty result must carry notes: %+v", out)
	}
	joined := strings.Join(out.Notes, " | ")
	for _, want := range []string{`metric_name="typo_total"`, "perfintDev.ide", "search_metric_names"} {
		if !strings.Contains(joined, want) {
			t.Errorf("notes %q missing %q", joined, want)
		}
	}
}

func TestSearchMetricValues_EmptyResultIncludesMachineFilter(t *testing.T) {
	t.Parallel()
	db := &fakeDriver{}
	db.push(fakeQueryResult{
		verify: func(sql string, args []any) error {
			if !strings.Contains(sql, "and machine like ?") || !slices.Contains(args, "no-such-machine%") {
				return errors.New("expected machine filter in query")
			}
			return nil
		},
	})

	svc := newTestService(db, []tableRef{{Database: "perfintDev", Table: "ide"}})
	cs := connectClient(t, svc)

	var out searchMetricValuesOutput
	res := callTool(t, cs, "search_metric_values", map[string]any{
		"project":     "kotlin",
		"metric_name": "startup_total",
		"machine":     "no-such-machine%",
	}, &out)
	if res.IsError {
		t.Fatalf("unexpected error: %s", errorText(t, res))
	}
	joined := strings.Join(out.Notes, " | ")
	if out.Count != 0 || !strings.Contains(joined, `machine="no-such-machine%"`) {
		t.Errorf("empty result must include the machine filter, got %+v", out)
	}
}

func TestSearchMetricValues_LimitHitIsNoted(t *testing.T) {
	t.Parallel()
	db := &fakeDriver{}
	db.push(fakeQueryResult{
		rows: [][]any{
			{"perfintDev", "ide", "2026-05-01 10:00:00", uint32(1001), 12.5, uint16(0), uint16(0), uint16(0), uint32(0)},
		},
	})

	svc := newTestService(db, []tableRef{{Database: "perfintDev", Table: "ide"}})
	cs := connectClient(t, svc)

	var out searchMetricValuesOutput
	callTool(t, cs, "search_metric_values", map[string]any{
		"project":     "kotlin",
		"metric_name": "startup_total",
		"limit":       1,
	}, &out)
	if !strings.Contains(strings.Join(out.Notes, " | "), "TRUNCATED at limit=1") {
		t.Errorf("a limit-sized answer must warn it is partial, got notes %v", out.Notes)
	}
}

// A limit smaller than the window discards the OLD end, so the answer can cover days when the
// caller asked for months while still looking complete. Without this the caller computes a
// "long-run baseline" from a window that may lie entirely after the change under investigation.
func TestSearchMetricValues_TruncatedWindowNamesWhatIsMissing(t *testing.T) {
	t.Parallel()
	db := &fakeDriver{}
	db.push(fakeQueryResult{
		rows: [][]any{
			{"perfintDev", "ide", "2026-09-09 03:50:53", uint32(1002), 121.0, uint16(0), uint16(0), uint16(0), uint32(0)},
			{"perfintDev", "ide", "2026-09-07 10:00:00", uint32(1001), 378.0, uint16(0), uint16(0), uint16(0), uint32(0)},
			{"perfintDev", "ide", "2026-09-03 20:44:26", uint32(1000), 96.0, uint16(0), uint16(0), uint16(0), uint32(0)},
		},
	})

	svc := newTestService(db, []tableRef{{Database: "perfintDev", Table: "ide"}})
	cs := connectClient(t, svc)

	var out searchMetricValuesOutput
	callTool(t, cs, "search_metric_values", map[string]any{
		"project":     "spring_boot/showIntentions",
		"metric_name": "test#max_awt_delay",
		"days":        90,
		"limit":       3,
	}, &out)

	if want := "2026-09-03..2026-09-09 (7 of 90 requested days)"; out.Covered != want {
		t.Errorf("covered = %q, want %q", out.Covered, want)
	}
	notes := strings.Join(out.Notes, " | ")
	for _, want := range []string{"TRUNCATED at limit=3", "2026-09-03..2026-09-09", "older data exists", `aggregate="daily"`} {
		if !strings.Contains(notes, want) {
			t.Errorf("truncation note missing %q, got %v", want, out.Notes)
		}
	}
}

func TestSearchMetricValues_UntruncatedAnswerStillReportsItsSpan(t *testing.T) {
	t.Parallel()
	db := &fakeDriver{}
	db.push(fakeQueryResult{
		rows: [][]any{
			{"perfintDev", "ide", "2026-09-09 03:50:53", uint32(1001), 121.0, uint16(0), uint16(0), uint16(0), uint32(0)},
			{"perfintDev", "ide", "2026-09-08 03:50:53", uint32(1000), 378.0, uint16(0), uint16(0), uint16(0), uint32(0)},
		},
	})

	svc := newTestService(db, []tableRef{{Database: "perfintDev", Table: "ide"}})
	cs := connectClient(t, svc)

	var out searchMetricValuesOutput
	callTool(t, cs, "search_metric_values", map[string]any{
		"project":     "kotlin",
		"metric_name": "startup_total",
		"days":        30,
		"limit":       10,
	}, &out)

	if want := "2026-09-08..2026-09-09 (2 of 30 requested days)"; out.Covered != want {
		t.Errorf("covered = %q, want %q", out.Covered, want)
	}
	if strings.Contains(strings.Join(out.Notes, " | "), "TRUNCATED") {
		t.Errorf("a complete answer must not claim truncation, got notes %v", out.Notes)
	}
}

// aggregate=daily must return day buckets, aggregate over DISTINCT builds (a build can appear
// twice and would otherwise double-weight the median), and span the window instead of the
// most recent rows.
func TestSearchMetricValues_DailyAggregate(t *testing.T) {
	t.Parallel()
	db := &fakeDriver{}
	db.push(fakeQueryResult{
		verify: func(sql string, _ []any) error {
			for _, want := range []string{"toDate(gen_time)", "quantileExact(0.5)(value)", "select distinct", "group by db_name, table_name, day"} {
				if !strings.Contains(sql, want) {
					return fmt.Errorf("daily sql missing %q: %s", want, sql)
				}
			}
			return nil
		},
		rows: [][]any{
			{"perfintDev", "ide", "2026-08-28", uint32(19), 418.0, 72.0, 631.0},
			{"perfintDev", "ide", "2026-08-27", uint32(19), 351.0, 99.0, 704.0},
			{"perfintDev", "ide", "2026-08-26", uint32(19), 116.0, 91.0, 725.0},
		},
	})

	svc := newTestService(db, []tableRef{{Database: "perfintDev", Table: "ide"}})
	cs := connectClient(t, svc)

	var out searchMetricValuesOutput
	res := callTool(t, cs, "search_metric_values", map[string]any{
		"project":     "spring_boot/showIntentions",
		"metric_name": "test#max_awt_delay",
		"days":        60,
		"aggregate":   "daily",
	}, &out)
	if res.IsError {
		t.Fatalf("unexpected error: %s", errorText(t, res))
	}
	if out.Count != 3 || len(out.Groups) != 1 {
		t.Fatalf("count/groups = %d/%d, want 3/1", out.Count, len(out.Groups))
	}
	g := out.Groups[0]
	if len(g.Rows) != 0 {
		t.Errorf("daily answer must not carry per-build rows, got %d", len(g.Rows))
	}
	if len(g.Daily) != 3 {
		t.Fatalf("daily buckets = %d, want 3", len(g.Daily))
	}
	// The step this whole feature exists to make visible: median triples, floor does not move.
	if g.Daily[1].Date != "2026-08-27" || g.Daily[1].Median != 351 || g.Daily[1].Min != 99 || g.Daily[1].Builds != 19 {
		t.Errorf("bucket lost its shape: %+v", g.Daily[1])
	}
	if want := "2026-08-26..2026-08-28 (3 of 60 requested days)"; out.Covered != want {
		t.Errorf("covered = %q, want %q", out.Covered, want)
	}
}

func TestSearchMetricValues_RejectsUnknownAggregate(t *testing.T) {
	t.Parallel()
	svc := newTestService(&fakeDriver{}, []tableRef{{Database: "perfintDev", Table: "ide"}})
	cs := connectClient(t, svc)

	res := callTool(t, cs, "search_metric_values", map[string]any{
		"project":     "kotlin",
		"metric_name": "startup_total",
		"aggregate":   "weekly",
	}, nil)
	if got := errorText(t, res); !strings.Contains(got, `aggregate must be "daily"`) {
		t.Errorf("unknown aggregate must be rejected by name, got %q", got)
	}
}

func TestListProjects_EmptyResultExplainsItself(t *testing.T) {
	t.Parallel()
	db := &fakeDriver{}
	db.push(fakeQueryResult{rows: [][]any{}})

	svc := newTestService(db, []tableRef{{Database: "perfintDev", Table: "ide"}})
	cs := connectClient(t, svc)

	var out listProjectsOutput
	callTool(t, cs, "list_projects", map[string]any{"machine": "no-such-machine%"}, &out)
	joined := strings.Join(out.Notes, " | ")
	if !strings.Contains(joined, "no project matched") || !strings.Contains(joined, "no-such-machine%") {
		t.Errorf("notes must name the filters that matched nothing, got %q", joined)
	}
}

func TestSearchMetricNames_EmptyResultExplainsItself(t *testing.T) {
	t.Parallel()
	db := &fakeDriver{}
	db.push(fakeQueryResult{rows: [][]any{}})

	svc := newTestService(db, []tableRef{{Database: "perfintDev", Table: "ide"}})
	cs := connectClient(t, svc)

	var out searchMetricNamesOutput
	callTool(t, cs, "search_metric_names", map[string]any{"project": "no-such-project"}, &out)
	joined := strings.Join(out.Notes, " | ")
	if !strings.Contains(joined, "no metric name matched") || !strings.Contains(joined, "list_projects") {
		t.Errorf("notes must name the filters and the way to check them, got %q", joined)
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
			{"perfintDev", "ide", true, true, true},
			{"perfintDev", "kotlin", false, false, false},
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
	if len(got) != 2 || got[0].Table != "ide" || !got[0].HasInstallerID || !got[0].HasBuildComponents {
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
