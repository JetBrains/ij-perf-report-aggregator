package server

import (
  "context"
  "errors"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/analyzer"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/util"
  "net/http"
)

func (t *StatsServer) handleInfoRequest(request *http.Request) ([]byte, error) {
  groupNames := t.machineInfo.GroupNames
  productNameToMachineMap, productNames, err := t.computeProductToMachines(request.Context())
  if err != nil {
    return nil, err
  }

  productToProjects, _, err := t.computeProductToProjects(request.Context())
  if err != nil {
    return nil, err
  }

  buffer := byteBufferPool.Get()
  defer byteBufferPool.Put(buffer)
  WriteInfo(buffer, productNames, groupNames, analyzer.MetricDescriptors, productNameToMachineMap, productToProjects)
  return CopyBuffer(buffer), nil
}

func (t *StatsServer) computeProductToMachines(taskContext context.Context) (map[string]map[string]*MachineGroup, []string, error) {
  var productNames []string
  rows, err := t.db.QueryContext(taskContext, "select distinct product, machine from report where machine != 'Dead agent' group by product, machine order by product, machine")
  if err != nil {
    return nil, nil, err
  }

  defer util.Close(rows, t.logger)

  productNameToMachineMap := make(map[string]map[string]*MachineGroup)
  for rows.Next() {
    var product string
    var machine string
    err := rows.Scan(&product, &machine)
    if err != nil {
      return nil, nil, err
    }

    groupName, ok := t.machineInfo.MachineToGroupName[machine]
    if !ok {
      return nil, nil, errors.New("Group is unknown machine: " + machine)
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
  return productNameToMachineMap, productNames, nil
}

func (t *StatsServer) computeProductToProjects(taskContext context.Context) (map[string]*[]string, []string, error) {
  var productNames []string
  rows, err := t.db.QueryContext(taskContext, "select distinct product, project from report where machine != 'Dead agent' group by product, project order by product, project")
  if err != nil {
    return nil, nil, err
  }

  defer util.Close(rows, t.logger)

  productNameToProjectsMap := make(map[string]*[]string)
  for rows.Next() {
    var product string
    var project string
    err := rows.Scan(&product, &project)
    if err != nil {
      return nil, nil, err
    }

    groupToProjects, ok := productNameToProjectsMap[product]
    if ok {
      *groupToProjects = append(*groupToProjects, project)
    } else {
      productNames = append(productNames, product)
      productNameToProjectsMap[product] = &[]string{project}
    }
  }
  return productNameToProjectsMap, productNames, nil
}

type MachineGroup struct {
  name     string
  machines []string
}
