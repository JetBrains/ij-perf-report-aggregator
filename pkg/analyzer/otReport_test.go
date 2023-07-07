package analyzer

import (
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
      if value, ok := got[name]; !ok {
        t.Errorf("Expected key '%s' not found in output", name)
      } else if value[0] != expected {
        t.Errorf("Incorrect value for %s: got %v, want %v", name, value[0], expected)
      }
    })
  }
}
