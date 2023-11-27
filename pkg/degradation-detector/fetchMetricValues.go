package degradation_detector

import (
  "context"
  "encoding/json"
  "fmt"
  dataQuery "github.com/JetBrains/ij-perf-report-aggregator/pkg/data-query"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/util"
  "github.com/alitto/pond"
  "io"
  "log/slog"
  "net/http"
  "strings"
  "sync"
  "time"
)

type responseData struct {
  timestamps []int64
  values     []int
  builds     []string
}

type responseDataWithSettings struct {
  responseData
  Settings
}

type dataQueryProducer interface {
  DataQuery() dataQuery.DataQuery
}

func fetchMetricsFromClickhouse(settings []Settings, client *http.Client, backendUrl string) chan responseDataWithSettings {
  dataChan := make(chan responseDataWithSettings, 5)
  go func() {
    defer close(dataChan)
    var wg sync.WaitGroup
    pool := pond.New(5, 1000)
    for _, setting := range settings {
      wg.Add(1)
      setting := setting
      pool.Submit(func() {
        defer wg.Done()
        ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
        defer cancel()
        data, err := getDataFromClickhouse(ctx, client, backendUrl, setting.DataQuery())
        if err != nil {
          slog.Error("error while getting responseData from clickhouse", "error", err, "settings", setting)
          return
        }
        slog.Info("fetched from clickhouse", "settings", setting)
        dataChan <- responseDataWithSettings{
          responseData: data,
          Settings:     setting,
        }
      })
    }
    wg.Wait()
  }()
  return dataChan
}

func getDataFromClickhouse(ctx context.Context, client *http.Client, backendUrl string, query dataQuery.DataQuery) (responseData, error) {
  response, err := getValuesFromServer(ctx, client, backendUrl, query)
  if err != nil {
    return responseData{}, err
  }
  data, err := extractDataFromRequest(response)
  return data, err
}

func getValuesFromServer(ctx context.Context, client *http.Client, backendURL string, query dataQuery.DataQuery) ([]byte, error) {
  url := backendURL + "/api/q/"
  queries := []dataQuery.DataQuery{query}
  jsonQuery, err := json.Marshal(queries)
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
    return nil, fmt.Errorf("failed to get responseData: %v", resp.Status)
  }
  body, err := io.ReadAll(resp.Body)
  if err != nil {
    return nil, fmt.Errorf("failed to read response body: %w", err)
  }
  return body, err
}

func (s StartupSettings) DataQuery() dataQuery.DataQuery {
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

  query := dataQuery.DataQuery{
    Database: "ij",
    Table:    "report",
    Fields:   fields,
    Filters:  filters,
    Order:    []string{"t"},
  }
  return query
}

func (s PerformanceSettings) DataQuery() dataQuery.DataQuery {
  fields := []dataQuery.DataQueryDimension{
    {Name: "t", Sql: "toUnixTimestamp(generated_time)*1000"},
    {Name: "measures", SubName: "value"},
  }
  if s.Db == "perfint" {
    fields = append(fields, dataQuery.DataQueryDimension{Name: "Build", Sql: "concat(toString(build_c1),'.',toString(build_c2))"})
  } else {
    fields = append(fields, dataQuery.DataQueryDimension{Name: "tc_build_id"})
  }

  query := dataQuery.DataQuery{
    Database: s.Db,
    Table:    s.Table,
    Fields:   fields,
    Filters: []dataQuery.DataQueryFilter{
      {Field: "branch", Value: s.Branch},
      {Field: "generated_time", Sql: ">subtractDays(now(),100)"},
      {Field: "project", Value: s.Project},
      {Field: "measures.name", Value: s.Metric},
      {Field: "machine", Value: s.Machine, Operator: "like"},
      {Field: "triggeredBy", Value: ""},
    },
    Order: []string{"t"},
  }
  return query
}

func extractDataFromRequest(response []byte) (responseData, error) {
  var data [][][]interface{}

  err := json.Unmarshal(response, &data)
  if err != nil {
    return responseData{}, fmt.Errorf("failed to decode JSON: %w", err)
  }
  if len(data) == 0 {
    return responseData{}, fmt.Errorf("no responseData")
  }
  if len(data[0]) < 3 {
    return responseData{}, fmt.Errorf("not enough responseData")
  }
  timestamps, err := SliceToSliceInt64(data[0][0])
  if err != nil {
    return responseData{}, fmt.Errorf("failed to convert values: %w", err)
  }
  values, err := SliceToSliceOfInt(data[0][1])
  if err != nil {
    return responseData{}, fmt.Errorf("failed to convert values: %w", err)
  }
  builds, err := SliceToSliceOfString(data[0][2])
  if err != nil {
    return responseData{}, fmt.Errorf("failed to convert values: %w", err)
  }
  return responseData{
    timestamps: timestamps,
    values:     values,
    builds:     builds,
  }, err
}
