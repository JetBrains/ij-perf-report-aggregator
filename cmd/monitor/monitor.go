package main

import (
  "flag"
  "fmt"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/server"
  "github.com/develar/errors"
  "github.com/json-iterator/go"
  "log"
  "net/http"
  "os"
)

// kubectl port-forward svc/clickhouse 9900:9000

func main() {
  serverUrl := flag.String("server", "http://localhost:9044", "The server URL.")

  err := analyzeTotal(*serverUrl)
  if err != nil {
    log.Fatal(fmt.Sprintf("%+v", err))
  }
}

type TeamCityActivity struct {
  value string

  startMarker string
  endMarker   string
}

func (t *TeamCityActivity) Start(value string) {
  if len(t.value) == 0 || t.value != value {
    t.End()
    t.value = value
    _, _ = fmt.Fprintf(os.Stdout, "##teamcity[%s name='%s']\n", t.startMarker, t.value)
  }
}

func (t *TeamCityActivity) End() {
  if len(t.value) != 0 {
    _, _ = fmt.Fprintf(os.Stdout, "##teamcity[%s name='%s']\n", t.endMarker, t.value)
  }
}

func analyzeTotal(serverUrl string) error {
  result, err := getResult(serverUrl, "2019-11-04")
  if err != nil {
    return errors.WithStack(err)
  }

  product := TeamCityActivity{
    startMarker: "testSuiteStarted",
    endMarker: "testSuiteFinished",
  }
  machineGroup := TeamCityActivity{
    startMarker: "testStarted",
    endMarker:   "testFinished",
  }

  for _, item := range result.GoldItems {
    product.Start(item.Product)
    testName := item.MachineGroup
    machineGroup.Start(testName)

    actualItems, ok := result.CurrentItems[item.MachineGroup]
    if !ok {
      _, _ = fmt.Fprintf(os.Stdout, "##teamcity[testStdErr name='%s' out='error text']\n", "no actual data")
      continue
    }

    actualItem := actualItems[item.Product]
    if len(actualItem) == 0 {
      _, _ = fmt.Fprintf(os.Stdout, "##teamcity[testStdErr name='%s' out='error text']\n", "no actual data")
      continue
    }

    diff := actualItem[0].Bootstrap - item.Bootstrap
    // 4 ms
    if diff < 4 {
      _, _ = fmt.Fprintf(os.Stdout, "##teamcity[testStdOut name='%s' out='gold: %.1f, actual: %.1f']\n", testName, item.Bootstrap, actualItem[0].Bootstrap)
      continue
    }

    _, _ = fmt.Fprintf(os.Stdout, "##teamcity[testFailed type='comparisonFailure' name='%s' message='Diff > 4 ms' expected='%.1f' actual='%.1f']\n", testName, item.Bootstrap, actualItem[0].Bootstrap)
  }

  machineGroup.End()
  product.End()

  return nil
}

func getResult(host string, goldWeekStart string) (*server.AnalyzeResult, error) {
  request, err := http.NewRequest("get", host+"/api/v1/compareMetrics?goldWeekStart="+goldWeekStart, nil)
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

  var result server.AnalyzeResult
  err = jsoniter.ConfigFastest.NewDecoder(response.Body).Decode(&result)
  if err != nil {
    return nil, errors.WithStack(err)
  }
  return &result, nil
}
