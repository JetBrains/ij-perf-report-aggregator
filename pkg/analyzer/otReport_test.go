package analyzer

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/JetBrains/ij-perf-report-aggregator/pkg/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func readTestTrace(t *testing.T) Trace {
	t.Helper()
	data, err := os.ReadFile("../../testData/opentelemetry.json")
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}
	var trace Trace
	if err := json.Unmarshal(data, &trace); err != nil {
		t.Fatalf("Failed to unmarshal trace: %v", err)
	}
	return trace
}

func TestFilter(t *testing.T) {
	t.Parallel()
	trace := readTestTrace(t)

	testCases := map[string]int32{
		"project.opening":   5035,
		"globalInspections": 141561,
	}

	got := analyzeOtJson(trace, []string{"globalInspections", "project.opening"})

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
	trace := readTestTrace(t)

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

	got := analyzeOtJson(trace, operationNames)

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
	trace := readTestTrace(t)

	got := aggregateOtJson(trace, []string{"Matching taint entries", "Running DFA engine", "Running reverse debug", "SpanThatDoesNotExist"})

	assert.Contains(t, got, "Matching taint entries")
	assert.Equal(t, int32(60), got["Matching taint entries"].sum)
	assert.Equal(t, int32(3), got["Matching taint entries"].count)

	assert.Contains(t, got, "Running DFA engine")
	assert.Equal(t, int32(120), got["Running DFA engine"].sum)
	assert.Equal(t, int32(2), got["Running DFA engine"].count)

	assert.Contains(t, got, "Running reverse debug")
	assert.Equal(t, int32(15), got["Running reverse debug"].sum)
	assert.Equal(t, int32(1), got["Running reverse debug"].count)

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

	occurrences := 0
	for _, n := range measureNames {
		if n == "Running DFA engine" {
			occurrences++
		}
	}
	assert.Equal(t, 1, occurrences)
}
