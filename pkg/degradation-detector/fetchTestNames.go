package degradation_detector

import (
  "context"
  "encoding/json"
  "errors"
  "fmt"
  dataQuery "github.com/JetBrains/ij-perf-report-aggregator/pkg/data-query"
  "log/slog"
  "net/http"
  "strings"
  "time"
)

func ExpandTestsByPattern(backendUrl string, client *http.Client, tests []string, baseSettings PerformanceSettings) []string {
  testsExpanded := make([]string, 0, len(tests)*5)
  for _, test := range tests {
    if strings.Contains(test, "%") {
      matchingTests, err := fetchTestsByPattern(backendUrl, client, baseSettings, test)
      if err != nil {
        slog.Error("error while fetching tests by pattern", "error", err, "pattern", test)
        continue
      }
      testsExpanded = append(testsExpanded, matchingTests...)
    } else {
      testsExpanded = append(testsExpanded, test)
    }
  }
  return testsExpanded
}

func fetchTestsByPattern(backendUrl string, client *http.Client, settings PerformanceSettings, pattern string) ([]string, error) {
  ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
  defer cancel()
  filters := []dataQuery.QueryFilter{
    {Field: "branch", Value: settings.Branch},
    {Field: "generated_time", Sql: ">subtractDays(now(),30)"},
    {Field: "machine", Value: settings.Machine, Operator: "like"},
    {Field: "triggeredBy", Value: ""},
  }
  if pattern != "" {
    filters = append(filters, dataQuery.QueryFilter{Field: "project", Value: pattern, Operator: "like"})
  }
  query := dataQuery.Query{
    Database: settings.Db,
    Table:    settings.Table,
    Fields:   []dataQuery.QueryDimension{{Name: "project", Sql: "distinct project"}},
    Flat:     true,
    Filters:  filters,
    Order:    []string{"project"},
  }
  response, err := getValuesFromServer(ctx, client, backendUrl, query)
  if err != nil {
    return nil, err
  }
  tests, err := extractValuesFromRequest(response)
  if err != nil {
    return nil, err
  }
  return tests, nil
}

func FetchAllTests(backendUrl string, client *http.Client, settings PerformanceSettings) ([]string, error) {
  return fetchTestsByPattern(backendUrl, client, settings, "")
}

func extractValuesFromRequest(response []byte) ([]string, error) {
  var data [][]interface{}

  err := json.Unmarshal(response, &data)
  if err != nil {
    return nil, fmt.Errorf("failed to decode JSON: %w", err)
  }
  if len(data) == 0 {
    return nil, errors.New("no data")
  }
  if len(data[0]) < 1 {
    return nil, errors.New("not enough data")
  }
  tests, err := SliceToSliceOfString(data[0])
  if err != nil {
    return nil, fmt.Errorf("failed to convert values: %w", err)
  }
  return tests, nil
}
