package server

import (
  "encoding/json"
  "fmt"
  "github.com/ClickHouse/clickhouse-go/v2"
  "github.com/ClickHouse/clickhouse-go/v2/lib/driver"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/util"
  "github.com/valyala/bytebufferpool"
  "github.com/valyala/quicktemplate"
  "net/http"
  "strings"
)

func (t *StatsServer) getBranchComparison(request *http.Request) (*bytebufferpool.ByteBuffer, bool, error) {

  type requestParams struct {
    Table        string   `json:"table"`
    MeasureNames []string `json:"measure_names"`
    Branch       string   `json:"branch"`
    Machine      string   `json:"machine"`
  }

  var params requestParams
  data, err := util.DecodeQuery(request.URL.Path[len("/api/compareBranches/"):])
  if err != nil {
    return nil, false, err
  }
  err = json.Unmarshal(data, &params)
  if err != nil {
    return nil, false, err
  }

  table := params.Table
  measureNames := params.MeasureNames
  branch := params.Branch
  quotedMeasureNames := make([]string, len(measureNames))
  for i, name := range measureNames {
    quotedMeasureNames[i] = "'" + name + "'"
  }
  measureNamesString := strings.Join(quotedMeasureNames, ",")
  machine := params.Machine

  sql := fmt.Sprintf("SELECT project as Project, measure_name as MeasureName, arraySlice(groupArray(measure_value), 1, 50) AS MeasureValues FROM (SELECT project, measures.name as measure_name, measures.value as measure_value FROM %s ARRAY JOIN measures WHERE branch = '%s' AND measure_name in (%s) AND machine like '%s' ORDER BY generated_time DESC)GROUP BY project, measure_name;", table, branch, measureNamesString, machine)
  db, err := clickhouse.Open(&clickhouse.Options{
    Addr: []string{t.dbUrl},
    Auth: clickhouse.Auth{
      Database: "ij",
    },
    Settings: map[string]interface{}{
      "readonly":         1,
      "max_query_size":   1000000,
      "max_memory_usage": 3221225472,
    },
  })
  defer func(db driver.Conn) {
    _ = db.Close()
  }(db)
  if err != nil {
    return nil, false, err
  }

  var queryResults []struct {
    Project       string
    MeasureName   string
    MeasureValues []int
  }

  err = db.Select(request.Context(), &queryResults, sql)
  if err != nil {
    return nil, false, err
  }

  type responseItem struct {
    Project     string
    MeasureName string
    Median      float64
  }

  response := make([]responseItem, len(queryResults))

  for i, result := range queryResults {
    indexes, err := GetChangePointIndexes(result.MeasureValues, 1)
    if err != nil {
      return nil, false, err
    }

    var valuesAfterLastChangePoint []int
    if len(indexes) == 0 {
      valuesAfterLastChangePoint = result.MeasureValues
    } else {
      lastIndex := indexes[len(indexes)-1]
      valuesAfterLastChangePoint = result.MeasureValues[lastIndex:]
    }
    median := CalculateMedian(valuesAfterLastChangePoint)
    response[i] = responseItem{
      Project:     result.Project,
      MeasureName: result.MeasureName,
      Median:      median,
    }
  }

  jsonData, err := json.Marshal(response)
  if err != nil {
    return nil, false, err
  }

  buffer := byteBufferPool.Get()
  _, err = buffer.Write(jsonData)
  if err != nil {
    return nil, false, err
  }

  return buffer, true, err
}

func (t *StatsServer) getDistinctHighlightingPasses(request *http.Request) (*bytebufferpool.ByteBuffer, bool, error) {
  sql := "SELECT DISTINCT arrayJoin((arrayFilter(x-> x LIKE 'highlighting/%', `metrics.name`))) as PassName from report where generated_time >subtractMonths(now(),12)"
  db, err := clickhouse.Open(&clickhouse.Options{
    Addr: []string{t.dbUrl},
    Auth: clickhouse.Auth{
      Database: "ij",
    },
    Settings: map[string]interface{}{
      "readonly":         1,
      "max_query_size":   1000000,
      "max_memory_usage": 3221225472,
    },
  })
  if err != nil {
    return nil, false, err
  }
  defer func(db driver.Conn) {
    _ = db.Close()
  }(db)

  var result []struct {
    PassName string
  }
  err = db.Select(request.Context(), &result, sql)
  if err != nil {
    return nil, false, err
  }

  buffer := byteBufferPool.Get()
  templateWriter := quicktemplate.AcquireWriter(buffer)
  defer quicktemplate.ReleaseWriter(templateWriter)
  jsonWriter := templateWriter.N()
  jsonWriter.S("[")
  for i, v := range result {
    if i != 0 {
      jsonWriter.S(",")
    }
    jsonWriter.Q(v.PassName)
  }
  jsonWriter.S("]")

  return buffer, true, err
}
