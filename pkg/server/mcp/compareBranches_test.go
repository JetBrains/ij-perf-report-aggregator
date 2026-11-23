package mcp

import (
	"errors"
	"maps"
	"reflect"
	"slices"
	"strings"
	"testing"
)

var unitReport = tableRef{Database: "perfUnitTests", Table: "report", HasMode: true}

// compareSeriesRow is one row of comparison.SQL: branch, project, metric, machine group, newest values.
func compareSeriesRow(branch, project string, values ...int) []any {
	return []any{branch, project, "attempt.mean.ms", "linux-unit-perf-blade", values}
}

func compareArgs(extra map[string]any) map[string]any {
	args := map[string]any{
		"database": "perfUnitTests",
		"table":    "report",
		"branch":   "my-branch",
		"metrics":  []string{"attempt.mean.ms"},
		"machines": []string{"linux-unit-perf-blade"},
	}
	maps.Copy(args, extra)
	return args
}

func TestCompareBranches_RequiresArguments(t *testing.T) {
	t.Parallel()
	cs := newTestClient(t, []tableRef{unitReport})
	for _, missing := range []string{"table", "branch", "metrics", "machines"} {
		args := compareArgs(nil)
		delete(args, missing)
		if res := callTool(t, cs, "compare_branches", args, nil); !res.IsError {
			t.Errorf("a call without %s must fail", missing)
		}
	}
}

func TestCompareBranches_RejectsTheBaseAsBranch(t *testing.T) {
	t.Parallel()
	cs := newTestClient(t, []tableRef{unitReport})
	res := callTool(t, cs, "compare_branches", compareArgs(map[string]any{"branch": "master"}), nil)
	if msg := errorText(t, res); !strings.Contains(msg, "both") {
		t.Errorf("unexpected error: %s", msg)
	}
}

func TestCompareBranches_BindsTheFilters(t *testing.T) {
	t.Parallel()
	cs := newTestClient(t, []tableRef{unitReport}, fakeQueryResult{
		verify: func(sql string, args []any) error {
			if !strings.Contains(sql, "FROM perfUnitTests.report") {
				return errors.New("expected the requested table")
			}
			if !strings.Contains(sql, "machine_group IN ('linux-unit-perf-blade')") {
				return errors.New("expected the machine group condition")
			}
			want := []any{[]string{"master", "my-branch"}, []string{"attempt.mean.ms"}, "", []string{"a.Test.test - attempt.mean.ms"}}
			if !reflect.DeepEqual(args, want) {
				return errors.New("unexpected args")
			}
			return nil
		},
	})

	var out compareBranchesOutput
	callToolOK(t, cs, "compare_branches", compareArgs(map[string]any{"projects": []string{"a.Test.test - attempt.mean.ms"}}), &out)
	if out.Compared != 0 || !strings.Contains(joinNotes(out.Notes), `"my-branch"`) {
		t.Errorf("an empty answer must explain itself, got %+v", out)
	}
}

func TestCompareBranches_NoModeFilterWithoutAModeColumn(t *testing.T) {
	t.Parallel()
	noMode := tableRef{Database: "perfUnitTests", Table: "report"}
	cs := newTestClient(t, []tableRef{noMode}, fakeQueryResult{
		verify: func(sql string, _ []any) error {
			if strings.Contains(sql, "mode = ?") {
				return errors.New("a table without a mode column must not filter by mode")
			}
			return nil
		},
	})
	callToolOK(t, cs, "compare_branches", compareArgs(nil), &compareBranchesOutput{})
}

func TestCompareBranches_KeepsOnlySignificantRows(t *testing.T) {
	t.Parallel()
	stable := []int{100, 101, 99, 100, 102, 98, 100, 101, 99, 100}
	slower := []int{150, 151, 149, 150, 152, 148, 150, 151, 149, 150}
	rows := [][]any{
		compareSeriesRow("master", "slower", stable...),
		compareSeriesRow("my-branch", "slower", slower...),
		compareSeriesRow("master", "same", stable...),
		compareSeriesRow("my-branch", "same", stable...),
		compareSeriesRow("master", "only-on-master", stable...),
	}
	cs := newTestClient(t, []tableRef{unitReport}, fakeQueryResult{rows: rows}, fakeQueryResult{rows: rows})

	var out compareBranchesOutput
	callToolOK(t, cs, "compare_branches", compareArgs(nil), &out)
	if out.Compared != 2 || out.Significant != 1 || len(out.Rows) != 1 {
		t.Fatalf("expected 1 significant row of 2 compared, got %+v", out)
	}
	row := out.Rows[0]
	if row.Project != "slower" || !row.Significant || row.DiffPercent != 50 || row.DiffLower <= 0 {
		t.Errorf("unexpected row: %+v", row)
	}

	var all compareBranchesOutput
	callToolOK(t, cs, "compare_branches", compareArgs(map[string]any{"include_noise": true}), &all)
	projects := make([]string, 0, len(all.Rows))
	for _, r := range all.Rows {
		projects = append(projects, r.Project)
	}
	if !slices.Equal(projects, []string{"slower", "same"}) {
		t.Errorf("include_noise must return all rows, significant first, got %v", projects)
	}
}
