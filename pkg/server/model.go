package server

import (
  "github.com/asaskevich/govalidator"
  "net/http"
  "net/url"
  "strings"
)

type MedianResult struct {
  metricName    string
  groupedValues []Value
}

type Value struct {
  group string
  value int
}

type BaseMetricQuery struct {
  product   string
  machines  []string
  eventType rune
}

func getProductAndMachine(query url.Values) (string, []string, rune, error) {
  product := query.Get("product")
  switch {
  case len(product) == 0:
    return "", nil, 0, NewHttpError(400, "The product parameter is required")

  case len(product) > 2:
    // prevent misuse of parameter
    return "", nil, 0, NewHttpError(400, "The product parameter is not correct")

  case !govalidator.IsAlpha(product):
    return "", nil, 0, NewHttpError(400, "The product parameter must contain only letters a-zA-Z")
  }

  machine := query.Get("machine")
  if len(machine) == 0 {
    return "", nil, 0, NewHttpError(400, "The machine parameter is required")
  }

  var normalizedEventTypeValue rune
  eventType := query.Get("eventType")
  if len(eventType) == 0 {
    normalizedEventTypeValue = 'd'
  } else {
    if len(eventType) != 1 {
      // prevent misuse of parameter
      return "", nil, 0, NewHttpError(400, "The eventType parameter must be one char")
    }
    normalizedEventTypeValue = rune(eventType[0])
  }

  return product, strings.Split(machine, ","), normalizedEventTypeValue, nil
}

func parseQuery(request *http.Request) (url.Values, error) {
  path := request.URL.Path
  index := strings.LastIndexByte(path, '/')
  var values url.Values
  if index != -1 {
    var err error
    values, err = url.ParseQuery(path[index+1:])
    if err != nil {
      return nil, err
    }
  }
  return values, nil
}
