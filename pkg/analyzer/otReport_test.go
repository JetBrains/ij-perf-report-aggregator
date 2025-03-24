package analyzer

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
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
