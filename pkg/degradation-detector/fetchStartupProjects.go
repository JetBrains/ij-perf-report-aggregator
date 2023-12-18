package degradation_detector

import (
  "context"
  dataQuery "github.com/JetBrains/ij-perf-report-aggregator/pkg/data-query"
  "net/http"
  "time"
)

func FetchAllProjects(backendUrl string, client *http.Client, settings StartupSettings) ([]string, error) {
  ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
  defer cancel()
  query := dataQuery.Query{
    Database: "ij",
    Table:    "report",
    Fields:   []dataQuery.QueryDimension{{Name: "project", Sql: "distinct project"}},
    Flat:     true,
    Filters: []dataQuery.QueryFilter{
      {Sql: "not endsWith(project, '(fast installer)')"},
      {Field: "branch", Value: settings.Branch},
      {Field: "generated_time", Sql: ">subtractDays(now(),100)"},
      {Field: "machine", Value: settings.Machine, Operator: "like"},
      {Field: "triggeredBy", Value: ""},
      {Field: "product", Value: settings.Product},
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
