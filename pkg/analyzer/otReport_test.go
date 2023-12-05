package analyzer

import (
  "github.com/stretchr/testify/assert"
  "os"
  "testing"
)

func TestFilter(t *testing.T) {
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
      assert.Contains(t, got, name)
      assert.Equal(t, expected, got[name][0])
    })
  }
}
