package server

import (
  "github.com/asaskevich/govalidator"
  "github.com/develar/errors"
  "github.com/jmoiron/sqlx"
  "github.com/json-iterator/go"
  "github.com/olekukonko/tablewriter"
  "io"
  "log"
  "net/http"
  "sort"
  "strconv"
)

type Item struct {
  Product   string
  Bootstrap float64

  MachineGroup string
}

type ActualData struct {
  Item

  InstallerId int `db:"tc_installer_build_id"`
  RunCount    int `db:"runCount"`
}

type AnalyzeResult struct {
  GoldItems    []*Item
  CurrentItems map[string]map[string][]ActualData
}

func (t *StatsServer) handleStatusRequest(request *http.Request) ([]byte, error) {
  urlQuery := request.URL.Query()

  branch := urlQuery.Get("branch")
  switch {
  case len(branch) == 0:
    branch = "master"
  case !govalidator.IsAlphanumeric(branch):
    return nil, NewHttpError(400, "The branch parameter must be alphanumeric")
  }

  goldWeekStart := urlQuery.Get("goldWeekStart")
  switch {
  case len(goldWeekStart) == 0:
    return nil, NewHttpError(400, "The goldWeekStart parameter is required")
  case !govalidator.IsTime(goldWeekStart, "2006-01-02"):
    return nil, NewHttpError(400, "The goldWeekStart parameter must be in format yyyy-mm-dd hh:ss:mm")
  }

  result, err := CompareMetrics(t.db, branch, goldWeekStart + " 00:00:00")
  if err != nil {
    return nil, err
  }

  if request.Header.Get("Accept") == "application/json" {
    bytes, err := jsoniter.ConfigFastest.Marshal(result)
    if err != nil {
      return nil, err
    }

    return bytes, nil
  }

  buffer := byteBufferPool.Get()
  defer byteBufferPool.Put(buffer)
  PrintResult(*result, buffer)
  return CopyBuffer(buffer), nil
}

func PrintResult(result AnalyzeResult, writer io.Writer) {
  table := tablewriter.NewWriter(writer)
  table.SetHeader([]string{"Product", "Machine", "Gold", "Actual"})
  table.SetAutoMergeCells(true)
  table.SetRowLine(true)
  table.SetAutoWrapText(false)

  for _, item := range result.GoldItems {
    actualItems, ok := result.CurrentItems[item.MachineGroup]
    if !ok {
      log.Printf("error: no actual data for product %s and machine %s", item.Product, item.MachineGroup)
      continue
    }

    actualItem := actualItems[item.Product]
    if len(actualItem) == 0 {
      log.Printf("error: no actual data for product %s and machine %s", item.Product, item.MachineGroup)
      continue
    }

    table.Append([]string{item.Product, item.MachineGroup, formatFloat(item.Bootstrap), formatFloat(actualItem[0].Bootstrap)})
  }

  table.Render()
}

func CompareMetrics(db *sqlx.DB, branch string, goldWeekStart string) (*AnalyzeResult, error) {
  goldItems, err := analyzeGold(db, branch, goldWeekStart)
  if err != nil {
    return nil, err
  }

  currentItems, err := analyzeCurrent(db, branch)
  if err != nil {
    return nil, err
  }

  return &AnalyzeResult{
    GoldItems:    goldItems,
    CurrentItems: currentItems,
  }, nil
}

func formatFloat(v float64) string {
  return strconv.FormatFloat(v, 'f', 1, 64)
}

func analyzeGold(db *sqlx.DB, branch string, goldWeekStart string) ([]*Item, error) {
  totalResult := make([]*Item, 0, 4)

  // We compute aggregated data for machine group, but db contains only machine and not group. So, have to execute request for each machine group separately.
  machineInfo := GetMachineInfo()
  for _, groupName := range machineInfo.GroupNames {
    query, args, err := sqlx.In(`
      select product, quantileTDigest(bootstrap_d) as bootstrap 
      from report 
      where branch = ? and generated_time between toDateTime(?, 'UTC') and addWeeks(toDateTime(?, 'UTC'), 1) and machine in (?)
      group by product
    `, branch, goldWeekStart, goldWeekStart, getMachineNames(machineInfo, groupName))
    if err != nil {
      return nil, errors.WithStack(err)
    }

    var result []*Item
    err = db.Select(&result, query, args...)
    if err != nil {
      return nil, errors.WithStack(err)
    }

    if len(result) == 0 {
      continue
    }

    for _, item := range result {
      item.MachineGroup = groupName
    }

    totalResult = append(totalResult, result...)
  }

  sort.Slice(totalResult, func(i, j int) bool {
    return totalResult[i].Product < totalResult[j].Product
  })

  return totalResult, nil
}

func analyzeCurrent(db *sqlx.DB, branch string) (map[string]map[string][]ActualData, error) {
  machineInfo := GetMachineInfo()

  totalResult := make(map[string]map[string][]ActualData, len(machineInfo.GroupNames))

  // We compute aggregated data for machine group, but db contains only machine and not group. So, have to execute request for each machine group separately.
  for _, machineGroupName := range machineInfo.GroupNames {
    // do not limit, because there are multiple products
    query, args, err := sqlx.In(`
      select product, tc_installer_build_id, count(tc_installer_build_id) as runCount, quantileTDigest(bootstrap_d) as bootstrap from report
      where tc_installer_build_id != 0 and branch = ? and generated_time > subtractDays(now(), 1) and machine in (?)
      group by product, tc_installer_build_id
      order by product, runCount desc, tc_installer_build_id desc
    `, branch, getMachineNames(machineInfo, machineGroupName))
    if err != nil {
      return nil, errors.WithStack(err)
    }

    var result []ActualData
    err = db.Select(&result, query, args...)
    if err != nil {
      return nil, errors.WithStack(err)
    }

    if len(result) == 0 {
      continue
    }

    for _, item := range result {
      item.MachineGroup = machineGroupName
    }

    productMap := totalResult[machineGroupName]
    if productMap == nil {
      productMap = make(map[string][]ActualData)
      totalResult[machineGroupName] = productMap
    }

    for _, item := range result {
      list := productMap[item.Product]
      if list == nil {
        productMap[item.Product] = []ActualData{item}
      } else {
        productMap[item.Product] = append(list, item)
      }
    }
  }

  return totalResult, nil
}

func getMachineNames(machineInfo MachineInfo, groupName string) []string {
  var machineNames []string
  for name, otherGroupName := range machineInfo.MachineToGroupName {
    if groupName == otherGroupName {
      machineNames = append(machineNames, name)
    }
  }
  return machineNames
}
