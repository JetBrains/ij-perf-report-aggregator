package degradation_detector

import (
  "context"
  "encoding/json"
  "fmt"
  dataQuery "github.com/JetBrains/ij-perf-report-aggregator/pkg/data-query"
  "net/http"
  "time"
)

func FetchAllTests(backendUrl string, client *http.Client, settings PerformanceSettings) ([]string, error) {
  ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
  defer cancel()
  query := dataQuery.Query{
    Database: settings.Db,
    Table:    settings.Table,
    Fields:   []dataQuery.QueryDimension{{Name: "project", Sql: "distinct project"}},
    Flat:     true,
    Filters: []dataQuery.QueryFilter{
      {Field: "branch", Value: settings.Branch},
      {Field: "generated_time", Sql: ">subtractDays(now(),100)"},
      {Field: "machine", Value: settings.Machine, Operator: "like"},
      {Field: "triggeredBy", Value: ""},
    },
    Order: []string{"project"},
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

func extractValuesFromRequest(response []byte) ([]string, error) {
  var data [][]interface{}

  err := json.Unmarshal(response, &data)
  if err != nil {
    return nil, fmt.Errorf("failed to decode JSON: %w", err)
  }
  if len(data) == 0 {
    return nil, fmt.Errorf("no data")
  }
  if len(data[0]) < 1 {
    return nil, fmt.Errorf("not enough data")
  }
  tests, err := SliceToSliceOfString(data[0])
  if err != nil {
    return nil, fmt.Errorf("failed to convert values: %w", err)
  }
  return tests, nil
}
