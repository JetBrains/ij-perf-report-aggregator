package degradation_detector

import (
  "context"
  "encoding/json"
  "fmt"
  dataQuery "github.com/JetBrains/ij-perf-report-aggregator/pkg/data-query"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/util"
  "io"
  "net/http"
  "strings"
)

func GetDataFromClickhouse(ctx context.Context, client *http.Client, backendURL string, analysisSettings Settings) ([]int64, []int, []string, error) {
  response, err := GetValuesFromServer(ctx, client, backendURL, analysisSettings.DataQuery())
  if err != nil {
    return nil, nil, nil, err
  }
  timestamps, values, builds, err := extractDataFromRequest(response)
  return timestamps, values, builds, err
}

func GetValuesFromServer(ctx context.Context, client *http.Client, backendURL string, query []dataQuery.DataQuery) ([]byte, error) {
  url := backendURL + "/api/q/"
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

func (s StartupSettings) DataQuery() []dataQuery.DataQuery {
  fields := []dataQuery.DataQueryDimension{
    {Name: "t", Sql: "toUnixTimestamp(generated_time)*1000"},
  }
  filters := []dataQuery.DataQueryFilter{
    {Field: "branch", Value: s.Branch},
    {Field: "generated_time", Sql: ">subtractDays(now(),100)"},
    {Field: "project", Value: s.Project},
    {Field: "product", Value: s.Product},
    {Field: "machine", Value: s.Machine, Operator: "like"},
    {Field: "triggeredBy", Value: ""},
  }
  if strings.Contains(s.Metric, "/") {
    filters = append(filters, dataQuery.DataQueryFilter{Field: "metrics.name", Value: s.Metric})
    fields = append(fields, dataQuery.DataQueryDimension{Name: "metrics", SubName: "value"})
  }
  if strings.HasSuffix(s.Metric, ".end") {
    metricName, _ := strings.CutSuffix(s.Metric, ".end")
    filters = append(filters, dataQuery.DataQueryFilter{Field: "measure.name", Value: metricName})
    fields = append(fields, dataQuery.DataQueryDimension{Name: "measure", SubName: "end", Sql: "(measure.start+measure.duration)"})
  }
  if !strings.HasSuffix(s.Metric, ".end") && !strings.Contains(s.Metric, "/") {
    fields = append(fields, dataQuery.DataQueryDimension{Name: s.Metric})
  }
  fields = append(fields, dataQuery.DataQueryDimension{Name: "Build", Sql: "concat(toString(build_c1),'.',toString(build_c2))"})

  queries := []dataQuery.DataQuery{
    {
      Database: "ij",
      Table:    "report",
      Fields:   fields,
      Filters:  filters,
      Order:    []string{"t"},
    },
  }
  return queries
}

func (s PerformanceSettings) DataQuery() []dataQuery.DataQuery {
  fields := []dataQuery.DataQueryDimension{
    {Name: "t", Sql: "toUnixTimestamp(generated_time)*1000"},
    {Name: "measures", SubName: "value"},
  }
  if s.Db == "perfint" {
    fields = append(fields, dataQuery.DataQueryDimension{Name: "Build", Sql: "concat(toString(build_c1),'.',toString(build_c2))"})
  } else {
    fields = append(fields, dataQuery.DataQueryDimension{Name: "tc_build_id"})
  }

  queries := []dataQuery.DataQuery{
    {
      Database: s.Db,
      Table:    s.Table,
      Fields:   fields,
      Filters: []dataQuery.DataQueryFilter{
        {Field: "branch", Value: s.Branch},
        {Field: "generated_time", Sql: ">subtractDays(now(),100)"},
        {Field: "project", Value: s.Branch},
        {Field: "measures.name", Value: s.Metric},
        {Field: "machine", Value: s.Machine, Operator: "like"},
        {Field: "triggeredBy", Value: ""},
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
  if len(data) == 0 {
    return nil, nil, nil, fmt.Errorf("no data")
  }
  if len(data[0]) < 3 {
    return nil, nil, nil, fmt.Errorf("not enough data")
  }
  timestamps, err := SliceToSliceInt64(data[0][0])
  if err != nil {
    return nil, nil, nil, fmt.Errorf("failed to convert values: %w", err)
  }
  values, err := SliceToSliceOfInt(data[0][1])
  if err != nil {
    return nil, nil, nil, fmt.Errorf("failed to convert values: %w", err)
  }
  builds, err := SliceToSliceOfString(data[0][2])
  if err != nil {
    return nil, nil, nil, fmt.Errorf("failed to convert values: %w", err)
  }
  return timestamps, values, builds, err
}
