package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func ownerResult(branch, project, metric, machine string, values []int) ownerQueryResult {
	return ownerQueryResult{
		measureQueryResult: measureQueryResult{Branch: branch, Project: project, MeasureName: metric, MeasureValues: values},
		Machine:            machine,
	}
}

// TestBuildComparisonResponseComparesPerMachine verifies the AT-4930 fix: values from
// distinct agents matching the same machine filter must be compared per-agent and never
// pooled together. Here each agent is stable across branches (no real regression), yet the
// fast and slow agents happen to be split unevenly between branches. Pooling would report a
// huge fake degradation; per-machine comparison must report ~0 for each agent instead.
func TestBuildComparisonResponseComparesPerMachine(t *testing.T) {
	t.Parallel()

	const project = "radler/curl/indexing"
	const metric = "indexingTimeWithoutPauses"
	const fastMachine = "Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)"
	const slowMachine = "Linux Blade (8 vCPU, 16 GB)"

	fast := []int{100, 102, 98, 101, 99, 100}
	slow := []int{400, 405, 398, 402, 401, 399}

	results := []ownerQueryResult{
		ownerResult("261", project, metric, fastMachine, fast),
		ownerResult("master", project, metric, fastMachine, fast),
		ownerResult("261", project, metric, slowMachine, slow),
		ownerResult("master", project, metric, slowMachine, slow),
	}

	dbTableMap := map[string]dbTableKey{project: {DbName: "perfintDev", TableName: "clion"}}

	response := buildComparisonResponse(results, dbTableMap, "261", "master")

	// One row per machine, each comparing the agent against itself.
	require.Len(t, response, 2)

	byMachine := make(map[string]comparisonResponseItem, len(response))
	for _, item := range response {
		byMachine[item.Machine] = item
	}

	require.Contains(t, byMachine, fastMachine)
	require.Contains(t, byMachine, slowMachine)

	for machine, item := range byMachine {
		assert.InDeltaf(t, 0, item.Diff, 1.0, "unexpected diff for machine %q", machine)
		assert.Equalf(t, project, item.Project, "project for machine %q", machine)
		assert.Equalf(t, metric, item.Metric, "metric for machine %q", machine)
	}

	// The link must carry the exact url-encoded agent (spaces -> '+', parens escaped).
	assert.Contains(t, byMachine[fastMachine].Link, "machine=Linux+EC2+C6id.8xlarge+%2832+vCPU+Xeon%2C+64+GB%29")
	assert.Contains(t, byMachine[slowMachine].Link, "machine=Linux+Blade+%288+vCPU%2C+16+GB%29")

	// Sanity: the fast agent's median must stay ~100, not be dragged toward the slow agent's ~400.
	assert.InDelta(t, 100, byMachine[fastMachine].BaseBranchValue, 5)
	assert.InDelta(t, 400, byMachine[slowMachine].BaseBranchValue, 10)
}

// TestBuildComparisonResponseSkipsSingleAgentBranch verifies that when a metric was recorded
// on an agent in only one of the two branches, no (misleading) comparison row is produced for
// that agent.
func TestBuildComparisonResponseSkipsSingleAgentBranch(t *testing.T) {
	t.Parallel()

	const project = "p"
	const metric = "m"
	const machine = "Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)"

	results := []ownerQueryResult{
		ownerResult("261", project, metric, machine, []int{100, 101, 99}),
		// no "master" data on this machine
	}

	response := buildComparisonResponse(results, map[string]dbTableKey{project: {DbName: "db", TableName: "t"}}, "261", "master")
	assert.Empty(t, response)
}
