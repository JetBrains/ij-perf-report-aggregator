package server

import (
  "net/http"
  "net/url"
  "strings"
)

type MedianResult struct {
  metricName   string
  buildToValue []Value
}

type Value struct {
  buildC1 int
  value   int
}

func getProductAndMachine(query url.Values) (string, string, error) {
  product := query.Get("product")
  if len(product) == 0 {
    return "", "", NewHttpError(400, "product parameter is required")
  } else if len(product) > 2 {
    // prevent misuse of parameter
    return "", "", NewHttpError(400, "product code is not correct")
  }

  machine := query.Get("machine")
  if len(machine) == 0 {
    return "", "", NewHttpError(400, "machine parameter is required")
  }

  for _, c := range machine {
    if c < '0' || c > '9' {
      return "", "", NewHttpError(400, "machine id is not numeric")
    }
  }
  return product, machine, nil
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
