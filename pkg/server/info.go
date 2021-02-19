package server

import (
  "context"
  "database/sql"
  "errors"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/analyzer"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/data-query"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/util"
  "github.com/jmoiron/sqlx"
  "github.com/valyala/quicktemplate"
  "net/http"
  "strings"
)

func (t *StatsServer) handleMetaMeasureRequest(request *http.Request) ([]byte, error) {
  dbName, err := data_query.ValidateDatabaseName(request.URL.Query().Get("db"))
  if err != nil {
    return nil, err
  }

  db, err := t.GetDatabase(dbName)
  if err != nil {
    return nil, err
  }

  buffer := byteBufferPool.Get()
  defer byteBufferPool.Put(buffer)

  var measureNames []string
  if dbName == "ij" {
    measureNames = make([]string, len(analyzer.IjMetricDescriptors))
    for index, descriptor := range analyzer.IjMetricDescriptors {
      measureNames[index] = descriptor.Name
    }
  } else {
    measureNames, err = t.computeAvailableMetrics(dbName, db, request.Context())
    if err != nil {
      return nil, err
    }
  }

  templateWriter := quicktemplate.AcquireWriter(buffer)
  defer quicktemplate.ReleaseWriter(templateWriter)
  jsonWriter := templateWriter.N()
  for index, name := range measureNames {
    if index != 0 {
      jsonWriter.S(",")
    }
    jsonWriter.Q(name)
  }
  jsonWriter.S("]")

  return CopyBuffer(buffer), nil
}

func (t *StatsServer) computeAvailableMetrics(dbName string, db *sqlx.DB, taskContext context.Context) ([]string, error) {
	if dbName == "fleet" {
		return []string{}, nil
	}

	var result []string
  //goland:noinspection SqlResolve
	err := db.SelectContext(taskContext, &result, "select distinct measures.name from report array join measures")
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (t *StatsServer) computeProductToMachines(dbName string, db *sqlx.DB, taskContext context.Context) (map[string]map[string]*MachineGroup, []string, error) {
  var productNames []string
  var rows *sql.Rows
  var err error
  hasProductField := dbName == "ij"
  if hasProductField {
    rows, err = db.QueryContext(taskContext, "select distinct product, machine from report group by product, machine order by product, machine")
  } else {
    rows, err = db.QueryContext(taskContext, "select distinct machine from report order by machine")
  }
  if err != nil {
    return nil, nil, err
  }

  defer util.Close(rows, t.logger)

  productNameToMachineMap := make(map[string]map[string]*MachineGroup)
  for rows.Next() {
    var product string
    var machine string
    if hasProductField {
      err = rows.Scan(&product, &machine)
    } else {
      product = dbName
      err = rows.Scan(&machine)
    }
    if err != nil {
      return nil, nil, err
    }

    var groupName string
    if strings.HasPrefix(machine, "intellij-linux-hw-blade-") {
      groupName = "linux-blade"
    } else if strings.HasPrefix(machine, "intellij-windows-hw-blade-") {
      groupName = "windows-blade"
    } else {
      var ok bool
      groupName, ok = t.machineInfo.MachineToGroupName[machine]
      if !ok {
        return nil, nil, errors.New("Group is unknown machine: " + machine)
      }
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

func (t *StatsServer) computeProductToProjects(dbName string, db *sqlx.DB, taskContext context.Context) (map[string]*[]string, []string, error) {
  var productNames []string
  var rows *sql.Rows
  var err error
  hasProductField := dbName == "ij"
  if hasProductField {
    rows, err = db.QueryContext(taskContext, "select distinct product, project from report group by product, project order by product, project")
  } else {
    rows, err = db.QueryContext(taskContext, "select distinct project from report order by project")
  }
  if err != nil {
    return nil, nil, err
  }

  defer util.Close(rows, t.logger)

  productNameToProjectsMap := make(map[string]*[]string)
  for rows.Next() {
    var product string
    var project string
    if hasProductField {
      err = rows.Scan(&product, &project)
    } else {
      product = dbName
      err = rows.Scan(&project)
    }
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
