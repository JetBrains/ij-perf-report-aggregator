package analyzer

import (
  "encoding/json"
)

type Trace struct {
  Data []struct {
    Spans []struct {
      OperationName string `json:"operationName"`
      Duration      int    `json:"duration"`
    } `json:"spans"`
  } `json:"data"`
}

func analyzeOtJson(ot []byte, operationNames []string) map[string][]int {
  var trace Trace
  err := json.Unmarshal(ot, &trace)
  if err != nil {
    return nil
  }
  durationMap := make(map[string][]int)
  for _, data := range trace.Data {
    for _, span := range data.Spans {
      for _, operationName := range operationNames {
        if span.OperationName == operationName {
          durationMap[operationName] = append(durationMap[operationName], span.Duration)
        }
      }
    }
  }
  return durationMap
}
