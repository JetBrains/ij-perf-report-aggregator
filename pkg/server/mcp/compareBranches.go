package mcp

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"math"
	"slices"

	"github.com/JetBrains/ij-perf-report-aggregator/pkg/comparison"
	"github.com/JetBrains/ij-perf-report-aggregator/pkg/machine"
	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

type compareBranchesInput struct {
	Database     string   `json:"database"                jsonschema:"Database of the table to compare, e.g. perfUnitTests"`
	Table        string   `json:"table"                   jsonschema:"Table to compare, e.g. report"`
	Branch       string   `json:"branch"                  jsonschema:"Branch to compare against base_branch"`
	BaseBranch   string   `json:"base_branch,omitempty"   jsonschema:"Base branch. Default: master"`
	Metrics      []string `json:"metrics"                 jsonschema:"Metric names to compare (exact match in measures.name), e.g. attempt.mean.ms"`
	Machines     []string `json:"machines"                jsonschema:"Machine groups of the dashboard, e.g. linux-unit-perf-blade, or exact agent names. Not LIKE patterns."`
	Mode         string   `json:"mode,omitempty"          jsonschema:"Test mode. Default: \"default\", the runs without a mode. Ignored for tables without a mode column."`
	Projects     []string `json:"projects,omitempty"      jsonschema:"Optional exact project names to keep. Omitted: every project with values on both branches."`
	IncludeNoise bool     `json:"include_noise,omitempty" jsonschema:"Also return the rows whose change is not significant. Default: false"`
}

type compareBranchesRow struct {
	Project      string  `json:"project"`
	Metric       string  `json:"metric"`
	Machine      string  `json:"machine"            jsonschema:"Machine group the values come from; runs of different groups are never pooled"`
	BaseMedian   float64 `json:"base_median"`
	BranchMedian float64 `json:"branch_median"`
	DiffPercent  float64 `json:"diff_percent"       jsonschema:"branch_median / base_median - 1, in percent. For a time metric, a negative value is a gain."`
	DiffLower    float64 `json:"diff_lower_percent" jsonschema:"Lower bound of the shift of the branch against the base, in percent of base_median"`
	DiffUpper    float64 `json:"diff_upper_percent" jsonschema:"Upper bound of the shift of the branch against the base, in percent of base_median"`
	Significant  bool    `json:"significant"        jsonschema:"The bounds exclude 0, so the change is not noise"`
}

type compareBranchesOutput struct {
	Rows        []compareBranchesRow `json:"rows"            jsonschema:"Significant rows first, then by the absolute diff_percent"`
	Compared    int                  `json:"compared"        jsonschema:"Rows with values on both branches, before the noise filter"`
	Significant int                  `json:"significant"     jsonschema:"Rows whose change is not noise"`
	Notes       []string             `json:"notes,omitempty" jsonschema:"Read these before drawing conclusions: they say why the result is empty. Missing when the answer is complete."`
}

func (s *service) compareBranches(ctx context.Context, _ *sdk.CallToolRequest, in compareBranchesInput) (*sdk.CallToolResult, compareBranchesOutput, error) {
	if in.Database == "" || in.Table == "" {
		return nil, compareBranchesOutput{}, errors.New("database and table are required")
	}
	if in.Branch == "" {
		return nil, compareBranchesOutput{}, errors.New("branch is required")
	}
	if len(in.Metrics) == 0 {
		return nil, compareBranchesOutput{}, errors.New("metrics are required")
	}
	if len(in.Machines) == 0 {
		return nil, compareBranchesOutput{}, errors.New("machines are required")
	}
	base := cmp.Or(in.BaseBranch, defaultBranch)
	if base == in.Branch {
		return nil, compareBranchesOutput{}, fmt.Errorf("branch and base_branch are both %q", base)
	}
	tables, err := s.resolveTables(ctx, in.Database, in.Table)
	if err != nil {
		return nil, compareBranchesOutput{}, err
	}
	table := tables[0]

	where := "branch IN ? AND measure_name IN ? AND (" + machine.Condition(in.Machines) + ")"
	args := []any{[]string{base, in.Branch}, in.Metrics}
	if table.HasMode {
		where += " AND mode = ?"
		args = append(args, modeValue(in.Mode))
	}
	if len(in.Projects) > 0 {
		where += " AND project IN ?"
		args = append(args, in.Projects)
	}

	rows, err := s.db.Query(ctx, comparison.SQL("branch", table.Database+"."+table.Table, where), args...)
	if err != nil {
		return nil, compareBranchesOutput{}, fmt.Errorf("compare_branches: %w", err)
	}
	defer rows.Close()
	var series []comparison.Series
	for rows.Next() {
		var r comparison.Series
		if err := rows.Scan(&r.Branch, &r.Project, &r.MeasureName, &r.Machine, &r.MeasureValues); err != nil {
			return nil, compareBranchesOutput{}, fmt.Errorf("scan: %w", err)
		}
		series = append(series, r)
	}
	if err := rows.Err(); err != nil {
		return nil, compareBranchesOutput{}, fmt.Errorf("rows: %w", err)
	}

	compared := comparison.Compare(series, base, in.Branch)
	out := compareBranchesOutput{Rows: []compareBranchesRow{}, Compared: len(compared)}
	for _, r := range compared {
		if r.Significant {
			out.Significant++
		} else if !in.IncludeNoise {
			continue
		}
		out.Rows = append(out.Rows, compareBranchesRow{
			Project: r.Project, Metric: r.MeasureName, Machine: r.Machine,
			BaseMedian: r.Median1, BranchMedian: r.Median2,
			DiffPercent: r.Diff, DiffLower: r.DiffLower, DiffUpper: r.DiffUpper,
			Significant: r.Significant,
		})
	}
	slices.SortFunc(out.Rows, func(a, b compareBranchesRow) int {
		if a.Significant != b.Significant {
			if a.Significant {
				return -1
			}
			return 1
		}
		return cmp.Or(cmp.Compare(math.Abs(b.DiffPercent), math.Abs(a.DiffPercent)), cmp.Compare(a.Project, b.Project),
			cmp.Compare(a.Metric, b.Metric), cmp.Compare(a.Machine, b.Machine))
	})

	if len(compared) == 0 {
		out.Notes = append(out.Notes, fmt.Sprintf(
			"no project has values on both %q and %q in %s.%s for these metrics and machines in the last month",
			in.Branch, base, table.Database, table.Table))
	}
	return nil, out, nil
}
