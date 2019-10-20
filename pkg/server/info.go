package server

import (
  "errors"
  "net/http"
  "report-aggregator/pkg/model"
  "report-aggregator/pkg/util"
)

// Mac mini Space Gray/3.0 GHz 6C/8GB/256GB; Model No. A1993; Part No. MRTT2RU/A; Serial No. C07XX9PFJYVX; Prod.12/2018, for code-sign (ADM-32069) -> ADM-35488
const macMini = "macMini 2018"

// Core i7-3770 16Gb, Intel SSD 535
const win = "Windows: i7-3770, 16Gb, Intel SSD 535"

// old RAM	RAM	RAM type	CPU	CPU CLOCK	MotherBoard	HDDs

// 16384 Mb	16384 Mb	2xDDR3-12800 1600MHz 8Gb(8192Mb)	Core i7-3770	3400 Mhz	Intel DH77EB	240 Gb
const linux = "Linux: i7-3770, 16Gb (12800 1600MHz), SSD"

// 16384 Mb	16384 Mb	2xDDR3-10600 1333MHz 8Gb(8192Mb)	Core i7-3770	3400 Mhz	Intel DH77EB	240 Gb
const linux2 = "Linux: i7-3770, 16Gb (10600 1333MHz), SSD"

// 16384 Mb	32768 Mb	4xDDR3-12800 1600MHz 8Gb(8192Mb)	Core i7-3770	3400 Mhz	Intel DH77EB	240 Gb
const linux32gb = "Linux: i7-3770, 32Gb (12800 1600MHz), SSD"

func (t *StatsServer) handleInfoRequest(request *http.Request) ([]byte, error) {
  var productNames []string
  groupNames := []string{macMini, linux, linux2, linux32gb, win}
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

    groupName, ok := t.machineToGroupName[machine]
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
