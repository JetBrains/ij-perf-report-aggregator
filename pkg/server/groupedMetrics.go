package server

import (
  "github.com/develar/errors"
  "github.com/json-iterator/go"
  "net/http"
  "report-aggregator/pkg/util"
  "strconv"
  "strings"
)

func (t *StatsServer) handleGroupedMetricsRequest(request *http.Request) ([]byte, error) {
  urlQuery, err := parseQuery(request)
  if err != nil {
    return nil, err
  }

  var query Query

  query.product, query.machine, err = getProductAndMachine(urlQuery)
  if err != nil {
    return nil, err
  }

  query.eventType = urlQuery.Get("eventType")
  if len(query.eventType) == 0 {
    query.eventType = "d"
  } else if len(query.eventType) != 1 {
    // prevent misuse of parameter
    return nil, NewHttpError(400, "eventType is not supported")
  }

  query.operator = urlQuery.Get("operator")
  if len(query.operator) == 0 {
    query.operator = "median"
  }

  operatorArg := urlQuery.Get("operatorArg")
  if query.operator == "quantile" {
    if len(operatorArg) == 0 {
      return nil, NewHttpError(400, "operatorArg parameter is required")
    }

    v, err := strconv.ParseInt(operatorArg, 10, 8)
    if err != nil {
      return nil, NewHttpError(400, "quantile is not correct")
    }
    query.quantile = float64(v) / 100
  }

  var metricNames []string
  if query.eventType == "d" {
    metricNames = essentialMetricNames
  } else {
    metricNames = instantMetricNames
  }

  results := make([]*MedianResult, len(metricNames))
  err = util.MapAsyncConcurrency(len(metricNames), 4, func(taskIndex int) (f func() error, err error) {
    return func() error {
      result, err := t.computeMedian(metricNames[taskIndex], query)
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

  buffer := byteBufferPool.Get()
  defer byteBufferPool.Put(buffer)
  WriteGroupedMetricList(buffer, results)
  return CopyBuffer(buffer), nil
}

type Query struct {
  product   string
  machine   string
  eventType string

  operator string
  quantile float64
}

func (t *StatsServer) computeMedian(metricName string, query Query) (*MedianResult, error) {
  result := &MedianResult{
    metricName: metricName,
  }

  var s strings.Builder
  s.WriteString(query.operator)
  s.WriteRune('(')
  if query.operator == "quantile" {
    s.WriteString(strconv.FormatFloat(query.quantile, 'f', 1, 32))
    s.WriteString(", ")
  }
  s.WriteString(metricName)
  s.WriteRune('_')
  s.WriteString(query.eventType)
  s.WriteString(`{product="` + query.product + `",machine="` + query.machine + `"}`)
  s.WriteString("[2y]")
  s.WriteRune(')')
  s.WriteString(` by (buildC1)`)

  response, err := t.performRequest(s.String())
  if err != nil {
    return nil, err
  }

  defer util.Close(response.Body, t.logger)
  err = t.readJson(response, result)
  if err != nil {
    return nil, err
  }

  return result, nil
}

func (t *StatsServer) readJson(response *http.Response, result *MedianResult) error {
  iterator := jsoniter.Parse(jsoniter.ConfigFastest, response.Body, 64*1024)
  for {
    field := iterator.ReadObject()
    switch field {
    case "status":
      status := iterator.ReadString()
      if status != "success" {
        return errors.Errorf("query status: %s", status)
      }

    case "data":
      err := readData(iterator, result)
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
