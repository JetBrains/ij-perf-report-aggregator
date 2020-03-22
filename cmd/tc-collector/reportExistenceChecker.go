package main

import (
  "context"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/analyzer"
  "github.com/develar/errors"
  "golang.org/x/tools/container/intsets"
  "strings"
  "time"
)

type ReportExistenceChecker struct {
  ids intsets.Sparse
}

func (t *ReportExistenceChecker) reset(buildTypeId string, reportAnalyzer *analyzer.ReportAnalyzer, taskContext context.Context, since time.Time) error {
  t.ids.Clear()

  var product string
  for code, name := range ProductCodeToBuildName {
    if strings.Contains(buildTypeId, "_"+name) {
      product = code
      break
    }
  }
  if len(product) == 0 {
    return errors.New("cannot infer product from " + buildTypeId)
  }

  // don't filter by machine - product is enough to reduce set
  rows, err := reportAnalyzer.Db.QueryContext(taskContext, "select tc_build_id from report where product = ? and generated_time > ? order by tc_build_id", product, since)
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