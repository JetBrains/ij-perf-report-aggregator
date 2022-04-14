package server

import (
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/analyzer"
  "github.com/develar/errors"
  "github.com/jmoiron/sqlx"
  "github.com/olekukonko/tablewriter"
  "io"
  "log"
  "sort"
  "strconv"
  "strings"
)

//noinspection SpellCheckingInspection
var projectNameToTitle = map[string]string{
  "/q9N7EHxr8F1NHjbNQnpqb0Q0fs": "joda-time",
  "73YWaW9bytiPDGuKvwNIYMK5CKI": "simple for IJ",
  "j1a8nhKJexyL/zyuOXJ5CFOHYzU": "simple for PS",
  "JeNLJFVa04IA+Wasc+Hjj3z64R0": "simple for WS",
  "nC4MRRFMVYUSQLNIvPgDt+B3JqA": "Idea",
}

type Item struct {
  Product string
  Project string

  InstallerId int `db:"tc_installer_build_id"`

  Bootstrap float64
  AppInit   float64 `db:"appInit"`

  MachineGroup string
}

type ActualData struct {
  Item

  BuildTime int `db:"build_time"`
  RunCount  int `db:"runCount"`
}

type AnalyzeResult struct {
  GoldItems    []*Item
  CurrentItems map[string]map[string][]ActualData
}

//func (t *StatsServer) handleStatusRequest(request *http.Request) ([]byte, error) {
//  urlQuery := request.URL.Query()
//
//  branch := urlQuery.Get("branch")
//  switch {
//  case len(branch) == 0:
//    branch = "master"
//  case !govalidator.IsAlphanumeric(branch):
//    return nil, http_error.NewHttpError(400, "The branch parameter must be alphanumeric")
//  }
//
//  goldWeekStart := urlQuery.Get("goldWeekStart")
//  switch {
//  case len(goldWeekStart) == 0:
//    return nil, http_error.NewHttpError(400, "The goldWeekStart parameter is required")
//  case !govalidator.IsTime(goldWeekStart, "2006-01-02"):
//    return nil, http_error.NewHttpError(400, "The goldWeekStart parameter must be in format yyyy-mm-dd hh:ss:mm")
//  }
//
//  db, err := t.AcquireDatabase("ij", context.Background())
//  if err != nil {
//    return nil, err
//  }
//
//  result, err := CompareMetrics(db, branch, goldWeekStart+" 00:00:00")
//  if err != nil {
//    return nil, err
//  }
//
//  if request.Header.Get("Accept") == "application/json" {
//    bytes, err := jsoniter.ConfigFastest.Marshal(result)
//    if err != nil {
//      return nil, err
//    }
//
//    return bytes, nil
//  }
//
//  buffer := byteBufferPool.Get()
//  defer byteBufferPool.Put(buffer)
//  PrintResult(*result, buffer)
//  return CopyBuffer(buffer), nil
//}

func PrintResult(result AnalyzeResult, writer io.Writer) {
  table := tablewriter.NewWriter(writer)
  table.SetHeader([]string{"Product", "Project", "Machine", "Bootstrap (g)", "Bootstrap (a)", "App Init (g)", "App Init (a)", "Installer Id"})
  table.SetAutoMergeCells(true)
  table.SetRowLine(true)
  table.SetAutoWrapText(false)

  for _, goldItem := range result.GoldItems {
    actualItems, ok := result.CurrentItems[goldItem.MachineGroup]
    if !ok {
      log.Printf("error: no actual data for product %s and machine %s", goldItem.Product, goldItem.MachineGroup)
      continue
    }

    actualItem := actualItems[getItemKey(*goldItem)]
    if len(actualItem) == 0 {
      log.Printf("error: no actual data for product %s and machine %s", goldItem.Product, goldItem.MachineGroup)
      continue
    }

    projectLabel := projectNameToTitle[goldItem.Project]
    if len(projectLabel) == 0 {
      projectLabel = goldItem.Project
    }
    actualData := actualItem[0]
    table.Append([]string{goldItem.Product, projectLabel, goldItem.MachineGroup,
      formatFloat(goldItem.Bootstrap), formatFloat(actualData.Bootstrap),
      formatFloat(goldItem.AppInit), formatFloat(actualData.AppInit),
      strconv.Itoa(actualData.InstallerId),
      //time.Unix(int64(actualData.BuildTime), 0).String(),
    })
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
  machineInfo := analyzer.GetMachineInfo()
  for _, groupName := range machineInfo.GroupNames {
    query, args, err := sqlx.In(`
      select product, project, quantileTDigest(bootstrap_d) as bootstrap, quantileTDigest(appInit_d) as appInit 
      from report 
      where branch = ? and generated_time between toDateTime(?, 'UTC') and addWeeks(toDateTime(?, 'UTC'), 1) and machine in (?)
      group by product, project
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
    a := totalResult[i]
    b := totalResult[j]
    if a.Product != b.Product {
      return a.Product < b.Product
    }
    if a.Project != b.Project {
      return a.Project < b.Project
    }

    aM := machineGroupOrder(a.MachineGroup)
    bM := machineGroupOrder(b.MachineGroup)
    if aM == bM {
      return a.MachineGroup < b.MachineGroup
    } else {
      return aM < bM
    }
  })

  return totalResult, nil
}

// wanted order: mac, linux, windows
func machineGroupOrder(s string) int {
  switch {
  case strings.HasPrefix(s, "mac"):
    return 1

  case strings.HasPrefix(s, "Linux"):
    return 2

  default:
    return 3
  }
}

func analyzeCurrent(db *sqlx.DB, branch string) (map[string]map[string][]ActualData, error) {
  machineInfo := analyzer.GetMachineInfo()

  totalResult := make(map[string]map[string][]ActualData, len(machineInfo.GroupNames))

  // We compute aggregated data for machine group, but db contains only machine and not group. So, have to execute request for each machine group separately.
  for _, machineGroupName := range machineInfo.GroupNames {
    // do not limit, because there are multiple products
    // for last 3 days — to ensure that we can find installer that was enough tested (at least 4 times).
    query, args, err := sqlx.In(`
      select * 
      from (
            select product, project, tc_installer_build_id, count(tc_installer_build_id) as runCount, quantileTDigest(bootstrap_d) as bootstrap, quantileTDigest(appInit_d) as appInit
            from report
            where tc_installer_build_id != 0 and branch = ? and generated_time > subtractDays(now(), 3) and machine in (?)
            group by product, project, tc_installer_build_id
           )
      where runCount >= 4
      order by product, project, tc_installer_build_id desc
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

    itemMap := totalResult[machineGroupName]
    if itemMap == nil {
      itemMap = make(map[string][]ActualData)
      totalResult[machineGroupName] = itemMap
    }

    for _, item := range result {
      if item.RunCount <= 3 {
        // skip
        continue
      }

      key := getItemKey(item.Item)
      list := itemMap[key]
      if list == nil {
        itemMap[key] = []ActualData{item}
      } else {
        itemMap[key] = append(list, item)
      }
    }
  }

  return totalResult, nil
}

func getItemKey(item Item) string {
  return item.Product + ":" + item.Project
}

func getMachineNames(machineInfo analyzer.MachineInfo, groupName string) []string {
  var machineNames []string
  for name, otherGroupName := range machineInfo.MachineToGroupName {
    if groupName == otherGroupName {
      machineNames = append(machineNames, name)
    }
  }
  return machineNames
}
