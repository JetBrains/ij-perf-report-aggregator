package analyzer

import (
	"os"
	"testing"
	"time"

	"github.com/JetBrains/ij-perf-report-aggregator/pkg/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFilter(t *testing.T) {
	t.Parallel()
	data, err := os.ReadFile("../../testData/opentelemetry.json")
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	testCases := map[string]int32{
		"project.opening":   5035,
		"globalInspections": 141561,
	}

	got := analyzeOtJson(data, []string{"globalInspections", "project.opening"})

	for name, expected := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assert.Contains(t, got, name)
			assert.Equal(t, expected, got[name][0])
		})
	}
}

func TestFilterNewOpenGrepSpans(t *testing.T) {
	t.Parallel()
	data, err := os.ReadFile("../../testData/opentelemetry.json")
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	testCases := map[string]int32{
		"OpenGrepGlobalInspection":         500,
		"Building IR":                      300,
		"Collecting Opengrep AST/Problems": 100,
		"Building IR library symbols":      50,
		"Building IR project":              200,
		"Running RML + IFDS":               150,
	}

	operationNames := make([]string, 0, len(testCases))
	for name := range testCases {
		operationNames = append(operationNames, name)
	}

	got := analyzeOtJson(data, operationNames)

	for name, expected := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assert.Contains(t, got, name)
			assert.Len(t, got[name], 1)
			assert.Equal(t, expected, got[name][0])
		})
	}
}

func TestAggregateOtJson(t *testing.T) {
	t.Parallel()
	data, err := os.ReadFile("../../testData/opentelemetry.json")
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	got := aggregateOtJson(data, []string{"Matching taint entries", "Running DFA engine", "Running reverse debug", "SpanThatDoesNotExist"})

	assert.Contains(t, got, "Matching taint entries")
	assert.Equal(t, int32(60), got["Matching taint entries"].sum)
	assert.Equal(t, int32(3), got["Matching taint entries"].count)

	assert.Contains(t, got, "Running DFA engine")
	assert.Equal(t, int32(120), got["Running DFA engine"].sum)
	assert.Equal(t, int32(2), got["Running DFA engine"].count)

	assert.Contains(t, got, "Running reverse debug")
	assert.Equal(t, int32(15), got["Running reverse debug"].sum)
	assert.Equal(t, int32(1), got["Running reverse debug"].count)

	// a name with zero occurrences must be absent from the map, not present with a zero value
	assert.NotContains(t, got, "SpanThatDoesNotExist")
}

func TestAnalyzeQodanaReportAggregatesRepeatedSpans(t *testing.T) {
	t.Parallel()
	rawReport, err := os.ReadFile("../../testData/opentelemetry.json")
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	runResult := &RunResult{RawReport: rawReport}
	extraData := model.ExtraData{
		TcBuildProperties: []byte("{}"),
		CurrentBuildTime:  time.Now(),
	}

	ignore := analyzeQodanaReport(runResult, extraData)
	require.False(t, ignore)

	fields := runResult.ExtraFieldData
	measureNames, ok := fields[0].([]string)
	require.True(t, ok)
	measureValues, ok := fields[1].([]int32)
	require.True(t, ok)
	measureTypes, ok := fields[2].([]string)
	require.True(t, ok)

	indexOf := func(name string) int {
		for i, n := range measureNames {
			if n == name {
				return i
			}
		}
		return -1
	}

	dfaIndex := indexOf("Running DFA engine")
	require.GreaterOrEqual(t, dfaIndex, 0)
	assert.Equal(t, int32(120), measureValues[dfaIndex])
	assert.Equal(t, "d", measureTypes[dfaIndex])

	dfaCountIndex := indexOf("Running DFA engine.count")
	require.GreaterOrEqual(t, dfaCountIndex, 0)
	assert.Equal(t, int32(2), measureValues[dfaCountIndex])
	assert.Equal(t, "c", measureTypes[dfaCountIndex])

	// "Running DFA engine" must not appear as raw per-occurrence entries alongside the aggregate
	occurrences := 0
	for _, n := range measureNames {
		if n == "Running DFA engine" {
			occurrences++
		}
	}
	assert.Equal(t, 1, occurrences)
}
