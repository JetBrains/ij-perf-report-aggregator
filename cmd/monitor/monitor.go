package main

import (
  "flag"
  "fmt"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/server"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/util"
  "github.com/develar/errors"
  "github.com/json-iterator/go"
  "go.uber.org/zap"
  "io"
  "log"
  "net/http"
  "net/url"
  "strings"
)

// kubectl port-forward svc/clickhouse 9900:9000

func main() {
  logger := util.CreateLogger()
  defer func() {
    _ = logger.Sync()
  }()

  serverUrl := flag.String("server", "http://localhost:9044", "The server URL.")

  flag.Parse()

  err := analyzeTotal(*serverUrl, logger)
  if err != nil {
    log.Fatalf("%+v", err)
  }
}

func analyzeTotal(serverUrl string, logger *zap.Logger) error {
  result, err := getResult(serverUrl, "2019-11-04", logger)
  if err != nil {
    return errors.WithStack(err)
  }

  product := TeamCityActivity{
    startMarker: "testSuiteStarted",
    endMarker:   "testSuiteFinished",
  }
  machineGroup := TeamCityTest{
    TeamCityActivity: TeamCityActivity{
      startMarker: "testStarted",
      endMarker:   "testFinished",
    },
  }

  for _, item := range result.GoldItems {
    product.Start(item.Product)
    machineGroup.Start(item.MachineGroup)

    actualItems, ok := result.CurrentItems[item.MachineGroup]
    if !ok {
      machineGroup.Failed("no actual data")
      continue
    }

    actualItem := actualItems[item.Product]
    if len(actualItem) == 0 {
      machineGroup.Failed("no actual data")
      continue
    }

    diff := actualItem[0].Bootstrap - item.Bootstrap
    // 4 ms
    if diff < 4 {
      machineGroup.Output(fmt.Sprintf("gold: %.1f, actual: %.1f", item.Bootstrap, actualItem[0].Bootstrap))
      continue
    }

    machineGroup.CompareFailed("Diff > 4 ms", actualItem[0].Bootstrap, item.Bootstrap)
  }

  machineGroup.End()
  product.End()

  return nil
}

func getResult(host string, goldWeekStart string, logger *zap.Logger) (*server.AnalyzeResult, error) {
  serverUrl, err := url.Parse(strings.TrimSuffix(host, "/") + "/api/v1/compareMetrics")
  if err != nil {
    return nil, err
  }

  q := serverUrl.Query()
  q.Set("goldWeekStart", goldWeekStart)
  serverUrl.RawQuery = q.Encode()
  urlString := serverUrl.String()
  logger.Info("get data", zap.String("url", urlString))
  request, err := http.NewRequest("GET", urlString, nil)
  if err != nil {
    return nil, errors.WithStack(err)
  }

  request.Header.Set("Accept", "application/json")
  response, err := http.DefaultClient.Do(request)
  if err != nil {
    return nil, errors.WithStack(err)
  }

  defer func() {
    _ = response.Body.Close()
  }()

  if response.StatusCode >= 400 {
    responseBytes, _ := io.ReadAll(response.Body)
    return nil, errors.Errorf("status: %s, response: %s", response.Status, responseBytes)
  }

  var result server.AnalyzeResult
  err = jsoniter.ConfigFastest.NewDecoder(response.Body).Decode(&result)
  if err != nil {
    return nil, errors.WithStack(err)
  }
  return &result, nil
}
