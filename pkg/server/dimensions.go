package server

import (
  "context"
  "database/sql"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/analyzer"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/data-query"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/util"
  "github.com/jmoiron/sqlx"
  "net/http"
)

func (t *StatsServer) handleDimensionsRequest(request *http.Request) ([]byte, error) {
  dbName, err := data_query.ValidateDatabaseName(request.URL.Query().Get("db"))
  if err != nil {
    return nil, err
  }

  db, err := t.GetDatabase(dbName)
  if err != nil {
    return nil, err
  }

  productNameToMachineMap, productNames, err := t.computeProductToMachines(dbName, db, request.Context())
  if err != nil {
    return nil, err
  }

  productToProjects, _, err := t.computeProducts(dbName, db, request.Context())
  if err != nil {
    return nil, err
  }

  buffer := byteBufferPool.Get()
  defer byteBufferPool.Put(buffer)

  var metricDescriptors []*analyzer.Metric
  if dbName == "ij" {
    metricDescriptors = analyzer.MetricDescriptors
  } else {
    metrics, err := t.computeAvailableMetrics(dbName, db, request.Context())
    if err != nil {
      return nil, err
    }
    metricDescriptors = metrics
  }
  WriteInfo(buffer, productNames, t.machineInfo.GroupNames, metricDescriptors, productNameToMachineMap, productToProjects)
  return CopyBuffer(buffer), nil
}

func (t *StatsServer) computeProducts(dbName string, db *sqlx.DB, taskContext context.Context) (map[string]*[]string, []string, error) {
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