package main

import (
  "context"
  "database/sql"
  "fmt"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/analyzer"
  "github.com/develar/errors"
  "golang.org/x/tools/container/intsets"
  "strings"
  "time"
)

type ReportExistenceChecker struct {
  ids intsets.Sparse
}

func (t *ReportExistenceChecker) reset(dbName string, buildTypeId string, reportAnalyzer *analyzer.ReportAnalyzer, taskContext context.Context, since time.Time) error {
  t.ids.Clear()

  var rows *sql.Rows
  var err error
  if dbName == "ij" {
    var product string
    for code, name := range productCodeToBuildName {
      if strings.Contains(buildTypeId, "_"+name) {
        product = code
        break
      }
    }
    if len(product) == 0 {
      return errors.New("cannot infer product from " + buildTypeId)
    }

    // don't filter by machine - product is enough to reduce set
    rows, err = reportAnalyzer.Db.QueryContext(taskContext, "select tc_build_id from report where product = ? and generated_time > ? order by tc_build_id", product, since)
  } else {
    table := "report"
    if reportAnalyzer.InsertReportManager.TableName != "" {
      table = reportAnalyzer.InsertReportManager.TableName
    }
    query := fmt.Sprintf("select tc_build_id from %s where generated_time > ? order by tc_build_id", table)
    rows, err = reportAnalyzer.Db.QueryContext(taskContext, query, since)
  }

  if err != nil {
    return errors.WithStack(err)
  }

  for rows.Next() {
    var id int
    err = rows.Scan(&id)
    if err != nil {
      return errors.WithStack(err)
    }

    t.ids.Insert(id)
  }

  if rows.Err() != nil {
    return errors.WithStack(rows.Err())
  }
  return nil
}

func (t *ReportExistenceChecker) has(id int) bool {
  return t.ids.Has(id)
}
