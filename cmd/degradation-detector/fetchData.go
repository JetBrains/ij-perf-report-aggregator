package main

import (
  "context"
  "encoding/json"
  "fmt"
  dataQuery "github.com/JetBrains/ij-perf-report-aggregator/pkg/data-query"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/util"
  "io"
  "log"
  "net/http"
)

func getDataFromClickhouse(ctx context.Context, analysisSettings AnalysisSettings) ([]int64, []int, []string, error) {
  response, err := getValuesFromServer(ctx, analysisSettings)
  if err != nil {
    log.Printf("%v", err)
  }
  timestamps, values, builds, err := extractDataFromRequest(response)
  if err != nil {
    log.Printf("%v", err)
  }
  return timestamps, values, builds, err
}

func getValuesFromServer(ctx context.Context, analysisSettings AnalysisSettings) ([]byte, error) {
  url := URL + "/api/q/"
  query := getDataQuery(analysisSettings)
  jsonQuery, err := json.Marshal(query)
  if err != nil {
    return nil, fmt.Errorf("failed to marshal query: %w", err)
  }

  encoded, err := util.EncodeQuery(jsonQuery)
  if err != nil {
    return nil, fmt.Errorf("failed to encode query: %w", err)
  }

  req, err := http.NewRequestWithContext(ctx, http.MethodGet, url+encoded, nil)
  if err != nil {
    return nil, fmt.Errorf("failed to create request: %w", err)
  }

  client := &http.Client{}
  resp, err := client.Do(req)
  if err != nil {
    return nil, fmt.Errorf("failed to send GET request: %w", err)
  }
  defer resp.Body.Close()

  if resp.StatusCode != http.StatusOK {
    return nil, fmt.Errorf("failed to get data: %v", resp.Status)
  }
  body, err := io.ReadAll(resp.Body)
  if err != nil {
    return nil, fmt.Errorf("failed to read response body: %w", err)
  }
  return body, err
}

func getDataQuery(settings AnalysisSettings) []dataQuery.DataQuery {
  fields := []dataQuery.DataQueryDimension{
    {Name: "t", Sql: "toUnixTimestamp(generated_time)*1000"},
    {Name: "measures", SubName: "value"},
  }
  if settings.db == "perfint" {
    fields = append(fields, dataQuery.DataQueryDimension{Name: "build", Sql: "concat(toString(build_c1),'.',toString(build_c2))"})
  } else if settings.db == "perfintDev" {
    fields = append(fields, dataQuery.DataQueryDimension{Name: "tc_build_id"})
  }

  queries := []dataQuery.DataQuery{
    {
      Database: settings.db,
      Table:    settings.table,
      Fields:   fields,
      Filters: []dataQuery.DataQueryFilter{
        {Field: "branch", Value: "master"},
        {Field: "generated_time", Sql: ">subtractDays(now(),100)"},
        {Field: "project", Value: settings.test},
        {Field: "measures.name", Value: settings.metric},
        {Field: "machine", Value: "intellij-linux-performance-aws-%", Operator: "like"},
        {Field: "triggeredBy", Value: ""},
        {Field: "build_c3", Value: 0, Operator: "="},
      },
      Order: []string{"t"},
    },
  }
  return queries
}

func extractDataFromRequest(response []byte) ([]int64, []int, []string, error) {
  var data [][][]interface{}

  err := json.Unmarshal(response, &data)
  if err != nil {
    return nil, nil, nil, fmt.Errorf("failed to decode JSON: %w", err)
  }
  timestamps, err := sliceToSliceInt64(data[0][0])
  if err != nil {
    return nil, nil, nil, fmt.Errorf("failed to convert values: %w", err)
  }
  values, err := sliceToSliceOfInt(data[0][1])
  if err != nil {
    return nil, nil, nil, fmt.Errorf("failed to convert values: %w", err)
  }
  builds, err := sliceToSliceOfString(data[0][2])
  if err != nil {
    return nil, nil, nil, fmt.Errorf("failed to convert values: %w", err)
  }
  return timestamps, values, builds, err
}
