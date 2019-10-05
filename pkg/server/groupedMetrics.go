package server

import (
  "github.com/develar/errors"
  "github.com/json-iterator/go"
  "github.com/valyala/quicktemplate"
  "net/http"
  "net/url"
  "report-aggregator/pkg/util"
  "strconv"
  "time"
)

var httpClient = &http.Client{Timeout: 30 * time.Second}

func (t *StatsServer) handleGroupedMetricsRequest(request *http.Request) ([]byte, error) {
  query, err := parseQuery(request)
  if err != nil {
    return nil, err
  }

  product, machine, err := getProductAndMachine(query)
  if err != nil {
    return nil, err
  }

  eventType := query.Get("eventType")
  if len(eventType) == 0 {
    eventType = "d"
  } else if len(eventType) != 1 {
    // prevent misuse of parameter
    return nil, NewHttpError(400, "eventType is not supported")
  }

  operator := query.Get("operator")
  if len(operator) == 0 {
    operator = "median"
  } else if len(operator) > 6 {
    // prevent misuse of parameter
    return nil, NewHttpError(400, "operator is not supported")
  }

  var metricNames []string
  if eventType == "d" {
    metricNames = essentialMetricNames
  } else {
    metricNames = instantMetricNames
  }

  results := make([]*MedianResult, len(metricNames))
  err = util.MapAsyncConcurrency(len(metricNames), 4, func(taskIndex int) (f func() error, err error) {
    return func() error {
      result, err := t.computeMedian(metricNames[taskIndex], product, machine, eventType, operator, httpClient)
      if err != nil {
        return err
      }
      results[taskIndex] = result
      return nil
    }, nil
  })
  if err != nil {
    return nil, err
  }

  buffer := quicktemplate.AcquireByteBuffer()
  defer quicktemplate.ReleaseByteBuffer(buffer)
  WriteGroupedMetricList(buffer, results)
  return CopyBuffer(buffer), nil
}

func (t *StatsServer) computeMedian(metricName string, product string, machine string, eventType string, operator string, httpClient *http.Client) (*MedianResult, error) {
  result := &MedianResult{
    metricName: metricName,
  }

  u, err := url.Parse("http://localhost:8428/api/v1/query")
  if err != nil {
    return nil, err
  }

  q := u.Query()

  q.Set("query", operator+`(`+metricName+`_`+eventType+`{product="`+product+`",machine="`+machine+`"}[2y]) by (buildC1)`)
  u.RawQuery = q.Encode()

  err = t.getJson(u.String(), httpClient, result)
  if err != nil {
    return nil, err
  }

  return result, nil
}

func (t *StatsServer) getJson(url string, httpClient *http.Client, result *MedianResult) error {
  r, err := httpClient.Get(url)
  if err != nil {
    return err
  }

  defer util.Close(r.Body, t.logger)

  if r.StatusCode >= 400 {
    return errors.Errorf("Request failed: %s", r.Status)
  }

  iterator := jsoniter.Parse(jsoniter.ConfigFastest, r.Body, 64*1024)
  for {
    field := iterator.ReadObject()
    switch field {
    case "status":
      status := iterator.ReadString()
      if status != "success" {
        return errors.Errorf("query status: %s", status)
      }

    case "data":
      err = readData(iterator, result)
      if err != nil {
        return err
      }

    case "":
      return nil

    default:
      iterator.Skip()
    }
  }
}

func readData(iterator *jsoniter.Iterator, result *MedianResult) error {
  for {
    field := iterator.ReadObject()
    switch field {
    case "resultType":
      resultType := iterator.ReadString()
      if resultType != "vector" {
        return errors.Errorf("unexpected resultType: %s", resultType)
      }

    case "result":
      for iterator.ReadArray() {
        err := readResultItem(iterator, result)
        if err != nil {
          return err
        }
      }

    case "":
      return nil

    default:
      iterator.Skip()
    }
  }
}

func readResultItem(iterator *jsoniter.Iterator, result *MedianResult) error {
  var err error

  // use -2 as null because sometimes value -1 is a valid value
  v := Value{
    buildC1: -2,
    value:   -2,
  }

readResultItem:
  for {
    field := iterator.ReadObject()
    switch field {
    case "metric":
      v.buildC1, err = readMetric(iterator)
      if err != nil {
        return err
      }

    case "value":
      v.value, err = readValue(iterator)
      if err != nil {
        return err
      }

    case "":
      break readResultItem

    default:
      iterator.Skip()
    }
  }

  if v.buildC1 == -2 {
    return errors.New("buildC1 not found")
  }
  if v.value == -2 {
    return errors.New("value not found")
  }

  result.buildToValue = append(result.buildToValue, v)

  return nil
}

func readMetric(iterator *jsoniter.Iterator) (int, error) {
  var err error

  buildC1 := -2

readMetric:
  for {
    field := iterator.ReadObject()
    switch field {
    case "buildC1":
      buildC1, err = strconv.Atoi(iterator.ReadString())
      if err != nil {
        return -1, err
      }

    case "":
      break readMetric

    default:
      iterator.Skip()
    }
  }
  return buildC1, nil
}

func readValue(iterator *jsoniter.Iterator) (int, error) {
  if !iterator.ReadArray() {
    return -1, errors.New("2 values are expected")
  }

  // skip timestamp
  iterator.Skip()

  if !iterator.ReadArray() {
    return -1, errors.New("2 values are expected")
  }

  value, err := strconv.Atoi(iterator.ReadString())
  if iterator.ReadArray() {
    return -1, errors.New("only 2 values are expected")
  }
  return value, err
}
