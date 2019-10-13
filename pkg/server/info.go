package server

import (
  "net/http"
  "report-aggregator/pkg/util"
)

var EssentialDurationMetricNames = []string{"bootstrap", "appInitPreparation", "appInit", "pluginDescriptorLoading", "appComponentCreation", "projectComponentCreation"}
var DurationMetricNames = append(EssentialDurationMetricNames, "moduleLoading")
var InstantMetricNames = []string{"splash", "startUpCompleted"}

func ProcessMetricName(handler func(name string, isInstant bool)) {
  for _, name := range DurationMetricNames {
    handler(name, false)
  }
  for _, name := range InstantMetricNames {
    handler(name, true)
  }
}

func (t *StatsServer) handleInfoRequest(request *http.Request) ([]byte, error) {
  var productNames []string
  rows, err := t.db.QueryContext(request.Context(), "select distinct product, machine from report group by product, machine order by product, machine")
  if err != nil {
    return nil, err
  }

  defer util.Close(rows, t.logger)

  productNameToMachineNames := make(map[string][]string)
  for rows.Next() {
    var product string
    var machine string
    err := rows.Scan(&product, &machine)
    if err != nil {
      return nil, err
    }

    machines, ok := productNameToMachineNames[product]
    if ok {
      productNameToMachineNames[product] = append(machines, machine)
    } else {
      productNames = append(productNames, product)
      productNameToMachineNames[product] = []string{machine}
    }
  }

  buffer := byteBufferPool.Get()
  defer byteBufferPool.Put(buffer)
  WriteInfo(buffer, productNames, DurationMetricNames, InstantMetricNames, productNameToMachineNames)
  return CopyBuffer(buffer), nil
}