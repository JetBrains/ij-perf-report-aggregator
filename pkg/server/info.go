package server

import (
  "errors"
  "net/http"
  "report-aggregator/pkg/model"
  "report-aggregator/pkg/util"
)

func (t *StatsServer) handleInfoRequest(request *http.Request) ([]byte, error) {
  var productNames []string
  groupNames := t.machineInfo.groupNames
  rows, err := t.db.QueryContext(request.Context(), "select distinct product, machine from report where machine != 'Dead agent' group by product, machine order by product, machine")
  if err != nil {
    return nil, err
  }

  defer util.Close(rows, t.logger)

  productNameToMachineMap := make(map[string]map[string]*MachineGroup)
  for rows.Next() {
    var product string
    var machine string
    err := rows.Scan(&product, &machine)
    if err != nil {
      return nil, err
    }

    groupName, ok := t.machineInfo.machineToGroupName[machine]
    if !ok {
      return nil, errors.New("Group is unknown machine: " + machine)
    }

    groupToMachine, ok := productNameToMachineMap[product]
    if !ok {
      productNames = append(productNames, product)
      groupToMachine = make(map[string]*MachineGroup)
      productNameToMachineMap[product] = groupToMachine
    }

    machineGroup := groupToMachine[groupName]
    if machineGroup == nil {
      machineGroup = &MachineGroup{name: groupName}
      groupToMachine[groupName] = machineGroup
    }

    machineGroup.machines = append(machineGroup.machines, machine)
  }

  buffer := byteBufferPool.Get()
  defer byteBufferPool.Put(buffer)
  WriteInfo(buffer, productNames, groupNames, model.DurationMetricNames, model.InstantMetricNames, productNameToMachineMap)
  return CopyBuffer(buffer), nil
}

type MachineGroup struct {
  name     string
  machines []string
}
